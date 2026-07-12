# External Converter Boundary

Design reference: [`docs/Designs/Logging/tools-converters-external-converter-boundary.md`](../../Designs/Logging/tools-converters-external-converter-boundary.md)

## Purpose
Tracks the planned boundary that keeps foreign-format conversion out of GoDriveLog core runtime.

## Implementation Status
Status: **Not implemented**.

There is no `tools/converters` implementation or converter boundary in the current tree.

## Packages and Files
- [`internal/logger/event_jsonl.go`](../../../internal/logger/event_jsonl.go)

## Types
- None in current code.

## Functions and Methods
- `JSONLEventWriter` can emit the current event stream, but there are no first-party converter entrypoints.

## Runtime Flow
The runtime only produces native live events and optional JSONL logs. It does not import or convert external formats.

## Configuration
No converter registration, import CLI, or external-format mapping layer exists.

## Behaviour
Any foreign-format conversion remains outside the repo or must be written ad hoc.

## Rendering
No rendering impact.

## Tests
- None in current code.

## Limitations
The converter boundary depends on a formal canonical event log and companion tooling that are not present yet.

## Deviations from Design
The design names an explicit tools area. Current code has not created it.

## Remaining Work
Add the `tools/converters` layout, define import/export boundaries, and target the canonical event log format.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
