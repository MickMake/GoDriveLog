# GoDriveLog v3 migration state

Status: active
Last updated: 2026-06-14
State owner: migration implementor

## Purpose

This file is the repo-owned state tracker for the v3 migration.

## Current migration position

Current version: `v3.0.11`
Current phase: retirement audit under review
Current branch prefix: `v3.0.11`
Current PR: `pending`
Current PR branch: `v3.0.11-retirement-audit`

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
- v3.0.8 smallest selected dashboard has been merged.
- v3.0.9 richer asset registry has been merged.
- v3.0.10 richer dashboard widgets have been merged.
- v3.0.11 retirement audit is open for verification.
- No code has been removed, moved, or archived by the v3.0.11 audit slice.

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
| v3.0.8 | complete | #46 | Added smallest selected dashboard scene runtime. |
| v3.0.9 | complete | #47 | Added richer asset registry support for bar and frame assets. |
| v3.0.10 | complete | #48 | Added richer dashboard widget rendering for `bar_display` and `frame_gauge`. |

## Next target

Next version: `v3.0.11`
Next action: verify the v3.0.11 retirement audit PR and review `docs/v3/RetirementAudit.md`.

After v3.0.11 is merged, manually review the retirement audit before creating any removal or archive slices.

## Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.0.0 | working-code inventory and seam plan | complete |
| v3.0.1 | frozen v3 docs and schema target | complete |
| v3.0.2 | strict config load and validation | complete |
| v3.0.3 | RuntimePlan resolution | complete |
| v3.0.4 | endpoint abstraction with serial/TCP simulator support | complete |
| v3.0.5 | sensor event spine and latest-state store | complete |
| v3.0.6 | selected JSONL logging | complete |
| v3.0.7 | minimal asset registry: image, digit, indicator | complete |
| v3.0.8 | smallest selected dashboard: image plus digit_display plus indicator | complete |
| v3.0.9 | richer asset registry: bar and frame assets | complete |
| v3.0.10 | richer dashboard widgets: bar_display and frame_gauge | complete |
| v3.0.11 | retirement audit: review what can be removed or archived later | PR under review |

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
- v3.0.9-richer-asset-registry
- v3.0.10-richer-dashboard-widgets
- v3.0.11-retirement-audit

## Notes for current PR

The current PR is a v3.0.11 docs-only retirement audit slice.

Expected verification focus:

- `docs/v3/RetirementAudit.md` exists.
- The audit identifies old/current paths that may be removed or archived later.
- The audit includes recommendations, confidence, removal conditions, and risks.
- The audit explicitly keeps shared v3 seam/foundation paths.
- The slice does not remove code.
- The slice does not move code.
- The slice does not archive files.
- The slice does not change runtime behaviour.
- The slice does not change tests.
- The slice does not change v3 schema.

Carried follow-up from v3.0.8:

- `ApplyEvent()` currently rebuilds selected dashboard scenes via `Snapshot()` before detecting unchanged rendered output by scene signature. This remains acceptable for the richer widget slice, but future dashboard performance work should avoid rebuilding unaffected widgets or dashboards on every sensor event.
- Do not solve this by adding dashboard polling, YAML formulas, widget-owned sensor reads, or endpoint access from dashboard code.
