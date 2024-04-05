package types

type Player struct {
  Id string
	Name string
}

type Point struct {
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
	Paths    map[string]Path
	Players  []Player
	Finished bool
	Turn     string
	Winner   string
}
