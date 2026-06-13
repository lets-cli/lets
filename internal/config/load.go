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

	return loadConfigFromFile(configPath.AbsPath, configPath.WorkDir, configPath.DotLetsDir, configPath.Filename, version)
}

// LoadRemote downloads (or loads from cache) a remote lets.yaml at url and
// returns a Config with the working directory set to the caller's CWD.
func LoadRemote(ctx context.Context, url string, noCache bool, version string) (*config.Config, error) {
	cachedPath, err := ensureRemoteConfig(ctx, url, noCache)
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

	c, err := loadConfigFromFile(cachedPath, cwd, dotLetsDir, url, version)
	if err != nil {
		return nil, fmt.Errorf("%w (use --no-cache to re-download)", err)
	}

	c.RemoteSource = url

	return c, nil
}

// loadConfigFromFile is shared by Load and LoadRemote: opens the file at absPath,
// decodes YAML, validates, and sets up env. displayName appears in parse error messages.
func loadConfigFromFile(absPath, workDir, dotLetsDir, displayName, version string) (*config.Config, error) {
	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := config.NewConfig(workDir, absPath, dotLetsDir)
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", displayName, err)
	}

	if err = validate(c, version); err != nil {
		return nil, err
	}

	if err := c.SetupEnv(); err != nil {
		return nil, err
	}

	return c, nil
}

func ensureRemoteConfig(ctx context.Context, url string, noCache bool) (string, error) {
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

	data, downloadErr := fetch.Download(ctx, url)
	if downloadErr != nil {
		if util.FileExists(cachePath) {
			log.Warnf("failed to download remote config (%v), falling back to cached version", downloadErr)
			return cachePath, nil
		}

		return "", fmt.Errorf("failed to download remote config: %w", downloadErr)
	}

	if err := writeCacheAtomic(cachePath, data); err != nil {
		return "", err
	}

	return cachePath, nil
}

// writeCacheAtomic writes data to a sibling temp file then renames it to dst,
// ensuring the cache path is never left in a partially-written state.
func writeCacheAtomic(dst string, data []byte) error {
	dir := filepath.Dir(dst)

	tmp, err := os.CreateTemp(dir, "*.yaml.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp cache file: %w", err)
	}

	tmpPath := tmp.Name()

	_, writeErr := tmp.Write(data)
	closeErr := tmp.Close()

	if writeErr != nil || closeErr != nil {
		os.Remove(tmpPath)
		if writeErr != nil {
			return fmt.Errorf("failed to write temp cache file: %w", writeErr)
		}

		return fmt.Errorf("failed to close temp cache file: %w", closeErr)
	}

	//#nosec G306
	if err := os.Chmod(tmpPath, 0o644); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to chmod temp cache file: %w", err)
	}

	if err := os.Rename(tmpPath, dst); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to cache remote config: %w", err)
	}

	return nil
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
