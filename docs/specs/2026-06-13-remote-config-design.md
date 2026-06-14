---
name: remote-config-design
description: Design for lets -c https://... remote config fetching with caching (issue #351)
---

# Remote Config Design

**Issue:** [#351](https://github.com/lets-cli/lets/issues/351)

## Summary

Allow `lets -c https://url` to download a remote `lets.yaml`, cache it locally, and run commands from it with the invocation directory as the working directory. Add `--no-cache` flag to force re-download.

## Architecture

### New flag: `--no-cache`

Added to `initRootFlags` in `internal/cmd/root.go` and parsed in `parseRootFlags` in `internal/cli/cli.go`, following the same pattern as `--debug` and `--init`.

### URL detection in `cli.go`

After parsing root flags, if `rootFlags.config` starts with `http://` or `https://`, call `config.LoadRemote(url, noCache, version)` instead of `config.Load(...)`. The existing `Load` path is untouched.

### New `internal/fetch` package

Extract HTTP download + content-type validation from `RemoteMixin.download()` into `internal/fetch/fetch.go`:

```go
func Download(ctx context.Context, url string) ([]byte, error)
```

Same 15-minute timeout and content-type whitelist (`text/plain`, `text/yaml`, `text/x-yaml`, `application/yaml`, `application/x-yaml`). `RemoteMixin` becomes a thin wrapper calling `fetch.Download`. This eliminates duplication.

### New `LoadRemote` in `internal/config/load.go`

```go
func LoadRemote(url string, noCache bool, version string) (*config.Config, error)
```

- Cache path: `util.LetsUserDir()` + `remote-configs/<hex(sha256(url))>.yaml`
  - Resolves to `~/.config/lets/remote-configs/<sha256>.yaml`
- Creates `~/.config/lets/remote-configs/` if needed
- Working directory: `os.Getwd()` (CWD at invocation time), passed as `configDir` to `Load`

## Data Flow

### Normal flow (no `--no-cache`)

1. Cache file exists → load directly from cached path, skip HTTP
2. Cache missing → download via `fetch.Download` → persist → load from cached path

### `--no-cache` flow

1. Attempt download → persist (overwriting cache) → load
2. Download fails + cache exists → `log.Warnf("failed to refresh remote config, using cached version: %s", err)` → load from cache
3. Download fails + no cache → return error

### Working directory

`LoadRemote` calls `Load(cachedPath, cwd, version)` where `cwd = os.Getwd()`. Commands in the remote config execute in the invocation directory unless they specify `work_dir`.

## Files Changed

| File | Change |
|------|--------|
| `internal/fetch/fetch.go` | New — extracted `Download` func |
| `internal/fetch/fetch_test.go` | New — unit tests for fetch |
| `internal/config/config/mixin.go` | Refactor `download()` to use `fetch.Download` |
| `internal/config/load.go` | Add `LoadRemote` |
| `internal/config/load_test.go` | Add `LoadRemote` tests |
| `internal/cli/cli.go` | URL detection, call `LoadRemote`, thread `noCache` |
| `internal/cmd/root.go` | Add `--no-cache` flag |

## Testing

### `internal/fetch/fetch_test.go`
- Valid YAML content-type → success
- Unsupported content-type → error
- 404 → error
- Non-2xx → error

### `internal/config/load_test.go`
- `LoadRemote` with `httptest.NewServer` serving valid YAML → loads, caches, runs from CWD
- Cache hit: serve once, call twice, assert server receives only one request
- `--no-cache` refresh: serve two different configs, assert second call gets new content
- `--no-cache` + server down + cache exists: assert warning logged + falls back to cache
- `--no-cache` + server down + no cache: assert error returned
