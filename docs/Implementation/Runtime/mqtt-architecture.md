# MQTT Architecture Notes — Implementation

## Purpose
Audits whether the repository has implemented the MQTT architecture note.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `internal/runtime/v3runtime/run.go`
- `internal/config/v3config/resolve.go`

## Types
- `Options`
- `RuntimePlan`

## Functions and Methods
- `Run`
- `Resolve`

## Runtime Flow
`Run` loads a config file, resolves a `RuntimePlan`, connects a vehicle reader, starts the polling runtime, and delivers events directly to log subscribers and the dashboard sink. No MQTT publish or subscribe path was found.

## Configuration
No MQTT broker, topic, or transport config was found in current runtime configuration types.

## Behaviour
Current telemetry flow is direct and in-process.

## Rendering
Dashboard rendering receives direct runtime scenes, not MQTT-delivered messages.

## Tests
No feature-specific tests found.

## Limitations
This record covers current repository code only.

## Deviations from Design
The design note describes a future MQTT-based architecture. Current code does not implement it.

## Remaining Work
Add MQTT transport only if that architecture work is explicitly scheduled.

## Verification Notes

Files inspected:
- `internal/runtime/v3runtime/run.go`
- `internal/config/v3config/resolve.go`

Symbols verified:
- `Options`
- `RuntimePlan`
- `Run`
- `Resolve`

Searches performed:
- `mqtt`
- `broker`
- `topic`
