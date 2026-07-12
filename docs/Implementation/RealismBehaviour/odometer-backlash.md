# Candidate: Odometer Backlash — Implementation

## Purpose
Audits the odometer-backlash planning note against current code.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`

## Functions and Methods
- `UnmarshalYAML`
- `validateRealism`
- `resolveOdometerMovementState`

## Runtime Flow
No backlash-specific runtime path was found. Current odometer realism uses carry drag, snap settle, drum slop, and movement curves only.

## Configuration
`Realism` does not declare a `Backlash` field, and `(*Realism).UnmarshalYAML` does not accept `backlash`.

## Behaviour
No direction-reversal slack behaviour matching this planning note was found.

## Rendering
No backlash-specific odometer render path was found.

## Tests
No feature-specific tests found.

## Limitations
This record is limited to current repository code.

## Deviations from Design
No verified deviation found between the planning note and current code truth: both indicate backlash is absent.

## Remaining Work
No code work is implied unless backlash is deliberately scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `UnmarshalYAML`
- `validateRealism`
- `resolveOdometerMovementState`

Searches performed:
- `backlash`
- `realism.backlash`
- `odometer backlash`
