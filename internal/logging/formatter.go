package logging

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/fatih/color"
	"github.com/lets-cli/fang"
	log "github.com/sirupsen/logrus"
)

// LogRepresenter is an interface for objects that can format themselves for
// logging.
type LogRepresenter interface {
	Repr() string
}

type errorStyles struct {
	header lipgloss.Style
	text   lipgloss.Style
}

// Formatter formats a log entry in a human readable way.
type Formatter struct {
	errorStyles *errorStyles // nil when output is not a TTY or no scheme given
}

// newFormatter builds a Formatter, enabling lipgloss error styling when
// errWriter is a terminal and cs is non-nil.
func newFormatter(errWriter io.Writer, cs fang.ColorSchemeFunc) *Formatter {
	f := &Formatter{}
	if cs == nil {
		return f
	}

	file, ok := errWriter.(term.File)
	if !ok || !term.IsTerminal(file.Fd()) {
		return f
	}

	isDark := lipgloss.HasDarkBackground(os.Stdin, file)
	scheme := cs(lipgloss.LightDark(isDark))

	w, _, err := term.GetSize(file.Fd())
	if err != nil || w == 0 {
		w = 160
	}
	if w > 160 {
		w = 160
	}

	f.errorStyles = &errorStyles{
		header: lipgloss.NewStyle().
			Foreground(scheme.ErrorHeader[0]).
			Background(scheme.ErrorHeader[1]).
			Bold(true).
			Padding(0, 1).
			Margin(1).
			MarginLeft(2).
			SetString("ERROR"),
		text: lipgloss.NewStyle().
			MarginLeft(2).
			Width(w - 2),
	}

	return f
}

// Format implements the log.Formatter interface.
func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	if entry.Level == log.ErrorLevel && f.errorStyles != nil {
		return f.formatStyledError(entry), nil
	}

	buff := &bytes.Buffer{}
	parts := []string{formatPrefix(entry)}

	if data := writeData(entry.Data); data != "" {
		parts = append(parts, data)
	}

	parts = append(parts, formatMessage(entry))

	buff.WriteString(strings.Join(parts, " "))
	buff.WriteString("\n")

	return buff.Bytes(), nil
}

func (f *Formatter) formatStyledError(entry *log.Entry) []byte {
	var buf bytes.Buffer
	buf.WriteString(f.errorStyles.header.String())
	buf.WriteString("\n")
	buf.WriteString(f.errorStyles.text.Render(capitalizeFirst(entry.Message) + "."))
	buf.WriteString("\n\n")
	return buf.Bytes()
}

func formatPrefix(entry *log.Entry) string {
	if entry.Level == log.DebugLevel {
		return color.BlueString("lets:")
	}

	return "lets:"
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

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
