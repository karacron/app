---
name: web-clean-code
description: "Use when you need clean-code guidance, refactors, reviews, or code generation for apps/web (Next.js App Router, React, TypeScript). Trigger phrases: web, Next.js, App Router, React, lint, refactor, clean code, review."
---

# Web Clean Code Skill

Use this skill when working in [apps/web](/Volumes/developer/authuser/karacron/apps/web) and you want a consistent syntax for writing, reviewing, or refactoring clean code.

## Common Syntax

Use this syntax to request work consistently:

1. Objective: state whether you want refactoring, review, or generation, and whether behavior must remain unchanged.
2. Scope: limit the request to the smallest surface in `apps/web` (specific route or component).
3. Constraints: specify limits such as no new dependencies, no out-of-scope UI changes, and no mixed functional changes.
4. Clean-code criteria: prioritize clear naming, small functions, low coupling, and removal of real duplication.
5. Validation: require lint and typecheck for the impacted area.
6. Expected output: request a short summary of changes, validations executed, and remaining risks.

Recommended request format:

- Objective: <what should be improved>
- Scope: <specific file or module in apps/web>
- Constraints: <non-negotiable limits>
- Criteria: <naming, function size, errors, duplication, comments>
- Validation: <commands that must pass>
- Output: <what the final summary must include>

## Stack Rules

- Follow Next.js 16 with App Router and strict TypeScript.
- Prefer small components and explicit props.
- Keep server and client component separation only when needed.
- Avoid duplicated logic in `src/app`; extract helpers or components as blocks grow.
- Respect existing ESLint rules and the `@/*` alias.

## Best Practices by Scope (Web)

- Small scope (1 file): prioritize renaming, function extraction, and condition simplification without changing public contracts.
- Medium scope (multiple files in one module): extract shared utilities, define explicit types, and remove duplication across nearby components.
- Large scope (full route or feature): split into phases, stabilize behavior first, then refactor structure; avoid mixing visual redesign with internal cleanup.
- Any scope: if you change naming or structure, preserve functional semantics and validate with lint + typecheck.

## Naming (Web)

- Components and types: `PascalCase` (`UserCard`, `ProfileFormProps`).
- Functions and variables: `camelCase` (`loadUserProfile`, `isSaving`, `errorMessage`).
- Booleans: use `is`, `has`, `can`, `should` prefixes.
- Event handlers: `handleX` (`handleSubmit`, `handleFilterChange`).
- Avoid generic names: `data`, `info`, `temp`, `value2`.

## Code Examples (Web)

Good:

```tsx
type UserCardProps = {
	userName: string;
	isActive: boolean;
};

function UserCard({ userName, isActive }: UserCardProps) {
	const statusLabel = isActive ? "Active" : "Inactive";
	return (
		<span>
			{userName} - {statusLabel}
		</span>
	);
}
```

Bad:

```tsx
function card(props: any) {
	const v = props.a ? "A" : "B";
	return (
		<span>
			{props.x} - {v}
		</span>
	);
}
```

Long function refactor:

```tsx
// Better: split filtering and mapping into small focused functions.
function buildVisibleUsers(users: User[]) {
	return users.filter(isVisibleUser).map(toUserListItem);
}
```

## Clean-Code Checklist

- Names reveal intent.
- Each function does one thing.
- There is no unnecessary local state.
- There is no bloated JSX when a helper component would be clearer.
- Do not introduce abstractions before they solve real repetition.

## Recommended Validation

- `npm --workspace web run lint`
- `npm --workspace web run check`

## Anti-Patterns

- Rewriting more than needed for the objective.
- Mixing UI, data, and refactor changes in one pass without reason.
- Ignoring a small improvement because it is not perfect.
- Introducing new dependencies for a simple problem.
