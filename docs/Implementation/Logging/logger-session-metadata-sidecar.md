# Session Metadata Sidecar — Implementation

## Purpose
Audits whether JSONL log writes also create a session metadata sidecar.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `internal/logger/event_jsonl.go`
- `internal/runtime/v3runtime/run.go`

## Types
- `JSONLEventWriter`
- `JSONLSubscriber`

## Functions and Methods
- `NewJSONLEventWriter`
- `WriteEvent`
- `Run`

## Runtime Flow
`Run` builds JSONL subscribers and drains events into `WriteEvent`. No second writer, sidecar path, or metadata emission path was found.

## Configuration
`LogConfig` exposes `path` and `sensors` only. No sidecar path or metadata fields were found.

## Behaviour
Current code writes only the event stream.

## Rendering
Not applicable.

## Tests
No feature-specific tests found.

## Limitations
No repository code was found for metadata capture, version stamping, or provenance sidecar output.

## Deviations from Design
The design calls for a separate metadata file next to each event log. Current code writes JSONL lines only.

## Remaining Work
Add sidecar schema, output path rules, and writer integration if the design remains active.

## Verification Notes

Files inspected:
- `internal/logger/event_jsonl.go`
- `internal/runtime/v3runtime/run.go`
- `internal/config/v3config/config.go`

Symbols verified:
- `JSONLEventWriter`
- `JSONLSubscriber`
- `NewJSONLEventWriter`
- `WriteEvent`
- `Run`
- `LogConfig`

Configuration verified:
- `path`
- `sensors`

Tests inspected:
- `TestRunLoadsResolvedVehicleAndWritesSelectedJSONLLog`

Searches performed:
- `sidecar`
- `metadata`
- `.gdl.meta.json`
- `session metadata`
