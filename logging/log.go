package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Log is the main application logger.
var Log = log.New()

func InitLogging(verbose bool) {
	logger := Log

	logger.Level = log.InfoLevel

	if verbose {
		logger.Level = log.DebugLevel
	}

	logger.Out = os.Stderr

	formatter := &Formatter{}
	log.SetFormatter(formatter)
	logger.Formatter = formatter
}
