package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/codeclysm/extract"
)

var archAdaptMap = map[string]string{
	"386":   "i386",
	"amd64": "x86_64",
}

type RepoRegistry interface {
	GetLatestRelease() (string, error)
	DownloadReleaseBinary(packageName string, version string, dstPath string) error
	GetPackageName(os string, arch string) (string, error)
	GetDownloadURL(repoURI string, packageName string, version string) string
}

type GithubRegistry struct {
	client                 *http.Client
	ctx                    context.Context
	repoURI                string
	downloadURL            string
	downloadPackageTimeout time.Duration
}

func NewGithubRegistry(ctx context.Context) *GithubRegistry {
	client := &http.Client{
		Timeout: 15 * 60 * time.Second, // global timeout
	}

	reg := &GithubRegistry{
		client:                 client,
		ctx:                    ctx,
		repoURI:                "https://github.com/lets-cli/lets",
		downloadURL:            "",
		downloadPackageTimeout: 60 * 5 * time.Second, // TODO 5 minute timeout is enough?
	}

	return reg
}

func (reg *GithubRegistry) GetDownloadURL(repoURI string, packageName string, version string) string {
	return fmt.Sprintf("%s/releases/download/%s/%s", repoURI, version, packageName)
}

func (reg *GithubRegistry) GetPackageName(os string, arch string) (string, error) {
	os = strings.Title(os)

	arch, archExists := archAdaptMap[arch]
	if !archExists {
		return "", fmt.Errorf("arch %s is not supported", arch)
	}

	return fmt.Sprintf("lets_%s_%s", os, arch), nil
}

func (reg *GithubRegistry) DownloadReleaseBinary( //nolint:cyclop
	packageName string,
	version string,
	dstPath string,
) error {
	downloadURL := reg.GetDownloadURL(reg.repoURI, fmt.Sprintf("%s.tar.gz", packageName), version)

	ctx, cancel := context.WithTimeout(reg.ctx, reg.downloadPackageTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		downloadURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := reg.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no such package: %s", packageName)
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("network error: %s", resp.Status)
	}

	dstDir := fmt.Sprintf("%s.dir", dstPath)
	// cleanup if something abd happens during download/extract/rename flow
	defer os.RemoveAll(dstDir)

	err = os.RemoveAll(dstDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove download dir: %w", err)
		}
	}

	err = os.RemoveAll(dstPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove download path: %w", err)
		}
	}

	// TODO drop extract dependency, replace with own code
	err = extract.Gz(reg.ctx, resp.Body, dstDir, nil)
	if err != nil {
		return fmt.Errorf("failed to extract package: %w", err)
	}

	// since we do not need all content from tar, we take only binary and delete the rest
	err = os.Rename(path.Join(dstDir, "lets"), dstPath)
	if err != nil {
		return fmt.Errorf("failed to extract binary from package: %w", err)
	}

	return nil
}

type release struct {
	TagName string `json:"tag_name"`
}

func (reg *GithubRegistry) GetLatestRelease() (string, error) {
	ctx, cancel := context.WithTimeout(reg.ctx, 30*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/releases/latest", reg.repoURI)

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		url,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")

	resp, err := reg.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read package body: %w", err)
	}

	var release release
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("failed to decode package body: %w", err)
	}

	return release.TagName, nil
}
