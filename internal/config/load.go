package config

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lets-cli/lets/internal/config/config"
	"github.com/lets-cli/lets/internal/fetch"
	"github.com/lets-cli/lets/internal/util"
	"github.com/lets-cli/lets/internal/workdir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func Load(configName string, configDir string, version string) (*config.Config, error) {
	configPath, err := FindConfig(configName, configDir)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(configPath.AbsPath)
	if err != nil {
		return nil, err
	}

	c := config.NewConfig(
		configPath.WorkDir,
		configPath.AbsPath,
		configPath.DotLetsDir,
	)
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", configPath.Filename, err)
	}

	if err = validate(c, version); err != nil {
		return nil, err
	}

	if err := c.SetupEnv(); err != nil {
		return nil, err
	}

	return c, nil
}

// LoadRemote downloads (or loads from cache) a remote lets.yaml at url and
// returns a Config with the working directory set to the caller's CWD.
func LoadRemote(url string, noCache bool, version string) (*config.Config, error) {
	cachedPath, err := ensureRemoteConfig(url, noCache)
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	dotLetsDir, err := workdir.GetDotLetsDir(cwd)
	if err != nil {
		return nil, fmt.Errorf("can not get .lets path: %w", err)
	}

	if err := util.SafeCreateDir(dotLetsDir); err != nil {
		return nil, fmt.Errorf("can not create .lets dir: %w", err)
	}

	f, err := os.Open(cachedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cached remote config: %w", err)
	}
	defer f.Close()

	c := config.NewConfig(cwd, cachedPath, dotLetsDir)
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, fmt.Errorf("failed to parse remote config %s: %w", url, err)
	}

	if err = validate(c, version); err != nil {
		return nil, err
	}

	if err := c.SetupEnv(); err != nil {
		return nil, err
	}

	return c, nil
}

func ensureRemoteConfig(url string, noCache bool) (string, error) {
	cacheDir, err := remoteConfigCacheDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", fmt.Errorf("can not create remote config cache dir: %w", err)
	}

	cachePath := remoteConfigCachePath(cacheDir, url)

	if !noCache && util.FileExists(cachePath) {
		return cachePath, nil
	}

	data, downloadErr := fetch.Download(context.Background(), url)
	if downloadErr != nil {
		if util.FileExists(cachePath) {
			log.Warnf("failed to refresh remote config, using cached version: %s", downloadErr)
			return cachePath, nil
		}

		return "", fmt.Errorf("failed to download remote config: %w", downloadErr)
	}

	//#nosec G306
	if err := os.WriteFile(cachePath, data, 0o644); err != nil {
		return "", fmt.Errorf("failed to cache remote config: %w", err)
	}

	return cachePath, nil
}

func remoteConfigCacheDir() (string, error) {
	userDir, err := util.LetsUserDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userDir, "remote-configs"), nil
}

func remoteConfigCachePath(cacheDir, url string) string {
	hash := sha256.Sum256([]byte(url))
	return filepath.Join(cacheDir, hex.EncodeToString(hash[:])+".yaml")
}
