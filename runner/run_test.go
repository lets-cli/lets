package runner

import (
	"strings"
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

func TestConvertChecksumMapToEnvList(t *testing.T) {
	findEnv := func(key string, list []string) bool {
		found := false

		for _, item := range list {
			if item == key {
				found = true
			}
		}

		return found
	}

	t.Run("should convert map to list of key=val", func(t *testing.T) {
		env := make(map[string]string)
		env["one"] = "111"
		env["two-two"] = "222"
		env["three_three"] = "333"
		envList := convertChecksumMapToEnvForCmd(env)

		joinedEnv := strings.Join(envList, ";")
		if !findEnv("LETS_CHECKSUM_ONE=111", envList) {
			t.Errorf("failed to convert env map to list. \nexp: %s\ngot: %s", "LETS_CHECKSUM_ONE=1", joinedEnv)
		}

		if !findEnv("LETS_CHECKSUM_TWO_TWO=222", envList) {
			t.Errorf("failed to convert env map to list. \nexp: %s\ngot: %s", "LETS_CHECKSUM_TWO_TWO=222", joinedEnv)
		}

		if !findEnv("LETS_CHECKSUM_THREE_THREE=333", envList) {
			t.Errorf("failed to convert env map to list. \nexp: %s\ngot: %s", "LETS_CHECKSUM_THREE_THREE=333", joinedEnv)
		}
	})
}

func TestComposeEnv(t *testing.T) {
	t.Run("should compose env", func(t *testing.T) {
		toCompose := []string{"A=1"}
		toCompose1 := []string{"B=2"}
		total := len(toCompose) + len(toCompose1)
		env := composeEnvs(toCompose, toCompose1)
		if len(env) != total {
			t.Errorf("composed env len different from expected: exp: %d, got: %d", total, len(env))
		}
		if env[0] != "A=1" {
			t.Errorf("first element from composed env different from expected: exp: %s, got: %s", toCompose[0], env[0])
		}
		if env[1] != "B=2" {
			t.Errorf("first element from composed env different from expected: exp: %s, got: %s", toCompose1[0], env[1])
		}
	})
}
