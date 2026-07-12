# Bar Gauge Overshoot Follow-Up — Implementation

## Purpose
Audits whether bar gauges now support the overshoot feature that this follow-up note described.

## Implementation Status
Implemented.

Verified current code provides the behaviour described in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `Realism`
- `OvershootConfig`

## Functions and Methods
- `validateRealism`
- `resolveBarMovementState`
- `radialOvershootTarget`
- `radialOvershootTravelDuration`

## Runtime Flow
Bar overshoot is handled by `resolveBarMovementState`, which computes an overshoot target and a bounded settle phase before the display returns to the final value.

## Configuration
`Realism` declares `Overshoot *OvershootConfig`. `validateRealism` accepts `realism.overshoot` for bar gauges and rejects radial-only settle fields such as `settle_cycles` and `settle_damping` for bars.

## Behaviour
Bar gauges can move past the target briefly and settle back.

## Rendering
The bar scene renders the resolved display value; there is no separate overshoot-only render layer.

## Tests
- `TestLoadPackageAcceptsBarOvershoot`
- `TestRuntimeBarGaugeOvershootDefaultDisabledStaysImmediate`
- `TestRuntimeBarGaugeOvershootAnimatesRisingReveal`
- `TestRuntimeBarGaugeOvershootAnimatesFallingReveal`
- `TestRuntimeBarGaugeOvershootStaysBoundedAndSettlesOnTarget`
- `TestRuntimeBarGaugeOvershootSettlesAtFinalReveal`

## Limitations
The current bar implementation shares the general overshoot machinery; the follow-up note does not have a dedicated package or command boundary.

## Deviations from Design
No verified deviation found in the audited scope.

## Remaining Work
No remaining work was proven by this audit for the current design note.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/package.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `Realism`
- `OvershootConfig`
- `validateRealism`
- `resolveBarMovementState`
- `radialOvershootTarget`
- `radialOvershootTravelDuration`

Configuration verified:
- `realism.overshoot`
- `ratio`
- `min_change_ratio`
- `max_span_ratio`
- `settle_mode`
- `allow_extremes`

Tests inspected:
- `TestLoadPackageAcceptsBarOvershoot`
- `TestRuntimeBarGaugeOvershootDefaultDisabledStaysImmediate`
- `TestRuntimeBarGaugeOvershootAnimatesRisingReveal`
- `TestRuntimeBarGaugeOvershootAnimatesFallingReveal`
- `TestRuntimeBarGaugeOvershootStaysBoundedAndSettlesOnTarget`
- `TestRuntimeBarGaugeOvershootSettlesAtFinalReveal`

Searches performed:
- `overshoot`
- `realism.overshoot`
- `settle_mode`
