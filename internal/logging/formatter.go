package logging

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

// LogRepresenter is an interface for objects that can format themselves for
// logging.
type LogRepresenter interface {
	Repr() string
}

// Formatter formats a log entry in a human readable way.
type Formatter struct{}

// Format implements the log.Formatter interface.
func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	buff := &bytes.Buffer{}
	parts := []string{color.BlueString("lets:")}

	if data := writeData(entry.Data); data != "" {
		parts = append(parts, data)
	}

	parts = append(parts, formatMessage(entry))

	buff.WriteString(strings.Join(parts, " "))
	buff.WriteString("\n")

	return buff.Bytes(), nil
}

func formatMessage(entry *log.Entry) string {
	if entry.Level == log.DebugLevel {
		return color.BlueString(entry.Message)
	}

	return entry.Message
}

func writeData(fields log.Fields) string {
	var buff []string

	for key, value := range fields {
		switch value := value.(type) {
		case LogRepresenter:
			buff = append(buff, value.Repr())
		default:
			buff = append(buff, fmt.Sprintf("%v=%v", key, value))
		}
	}

	return strings.Join(buff, " ")
}
