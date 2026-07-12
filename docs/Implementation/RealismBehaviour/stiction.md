# `stiction`

Design reference: [`docs/Designs/RealismBehaviour/stiction.md`](../../Designs/RealismBehaviour/stiction.md)

## Purpose
Tracks thresholded release behaviour for small radial and bar changes.

## Implementation Status
Status: **Implemented**.

Radial and bar packages support `realism.stiction`, and runtime can hold tiny changes until the threshold is exceeded.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `resolveMovementState`

## Runtime Flow
The movement resolver can hold the displayed state against small source changes, then release it into a catch-up movement when the configured threshold is crossed.

## Configuration
Radial and bar packages accept `realism.stiction` as a thresholded display-only effect.

## Behaviour
Small changes can appear to stick, while larger changes still move promptly.

## Rendering
Rendering uses the currently held or released displayed state; no extra art layers are involved.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
Only radial and bar families implement stiction.

## Deviations from Design
The implementation matches the documented scope.

## Remaining Work
No design-specific work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
