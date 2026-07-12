# Radial Movement Options — Implementation

## Purpose
Audits current radial movement configuration against the radial movement options design.

## Implementation Status
Partially implemented.

Verified current code implements part of the design, but the audited scope also has missing or different behaviour.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`

## Functions and Methods
- `normalizePackage`
- `validateRealism`
- `resolveMovementState`
- `effectiveMovementPolicy`
- `applyMovementPolicy`

## Runtime Flow
Radial movement is handled by `resolveMovementState` and shaped by `applyMovementPolicy`.

## Configuration
Current code uses `Realism.MovementPolicy`, not a scalar `movement` key. `validateRealism` accepts `immediate`, `linear`, and `ease_out` only. `normalizePackage` defaults the field to `immediate`.

## Behaviour
Radial gauges can move immediately, linearly, or with an ease-out curve when movement is active. `effectiveMovementPolicy` forces immediate behaviour when no radial movement feature is active, and promotes active `immediate` movement to `linear` when damping, overshoot, or peg bounce need visible travel.

## Rendering
The render path uses the current resolved display value; radial movement is not implemented inside scene rendering.

## Tests
- `TestLoadPackageAcceptsSharedMovementPolicies`
- `TestLoadPackageRejectsInvalidSharedMovementPolicy`
- `TestLoadPackageRejectsMisspelledSharedMovementPolicyKey`
- `TestRuntimeRadialGaugeMovementDefaultsToImmediateWithoutDamping`
- `TestRuntimeGaugeMovementEaseOutPolicyAdvancesFurtherThanLinear`

## Limitations
The design proposes a scalar `movement` contract and includes `bell`. Current radial implementation does not expose that contract.

## Deviations from Design
Current code uses `movement_policy` instead of `movement`, and it does not accept `bell` as a radial configuration value.

## Remaining Work
Adopt the design contract or update the design only if this feature is still active.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `normalizePackage`
- `validateRealism`
- `resolveMovementState`
- `effectiveMovementPolicy`
- `applyMovementPolicy`

Configuration verified:
- `realism.movement_policy`
- `immediate`
- `linear`
- `ease_out`

Tests inspected:
- `TestLoadPackageAcceptsSharedMovementPolicies`
- `TestLoadPackageRejectsInvalidSharedMovementPolicy`
- `TestLoadPackageRejectsMisspelledSharedMovementPolicyKey`
- `TestRuntimeRadialGaugeMovementDefaultsToImmediateWithoutDamping`
- `TestRuntimeGaugeMovementEaseOutPolicyAdvancesFurtherThanLinear`

Searches performed:
- `movement_policy`
- `MovementPolicyImmediate`
- `MovementPolicyEaseOut`
