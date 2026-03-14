package config

import "testing"

func TestEnvsExecute(t *testing.T) {
	cfg := Config{
		Shell:   "bash",
		WorkDir: ".",
	}

	t.Run("resolves env entries sequentially", func(t *testing.T) {
		envs := &Envs{}
		envs.Set("ENGINE", Env{Name: "ENGINE", Value: "docker"})
		envs.Set("COMPOSE", Env{Name: "COMPOSE", Sh: `echo "${ENGINE}-compose"`})

		err := envs.Execute(cfg, nil)
		if err != nil {
			t.Fatalf("unexpected execute error: %s", err)
		}

		if got := envs.Mapping["COMPOSE"].Value; got != "docker-compose" {
			t.Fatalf("expected COMPOSE=docker-compose, got %q", got)
		}
	})

	t.Run("uses base env for sh evaluation", func(t *testing.T) {
		envs := &Envs{}
		envs.Set("COMPOSE", Env{Name: "COMPOSE", Sh: `echo "${ENGINE}-compose"`})

		err := envs.Execute(cfg, map[string]string{"ENGINE": "docker"})
		if err != nil {
			t.Fatalf("unexpected execute error: %s", err)
		}

		if got := envs.Mapping["COMPOSE"].Value; got != "docker-compose" {
			t.Fatalf("expected COMPOSE=docker-compose, got %q", got)
		}
	})

	t.Run("resolved lets env overrides process env", func(t *testing.T) {
		t.Setenv("ENGINE", "podman")

		envs := &Envs{}
		envs.Set("ENGINE", Env{Name: "ENGINE", Value: "docker"})
		envs.Set("COMPOSE", Env{Name: "COMPOSE", Sh: `echo "${ENGINE}-compose"`})

		err := envs.Execute(cfg, nil)
		if err != nil {
			t.Fatalf("unexpected execute error: %s", err)
		}

		if got := envs.Mapping["COMPOSE"].Value; got != "docker-compose" {
			t.Fatalf("expected COMPOSE=docker-compose, got %q", got)
		}
	})

	t.Run("keeps cached values after first execution", func(t *testing.T) {
		envs := &Envs{}
		envs.Set("COMPOSE", Env{Name: "COMPOSE", Sh: `echo "${ENGINE}-compose"`})

		err := envs.Execute(cfg, map[string]string{"ENGINE": "docker"})
		if err != nil {
			t.Fatalf("unexpected execute error: %s", err)
		}

		err = envs.Execute(cfg, map[string]string{"ENGINE": "podman"})
		if err != nil {
			t.Fatalf("unexpected execute error: %s", err)
		}

		if got := envs.Mapping["COMPOSE"].Value; got != "docker-compose" {
			t.Fatalf("expected cached COMPOSE=docker-compose, got %q", got)
		}
	})
}
