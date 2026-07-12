# Odometer Movement Cleanup Candidates

Design reference: [`docs/Designs/RealismBehaviour/odometer-movement-cleanup-candidates.md`](../../Designs/RealismBehaviour/odometer-movement-cleanup-candidates.md)

## Purpose
Tracks the cleanup note for reserved odometer movement values such as `smooth` and `click`.

## Implementation Status
Status: **Partially implemented**.

Current loading recognises the reserved values and falls back to `instant`, but they do not have distinct runtime meanings.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Types
- None in current code.

## Functions and Methods
- `validateOdometerMovementMode`

## Runtime Flow
Odometer runtime implements `instant`, `linear`, `ease_out`, and `bell`. Reserved `smooth` and `click` values warn and behave like immediate movement.

## Configuration
The parser accepts the reserved values only as compatibility fallbacks; they are not separate features.

## Behaviour
Users can load older or speculative configs without a hard failure, but they do not get special smooth/click semantics.

## Rendering
Rendering stays on the existing odometer movement curves; there is no extra click-phase or alternate curve for the reserved names.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The reserved values still exist in a limbo state that can confuse readers.

## Deviations from Design
The note argues against inventing meaning without a dedicated slice. Current code follows that advice by warning and falling back.

## Remaining Work
Either remove the compatibility names later or define them explicitly in a focused odometer movement slice.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
