package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"github.com/kilianmandscharo/papierraser/components"
	"github.com/kilianmandscharo/papierraser/race"
	"github.com/kilianmandscharo/papierraser/state"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

func main() {
	ch := make(chan state.ActionRequest)

	go state.Handler(ch)

	http.Handle("/", templ.Handler(components.Index()))
	http.HandleFunc("/ws", websocketHandler(ch))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getGameId(r *http.Request) (string, bool) {
	queryStrings := r.URL.Query()["id"]

	if len(queryStrings) == 0 {
		return "", false
	}

	return queryStrings[0], true
}

func renderLobby(race *race.Race) []byte {
	var buf bytes.Buffer
	err := components.Lobby(race.Players).Render(
		context.Background(),
		&buf,
	)
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()
}

func websocketHandler(ch chan<- state.ActionRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		gameId, ok := getGameId(r)
		if !ok {
			log.Println("no game id provided")
			return
		}

		addr := r.RemoteAddr

		defer func() {
			ch <- state.ActionRequest{
				GameId: gameId,
				Updater: func(race *race.Race) []byte {
					race.DisconnectPlayer(addr)
					return renderLobby(race)
				},
			}
			conn.Close()
		}()

		ch <- state.ActionRequest{
			GameId: gameId,
			Updater: func(race *race.Race) []byte {
				race.ConnectPlayer(addr, conn)
				return renderLobby(race)
			},
		}

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(messageType, p)
		}
	}
}
