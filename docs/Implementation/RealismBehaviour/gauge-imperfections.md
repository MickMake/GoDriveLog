# Gauge Imperfections — Implementation

## Purpose
Audits the design for a gauge-level `realism.imperfections` layer and related umbrella behaviour.

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

## Runtime Flow
No gauge-level imperfections runtime path was found. The current runtime implements individual behaviours such as `needle_shadow`, `thermal_fade`, and `calibration_offset` separately.

## Configuration
`Realism` does not declare an `Imperfections` field, and `(*Realism).UnmarshalYAML` does not accept `imperfections`.

## Behaviour
No umbrella imperfections feature matching this design was found.

## Rendering
No shared imperfections render layer was found.

## Tests
No feature-specific tests found.

## Limitations
Separate implemented behaviours were not treated as evidence for this umbrella design.

## Deviations from Design
The design calls for a gauge-level `realism.imperfections` feature. Current code implements only separate named behaviours.

## Remaining Work
Add a dedicated `realism.imperfections` model only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `UnmarshalYAML`
- `validateRealism`

Searches performed:
- `imperfections`
- `realism.imperfections`
- `gauge imperfections`
