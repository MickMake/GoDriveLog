# GoDriveLog v3 migration state

Status: active  
Last updated: 2026-06-13  
State owner: migration verifier

## 1. Purpose

This file is the repo-owned state tracker for the v3 migration.

Every implementor and verifier chat should read this file before doing work. Do not infer migration position from memory when this file exists.

## 2. Current migration position

Current version: `v3.0.2`  
Current phase: strict v3 config load and validation under review  
Current branch prefix: `v3.0.2`  
Current PR: `#40`  
Current PR branch: `v3.0.2-config-loader-validation`

Current state:

```text
v3.0.0 process scaffolding has been merged.
Chat prompt workflow has been merged.
v3.0.0 working-code inventory and seam plan has been merged.
v3.0.1 frozen docs/schema target has been merged.
v3.0.2 strict config load and validation is open for verification.
Runtime implementation work beyond config validation has not started yet.
```

## 3. Completed versions

| Version | Status | PR | Notes |
|---|---|---|---|
| v3.0.0 | complete | #36 | Defined versioned migration process, seam plan, branch naming rules, and this state tracker. |
| v3.0.0 | complete | #37 | Added reusable implementation/verification chat prompts for the v3 migration workflow. |
| v3.0.0 | complete | #38 | Added working-code inventory and seam plan before runtime implementation. |
| v3.0.1 | complete | #39 | Added frozen v3 docs and schema target before strict config loading. |

## 4. Next target

Next version: `v3.0.2`  
Next action: verify PR `#40` against the v3.0.2 strict v3 config load and validation prompt.

If PR `#40` passes verification and is merged, the next action should be:

```text
Create the v3.0.3 RuntimePlan resolution implementation slice using the v3.0.3 implementation prompt in docs/v3/ChatPrompts.md.
```

## 5. Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.0.0 | working-code inventory and seam plan | complete |
| v3.0.1 | frozen v3 docs and schema target | complete |
| v3.0.2 | strict v3 config load/validation | PR #40 under review |
| v3.0.3 | RuntimePlan resolution | pending |
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

The current PR is a v3.0.2 strict config load and validation slice.

Expected verification focus:

- `internal/config/v3config` exists beside the current config path.
- It implements strict v3 YAML loading without changing the old/current config loader.
- It rejects unknown root and nested fields.
- It validates docs/v3/config.example.yaml.
- It validates docs/v3/config.full.yaml.
- It validates active docs/v3/examples/*.yaml files.
- It validates vehicle, log, dashboard, widget, and asset references.
- It keeps sensors and assets as global catalogues.
- It does not wire the full runtime.
- It does not add compatibility aliases.
- It does not auto-convert old/current config shapes.
- It does not start v3.0.3 RuntimePlan implementation.
