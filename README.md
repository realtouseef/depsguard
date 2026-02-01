# DepGuard

DepGuard is a fast, offline-first Go CLI that blocks new dependencies unless at least 40% of existing dependencies are explained.

Repository: https://github.com/realtouseef/depsguard

## Install

```bash
go install github.com/realtouseef/depsguard@latest
```

## Quickstart

From the root of a Node.js project with a package.json:

```bash
depguard init
```

Add this script to package.json:

```json
{
  "scripts": {
    "preinstall": "depguard verify"
  }
}
```

Explain dependencies as needed:

```bash
depguard explain <dependency-name>
```

Check coverage:

```bash
depguard audit
```

## Commands

```bash
depguard init
depguard verify
depguard explain <dependency-name>
depguard audit
```

## How It Works

- On init, DepGuard snapshots dependencies into .depguard/baseline.json and creates .depguard/knowledge.json.
- On verify, if new dependencies are detected, DepGuard deterministically selects 40% of all dependencies and requires valid explanations.
- Explanations expire after 90 days.

## Knowledge Format

Each entry in .depguard/knowledge.json:

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

DepGuard seeds selection using a commit SHA when running in CI, falling back to local time otherwise. It checks these environment variables:

- GITHUB_SHA
- CI_COMMIT_SHA
- CIRCLE_SHA1
- TRAVIS_COMMIT
- BUILDKITE_COMMIT

## Output Files

- .depguard/baseline.json
- .depguard/knowledge.json

## Example Files

Example JSON files are included in examples/:

- examples/baseline.json
- examples/knowledge.json
