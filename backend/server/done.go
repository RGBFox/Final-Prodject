package server

import (
	"RGBFox/backend/rules"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {

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
	if arg.Repeat == "" {
		_, err := db.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", arg.ID))
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusOK, struct{}{})
		return
	}

	now := time.Now()
	arg.Date, err = rules.NextDate(now, arg.Date, arg.Repeat)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	_, err = updateTask(&arg)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

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
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("не указан id"))
		return
	}
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", search))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	n, err := res.RowsAffected()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	if n == 0 {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("задача не найдена"))
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}
