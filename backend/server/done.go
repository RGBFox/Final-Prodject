package server

import (
	"RGBFox/backend/database"
	"RGBFox/backend/rules"
	"fmt"
	"net/http"
	"strconv"
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
	id, err := strconv.Atoi(search)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("неверный формат id"))
		return
	}
	arg, err := database.SearchByIDTask(s.db, id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	if arg.Repeat == "" {
		res, err := database.DeleteTask(s.db, id)
		if !checkDatabaseErrors(w, r, res, err) {
			return
		}
		writeJSON(w, res, struct{}{})
		return
	}

	now := time.Now()
	// Новую дату считаем от сохранённой даты задачи с учётом сегодняшнего дня.
	arg.Date, err = rules.NextDate(now, arg.Date, arg.Repeat)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("ошибка поиска даты"))
		return
	}
	res, err := database.UpdateTask(&arg, s.db)
	if !checkDatabaseErrors(w, r, res, err) {
		return
	}
	writeJSON(w, res, struct{}{})
}

// deleteTaskHandler удаляет задачу по ID из URL и сообщает, если запись не найдена.
func (s *Server) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	search := q.Get("id")
	if search == "" {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("не указан id"))
		return
	}
	id, err := strconv.Atoi(search)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("неверный формат id"))
		return
	}
	res, err := database.DeleteTask(s.db, id)
	if !checkDatabaseErrors(w, r, res, err) {
		return
	}
	writeJSON(w, res, struct{}{})
}
