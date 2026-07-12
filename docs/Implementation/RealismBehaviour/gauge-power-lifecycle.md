# Gauge Power Lifecycle

Design reference: [`docs/Designs/RealismBehaviour/gauge-power-lifecycle.md`](../../Designs/RealismBehaviour/gauge-power-lifecycle.md)

## Purpose
Tracks the planned gauge-level power-on and power-off realism driven by a dashboard power signal.

## Implementation Status
Status: **Not implemented**.

Current runtime has no dashboard power-state signal or per-gauge power lifecycle support.

## Packages and Files
- [`internal/dashboard/v3dashboard/dashboard.go`](../../../internal/dashboard/v3dashboard/dashboard.go)
- [`internal/runtime/v3runtime/run.go`](../../../internal/runtime/v3runtime/run.go)
- [`internal/dashboard/gauges/package.go`](../../../internal/dashboard/gauges/package.go)

## Types
- None in current code.

## Functions and Methods
- `Run` wires live sensor and dashboard streams, but not dashboard power-state events.

## Runtime Flow
The dashboard runtime reacts to live sensor snapshots and gauge realism movement only. There is no ACC-style power event path into gauges.

## Configuration
No package config keys exist for power-up sweep, blanking, drop, or shutdown behaviour.

## Behaviour
Gauges do not wake, blank, or settle differently when external dashboard power changes because no such signal is modeled.

## Rendering
Rendered state only reflects current sensor and movement state.

## Tests
- None in current code.

## Limitations
Both the runtime signal and gauge-level contracts are absent.

## Deviations from Design
The design expects a new runtime event axis that current code has not started.

## Remaining Work
Add runtime power-state inputs, per-gauge config, and lifecycle render behaviours.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
