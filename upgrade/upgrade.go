package upgrade

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/lets-cli/lets/logging"
	"github.com/lets-cli/lets/upgrade/registry"
)

func Upgrade(reg registry.RepoRegistry) error {
	executablePath, err := binaryPath()
	if err != nil {
		return err
	}

	latestVersion, err := reg.GetLatestRelease()
	if err != nil {
		return err
	}

	packageName, err := reg.GetPackageName(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	// TODO rewrite with home dir
	downloadPath := path.Join(os.TempDir(), fmt.Sprintf("%s_%s", packageName, latestVersion))

	logging.Log.Printf("Downloading latest release %s...", latestVersion)

	err = reg.DownloadReleaseBinary(
		packageName,
		latestVersion,
		downloadPath,
	)
	if err != nil {
		return err
	}

	backupPath := path.Join(os.TempDir(), "lets.backup")

	err = backupExecutable(executablePath, backupPath)
	if err != nil {
		return err
	}

	err = replaceBinaries(downloadPath, executablePath, backupPath)
	if err != nil {
		return err
	}

	logging.Log.Printf("Upgraded to version %s", latestVersion)
	return nil
}

func binaryPath() (string, error) {
	// TODO after implementing $HOME/.lets/bin, deny upgrading in other places
	return os.Executable()
}

func backupExecutable(executablePath string, backupPath string) error {
	errFmt := func(err error) error {
		return fmt.Errorf("failed to backup current lets binary: %s", err)
	}

	executableFile, err := os.Open(executablePath)
	if err != nil {
		return errFmt(err)
	}

	// TODO remove old backup
	// TODO in edge db there is a hard link from original executable to backup path
	backupFile, err := os.Create(backupPath)
	if err != nil {
		return errFmt(err)
	}

	// in edgedb there is hard_link
	_, err = io.Copy(backupFile, executableFile)
	if err != nil {
		return errFmt(err)
	}

	return nil
}

func replaceBinaries(downloadPath string, executablePath string, backupPath string) error {
	defer os.RemoveAll(downloadPath)
	defer os.RemoveAll(backupPath)

	newExecutablePath := path.Join(downloadPath, "lets")
	err := os.Rename(newExecutablePath, executablePath)
	if err != nil {
		// TODO handle backupPath
		return fmt.Errorf("failed to update lets binary: %s", err)
	}

	return nil
}
