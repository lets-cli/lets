package upgrade

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/lets-cli/lets/internal/upgrade/registry"
)

type mockNotifierRegistry struct {
	release *registry.ReleaseInfo
	calls   int
}

func (m *mockNotifierRegistry) GetLatestReleaseInfo(ctx context.Context) (*registry.ReleaseInfo, error) {
	m.calls++
	return m.release, nil
}

func (m *mockNotifierRegistry) GetLatestRelease() (string, error) {
	m.calls++
	if m.release == nil {
		return "", nil
	}

	return m.release.TagName, nil
}

func (m *mockNotifierRegistry) DownloadReleaseBinary(packageName string, version string, dstPath string) error {
	return nil
}

func (m *mockNotifierRegistry) GetPackageName(os string, arch string) (string, error) {
	return "", nil
}

func (m *mockNotifierRegistry) GetDownloadURL(repoURI string, packageName string, version string) string {
	return ""
}

func TestUpdateNotifierCheck(t *testing.T) {
	now := time.Date(2026, 3, 19, 12, 0, 0, 0, time.UTC)

	t.Run("should use cached state without network call", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.yaml")
		reg := &mockNotifierRegistry{
			release: &registry.ReleaseInfo{
				TagName:     "v0.0.9",
				PublishedAt: now.Add(-48 * time.Hour),
			},
		}

		notifier := newUpdateNotifier(reg, statePath, "/usr/local/bin/lets", func() time.Time { return now })
		err := notifier.writeState(notifierState{
			CheckedAt:         now.Add(-time.Hour),
			LatestVersion:     "v0.0.9",
			LatestPublishedAt: now.Add(-48 * time.Hour),
		})
		if err != nil {
			t.Fatalf("writeState() error = %v", err)
		}

		notice, err := notifier.Check(context.Background(), "0.0.8")
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		if notice == nil {
			t.Fatal("expected cached notice")
		}
		if reg.calls != 0 {
			t.Fatalf("expected no network calls, got %d", reg.calls)
		}
	})

	t.Run("should skip dev builds", func(t *testing.T) {
		tmpDir := t.TempDir()
		reg := &mockNotifierRegistry{
			release: &registry.ReleaseInfo{
				TagName:     "v0.0.9",
				PublishedAt: now.Add(-48 * time.Hour),
			},
		}

		notifier := newUpdateNotifier(reg, filepath.Join(tmpDir, "state.yaml"), "/usr/local/bin/lets", func() time.Time { return now })

		notice, err := notifier.Check(context.Background(), "0.0.8-dev")
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		if notice != nil {
			t.Fatal("expected no notice for dev build")
		}
		if reg.calls != 0 {
			t.Fatalf("expected no registry calls, got %d", reg.calls)
		}
	})

	t.Run("should persist latest release after successful check", func(t *testing.T) {
		tmpDir := t.TempDir()
		statePath := filepath.Join(tmpDir, "state.yaml")
		reg := &mockNotifierRegistry{
			release: &registry.ReleaseInfo{
				TagName:     "v0.0.9",
				PublishedAt: now.Add(-48 * time.Hour),
			},
		}

		notifier := newUpdateNotifier(reg, statePath, "/usr/local/bin/lets", func() time.Time { return now })

		notice, err := notifier.Check(context.Background(), "0.0.8")
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		if notice == nil {
			t.Fatal("expected update notice")
		}

		state, err := notifier.readState()
		if err != nil {
			t.Fatalf("readState() error = %v", err)
		}
		if state.LatestVersion != "v0.0.9" {
			t.Fatalf("expected latest version to be persisted, got %q", state.LatestVersion)
		}
		if state.CheckedAt != now {
			t.Fatalf("expected checkedAt %s, got %s", now, state.CheckedAt)
		}
	})

	t.Run("should suppress homebrew notices until release ages", func(t *testing.T) {
		tmpDir := t.TempDir()
		reg := &mockNotifierRegistry{
			release: &registry.ReleaseInfo{
				TagName:     "v0.0.9",
				PublishedAt: now.Add(-2 * time.Hour),
			},
		}

		notifier := newUpdateNotifier(
			reg,
			filepath.Join(tmpDir, "state.yaml"),
			"/opt/homebrew/Cellar/lets/0.0.8/bin/lets",
			func() time.Time { return now },
		)

		notice, err := notifier.Check(context.Background(), "0.0.8")
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		if notice != nil {
			t.Fatal("expected no notice during homebrew delay window")
		}
	})

	t.Run("should suppress repeated notices after mark notified", func(t *testing.T) {
		tmpDir := t.TempDir()
		reg := &mockNotifierRegistry{
			release: &registry.ReleaseInfo{
				TagName:     "v0.0.9",
				PublishedAt: now.Add(-48 * time.Hour),
			},
		}

		notifier := newUpdateNotifier(reg, filepath.Join(tmpDir, "state.yaml"), "/usr/local/bin/lets", func() time.Time { return now })

		notice, err := notifier.Check(context.Background(), "0.0.8")
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		if notice == nil {
			t.Fatal("expected update notice")
		}

		if err := notifier.MarkNotified(notice); err != nil {
			t.Fatalf("MarkNotified() error = %v", err)
		}

		secondNotice, err := notifier.Check(context.Background(), "0.0.8")
		if err != nil {
			t.Fatalf("second Check() error = %v", err)
		}
		if secondNotice != nil {
			t.Fatal("expected repeated notice to be suppressed")
		}
	})
}

func TestLetsStatePath(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	t.Setenv("HOME", tmpDir)

	path, err := letsStatePath()
	if err != nil {
		t.Fatalf("letsStatePath() error = %v", err)
	}
	if filepath.Base(path) != "state.yaml" {
		t.Fatalf("unexpected state file %q", path)
	}
	if filepath.Base(filepath.Dir(path)) != "lets" {
		t.Fatalf("unexpected state dir %q", filepath.Dir(path))
	}
	if !filepath.IsAbs(path) {
		t.Fatalf("expected absolute state path, got %q", path)
	}
}

func TestIsHomebrewInstall(t *testing.T) {
	if !isHomebrewInstall("/opt/homebrew/Cellar/lets/0.0.1/bin/lets") {
		t.Fatal("expected homebrew path to be detected")
	}
	if isHomebrewInstall("/usr/local/bin/lets") {
		t.Fatal("did not expect generic install path to be detected as homebrew")
	}
}
