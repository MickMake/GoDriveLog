# Gauge Presets

Design reference: [`docs/Designs/RealismBehaviour/gauge-presets.md`](../../Designs/RealismBehaviour/gauge-presets.md)

## Purpose
Tracks the planned reusable preset/profile layer for gauge visuals and realism.

## Implementation Status
Status: **Not implemented**.

Gauge packages are loaded as self-contained YAML files today; there is no preset inheritance or named profile system.

## Packages and Files
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)

## Types
- None in current code.

## Functions and Methods
- Package loading validates a single package definition with local realism config only.

## Runtime Flow
Packages are resolved and rendered directly from their own YAML definitions.

## Configuration
There is no `preset`, profile include, or pre-merge config stage in current package loading.

## Behaviour
Every gauge repeats its own realism and asset settings instead of reusing named preset blocks.

## Rendering
Rendering only sees the final explicit package config loaded from each package.

## Tests
- [`internal/dashboard/gauges/package_test.go`](../../../internal/dashboard/gauges/package_test.go)

## Limitations
Introducing presets would affect loading order, validation, and override precedence.

## Deviations from Design
The design describes a preset/profile system that does not exist in the current loader.

## Remaining Work
Define preset resolution, merging rules, and validation boundaries before implementation.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
