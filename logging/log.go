package logging

import (
	"io"

	log "github.com/sirupsen/logrus"
)

// Log is the main application logger.
var Log = log.New()

// InitLogging for logrus.
func InitLogging(
	verbose bool,
	stdWriter io.Writer,
	errWriter io.Writer,
) {
	Log.SetOutput(io.Discard)

	logger := Log

	logger.Level = log.InfoLevel

	if verbose {
		logger.Level = log.DebugLevel
	}

	Log.AddHook(&WriterHook{
		Writer: stdWriter,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
			log.WarnLevel,
		},
	})

	Log.AddHook(&WriterHook{
		Writer: errWriter,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		},
	})

	formatter := &Formatter{}
	Log.SetFormatter(formatter)
	logger.Formatter = formatter
}
