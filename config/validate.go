package config

import (
	"fmt"
	"strings"
)

const (
	NoticeColor = "\033[1;36m%s\033[0m"
)

// Validate loaded config
func Validate(config *Config) error {
	return validateCircularDepends(config)
}

func validateCircularDepends(cfg *Config) error {
	for _, cmdA := range cfg.Commands {
		for _, cmdB := range cfg.Commands {
			depsA := strings.Join(cmdA.Depends, " ")
			depsB := strings.Join(cmdB.Depends, " ")
			if strings.Contains(depsB, cmdA.Name) &&
				strings.Contains(depsA, cmdB.Name) {
				return fmt.Errorf(
					"command '%s' have circular depends on command '%s'",
					fmt.Sprintf(NoticeColor, cmdA.Name),
					fmt.Sprintf(NoticeColor, cmdB.Name),
				)
			}
		}
	}
	return nil
}

func validateTopLevelFields(rawKeyValue map[string]interface{}, validFields string) error {
	for k := range rawKeyValue {
		if !strings.Contains(validFields, k) {
			return fmt.Errorf("unknown top-level field '%s'", k)
		}
	}
	return nil
}
