package state

import (
	"encoding/json"
	"log"

	"github.com/kilianmandscharo/papierraser/game"
)

type State = map[string]*game.Race

type ActionRequest struct {
	GameId     string
	UpdateFunc UpdateFunc
	RenderFunc RenderFunc
}

type UpdateFunc = func(*game.Race)
type RenderFunc func(*game.Race, string) (string, []byte)

type MessagePayload struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func newPayload(messageType string, data []byte) ([]byte, error) {
	return json.Marshal(MessagePayload{
		Type: messageType,
		Data: string(data),
	})
}

func newState() State {
	state := make(State)
	return state
}

func broadcast(state State, gameId string, renderFunc RenderFunc) {
	if race, ok := state[gameId]; ok {
		for addr, player := range race.Players {
			if player.Conn != nil {
				messageType, html := renderFunc(state[gameId], player.Name)
				payload, err := newPayload(messageType, html)
				if err != nil {
					log.Println("failed to create payload", err)
					return
				}
				err = player.Conn.WriteMessage(1, payload)
				if err != nil {
					log.Printf("failed to write message to %s\n", addr)
					return
				}
			}
		}
	}
}

func Handler(ch <-chan ActionRequest) {
	state := newState()

	for message := range ch {
		gameId := message.GameId
		updateFunc := message.UpdateFunc
		renderFunc := message.RenderFunc

		if _, ok := state[gameId]; !ok {
			state[gameId] = game.NewRace()
			log.Printf("Created new race for %s\n", gameId)
		}

		if updateFunc != nil {
			log.Printf("updating %s\n", gameId)
			updateFunc(state[gameId])
		}

		broadcast(state, gameId, renderFunc)
	}
}
