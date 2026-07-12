# `backlash` — Implementation

## Purpose
Audits whether odometer backlash exists in current code.

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
No backlash-specific runtime phase or state was found in odometer movement handling. Odometer movement currently routes through `resolveOdometerMovementState` and the existing carry-drag and snap-settle helpers only.

## Configuration
`Realism` does not declare a `Backlash` field, and `(*Realism).UnmarshalYAML` does not accept a `backlash` key.

## Behaviour
No direction-change slack behaviour matching this design was found.

## Rendering
No backlash-specific rendering path was found in odometer scene code.

## Tests
No feature-specific tests found.

## Limitations
This audit did not treat older planning documents as implementation evidence.

## Deviations from Design
The design describes odometer backlash. Current code has no parser key, runtime state, or scene logic for it.

## Remaining Work
Add `realism.backlash` parsing, validation, runtime state, and tests if this design is scheduled.

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
- `direction-change slack`
