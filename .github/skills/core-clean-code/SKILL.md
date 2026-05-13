---
name: core-clean-code
description: "Use when you need clean-code guidance, refactors, reviews, or code generation for apps/core (Go). Trigger phrases: core, Go, go vet, go test, refactor, clean code, review."
---

# Core Clean Code Skill

Use this skill when working in [apps/core](/Volumes/developer/authuser/karacron/apps/core) and you want a consistent syntax for writing and reviewing clean Go code.

## Common Syntax

Use this syntax to request work consistently:

1. Objective: state whether you want refactoring, review, or generation, and whether behavior must remain unchanged.
2. Scope: limit the request to the minimum package, file, or handler within `apps/core`.
3. Constraints: specify boundaries such as no unnecessary layers and no public contract changes.
4. Clean-code criteria: prioritize intention-revealing names, small functions, explicit error handling, and simplicity.
5. Validation: require checks with `go vet`, tests, and build for the impacted area.
6. Expected output: request a short summary of changes, validations executed, and remaining technical debt.

Recommended request format:

- Objective: <what should be improved>
- Scope: <specific file, package, or handler in apps/core>
- Constraints: <non-negotiable limits>
- Criteria: <naming, function size, errors, duplication, comments>
- Validation: <commands that must pass>
- Output: <what the final summary must include>

## Stack Rules

- Use idiomatic Go and avoid premature abstractions.
- Keep `main.go` simple; if it grows, move logic into small functions or packages.
- Prefer explicit and manageable errors.
- Use short but precise names.
- Keep handlers small and predictable.

## Best Practices by Scope (Core Go)

- Small scope (1 function): apply single responsibility, improve naming, and return contextual errors.
- Medium scope (handler or file): separate parsing, logic, and HTTP response into distinct functions; avoid duplicated validations.
- Large scope (package): move responsibilities into lightweight layers, preserve existing contracts, and add test coverage around the change.
- Any scope: avoid over-engineering, prioritize clarity over abstractions, and validate with vet + tests + build.

## Naming (Core Go)

- Exported functions: `PascalCase` (`BuildHealthResponse`).
- Internal functions: `camelCase` (`buildHealthResponse`).
- Intention-revealing variables: `serverAddr`, `healthPayload`, `requestID`.
- Contextual errors: `fmt.Errorf("read config: %w", err)`.
- Avoid opaque abbreviations: `cfg2`, `tmp`, `x1`.

## Code Examples (Core Go)

Good:

```go
func buildHealthResponse(serviceName string) []byte {
	return []byte(fmt.Sprintf(`{"status":"ok","service":"%s"}`, serviceName))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(buildHealthResponse("kara-core"))
}
```

Bad:

```go
func h(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
```

Correct error handling:

```go
func loadConfig(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}
	return b, nil
}
```

## Clean-Code Checklist

- Logic is placed where it belongs.
- Errors are handled explicitly.
- There is no obvious duplication.
- Functions are short and easy to test.
- Code remains clear without superfluous comments.

## Recommended Validation

- `npm --workspace core run lint`
- `npm --workspace core run check`
- `npm --workspace core run build`

## Anti-Patterns

- Putting business logic in `main` without need.
- Hiding errors with `panic` when a reasonable handling path exists.
- Creating extra interfaces or layers for a simple flow.
- Optimizing before it is needed.
