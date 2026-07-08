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

type loadOptions struct {
	progress fetch.ProgressObserver
	noCache  bool
}

type LoadOption func(*loadOptions)

func WithProgress(progress fetch.ProgressObserver) LoadOption {
	return func(opts *loadOptions) {
		opts.progress = progress
	}
}

func WithNoCache() LoadOption {
	return func(opts *loadOptions) {
		opts.noCache = true
	}
}

func newLoadOptions(options []LoadOption) loadOptions {
	opts := loadOptions{}
	for _, option := range options {
		option(&opts)
	}

	return opts
}

func Load(configName string, configDir string, version string) (*config.Config, error) {
	return LoadWithContext(context.Background(), configName, configDir, version)
}

func LoadWithContext(ctx context.Context, configName string, configDir string, version string, options ...LoadOption) (*config.Config, error) {
	opts := newLoadOptions(options)

	configPath, err := FindConfig(configName, configDir)
	if err != nil {
		return nil, err
	}

	return loadConfigFromFile(ctx, configPath.AbsPath, configPath.WorkDir, configPath.DotLetsDir, configPath.Filename, version, opts)
}

// LoadRemote downloads (or loads from cache) a remote lets.yaml at url and
// returns a Config with the working directory set to the caller's CWD.
func LoadRemote(ctx context.Context, url string, noCache bool, version string, options ...LoadOption) (*config.Config, error) {
	opts := newLoadOptions(options)
	if noCache {
		opts.noCache = true
	}

	cachedPath, err := ensureRemoteConfig(ctx, url, opts.noCache, opts.progress)
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

	c, err := loadConfigFromFile(ctx, cachedPath, cwd, dotLetsDir, url, version, opts)
	if err != nil {
		return nil, fmt.Errorf("%w (use --no-cache to re-download)", err)
	}

	c.RemoteSource = url

	return c, nil
}

// loadConfigFromFile is shared by Load and LoadRemote: opens the file at absPath,
// decodes YAML, validates, and sets up env. displayName appears in parse error messages.
func loadConfigFromFile(
	ctx context.Context,
	absPath, workDir, dotLetsDir, displayName, version string,
	opts loadOptions,
) (*config.Config, error) {
	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := config.NewConfig(workDir, absPath, dotLetsDir)
	c.SetDownloadOptions(ctx, opts.progress, opts.noCache)

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

func ensureRemoteConfig(ctx context.Context, url string, noCache bool, progress fetch.ProgressObserver) (string, error) {
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

	data, downloadErr := fetch.Download(ctx, url, fetch.WithProgress(fetch.SourceRemoteConfig, progress))
	if downloadErr != nil {
		if util.FileExists(cachePath) {
			log.Warnf("failed to download remote config (%v), falling back to cached version", downloadErr)
			return cachePath, nil
		}

		return "", fmt.Errorf("failed to download remote config: %w", downloadErr)
	}

	if err := util.WriteFileAtomic(cachePath, data); err != nil {
		return "", err
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
