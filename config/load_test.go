package config

import (
	"testing"

)

func TestLoadConfig(t *testing.T) {
	t.Run("just read config", func(t *testing.T) {
		_, err := Load("0.0.0-test")
		if err != nil {
			t.Errorf("can not load test config: %s", err)
		}
	})
}
