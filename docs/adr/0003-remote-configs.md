# ADR-0003 — Support Remote configs with local caching

**Date:** 2026-06-14
**Status:** Accepted

## Context

This decision was originally captured as a design spec on 2026-06-13 and later promoted to an ADR. It originated from Issue [#351](https://github.com/lets-cli/lets/issues/351).

`lets` could load a local **Project config** from `lets.yaml`, but users also wanted to run reusable project workflows published at a URL:

```bash
lets -c https://example.com/lets.yaml build
```

A **Remote config** needed different behavior from a local file:

- avoid re-downloading unchanged configs on every invocation
- let users force a refresh when the remote source changes
- preserve the invocation directory as the **Work dir** for commands
- share HTTP download behavior with **Remote mixins** instead of duplicating it
- fail safely when the network is unavailable

## Decision

Treat a root `--config` / `-c` value that starts with `http://` or `https://` as a **Remote config**.

Add a root `--no-cache` flag. For remote sources, this asks lets to re-download instead of using an existing cached copy. The same no-cache choice is passed through config loading so **Remote mixins** refresh consistently with **Remote configs**.

Cache downloaded Remote configs at:

```text
~/.config/lets/remote-configs/<sha256(url)>.yaml
```

The URL hash is the cache identity. Cache writes are atomic: write a sibling temp file, set file mode, then rename into place.

Remote config loading follows this flow:

1. Without `--no-cache`, use the cache if present.
2. Otherwise download the URL, write it to the cache, then load the cached file.
3. If the download fails and a cache file exists, warn and fall back to the cached file.
4. If the download fails with no cache file, return an error.

Load the cached YAML through the same config validation and setup path as local Project configs, but set the config **Work dir** and `.lets` directory from the invocation CWD. Commands from a Remote config therefore run from where the user invoked `lets`, unless a command declares `work_dir`.

Extract shared HTTP download behavior into `internal/fetch` so Remote configs and Remote mixins use the same implementation. `fetch.Download` accepts only `text/plain`, `text/yaml`, `text/x-yaml`, `application/yaml`, and `application/x-yaml`; reports non-2xx responses; supports progress observers; and is cancellable through context.

If `LETS_CONFIG_DIR` is set while using a Remote config URL, warn and ignore it because discovery relative to a local config directory does not apply.

## Consequences

- **Positive:** Users can publish and consume shared Project configs by URL.
- **Positive:** Cached configs make repeated invocations fast and provide an offline fallback.
- **Positive:** `--no-cache` gives users a direct way to refresh Remote configs and Remote mixins.
- **Positive:** Remote config downloads and Remote mixin downloads share one fetch implementation and progress-reporting path.
- **Neutral:** A cached file that no longer parses produces a normal config error with guidance to retry using `--no-cache`.
- **Neutral:** The invocation CWD, not the cache directory, defines the default Work dir for Remote config commands.
- **Negative:** Remote configs are executable project workflow definitions; users must trust the URL they run.
- **Negative:** The cache is keyed only by URL. Different content at the same URL replaces the prior cached config when refreshed.

## Related implementation

- `internal/cli/cli.go`
- `internal/cmd/root.go`
- `internal/config/load.go`
- `internal/config/config/mixin.go`
- `internal/fetch/fetch.go`
- `internal/config/load_test.go`
- `internal/fetch/fetch_test.go`
