# GoDriveLog v3 migration state

Status: active  
Last updated: 2026-06-12  
State owner: migration verifier

## 1. Purpose

This file is the repo-owned state tracker for the v3 migration.

Every implementor and verifier chat should read this file before doing work.

It answers:

- where the migration is up to
- what versioned slice is current
- what branch prefix should be used
- which PR is in review
- what should happen next

Do not infer migration position from memory when this file exists. Read this file first.

## 2. Current migration position

Current version: `v3.0.0`  
Current phase: chat prompt setup for phone-friendly implementation/verification workflow  
Current branch prefix: `v3.0.0`  
Current PR: pending  
Current PR branch: `v3.0.0-chat-prompts`

Current state:

```text
v3.0.0 process scaffolding has been merged.
Chat prompts are being added before implementation work starts.
Implementation work has not started yet.
```

## 3. Completed versions

| Version | Status | PR | Notes |
|---|---|---|---|
| v3.0.0 | partial | #36 | Defined versioned migration process, seam plan, branch naming rules, and this state tracker. |

## 4. Next target

Next version: `v3.0.0`  
Next action: verify and merge the chat prompts PR.

After the chat prompts PR is merged, the next action should be:

```text
Create the v3.0.0 working-code inventory and seam-plan implementation slice using the v3.0.0 implementation prompt in docs/v3/ChatPrompts.md.
```

That slice should map current config, runtime, OBD, logging, dashboard, renderer, and asset code to v3 roles, then record reuse/refactor/wrap/replace/archive decisions.

## 5. Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.0.0 | working-code inventory and seam plan | prompts being added |
| v3.0.1 | frozen v3 docs and schema target | pending |
| v3.0.2 | strict v3 config load/validation | pending |
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

Implementor chats may propose updates to this file.

Verifier chats decide whether migration state actually advances.

Do not mark a version complete just because a PR exists.

A version is complete only when:

- the verifier has reviewed the PR against the target version
- required fixes are complete
- the PR has been merged or the verifier explicitly says the state may advance
- the next target version/action is clear

Useful rule:

```text
Implementation proposes state.
Verification advances state.
```

## 7. Implementor chat workflow

Implementor chats must:

1. Read this file first.
2. Read `docs/v3/ChatPrompts.md`.
3. Read `docs/v3/MigrationGuardrails.md`.
4. Use the matching implementation prompt from `docs/v3/ChatPrompts.md`.
5. Confirm the current target version.
6. Use a branch name starting with the target version.
7. Work only on the current or next target version.
8. Avoid later-version work unless explicitly required and documented.
9. Update this file in the PR only if the migration state changes or a new PR enters review.
10. Open a PR to `main`.
11. Do not merge the PR.

Before coding, implementor chats should report:

- current target version
- branch name
- files/docs read
- scope
- non-goals
- expected tests or docs-only justification

After coding, implementor chats should report:

- files changed
- tests added or updated
- tests run
- state-file changes, if any
- PR number and branch
- known deferrals

## 8. Verifier chat workflow

Verifier chats must:

1. Read this file first.
2. Read `docs/v3/ChatPrompts.md`.
3. Fetch the PR under review.
4. Confirm the PR branch starts with the target version.
5. Confirm the PR scope matches the relevant implementation prompt.
6. Confirm the PR does not perform later-version work without justification.
7. Confirm this file was updated correctly, if the PR changes migration state.
8. Return a clear verdict.

Verifier verdicts:

```text
PASS — ready to merge
PASS WITH NOTES — acceptable, but follow-up work should be tracked
FAIL — changes required before merge
BLOCKED — cannot verify due to missing information/tooling
```

Verifier chats should include:

- verdict
- target version
- PR branch
- files reviewed
- checks passed
- issues found
- required fixes, if any
- whether `MigrationState.md` should advance
- next target version/action if the PR passes

## 9. Branch naming reminder

Branches for v3 migration work must start with the target version number.

Examples:

```text
v3.0.0-docs-migration-seams
v3.0.0-chat-prompts
v3.0.0-working-code-inventory
v3.0.2-config-loader-validation
v3.0.3-runtime-plan
v3.0.4-endpoint-abstraction
v3.0.8-smallest-dashboard
```

If the target version is unclear, decide the target version before creating the branch.

## 10. Notes for current PR

The current chat prompts PR is process setup, not implementation.

It should be verified as `v3.0.0` scaffolding.

Expected verification focus:

- `docs/v3/ChatPrompts.md` exists.
- It contains stage-specific implementation prompts.
- It contains one generic verifier prompt.
- It supports short phone-friendly commands.
- `MigrationState.md` tells future chats to use `ChatPrompts.md`.
- Implementation work has not started prematurely.
