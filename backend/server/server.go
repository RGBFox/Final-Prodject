package server

import (
	"net/http"
	"os"
)

func Init() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	return port
}
