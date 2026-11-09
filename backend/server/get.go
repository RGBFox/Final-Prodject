package server

import (
	"database/sql"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	dbFile := "scheduler.db"
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	q := r.URL.Query()
	search := q.Get("search")
	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	var args []any

	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			// поиск по дате
			query += ` WHERE date = :date`
			args = append(args, sql.Named("date", t.Format("20060102")))
		} else {
			// поиск по тексту
			query += ` WHERE title LIKE :pattern OR comment LIKE :pattern`
			args = append(args, sql.Named("pattern", "%"+search+"%"))
		}
	}
	query += ` ORDER BY date LIMIT 50`

	res, err := db.Query(query, args...)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	defer res.Close()
	tasks := []Task{}
	for res.Next() {
		var arg Task
		err = res.Scan(&arg.ID, &arg.Date, &arg.Title, &arg.Comment, &arg.Repeat)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, err)
			return
		}
		tasks = append(tasks, arg)
	}
	if err := res.Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	type TasksResp struct {
		Tasks []Task `json:"tasks"`
	}
	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
