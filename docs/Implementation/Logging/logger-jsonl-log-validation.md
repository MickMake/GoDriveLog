# JSONL Log Validation — Implementation

## Purpose
Audits whether the repository can validate JSONL event logs as a first-party feature.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `cmd/GoDriveLog/main_ebiten.go`
- `internal/logger/event_jsonl.go`

## Types
None found in current code.

## Functions and Methods
- `runDashboardCLI`
- `runDashboardValidateCommand`

## Runtime Flow
No feature-specific runtime path was found. Current code can write logs and validate dashboard configs, but it does not validate log files.

## Configuration
No `logs` command tree, no `logs validate` subcommand, and no log-schema validation entrypoint were found.

## Behaviour
Current code does not provide a first-party log validation command.

## Rendering
Not applicable.

## Tests
No feature-specific tests found.

## Limitations
This record only covers current repository code.

## Deviations from Design
The design proposes `godrivelog logs validate drive.gdl.jsonl`. No matching command or validator was found.

## Remaining Work
Add a log-validation command and schema checks if this design is still wanted.

## Verification Notes

Files inspected:
- `cmd/GoDriveLog/main_ebiten.go`
- `internal/logger/event_jsonl.go`

Symbols verified:
- `runDashboardCLI`
- `runDashboardValidateCommand`
- `JSONLEventWriter`

Tests inspected:
- `TestDashboardValidateRejectsPositionalAndFlagConfig`
- `TestDashboardValidateDiscoversMultiVehicleConfigWithoutVehicle`

Searches performed:
- `logs validate`
- `dashboard validate`
- `.gdl.jsonl`
