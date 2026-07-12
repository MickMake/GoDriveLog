# `ghosting`

Design reference: [`docs/Designs/RealismBehaviour/numeric-ghosting.md`](../../Designs/RealismBehaviour/numeric-ghosting.md)

## Purpose
Tracks the planned residual afterimage effect for numeric and segmented displays.

## Implementation Status
Status: **Not implemented**.

Current numeric and segmented rendering has no previous-character persistence or ghost image path.

## Packages and Files
- [`internal/dashboard/gauges/scene.go`](../../../internal/dashboard/gauges/scene.go)
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)

## Types
- None in current code.

## Functions and Methods
- None in current code.

## Runtime Flow
The dashboard runtime stores the current displayed state only. It does not preserve previous displayed glyphs for decay rendering.

## Configuration
There is no `realism.ghosting` key in package loading.

## Behaviour
Digit changes replace the previous display immediately within the current formatting and scene-signature pipeline.

## Rendering
Numeric and segmented scenes render active output only; no faded previous-character overlay is composed.

## Tests
- [`internal/dashboard/gauges/scene_test.go`](../../../internal/dashboard/gauges/scene_test.go)
- [`internal/dashboard/v3dashboard/dashboard_test.go`](../../../internal/dashboard/v3dashboard/dashboard_test.go)

## Limitations
Adding ghosting would need glyph-history retention and family-specific render treatment.

## Deviations from Design
The design remains a candidate and has not started in code.

## Remaining Work
Define history timing, alpha rules, and asset expectations before implementation.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
