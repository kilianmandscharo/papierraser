package utils

import (
	"fmt"
	"strconv"

	"github.com/a-h/templ"
	"github.com/kilianmandscharo/papierraser/types"
)

func GetPathString(path types.Path) string {
	pathString := ""

	for _, point := range path {
		pathString += fmt.Sprintf(" %d,%d", point.X, point.Y)
	}

	return pathString
}

func GetVerticalLineAttrs(x, height int) templ.Attributes {
	return templ.Attributes{
		"x1": strconv.Itoa(x),
		"x2": strconv.Itoa(x),
		"y1": "0",
		"y2": strconv.Itoa(height),
	}
}

func GetHorizontalLineAttrs(y, width int) templ.Attributes {
	return templ.Attributes{
		"x1": "0",
		"x2": strconv.Itoa(width),
		"y1": strconv.Itoa(y),
		"y2": strconv.Itoa(y),
	}
}
