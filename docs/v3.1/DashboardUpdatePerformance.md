# GoDriveLog v3.1.3 dashboard update performance

Status: v3.1.3 implementation

## Target

- Preferred visible dashboard cadence: `50ms`, about 20Hz.
- Minimum acceptable visible dashboard cadence: `100ms`, about 10Hz.
- Sensor polling and logging correctness take priority over dashboard freshness.

## Implementation

The visible v3 display paths use a coalescing scene sink.

```text
v3 dashboard scenes
-> latest-scene sink
-> Fyne display adapter
```

The sink keeps only the latest pending scene while rendering is still processing an earlier scene. If several updates arrive before display rendering catches up, older pending scenes are replaced.

## Scope

This slice does not change the v3 schema, add dashboard polling, or let widgets read sensors directly. It is a local display-path optimisation below the dashboard scene boundary.

## Follow-up

`v3.1.7` remains the deeper dashboard event efficiency slice for dirty-widget updates, lower allocation paths, and more granular redraw behaviour.

## Manual check

```bash
go run ./cmd/GoDriveLog \
  --v3 \
  --harness \
  --config ./docs/v3/config.example.yaml \
  --vehicle vw_caddy \
  --pattern sweep \
  --interval 50ms
```

If `50ms` is visually unreliable, retry with `--interval 100ms` and keep deeper renderer work queued for `v3.1.7`.
