package main

import (
	"bytes"
	"context"
	"encoding/json"
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

var track = types.Track{
	Width:  20,
	Height: 10,
	Outer: []types.Point{
		{X: 0, Y: 0}, {X: 20, Y: 0}, {X: 20, Y: 10}, {X: 0, Y: 10},
	},
	Inner: []types.Point{
		{X: 4, Y: 3}, {X: 16, Y: 3}, {X: 16, Y: 7}, {X: 4, Y: 7},
	},
	Finish: [2]types.Point{
		{X: 0, Y: 5}, {X: 4, Y: 5},
	},
}

func main() {
	race := types.NewRace(track)
	race.AddPlayer("Mason")

	http.Handle("/", templ.Handler(components.Index(race)))
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn := unwrap(upgrader.Upgrade(w, r, nil))

	race := types.NewRace(track)
	race.AddPlayer("Mason")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

    var point types.Point
    err = json.Unmarshal(p, &point)
    if err != nil {
      fmt.Println("ERROR:", err)
      return
    }
    
    race.Move(point)

    var buf bytes.Buffer
    err = components.Track(race).Render(context.Background(), &buf)
    if err != nil {
      fmt.Println(err)
      return;
    }

    if err := conn.WriteMessage(messageType, buf.Bytes()); err != nil {
      fmt.Println(err)
      return;
    }
	}
}

func unwrap[T any](value T, err error) T {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err)
	}
	return value
}
