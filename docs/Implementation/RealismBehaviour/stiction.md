# `stiction` — Implementation

## Purpose
Audits current stiction support for radial and bar gauges.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`
- `ValueMap`

## Functions and Methods
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `stictionShouldHold`

## Runtime Flow
Stiction is evaluated by the movement resolvers through `stictionShouldHold` before a new movement phase is started.

## Configuration
`Realism` declares `Stiction *float64`. `validateRealism` restricts the field to radial and bar gauges, requires a finite positive threshold, and bounds the value by the `ValueMap` span.

## Behaviour
Small changes can be held at the previous display value until the threshold is exceeded.

## Rendering
The render path uses the current resolved display value; stiction is handled before scene generation.

## Tests
- `TestLoadPackageAcceptsRadialStiction`
- `TestLoadPackageAcceptsBarStiction`
- `TestLoadPackageRejectsInvalidStiction`
- `TestRuntimeRadialGaugeStictionBelowThresholdHoldsDisplay`
- `TestRuntimeRadialGaugeStictionReleasesAboveThreshold`
- `TestRuntimeBarGaugeStictionBelowThresholdHoldsDisplay`
- `TestRuntimeBarGaugeStictionReleasesAboveThreshold`

## Limitations
Only radial and bar gauges implement stiction.

## Deviations from Design
No verified deviation found in the audited scope.

## Remaining Work
No remaining work was proven by this audit.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `ValueMap`
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `stictionShouldHold`

Configuration verified:
- `realism.stiction`

Tests inspected:
- `TestLoadPackageAcceptsRadialStiction`
- `TestLoadPackageAcceptsBarStiction`
- `TestLoadPackageRejectsInvalidStiction`
- `TestRuntimeRadialGaugeStictionBelowThresholdHoldsDisplay`
- `TestRuntimeRadialGaugeStictionReleasesAboveThreshold`
- `TestRuntimeBarGaugeStictionBelowThresholdHoldsDisplay`
- `TestRuntimeBarGaugeStictionReleasesAboveThreshold`

Searches performed:
- `stiction`
- `stictionShouldHold`
