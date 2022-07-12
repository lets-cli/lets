package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type ConfigPlugin struct {
	Name string
	// if Repo not specified, then it is lets own plugins, lets-cli/lets-plugin-<name>
	// if Repo in format <repo>, then we append Name to Repo, <repo>/<name> # TODO or maybe <repo>/lets-plugin-<name>
	Repo    string
	Version string
	Url     string
	Bin     string
}

type pluginResult struct {
	ResultType string `json:"type"`
	Result     string `json:"result"`
}

func (p ConfigPlugin) Exec(commandPlugin CommandPlugin) error {
	config, err := commandPlugin.SerializeConfig()
	if err != nil {
		return err
	}

	cmd := exec.Command(
		p.Bin,
		string(config), // TODo how to encode
	) // #nosec G204

	var out bytes.Buffer
	// TODO how to show plugin logs ?
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
		//return &RunErr{err: fmt.Errorf("failed to run child command '%s' from 'depends': %w", r.cmd.Name, err)}
	}

	// TODO check exit code
	output := out.String()
	fmt.Printf("output %s", output)

	result := pluginResult{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return err
	}
	fmt.Printf("result %v", result)

	if result.ResultType == "json" {
		pluginResponse := map[string]interface{}{}

		if err := json.Unmarshal([]byte(result.Result), &pluginResponse); err != nil {
			return err
		}

		fmt.Printf("pluginResponse %v", pluginResponse)
	}
	return nil
}

// TODO maybe plugin must have lifecycle
type CommandPlugin struct {
	Name   string
	Config map[string]interface{}
}

func (p CommandPlugin) Run(cmd *Command, cfg *Config) error {
	plugin := cfg.Plugins[p.Name]
	return plugin.Exec(p)
}

func (p CommandPlugin) SerializeConfig() ([]byte, error) {
	return json.Marshal(p.Config)
}
