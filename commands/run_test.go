package commands

import (
	"testing"
)

func TestConvertEnvMapToList(t *testing.T) {
	t.Run("should convert map to list of key=val", func(t *testing.T) {
		env := make(map[string]string)
		env["ONE"] = "1"
		envList := convertEnvMapToList(env)

		exp := "ONE=1"
		if envList[0] != exp {
			t.Errorf("failed to convert env map to list. \nexp: %s\ngot: %s", exp, envList[0])
		}
	})
}

func TestConvertChecksumMapToEnvList(t *testing.T) {
	t.Run("should convert map to list of key=val", func(t *testing.T) {
		env := make(map[string]string)
		env["one"] = "111"
		env["two-two"] = "222"
		env["three_three"] = "333"
		envList := convertChecksumMapToEnvForCmd(env)

		if envList[0] != "LETS_CHECKSUM_ONE=111" {
			t.Errorf("failed to convert env map to list. \nexp: %s\ngot: %s", "LETS_CHECKSUM_ONE=1", envList[0])
		}

		if envList[1] != "LETS_CHECKSUM_TWO_TWO=222" {
			t.Errorf("failed to convert env map to list. \nexp: %s\ngot: %s", "LETS_CHECKSUM_TWO_TWO=222", envList[1])
		}

		if envList[2] != "LETS_CHECKSUM_THREE_THREE=333" {
			t.Errorf("failed to convert env map to list. \nexp: %s\ngot: %s", "LETS_CHECKSUM_THREE_THREE=333", envList[2])
		}

	})
}
