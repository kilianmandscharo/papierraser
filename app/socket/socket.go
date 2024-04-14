package socket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kilianmandscharo/papierraser/routine"
	"github.com/kilianmandscharo/papierraser/state"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

func Handler(routineChan routine.ChannelSend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		addr := r.RemoteAddr
		gameId, ok := getGameId(r)
		if !ok {
			log.Println("no game id provided")
			return
		}

		receiveChan := make(chan state.Channel)
		routineChan <- routine.NewRequest(gameId, receiveChan)
		stateChan := <-receiveChan

		defer handleReceiveActionDisconnectPlayer(stateChan, gameId, addr, conn)

		handleReceiveActionConnectPlayer(stateChan, gameId, addr, conn)

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
			case "ServerActionNameChange":
				handleReceiveActionNameChange(stateChan, gameId, addr, message)
			case "ServerActionToggleReady":
				handleReceiveActionToggleReady(stateChan, gameId, addr)
			case "ServerActionChooseStartingPosition":
				handleReceiveActionChooseStartingPosition(stateChan, gameId, addr, message)
			case "ServerActionMakeMove":
				handleReceiveActionMakeMove(stateChan, gameId, message)
			case "ServerActionMoveAnimationDone":
				handleReceiveActionMoveAnimationDone(stateChan, gameId)
			default:
				log.Printf("unknown message type '%s' provided by client\n", message.Type)
			}
		}
	}
}

func getGameId(r *http.Request) (string, bool) {
	queryStrings := r.URL.Query()["id"]

	if len(queryStrings) == 0 {
		return "", false
	}

	return queryStrings[0], true
}

type messageReceive struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func parseMessage(payload []byte) (messageReceive, error) {
	var message messageReceive
	return message, json.Unmarshal(payload, &message)
}
