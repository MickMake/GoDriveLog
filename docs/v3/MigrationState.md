# GoDriveLog v3 migration state

Status: active
Last updated: 2026-06-14
State owner: migration implementor

## Purpose

This file is the repo-owned state tracker for the v3 migration.

## Current migration position

Current version: `v3.0.7`
Current phase: minimal asset registry under review
Current branch prefix: `v3.0.7`
Current PR: `pending`
Current PR branch: `v3.0.7-asset-registry`

## Current state

- v3.0.0 process scaffolding has been merged.
- Chat prompt workflow has been merged.
- v3.0.0 working-code inventory and seam plan has been merged.
- v3.0.1 frozen v3 docs/schema target has been merged.
- v3.0.2 strict config load and validation has been merged.
- v3.0.3 RuntimePlan resolution has been merged.
- v3.0.4 endpoint abstraction has been merged.
- v3.0.5 sensor event spine and latest-state store has been merged.
- v3.0.6 selected JSONL logging has been merged.
- v3.0.7 minimal asset registry is open for verification.
- Smallest selected dashboard has not started yet.

## Completed versions

| Version | Status | PR | Notes |
|---|---|---|---|
| v3.0.0 | complete | #36 | Defined versioned migration process, seam plan, branch naming rules, and this state tracker. |
| v3.0.0 | complete | #37 | Added reusable implementation and verification chat prompts for the v3 migration workflow. |
| v3.0.0 | complete | #38 | Added working-code inventory and seam plan before runtime implementation. |
| v3.0.1 | complete | #39 | Added frozen v3 docs and schema target before strict config loading. |
| v3.0.2 | complete | #40 | Added strict v3 config load and validation. |
| v3.0.3 | complete | #41 | Added RuntimePlan resolution. |
| v3.0.4 | complete | #42 | Added endpoint abstraction with serial and TCP simulator support. |
| v3.0.5 | complete | #43 | Added sensor event spine and latest-state store. |
| v3.0.6 | complete | #44 | Added selected JSONL logging. |

## Next target

Next version: `v3.0.7`
Next action: verify the v3.0.7 minimal asset registry PR against the v3.0.7 implementation prompt in `docs/v3/ChatPrompts.md`.

After v3.0.7 is merged, create the v3.0.8 smallest selected dashboard implementation slice.

## Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.0.0 | working-code inventory and seam plan | complete |
| v3.0.1 | frozen v3 docs and schema target | complete |
| v3.0.2 | strict config load and validation | complete |
| v3.0.3 | RuntimePlan resolution | complete |
| v3.0.4 | endpoint abstraction with serial and TCP simulator support | complete |
| v3.0.5 | sensor event spine and latest-state store | complete |
| v3.0.6 | selected JSONL logging | complete |
| v3.0.7 | minimal asset registry: image, digit, indicator | PR under review |
| v3.0.8 | smallest selected dashboard: image plus digit_display plus indicator | pending |
| v3.0.9 | richer asset registry: bar and frame assets | pending |
| v3.0.10 | richer dashboard widgets: bar_display and frame_gauge | pending |
| v3.0.11 | retire or archive replaced current paths | pending |

## Branch naming reminder

Branches for v3 migration work must start with the target version number.

Examples:

- v3.0.0-working-code-inventory
- v3.0.1-freeze-v3-docs-schema
- v3.0.2-config-loader-validation
- v3.0.3-runtime-plan
- v3.0.4-endpoint-abstraction
- v3.0.5-sensor-event-spine
- v3.0.6-jsonl-subscriber
- v3.0.7-asset-registry

## Notes for current PR

The current PR is a v3.0.7 minimal asset registry slice.

Expected verification focus:

- `internal/assets` contains a v3 asset registry for image, digit, and indicator assets.
- Asset paths resolve as repository-root relative paths.
- Missing assets fail clearly.
- Required indicator states `off`, `on`, and `unknown` are loaded and validated by the registry boundary.
- Digit character images, decimal point, and optional layer images can be decoded and reused.
- Decoded images are loaded once into reusable registry structs for later widgets/renderers.
- The registry does not implement bar or frame assets.
- The registry does not implement dashboard rendering.
- The registry does not add vehicle-owned asset lists.
- The registry does not add YAML rules, scripts, formulas, or conditions.
