# `load_sag` — Implementation

## Purpose
Audits whether numeric or segmented load sag exists in current code.

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
No display-load calculation or brightness sag state was found.

## Configuration
No `realism.load_sag` key was found in current code.

## Behaviour
No value-dependent dimming behaviour matching this design was found.

## Rendering
No brightness sag render path was found.

## Tests
No feature-specific tests found.

## Limitations
This record only covers current repository code.

## Deviations from Design
The design describes load-dependent display sag. Current code does not implement it.

## Remaining Work
Add the feature only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/scene.go`

Symbols verified:
- `NumericScene`
- `SegmentedScene`

Searches performed:
- `load_sag`
- `realism.load_sag`
