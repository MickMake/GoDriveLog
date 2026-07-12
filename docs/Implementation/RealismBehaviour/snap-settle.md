# `snap_settle` — Implementation

## Purpose
Audits current odometer snap-settle support.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`
- `Odometer`

## Functions and Methods
- `validateRealism`
- `OdometerSnapSettleWheelOffsets`
- `odometerSnapSettleEnabled`
- `resolveOdometerMovementState`
- `applyOdometerMovementRealism`

## Runtime Flow
`resolveOdometerMovementState` computes movement timing and `applyOdometerMovementRealism` applies `OdometerSnapSettleWheelOffsets` during the settle phase when enabled.

## Configuration
`Realism` declares `SnapSettle *bool`. `validateRealism` restricts it to odometer gauges only.

## Behaviour
After the main odometer travel, wheels can move through a short bounded settle phase and then land exactly on target.

## Rendering
Odometer scene rendering consumes the adjusted wheel offsets from movement and snap-settle logic.

## Tests
- `TestLoadPackageLoadsOdometerSnapSettleRealism`
- `TestLoadPackageRejectsSnapSettleOnNonOdometerGauge`
- `TestOdometerSnapSettleDisabledKeepsBaseWheelOffsets`
- `TestOdometerSnapSettleEnabledAddsSmallForwardSettleAndReturnsToTarget`
- `TestOdometerSnapSettleDoesNotOvershootBelowZeroAtLowerBoundary`
- `TestRuntimeOdometerGaugeSnapSettleAddsShortTailAndSettlesExactlyOnTarget`

## Limitations
The feature is odometer-only.

## Deviations from Design
No verified deviation found in the audited scope.

## Remaining Work
No remaining work was proven by this audit.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `Odometer`
- `validateRealism`
- `OdometerSnapSettleWheelOffsets`
- `odometerSnapSettleEnabled`
- `resolveOdometerMovementState`
- `applyOdometerMovementRealism`

Configuration verified:
- `realism.snap_settle`

Tests inspected:
- `TestLoadPackageLoadsOdometerSnapSettleRealism`
- `TestLoadPackageRejectsSnapSettleOnNonOdometerGauge`
- `TestOdometerSnapSettleDisabledKeepsBaseWheelOffsets`
- `TestOdometerSnapSettleEnabledAddsSmallForwardSettleAndReturnsToTarget`
- `TestOdometerSnapSettleDoesNotOvershootBelowZeroAtLowerBoundary`
- `TestRuntimeOdometerGaugeSnapSettleAddsShortTailAndSettlesExactlyOnTarget`

Searches performed:
- `snap_settle`
- `OdometerSnapSettleWheelOffsets`
