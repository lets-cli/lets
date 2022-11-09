package config

import (
	"testing"

	"github.com/lets-cli/lets/config/config"
)

func TestValidateCommandInDependsExists(t *testing.T) {
	t.Run("command depends on non-existing command", func(t *testing.T) {
		testCfg := &config.Config{
			Commands: make(map[string]*config.Command),
		}
		deps := &config.Deps{}
		deps.Set("bar", config.Dep{Name: "bar"})
		testCfg.Commands["foo"] = &config.Command{
			Name:    "foo",
			Depends: deps,
		}
		err := validateDepends(testCfg)
		if err == nil {
			t.Error("command foo depends on non-existing command bar. Must fail")
		}
	})
}

func TestValidateCircularDeps(t *testing.T) {
	t.Run("command skip itself", func(t *testing.T) {
		testCfg := &config.Config{
			Commands: make(map[string]*config.Command),
		}
		depsA := &config.Deps{}
		testCfg.Commands["a"] = &config.Command{
			Name:    "a",
			Depends: depsA,
		}

		depsB := &config.Deps{}
		testCfg.Commands["b"] = &config.Command{
			Name:    "b",
			Depends: depsB,
		}

		err := validateDepends(testCfg)
		if err != nil {
			t.Errorf("checked itself when validation circular depends. got:  %s", err)
		}
	})

	t.Run("command with similar name should not fail validation", func(t *testing.T) {
		testCfg := &config.Config{
			Commands: make(map[string]*config.Command),
		}
		depsA := &config.Deps{}
		depsA.Set("b1", config.Dep{Name: "b1"})
		testCfg.Commands["a"] = &config.Command{
			Name:    "a",
			Depends: depsA,
		}

		depsB := &config.Deps{}
		depsB.Set("a", config.Dep{Name: "a"})
		testCfg.Commands["b"] = &config.Command{
			Name:    "b",
			Depends: depsB,
		}

		depsB1 := &config.Deps{}
		testCfg.Commands["b1"] = &config.Command{
			Name:    "b1",
			Depends: depsB1,
		}

		err := validateDepends(testCfg)
		if err != nil {
			t.Errorf("checked itself when validation circular depends. got:  %s", err)
		}
	})

	t.Run("validation should fail", func(t *testing.T) {
		testCfg := &config.Config{
			Commands: make(map[string]*config.Command),
		}
		depsA := &config.Deps{}
		depsA.Set("b", config.Dep{Name: "b"})
		testCfg.Commands["a"] = &config.Command{
			Name:    "a",
			Depends: depsA,
		}

		depsB := &config.Deps{}
		depsB.Set("a", config.Dep{Name: "a"})
		testCfg.Commands["b"] = &config.Command{
			Name:    "b",
			Depends: depsB,
		}

		err := validateDepends(testCfg)

		if err == nil {
			t.Errorf("validation should fail. got: %s", err)
		}
	})
}
