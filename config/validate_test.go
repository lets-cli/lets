package config

import (
	"testing"

	"github.com/lets-cli/lets/config/config"
)

func TestValidateCommandInDependsExists(t *testing.T) {
	t.Run("command depends on non-existing command", func(t *testing.T) {
		testCfg := &config.Config{
			Commands: make(map[string]config.Command),
		}
		testCfg.Commands["foo"] = config.Command{
			Name: "foo",
			Depends: map[string]config.Dep{
				"bar": {Name: "bar"},
			},
		}
		err := validateCommandInDependsExists(testCfg)
		if err == nil {
			t.Error("command foo depends on non-existing command bar. Must fail")
		}
	})
}

func TestValidateCircularDeps(t *testing.T) {
	t.Run("command skip itself", func(t *testing.T) {
		testCfg := &config.Config{
			Commands: make(map[string]config.Command),
		}
		testCfg.Commands["a-cmd"] = config.Command{
			Name:    "a-cmd",
			Depends: map[string]config.Dep{},
		}
		testCfg.Commands["b-cmd"] = config.Command{
			Name:    "b-cmd",
			Depends: map[string]config.Dep{},
		}
		err := validateCircularDepends(testCfg)
		if err != nil {
			t.Errorf("checked itself when validation circular depends. got:  %s", err)
		}
	})

	t.Run("command with similar name should not fail validation", func(t *testing.T) {
		testCfg := &config.Config{
			Commands: make(map[string]config.Command),
		}
		testCfg.Commands["a-cmd"] = config.Command{
			Name: "a-cmd",
			Depends: map[string]config.Dep{
				"b1-cmd": {Name: "b1-cmd"},
			},
		}
		testCfg.Commands["b"] = config.Command{
			Name: "b",
			Depends: map[string]config.Dep{
				"a-cmd": {Name: "a-cmd"},
			},
		}
		testCfg.Commands["b1-cmd"] = config.Command{
			Name:    "b1-cmd",
			Depends: map[string]config.Dep{},
		}
		err := validateCircularDepends(testCfg)
		if err != nil {
			t.Errorf("checked itself when validation circular depends. got:  %s", err)
		}
	})

	t.Run("validation should fail", func(t *testing.T) {
		testCfg := &config.Config{
			Commands: make(map[string]config.Command),
		}
		testCfg.Commands["a-cmd"] = config.Command{
			Name: "a-cmd",
			Depends: map[string]config.Dep{
				"b1-cmd": {Name: "b1-cmd"},
			},
		}
		testCfg.Commands["b"] = config.Command{
			Name: "b",
			Depends: map[string]config.Dep{
				"a-cmd": {Name: "a-cmd"},
			},
		}
		testCfg.Commands["b1-cmd"] = config.Command{
			Name: "b1-cmd",
			Depends: map[string]config.Dep{
				"a-cmd": {Name: "a-cmd"},
			},
		}
		err := validateCircularDepends(testCfg)

		if err == nil {
			t.Errorf("validation should fail. got: %s", err)
		}
	})
}
