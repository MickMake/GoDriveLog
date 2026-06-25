# GoDriveLog v3.4 prompts

This directory contains one prompt per planned v3.4 implementation slice.

Each prompt should be used in a fresh implementation chat after reading:

1. `docs/v3.4/ReleasePlan.md`
2. `docs/v3.4/ImplementationState.md`
3. the relevant prompt file

## Prompt list

| Slice | Prompt |
|---|---|
| v3.4.0 | `v3.4.0-gauge-type-docs.md` |
| v3.4.1 | `v3.4.1-numeric-rename.md` |
| v3.4.2 | `v3.4.2-odometer-gauge.md` |
| v3.4.3 | `v3.4.3-indicator-gauge.md` |
| v3.4.4 | `v3.4.4-bar-gauge.md` |
| v3.4.5 | `v3.4.5-segmented-gauge.md` |
| v3.4.6 | `v3.4.6-example-asset-framework.md` |
| v3.4.7 | `v3.4.7-ornate-timber-dashboard.md` |
| v3.4.8 | `v3.4.8-neon-grid-dashboard.md` |
| v3.4.9 | `v3.4.9-steam-scrap-dashboard.md` |
| v3.4.10 | `v3.4.10-dashboard-cli.md` |

## Standard workflow

1. Confirm previous relevant PR is merged into `main`.
2. Confirm no blocking open PRs exist.
3. Branch from latest `main`.
4. Implement only the named slice.
5. Update `CHANGES.md` and `docs/v3.4/ImplementationState.md`.
6. Open a PR.
7. Stop.

For v3.4.6 through v3.4.9, keep the work deterministic and local. Do not use remote image generation, downloaded stock art, or hand-edited opaque PNGs as the source of truth.

For v3.4.10, keep the work focused on dashboard CLI shape and related runtime/harness wiring. Do not add new gauge behaviour types or renderer models.

No wandering into the shed and accidentally building a font engine. That way lie tiny layout gremlins with measuring tapes.
