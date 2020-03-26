package config

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Run("just read config", func(t *testing.T) {
		_, err := Load("lets.yaml", "..", "0.0.0-test")
		if err != nil {
			t.Errorf("can not read test config: %s", err)
		}
	})
}
