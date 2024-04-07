package state

import (
	"log"

	"github.com/kilianmandscharo/papierraser/race"
)

type State = map[string]*race.Race

type ActionRequest struct {
	GameId  string
	Updater func(*race.Race) []byte
}

func newState() State {
	state := make(State)
	return state
}

func broadcast(state State, gameId string, payload []byte) {
	if race, ok := state[gameId]; ok {
		for addr, player := range race.Players {
			if player.Conn != nil {
				if err := player.Conn.WriteMessage(1, payload); err != nil {
					log.Printf("failed to write message to %s\n", addr)
				}
			}
		}
	}
}

func Handler(ch <-chan ActionRequest) {
	state := newState()

	for message := range ch {
		gameId := message.GameId
		updater := message.Updater

		if _, ok := state[gameId]; !ok {
			state[gameId] = race.New()
		}

		log.Printf("updating %s\n", gameId)
		html := updater(state[gameId])

		broadcast(state, gameId, html)
	}
}
