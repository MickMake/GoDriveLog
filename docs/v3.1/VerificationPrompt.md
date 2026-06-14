# GoDriveLog v3.1 verification prompt

Status: planning stub
Owner: migration implementor

## Purpose

Use this file as the reusable review checklist for v3.1 slices.

## Required checks

1. Confirm the PR branch name matches the target version.
2. Confirm the PR targets `main`.
3. Confirm the slice scope matches the planned version.
4. Confirm `CHANGES.md` is updated.
5. Confirm `docs/v3.1/MigrationState.md` is updated.
6. Confirm tests were added or run when runtime behaviour changed.
7. Confirm old paths are not removed unless the slice explicitly allows retirement.
8. Confirm dashboard code does not read sensors or endpoints directly.
9. Confirm config remains data, not embedded scripting.

## Review summary format

- Scope match: yes/no
- Runtime behaviour changed: yes/no
- Tests: listed or not applicable
- Docs updated: yes/no
- Risks found: list
- Recommendation: approve, request changes, or comment
