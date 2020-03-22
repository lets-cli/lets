package test

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type SerializableTestConfig struct {
	Shell    string
	Commands map[string]map[string]string
}

// NewTestConfig creates config, write it to temp dir, set LETS_CONFIG_DIR with LETS_CONFIG and return cleanup func
func NewTestConfig(configRaw *SerializableTestConfig) func() {
	tempDir := os.TempDir()
	testConfigFile := CreateTempFile(tempDir, "lets_*.yaml")

	err := yaml.NewEncoder(testConfigFile).Encode(configRaw)
	if err != nil {
		log.Fatalf("can not create test config: %s", err)
	}

	err = os.Setenv("LETS_CONFIG_DIR", tempDir)
	if err != nil {
		log.Fatalf("can not set LETS_CONFIG_DIR during test: %s", err)
	}

	err = os.Setenv("LETS_CONFIG", testConfigFile.Name())
	if err != nil {
		log.Fatalf("can not set LETS_CONFIG during test: %s", err)
	}

	return func() {
		err := os.Unsetenv("LETS_CONFIG_DIR")
		if err != nil {
			log.Fatalf("can not unset LETS_CONFIG_DIR after test: %s", err)
		}

		err = os.Unsetenv("LETS_CONFIG")
		if err != nil {
			log.Fatalf("can not unset LETS_CONFIG after test: %s", err)
		}

		err = os.Remove(testConfigFile.Name())
		if err != nil {
			log.Fatalf("can not remove temp config file %s after test: %s", testConfigFile.Name(), err)
		}
	}
}
