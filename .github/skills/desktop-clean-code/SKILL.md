---
name: desktop-clean-code
description: "Use when you need clean-code guidance, refactors, reviews, or code generation for apps/desktop (Electron main process, preload, packaging). Trigger phrases: desktop, Electron, preload, IPC, main process, refactor, clean code, review."
---

# Desktop Clean Code Skill

Use this skill when working in [apps/desktop](/Volumes/developer/authuser/karacron/apps/desktop) and you want a consistent syntax for clean Electron code.

## Common Syntax

Use this syntax to request work consistently:

1. Objective: state whether you want refactoring, review, or generation, and whether it impacts security, IPC, packaging, or UI.
2. Scope: limit the request to `main.js`, `preload.js`, or a specific script in `apps/desktop`.
3. Constraints: specify limits such as keeping context isolation and not expanding exposed APIs without reason.
4. Clean-code criteria: prioritize small functions, minimal IPC surface, clear paths, and separated responsibilities.
5. Validation: require desktop build or packaging checks.
6. Expected output: request a short summary of changes, validations executed, and security/compatibility risks.

Recommended request format:

- Objective: <what should be improved>
- Scope: <specific file or script in apps/desktop>
- Constraints: <non-negotiable limits>
- Criteria: <naming, function size, errors, duplication, comments>
- Validation: <commands that must pass>
- Output: <what the final summary must include>

## Stack Rules

- Keep `contextIsolation: true` and `nodeIntegration: false`.
- Use `preload.js` as the only access bridge from the renderer.
- Expose the minimum possible API in `contextBridge`.
- Avoid mixing path resolution, window startup, and business logic in one function when it can be separated.
- Keep packaging reproducible and resolve `core-bin` and `web-dist` correctly.

## Best Practices by Scope (Desktop Electron)

- Small scope (1 file): improve naming and separate resolution functions from side-effect functions without expanding IPC surface.
- Medium scope (main + preload): review security boundaries, keep APIs explicit, and reduce data exposed to the renderer.
- Large scope (full desktop feature): split startup, event wiring, and resource access into separate steps; validate security first, then packaging.
- Any scope: keep `contextIsolation` enabled, do not expose `require`/`process`, and check build/package before closing.

## Naming (Desktop Electron)

- Resolution functions: `resolveX` (`resolveWebEntry`, `resolveCoreBinary`).
- Side-effect functions: clear verbs (`createWindow`, `registerAppEvents`).
- Intentional window/process variables (`mainWindow`, `coreBinaryPath`).
- Explicit names for preload API (`getVersion`, `getPlatform`).
- Avoid ambiguous names: `doThing`, `x`, `tmpData`.

## Code Examples (Desktop)

Good (`main.js`):

```js
function resolveCoreBinaryPath() {
	const executableName =
		process.platform === "win32" ? "kara-core.exe" : "kara-core";
	return app.isPackaged
		? path.join(process.resourcesPath, "core-bin", executableName)
		: path.join(__dirname, "core-bin", executableName);
}
```

Bad (`main.js`):

```js
function x() {
	return process.platform === "win32" ? "a" : "b";
}
```

Good (`preload.js`):

```js
contextBridge.exposeInMainWorld("kara", {
	getPlatform: () => process.platform,
	getMode: () => process.env.NODE_ENV || "production",
});
```

Bad (`preload.js`):

```js
contextBridge.exposeInMainWorld("api", {
	require,
	process,
});
```

## Clean-Code Checklist

- The main process remains small and readable.
- The IPC surface is minimal and explicit.
- No more power than needed is exposed to the renderer.
- Paths and resources resolve predictably.
- The change works in both local and packaged modes.

## Recommended Validation

- `npm --workspace desktop run build`
- `npm --workspace desktop run package`

## Anti-Patterns

- Exposing full Node access to the renderer.
- Adding broad IPC without real need.
- Mixing security, UI, and packaging in one hard-to-follow block.
- Using complex solutions for a simple bridge.
