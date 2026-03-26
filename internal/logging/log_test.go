package logging

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

var colorNoColorMu sync.Mutex

func setNoColorForTest(t *testing.T, noColor bool) {
	t.Helper()

	colorNoColorMu.Lock()
	prevNoColor := color.NoColor
	color.NoColor = noColor

	t.Cleanup(func() {
		color.NoColor = prevNoColor
		colorNoColorMu.Unlock()
	})
}

func TestLoggingToStd(t *testing.T) {
	t.Run("should write log to correct std descriptor", func(t *testing.T) {
		stdOutMsg := "Log in std out"
		stdErrMsg := "Log in std err"

		var stdBuff bytes.Buffer

		var errBuff bytes.Buffer

		setNoColorForTest(t, true)

		InitLogging(&stdBuff, &errBuff)

		log.Info(stdOutMsg)
		log.Error(stdErrMsg)

		// coz log adds line break for output
		if stdBuff.String() != "lets: "+stdOutMsg+"\n" {
			t.Errorf("stdBuff != stdOutMsg plz check your init stdWriter")
		}

		if errBuff.String() != "lets: "+stdErrMsg+"\n" {
			t.Errorf("errBuff != stdErrMsg plz check your init errWriter")
		}
	})
}

func TestFormatterColorsDebugMessages(t *testing.T) {
	setNoColorForTest(t, false)

	line, err := (&Formatter{}).Format(&log.Entry{
		Level:   log.DebugLevel,
		Message: "debug message",
	})
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	expected := color.BlueString("lets:") + " " + color.BlueString("debug message") + "\n"
	if string(line) != expected {
		t.Fatalf("unexpected debug line: %q", string(line))
	}
}

func TestFormatterFormatsLevelsAndFields(t *testing.T) {
	setNoColorForTest(t, true)

	tests := []struct {
		name   string
		entry  *log.Entry
		fields []string
	}{
		{
			name: "info_no_data",
			entry: &log.Entry{
				Level:   log.InfoLevel,
				Message: "info message",
				Data:    log.Fields{},
			},
		},
		{
			name: "warn_with_single_field",
			entry: &log.Entry{
				Level:   log.WarnLevel,
				Message: "warn message",
				Data: log.Fields{
					"foo": "bar",
				},
			},
			fields: []string{"foo=bar"},
		},
		{
			name: "error_with_multiple_fields",
			entry: &log.Entry{
				Level:   log.ErrorLevel,
				Message: "error message",
				Data: log.Fields{
					"alpha": "one",
					"beta":  "two",
				},
			},
			fields: []string{"alpha=one", "beta=two"},
		},
		{
			name: "trace_no_data",
			entry: &log.Entry{
				Level:   log.TraceLevel,
				Message: "trace message",
				Data:    log.Fields{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lineBytes, err := (&Formatter{}).Format(tt.entry)
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			line := string(lineBytes)
			if !strings.HasPrefix(line, "lets: ") {
				t.Fatalf("line does not start with expected prefix: %q", line)
			}

			if strings.Contains(line, "\x1b[") {
				t.Fatalf("expected non-colorized output for non-debug levels, got: %q", line)
			}

			if len(tt.fields) == 0 {
				expected := "lets: " + tt.entry.Message + "\n"
				if line != expected {
					t.Fatalf("unexpected formatted line for empty Data.\nexpected: %q\ngot:      %q", expected, line)
				}

				return
			}

			msgIdx := strings.LastIndex(line, tt.entry.Message)
			if msgIdx == -1 {
				t.Fatalf("message %q not found in line: %q", tt.entry.Message, line)
			}

			if strings.Contains(line, "  ") {
				t.Fatalf("unexpected double spaces in line: %q", line)
			}

			for _, field := range tt.fields {
				fieldIdx := strings.Index(line, field)
				if fieldIdx == -1 {
					t.Fatalf("field %q not found in line: %q", field, line)
				}

				if fieldIdx <= len("lets:") || fieldIdx >= msgIdx {
					t.Fatalf("field %q not positioned between prefix and message in line: %q", field, line)
				}
			}
		})
	}
}
