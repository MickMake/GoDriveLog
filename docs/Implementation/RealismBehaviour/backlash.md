# `backlash`

Design reference: [`docs/Designs/RealismBehaviour/backlash.md`](../../Designs/RealismBehaviour/backlash.md)

## Purpose
Tracks the planned odometer backlash effect for direction-change slack.

## Implementation Status
Status: **Not implemented**.

No `realism.backlash` field, parser allowance, or odometer runtime behaviour exists on `main`.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism` does not contain a backlash field.

## Functions and Methods
- `validateRealismForGaugeFamily` rejects unsupported realism keys.

## Runtime Flow
Odometer runtime movement uses existing movement curves, carry drag, snap settle, wraparound, and drum slop only.

## Configuration
`gauge.yaml` cannot declare `realism.backlash` without failing validation.

## Behaviour
Direction changes do not show any explicit slack phase before the wheel follows the new direction.

## Rendering
Rendered odometer wheel positions move directly according to the implemented movement pipeline.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
There is no reserved config shape or runtime placeholder for backlash.

## Deviations from Design
The design is still future work, and the code agrees with that.

## Remaining Work
Add config support, direction-change slack rules, and odometer tests for reversal behaviour.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
