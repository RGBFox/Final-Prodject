package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	Port := os.Getenv("TODO_PORT")
	if Port == "" {
		Port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Fatal(http.ListenAndServe(":"+Port, nil))
}
