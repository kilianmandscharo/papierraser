package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint(t *testing.T) {
	testCases := []struct {
		line            Line
		other           Line
		shouldIntersect bool
		intersection    Point
	}{
		// 90deg intersection
		{
			line:            Line{from: Point{X: 1, Y: 1}, to: Point{X: 5, Y: 1}},
			other:           Line{from: Point{X: 2, Y: 4}, to: Point{X: 2, Y: 0}},
			shouldIntersect: true,
			intersection:    Point{X: 2, Y: 1},
		},
		// parallel
		{
			line:            Line{from: Point{X: 1, Y: 1}, to: Point{X: 5, Y: 1}},
			other:           Line{from: Point{X: 1, Y: 4}, to: Point{X: 5, Y: 4}},
			shouldIntersect: false,
			intersection:    Point{},
		},
		// starting point on starting point
		{
			line:            Line{from: Point{X: 1, Y: 1}, to: Point{X: 5, Y: 5}},
			other:           Line{from: Point{X: 1, Y: 1}, to: Point{X: 1, Y: 5}},
			shouldIntersect: true,
			intersection:    Point{X: 1, Y: 1},
		},
		// end point on end point
		{
			line:            Line{from: Point{X: 1, Y: 5}, to: Point{X: 5, Y: 5}},
			other:           Line{from: Point{X: 1, Y: 1}, to: Point{X: 5, Y: 5}},
			shouldIntersect: true,
			intersection:    Point{X: 5, Y: 5},
		},
	}

	for _, test := range testCases {
		intersection, ok := test.line.intersects(test.other)
		assert.Equal(t, test.shouldIntersect, ok)
		assert.Equal(t, test.intersection, intersection)
	}
}
