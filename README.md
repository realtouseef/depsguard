# depsguard

depsguard is a fast, offline-first Go CLI that blocks new dependencies unless at least 40% of existing dependencies are explained.

Repository: https://github.com/realtouseef/depsguard

## Install

```bash
go install github.com/realtouseef/depsguard/cmd/depsguard@latest
```

## Quickstart

From the root of a Node.js project with a package.json:

```bash
depsguard init
```

Add this script to package.json:

```json
{
  "scripts": {
    "preinstall": "depsguard verify"
  }
}
```

Explain dependencies as needed:

```bash
depsguard explain <dependency-name>
```

Check coverage:

```bash
depsguard audit
```

## Commands

```bash
depsguard init
depsguard verify
depsguard explain <dependency-name>
depsguard audit
```

## How It Works

- On init, depsguard snapshots dependencies into .depsguard/baseline.json and creates .depsguard/knowledge.json.
- On verify, if new dependencies are detected, depsguard deterministically selects 40% of all dependencies and requires valid explanations.
- Explanations expire after 90 days.

## Knowledge Format

Each entry in .depsguard/knowledge.json:

```json
{
  "dependency-name": {
    "summary": "string",
    "explained_by": "string",
    "expires_at": "RFC3339 timestamp"
  }
}
```

## Deterministic Selection

depsguard seeds selection using a commit SHA when running in CI, falling back to local time otherwise. It checks these environment variables:

- GITHUB_SHA
- CI_COMMIT_SHA
- CIRCLE_SHA1
- TRAVIS_COMMIT
- BUILDKITE_COMMIT

## Output Files

- .depsguard/baseline.json
- .depsguard/knowledge.json

## Example Files

Example JSON files are included in examples/:

- examples/baseline.json
- examples/knowledge.json
