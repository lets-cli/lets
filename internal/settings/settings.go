package settings

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lets-cli/lets/internal/theme"
	"github.com/lets-cli/lets/internal/util"
	"gopkg.in/yaml.v3"
)

type FileSettings struct {
	NoColor       *bool   `yaml:"no_color"`
	Theme         *string `yaml:"theme"`
	UpgradeNotify *bool   `yaml:"upgrade_notify"`
}

type Settings struct {
	NoColor       bool
	Theme         string
	UpgradeNotify bool
}

func Default() Settings {
	return Settings{
		NoColor:       false,
		Theme:         theme.DefaultName,
		UpgradeNotify: true,
	}
}

func Load() (Settings, error) {
	path, err := util.LetsUserFile("config.yaml")
	if err != nil {
		return Settings{}, err
	}

	return LoadFile(path)
}

func LoadFile(path string) (Settings, error) {
	cfg := Default()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			applyEnvOverrides(&cfg)
			return cfg, nil
		}

		return Settings{}, fmt.Errorf("failed to open settings file: %w", err)
	}

	defer file.Close()

	var fileSettings FileSettings

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	if err := decoder.Decode(&fileSettings); err != nil {
		return Settings{}, fmt.Errorf("failed to decode settings file: %w", err)
	}

	if fileSettings.NoColor != nil {
		cfg.NoColor = *fileSettings.NoColor
	}

	if fileSettings.Theme != nil {
		cfg.Theme = *fileSettings.Theme
	}

	if fileSettings.UpgradeNotify != nil {
		cfg.UpgradeNotify = *fileSettings.UpgradeNotify
	}

	applyEnvOverrides(&cfg)

	if !theme.ValidName(cfg.Theme) {
		return Settings{}, fmt.Errorf(
			"invalid theme %q: must be one of %q, %q, %q",
			cfg.Theme,
			theme.DefaultName,
			theme.ANSIName,
			theme.SynthwaveName,
		)
	}

	return cfg, nil
}

func (s Settings) Apply() {
	if s.NoColor {
		color.NoColor = true
	}
}

func applyEnvOverrides(cfg *Settings) {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		cfg.NoColor = true
	}

	if _, ok := os.LookupEnv("LETS_CHECK_UPDATE"); ok {
		cfg.UpgradeNotify = false
	}
}
