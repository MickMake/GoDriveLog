# Bar Realism Scope

Design reference: [`docs/Designs/RealismBehaviour/bar-realism-scope.md`](../../Designs/RealismBehaviour/bar-realism-scope.md)

## Purpose
Tracks which planned realism behaviours for bar gauges have landed and which remain backlog items.

## Implementation Status
Status: **Partially implemented**.

Bar gauges support several shared movement behaviours, but the full scope document is not complete.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `Realism`
- `DampingConfig`
- `OvershootConfig`
- `PointerMarkersConfig`

## Functions and Methods
- `validateRealismForGaugeFamily`
- `resolveMovementState`
- `barDampingDuration`

## Runtime Flow
Bar runtime supports damping, hysteresis, stiction, overshoot, peg bounce, and pointer markers through the shared movement path.

## Configuration
Implemented bar keys are accepted through `realism`. `stepped_fill` and `quantized_fill` remain unsupported and are rejected as unknown realism keys.

## Behaviour
Supported features animate the displayed reveal height only. Planned block/quantized fill behaviours are still absent.

## Rendering
Current rendering is still a continuous reveal model driven by package geometry and final resolved extent.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The document includes behaviours that have no parser, runtime, or rendering support yet.

## Deviations from Design
The scope doc now spans a mixed state: some bar realism options are live, while fill-quantisation concepts remain backlog.

## Remaining Work
Implement `stepped_fill` and `quantized_fill` if they remain desirable, or narrow the scope document to current support.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
