## Summary

Describe what changed and why.

## Scope

- [ ] web
- [ ] core
- [ ] desktop
- [ ] docs
- [ ] ci/cd

## Checklist

- [ ] I kept the change focused and minimal
- [ ] I followed clean-code naming and function-size conventions
- [ ] I did not mix unrelated refactors
- [ ] I validated impacted workspaces locally

## Validation

List commands executed and results.

- [ ] npm --workspace web run lint
- [ ] npm --workspace web run check
- [ ] npm --workspace core run lint
- [ ] npm --workspace core run check
- [ ] npm --workspace desktop run build
- [ ] npm --workspace desktop run package

## Risk / Rollback

- Risk level: low / medium / high
- Rollback plan:
