package server

import (
	"RGBFox/backend/database"
	"fmt"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

// getTasksHandler возвращает до 50 задач по возрастанию даты с необязательным поиском.
func (s *Server) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := q.Get("search")
	tasks, res, err := database.SearchTask(s.db, search)
	if !checkDatabaseErrors(w, r, res, err) {
		return
	}

	type TasksResp struct {
		Tasks []Task `json:"tasks"`
	}
	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}

// getTaskHandler возвращает одну задачу по параметру id в URL.
func (s *Server) getTaskHandler(w http.ResponseWriter, r *http.Request) {
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
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("задача не найдена"))
		return
	}

	writeJSON(w, http.StatusOK, arg)
}
