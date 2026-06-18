# GoDriveLog v3.1 migration state

Status: implementation
Last updated: 2026-06-18
State owner: migration implementor

## Purpose

This file tracks the v3.1 release series.

## Baseline from v3.0

v3.0 established the migration process, strict v3 config loading, RuntimePlan resolution, endpoint abstraction, sensor event spine, selected JSONL logging, asset registry, selected dashboard scene rendering, retirement audit, and inverse implementation audit.

## Current migration position

Current version: `v3.1.4`
Current phase: JSONL daily rotation
Current branch prefix: `v3.1.4`
Current PR: pending
Current PR branch: `v3.1.4-jsonl-daily-rotation`

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
- `v3.1.4` keeps JSONL logging simple: v3 event logs rotate daily by default.
- The configured log path is treated as a base path, and the active file path inserts the date before the extension.
- Example: `logs/vw_caddy.jsonl` writes to `logs/vw_caddy-2026-06-18.jsonl` for that day.
- No configurable rotation modes, retention policy, compression, upload, or logging architecture rewrite are part of `v3.1.4`.

## Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.1.0 | runnable command path | implemented |
| v3.1.1 | display adapter | implemented |
| v3.1.2 | dashboard and gauge test harness | implemented |
| v3.1.3 | dashboard update performance target | implemented |
| v3.1.4 | JSONL daily rotation | in progress |
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
- `v3.1.4-jsonl-daily-rotation`

## Notes for current slice

The current slice is implementation-only for v3 JSONL daily rotation.

Design intent:

- Keep logging useful for the current van logger use case.
- Use daily JSONL rotation unconditionally for v3 event logs.
- Treat the configured log path as the base path.
- Insert the active date before the extension.
- Roll to a new file when the logger write date changes.
- Keep selected-log and selected-sensor behaviour unchanged.
- Avoid new config schema fields or rotation mode choices.

Expected verification focus:

- `go test ./...` passes.
- Existing selected-log subscriber tests still pass with daily paths.
- Daily path generation is tested.
- JSONL event writer rollover is tested across a date boundary.
- `ActivePath()` reports the concrete daily file path, not the configured base path.
- No dashboard, sensor polling, or v3 schema changes are included.
