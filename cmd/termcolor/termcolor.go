package termcolor

type TermColor interface {
	ColorInfo(i ...interface{}) string
	ColorStatus(i ...interface{}) string
	ColorWarning(i ...interface{}) string
	ColorError(i ...interface{}) string
	ColorBold(i ...interface{}) string
	ColorAnswer(i ...interface{}) string
}

