package logging

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/lets-cli/lets/env"
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

type ExecLogger struct {
	log *log.Logger
	// command name
	name string
	// lets: [a=>b]
	prefix string
	cache  map[string]*ExecLogger
}

func NewExecLogger() *ExecLogger {
	l := log.New()

	if env.IsDebug() {
		l.SetLevel(log.DebugLevel)
	}

	l.SetFormatter(&Formatter{})

	return &ExecLogger{
		log:    l,
		prefix: color.BlueString("lets:"),
		cache:  make(map[string]*ExecLogger),
	}
}

func (l *ExecLogger) Child(name string) *ExecLogger {
	if _, ok := l.cache[name]; ok {
		return l.cache[name]
	}

	if l.name != "" {
		name = fmt.Sprintf("%s => %s", l.name, name)
	}

	l.cache[name] = &ExecLogger{
		log:    l.log,
		name:   name,
		prefix: color.BlueString("lets: %s", color.GreenString("[%s]", name)),
		cache:  make(map[string]*ExecLogger),
	}

	return l.cache[name]
}

func (l *ExecLogger) Info(format string, a ...interface{}) {
	format = fmt.Sprintf("%s %s", l.prefix, color.BlueString(format))
	l.log.Logf(log.InfoLevel, format, a...)
}

func (l *ExecLogger) Debug(format string, a ...interface{}) {
	format = fmt.Sprintf("%s %s", l.prefix, color.BlueString(format))
	l.log.Logf(log.DebugLevel, format, a...)
}
