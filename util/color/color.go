package color

var Color = &colorUtil{}

type colorUtil struct {
}

var (
	Default = "\u001B[39m"
	Reset   = "\u001B[30m"
)

func Black(str string) string {
	return ColorMap[black] + str + Default
}

func Red(str string) string {
	return ColorMap[red] + str + Default
}

func Green(str string) string {
	return ColorMap[green] + str + Default
}

func Yellow(str string) string {
	return ColorMap[yellow] + str + Default
}

func Blue(str string) string {
	return ColorMap[blue] + str + Default
}

func Magenta(str string) string {
	return ColorMap[magenta] + str + Default
}

func Cyan(str string) string {
	return ColorMap[cyan] + str + Default
}

func White(str string) string {
	return ColorMap[white] + str + Default
}
