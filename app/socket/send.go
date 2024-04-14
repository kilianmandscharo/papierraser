package socket

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/kilianmandscharo/papierraser/components"
	"github.com/kilianmandscharo/papierraser/game"
)

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
	return "ClientActionLobby", buf.Bytes()
}

type newPositionPayload struct {
	Point    game.Point `json:"point"`
	PlayerId int        `json:"playerId"`
}

func renderNewPosition(playerToMove *game.Player, newPos game.Point) (string, []byte) {
	data, err := json.Marshal(newPositionPayload{
		Point:    newPos,
		PlayerId: playerToMove.Id,
	})
	if err != nil {
		log.Println(err)
		return "", []byte{}
	}
	return "ClientActionMove", data
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
	return "ClientActionTrack", buf.Bytes()
}
