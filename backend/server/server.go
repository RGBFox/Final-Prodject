package server

import (
	"RGBFox/backend/rules"
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

	return port
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
