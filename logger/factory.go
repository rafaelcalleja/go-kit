package logger

func New() Logger {
	return newJxLogger()
}

func NewNullLogger() Logger {
	return emptyLogger
}
