package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"github.com/kilianmandscharo/papierraser/components"
	"github.com/kilianmandscharo/papierraser/game"
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

func renderLobby(race *game.Race, target string) (string, []byte) {
	var buf bytes.Buffer
	err := components.Lobby(race.GetPlayersSorted(), target).Render(
		context.Background(),
		&buf,
	)
	if err != nil {
		log.Println(err)
		return "", []byte{}
	}
	return "Lobby", buf.Bytes()
}

func renderTrack(race *game.Race, target string) (string, []byte) {
	var buf bytes.Buffer
	err := components.Track(race, target).Render(
		context.Background(),
		&buf,
	)
	if err != nil {
		log.Println(err)
		return "", []byte{}
	}
	return "Track", buf.Bytes()
}

type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func parseMessage(payload []byte) (Message, error) {
	var message Message
	return message, json.Unmarshal(payload, &message)
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
				UpdateFunc: func(race *game.Race) {
					race.DisconnectPlayer(addr)
				},
				RenderFunc: func(race *game.Race, target string) (string, []byte) {
					if race.Started {
						return renderTrack(race, target)
					}
					return renderLobby(race, target)
				},
			}
			conn.Close()
		}()

		ch <- state.ActionRequest{
			GameId: gameId,
			UpdateFunc: func(race *game.Race) {
				race.ConnectPlayer(addr, conn)
			},
			RenderFunc: func(race *game.Race, target string) (string, []byte) {
				if race.Started {
					return renderTrack(race, target)
				}
				return renderLobby(race, target)
			},
		}

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			message, err := parseMessage(p)
			if err != nil {
				log.Println(err)
				return
			}

			switch message.Type {
			case "ActionNameChange":
				ch <- state.ActionRequest{
					GameId: gameId,
					UpdateFunc: func(race *game.Race) {
						race.UpdatePlayerName(addr, message.Data.(string))
					},
					RenderFunc: func(race *game.Race, target string) (string, []byte) {
						return renderLobby(race, target)
					},
				}
			case "ActionStart":
				ch <- state.ActionRequest{
					GameId: gameId,
					UpdateFunc: func(race *game.Race) {
						race.Start()
					},
					RenderFunc: func(race *game.Race, target string) (string, []byte) {
						return renderTrack(race, target)
					},
				}
			case "ActionChooseStartingPosition":
				ch <- state.ActionRequest{
					GameId: gameId,
					UpdateFunc: func(race *game.Race) {
						data := message.Data.(map[string]any)
						point := game.Point{X: int(data["x"].(float64)), Y: int(data["y"].(float64))}
						race.UpdateStartingPosition(addr, point)
					},
					RenderFunc: func(race *game.Race, target string) (string, []byte) {
						return renderTrack(race, target)
					},
				}
			default:
				log.Printf("unknown message type '%s' provided by client\n", message.Type)
			}
		}
	}
}
