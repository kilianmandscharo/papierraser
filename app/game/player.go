package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	Id      int
	Name    string
	Path    Path
	Addr    string
	Conn    *websocket.Conn
	Ready   bool
	Crashed bool
	Color   string
}

type Velocity struct {
	X, Y int
}

func NewPlayer(id int, conn *websocket.Conn, addr string) Player {
	return Player{
		Id:    id,
		Conn:  conn,
		Addr:  addr,
		Name:  fmt.Sprintf("Spieler %d", id),
		Color: randomColorHex(),
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

func randomColorHex() string {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))

	red := generator.Intn(256)
	green := generator.Intn(256)
	blue := generator.Intn(256)

	redHex := strconv.FormatInt(int64(red), 16)
	greenHex := strconv.FormatInt(int64(green), 16)
	blueHex := strconv.FormatInt(int64(blue), 16)

	if len(redHex) == 1 {
		redHex = "0" + redHex
	}
	if len(greenHex) == 1 {
		greenHex = "0" + greenHex
	}
	if len(blueHex) == 1 {
		blueHex = "0" + blueHex
	}

	colorHex := strings.ToUpper("#" + redHex + greenHex + blueHex)

	return colorHex
}
