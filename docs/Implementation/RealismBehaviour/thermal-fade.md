# `thermal_fade`

Design reference: [`docs/Designs/RealismBehaviour/thermal-fade.md`](../../Designs/RealismBehaviour/thermal-fade.md)

## Purpose
Tracks incandescent-style warm-up and cool-down for indicators.

## Implementation Status
Status: **Implemented**.

Indicator packages support `realism.thermal_fade`, and runtime resolves separate on/off fade timing.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`
- `ThermalFadeConfig`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `resolveIndicatorThermalFadeState`

## Runtime Flow
Indicator movement state stores fade timing and updates display alpha across off-to-on and on-to-off transitions.

## Configuration
Indicator packages accept `realism.thermal_fade` with separate rise and fall timing.

## Behaviour
Indicator lamps fade on and off softly rather than switching instantly.

## Rendering
Scene composition blends `off` and `on` layers according to the resolved fade alpha.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The feature models timing only; it does not add random flicker or ageing effects.

## Deviations from Design
The implementation matches the design intent closely.

## Remaining Work
No known design work remains.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
