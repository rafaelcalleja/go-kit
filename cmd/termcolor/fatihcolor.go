package termcolor

import "github.com/fatih/color"

var colorMap = map[string]color.Attribute{
	"bold": color.Bold,

	"black":   color.FgBlack,
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
	"white":   color.FgWhite,
}

type fatihColor struct {
}

func (f fatihColor) ColorInfo(i ...interface{}) string {
	return color.New(colorMap["green"]).Sprint(i...)
}

func (f fatihColor) ColorStatus(i ...interface{}) string {
	return color.New(colorMap["blue"]).Sprint(i...)
}

func (f fatihColor) ColorWarning(i ...interface{}) string {
	return color.New(colorMap["yellow"]).Sprint(i...)
}

func (f fatihColor) ColorError(i ...interface{}) string {
	return color.New(colorMap["red"]).Sprint(i...)
}

func (f fatihColor) ColorBold(i ...interface{}) string {
	return color.New(colorMap["bold"]).Sprint(i...)
}

func (f fatihColor) ColorAnswer(i ...interface{}) string {
	return color.New(colorMap["cyan"]).Sprint(i...)
}

func newFatihColor() fatihColor {
	return *new(fatihColor)
}
