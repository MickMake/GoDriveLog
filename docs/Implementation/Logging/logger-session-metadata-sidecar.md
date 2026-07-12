# Session Metadata Sidecar

Design reference: [`docs/Designs/Logging/logger-session-metadata-sidecar.md`](../../Designs/Logging/logger-session-metadata-sidecar.md)

## Purpose
Tracks the proposed sidecar file that would capture replay and provenance metadata next to each event log.

## Implementation Status
Status: **Not implemented**.

Current logging writes event lines only; it does not emit a session metadata sidecar.

## Packages and Files
- [`internal/logger/event_jsonl.go`](../../../internal/logger/event_jsonl.go)
- [`internal/runtime/v3runtime/run.go`](../../../internal/runtime/v3runtime/run.go)

## Types
- None in current code.

## Functions and Methods
- `NewJSONLEventWriter`
- `Run`

## Runtime Flow
Runtime logging creates writers for event lines only. No second output file is opened or populated with session metadata.

## Configuration
The runtime plan can choose log destinations, but no metadata schema or sidecar path is configured.

## Behaviour
Captured logs cannot self-describe the config, asset set, vehicle selection, or runtime context needed for later replay.

## Rendering
No rendering impact.

## Tests
- [`internal/logger/event_jsonl_test.go`](../../../internal/logger/event_jsonl_test.go)
- [`internal/runtime/v3runtime/run_test.go`](../../../internal/runtime/v3runtime/run_test.go)

## Limitations
Without the sidecar, replay and audit tooling would need external context or ad hoc conventions.

## Deviations from Design
The design wants replay-safe provenance next to every log. Current code writes only the event stream.

## Remaining Work
Define the metadata schema, emit the file alongside logs, and thread config provenance into the writer setup.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
