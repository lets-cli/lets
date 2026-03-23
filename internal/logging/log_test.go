package logging

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

func TestLoggingToStd(t *testing.T) {
	t.Run("should write log to correct std descriptor", func(t *testing.T) {
		stdOutMsg := "Log in std out"
		stdErrMsg := "Log in std err"

		var stdBuff bytes.Buffer

		var errBuff bytes.Buffer

		prevNoColor := color.NoColor
		color.NoColor = true
		defer func() {
			color.NoColor = prevNoColor
		}()

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
	prevNoColor := color.NoColor
	color.NoColor = false
	defer func() {
		color.NoColor = prevNoColor
	}()

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
