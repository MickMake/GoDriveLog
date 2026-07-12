# `segment_bleed` / `digit_bleed`

Design reference: [`docs/Designs/RealismBehaviour/segment-bleed-digit-bleed.md`](../../Designs/RealismBehaviour/segment-bleed-digit-bleed.md)

## Purpose
Tracks the planned faint inactive-segment visibility for numeric and segmented displays.

## Implementation Status
Status: **Not implemented**.

There is no `segment_bleed` or `digit_bleed` config, history, or render layer on `main`.

## Packages and Files
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- None in current code.

## Functions and Methods
- None in current code.

## Runtime Flow
Runtime passes current formatted output to display scenes only; it does not track inactive-element persistence.

## Configuration
Package loading does not accept `realism.segment_bleed` or `realism.digit_bleed`.

## Behaviour
Inactive segments are not rendered as a deliberate faint mask unless artwork already bakes that in.

## Rendering
Scene composition renders active output only, with no inactive-segment overlay pass.

## Tests
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)

## Limitations
Any future implementation would need asset strategy for inactive masks across families.

## Deviations from Design
Still a candidate only.

## Remaining Work
Define asset expectations and alpha rules before implementation.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
