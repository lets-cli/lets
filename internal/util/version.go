package util

import (
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

func ParseVersion(version string) (*semver.Version, error) {
	version = strings.TrimPrefix(version, "v")

	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, fmt.Errorf("can not create semver version from %s: %w", version, err)
	}

	return v, nil
}
