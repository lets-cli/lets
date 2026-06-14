---
id: settings
title: Settings
---

`lets` settings control the behavior of `lets` itself.

Use settings for things like colored output, theming, or update notifications. Do not use this file for project commands or runtime env. Project behavior still belongs in `lets.yaml`.

## Settings file location

`lets` reads settings from:

```text
~/.config/lets/config.yaml
```

This file is per-user and applies to all projects on the machine.

Print the path:

```bash
lets self config path
```

Open the file in `$EDITOR`:

```bash
lets self config edit
```

## Precedence

Settings are resolved in this order:

1. environment variables
2. settings file
3. built-in defaults

This means env vars always win over `config.yaml`.

## Supported settings

### `no_color`

Disable colored output from `lets`.

Example:

```yaml
no_color: true
```

Environment override:

- `NO_COLOR` disables colors even if `no_color` is not set

Note:

- this affects `lets` output itself
- it does not inject `NO_COLOR` into commands from `lets.yaml`

### `theme`

Choose the theme for `lets` styled help and error output.

Supported values:

- `default`
- `ansi`
- `synthwave`

Example:

```yaml
theme: synthwave
```

Environment override:

- none

Default:

- `theme: default`

Notes:

- this affects `lets` output itself
- project commands still control their own colors

### `upgrade_notify`

Enable or disable background update notifications for interactive sessions.

Example:

```yaml
upgrade_notify: false
```

Environment override:

- `LETS_CHECK_UPDATE` disables update checks and notifications regardless of the settings file

Default:

- `upgrade_notify: true`

## Example

```yaml
no_color: false
theme: default
upgrade_notify: false
```

## Invalid settings

Unknown keys and invalid YAML cause `lets` startup to fail with an error. Keep this file limited to supported settings only.
