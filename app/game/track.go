package game

type Track struct {
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Inner  Path     `json:"inner"`
	Outer  Path     `json:"outer"`
	Finish [2]Point `json:"finish"`
}

var track = Track{
	Width:  20,
	Height: 10,
	Outer: []Point{
		{X: 1, Y: 1}, {X: 19, Y: 1}, {X: 19, Y: 9}, {X: 1, Y: 9},
	},
	Inner: []Point{
		{X: 4, Y: 3}, {X: 16, Y: 3}, {X: 16, Y: 7}, {X: 4, Y: 7},
	},
	Finish: [2]Point{
		{X: 1, Y: 5}, {X: 4, Y: 5},
	},
}

func (t *Track) GetStartingPositionOptions() []Point {
	options := []Point{}

	x1 := t.Finish[0].X
	y1 := t.Finish[0].Y
	x2 := t.Finish[1].X
	y2 := t.Finish[1].Y

	for x1 < x2-1 || y1 < y2-1 {
		if x1 < x2 {
			x1++
		}
		if y1 < y1 {
			y1++
		}
		options = append(options, Point{X: x1, Y: y1})
	}

	return options
}

func (t *Track) getOuterLines() []Line {
	lines := []Line{}
	for i := 0; i < len(t.Outer)-1; i++ {
		lines = append(lines, Line{from: t.Outer[i], to: t.Outer[i+1]})
	}
	lines = append(lines, Line{from: t.Outer[len(t.Outer)-1], to: t.Outer[0]})
	return lines
}

func (t *Track) getInnerLines() []Line {
	lines := []Line{}
	for i := 0; i < len(t.Inner)-1; i++ {
		lines = append(lines, Line{from: t.Inner[i], to: t.Inner[i+1]})
	}
	lines = append(lines, Line{from: t.Inner[len(t.Inner)-1], to: t.Inner[0]})
	return lines
}

func (t *Track) lineIntersectsTrack(line Line) (Point, bool) {
	for _, other := range t.getInnerLines() {
		if intersection, ok := line.intersects(other); ok {
			return intersection, ok
		}
	}

	for _, other := range t.getOuterLines() {
		if intersection, ok := line.intersects(other); ok {
			return intersection, ok
		}
	}

	return Point{}, false
}

func (t *Track) pointInsideTrack(point Point) bool {
	return t.pointInsideOuter(point) && t.pointOutsideInner(point)
}

func (t *Track) pointInsideOuter(point Point) bool {
	line := Line{from: point, to: Point{X: 0, Y: point.Y}}
	intersections := 0
	for _, other := range t.getOuterLines() {
		if _, ok := line.intersects(other); ok {
			intersections++
		}
	}
	return intersections%2 != 0
}

func (t *Track) pointOutsideInner(point Point) bool {
	line := Line{from: point, to: Point{X: -1, Y: point.Y}}
	intersections := 0
	for _, other := range t.getInnerLines() {
		if _, ok := line.intersects(other); ok {
			intersections++
		}
	}
	return intersections%2 == 0
}

func (t *Track) pointOnTrackLines(point Point) bool {
	for _, line := range t.getInnerLines() {
		if point.liesOnLine(line) {
			return true
		}
	}

	for _, line := range t.getOuterLines() {
		if point.liesOnLine(line) {
			return true
		}
	}

	return false
}
