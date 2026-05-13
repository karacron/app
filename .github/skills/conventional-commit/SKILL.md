---
name: conventional-commit
description: "Use when you need to write git commit messages in English that pass husky + commitlint (Conventional Commits). Trigger phrases: commit message, conventional commit, husky commitlint, commit format, chore, feat, fix, breaking change."
---

# Conventional Commit Skill

Use this skill when you need commit messages that pass the repository's Husky + Commitlint checks.

## Goal

- Produce commit messages in English.
- Follow the Conventional Commits format required by `@commitlint/config-conventional`.
- Keep messages short, clear, and release-friendly.

## Required Format

Use this structure:

- `<type>: <subject>`
- `<type>(<scope>): <subject>`

Examples:

- `feat: add user profile endpoint`
- `fix(core): handle nil config in loader`
- `chore(ci): update release workflow`

## Allowed Types

Use standard Conventional Commit types:

- `feat`
- `fix`
- `docs`
- `style`
- `refactor`
- `perf`
- `test`
- `build`
- `ci`
- `chore`
- `revert`

## English Writing Rules

- Subject must be in English.
- Use imperative mood: `add`, `fix`, `update`, `remove`.
- Keep subject concise and specific.
- Start with lowercase after `type:`.
- Do not end the subject with a period.

Good:

- `feat(web): add locale switcher`
- `fix(desktop): prevent preload api overexposure`
- `chore: add codeowners`

Bad:

- `Added locale switcher`
- `Fix bug`
- `chore: Added CODEOWNERS.`

## Scope Guidelines

Use scope when it improves clarity:

- `web` for `apps/web`
- `core` for `apps/core`
- `desktop` for `apps/desktop`
- `ci` for workflows and automation
- `docs` for documentation-only changes

Examples:

- `feat(web): add chat input validation`
- `fix(core): return contextual error on config read`
- `chore(desktop): align preload api naming`

## Breaking Changes

For breaking changes, use one of these options:

- `feat!: remove legacy auth flow`
- Add a footer:

```text
BREAKING CHANGE: removed legacy auth flow and old token format
```

## Multi-line Commit Template

Use this template when extra context is needed:

```text
type(scope): short imperative summary

Why this change was needed.
What was changed at a high level.

Refs: #123
```

## Validation

Before finalizing commits:

- Ensure message follows Conventional Commit format.
- Ensure subject is in English.
- Ensure the message would pass Husky/Commitlint hooks.

## Quick Prompt Syntax

Use this request format with the assistant:

- Change summary: <what changed>
- Scope: <web|core|desktop|ci|docs|none>
- Type hint: <feat|fix|chore|...>
- Breaking: <yes|no>
- Output: <single commit or list of commits>
