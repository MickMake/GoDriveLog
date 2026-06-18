# GoDriveLog v3.1 migration state

Status: implementation
Last updated: 2026-06-18
State owner: migration implementor

## Purpose

This file tracks the v3.1 release series.

## Baseline from v3.0

v3.0 established the migration process, strict v3 config loading, RuntimePlan resolution, endpoint abstraction, sensor event spine, selected JSONL logging, asset registry, selected dashboard scene rendering, retirement audit, and inverse implementation audit.

## Current migration position

Current version: `v3.1.3`
Current phase: dashboard update performance target
Current branch prefix: `v3.1.3`
Current PR: pending
Current PR branch: `v3.1.3-dashboard-update-performance`

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
- `v3.1.2` added a dashboard harness behind `--v3 --harness` so selected v3 dashboards can be exercised without OBD hardware.
- The harness feeds fake `sensors.SensorEvent` values through the real `v3dashboard.Runtime.ApplyEvent` path and then into the Fyne display adapter.
- Relative v3 asset paths are resolved without a CLI directory flag, using config-dir/vehicle, pwd/vehicle, config-dir, then pwd.
- The harness supports explicit `sweep`, `heartbeat`, and `fixed` patterns and rejects unknown pattern names.
- `v3.1.3` addresses the Raspberry Pi 4 2GB dashboard performance problem by reducing Fyne canvas/image object churn in the visible display path.
- The Fyne adapter now reuses existing `canvas.Image` objects when the rendered part count is stable instead of rebuilding the full object tree every update.
- `v3.1.3` also adds a coalescing dashboard scene sink for the visible v3 paths.
- The scene sink keeps the latest pending scene and drops stale intermediate display frames when rendering cannot keep up.
- Scene sink submission returns when the submitted frame renders, is superseded by a newer frame, or hits a render error, so render errors are still surfaced.
- Sensor polling and JSONL subscribers remain upstream of display rendering; dashboard freshness yields to runtime/logging correctness.
- The v3 window is now sized from selected dashboard config before startup instead of starting at a hard-coded `800x480` and relying on later window resize behaviour.
- Fyne UI work remains inside `fyne.DoAndWait`, and shutdown avoids window manipulation during drain/close to prevent Ctrl-C Fyne thread warnings.

## Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.1.0 | runnable command path | implemented |
| v3.1.1 | display adapter | implemented |
| v3.1.2 | dashboard and gauge test harness | implemented |
| v3.1.3 | dashboard update performance target | in progress |
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
- `v3.1.3-dashboard-update-performance`

## Notes for current slice

The current slice is implementation-only for the v3 dashboard update performance target.

Design intent:

- Fix visible Fyne display memory churn on Raspberry Pi 4 2GB class hardware.
- Keep the v3 schema unchanged.
- Keep dashboards as consumers of v3 dashboard scenes, not direct sensor readers.
- Keep sensor polling and JSONL logging upstream of display rendering.
- Prefer latest visible dashboard state over queued stale display frames.
- Preserve render error propagation instead of hiding display adapter failures.
- Leave deeper dirty-widget rendering, event pipeline optimisation, and remaining sustained-render-backpressure work to `v3.1.7`.

Expected verification focus:

- `go test ./...` passes.
- Coalescing scene sink tests prove stale pending frames are dropped in favour of the latest frame.
- Scene sink submission returns when its frame renders, is superseded, or hits a render error.
- `go run ./cmd/GoDriveLog --v3 --harness --config CONFIG --vehicle VEHICLE_ID --pattern sweep --interval 50ms` can be used as the preferred visual cadence check.
- RSS should not grow rapidly due to repeated full Fyne object-tree rebuilds.
- Ctrl-C should stop the harness/runtime cleanly and avoid Fyne thread warnings.
- `100ms` remains the minimum acceptable fallback if 50ms is not reliable on Raspberry Pi 4 class hardware.
- Dashboard event efficiency remains a later slice for deeper dirty-widget optimisation and any remaining sustained-render-backpressure work.
