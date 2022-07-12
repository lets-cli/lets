package plugins

import (
	"fmt"
	"strings"

	"github.com/lets-cli/lets/config/config"
)

func Load(cfg *config.Config) error {
	// 1. check if plugin exist on .lets
	// 2. check if version is downloaded in .lets
	// 3. downloading progress bar
	for _, plugin := range cfg.Plugins {
		if plugin.Bin != "" {
			// if bin specified, skip downloading new version
			// TODO do we need to copypaste binary to .lets ?
			continue
		}

		// TODO validate repo and url
		if plugin.Url == "" {
			plugin.Url = getDefaultDownloadUrl(plugin)
		} else {
			plugin.Url = expandUrl(plugin, cfg)
		}

		// TODO download from url

	}
	return nil
}

func getDefaultDownloadUrl(plugin config.ConfigPlugin) string {
	repo := plugin.Repo
	if repo == "" {
		repo = fmt.Sprintf("lets-cli/lets-plugin-%s", plugin.Name)
	} else if !strings.Contains(repo, "/") {
		repo = fmt.Sprintf("%s/lets-plugin-%s", repo, plugin.Name)
	}

	//https://github.com/lets-cli/lets/releases/download/{{.Version}}/lets_{{.Os}}_{{.Arch}}.tar.gz
	os := "linux"
	arch := "amd64"
	bin := fmt.Sprintf("lets_plugin_%s_%s_%s", plugin.Name, os, arch)
	// TODO require bin or tar.gz ?
	version := plugin.Version // TODO what if latest ?
	return fmt.Sprintf(
		"https://github.com/%s/releases/download/%s/%s",
		repo, version, bin,
	)
}

func expandUrl(plugin config.ConfigPlugin, cfg *config.Config) string {
	url := plugin.Url

	if strings.Contains(url, "{{.Version}}") {
		// TODO well we must use go templates here ))
		url = strings.Replace(url, "{{.Version}}", plugin.Version, 1)
	}

	return url
}
