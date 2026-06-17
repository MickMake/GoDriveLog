# GoDriveLog v3.1 migration state

Status: implementation
Last updated: 2026-06-17
State owner: migration implementor

## Purpose

This file tracks the v3.1 release series.

## Baseline from v3.0

v3.0 established the migration process, strict v3 config loading, RuntimePlan resolution, endpoint abstraction, sensor event spine, selected JSONL logging, asset registry, selected dashboard scene rendering, retirement audit, and inverse implementation audit.

## Current migration position

Current version: `v3.1.1`
Current phase: display adapter
Current branch prefix: `v3.1.1`
Current PR: pending
Current PR branch: `v3.1.1-display-adapter`

## Current state

- v3.1 planning docs exist under `docs/v3.1/`.
- v3.1 starts from the merged v3.0 foundation.
- v3.1 focuses on the remaining implementation needed for a runnable, visible, independently testable, performant app path.
- `v3.1.0` added the first runnable v3 command path behind `--v3`.
- The old command/runtime path remains the default when `--v3` is not supplied.
- The v3 command path loads v3 config, selects one vehicle, resolves `RuntimePlan`, connects the configured endpoint, starts the sensor polling runtime, wires selected JSONL logs to sensor events, exposes the dashboard scene boundary, and shuts down on SIGINT/SIGTERM.
- `v3.1.1` adds a Fyne display adapter for v3 dashboard scene output.
- The adapter consumes dashboard scenes and repo-relative asset paths; it does not read sensors or OBD endpoints.
- The visible adapter is wired into the `--v3` command path while the old runtime remains available without `--v3`.

## Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.1.0 | runnable command path | implemented |
| v3.1.1 | display adapter | in progress |
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

The current slice is implementation-only for the v3 display adapter.

Expected verification focus:

- `go test ./...` passes.
- `go run ./cmd/GoDriveLog --v3 --config CONFIG --vehicle VEHICLE_ID --repo-root REPO_ROOT` opens a Fyne window for selected v3 dashboard output.
- The display adapter consumes v3 dashboard scenes and repo-root-relative asset paths.
- Display code does not read sensors or OBD endpoints.
- Dashboard runtime remains responsible for scene generation.
- Old runtime remains available without `--v3`.
