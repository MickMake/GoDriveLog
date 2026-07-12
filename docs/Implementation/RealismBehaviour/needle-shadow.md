# `needle_shadow`

Design reference: [`docs/Designs/RealismBehaviour/needle-shadow.md`](../../Designs/RealismBehaviour/needle-shadow.md)

## Purpose
Tracks the static shadow/depth cue for radial needles.

## Implementation Status
Status: **Implemented**.

Radial packages support `realism.needle_shadow`, including default alpha handling and scene-layer placement.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- `Realism`
- `NeedleShadowConfig`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `buildNeedleShadowParts`

## Runtime Flow
The runtime does not animate the shadow separately; it resolves the displayed angle and scene composition adds the shadow part before the live needle.

## Configuration
Radial packages accept the shadow config and default a sensible alpha when enabled without custom tuning.

## Behaviour
The effect is a static display-only depth cue tied to the rendered needle angle.

## Rendering
Scene composition inserts the shadow layer before the needle and clamps/positions it with the same radial geometry.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
There is no dynamic lighting or parallax model.

## Deviations from Design
The implementation matches the design intent closely.

## Remaining Work
No known design work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
