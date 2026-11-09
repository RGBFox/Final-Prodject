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

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	dbFile := "scheduler.db"
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	q := r.URL.Query()
	search := q.Get("id")
	if search == "" {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("некоректное поле"))
		return
	}
	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	var args []any

	if search != "" {
		query += ` WHERE id = :id`
		args = append(args, sql.Named("id", search))
	}
	res := db.QueryRow(query, args...)
	var arg Task
	err = res.Scan(&arg.ID, &arg.Date, &arg.Title, &arg.Comment, &arg.Repeat)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, arg)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("поле Title не может быть пустым"))
		return
	}
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest,
			fmt.Errorf("некорректный id"))
		return
	}

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	_, err = time.Parse("20060102", task.Date)
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

	id, err = updateTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func updateTask(task *Task) (int64, error) {
	var id int64
	dbFile := "scheduler.db"
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := `UPDATE scheduler 
	SET date = :date,
	  title= :title,
	  comment= :comment,
	  repeat=  :repeat 
	WHERE id = :id;
	`
	res, err := db.Exec(query, sql.Named("id", task.ID), sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, fmt.Errorf("error")
	}
	return id, err
}
