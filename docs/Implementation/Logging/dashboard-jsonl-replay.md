# JSONL Dashboard Replay

Design reference: [`docs/Designs/Logging/dashboard-jsonl-replay.md`](../../Designs/Logging/dashboard-jsonl-replay.md)

## Purpose
Tracks the missing replay mode for feeding recorded event logs back through the dashboard runtime.

## Implementation Status
Status: **Not implemented**.

No `dashboard replay` command or file-backed replay pipeline exists on `main`.

## Packages and Files
- [`cmd/GoDriveLog/main_ebiten.go`](../../../cmd/GoDriveLog/main_ebiten.go)
- [`internal/runtime/v3runtime/run.go`](../../../internal/runtime/v3runtime/run.go)

## Types
- None in current code.

## Functions and Methods
- `Run` starts live runtime wiring only; there is no replay entrypoint.
- `main` exposes `dashboard run`, `harness`, `examples`, `validate`, and `preview`, but not `replay`.

## Runtime Flow
The current runtime only subscribes to live sensor updates and selected JSONL logging sinks. No component reads a `.jsonl` log and re-emits events into the dashboard path.

## Configuration
There is no replay-specific CLI flag, config struct, or log source selection for recorded sessions.

## Behaviour
Recorded sessions cannot currently drive the dashboard without an external custom tool.

## Rendering
Rendering remains tied to live runtime events or preview/harness flows, not persisted logs.

## Tests
- [`cmd/GoDriveLog/main_ebiten_test.go`](../../../cmd/GoDriveLog/main_ebiten_test.go)

## Limitations
The design depends on a canonical event log format, replay clocking, and CLI plumbing that are all absent.

## Deviations from Design
The design calls this a core validation feature. Current code stops at live logging.

## Remaining Work
Add a replay command, log reader, event scheduler, and dashboard-runtime adapter for recorded sessions.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
