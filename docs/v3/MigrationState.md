# GoDriveLog v3 migration state

Status: active
Last updated: 2026-06-14
State owner: migration implementor

## Purpose

This file is the repo-owned state tracker for the v3 migration.

## Current migration position

Current version: `v3.0.8`
Current phase: smallest selected dashboard under review
Current branch prefix: `v3.0.8`
Current PR: `pending`
Current PR branch: `v3.0.8-smallest-dashboard`

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
- v3.0.7 minimal asset registry has been merged.
- v3.0.8 smallest selected dashboard is open for verification.
- Richer asset registry has not started yet.

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
| v3.0.7 | complete | #45 | Added minimal asset registry. |

## Next target

Next version: `v3.0.8`
Next action: verify the v3.0.8 smallest selected dashboard PR against the v3.0.8 implementation prompt in `docs/v3/ChatPrompts.md`.

After v3.0.8 is merged, create the v3.0.9 richer asset registry implementation slice.

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
| v3.0.7 | minimal asset registry: image, digit, indicator | complete |
| v3.0.8 | smallest selected dashboard: image plus digit_display plus indicator | PR under review |
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
- v3.0.8-smallest-dashboard

## Notes for current PR

The current PR is a v3.0.8 smallest selected dashboard slice.

Expected verification focus:

- `internal/dashboard/v3dashboard` contains a v3 selected-dashboard scene runtime.
- Selected dashboards come from `v3config.RuntimePlan.Dashboards`, resolved from `vehicles.<id>.dashboards`.
- Image, digit display, and indicator widgets render from the v3 asset registry.
- Dashboard state is driven by `sensors.SensorState` and `sensors.SensorEvent` values.
- Dashboard code does not read endpoint, OBD, or sensor reader code directly.
- Static image widgets can render.
- Numeric values render through digit display slots.
- Decimal separators do not consume digit slots.
- Indicators use `unknown` for stale, error, missing, or otherwise non-`ok` sensor status.
- Unchanged formatted digit output does not trigger a changed scene result.
- The slice does not implement `bar_display` or `frame_gauge`.
- The slice does not add dashboard-level polling cadence.
- The slice does not add YAML rules, scripts, formulas, or conditions.

Verification follow-up:

- PR #46 is acceptable for the v3.0.8 smallest selected dashboard slice.
- `ApplyEvent()` currently rebuilds selected dashboard scenes via `Snapshot()` before detecting unchanged rendered output by scene signature. This is acceptable for the first dashboard seam, but later dashboard work should avoid rebuilding unaffected widgets or dashboards on every sensor event.
- Track this before or during v3.0.10 richer dashboard widgets, or in any earlier performance-focused dashboard refinement.
- Do not solve this by adding dashboard polling, YAML formulas, widget-owned sensor reads, or endpoint access from dashboard code.
