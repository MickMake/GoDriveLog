# Indicator Realism Scope

Design reference: [`docs/Designs/RealismBehaviour/indicator-realism-scope.md`](../../Designs/RealismBehaviour/indicator-realism-scope.md)

## Purpose
Tracks the limited current realism support for indicator gauges.

## Implementation Status
Status: **Partially implemented**.

Indicator realism currently consists of `thermal_fade`; the broader scope remains backlog.

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
Indicator runtime can interpolate on/off transitions using separate rise and fall timing, but no other realism path exists.

## Configuration
Indicator packages accept `realism.thermal_fade`. No additional indicator realism keys are implemented.

## Behaviour
Indicators can fade like incandescent lamps, but they do not support extra stateful lighting, ageing, or display artefacts.

## Rendering
Indicator scene composition blends `off` and `on` layers according to resolved fade alpha.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The document is explicitly planning-oriented and extends beyond current code.

## Deviations from Design
The implementation matches the narrow supported subset described in the planning notes.

## Remaining Work
Add more indicator-specific realism only if later designs define concrete behaviour beyond `thermal_fade`.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
