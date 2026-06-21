# GoDriveLog v3.3 implementation prompts

Use one prompt per implementation slice. Each prompt is intentionally narrow so a chat can create one branch, do one slice, open one PR, and stop before the yak notices the clippers.

## Prompt files

| Version | File | Purpose |
|---|---|---|
| v3.3.0 | `v3.3.0-renderer-checkpoint.md` | Create v3.3 planning docs and reusable examples structure. |
| v3.3.1 | `v3.3.1-ebiten-renderer-spike.md` | Add an experimental Ebiten renderer beside Fyne. |
| v3.3.2 | `v3.3.2-renderer-ab-comparison.md` | Compare Fyne and Ebiten using the same baseline path. |
| v3.3.3 | `v3.3.3-renderer-decision.md` | Decide whether to continue, promote, pause, or abandon Ebiten. |

## Required pre-flight

Every implementation chat must:

1. Confirm the previous relevant PR is merged into `main`.
2. Confirm there are no blocking open PRs.
3. Create a branch from latest `main` using the full target version prefix.
4. Keep the slice narrow.
5. Update `CHANGES.md` and `docs/v3.3/ImplementationState.md`.
6. Open a PR and stop.

Branch names should sort by version, for example:

```text
v3.3.1-ebiten-renderer-spike
```
