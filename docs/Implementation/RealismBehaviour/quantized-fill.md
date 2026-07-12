# `quantized_fill` — Implementation

## Purpose
Audits whether quantized fill exists in current code.

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
No quantized-fill runtime path was found.

## Configuration
`Realism` does not declare a `QuantizedFill` field, and `(*Realism).UnmarshalYAML` does not accept `quantized_fill`.

## Behaviour
Bar rendering remains continuous and reveal-based.

## Rendering
No quantized-fill render path was found.

## Tests
No feature-specific tests found.

## Limitations
This record only covers current repository code.

## Deviations from Design
The design describes quantized fill. Current code does not implement it.

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
- `quantized_fill`
- `realism.quantized_fill`
