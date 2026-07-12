# Gauge Lighting Mode

Design reference: [`docs/Designs/RealismBehaviour/lighting-mode.md`](../../Designs/RealismBehaviour/lighting-mode.md)

## Purpose
Tracks the planned per-gauge reaction to dashboard lights-state changes.

## Implementation Status
Status: **Not implemented**.

There is no dashboard lights-state signal or gauge asset switching tied to vehicle lighting on `main`.

## Packages and Files
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)
- [`internal/runtime/v3runtime/run.go`](../../../internal/runtime/v3runtime/run.go)
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)

## Types
- None in current code.

## Functions and Methods
- `Run` and dashboard runtime do not expose a lighting-state event path.

## Runtime Flow
The runtime has no dedicated concept of headlights, illumination, or day/night mode for gauges.

## Configuration
Package loading does not support alternate light-on asset sets or a `lighting_mode` config block.

## Behaviour
Gauge appearance does not change when dashboard lighting would change in a real vehicle.

## Rendering
All art selection remains static except for currently supported indicator state layers and bar zones.

## Tests
- None in current code.

## Limitations
The feature depends on both runtime signal delivery and gauge-level asset contracts that do not exist.

## Deviations from Design
The design expects a new dashboard event axis and asset-switching logic that current code lacks.

## Remaining Work
Add lighting-state inputs, package config, and scene selection rules for alternate assets.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
