package routine

import (
	"log"

	"github.com/kilianmandscharo/papierraser/state"
)

type Request struct {
	gameId string
	ch     chan state.Channel
}

func NewRequest(gameId string, ch chan state.Channel) Request {
	return Request{gameId: gameId, ch: ch}
}

type Channel = chan Request
type ChannelSend = chan<- Request
type ChannelReceive = <-chan Request

func Handler(ch ChannelReceive) {
	routines := make(map[string]state.Channel)

	for request := range ch {
		if routineChan, ok := routines[request.gameId]; ok {
			request.ch <- routineChan
		} else {
			newChan := make(state.Channel)
			routines[request.gameId] = newChan

			go state.Handler(newChan)

			log.Println("started new state handler for", request.gameId)

			request.ch <- newChan
		}
	}
}
