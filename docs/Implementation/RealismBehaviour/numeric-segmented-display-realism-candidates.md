# Numeric and Segmented Display Realism Candidates

Design reference: [`docs/Designs/RealismBehaviour/numeric-segmented-display-realism-candidates.md`](../../Designs/RealismBehaviour/numeric-segmented-display-realism-candidates.md)

## Purpose
Tracks the backlog of candidate realism behaviours for numeric and segmented displays.

## Implementation Status
Status: **Not implemented**.

The candidate realism set named by this document has not been implemented on `main`.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)

## Types
- None in current code.

## Functions and Methods
- None in current code.

## Runtime Flow
Numeric and segmented gauges render correctly as families, but the candidate realism options listed by this backlog note remain absent.

## Configuration
Package loading does not accept the candidate keys from this backlog note.

## Behaviour
The current families provide base rendering, not the richer realism catalogue described here.

## Rendering
Rendering focuses on current digits, sparse thresholds, and package-owned layout rather than display artefact realism.

## Tests
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/dashboard_test.go`](../../../internal/dashboard/v3dashboard/dashboard_test.go)

## Limitations
Basic numeric/segmented support can be mistaken for realism support; they are separate concerns.

## Deviations from Design
This file remains backlog context rather than an implemented contract.

## Remaining Work
Implement concrete candidate slices individually if they remain valuable.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
