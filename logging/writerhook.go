package logging

import (
	"io"

	log "github.com/sirupsen/logrus"
)

// WriterHook struct for routing std depending on lvl
type WriterHook struct {
	Writer    io.Writer
	LogLevels []log.Level
}

// Fire method prosees entry for Writer
func (hook *WriterHook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels geter for list of lvls
func (hook *WriterHook) Levels() []log.Level {
	return hook.LogLevels
}
