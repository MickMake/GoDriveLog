# `leading_zero_behaviour`

Design reference: [`docs/Designs/RealismBehaviour/numeric-leading-zero-behaviour.md`](../../Designs/RealismBehaviour/numeric-leading-zero-behaviour.md)

## Purpose
Tracks the planned deliberate handling of leading zero slots for numeric and segmented displays.

## Implementation Status
Status: **Not implemented**.

Current formatting does not expose a dedicated leading-zero realism control.

## Packages and Files
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- None in current code.

## Functions and Methods
- None in current code.

## Runtime Flow
Numeric output is formatted according to existing package display rules, without a separate realism decision for leading zeroes.

## Configuration
There is no `realism.leading_zero_behaviour` key in package loading.

## Behaviour
Leading-zero presentation is determined by existing formatting behaviour, not by a dedicated realism layer.

## Rendering
Digit scene composition uses whatever formatted output it receives and does not add special dim/blank/ghost states for leading zero slots.

## Tests
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/dashboard_test.go`](../../../internal/dashboard/v3dashboard/dashboard_test.go)

## Limitations
The current code may already hide or show zeroes through formatting, but not in the explicit, configurable way described by the design.

## Deviations from Design
The candidate design remains unimplemented.

## Remaining Work
Define the config model and interaction with numeric formatting before implementation.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
