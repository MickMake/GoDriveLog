# `uneven_brightness` — Implementation

## Purpose
Audits whether uneven brightness exists in current code.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/scene.go`

## Types
None found in current code.

## Functions and Methods
- `NumericScene`
- `SegmentedScene`

## Runtime Flow
No brightness-variation state was found.

## Configuration
No `realism.uneven_brightness` key was found in current code.

## Behaviour
No stable per-slot brightness variation matching this design was found.

## Rendering
No uneven-brightness render path was found.

## Tests
No feature-specific tests found.

## Limitations
This record only covers current repository code.

## Deviations from Design
The design describes uneven brightness. Current code does not implement it.

## Remaining Work
Add the feature only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/scene.go`

Symbols verified:
- `NumericScene`
- `SegmentedScene`

Searches performed:
- `uneven_brightness`
- `realism.uneven_brightness`
