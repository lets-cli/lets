package logging

import (
	"bytes"
	"testing"
)

func TestLoggingToStd(t *testing.T) {
	t.Run("should write log to correct std descriptor", func(t *testing.T) {
		stdOutMsg := "Log in std out"
		stdErrMsg := "Log in std err"

		var stdBuff bytes.Buffer

		var errBuff bytes.Buffer

		InitLogging(false, &stdBuff, &errBuff)

		Log.Info(stdOutMsg)
		Log.Error(stdErrMsg)

		// coz log adds line break for output
		if stdBuff.String() != stdOutMsg+"\n" {
			t.Errorf("stdBuff != stdOutMsg plz check your init stdWriter")
		}

		if errBuff.String() != stdErrMsg+"\n" {
			t.Errorf("errBuff != stdErrMsg plz check your init errWriter")
		}
	})
}
