# GoDriveLog v3 migration state

Status: active
Last updated: 2026-06-14
State owner: migration implementor

## Purpose

This file is the repo-owned state tracker for the v3 migration.

## Current migration position

Current version: `v3.0.6`
Current phase: selected JSONL logging under review
Current branch prefix: `v3.0.6`
Current PR: `pending`
Current PR branch: `v3.0.6-jsonl-subscriber`

## Current state

- v3.0.0 process scaffolding has been merged.
- Chat prompt workflow has been merged.
- v3.0.0 working-code inventory and seam plan has been merged.
- v3.0.1 frozen v3 docs/schema target has been merged.
- v3.0.2 strict config load and validation has been merged.
- v3.0.3 RuntimePlan resolution has been merged.
- v3.0.4 endpoint abstraction has been merged.
- v3.0.5 sensor event spine and latest-state store has been merged.
- v3.0.6 selected JSONL logging is open for verification.
- Minimal asset registry has not started yet.

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

## Next target

Next version: `v3.0.6`
Next action: verify the v3.0.6 selected JSONL logging PR against the v3.0.6 implementation prompt in `docs/v3/ChatPrompts.md`.

After v3.0.6 is merged, create the v3.0.7 minimal asset registry implementation slice.

## Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.0.0 | working-code inventory and seam plan | complete |
| v3.0.1 | frozen v3 docs and schema target | complete |
| v3.0.2 | strict config load and validation | complete |
| v3.0.3 | RuntimePlan resolution | complete |
| v3.0.4 | endpoint abstraction with serial and TCP simulator support | complete |
| v3.0.5 | sensor event spine and latest-state store | complete |
| v3.0.6 | selected JSONL logging | PR under review |
| v3.0.7 | minimal asset registry: image, digit, indicator | pending |
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

## Notes for current PR

The current PR is a v3.0.6 selected JSONL logging slice.

Expected verification focus:

- internal/logger contains a v3 JSONL event subscriber.
- Selected logs come from RuntimePlan.Logs, resolved from vehicles.<id>.logs.
- Logged sensors come from each selected logs.<id>.sensors definition.
- JSONL output is driven only by sensors.SensorEvent values.
- First readings, value changes, and status transitions are written.
- Unchanged duplicate values do not spam logs.
- Sensor read timestamp and status are included in every event record.
- The logger does not poll sensors.
- The logger does not own sensor cadence.
- The logger does not cause extra endpoint reads.
- It does not implement dashboard rendering.
- It does not change the v3 YAML schema.
