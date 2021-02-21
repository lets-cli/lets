package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	getDownloadUrl(repoUri string, packageName string, version string) string
}

type githubRegistry struct {
	client                 *http.Client
	ctx                    context.Context
	repoUri                string
	downloadUrl            string
	downloadPackageTimeout time.Duration
}

func NewGithubRegistry() RepoRegistry {
	client := &http.Client{}

	ctx := context.Background()

	reg := githubRegistry{
		client:  client,
		ctx:     ctx,
		repoUri: "https://github.com/lets-cli/lets",
	}

	return reg
}

func (reg githubRegistry) getDownloadUrl(repoUri string, packageName string, version string) string {
	return fmt.Sprintf("%s/releases/download/%s/%s", repoUri, version, packageName)
}

func (reg githubRegistry) GetPackageName(os string, arch string) (string, error) {
	os = strings.Title(os)
	arch, archExists := archAdaptMap[arch]
	if !archExists {
		return "", fmt.Errorf("arch %s is not supported", arch)
	}
	return fmt.Sprintf("lets_%s_%s", os, arch), nil
}

func (reg githubRegistry) DownloadReleaseBinary(packageName string, version string, dstPath string) error {
	downloadUrl := reg.getDownloadUrl(reg.repoUri, fmt.Sprintf("%s.tar.gz", packageName), version)

	errFmt := func(err error) error {
		return fmt.Errorf("failed to download release %s version %s: %s", packageName, version, err)
	}

	ctx, cancel := context.WithTimeout(reg.ctx, 60*5*time.Second) // TODO 5 minute timeout is enough?
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		downloadUrl,
		nil,
	)
	if err != nil {
		return errFmt(err)
	}

	resp, err := reg.client.Do(req)
	if err != nil {
		return errFmt(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errFmt(fmt.Errorf("no such package: %s", packageName))
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errFmt(fmt.Errorf("network error: %s", resp.Status))
	}

	err = os.RemoveAll(dstPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errFmt(err)
		}
	}

	// TODO drop extract dependency, replace with own code
	err = extract.Gz(reg.ctx, resp.Body, dstPath, nil)
	if err != nil {
		return errFmt(err)
	}

	return nil
}

type release struct {
	TagName string `json:"tag_name"`
}

func (reg githubRegistry) GetLatestRelease() (string, error) {
	errFmt := func(err error) error {
		return fmt.Errorf("failed to get latest release version: %s", err)
	}

	ctx, cancel := context.WithTimeout(reg.ctx, 30*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/releases/latest", reg.repoUri)

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		url,
		nil,
	)
	if err != nil {
		return "", errFmt(err)
	}

	req.Header.Add("Accept", "application/json")

	resp, err := reg.client.Do(req)
	if err != nil {
		return "", errFmt(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errFmt(err)
	}

	var release release
	if err := json.Unmarshal(body, &release); err != nil {
		return "", errFmt(err)
	}

	return release.TagName, nil
}
