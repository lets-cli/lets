package upgrade

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/lets-cli/lets/internal/upgrade/registry"
	"github.com/lets-cli/lets/internal/util"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	updateCheckInterval  = 24 * time.Hour
	updateNotifyInterval = 24 * time.Hour
	homebrewNoticeDelay  = 24 * time.Hour
)

type UpdateNotice struct {
	CurrentVersion string
	LatestVersion  string
	command        string
}

func (n *UpdateNotice) Message() string {
	return fmt.Sprintf(
		"\n%s: %s -> %s\n%s",
		color.YellowString("lets: new version been released"),
		color.RedString(n.CurrentVersion),
		color.GreenString(n.LatestVersion),
		color.YellowString("Run '%s' or see https://lets-cli.org/docs/installation", n.command),
	)
}

type notifierState struct {
	CheckedAt         time.Time `yaml:"checked_at"`
	LatestVersion     string    `yaml:"latest_version"`
	LatestPublishedAt time.Time `yaml:"latest_published_at"`
	NotifiedAt        time.Time `yaml:"notified_at"`
}

type UpdateNotifier struct {
	registry       registry.RepoRegistry
	statePath      string
	executablePath string
	now            func() time.Time
}

func NewUpdateNotifier(reg registry.RepoRegistry) (*UpdateNotifier, error) {
	statePath, err := letsStatePath()
	if err != nil {
		return nil, err
	}

	executablePath, err := binaryPath()
	if err != nil {
		return nil, err
	}

	return newUpdateNotifier(reg, statePath, executablePath, time.Now), nil
}

func newUpdateNotifier(
	reg registry.RepoRegistry,
	statePath string,
	executablePath string,
	now func() time.Time,
) *UpdateNotifier {
	return &UpdateNotifier{
		registry:       reg,
		statePath:      statePath,
		executablePath: executablePath,
		now:            now,
	}
}

func (n *UpdateNotifier) Check(ctx context.Context, currentVersion string) (*UpdateNotice, error) {
	current, ok := parseStableVersion(currentVersion)
	if !ok {
		return nil, nil
	}

	state, err := n.readState()
	if err != nil {
		return nil, err
	}

	now := n.now()
	if now.Sub(state.CheckedAt) < updateCheckInterval {
		log.Debugf("lets: skip update check: next check at %s", state.CheckedAt.Add(updateCheckInterval))
		return n.noticeFromState(state, currentVersion, current, now), nil
	}

	release, err := n.registry.GetLatestReleaseInfo(ctx)
	if err != nil {
		return n.noticeFromState(state, currentVersion, current, now), err
	}

	state.CheckedAt = now
	state.LatestVersion = release.TagName
	state.LatestPublishedAt = release.PublishedAt

	if err := n.writeState(state); err != nil {
		return nil, err
	}

	return n.noticeFromState(state, currentVersion, current, now), nil
}

func (n *UpdateNotifier) MarkNotified(notice *UpdateNotice) error {
	if notice == nil {
		return nil
	}

	state, err := n.readState()
	if err != nil {
		return err
	}

	if state.LatestVersion != notice.LatestVersion {
		return nil
	}

	state.NotifiedAt = n.now()

	return n.writeState(state)
}

func (n *UpdateNotifier) noticeFromState(
	state notifierState,
	currentVersion string,
	current *semver.Version,
	now time.Time,
) *UpdateNotice {
	latest, ok := parseStableVersion(state.LatestVersion)
	if !ok {
		return nil
	}

	if !current.LessThan(*latest) {
		return nil
	}

	if now.Sub(state.NotifiedAt) < updateNotifyInterval {
		return nil
	}

	command := "lets self upgrade"
	if isHomebrewInstall(n.executablePath) {
		if !state.LatestPublishedAt.IsZero() && now.Sub(state.LatestPublishedAt) < homebrewNoticeDelay {
			return nil
		}

		command = "brew upgrade lets-cli/tap/lets"
	}

	return &UpdateNotice{
		CurrentVersion: currentVersion,
		LatestVersion:  state.LatestVersion,
		command:        command,
	}
}

func (n *UpdateNotifier) readState() (notifierState, error) {
	var state notifierState

	file, err := os.Open(n.statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return state, nil
		}

		return state, fmt.Errorf("failed to open update state file: %w", err)
	}

	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&state); err != nil {
		return notifierState{}, fmt.Errorf("failed to decode update state file: %w", err)
	}

	return state, nil
}

func (n *UpdateNotifier) writeState(state notifierState) error {
	dir := filepath.Dir(n.statePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create update state dir: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, "state.*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create update state temp file: %w", err)
	}

	tmpPath := tmpFile.Name()

	defer os.Remove(tmpPath)

	if err := yaml.NewEncoder(tmpFile).Encode(state); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to encode update state file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close update state temp file: %w", err)
	}

	if err := os.Rename(tmpPath, n.statePath); err != nil {
		return fmt.Errorf("failed to replace update state file: %w", err)
	}

	return nil
}

func letsStatePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}

	return filepath.Join(homeDir, ".config", "lets", "state.yaml"), nil
}

func parseStableVersion(version string) (*semver.Version, bool) {
	parsed, err := util.ParseVersion(version)
	if err != nil {
		return nil, false
	}

	if parsed.PreRelease != "" {
		return nil, false
	}

	return parsed, true
}

func isHomebrewInstall(path string) bool {
	return strings.Contains(path, "/Cellar/lets/")
}

func LogUpdateCheckError(err error) {
	if err == nil {
		return
	}

	log.Debugf("lets: update notifier error: %s", err)
}
