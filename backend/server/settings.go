package main

import (
	"log"
	"net/http"
)

func server() {
	http.Handle("/", http.FileServer(http.Dir("./web/index.html")))
	http.Handle("/css/style.css", http.FileServer(http.Dir("./web/css/style.css")))
	http.Handle("/favicon.ico", http.FileServer(http.Dir("./web/favicon.ico")))
	http.Handle("/js/scripts.min.js", http.FileServer(http.Dir("./web/js/scripts.min.js")))

	log.Fatal(http.ListenAndServe(":7540", nil))
}

func another() {

}
