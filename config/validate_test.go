package config

import (
	"testing"
)

func TestValidateCircularDeps(t *testing.T) {
	t.Run("command skip itself", func(t *testing.T) {
		testCfg := &Config{
			Commands: make(map[string]Command),
		}
		testCfg.Commands["a-cmd"] = Command{
			Name:    "a-cmd",
			Depends: []string{"noop"},
		}
		testCfg.Commands["b-cmd"] = Command{
			Name:    "b-cmd",
			Depends: []string{"noop"},
		}
		err := validateCircularDepends(testCfg)
		if err != nil {
			t.Errorf("checked itself when validation circular depends. got:  %s", err)
		}
	})

	t.Run("command with similar name should not fail validation", func(t *testing.T) {
		testCfg := &Config{
			Commands: make(map[string]Command),
		}
		testCfg.Commands["a-cmd"] = Command{
			Name:    "a-cmd",
			Depends: []string{"b1-cmd"},
		}
		testCfg.Commands["b"] = Command{
			Name:    "b",
			Depends: []string{"a-cmd"},
		}
		testCfg.Commands["b1-cmd"] = Command{
			Name:    "b1-cmd",
			Depends: []string{"noop"},
		}
		err := validateCircularDepends(testCfg)
		if err != nil {
			t.Errorf("checked itself when validation circular depends. got:  %s", err)
		}
	})

	t.Run("validation should fail", func(t *testing.T) {
		testCfg := &Config{
			Commands: make(map[string]Command),
		}
		testCfg.Commands["a-cmd"] = Command{
			Name:    "a-cmd",
			Depends: []string{"b1-cmd"},
		}
		testCfg.Commands["b"] = Command{
			Name:    "b",
			Depends: []string{"a-cmd"},
		}
		testCfg.Commands["b1-cmd"] = Command{
			Name:    "b1-cmd",
			Depends: []string{"a-cmd"},
		}
		err := validateCircularDepends(testCfg)

		if err == nil {
			t.Errorf("validation should fail. got: %s", err)
		}
	})
}
