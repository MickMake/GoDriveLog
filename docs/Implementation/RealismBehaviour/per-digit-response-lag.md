# `per_digit_response_lag`

Design reference: [`docs/Designs/RealismBehaviour/per-digit-response-lag.md`](../../Designs/RealismBehaviour/per-digit-response-lag.md)

## Purpose
Tracks the planned slot-by-slot update lag for numeric and segmented displays.

## Implementation Status
Status: **Not implemented**.

Current digit updates are not staggered per slot by a dedicated realism layer.

## Packages and Files
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- None in current code.

## Functions and Methods
- None in current code.

## Runtime Flow
The dashboard runtime keeps one current displayed value per widget rather than a timed lag state per digit slot.

## Configuration
There is no `realism.per_digit_response_lag` key in package loading.

## Behaviour
Multi-digit displays update as a single formatted output rather than slot-by-slot delayed transitions.

## Rendering
Numeric and segmented scenes render the current display state without per-slot timing offsets.

## Tests
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/dashboard_test.go`](../../../internal/dashboard/v3dashboard/dashboard_test.go)

## Limitations
Implementing this would need slot history and timing per digit, not just current formatted output.

## Deviations from Design
Still a candidate only.

## Remaining Work
Define the lag model and state retention rules before implementation.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
