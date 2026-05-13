# Contributing to Kara

Thanks for contributing to Kara.

This guide is technical-only and covers local setup, development commands, quality checks, commit conventions, and pull request expectations.

## Prerequisites

- Node.js 22+
- npm 10+
- pnpm available in your environment
- Git

## Project Structure

Kara is a monorepo managed with Turbo and npm workspaces.

- apps/web: Next.js web app
- apps/desktop: Electron desktop app
- apps/core: Go service

## Setup

1. Clone the repository.
2. Install dependencies:

```bash
pnpm install
```

## Development Commands

Run all apps in development mode:

```bash
pnpm dev
```

Run full monorepo build:

```bash
pnpm build
```

Clean build artifacts through Turbo:

```bash
pnpm clean
```

## Quality Gates

Before opening a PR, run:

```bash
pnpm lint
pnpm check
```

What these do:

- lint: runs Turbo lint tasks across workspaces
- check: runs Turbo checks across workspaces

Husky pre-commit hook runs:

```bash
npm run precommit
```

And precommit runs:

```bash
npm run lint && npm run check
```

## Commit Convention

This repository uses Conventional Commits through commitlint.

Use concise English commit messages, for example:

- feat: add desktop startup health check
- fix: resolve core health endpoint timeout
- docs: update setup instructions
- chore: update workspace tooling
- refactor: simplify desktop bootstrap flow

Recommended:

- Keep commits small and focused.
- Group related changes in one commit.
- Avoid mixing refactors and behavior changes in the same commit when possible.

## Branch and PR Workflow

1. Create a feature branch from develop.
2. Implement your change.
3. Run quality checks locally:

```bash
pnpm lint
pnpm check
pnpm build
```

4. Commit using Conventional Commits.
5. Open a PR to develop.

## Pull Request Checklist

Before requesting review, confirm:

- The change is scoped and clearly described.
- Local quality gates pass.
- Build passes for touched surfaces.
- Commit messages follow Conventional Commits.
- PR description includes:
    - Summary of the change
    - Why the change is needed
    - Verification steps executed locally

## Notes

- Release automation is handled by CI from main.
- Keep this file focused on contribution mechanics. Product and user documentation belongs in README and docs.
