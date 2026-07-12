# `movement` — Implementation

## Purpose
Audits current movement support across gauge families.

## Implementation Status
Partially implemented.

Verified current code implements part of the design, but the audited scope also has missing or different behaviour.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`

## Types
- `Realism`
- `Odometer`

## Functions and Methods
- `normalizePackage`
- `validateOdometer`
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `resolveOdometerMovementState`
- `effectiveMovementPolicy`
- `applyMovementPolicy`
- `applyOdometerMovementCurve`

## Runtime Flow
Odometer movement is handled by `resolveOdometerMovementState`. Radial movement is handled by `resolveMovementState`. Bar movement is handled by `resolveBarMovementState`.

## Configuration
Odometer movement uses `Odometer.Movement` and accepts `instant`, `linear`, `ease_out`, `bell`, `smooth`, and `click`. `normalizePackage` defaults odometers to `instant` and converts `smooth` and `click` to `instant` with a log message. Radial and bar movement use `Realism.MovementPolicy`, which `validateRealism` restricts to `immediate`, `linear`, and `ease_out`. `normalizePackage` defaults the policy to `immediate`.

## Behaviour
Odometers implement finite movement curves. Radial and bar gauges use movement policy only when damping, overshoot, or peg bounce make movement active; otherwise `effectiveMovementPolicy` forces immediate display updates.

## Rendering
Movement is resolved before scene generation. Scenes render the current movement state rather than performing interpolation themselves.

## Tests
- `TestLoadPackageAcceptsImplementedOdometerMovementValues`
- `TestLoadPackageWarnsAndFallsBackForRecognizedOdometerMovementValues`
- `TestLoadPackageAcceptsSharedMovementPolicies`
- `TestLoadPackageRejectsInvalidSharedMovementPolicy`
- `TestRuntimeGaugeMovementLifecycle`
- `TestRuntimeGaugeMovementEaseOutPolicyAdvancesFurtherThanLinear`
- `TestRuntimeOdometerGaugeBellMovementStartsSlowerThanLinearAndSettlesExactlyOnTarget`
- `TestRuntimeOdometerGaugeRecognizedMovementFallbacksStayInstant`

## Limitations
The scalar `movement` design is not implemented uniformly across families.

## Deviations from Design
Radial and bar gauges use `movement_policy`, not the scalar `movement` key described by the design. Radial movement does not accept `bell` as a configured policy.

## Remaining Work
Unify or explicitly separate the movement contracts if this design remains active.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/gauges/scene.go`

Symbols verified:
- `Realism`
- `Odometer`
- `normalizePackage`
- `validateOdometer`
- `validateRealism`
- `resolveMovementState`
- `resolveBarMovementState`
- `resolveOdometerMovementState`
- `effectiveMovementPolicy`
- `applyMovementPolicy`
- `applyOdometerMovementCurve`

Configuration verified:
- `movement`
- `movement_policy`
- `instant`
- `linear`
- `ease_out`
- `bell`
- `smooth`
- `click`
- `immediate`

Tests inspected:
- `TestLoadPackageAcceptsImplementedOdometerMovementValues`
- `TestLoadPackageWarnsAndFallsBackForRecognizedOdometerMovementValues`
- `TestLoadPackageAcceptsSharedMovementPolicies`
- `TestLoadPackageRejectsInvalidSharedMovementPolicy`
- `TestRuntimeGaugeMovementLifecycle`
- `TestRuntimeGaugeMovementEaseOutPolicyAdvancesFurtherThanLinear`
- `TestRuntimeOdometerGaugeBellMovementStartsSlowerThanLinearAndSettlesExactlyOnTarget`
- `TestRuntimeOdometerGaugeRecognizedMovementFallbacksStayInstant`

Searches performed:
- `movement_policy`
- `MovementBell`
- `MovementSmooth`
- `MovementClick`
