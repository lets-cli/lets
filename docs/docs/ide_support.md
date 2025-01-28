---
id: ide_support
title: IDE/Text editors support
---

### Jet Brains plugin (official)

Provides autocomplete and filetype support.

[Plugin site](https://plugins.jetbrains.com/plugin/14639-lets)
[Plugin repo](https://github.com/lets-cli/intellij-lets)

### VSCode plugin (official)

Provides autocomplete and filetype support.

[Plugin site](https://marketplace.visualstudio.com/items?itemName=kindritskyimax.vscode-lets)
[Plugin repo](https://github.com/lets-cli/vscode-lets)

### Emacs plugin (community)

Provides autocomplete and filetype support.

[Plugin repo](https://github.com/mpanarin/lets-mode)

### LSP

`LSP` stands for `Language Server Protocol`

Starting from `0.0.55` version lets comes with builtin `lsp` server under `lets self lsp` command.

Lsp support includes:

[x] Goto definition
  - Navigate to definitions of mixins files
  - Navigate to definitions of command from `depends`
[x] Completion
  - Complete commands in depends
[ ] Diagnostics
[ ] Hover
[ ] Formatting
[ ] Signature help
[ ] Code action

`lsp` server works with JSON Schema (see bellow).

#### VSCode

VSCode plugin supports lsp out of the box, you only want to make sure you have lets >= `0.0.55`.

#### Neovim

Neovim support for `lets self lsp` can be added manually:

1. Add new filetype:

```lua
vim.filetype.add({
  filename = {
    ["lets.yaml"] = "yaml.lets",
  },
})
```

2. In your `neovim/nvim-lspconfig` servers configuration:

In order for `nvim-lspconfig` to recognize `lets lsp` we must define config for `lets_ls` (lets_ls is just a conventional name because we are not officially added to `neovim/nvim-lspconfig`)

```lua
require("lspconfig.configs").lets_ls = {
  default_config = {
    cmd = { 
      "lets self lsp",
    },
    filetypes = { "yaml.lets" },
    root_dir = util.root_pattern("lets.yaml"),
    settings = {},
  },
}
```

3. And then enable lets_ls in then servers section:

```lua
return {
  "neovim/nvim-lspconfig",
  opts = {
    servers = {
      lets_ls = {},
      pyright = {},  -- pyright here just as hint to where we should add lets_ls
    },
  },
}
```

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

