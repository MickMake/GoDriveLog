# GoDriveLog v3.1 implementation guardrails

Status: planning stub
Owner: migration implementor

## Purpose

These guardrails define how v3.1 implementation slices should be planned and reviewed.

## Rules

1. Start every slice from latest `main`.
2. Do not stack work on an unmerged PR unless explicitly approved.
3. Keep each branch focused on one versioned slice.
4. Keep branch names prefixed with the target version.
5. Update `CHANGES.md` for every slice.
6. Update `docs/v3.1/MigrationState.md` for every slice.
7. Prefer small wiring changes over broad rewrites.
8. Preserve old runtime paths until the v3.1 replacement is runnable and verified.
9. Do not add YAML formulas or dashboard-owned sensor reads.
10. Do not let dashboard code access OBD endpoints directly.

## v3.1 implementation focus

The main v3.1 goal is to wire the existing v3 foundation into a practical app path.

The preferred pipeline remains:

```text
selected vehicle -> endpoint -> sensor polling runtime -> events -> selected logs and dashboards
```

## Done means

A slice is done only when it has a clear review target, small diff, updated docs, and no accidental runtime scope outside its purpose.
