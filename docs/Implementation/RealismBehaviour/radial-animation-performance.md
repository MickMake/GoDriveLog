# Radial Animation Performance

Design reference: [`docs/Designs/RealismBehaviour/radial-animation-performance.md`](../../Designs/RealismBehaviour/radial-animation-performance.md)

## Purpose
Tracks the reliability of subtle radial animations on slower render targets.

## Implementation Status
Status: **Partially implemented**.

The repo has a latest-frame scenesink to reduce stale renders, but it does not implement the full radial-performance slice described here.

## Packages and Files
- [`internal/dashboard/scenesink/latest.go`](../../../internal/dashboard/scenesink/latest.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- `LatestSink`

## Functions and Methods
- `SubmitLatest` and the latest-sink submission path coalesce pending frames instead of queueing stale work.

## Runtime Flow
Dashboard rendering already has a latest-only sink path that helps short animations survive backlog pressure by discarding stale pending frames.

## Configuration
There is no radial-animation-specific config or tuning surface.

## Behaviour
Short radial animations may benefit from latest-frame coalescing, but there is no dedicated radial-tail retention or effect-specific scheduling improvement.

## Rendering
The scenesink can prefer the newest frame, which is helpful under load, but scene composition itself is unchanged.

## Tests
- [`internal/dashboard/scenesink/latest_test.go`](../../../internal/dashboard/scenesink/latest_test.go)
- [`internal/dashboard/v3dashboard/gauge_widget_test.go`](../../../internal/dashboard/v3dashboard/gauge_widget_test.go)

## Limitations
The current support is an infrastructure improvement, not a completed implementation of the design investigation.

## Deviations from Design
The design discusses targeted reliability work for subtle radial effects. Current code only covers the broader latest-frame infrastructure piece.

## Remaining Work
Measure the remaining radial failure modes and implement any effect-specific scheduling or minimum-visibility logic only if still needed.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
