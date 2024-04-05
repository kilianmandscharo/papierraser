package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"github.com/kilianmandscharo/papierraser/components"
	"github.com/kilianmandscharo/papierraser/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func main() {
	track := types.Track{
		Width:  50,
		Height: 50,
		Outer: []types.Point{
			{X: 0, Y: 0}, {X: 50, Y: 0}, {X: 50, Y: 50}, {X: 0, Y: 50},
		},
		Inner: []types.Point{
			{X: 15, Y: 15}, {X: 35, Y: 15}, {X: 35, Y: 35}, {X: 15, Y: 35},
		},
		Finish: [2]types.Point{
			{X: 0, Y: 25}, {X: 15, Y: 25},
		},
	}

	paths := make(map[string]types.Path)
	paths["1"] = types.Path{
		{X: 5, Y: 25},
	}
	paths["2"] = types.Path{
		{X: 10, Y: 25},
	}

	race := types.Race{
		Players: []types.Player{
			{Id: "1", Name: "Mason"},
			{Id: "2", Name: "Dixon"},
		},
		Paths: paths,
		Track: track,
	}

	http.Handle("/", templ.Handler(components.Index(race)))
	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn := unwrap(upgrader.Upgrade(w, r, nil))

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(messageType, string(p))
	}
}

func unwrap[T any](value T, err error) T {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err)
	}
	return value
}
