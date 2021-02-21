package upgrade

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/lets-cli/lets/upgrade/registry"
)

type MockRegistry struct {
	latestVersion string
}

func (m MockRegistry) GetLatestRelease() (string, error) {
	return m.latestVersion, nil
}

func (m MockRegistry) DownloadReleaseBinary(packageName string, version string, dstPath string) error {
	file, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	latest, _ := m.GetLatestRelease()

	_, err = fmt.Fprint(file, latest)
	if err != nil {
		return err
	}

	return nil
}

func (m MockRegistry) GetPackageName(os string, arch string) (string, error) {
	return "lets_test_package", nil
}

func (m MockRegistry) GetDownloadUrl(repoUri string, packageName string, version string) string {
	return ""
}

func createTempBinary(content string) (*os.File, error) {
	binary, err := os.CreateTemp("", "lets.*.current")
	if err != nil {
		return nil, err
	}

	_, err = fmt.Fprint(binary, content)
	if err != nil {
		return nil, err
	}

	return binary, err
}


func newMockUpgrader(reg registry.RepoRegistry, version string) (*BinaryUpgrader, error) {
	binary, err := createTempBinary(version)
	if err != nil {
		return nil, err
	}

	return &BinaryUpgrader{
		registry: reg,
		currentVersion: version,
		binaryPath: binary.Name(),
		downloadPath: path.Join(os.TempDir(), "lets.download"),
		backupPath: path.Join(os.TempDir(), "lets.backup"),
	}, nil
}

func testVersion(filePath string, version string) bool {
	data, _ := os.ReadFile(filePath)
	return string(data) == version
}

func getFileModTime(filePath string) time.Time {
	file, _ := os.Open(filePath)
	stats, _ := file.Stat()
	return stats.ModTime()
}

func TestSelfUpgrade(t *testing.T) {
	t.Run("should self-upgrade to latest version", func(t *testing.T) {
		currentVersion := "v0.0.1"
		latestVersion := "v0.0.2"

		upgrader, err := newMockUpgrader(&MockRegistry{latestVersion: latestVersion}, currentVersion)
		if err != nil {
			t.Errorf("failed to create upgrader: %s", err)
		}

		if !testVersion(upgrader.binaryPath, currentVersion) {
			t.Errorf("expected version %s", currentVersion)
		}

		err = upgrader.Upgrade()
		if err != nil {
			t.Errorf("failed to upgrade: %s", err)
		}

		if !testVersion(upgrader.binaryPath, latestVersion) {
			t.Errorf("expected version %s", latestVersion)
		}
	})

	t.Run("should not self-upgrade same version", func(t *testing.T) {
		currentVersion := "v0.0.1"
		latestVersion := "v0.0.1"

		upgrader, err := newMockUpgrader(&MockRegistry{latestVersion: latestVersion}, currentVersion)
		if err != nil {
			t.Errorf("failed to create upgrader: %s", err)
		}

		if !testVersion(upgrader.binaryPath, currentVersion) {
			t.Errorf("expected version %s", currentVersion)
		}

		binaryModTime := getFileModTime(upgrader.binaryPath)

		// sleep to be sure files not created at the same time
		time.Sleep(10*time.Millisecond)

		err = upgrader.Upgrade()
		if err != nil {
			t.Errorf("failed to upgrade: %s", err)
		}

		if !testVersion(upgrader.binaryPath, currentVersion) {
			t.Errorf("expected version %s", currentVersion)
		}

		binaryUpdatedModTime := getFileModTime(upgrader.binaryPath)

		if binaryModTime != binaryUpdatedModTime {
			t.Errorf("binary must not been updated")
		}
	})
}
