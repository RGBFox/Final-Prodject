package server

import (
	"RGBFox/backend/rules"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func Init() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)

	return port
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	now := time.Now()

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ошибка десериализации JSON"})
		return
	}
	if task.Date != "" {
		_, err := time.Parse("20060102", task.Date)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "ошибка формата времени"})
			return
		}
	} else {
		task.Date = now.Format("20060102")
	}
	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ошибка формата времени"})
		return
	}
	if now.After(date) {
		task.Date = now.Format("20060102")
	} else if strings.Content(task.repeat) {

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]any{"id": "42"})
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if !q.Has("now") {
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

	w.Write([]byte(next))

}
