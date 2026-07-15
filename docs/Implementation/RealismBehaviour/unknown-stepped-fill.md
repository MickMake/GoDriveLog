# `stepped_fill` — Implementation

## Purpose
Audits whether stepped fill exists in current code.

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
- `BarSceneWithPointerMarkers`

## Runtime Flow
No stepped-fill runtime path was found.

## Configuration
`Realism` does not declare a `SteppedFill` field, and `(*Realism).UnmarshalYAML` does not accept `stepped_fill`.

## Behaviour
Current bar rendering remains reveal-based rather than step-based under a dedicated realism option.

## Rendering
No stepped-fill render path was found.

## Tests
No feature-specific tests found.

## Limitations
This record only covers current repository code.

## Deviations from Design
The design describes stepped fill. Current code does not implement it.

## Remaining Work
Add the feature only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `UnmarshalYAML`
- `BarSceneWithPointerMarkers`

Searches performed:
- `stepped_fill`
- `realism.stepped_fill`
