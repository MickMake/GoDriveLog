# JSONL Log Validation

Design reference: [`docs/Designs/Logging/logger-jsonl-log-validation.md`](../../Designs/Logging/logger-jsonl-log-validation.md)

## Purpose
Tracks the absent validator for canonical GoDriveLog event logs.

## Implementation Status
Status: **Not implemented**.

There is no `logs validate` command or schema-validation package on `main`.

## Packages and Files
- [`cmd/GoDriveLog/main_ebiten.go`](../../../cmd/GoDriveLog/main_ebiten.go)
- [`internal/logger/event_jsonl.go`](../../../internal/logger/event_jsonl.go)

## Types
- None in current code.

## Functions and Methods
- `main` does not register a `logs` command tree.
- `JSONLEventWriter` emits events but does not expose a standalone validator.

## Runtime Flow
Logs are produced during runtime, but there is no offline validation step before replay or conversion.

## Configuration
No validator flags, schema options, or CLI entrypoints exist.

## Behaviour
Malformed or incomplete log files are not checked by first-party tooling.

## Rendering
No rendering impact.

## Tests
- [`cmd/GoDriveLog/main_ebiten_test.go`](../../../cmd/GoDriveLog/main_ebiten_test.go)

## Limitations
Validation depends on a formal schema and canonical log contract that have not been finished.

## Deviations from Design
The design expects an explicit CLI validation pass. Current code has none.

## Remaining Work
Add a `logs validate` command, parse-and-check logic, and tests against good and bad canonical logs.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
