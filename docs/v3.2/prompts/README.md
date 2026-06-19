# v3.2 prompt index

Use these prompts one slice at a time.

Before implementing a slice:

1. Read `docs/v3.2/README.md`.
2. Read `docs/v3.2/ReleasePlan.md`.
3. Read `docs/v3.2/ImplementationState.md`.
4. Read `docs/v3.2/OpenDecisions.md`.
5. Read `docs/v3.2/CarryForward.md`.
6. Read the target prompt file.
7. Confirm the previous relevant PR is merged into `main`.
8. Confirm there are no blocking open PRs.
9. Branch from latest `main`.
10. Implement only the requested slice.
11. Update `CHANGES.md` and `docs/v3.2/ImplementationState.md`.
12. Update `OpenDecisions.md` only when needed.
13. Open a PR and stop.

## Prompts

| Version | Prompt |
|---|---|
| v3.2.1 | `v3.2.1-gauge-package-loader.md` |
| v3.2.2 | `v3.2.2-gauge-widget-support.md` |
| v3.2.3 | `v3.2.3-radial-gauge-scene-model.md` |
| v3.2.4 | `v3.2.4-fyne-radial-rendering.md` |
| v3.2.5 | `v3.2.5-example-gauge-package.md` |
| v3.2.6 | `v3.2.6-harness-verification.md` |
| v3.2.7 | `v3.2.7-checkpoint.md` |

## Common guardrails

- Do not add clusters in v3.2 unless a later checkpoint explicitly decides to.
- Do not add gauge inheritance.
- Do not add widget-level sensor override.
- Do not infer gauge type from directory names.
- Do not break existing widget types.
- Do not render fake live gauge values for non-`ok` sensor states.
- Keep Fyne-specific work in the display adapter.
