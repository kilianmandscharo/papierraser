package types

import "github.com/gorilla/websocket"

type Player struct {
	Id        int
	Name      string
	Path      Path
	Conn      *websocket.Conn
	Connected bool
}

func NewPlayer(id int, name string, startingPosition Point, conn *websocket.Conn) Player {
	return Player{
		Id:   id,
		Name: name,
		Path: []Point{startingPosition},
		Conn: conn,
	}
}

func (p *Player) GetPosition() Point {
	return p.Path[len(p.Path)-1]
}

func (p *Player) GetOptions() []Point {
	pos := p.GetPosition()
	vel := p.GetVelocity()

	newPos := Point{X: pos.X + vel.X, Y: pos.Y + vel.Y}

	options := []Point{}

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			optionX := newPos.X + i
			optionY := newPos.Y + j
			if !(optionX == pos.X && optionY == pos.Y) {
				options = append(options, Point{X: optionX, Y: optionY})
			}
		}
	}

	return options
}

func (p *Player) GetVelocity() Velocity {
	if len(p.Path) <= 1 {
		return Velocity{X: 0, Y: 0}
	}

	p1 := p.Path[len(p.Path)-1]
	p2 := p.Path[len(p.Path)-2]

	return Velocity{X: p1.X - p2.X, Y: p1.Y - p2.Y}
}

func (p *Player) Move(pos Point) {
	p.Path = append(p.Path, pos)
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Velocity struct {
	X, Y int
}

type Path = []Point

type Track struct {
	Width, Height int
	Inner, Outer  Path
	Finish        [2]Point
}

type Race struct {
	Track    Track
	Players  []Player
	Finished bool
	Turn     int
	Winner   int
}

func NewRace(track Track) Race {
	return Race{Track: track}
}

func (r *Race) AddPlayer(name string, conn *websocket.Conn) {
	if len(r.Players) == 0 {
		r.Turn = 1
	}

	id := r.nextId()
	pos := r.GetStartingPosition()
	r.Players = append(r.Players, NewPlayer(id, name, pos, conn))
}

func (r *Race) GetStartingPosition() Point {
	p1 := r.Track.Finish[0]
	p2 := r.Track.Finish[1]

	midpointX := (p1.X + p2.X) / 2
	midpointY := (p1.Y + p2.Y) / 2

	return Point{X: midpointX, Y: midpointY}
}

func (r *Race) nextId() int {
	id := 1

	for _, player := range r.Players {
		if player.Id > id {
			id = player.Id
		}
	}

	return id
}

func (r *Race) Move(pos Point) {
	if index := r.currentPlayer(); index > -1 {
		r.Players[index].Move(pos)
	}
}

func (r *Race) currentPlayer() int {
	for index, player := range r.Players {
		if player.Id == r.Turn {
			return index
		}
	}

	return -1
}

type Connections = map[string]*websocket.Conn
