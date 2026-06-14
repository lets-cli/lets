package logging

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/lets-cli/fang"
	log "github.com/sirupsen/logrus"
)

// InitLogging configures the global logrus logger. Pass a non-nil cs to enable
// lipgloss-styled error output when errWriter is a terminal.
func InitLogging(
	stdWriter io.Writer,
	errWriter io.Writer,
	cs fang.ColorSchemeFunc,
) {
	log.SetOutput(io.Discard)

	log.SetLevel(log.InfoLevel)

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

	log.SetFormatter(newFormatter(errWriter, cs))
}

// ExecLogger is used in Executor.
// If adds command chain in message like this:
// [foo=>bar] message.
type ExecLogger struct {
	log *log.Logger
	// command name
	name string
	// [a=>b]
	prefix string
	cache  map[string]*ExecLogger
}

func NewExecLogger() *ExecLogger {
	return &ExecLogger{
		log:   log.StandardLogger(),
		cache: make(map[string]*ExecLogger),
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
		prefix: color.GreenString("[%s]", name),
		cache:  make(map[string]*ExecLogger),
	}

	return l.cache[name]
}

func (l *ExecLogger) Info(format string, a ...any) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s %s", l.prefix, format)
	}

	l.log.Logf(log.InfoLevel, format, a...)
}

func (l *ExecLogger) Debug(format string, a ...any) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s %s", l.prefix, format)
	}

	l.log.Logf(log.DebugLevel, format, a...)
}
