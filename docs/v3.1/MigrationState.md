# GoDriveLog v3.1 migration state

Status: implementation
Last updated: 2026-06-17
State owner: migration implementor

## Purpose

This file tracks the v3.1 release series.

## Baseline from v3.0

v3.0 established the migration process, strict v3 config loading, RuntimePlan resolution, endpoint abstraction, sensor event spine, selected JSONL logging, asset registry, selected dashboard scene rendering, retirement audit, and inverse implementation audit.

## Current migration position

Current version: `v3.1.2`
Current phase: dashboard and gauge test harness
Current branch prefix: `v3.1.2`
Current PR: pending
Current PR branch: `v3.1.2-dashboard-gauge-test-harness`

## Current state

- v3.1 planning docs exist under `docs/v3.1/`.
- v3.1 starts from the merged v3.0 foundation.
- v3.1 focuses on the remaining implementation needed for a runnable, visible, independently testable, performant app path.
- `v3.1.0` added the first runnable v3 command path behind `--v3`.
- The old command/runtime path remains the default when `--v3` is not supplied.
- The v3 command path loads v3 config, selects one vehicle, resolves `RuntimePlan`, connects the configured endpoint, starts the sensor polling runtime, wires selected JSONL logs to sensor events, exposes the dashboard scene boundary, and shuts down on SIGINT/SIGTERM.
- `v3.1.1` added a Fyne display adapter for v3 dashboard scene output.
- The adapter consumes dashboard scenes and resolved asset paths; it does not read sensors or OBD endpoints.
- The visible adapter is wired into the `--v3` command path while the old runtime remains available without `--v3`.
- `v3.1.2` adds a dashboard harness behind `--v3 --harness` so selected v3 dashboards can be exercised without OBD hardware.
- The harness feeds fake `sensors.SensorEvent` values through the real `v3dashboard.Runtime.ApplyEvent` path and then into the Fyne display adapter.
- Relative v3 asset paths are resolved without a CLI repo-root flag, using config-dir/vehicle, pwd/vehicle, config-dir, then pwd.
- The harness supports explicit `sweep`, `heartbeat`, and `fixed` patterns and rejects unknown pattern names.

## Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.1.0 | runnable command path | implemented |
| v3.1.1 | display adapter | implemented |
| v3.1.2 | dashboard and gauge test harness | in progress |
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

The current slice is implementation-only for the v3 dashboard and gauge test harness.

Expected verification focus:

- `go test ./...` passes.
- `go run ./cmd/GoDriveLog --v3 --harness --config CONFIG --vehicle VEHICLE_ID --pattern sweep --interval 100ms` opens a Fyne window for selected v3 dashboard output without OBD.
- The harness uses fake sensor events but the real v3 dashboard event/state path.
- The harness uses the real Fyne display adapter rather than a parallel renderer.
- Unknown pattern names are rejected.
- Old runtime remains available without `--v3`.
