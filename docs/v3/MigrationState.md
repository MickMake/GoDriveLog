# GoDriveLog v3 migration state

Status: active  
Last updated: 2026-06-13  
State owner: migration verifier

## 1. Purpose

This file is the repo-owned state tracker for the v3 migration.

Every implementor and verifier chat should read this file before doing work. Do not infer migration position from memory when this file exists.

## 2. Current migration position

Current version: `v3.0.3`  
Current phase: RuntimePlan resolution under review  
Current branch prefix: `v3.0.3`  
Current PR: pending  
Current PR branch: `v3.0.3-runtime-plan`

Current state:

```text
v3.0.0 process scaffolding has been merged.
Chat prompt workflow has been merged.
v3.0.0 working-code inventory and seam plan has been merged.
v3.0.1 frozen docs/schema target has been merged.
v3.0.2 strict config load and validation has been merged.
v3.0.3 RuntimePlan resolution is being prepared for review.
Runtime implementation beyond config resolution has not started yet.
```

## 3. Completed versions

| Version | Status | PR | Notes |
|---|---|---|---|
| v3.0.0 | complete | #36 | Defined versioned migration process, seam plan, branch naming rules, and this state tracker. |
| v3.0.0 | complete | #37 | Added reusable implementation/verification chat prompts for the v3 migration workflow. |
| v3.0.0 | complete | #38 | Added working-code inventory and seam plan before runtime implementation. |
| v3.0.1 | complete | #39 | Added frozen v3 docs and schema target before strict config loading. |
| v3.0.2 | complete | #40 | Added strict v3 config load and validation. |

## 4. Next target

Next version: `v3.0.3`  
Next action: verify the v3.0.3 RuntimePlan resolution PR against the implementation prompt in `docs/v3/ChatPrompts.md`.

If the v3.0.3 PR passes verification and is merged, the next action should be:

```text
Create the v3.0.4 endpoint abstraction implementation slice using the v3.0.4 implementation prompt in docs/v3/ChatPrompts.md.
```

## 5. Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.0.0 | working-code inventory and seam plan | complete |
| v3.0.1 | frozen v3 docs and schema target | complete |
| v3.0.2 | strict v3 config load/validation | complete |
| v3.0.3 | RuntimePlan resolution | branch `v3.0.3-runtime-plan` in progress |
| v3.0.4 | endpoint abstraction with serial/TCP simulator support | pending |
| v3.0.5 | sensor event spine and latest-state store | pending |
| v3.0.6 | selected JSONL logging | pending |
| v3.0.7 | minimal asset registry: image, digit, indicator | pending |
| v3.0.8 | smallest selected dashboard: image + digit_display + indicator | pending |
| v3.0.9 | richer asset registry: bar and frame assets | pending |
| v3.0.10 | richer dashboard widgets: bar_display and frame_gauge | pending |
| v3.0.11 | retire or archive replaced current paths | pending |

## 6. State advancement rule

Implementation proposes state. Verification advances state.

A version is complete only when the verifier has reviewed the PR, required fixes are complete, the PR has been merged or explicitly approved for advancement, and the next target/action is clear.

## 7. Implementor chat workflow

Implementor chats must read this file, read the v3 prompts and guardrails, use a branch starting with the target version, keep work scoped to one migration slice, open a PR to `main`, and not merge it.

## 8. Verifier chat workflow

Verifier chats must read this file, read the v3 prompts, fetch the PR under review, confirm branch/scope/state, and return one of:

```text
PASS - ready to merge
PASS WITH NOTES - acceptable, but follow-up work should be tracked
FAIL - changes required before merge
BLOCKED - cannot verify due to missing information/tooling
```

## 9. Branch naming reminder

Branches for v3 migration work must start with the target version number.

Examples:

```text
v3.0.0-working-code-inventory
v3.0.1-freeze-v3-docs-schema
v3.0.2-config-loader-validation
v3.0.3-runtime-plan
v3.0.4-endpoint-abstraction
```

## 10. Notes for current PR

The current PR is a v3.0.3 RuntimePlan resolution slice.

Expected verification focus:

- `internal/config/v3config.Resolve` exists beside the strict config loader.
- It resolves one selected vehicle into an explicit runtime plan.
- It resolves endpoint config from the selected vehicle.
- It resolves only selected log definitions.
- It resolves only selected dashboard definitions.
- Multiple vehicles require explicit selection.
- A single vehicle can resolve by default.
- A single log/dashboard can be used by default when omitted by the vehicle.
- Display collision validation applies to selected dashboards.
- Unselected dashboards are inert and may share displays as alternatives.
- Runtime consumers can use the resolved plan instead of walking raw config maps.
- It does not implement endpoint connectors.
- It does not implement sensor polling.
- It does not implement dashboard rendering.
- It does not let logs or dashboards choose their own vehicle.
