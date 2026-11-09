package server

import (
	"RGBFox/backend/rules"
	"encoding/json"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

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

func writeJSON(w http.ResponseWriter, status int, v any) {
	if err, ok := v.(error); ok {
		v = map[string]string{"error": err.Error()}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)

}

func search() {

}
