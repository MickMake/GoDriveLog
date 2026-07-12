# `wraparound`

Design reference: [`docs/Designs/RealismBehaviour/wraparound.md`](../../Designs/RealismBehaviour/wraparound.md)

## Purpose
Tracks continuous odometer wheel routing through digit-strip boundaries.

## Implementation Status
Status: **Implemented**.

Odometer packages support `realism.wraparound`, and wheel rendering uses circular digit routing instead of disconnected digit jumps.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`

## Functions and Methods
- `wrappedDigitIndex` and the wheel-offset helpers in `scene.go` implement circular routing.

## Runtime Flow
Odometer movement resolves the target wheel offsets, and scene composition maps those offsets through a continuous 0-9 wheel strip.

## Configuration
Odometer packages accept `realism.wraparound` as a display-only realism option.

## Behaviour
Transitions such as `9 -> 0` and `0 -> 9` roll through adjacent wheel positions instead of jumping across the strip.

## Rendering
Wheel-strip slices are taken from virtual slots that wrap around the digit strip continuously.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
This effect is odometer-only.

## Deviations from Design
The implementation matches the design intent closely.

## Remaining Work
No known design work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
