package game

import "math"

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Path = []Point

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((y2-y1)*(y2-y1) + (x2-x1)*(x2-x1))
}

func (p *Point) liesOnLine(line Line) bool {
	x1, y1 := float64(line.from.X), float64(line.from.Y)
	x2, y2 := float64(line.to.X), float64(line.to.Y)
	x3, y3 := float64(p.X), float64(p.Y)
	return distance(x1, y1, x3, y3)+distance(x3, y3, x2, y2) ==
		distance(x1, y1, x2, y2)
}

type Line struct {
	from, to Point
}

func (l *Line) intersects(other Line) (Point, bool) {
	p1x, p1y := float64(l.from.X), float64(l.from.Y)
	p2x, p2y := float64(l.to.X), float64(l.to.Y)
	p3x, p3y := float64(other.from.X), float64(other.from.Y)
	p4x, p4y := float64(other.to.X), float64(other.to.Y)

	s1x, s1y := p2x-p1x, p2y-p1y
	s2x, s2y := p4x-p3x, p4y-p3y

	den := (-s2x*s1y + s1x*s2y)

	s := (-s1y*(p1x-p3x) + s1x*(p1y-p3y)) / den
	t := (s2x*(p1y-p3y) - s2y*(p1x-p3x)) / den

	if s >= 0 && s <= 1 && t >= 0 && t <= 1 {
		intersection := Point{X: int(p1x + (t * s1x)), Y: int(p1y + (t * s1y))}
		return intersection, true
	}

	return Point{}, false
}

func CastPoint(data map[string]any) Point {
	return Point{X: int(data["x"].(float64)), Y: int(data["y"].(float64))}
}
