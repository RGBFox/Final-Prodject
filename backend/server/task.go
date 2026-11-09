package server

import (
	"RGBFox/backend/rules"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	now := time.Now()

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("поле Title не может быть пустым"))
		return
	}
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}
	today := now.Format("20060102")

	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date, err = rules.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, err)
				return

			}
		}
	} else if task.Repeat != "" {
		// дату не меняем, только проверяем, что правило корректно
		if _, err := rules.NextDate(now, task.Date, task.Repeat); err != nil {
			writeJSON(w, http.StatusBadRequest, err)
			return

		}
	}

	id, err := AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func AddTask(task *Task) (int64, error) {
	var id int64
	dbFile := "scheduler.db"
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := `INSERT INTO scheduler (date, title, comment, repeat)
	          VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}
