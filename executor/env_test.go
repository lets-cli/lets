package executor

import (
	"testing"
)

func TestConvertEnvMapToList(t *testing.T) {
	t.Run("should convert map to list of key=val", func(t *testing.T) {
		env := make(map[string]string, 1)
		env["ONE"] = "1"
		envList := convertEnvMapToList(env)
		exp := "ONE=1"
		if envList[0] != exp {
			t.Errorf("failed to convert env map to list. \nexp: %s\ngot: %s", exp, envList[0])
		}
	})
}
