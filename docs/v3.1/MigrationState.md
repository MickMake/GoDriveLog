# GoDriveLog v3.1 migration state

Status: implementation
Last updated: 2026-06-18
State owner: migration implementor

## Purpose

This file tracks the v3.1 release series.

## Baseline from v3.0

v3.0 established the migration process, strict v3 config loading, RuntimePlan resolution, endpoint abstraction, sensor event spine, selected JSONL logging, asset registry, selected dashboard scene rendering, retirement audit, and inverse implementation audit.

## Current migration position

Current version: `v3.1.5`
Current phase: typed sensor values
Current branch prefix: `v3.1.5`
Current PR: pending
Current PR branch: `v3.1.5-typed-sensor-values`

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
- `v3.1.5` introduces an explicit typed sensor value payload alongside the existing numeric compatibility field.
- Typed sensor values require a non-empty supported `kind`; the zero value `Value{}` is invalid and must not be treated as a valid sensor reading.
- Sensor config may declare `value_kind`; if omitted, the runtime derives it from the selected parser/output contract.
- OBD parser output currently derives `numeric`; configured `value_kind` must match that parser output instead of relying on hard-coded sensor-id knowledge.
- JSONL event records now write a typed `value` object instead of a bare number.

## Version queue

| Version | Purpose | Status |
|---|---|---|
| v3.1.0 | runnable command path | implemented |
| v3.1.1 | display adapter | implemented |
| v3.1.2 | dashboard and gauge test harness | implemented |
| v3.1.3 | dashboard update performance target | implemented |
| v3.1.4 | JSONL daily rotation | implemented |
| v3.1.5 | typed sensor values | in progress |
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
- `v3.1.5-typed-sensor-values`

## Notes for current slice

The current slice is implementation-only for v3 typed sensor values.

Design intent:

- Make sensor value type explicit with a mandatory `kind`.
- Reject empty or unknown typed value kinds.
- Do not infer value type from populated fields.
- Avoid hard-coding sensor IDs such as `rpm` or `speed` to kinds.
- Derive default kind from parser/output contract; for current OBD readers this is `numeric`.
- Fail selected runtime setup when configured `value_kind` is invalid or incompatible with parser output.
- Preserve runtime operation for live read/parse errors by emitting typed error values/events instead of coercing to zero.
- Keep unsupported/missing/stale visual policy for `v3.1.6`.

Expected verification focus:

- `go test ./...` passes.
- Numeric sensor reads produce typed numeric values.
- Numeric zero remains a valid numeric value, not missing/error.
- `Value{}` and unknown kinds are invalid.
- Invalid configured `value_kind` fails selected runtime setup.
- Configured kind mismatch with parser output fails selected runtime setup.
- JSONL records include typed value objects.
