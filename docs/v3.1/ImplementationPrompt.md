# GoDriveLog v3.1 implementation prompt

Status: planning stub
Owner: migration implementor

## Purpose

Use this file as the reusable implementation prompt shape for v3.1 slices.

## Required starting checks

1. Confirm the previous v3.1 PR is merged into `main`.
2. Confirm there are no blocking open PRs.
3. Create the next branch from latest `main`.
4. Keep the branch name prefixed with the target version.

## Required slice behaviour

- Keep the implementation small.
- Implement only the named slice.
- Update `CHANGES.md`.
- Update `docs/v3.1/MigrationState.md`.
- Add or update tests only when the slice changes runtime behaviour.
- Do not remove old code unless the slice is explicitly a retirement slice.

## Standard branch format

```text
v3.1.x-short-purpose
```

## Standard completion summary

- Branch created from latest `main`.
- Files changed.
- Behaviour changed or not changed.
- Tests run or not run.
- PR opened.
