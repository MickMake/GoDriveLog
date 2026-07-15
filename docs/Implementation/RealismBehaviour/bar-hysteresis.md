# `hysteresis` — Implementation

## Purpose
Audits current hysteresis support for radial and bar gauges.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`

## Types
- `Realism`
- `ValueMap`

## Functions and Methods
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `radialHysteresisDisplayTarget`
- `barHysteresisDisplayTarget`

## Runtime Flow
Radial hysteresis is applied in `resolveMovementState`. Bar hysteresis is applied in `resolveBarMovementState`. Both use approach direction and value-map span to shift the displayed target.

## Configuration
`Realism` declares `Hysteresis *bool`. `validateRealism` accepts the field for radial and bar gauges only.

## Behaviour
The displayed target can shift slightly depending on whether the value approached from above or below. Stored raw source values remain unchanged.

## Rendering
The renderer uses the resolved display value. Hysteresis itself is handled before scene construction.

## Tests
- `TestLoadPackageAcceptsRadialHysteresis`
- `TestLoadPackageAcceptsBarHysteresis`
- `TestLoadPackageRejectsHysteresisOnUnsupportedGaugeType`
- `TestRuntimeRadialGaugeHysteresisOffsetsRisingApproach`
- `TestRuntimeRadialGaugeHysteresisOffsetsFallingApproach`
- `TestRuntimeBarGaugeHysteresisOffsetsRisingApproach`
- `TestRuntimeBarGaugeHysteresisOffsetsFallingApproach`

## Limitations
Only radial and bar gauges implement hysteresis.

## Deviations from Design
No verified deviation found in the audited scope.

## Remaining Work
No remaining work was proven by this audit.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`

Symbols verified:
- `Realism`
- `ValueMap`
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `radialHysteresisDisplayTarget`
- `barHysteresisDisplayTarget`

Configuration verified:
- `realism.hysteresis`

Tests inspected:
- `TestLoadPackageAcceptsRadialHysteresis`
- `TestLoadPackageAcceptsBarHysteresis`
- `TestLoadPackageRejectsHysteresisOnUnsupportedGaugeType`
- `TestRuntimeRadialGaugeHysteresisOffsetsRisingApproach`
- `TestRuntimeRadialGaugeHysteresisOffsetsFallingApproach`
- `TestRuntimeBarGaugeHysteresisOffsetsRisingApproach`
- `TestRuntimeBarGaugeHysteresisOffsetsFallingApproach`

Searches performed:
- `hysteresis`
- `realism.hysteresis`
- `radialHysteresisDisplayTarget`
- `barHysteresisDisplayTarget`
