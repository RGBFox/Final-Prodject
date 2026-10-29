package main

import (
	"RGBFox/backend/database"
	"RGBFox/backend/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := server.Init()
	err := database.Init()
	if err != nil {
		log.Fatal(err)
		fmt.Print(err)
	}
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
