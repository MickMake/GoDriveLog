# `realism.imperfections`

Design reference: [`docs/Designs/RealismBehaviour/imperfections.md`](../../Designs/RealismBehaviour/imperfections.md)

## Purpose
Tracks the proposed umbrella `realism.imperfections` config layer.

## Implementation Status
Status: **Not implemented**.

There is no `realism.imperfections` field or shared imperfection subsystem on `main`.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)

## Types
- `Realism` does not expose an `imperfections` field.

## Functions and Methods
- `validateRealismForGaugeFamily` rejects unsupported realism keys.

## Runtime Flow
Current realism support is implemented as separate concrete keys, not a nested imperfection layer.

## Configuration
YAML cannot declare `realism.imperfections` without failing package validation.

## Behaviour
Wear, noise, ageing, and artefact groups described by this document are not modeled as a single feature.

## Rendering
No shared imperfection pass exists in scene composition.

## Tests
- None in current code.

## Limitations
Related point features such as needle shadow or thermal fade do not add up to this umbrella contract.

## Deviations from Design
The design intentionally marks this as future work, and the code reflects that.

## Remaining Work
Define the config shape and decide which concrete artefacts belong under the umbrella before implementing anything.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
