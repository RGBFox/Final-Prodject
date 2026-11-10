package main

import (
	"RGBFox/backend/database"
	"RGBFox/backend/server"
	"database/sql"
	"log"
	"net/http"
)

// main подготавливает БД, регистрирует маршруты и запускает HTTP-сервер.
func main() {
	err, path := database.Init()
	if err != nil || path == "" {
		log.Fatal(err)
	}
	// Этот пул используется всеми обработчиками в течение работы сервера.
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	port := server.Init(db)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
