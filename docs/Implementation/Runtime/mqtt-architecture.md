# MQTT Architecture Notes

Design reference: [`docs/Designs/Runtime/mqtt-architecture.md`](../../Designs/Runtime/mqtt-architecture.md)

## Purpose
Tracks the future architectural direction for decoupling telemetry producers and consumers through MQTT.

## Implementation Status
Status: **Not implemented**.

Current runtime remains direct and in-process; it does not use MQTT anywhere in the main execution path.

## Packages and Files
- [`internal/runtime/v3runtime/run.go`](../../../internal/runtime/v3runtime/run.go)
- [`internal/config/v3config/resolve.go`](../../../internal/config/v3config/resolve.go)

## Types
- `Options`

## Functions and Methods
- `Run`
- `Resolve`

## Runtime Flow
`v3runtime.Run` resolves the selected vehicle, constructs the connector, polling runtime, and subscribers directly, and then feeds dashboard/log sinks in-process.

## Configuration
There is no MQTT broker config, topic schema, publish/subscribe client, or daemon split in runtime config.

## Behaviour
Telemetry producers and consumers remain tightly coupled inside one process runtime.

## Rendering
Dashboard rendering consumes direct runtime events, not MQTT messages.

## Tests
- [`internal/runtime/v3runtime/run_test.go`](../../../internal/runtime/v3runtime/run_test.go)

## Limitations
The design is architectural guidance only. The code has not started the MQTT slice.

## Deviations from Design
The future architecture note is intentionally ahead of current implementation.

## Remaining Work
Introduce a first slice such as daemon-to-MQTT-to-dashboard only when the architecture work is explicitly scheduled.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
