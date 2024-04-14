package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/kilianmandscharo/papierraser/components"
	"github.com/kilianmandscharo/papierraser/routine"
	"github.com/kilianmandscharo/papierraser/socket"
)

func main() {
	ch := make(chan routine.Request)

	go routine.Handler(ch)

	http.Handle("/", templ.Handler(components.Index()))
	http.HandleFunc("/ws", socket.Handler(ch))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
