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

// updateTaskHandler проверяет данные и обновляет задачу с указанным ID.
func (s *Server) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	// Числовой ID ещё не подтверждает наличие задачи: это проверяется после UPDATE.
	if id, err := strconv.ParseInt(task.ID, 10, 64); err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("некорректный id"))
		return
	}

	err = s.updateTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("ошибка обновления данных"))
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

// updateTask записывает поля задачи и возвращает sql.ErrNoRows, если её ID не найден.
func (s *Server) updateTask(task *Task) error {

	query := `UPDATE scheduler 
	SET date = :date,
	  title= :title,
	  comment= :comment,
	  repeat=  :repeat 
	WHERE id = :id;
	`
	res, err := s.db.Exec(query, sql.Named("id", task.ID), sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}
	// UPDATE без совпадений не считается ошибкой SQL, поэтому проверяем число строк.
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
