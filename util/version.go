package util

import (
	"github.com/coreos/go-semver/semver"
)

func ParseVersion(version string) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	return v, nil
}
