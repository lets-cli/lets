---
name: create-release
description: Create a lets repo release or prerelease by selecting a version, testing, tagging, and verifying GitHub release output. Use when the user asks to create a release, prerelease, release candidate, rc, or tag for lets.
---

# Create Release

## Workflow

1. Establish intent: stable release or prerelease. If unclear, ask.
2. Inspect repo and release task:

   ```bash
   pwd
   git remote -v
   lets help release
   git status --short --branch
   git fetch --tags --prune
   git tag --sort=-v:refname | head -20
   ```

   Completion: repo is `lets-cli/lets`; branch/sync/cleanliness are known; latest stable and RC tags are known. If not on `master`, not up to date with `origin/master`, or dirty, report it and stop unless the user explicitly continues.
3. Select a candidate version using the rules below.
4. Planning checkpoint: tell the user the discovered latest stable, relevant RCs, proposed version, release message, and exact `lets release ...` command. Ask for confirmation before tests or tagging.
5. Preflight after confirmation:

   ```bash
   go test ./...
   ```

   For stable releases, also verify `docs/docs/changelog.md` has `[<version>]` because `lets release` enforces it. Completion: tests and stable changelog check pass, or release stops.
6. Remote-changing checkpoint: ask explicit approval to run the exact command:

   ```bash
   lets release <version> -m "Release <version>"
   lets release <version>-rcN -m "Prerelease <version>-rcN"
   ```

7. After approval, run `lets release ...`.
8. If `gh` is available and authenticated, wait for the Release workflow and verify the GitHub release:

   ```bash
   gh run list --workflow Release --limit 3 --json databaseId,status,conclusion,headBranch,headSha,displayTitle,createdAt,url
   gh run watch <run-id> --exit-status
   gh release view v<version> --json tagName,isPrerelease,publishedAt,url,assets
   ```

   If `gh` is unavailable or unauthenticated, report that verification was skipped.
9. Report tag, workflow result if checked, release URL if checked, and working tree status.

## Version-selection rules

Definitions:

- Stable tag: `vX.Y.Z`.
- RC tag: `vX.Y.Z-rcN`, where `N` is numeric.
- Base version: the `X.Y.Z` part of an RC.
- Unstabilized RC base: an RC base with at least one RC tag and no stable tag for the same base.
- `rc(next)`: one more than the highest numeric RC for the same base.
- Compare versions semver-numerically, not lexically.
- Never reuse an existing local or remote tag.
- If no stable tags exist, ask for the base version instead of guessing.

Prerelease:

- Explicit `X.Y.Z-rcN`: use it only if `vX.Y.Z` and `vX.Y.Z-rcN` do not exist.
- Explicit base `X.Y.Z`: if `vX.Y.Z` exists, stop and propose the next patch base. Otherwise propose `X.Y.Z-rc1` or the next RC for that base.
- No version: find the latest stable `S`. If there is an unstabilized RC base `B` greater than `S`, propose `B-rc(next)`. Otherwise propose `next-patch(S)-rc1`.

Stable release:

- Explicit `X.Y.Z`: use it only if `vX.Y.Z` does not exist and the changelog has `[X.Y.Z]`.
- No version: find the highest unstabilized RC base `B` greater than the latest stable and propose `B`. If no such RC base exists, ask which version to release.
- Do not infer a new stable patch without an RC unless the user explicitly asks for that version.

## Discovery helpers

```bash
# Stable tags only
git tag --list 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | grep -v -- '-rc' | head -10

# RC tags only
git tag --list 'v*-rc*' --sort=-v:refname | head -20

# Existing tags for a base version
git tag --list 'v<base>*' --sort=-v:refname

# Remote/local existence check for a candidate tag
git tag --list 'v<version>'
git ls-remote --tags origin 'refs/tags/v<version>'
```
