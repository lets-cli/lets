package util

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
)

func ParseVersion(version string) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, fmt.Errorf("can not create semver version from %s: %w", version, err)
	}

	return v, nil
}
