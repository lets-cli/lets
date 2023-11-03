package upgrade

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/lets-cli/lets/upgrade/registry"
	log "github.com/sirupsen/logrus"
)

type Upgrader interface {
	Upgrade() error
}

type BinaryUpgrader struct {
	registry       registry.RepoRegistry
	currentVersion string
	binaryPath     string
	downloadPath   string
	backupPath     string
}

func NewBinaryUpgrader(reg registry.RepoRegistry, currentVersion string) (*BinaryUpgrader, error) {
	executablePath, err := binaryPath()
	if err != nil {
		return nil, err
	}

	return &BinaryUpgrader{
		registry:       reg,
		currentVersion: currentVersion,
		// TODO rewrite all paths with home dir
		binaryPath:   executablePath,
		downloadPath: path.Join(os.TempDir(), "lets.download"),
		backupPath:   path.Join(os.TempDir(), "lets.backup"),
	}, nil
}

func (up *BinaryUpgrader) Upgrade() error {
	latestVersion, err := up.registry.GetLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get latest release version: %w", err)
	}

	if up.currentVersion == latestVersion {
		log.Printf("Lets is up-to-date")

		return nil
	}

	packageName, err := up.registry.GetPackageName(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("failed to get package name: %w", err)
	}

	log.Printf("Downloading latest release %s...", latestVersion)

	err = up.registry.DownloadReleaseBinary(
		packageName,
		latestVersion,
		up.downloadPath,
	)
	if err != nil {
		return fmt.Errorf("failed to download release %s version %s: %w", packageName, latestVersion, err)
	}

	err = backupExecutable(up.binaryPath, up.backupPath)
	if err != nil {
		return err
	}

	err = replaceBinaries(up.downloadPath, up.binaryPath, up.backupPath)
	if err != nil {
		return err
	}

	log.Printf("Upgraded to version %s", latestVersion)

	return nil
}

func binaryPath() (string, error) {
	// TODO after implementing $HOME/.lets/bin, deny upgrading in other places
	return os.Executable()
}

func backupExecutable(executablePath string, backupPath string) error {
	errFmt := func(err error) error {
		return fmt.Errorf("failed to backup current lets binary: %w", err)
	}

	executableFile, err := os.Open(executablePath)
	if err != nil {
		return errFmt(err)
	}

	err = os.RemoveAll(backupPath)
	if err != nil {
		return errFmt(err)
	}

	backupFile, err := os.Create(backupPath)
	if err != nil {
		return errFmt(err)
	}

	_, err = io.Copy(backupFile, executableFile)
	if err != nil {
		return errFmt(err)
	}

	return nil
}

func replaceBinaries(downloadPath string, executablePath string, backupPath string) error {
	defer os.RemoveAll(downloadPath)
	defer os.RemoveAll(backupPath)

	newBinary, err := os.Open(downloadPath)
	defer newBinary.Close()

	if err != nil {
		return fmt.Errorf("failed to open new lets binary: %w", err)
	}

	currentBinary, err := os.OpenFile(executablePath, os.O_WRONLY, 0755)
	defer currentBinary.Close()

	if err != nil {
		return fmt.Errorf("failed to open current lets binary: %w", err)
	}

	_, err = io.Copy(currentBinary, newBinary)
	if err != nil {
		backupBinary, err := os.Open(backupPath)
		defer backupBinary.Close()

		if err != nil {
			return fmt.Errorf("failed to open backup lets binary: %w", err)
		}

		// restore original file from backup
		_, err = io.Copy(currentBinary, backupBinary)
		return fmt.Errorf("failed to update lets binary: %w", err)
	}

	return nil
}
