# Issue tracker: GitHub

The **Issue tracker** for this repo is GitHub. Work items and PRDs live as **GitHub Issues** in `lets-cli/lets`. Use the `gh` CLI for all operations.

## Conventions

- **Create a GitHub Issue**: `gh issue create --title "..." --body "..."`. Use a heredoc for multi-line bodies.
- **Read a GitHub Issue**: `gh issue view <number> --comments`, and fetch labels when they matter to the workflow.
- **List GitHub Issues**: `gh issue list --state open --json number,title,body,labels,comments --jq '[.[] | {number, title, body, labels: [.labels[].name], comments: [.comments[].body]}]'` with appropriate `--label` and `--state` filters.
- **Comment on a GitHub Issue**: `gh issue comment <number> --body "..."`
- **Apply / remove Triage labels**: `gh issue edit <number> --add-label "..."` / `--remove-label "..."`
- **Close a GitHub Issue**: `gh issue close <number> --comment "..."`

Infer the repo from `git remote -v` — `gh` does this automatically when run inside a clone.

## When a skill says "publish to the issue tracker"

Create a **GitHub Issue**.

## When a skill says "fetch the relevant ticket"

Fetch the relevant **GitHub Issue** with `gh issue view <number> --comments`.
