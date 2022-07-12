package parser

import (
	"fmt"

	"github.com/lets-cli/lets/config/config"
)

func parsePlugins(rawPlugins map[string]interface{}, newCmd *config.Command) error {
	plugins := make(map[string]config.CommandPlugin)

	for key, value := range rawPlugins {
		// TODO validate if plugin declared here is declared in config at the top
		pluginConfig, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("plugin %s configuration must be a mapping", key)
		}

		plugin := config.CommandPlugin{Name: key, Config: pluginConfig}
		plugins[key] = plugin
	}

	newCmd.Plugins = plugins

	return nil
}
