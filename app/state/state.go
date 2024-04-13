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

type UpdateFunc = func(*game.Race) RenderFunc
type RenderFunc func(*game.Player) (string, []byte)

type MessageSend struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func newPayload(messageType string, data []byte) ([]byte, error) {
	return json.Marshal(MessageSend{
		Type: messageType,
		Data: string(data),
	})
}

func newState() State {
	state := make(State)
	return state
}

func Handler(ch <-chan ActionRequest) {
	state := newState()

	for message := range ch {
		var race *game.Race

		if r, ok := state[message.GameId]; !ok {
			newRace := game.NewRace()
			state[message.GameId] = newRace
			race = newRace
			log.Printf("created new race for %s\n", message.GameId)
		} else {
			race = r
		}

		var renderFunc RenderFunc

		if message.UpdateFunc != nil {
			log.Printf("updating %s\n", message.GameId)
			renderFunc = message.UpdateFunc(race)
		}

		for _, player := range race.Players {
			if player.Conn == nil {
				continue
			}

			messageType, html := renderFunc(player)

			payload, err := newPayload(messageType, html)
			if err != nil {
				log.Println("failed to create payload:", err)
				return
			}

			err = player.Conn.WriteMessage(1, payload)
			if err != nil {
				log.Printf("failed to write message to %s\n", player.Name)
				return
			}
		}
	}
}
