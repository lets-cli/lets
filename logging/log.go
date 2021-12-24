package logging

import (
	"io"

	log "github.com/sirupsen/logrus"
)

// Log is the main application logger.
// InitLogging for logrus.
func InitLogging(
	verbose bool,
	stdWriter io.Writer,
	errWriter io.Writer,
) {
	log.SetOutput(io.Discard)

	log.SetLevel(log.InfoLevel)

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.AddHook(&WriterHook{
		Writer: stdWriter,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
			log.WarnLevel,
		},
	})

	log.AddHook(&WriterHook{
		Writer: errWriter,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		},
	})

	log.SetFormatter(&Formatter{})
}
