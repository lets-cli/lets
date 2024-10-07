---
id: ide_support
title: IDE/Text editors support
---

### Jet Brains plugin (official)

Provides autocomplete and filetype support.

[Plugin site](https://plugins.jetbrains.com/plugin/14639-lets)

### Emacs plugin (community)

Provides autocomplete and filetype support.

[Plugin site](https://github.com/mpanarin/lets-mode)

### JSON Schema

In order to get autocomplete and filetype support in any editor, you can use the JSON schema file provided by Lets.

#### VSCode

To use the JSON schema in VSCode, you can use the [YAML extension](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml).

Add the following to your `settings.json`:

```json
{
  "yaml.schemas": {
    "https://lets-cli.org/schema.json": [
      "**/lets.yaml",
      "**/lets*.yaml",
    ]
  }
}
```

#### Neovim

To use the JSON schema in Neovim, you can use the `nvim-lspconfig` with `SchemaStore` plugin.

In your `nvim-lspconfig` configuration, add the following:

```lua
servers = {
  yamlls = {
    on_new_config = function(new_config)
      local yaml_schemas = require("schemastore").yaml.schemas({
        extra = {
          {
            description = "Lets JSON schema",
            fileMatch = { "lets.yaml", "lets*.yaml" },
            name = "lets.schema.json",
            url = "https://lets-cli.org/schema.json",
          },
        },
      })
      new_config.settings.yaml.schemas = vim.tbl_deep_extend("force", new_config.settings.yaml.schemas or {}, yaml_schemas)
    end,
  },
}
```

