package logger

import (
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"os"
)

func newJxLogger() Logger {
	_ = log.SetLevel(os.Getenv("LOG_LEVEL"))

	return Logger(log.Logger())
}
