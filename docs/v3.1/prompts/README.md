# GoDriveLog v3.1 implementation prompts

Use one prompt file per implementation slice.

Planning work is not a numbered v3.1 slice. Implementation starts at `v3.1.0`.

## Prompt index

- `v3.1.0-dashboard-gauge-test-harness.md`
- `v3.1.1-runnable-command-path.md`
- `v3.1.2-display-adapter.md`
- `v3.1.3-dashboard-update-performance.md`
- `v3.1.4-jsonl-rotation-decision.md`
- `v3.1.5-typed-sensor-values.md`
- `v3.1.6-unsupported-missing-semantics.md`
- `v3.1.7-dashboard-event-efficiency.md`
- `v3.1.8-retirement-readiness.md`

## Branch-chat workflow

A slice implementation chat must:

1. Confirm the target version.
2. Read `docs/v3.1/ReleasePlan.md`.
3. Read this file.
4. Read the matching prompt file in this directory.
5. Confirm the previous relevant PR is merged into `main`.
6. Confirm there are no blocking open PRs.
7. Create a branch from latest `main` using the target version prefix.
8. Implement only the named slice.
9. Update `CHANGES.md`.
10. Update `docs/v3.1/MigrationState.md`.
11. Update `docs/v3.1/OpenDecisions.md` only when the slice resolves, changes, adds, or explicitly defers a decision.
12. Open a PR.
13. Stop.

Do not redesign the release plan inside a slice chat.

## Required docs update per slice

Every implementation slice must update:

- `CHANGES.md`
- `docs/v3.1/MigrationState.md`

A slice must also update `docs/v3.1/OpenDecisions.md` when it:

- resolves a decision
- changes a default position
- adds a new decision
- explicitly defers an existing decision with a reason

Do not edit `OpenDecisions.md` just to say there was no change.

## Global implementation guardrails

- Keep the implementation small.
- Start every slice from latest `main`.
- Do not stack work on an unmerged PR unless explicitly approved.
- Keep each branch focused on one versioned slice.
- Keep branch names prefixed with the target version.
- Use existing v3 foundation packages before adding new foundations.
- Prefer small wiring changes over broad rewrites.
- Keep config as data.
- Do not add YAML formulas.
- Do not let dashboard code access OBD endpoints.
- Do not let widgets own sensor reads.
- Do not remove old code unless the target slice explicitly allows retirement work.
- Preserve old runtime paths until the v3.1 replacement is runnable and verified.
- Add or update tests when runtime behaviour changes.

## v3.1 implementation focus

The main v3.1 goal is to wire the existing v3 foundation into a practical app path.

The preferred pipeline remains:

```text
selected vehicle -> endpoint -> sensor polling runtime -> events -> selected logs and dashboards
```

## Review checklist

Before approving a slice PR, confirm:

1. The PR branch name matches the target version.
2. The PR targets `main`.
3. The slice scope matches `ReleasePlan.md`.
4. The slice follows its prompt file.
5. `CHANGES.md` is updated.
6. `docs/v3.1/MigrationState.md` is updated.
7. `docs/v3.1/OpenDecisions.md` is updated only when the slice resolves, changes, adds, or explicitly defers a decision.
8. Tests were added or run when behaviour changed.
9. Old paths are not removed unless the slice explicitly allows retirement.
10. Dashboard code does not read sensors or endpoints directly.
11. Config remains data, not embedded scripting.
12. The slice chat did not redesign the release plan without planning approval.

## Review summary format

- Target version:
- Scope match: yes/no
- Behaviour changed: yes/no
- Tests: listed or not applicable
- Docs updated: yes/no
- OpenDecisions changed or not applicable:
- Risks found: list
- Recommendation: approve, request changes, or comment

## Done means

A slice is done only when it has a clear review target, small diff, updated docs, and no accidental runtime scope outside its purpose.