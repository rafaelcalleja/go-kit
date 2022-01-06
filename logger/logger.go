package logger

type Logger interface {
	Printf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

var emptyLogger Logger = &nullLogger{}

type nullLogger struct {
}

func (n nullLogger) Printf(format string, args ...interface{}) {
	return
}

func (n nullLogger) Infof(format string, args ...interface{}) {
	return
}

func (n nullLogger) Debugf(format string, args ...interface{}) {
	return
}

func (n nullLogger) Errorf(format string, args ...interface{}) {
	return
}

func (n nullLogger) Fatalf(format string, args ...interface{}) {
	return
}
