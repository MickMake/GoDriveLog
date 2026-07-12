# `decimal_point_behaviour`

Design reference: [`docs/Designs/RealismBehaviour/numeric-decimal-point-behaviour.md`](../../Designs/RealismBehaviour/numeric-decimal-point-behaviour.md)

## Purpose
Tracks the planned independent behaviour rules for decimal points in numeric and segmented displays.

## Implementation Status
Status: **Not implemented**.

Numeric formatting can render decimal points, but there is no dedicated decimal-point realism behaviour layer.

## Packages and Files
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard_test.go`](../../../internal/dashboard/v3dashboard/dashboard_test.go)

## Types
- None in current code.

## Functions and Methods
- Numeric scene formatting handles decimal point placement but not decimal-point-specific realism.

## Runtime Flow
Numeric rendering uses formatted output directly and does not model separate point fade, bleed, lag, or brightness behaviour.

## Configuration
There is no `realism.decimal_point_behaviour` key in package loading.

## Behaviour
Decimal points render as part of current glyph/layout rules only.

## Rendering
Scene composition supports decimal point placement in numeric gauges, but not independent display artefacts for the point element.

## Tests
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/dashboard_test.go`](../../../internal/dashboard/v3dashboard/dashboard_test.go)

## Limitations
Basic decimal-point rendering exists, which can make the missing realism layer easy to over-assume.

## Deviations from Design
The design is about point-specific realism, not mere decimal formatting; that realism layer is absent.

## Remaining Work
Define the config model and separate point-level rendering behaviour if the feature stays desirable.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
