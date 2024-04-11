package game

import (
	"cmp"
	"log"
	"slices"

	"github.com/gorilla/websocket"
)

type Race struct {
	Track    Track
	Players  Players
	Started  bool
	Finished bool
	Turn     int
	Winner   int
}

var track = Track{
	Width:  20,
	Height: 10,
	Outer: []Point{
		{X: 0, Y: 0}, {X: 20, Y: 0}, {X: 20, Y: 10}, {X: 0, Y: 10},
	},
	Inner: []Point{
		{X: 4, Y: 3}, {X: 16, Y: 3}, {X: 16, Y: 7}, {X: 4, Y: 7},
	},
	Finish: [2]Point{
		{X: 0, Y: 5}, {X: 4, Y: 5},
	},
}

func NewRace() *Race {
	return &Race{Track: track, Players: make(Players)}
}

func (r *Race) ConnectPlayer(addr string, conn *websocket.Conn) {
	if player, ok := r.Players[addr]; ok {
		player.Conn = conn
		r.Players[addr] = player
	} else {
		id := r.nextId()
		r.Players[addr] = NewPlayer(id, conn)
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

func (r *Race) AllPlayersReady() bool {
	for _, player := range r.Players {
		if !player.Ready {
			return false
		}
	}
	return true
}

func (r *Race) StartIfReady() {
	if r.AllPlayersReady() {
		r.Started = true
		r.Turn = 1
	}
}

func (r *Race) PlayerReady(addr string) bool {
	if player, ok := r.Players[addr]; ok {
		return player.Ready
	}
	return false
}

func (r *Race) End() {
	r.Finished = true
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

func (r *Race) TogglePlayerReady(addr string) {
	if player, ok := r.Players[addr]; ok {
		player.Ready = !player.Ready
		r.Players[addr] = player
	}
}

func (r *Race) GetPlayersSorted() []Player {
	playersSorted := make([]Player, len(r.Players))

	i := 0
	for _, player := range r.Players {
		playersSorted[i] = player
		i++
	}

	slices.SortFunc(playersSorted, func(a, b Player) int {
		return cmp.Compare(a.Id, b.Id)
	})

	return playersSorted
}

func (r *Race) StartingPositionsSet() bool {
	for _, player := range r.Players {
		if len(player.Path) == 0 {
			return false
		}
	}
	return true
}

func (r *Race) PickPlayerForStartingPosition() *Player {
	for _, player := range r.GetPlayersSorted() {
		if len(player.Path) == 0 {
			return &player
		}
	}
	return nil
}

func (r *Race) GetPlayerById(id int) *Player {
	for _, player := range r.Players {
		if player.Id == id {
			return &player
		}
	}
	return nil
}

func (r *Race) SomePlayerHasPosition(pos Point) bool {
	for _, player := range r.Players {
		if player.GetPosition() == pos {
			return true
		}
	}
	return false
}

func (r *Race) GetPlayerOptions(id int) []Point {
	player := r.GetPlayerById(id)
	if player == nil {
		return []Point{}
	}

	options := player.GetOptions()
	filteredOptions := []Point{}

	for _, option := range options {
		if !r.SomePlayerHasPosition(option) {
			filteredOptions = append(filteredOptions, option)
		}
	}

	return filteredOptions
}

func (r *Race) GetStartingPositionOptions() map[Point]bool {
	options := make(map[Point]bool)
	allOptions := r.Track.GetStartingPositionOptions()

	for _, option := range allOptions {
		selectable := true
		for _, player := range r.Players {
			if len(player.Path) > 0 && player.Path[0] == option {
				selectable = false
				break
			}
		}
		options[option] = selectable
	}

	return options
}

func (r *Race) UpdateStartingPosition(addr string, pos Point) {
	if player, ok := r.Players[addr]; ok {
		player.Path = append(player.Path, pos)
		r.Players[addr] = player
	}
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
