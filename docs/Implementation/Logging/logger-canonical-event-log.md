# Canonical GoDriveLog Event Log — Implementation

## Purpose
Audits the current JSONL event log implementation against the design for a canonical GoDriveLog-owned event log.

## Implementation Status
Partially implemented.

Verified current code implements part of the design, but the audited scope also has missing or different behaviour.

## Packages and Files
- `internal/logger/event_jsonl.go`
- `internal/runtime/v3runtime/run.go`
- `internal/config/v3config/config.go`

## Types
- `JSONLEventRecord`
- `JSONLEventWriter`
- `JSONLSubscriber`
- `LogConfig`

## Functions and Methods
- `NewJSONLEventWriter`
- `DailyJSONLPath`
- `NewJSONLSubscribersFromPlan`
- `Run`

## Runtime Flow
`Run` resolves selected logs from the runtime plan, builds `JSONLSubscriber` instances with `NewJSONLSubscribersFromPlan`, and drains live `sensors.SensorEvent` values into `JSONLEventWriter`.

## Configuration
Selected logs use `LogConfig.Path` and `LogConfig.Sensors`. `DailyJSONLPath` rotates the output path by day. The writer uses the configured path and appends a date before the existing extension. No `.gdl.jsonl` requirement or schema/version field was found.

## Behaviour
The logger writes one JSON object per line, records timestamps and typed values, rotates daily, and suppresses unchanged duplicate events per sensor/status/value/error combination.

## Rendering
Not applicable. The feature writes logs only.

## Tests
- `TestNewJSONLSubscribersFromPlanUsesSelectedVehicleLogs`
- `TestJSONLSubscriberWritesSelectedSensorEvents`
- `TestJSONLEventWriterRotatesDaily`
- `TestDailyJSONLPathAddsDateBeforeExtension`
- `TestJSONLSubscriberSuppressesUnchangedDuplicateEvents`
- `TestJSONLSubscriberWritesUnavailableStatusTypedValues`
- `TestRunLoadsResolvedVehicleAndWritesSelectedJSONLLog`

## Limitations
The current format is repository code, not a formally versioned product contract. No schema marker, sidecar metadata, validator, or replay reader was found.

## Deviations from Design
The design requires a canonical `.gdl.jsonl` format that GoDriveLog writes, validates, and replays. Current code only implements the write side of a generic daily-rotated JSONL stream.

## Remaining Work
Add a formal file contract, versioning, validator, sidecar metadata, and replay consumer if the design remains active.

## Verification Notes

Files inspected:
- `internal/logger/event_jsonl.go`
- `internal/runtime/v3runtime/run.go`
- `internal/config/v3config/config.go`

Symbols verified:
- `JSONLEventRecord`
- `JSONLEventWriter`
- `JSONLSubscriber`
- `LogConfig`
- `NewJSONLEventWriter`
- `DailyJSONLPath`
- `NewJSONLSubscribersFromPlan`
- `Run`

Configuration verified:
- `path`
- `sensors`

Tests inspected:
- `TestNewJSONLSubscribersFromPlanUsesSelectedVehicleLogs`
- `TestJSONLSubscriberWritesSelectedSensorEvents`
- `TestJSONLEventWriterRotatesDaily`
- `TestDailyJSONLPathAddsDateBeforeExtension`
- `TestJSONLSubscriberSuppressesUnchangedDuplicateEvents`
- `TestJSONLSubscriberWritesUnavailableStatusTypedValues`
- `TestRunLoadsResolvedVehicleAndWritesSelectedJSONLLog`

Searches performed:
- `gdl.jsonl`
- `JSONLEventRecord`
- `DailyJSONLPath`
- `JSONLSubscriber`
