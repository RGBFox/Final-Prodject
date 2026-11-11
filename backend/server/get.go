package server

import (
	"RGBFox/backend/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

// getTasksHandler возвращает до 50 задач по возрастанию даты с необязательным поиском.
func (s *Server) getTasksHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	search := q.Get("search")
	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	var args []any

	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			// поиск по дате
			query += ` WHERE date = :date`
			// Дату из пользовательского формата приводим к формату хранения в SQLite.
			args = append(args, sql.Named("date", t.Format("20060102")))
		} else {
			// поиск по тексту
			query += ` WHERE title LIKE :pattern OR comment LIKE :pattern`
			// Знаки % вокруг параметра позволяют искать совпадение внутри текста.
			args = append(args, sql.Named("pattern", "%"+search+"%"))
		}
	}
	// Строки формата YYYYMMDD сортируются в порядке календарных дат.
	query += ` ORDER BY date LIMIT 50`

	res, err := s.db.Query(query, args...)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("Ошибка базы данных"))
		return
	}
	defer res.Close()
	// Инициализированный пустой срез кодируется в JSON как [], а не null.
	tasks := []Task{}
	for res.Next() {
		var arg Task
		err = res.Scan(&arg.ID, &arg.Date, &arg.Title, &arg.Comment, &arg.Repeat)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, fmt.Errorf("такой задачи не существует"))
			return
		}
		tasks = append(tasks, arg)
	}
	// Next может завершиться из-за ошибки чтения; проверяем её перед отправкой ответа.
	if err := res.Err(); err != nil {
		log.Printf("ошибки чтения: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("данные не полностью считались"))
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
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, arg)
}
