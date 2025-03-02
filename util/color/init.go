package color

import "fmt"

var ColorMap = getColorMap()

type color = int

const (
	black = iota
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

func getColorMap() map[color]string {
	colors := make(map[color]string)
	for i := 0; i <= 7; i++ {
		colors[i] = fmt.Sprintf("\u001B[%dm", 30+i)
	}
	return colors
}
