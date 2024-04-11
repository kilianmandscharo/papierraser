package game

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Players = map[string]Player

type Player struct {
	Id    int
	Name  string
	Path  Path
	Conn  *websocket.Conn
	Ready bool
}

type Velocity struct {
	X, Y int
}

func NewPlayer(id int, conn *websocket.Conn) Player {
	return Player{
		Id:   id,
		Conn: conn,
		Name: fmt.Sprintf("Spieler %d", id),
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
			options = append(options, Point{X: optionX, Y: optionY})
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
