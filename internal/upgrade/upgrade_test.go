package upgrade

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/lets-cli/lets/internal/upgrade/registry"
)

type MockRegistry struct {
	latestVersion string
}

func (m MockRegistry) GetLatestRelease(ctx context.Context) (string, error) {
	return m.latestVersion, nil
}

func (m MockRegistry) GetLatestReleaseInfo(ctx context.Context) (*registry.ReleaseInfo, error) {
	return &registry.ReleaseInfo{TagName: m.latestVersion}, nil
}

func (m MockRegistry) DownloadReleaseBinary(ctx context.Context, packageName string, version string, dstPath string) error {
	file, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	latest, _ := m.GetLatestRelease(ctx)

	_, err = fmt.Fprint(file, latest)
	if err != nil {
		return err
	}

	return nil
}

func (m MockRegistry) GetPackageName(os string, arch string) (string, error) {
	return "lets_test_package", nil
}

func (m MockRegistry) GetDownloadURL(repoURI string, packageName string, version string) string {
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
		registry:       reg,
		currentVersion: version,
		binaryPath:     binary.Name(),
		downloadPath:   path.Join(os.TempDir(), "lets.download"),
		backupPath:     path.Join(os.TempDir(), "lets.backup"),
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

		err = upgrader.Upgrade(context.Background())
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
		time.Sleep(10 * time.Millisecond)

		err = upgrader.Upgrade(context.Background())
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

	t.Run("should self-upgrade symlink target", func(t *testing.T) {
		currentVersion := "v0.0.1"
		latestVersion := "v0.0.2"

		tempDir := t.TempDir()
		targetPath := path.Join(tempDir, ".lets", "bin", "lets")
		symlinkPath := path.Join(tempDir, ".local", "bin", "lets")
		if err := os.MkdirAll(path.Dir(targetPath), 0o755); err != nil {
			t.Fatalf("failed to create target dir: %s", err)
		}
		if err := os.MkdirAll(path.Dir(symlinkPath), 0o755); err != nil {
			t.Fatalf("failed to create symlink dir: %s", err)
		}
		if err := os.WriteFile(targetPath, []byte(currentVersion), 0o755); err != nil {
			t.Fatalf("failed to write target binary: %s", err)
		}
		if err := os.Symlink(targetPath, symlinkPath); err != nil {
			t.Fatalf("failed to create binary symlink: %s", err)
		}

		upgrader := &BinaryUpgrader{
			registry:       &MockRegistry{latestVersion: latestVersion},
			currentVersion: currentVersion,
			binaryPath:     symlinkPath,
			downloadPath:   path.Join(tempDir, "lets.download"),
			backupPath:     path.Join(tempDir, "lets.backup"),
		}

		err := upgrader.Upgrade(context.Background())
		if err != nil {
			t.Fatalf("failed to upgrade symlink target: %s", err)
		}

		if !testVersion(targetPath, latestVersion) {
			t.Errorf("expected target version %s", latestVersion)
		}

		if linkTarget, err := os.Readlink(symlinkPath); err != nil || linkTarget != targetPath {
			t.Fatalf("expected symlink to remain pointed at %s, got %q, err %v", targetPath, linkTarget, err)
		}
	})

	t.Run("should not self-upgrade homebrew-managed binary", func(t *testing.T) {
		currentVersion := "v0.0.1"
		latestVersion := "v0.0.2"

		tempDir := t.TempDir()
		binaryPath := path.Join(tempDir, "Cellar", "lets", currentVersion, "bin", "lets")
		if err := os.MkdirAll(path.Dir(binaryPath), 0o755); err != nil {
			t.Fatalf("failed to create homebrew binary dir: %s", err)
		}

		if err := os.WriteFile(binaryPath, []byte(currentVersion), 0o755); err != nil {
			t.Fatalf("failed to write homebrew binary: %s", err)
		}

		upgrader := &BinaryUpgrader{
			registry:       &MockRegistry{latestVersion: latestVersion},
			currentVersion: currentVersion,
			binaryPath:     binaryPath,
			downloadPath:   path.Join(tempDir, "lets.download"),
			backupPath:     path.Join(tempDir, "lets.backup"),
		}

		err := upgrader.Upgrade(context.Background())
		if err == nil {
			t.Fatal("expected homebrew upgrade error")
		}

		if !strings.Contains(err.Error(), "brew upgrade lets-cli/tap/lets") {
			t.Fatalf("expected homebrew upgrade command in error, got %q", err.Error())
		}

		if !testVersion(binaryPath, currentVersion) {
			t.Errorf("expected version %s", currentVersion)
		}

		if _, err := os.Stat(upgrader.downloadPath); !os.IsNotExist(err) {
			t.Fatalf("expected no downloaded binary, got err %v", err)
		}
	})

	t.Run("should not self-upgrade system-managed binary", func(t *testing.T) {
		currentVersion := "v0.0.1"
		latestVersion := "v0.0.2"

		tempDir := t.TempDir()
		upgrader := &BinaryUpgrader{
			registry:       &MockRegistry{latestVersion: latestVersion},
			currentVersion: currentVersion,
			binaryPath:     "/usr/bin/lets",
			downloadPath:   path.Join(tempDir, "lets.download"),
			backupPath:     path.Join(tempDir, "lets.backup"),
		}

		err := upgrader.Upgrade(context.Background())
		if err == nil {
			t.Fatal("expected system-managed upgrade error")
		}

		if !strings.Contains(err.Error(), "system-managed lets install") {
			t.Fatalf("expected system-managed error, got %q", err.Error())
		}

		if _, err := os.Stat(upgrader.downloadPath); !os.IsNotExist(err) {
			t.Fatalf("expected no downloaded binary, got err %v", err)
		}
	})
}
