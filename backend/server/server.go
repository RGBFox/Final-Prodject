package server

import (
	"RGBFox/backend/database"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

// Server хранит общий пул соединений, используемый HTTP-обработчиками.
type Server struct {
	db *sql.DB
}

// Init регистрирует маршруты в http.DefaultServeMux и возвращает порт сервера.
// Переданный db должен оставаться открытым всё время работы обработчиков.
func Init(db *sql.DB) string {
	s := &Server{db: db}
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/signin", s.signinHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", s.auth(s.taskHandler))
	http.HandleFunc("/api/tasks", s.auth(s.tasksHandler))
	http.HandleFunc("/api/task/done", s.auth(s.doneHandler))

	return port
}

// Task описывает задачу в БД и JSON-ответах API.
// Date хранится в формате 20060102, Repeat содержит правило повторения.
type Task = database.Task

// taskHandler выбирает операцию над одной задачей по HTTP-методу.
func (s *Server) taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.addTaskHandler(w, r)
	case http.MethodGet:
		s.getTaskHandler(w, r)
	case http.MethodDelete:
		s.deleteTaskHandler(w, r)
	case http.MethodPut:
		s.updateTaskHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, fmt.Errorf("метод не поддерживается"))
	}
}

// doneHandler распределяет запросы маршрута /api/task/done по HTTP-методам.
func (s *Server) doneHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.doneTaskHandler(w, r)
	case http.MethodPut:
		s.updateTaskHandler(w, r)
	case http.MethodDelete:
		s.deleteTaskHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, fmt.Errorf("метод не поддерживается"))
	}
}

// tasksHandler обрабатывает GET-запросы списка задач.
func (s *Server) tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getTasksHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, fmt.Errorf("метод не поддерживается"))
	}
}
