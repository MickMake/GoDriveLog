# `carry_drag` — Implementation

## Purpose
Audits current odometer carry-drag support.

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
- `OdometerCarryDragWheelOffsets`
- `odometerCarryDragEnabled`
- `resolveOdometerMovementState`
- `applyOdometerMovementRealism`

## Runtime Flow
`resolveOdometerMovementState` updates wheel offsets, and `applyOdometerMovementRealism` applies `OdometerCarryDragWheelOffsets` when carry drag is enabled and forward movement is in flight.

## Configuration
`Realism` declares `CarryDrag *bool`. `validateRealism` restricts it to odometer gauges only.

## Behaviour
Higher wheels begin to advance before the lower wheel reaches rollover during qualifying forward movement.

## Rendering
Odometer scene rendering consumes the adjusted wheel offsets produced by movement and carry-drag logic.

## Tests
- `TestLoadPackageLoadsOdometerCarryDragRealism`
- `TestLoadPackageRejectsCarryDragOnNonOdometerGauge`
- `TestOdometerCarryDragDisabledKeepsBaseWheelOffsets`
- `TestOdometerCarryDragEnabledAdvancesHigherWheelNearRollover`
- `TestOdometerCarryDragStraddlingUpdateStartsBeforeLowerWheelPassesRollover`
- `TestRuntimeOdometerGaugeCarryDragAdvancesHigherWheelNearRolloverAndSettlesExactlyOnTarget`

## Limitations
The current implementation only applies carry drag to forward rollover movement.

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
- `OdometerCarryDragWheelOffsets`
- `odometerCarryDragEnabled`
- `resolveOdometerMovementState`
- `applyOdometerMovementRealism`

Configuration verified:
- `realism.carry_drag`

Tests inspected:
- `TestLoadPackageLoadsOdometerCarryDragRealism`
- `TestLoadPackageRejectsCarryDragOnNonOdometerGauge`
- `TestOdometerCarryDragDisabledKeepsBaseWheelOffsets`
- `TestOdometerCarryDragEnabledAdvancesHigherWheelNearRollover`
- `TestOdometerCarryDragStraddlingUpdateStartsBeforeLowerWheelPassesRollover`
- `TestRuntimeOdometerGaugeCarryDragAdvancesHigherWheelNearRolloverAndSettlesExactlyOnTarget`

Searches performed:
- `carry_drag`
- `OdometerCarryDragWheelOffsets`
