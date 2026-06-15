# GoDriveLog v3.1 verification prompt

Status: planning
Owner: migration implementor

## Purpose

Use this file as the reusable review checklist for v3.1 slices.

## Required checks

1. Confirm the PR branch name matches the target version.
2. Confirm the PR targets `main`.
3. Confirm the slice scope matches `ReleasePlan.md`.
4. Confirm the slice follows its section in `ImplementationPrompt.md`.
5. Confirm `CHANGES.md` is updated.
6. Confirm `docs/v3.1/MigrationState.md` is updated.
7. Confirm `docs/v3.1/OpenDecisions.md` is updated only when the slice resolves, changes, adds, or explicitly defers a decision.
8. Confirm tests were added or run when behaviour changed.
9. Confirm old paths are not removed unless the slice explicitly allows retirement.
10. Confirm dashboard code does not read sensors or endpoints directly.
11. Confirm config remains data, not embedded scripting.
12. Confirm a slice chat did not redesign the release plan without planning approval.

## Extra checks for v3.1.1

- Dummy data can exercise dashboard, gauge, widget, or display output without vehicle hardware.
- Sweep pattern covers min to max to min over 10 seconds.
- Heartbeat pattern exists or is clearly scaffolded.
- 50ms and 100ms cadence targets are represented or planned.
- `OpenDecisions.md` is updated only if the harness shape or cadence decision changes.

## Extra checks for v3.1.4

- 50ms preferred and 100ms minimum dashboard update targets are documented.
- Raspberry Pi 4 remains the reference hardware target.
- Display rendering is not designed to block polling or logging.
- `OpenDecisions.md` is updated if cadence targets are resolved, changed, or explicitly deferred.

## Review summary format

- Target version:
- Scope match: yes/no
- Behaviour changed: yes/no
- Tests: listed or not applicable
- Docs updated: yes/no
- OpenDecisions changed or not applicable:
- Risks found: list
- Recommendation: approve, request changes, or comment
