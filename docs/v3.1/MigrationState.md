# GoDriveLog v3.1 migration state

Status: planning
Last updated: 2026-06-15
State owner: migration implementor

## Purpose

This file tracks the v3.1 release series.

## Baseline from v3.0

v3.0 established the migration process, strict v3 config loading, RuntimePlan resolution, endpoint abstraction, sensor event spine, selected JSONL logging, asset registry, selected dashboard scene rendering, retirement audit, and inverse implementation audit.

## Current migration position

Current version: `v3.1.0`
Current phase: release planning stubs
Current branch prefix: `v3.1.0`
Current PR: pending
Current PR branch: `v3.1.0-release-planning-stubs`

## Current state

- v3.1 planning docs are being introduced under `docs/v3.1/`.
- v3.1 starts from the merged v3.0 foundation.
- v3.1 focuses on the remaining implementation needed for a runnable, visible, independently testable, performant app path.
- v3.1 uses one prompt file per planned implementation slice under `docs/v3.1/prompts/`.
- Each future branch chat should implement exactly one version slice from the plan.
- No runtime, test, schema, archive, move, or deletion changes are part of this planning PR.

## Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.1.0 | runnable command path | planned |
| v3.1.1 | display adapter | planned |
| v3.1.2 | dashboard and gauge test harness | planned |
| v3.1.3 | dashboard update performance target | planned |
| v3.1.4 | JSONL rotation decision | planned |
| v3.1.5 | typed sensor values | planned |
| v3.1.6 | unsupported and missing sensor semantics | planned |
| v3.1.7 | dashboard event efficiency | planned |
| v3.1.8 | retirement readiness review | planned |

## Branch naming reminder

Branches for v3.1 work must start with the target version number.

Examples:

- `v3.1.0-release-planning-stubs`
- `v3.1.0-runnable-command-path`
- `v3.1.1-display-adapter`
- `v3.1.2-dashboard-gauge-test-harness`

## Notes for current slice

The current slice is docs-only.

Expected verification focus:

- `docs/v3.1/` exists.
- Planning stub files exist.
- The structure follows the proven `docs/v3/` style.
- Completed v3.0 history is summarised, not bulk-copied.
- `docs/v3.1/prompts/` contains one prompt file per planned v3.1 implementation slice.
- The slice does not change code, tests, runtime behaviour, archives, or schema.
