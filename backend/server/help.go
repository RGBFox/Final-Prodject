package server

import (
	"RGBFox/backend/rules"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

// nextDateHandler вычисляет следующую дату по параметрам now, date и repeat.
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if !q.Has("now") || q.Get("now") == "" {
		http.Error(w, "нужен параметр now", http.StatusBadRequest)
		return
	}

	if !q.Has("date") {
		http.Error(w, "нужен параметр date", http.StatusBadRequest)
		return
	}

	if !q.Has("repeat") {
		http.Error(w, "нужен параметр repeat", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", q.Get("now"))
	if err != nil {
		http.Error(w, "неверный формат now: "+err.Error(), http.StatusBadRequest)
		return
	}
	date := q.Get("date")
	repeat := q.Get("repeat")

	next, err := rules.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "ошибка обработки: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Этот маршрут возвращает дату обычным текстом, без JSON-обёртки.
	w.Write([]byte(next))

}

// writeJSON отправляет JSON с заданным HTTP-статусом.
// Значение типа error преобразуется в объект с полем error.
func writeJSON(w http.ResponseWriter, status int, v any) {
	if err, ok := v.(error); ok {
		log.Printf("ошибка: %v\n", err)
		v = map[string]string{"error": err.Error()}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)

}

// prepareTask проверяет название, дату и повторение, изменяя task.Date при необходимости.
// Пустую дату заполняет сегодняшним днём; просроченную переносит с учётом повторения.
func prepareTask(task *Task, now time.Time) error {
	if task.Title == "" {
		return fmt.Errorf("поле Title не может быть пустым")
	}
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	today := now.Format("20060102")
	// Сравнение строк корректно после проверки формата 20060102.
	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date, err = rules.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
		}
	} else if task.Repeat != "" {
		// дату не меняем, только проверяем, что правило корректно
		if _, err := rules.NextDate(now, task.Date, task.Repeat); err != nil {
			return err
		}
	}
	return nil
}

func checkDatabaseErrors(w http.ResponseWriter, r *http.Request, typeerr int, err error) bool {
	if err != nil {
		if typeerr == http.StatusInternalServerError {
			log.Printf("ошибка обновления данных: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, fmt.Errorf("ошибка обновления данных на стороне сервера"))
			return false
		} else {
			writeJSON(w, http.StatusBadRequest, fmt.Errorf("ошибка обновления данных"))
			return false
		}
	}
	return true

}
