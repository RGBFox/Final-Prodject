package main

import (
	"RGBFox/backend/database"
	"RGBFox/backend/server"
	"log"
	"net/http"
)

// main подготавливает БД, регистрирует маршруты и запускает HTTP-сервер.
func main() {
	err, db := database.Init()
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
