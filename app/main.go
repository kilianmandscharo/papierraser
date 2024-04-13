package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func renderLobby(race *game.Race, target *game.Player) (string, []byte) {
	var buf bytes.Buffer
	err := components.Lobby(race, target).Render(
		context.Background(),
		&buf,
	)
	if err != nil {
		log.Println(err)
		return "", []byte{}
	}
	return "Lobby", buf.Bytes()
}

type NewPositionPayload struct {
	Point    game.Point `json:"point"`
	PlayerId int        `json:"playerId"`
}

func renderNewPosition(playerToMove *game.Player, newPos game.Point) (string, []byte) {
	data, err := json.Marshal(NewPositionPayload{
		Point:    newPos,
		PlayerId: playerToMove.Id,
	})
	if err != nil {
		log.Println(err)
		return "", []byte{}
	}
	return "Move", data
}

func renderTrack(race *game.Race, target *game.Player) (string, []byte) {
	var buf bytes.Buffer
	err := components.Race(race, target).Render(
		context.Background(),
		&buf,
	)
	if err != nil {
		log.Println(err)
		return "", []byte{}
	}
	return "Track", buf.Bytes()
}

type MessageReceive struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func parseMessage(payload []byte) (MessageReceive, error) {
	var message MessageReceive
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

		defer disconnectPlayer(ch, gameId, addr, conn)

		connectPlayer(ch, gameId, addr, conn)

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
				handleActionNameChange(ch, gameId, addr, message)
			case "ActionToggleReady":
				handleActionToggleReady(ch, gameId, addr)
			case "ActionChooseStartingPosition":
				handleActionChooseStartingPosition(ch, gameId, addr, message)
			case "ActionMakeMove":
				handleActionMakeMove(ch, gameId, message)
			case "ActionMoveAnimationDone":
				fmt.Println("move animation done")
				handleActionMoveAnimationDone(ch, gameId)
			default:
				log.Printf("unknown message type '%s' provided by client\n", message.Type)
			}
		}
	}
}

func disconnectPlayer(ch chan<- state.ActionRequest, gameId string, addr string, conn *websocket.Conn) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.DisconnectPlayer(addr)
			return func(target *game.Player) (string, []byte) {
				if race.Started {
					return renderTrack(race, target)
				}
				return renderLobby(race, target)
			}
		},
	}
	conn.Close()
}

func connectPlayer(ch chan<- state.ActionRequest, gameId string, addr string, conn *websocket.Conn) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.ConnectPlayer(addr, conn)
			return func(target *game.Player) (string, []byte) {
				if race.Started {
					return renderTrack(race, target)
				}
				return renderLobby(race, target)
			}
		},
	}
}

func handleActionNameChange(ch chan<- state.ActionRequest, gameId string, addr string, message MessageReceive) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.UpdatePlayerName(addr, message.Data.(string))
			return func(target *game.Player) (string, []byte) {
				return renderLobby(race, target)
			}
		},
	}
}

func handleActionToggleReady(ch chan<- state.ActionRequest, gameId string, addr string) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.TogglePlayerReady(addr)
			race.StartIfReady()
			return func(target *game.Player) (string, []byte) {
				if race.AllPlayersReady() {
					return renderTrack(race, target)
				}
				return renderLobby(race, target)
			}
		},
	}
}

func handleActionChooseStartingPosition(ch chan<- state.ActionRequest, gameId string, addr string, message MessageReceive) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			race.UpdateStartingPosition(addr, game.CastPoint(message.Data))

			return func(target *game.Player) (string, []byte) {
				return renderTrack(race, target)
			}
		},
	}
}

func handleActionMakeMove(ch chan<- state.ActionRequest, gameId string, message MessageReceive) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			playerToMove := race.CurrentPlayer()
			movedTo, hasMoved := race.MakeMove(game.CastPoint(message.Data))

			return func(target *game.Player) (string, []byte) {
				if hasMoved {
					return renderNewPosition(playerToMove, movedTo)
				}
				return renderTrack(race, target)
			}
		},
	}
}

func handleActionMoveAnimationDone(ch chan<- state.ActionRequest, gameId string) {
	ch <- state.ActionRequest{
		GameId: gameId,
		UpdateFunc: func(race *game.Race) state.RenderFunc {
			return func(target *game.Player) (string, []byte) {
				return renderTrack(race, target)
			}
		},
	}
}
