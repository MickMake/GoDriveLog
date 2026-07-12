# Canonical GoDriveLog Event Log

Design reference: [`docs/Designs/Logging/logger-canonical-event-log.md`](../../Designs/Logging/logger-canonical-event-log.md)

## Purpose
Records how far the current JSONL event stream has progressed toward a formal GoDriveLog-owned log format.

## Implementation Status
Status: **Partially implemented**.

Structured JSONL event output exists, but it is not yet a branded, versioned, replay-ready canonical format.

## Packages and Files
- [`internal/logger/event_jsonl.go`](../../../internal/logger/event_jsonl.go)
- [`internal/runtime/v3runtime/run.go`](../../../internal/runtime/v3runtime/run.go)

## Types
- `JSONLEventRecord`
- `JSONLEventWriter`
- `JSONLSubscriber`

## Functions and Methods
- `NewJSONLEventWriter`
- `DailyJSONLPath`
- `Run`

## Runtime Flow
`v3runtime.Run` can construct JSONL subscribers from the resolved plan and feed them live sensor and status updates during a run.

## Configuration
Logging is selected from runtime plan outputs. The emitted filenames rotate daily, but the code still writes generic `.jsonl` files instead of a dedicated GoDriveLog extension.

## Behaviour
Events are written as newline-delimited JSON records with typed values, timestamps, sensor IDs, and duplicate suppression for unchanged events.

## Rendering
This feature does not affect dashboard rendering directly.

## Tests
- [`internal/logger/event_jsonl_test.go`](../../../internal/logger/event_jsonl_test.go)
- [`internal/logger/event_jsonl_status_semantics_test.go`](../../../internal/logger/event_jsonl_status_semantics_test.go)
- [`internal/runtime/v3runtime/run_test.go`](../../../internal/runtime/v3runtime/run_test.go)

## Limitations
There is no explicit schema version, no `.gdl.jsonl` contract, no replay-facing guarantees, and no provenance sidecar.

## Deviations from Design
The design promotes the logger output into a formal product format. Current code provides the raw stream but not the finished contract.

## Remaining Work
Add schema/version markers, canonical extension and naming rules, metadata sidecar support, validator support, and replay consumers.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
