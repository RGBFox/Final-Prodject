package server

import (
	"RGBFox/backend/database"
	"encoding/json"
	"fmt"
	"log"
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
	id, err := database.AddTask(&task, s.db)
	if err != nil {
		log.Printf("ошибка добавления задачи: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("ошибка добавления задачи"))
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}
