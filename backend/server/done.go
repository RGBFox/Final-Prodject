package server

import (
	"RGBFox/backend/rules"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

// doneTaskHandler удаляет выполненную разовую задачу
// или переносит повторяющуюся задачу на следующую дату.
func (s *Server) doneTaskHandler(w http.ResponseWriter, r *http.Request) {

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
	res := s.db.QueryRow(query, args...)
	var arg Task
	err := res.Scan(&arg.ID, &arg.Date, &arg.Title, &arg.Comment, &arg.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("такой задачи не существует"))
		return
	}
	if arg.Repeat == "" {
		res, err := s.db.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", arg.ID))
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, fmt.Errorf("ошибка удаления задачи"))
			return
		}
		// Успешный DELETE может не удалить ни одной строки, если такого ID нет.
		n, err := res.RowsAffected()
		if err != nil {
			writeJSON(w, http.StatusBadRequest, fmt.Errorf("такой задачи не существует"))
			return
		}
		if n == 0 {
			writeJSON(w, http.StatusBadRequest, fmt.Errorf("задача не найдена"))
			return
		}

		writeJSON(w, http.StatusOK, struct{}{})
		return
	}

	now := time.Now()
	// Новую дату считаем от сохранённой даты задачи с учётом сегодняшнего дня.
	arg.Date, err = rules.NextDate(now, arg.Date, arg.Repeat)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("ошибка поиска даты"))
		return
	}
	err = s.updateTask(&arg)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("Ошибка обновления задачи"))
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}

// deleteTaskHandler удаляет задачу по ID из URL и сообщает, если запись не найдена.
func (s *Server) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	search := q.Get("id")
	if search == "" {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("не указан id"))
		return
	}
	res, err := s.db.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", search))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("Ошибка удаления задачи"))
		return
	}
	n, err := res.RowsAffected()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("Ошибка поиска задачи"))
		return
	}
	if n == 0 {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("задача не найдена"))
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}
