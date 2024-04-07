package race

import (
	"cmp"
	"log"
	"slices"

	"github.com/gorilla/websocket"
	"github.com/kilianmandscharo/papierraser/types"
)

type Race struct {
	Track    types.Track
	Players  types.Players
	Finished bool
	Turn     int
	Winner   int
}

var track = types.Track{
	Width:  20,
	Height: 10,
	Outer: []types.Point{
		{X: 0, Y: 0}, {X: 20, Y: 0}, {X: 20, Y: 10}, {X: 0, Y: 10},
	},
	Inner: []types.Point{
		{X: 4, Y: 3}, {X: 16, Y: 3}, {X: 16, Y: 7}, {X: 4, Y: 7},
	},
	Finish: [2]types.Point{
		{X: 0, Y: 5}, {X: 4, Y: 5},
	},
}

func New() *Race {
	return &Race{Track: track, Players: make(types.Players)}
}

func (r *Race) ConnectPlayer(addr string, conn *websocket.Conn) {
	if player, ok := r.Players[addr]; ok {
		player.Conn = conn
		r.Players[addr] = player
	} else {
		id := r.nextId()
		r.Players[addr] = types.NewPlayer(id, conn)
	}
	log.Printf("connected %s\n", addr)
}

func (r *Race) DisconnectPlayer(addr string) {
	if player, ok := r.Players[addr]; ok {
		player.Conn = nil
		r.Players[addr] = player
		log.Printf("disconnected %s\n", addr)
	}
}

func (r *Race) UpdatePlayerName(addr string, name string) {
	if player, ok := r.Players[addr]; ok {
		player.Name = name
		r.Players[addr] = player
	}
}

func (r *Race) numberOfPlayers() int {
	return len(r.Players)
}

func (r *Race) nextId() int {
	id := 0

	for _, player := range r.Players {
		if player.Id > id {
			id = player.Id
		}
	}

	return id + 1
}

func (r *Race) GetPlayersSorted() []types.Player {
	playersSorted := make([]types.Player, len(r.Players))

	i := 0
	for _, player := range r.Players {
		playersSorted[i] = player
		i++
	}

	slices.SortFunc(playersSorted, func(a, b types.Player) int {
		return cmp.Compare(a.Id, b.Id)
	})

	return playersSorted
}

// func (r *Race) Move(pos Point) {
// 	if index := r.currentPlayer(); index > -1 {
// 		r.Players[index].Move(pos)
// 	}
// }
//
// func (r *Race) currentPlayer() int {
// 	for index, player := range r.Players {
// 		if player.Id == r.Turn {
// 			return index
// 		}
// 	}
//
// 	return -1
// }

// func (r *Race) GetStartingPosition() Point {
// 	p1 := r.Track.Finish[0]
// 	p2 := r.Track.Finish[1]
//
// 	midpointX := (p1.X + p2.X) / 2
// 	midpointY := (p1.Y + p2.Y) / 2
//
// 	return Point{X: midpointX, Y: midpointY}
// }
