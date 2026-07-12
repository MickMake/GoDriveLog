# Needle Trail

Design reference: [`docs/Designs/RealismBehaviour/needle-trail.md`](../../Designs/RealismBehaviour/needle-trail.md)

## Purpose
Tracks the planned fading history of previous radial needle positions.

## Implementation Status
Status: **Not implemented**.

There is no `realism.needle_trail` parser support, history buffer, or ghost-needle rendering on `main`.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- None in current code.

## Functions and Methods
- Current runtime resolves only the active displayed state for the live needle.

## Runtime Flow
Movement keeps the current displayed state only; it does not retain a bounded visual history of prior positions.

## Configuration
Radial packages cannot declare `realism.needle_trail` without failing validation.

## Behaviour
Once the needle moves, previous positions disappear immediately instead of lingering as a controlled trail.

## Rendering
Radial scene composition emits one live needle plus optional shadow and pointer markers, not a stack of fading historical needles.

## Tests
- None in current code.

## Limitations
The design requires both motion history retention and new render-layer composition that do not exist.

## Deviations from Design
This remains future work.

## Remaining Work
Add config support, history sampling rules, fade composition, and performance limits for bounded trails.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
