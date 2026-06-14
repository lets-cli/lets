package config

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lets-cli/lets/internal/fetch"
	"github.com/lets-cli/lets/internal/util"
)

type Mixins []*Mixin

type Mixin struct {
	FileName string
	// e.g. .gitignored
	Ignored bool
	Remote  *RemoteMixin
}

type RemoteMixin struct {
	URL     string
	Version string

	mixinsDir string
}

// Filename is name of mixin file (hash from url).
func (rm *RemoteMixin) Filename() string {
	hasher := sha256.New()
	hasher.Write([]byte(rm.URL))

	if rm.Version != "" {
		hasher.Write([]byte(rm.Version))
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

// Path is abs path to mixin file (.lets/mixins/<filename>).
func (rm *RemoteMixin) Path() string {
	return filepath.Join(rm.mixinsDir, rm.Filename())
}

func (rm *RemoteMixin) persist(data []byte) error {
	f, err := os.OpenFile(rm.Path(), os.O_CREATE|os.O_WRONLY, 0o755) //nolint:nosnakecase
	if err != nil {
		return fmt.Errorf("can not open file %s to persist mixin: %w", rm.Path(), err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("can not write mixin to file %s: %w", rm.Path(), err)
	}

	return nil
}

func (rm *RemoteMixin) exists() bool {
	return util.FileExists(rm.Path())
}

func (rm *RemoteMixin) tryRead() ([]byte, error) {
	if !rm.exists() {
		return nil, nil
	}

	data, err := os.ReadFile(rm.Path())
	if err != nil {
		return nil, fmt.Errorf("can not read mixin config file at %s: %w", rm.Path(), err)
	}

	return data, nil
}

func (rm *RemoteMixin) download(ctx context.Context, progress fetch.ProgressObserver) ([]byte, error) {
	return fetch.Download(ctx, rm.URL, fetch.WithProgress(fetch.SourceRemoteMixin, progress))
}

// Trim `-` prefix.
// Using this prefix we allow to include non-existed mixins (git-ignored for example).
func normalizeMixinFilename(filename string) string {
	return strings.TrimPrefix(filename, "-")
}

// Ignored means that it is okay if minix does not exist.
// It can be a git-ignored file for example.
func isIgnoredMixin(filename string) bool {
	return strings.HasPrefix(filename, "-")
}

func (m *Mixin) UnmarshalYAML(unmarshal func(any) error) error {
	var filename string
	if err := unmarshal(&filename); err == nil {
		m.FileName = normalizeMixinFilename(filename)
		m.Ignored = isIgnoredMixin(filename)

		return nil
	}

	var remote struct {
		URL     string
		Version string
	}

	if err := unmarshal(&remote); err != nil {
		return err
	}

	m.Remote = &RemoteMixin{
		// TODO check if url is valid
		URL:     remote.URL,
		Version: remote.Version,
	}

	return nil
}

func (m *Mixin) IsRemote() bool {
	return m.Remote != nil
}

func (m *Mixin) Source() string {
	if m.IsRemote() {
		return m.Remote.URL
	}

	return m.FileName
}
