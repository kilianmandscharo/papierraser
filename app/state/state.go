package state

import (
	"encoding/json"
	"log"

	"github.com/kilianmandscharo/papierraser/enum"
	"github.com/kilianmandscharo/papierraser/game"
)

type UpdateFunc = func(*game.Race) RenderFunc
type RenderFunc func(*game.Player) (enum.ClientAction, []byte)

type ActionRequest struct {
	GameId     string
	UpdateFunc UpdateFunc
}

type messageSend struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func newPayload(messageType string, data []byte) ([]byte, error) {
	return json.Marshal(messageSend{
		Type: messageType,
		Data: string(data),
	})
}

type Channel = chan ActionRequest
type ChannelReceive = <-chan ActionRequest

func Handler(ch ChannelReceive) {
	race := game.NewRace()

	for request := range ch {
		if request.UpdateFunc == nil {
			return
		}

		renderFunc := request.UpdateFunc(race)

		for _, player := range race.Players {
			if player.Conn == nil || renderFunc == nil {
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
