package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

// addTaskHandler проверяет данные задачи из JSON и возвращает ID новой записи.
func (s *Server) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	now := time.Now()
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("ошибка формата данных"))
		return
	}
	err := prepareTask(&task, now)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}
	id, err := s.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("ошибка Добавления задачи"))
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}

// AddTask сохраняет подготовленную задачу и возвращает ID, назначенный SQLite.
func (s *Server) AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat)
	          VALUES (:date, :title, :comment, :repeat)`
	res, err := s.db.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}
