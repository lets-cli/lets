package config

import (
	"testing"
)

func TestFindConfig(t *testing.T) {
	t.Run("just find config", func(t *testing.T) {
		_, err := FindConfig("", "")
		if err != nil {
			t.Errorf("can not find test config: %s", err)
		}
	})
}
