# Gauge Lighting Mode — Implementation

## Purpose
Audits whether gauges react to a lighting-mode signal in current code.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/runtime/v3runtime/run.go`

## Types
None found in current code.

## Functions and Methods
- `Run`
- `NewRuntime`

## Runtime Flow
No lights-state event, lighting-mode state, or asset-selection path was found in current runtime code.

## Configuration
No gauge package lighting-mode config was found in `Package` or `Realism`.

## Behaviour
Current gauges do not implement a lights-on or lights-off visual mode from this design.

## Rendering
No lighting-mode render path was found.

## Tests
No feature-specific tests found.

## Limitations
This record only covers current repository code.

## Deviations from Design
The design requires a runtime lights signal and gauge-owned visual response. Current code has neither.

## Remaining Work
Add signal plumbing, package config, and tests only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/runtime/v3runtime/run.go`

Symbols verified:
- `Run`
- `NewRuntime`
- `Package`
- `Realism`

Searches performed:
- `lighting_mode`
- `lights`
- `lighting mode`
