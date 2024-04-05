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
		pathString += fmt.Sprintf(" %d,%d", point.X * 5, point.Y * 5)
	}

	return pathString
}

func GetVerticalLineAttrs(x, height int) templ.Attributes {
	return templ.Attributes{
		"x1": strconv.Itoa(x * 5),
		"x2": strconv.Itoa(x * 5),
		"y1": "0",
		"y2": strconv.Itoa(height * 5),
	}
}

func GetHorizontalLineAttrs(y, width int) templ.Attributes {
	return templ.Attributes{
		"x1": "0",
		"x2": strconv.Itoa(width * 5),
		"y1": strconv.Itoa(y * 5),
		"y2": strconv.Itoa(y * 5),
	}
}
