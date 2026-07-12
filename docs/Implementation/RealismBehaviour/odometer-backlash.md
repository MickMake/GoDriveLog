# Candidate: Odometer Backlash

Design reference: [`docs/Designs/RealismBehaviour/odometer-backlash.md`](../../Designs/RealismBehaviour/odometer-backlash.md)

## Purpose
Tracks the backlog note confirming that odometer backlash is not implemented on `main`.

## Implementation Status
Status: **Not implemented**.

The audit note is accurate: there is still no odometer backlash support in parser or runtime.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- None in current code.

## Functions and Methods
- `validateRealismForGaugeFamily`
- `validateOdometerMovementMode`

## Runtime Flow
Odometer movement uses implemented curves and realism helpers only; there is no direction-change slack phase.

## Configuration
No `realism.backlash` key exists for odometer packages.

## Behaviour
Reversals do not show the designed slack-before-follow behaviour.

## Rendering
Wheel rendering follows the implemented movement pipeline directly.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
This file is already a correction note rather than a feature request, so the main job is to keep documentation honest.

## Deviations from Design
None; the code matches the note that backlash is not implemented.

## Remaining Work
Only implement if a later slice defines concrete backlash behaviour.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
