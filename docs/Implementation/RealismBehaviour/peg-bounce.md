# `peg_bounce` — Implementation

## Purpose
Audits current peg-bounce support for radial and bar gauges.

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
- `radialPegBounceValues`
- `radialPegBounceDurations`

## Runtime Flow
Peg bounce is scheduled by the movement resolvers when the displayed target is a clamped stop and the movement conditions qualify.

## Configuration
`Realism` declares `PegBounce *bool`. `validateRealism` restricts the option to radial and bar gauges and requires a clamped `ValueMap` range when the option is enabled.

## Behaviour
The displayed value can rebound from the minimum or maximum stop and settle back. In-range targets do not trigger peg bounce.

## Rendering
The render path uses the current movement state; peg bounce itself is resolved before scene generation.

## Tests
- `TestLoadPackageAcceptsRadialPegBounce`
- `TestLoadPackageAcceptsBarPegBounce`
- `TestLoadPackageRejectsInvalidPegBounce`
- `TestRuntimeRadialGaugePegBounceDoesNotTriggerForInRangeTarget`
- `TestRuntimeRadialGaugePegBounceAtMaxStopSettlesBackToLimit`
- `TestRuntimeBarGaugePegBounceAtMaxStopSettlesBackToLimit`
- `TestRuntimeBarGaugePegBounceDoesNotTriggerForInRangeTarget`

## Limitations
Only radial and bar gauges implement peg bounce.

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
- `radialPegBounceValues`
- `radialPegBounceDurations`

Configuration verified:
- `realism.peg_bounce`

Tests inspected:
- `TestLoadPackageAcceptsRadialPegBounce`
- `TestLoadPackageAcceptsBarPegBounce`
- `TestLoadPackageRejectsInvalidPegBounce`
- `TestRuntimeRadialGaugePegBounceDoesNotTriggerForInRangeTarget`
- `TestRuntimeRadialGaugePegBounceAtMaxStopSettlesBackToLimit`
- `TestRuntimeBarGaugePegBounceAtMaxStopSettlesBackToLimit`
- `TestRuntimeBarGaugePegBounceDoesNotTriggerForInRangeTarget`

Searches performed:
- `peg_bounce`
- `radialPegBounceValues`
