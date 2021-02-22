package config

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Run("just read config", func(t *testing.T) {
		cp, err := FindConfig()
		if err != nil {
			t.Errorf("can not find test config: %s", err)
		}

		_, err = LoadFromFile(cp, "0.0.0-test")
		if err != nil {
			t.Errorf("can not read test config: %s", err)
		}
	})
}
