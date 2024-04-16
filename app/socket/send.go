package socket

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/kilianmandscharo/papierraser/components"
	"github.com/kilianmandscharo/papierraser/enum"
	"github.com/kilianmandscharo/papierraser/game"
)

type drawRacePayload struct {
	Track    game.Track     `json:"track"`
	Players  []*game.Player `json:"players"`
	TargetId int            `json:"targetId"`
}

func drawRace(race *game.Race, target *game.Player) (enum.ClientAction, []byte) {
	data, err := json.Marshal(drawRacePayload{
		Track:    race.Track,
		Players:  race.Players,
		TargetId: target.Id,
	})
	if err != nil {
		log.Println(err)
		return "", []byte{}
	}
	return enum.ClientActionDrawRace, data
}

func drawLobby(race *game.Race, target *game.Player) (enum.ClientAction, []byte) {
	var buf bytes.Buffer
	err := components.Lobby(race, target).Render(
		context.Background(),
		&buf,
	)
	if err != nil {
		log.Println(err)
		return "", []byte{}
	}
	return enum.ClientActionDrawLobby, buf.Bytes()
}

type newPositionPayload struct {
	Point    game.Point `json:"point"`
	PlayerId int        `json:"playerId"`
}

func drawNewPosition(playerToMove *game.Player, newPos game.Point) (enum.ClientAction, []byte) {
	data, err := json.Marshal(newPositionPayload{
		Point:    newPos,
		PlayerId: playerToMove.Id,
	})
	if err != nil {
		log.Println(err)
		return "", []byte{}
	}
	return enum.ClientActionDrawNewPosition, data
}

// func renderTrack(race *game.Race, target *game.Player) (string, []byte) {
// 	var buf bytes.Buffer
// 	err := components.Race(race, target).Render(
// 		context.Background(),
// 		&buf,
// 	)
// 	if err != nil {
// 		log.Println(err)
// 		return "", []byte{}
// 	}
// 	return "ClientActionTrack", buf.Bytes()
// }
