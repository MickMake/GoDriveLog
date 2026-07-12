# JSONL Dashboard Replay — Implementation

## Purpose
Audits whether the repository can replay recorded JSONL events back through the dashboard runtime.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `cmd/GoDriveLog/main_ebiten.go`
- `cmd/GoDriveLog/v3_preview_ebiten.go`
- `internal/runtime/v3runtime/run.go`

## Types
None found in current code.

## Functions and Methods
- `runDashboardCLI`
- `runDashboardPreviewCommand`
- `Run`

## Runtime Flow
No feature-specific runtime path was found. `Run` connects live readers, polling runtime subscribers, and the dashboard sink; it does not read a log file and emit replayed `sensors.SensorEvent` values.

## Configuration
No `dashboard replay` subcommand, `--log` flag, replay config struct, or replay source selection was found. The only replay-related code in the CLI is preview-local `replayPending` state for repeating the last manual preview transition.

## Behaviour
Recorded JSONL logs cannot be played back through the dashboard command tree in current code.

## Rendering
Not applicable. No log-replay rendering path was found.

## Tests
No feature-specific tests found.

## Limitations
This audit only covers the current repository. It does not treat design notes or historical plans as implementation evidence.

## Deviations from Design
The design calls for a dashboard replay mode. Current code exposes `dashboard run`, `dashboard harness`, `dashboard preview`, `dashboard examples`, and `dashboard validate`, but not replay.

## Remaining Work
Add a replay command, a log reader, event-to-runtime wiring, and replay-specific tests if this design is still wanted.

## Verification Notes

Files inspected:
- `cmd/GoDriveLog/main_ebiten.go`
- `cmd/GoDriveLog/v3_preview_ebiten.go`
- `internal/runtime/v3runtime/run.go`

Symbols verified:
- `runDashboardCLI`
- `runDashboardPreviewCommand`
- `Run`
- `replayPending`

Tests inspected:
- `TestDashboardHelpOutputsIncludeNewCommandTree`
- `TestDashboardPreviewAcceptsFileBeforeOrAfterFlags`

Searches performed:
- `dashboard replay`
- `replayPending`
- `.gdl.jsonl`
