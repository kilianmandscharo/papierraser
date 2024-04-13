package game

import (
	"log"

	"github.com/gorilla/websocket"
)

type Race struct {
	Track    Track
	Players  []*Player
	Started  bool
	Finished bool
	Turn     int
	Winner   int
}

func NewRace() *Race {
	return &Race{Track: track, Players: make([]*Player, 0)}
}

func (r *Race) ConnectPlayer(addr string, conn *websocket.Conn) {
	if player := r.getPlayerByAddr(addr); player != nil {
		player.Conn = conn
	} else {
		id := r.nextId()
		newPlayer := NewPlayer(id, conn, addr)
		r.Players = append(r.Players, &newPlayer)
	}
	log.Printf("connected %s\n", addr)
}

func (r *Race) DisconnectPlayer(addr string) {
	if player := r.getPlayerByAddr(addr); player != nil {
		player.Conn = nil
		log.Printf("disconnected %s\n", addr)
	}
}

func (r *Race) UpdatePlayerName(addr string, name string) {
	if player := r.getPlayerByAddr(addr); player != nil {
		player.Name = name
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
	if player := r.getPlayerByAddr(addr); player != nil {
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
	if player := r.getPlayerByAddr(addr); player != nil {
		player.Ready = !player.Ready
	}
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
	for _, player := range r.Players {
		if len(player.Path) == 0 {
			return player
		}
	}
	return nil
}

func (r *Race) getPlayerById(id int) *Player {
	for _, player := range r.Players {
		if player.Id == id {
			return player
		}
	}
	return nil
}

func (r *Race) getPlayerByAddr(addr string) *Player {
	for _, player := range r.Players {
		if player.Addr == addr {
			return player
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
	player := r.getPlayerById(id)
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
	if player := r.getPlayerByAddr(addr); player != nil {
		player.Path = append(player.Path, pos)
	}
}

func (r *Race) MakeMove(pos Point) {
	if player := r.CurrentPlayer(); player != nil {
		pathTaken := Line{from: player.GetPosition(), to: pos}

		point, isIntersection := r.Track.lineIntersectsTrack(pathTaken)

		if isIntersection && point != player.GetPosition() {
			player.Move(point)
			player.Move(point)
			player.Crashed = true
		} else if track.pointOnTrackLines(pos) ||
			!r.Track.pointInsideTrack(pos) {
			player.Move(player.GetPosition())
			player.Crashed = true
		} else {
			player.Move(pos)
		}
		r.SetNextTurn()
	}
}

func (r *Race) CurrentPlayer() *Player {
	return r.getPlayerById(r.Turn)
}

func (r *Race) TogglePlayerCrashed(id int) {
	if player := r.getPlayerById(id); player != nil {
		player.Crashed = !player.Crashed
	}
}

func (r *Race) AllPlayersCrashed() bool {
	for _, player := range r.Players {
		if !player.Crashed {
			return false
		}
	}
	return true
}

func (r *Race) UncrashAllPlayers() {
	for _, player := range r.Players {
		player.Crashed = false
	}
}

func (r *Race) SetNextTurn() {
	for {
		var nextTurn int
		if r.Turn+1 > len(r.Players) {
			nextTurn = 1
		} else {
			nextTurn = r.Turn + 1
		}
		if player := r.getPlayerById(nextTurn); player != nil {
			if player.Crashed {
				allCrashed := r.AllPlayersCrashed()
				r.TogglePlayerCrashed(player.Id)
				r.Turn = nextTurn
				if allCrashed {
					r.UncrashAllPlayers()
					break
				}
			} else {
				r.Turn = nextTurn
				break
			}
		}
	}
}
