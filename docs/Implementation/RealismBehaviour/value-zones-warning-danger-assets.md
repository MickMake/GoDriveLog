# Value Zones / Warning-Danger Assets

Design reference: [`docs/Designs/RealismBehaviour/value-zones-warning-danger-assets.md`](../../Designs/RealismBehaviour/value-zones-warning-danger-assets.md)

## Purpose
Tracks the planned asset-switching feature that selects warning or danger variants when values enter configured zones.

## Implementation Status
Status: **Not implemented**.

Current code does not switch gauge-package assets by value zone, although a separate bar-widget zone mechanism already exists.

## Packages and Files
- [`internal/dashboard/v3dashboard/dashboard_test.go`](../../../internal/dashboard/v3dashboard/dashboard_test.go)
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- None in current code.

## Functions and Methods
- None in current code.

## Runtime Flow
Widget-level bar display zones can choose different cells by value, but gauge packages do not select alternate warning/danger art sets from a `zones` contract.

## Configuration
There is no package-level `zones` schema for gauge assets. Existing `WidgetConfig.Zones` applies to simple bar-display widgets, not gauge-package realism or asset selection.

## Behaviour
Gauge packages cannot swap to warning/danger dial art when values cross configured ranges.

## Rendering
Bar-display widgets can pick zone cells, but gauge-package scenes render the same configured assets regardless of value zone.

## Tests
- [`internal/dashboard/v3dashboard/dashboard_test.go`](../../../internal/dashboard/v3dashboard/dashboard_test.go)
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)

## Limitations
The existence of bar widget zones is related but not equivalent to the feature described here.

## Deviations from Design
The design describes package-asset switching. Current code only has a narrower, separate zone system for simple bar-display widgets.

## Remaining Work
Add package-level zone config, asset selection rules, validation, and tests if the feature remains wanted.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
