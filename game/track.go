package game

type Track struct {
	Width, Height int
	Inner, Outer  Path
	Finish        [2]Point
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
