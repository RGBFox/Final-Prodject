package server

import (
	"RGBFox/backend/database"
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

	res, err := database.UpdateTask(&task, s.db)
	if !checkDatabaseErrors(w, r, res, err) {
		return
	}

	writeJSON(w, res, struct{}{})
}
