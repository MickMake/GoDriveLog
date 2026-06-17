# GoDriveLog v3.1 dashboard harness

Status: v3.1.2 implementation

## Purpose

The dashboard harness lets a developer exercise selected v3 dashboard output without connecting to OBD hardware.

It feeds fake sensor events through the real v3 dashboard event/state path, then sends the resulting v3 dashboard scenes to the normal Fyne display adapter.

The harness is not a separate demo renderer and does not read sensors or endpoints.

## Command

```bash
go run ./cmd/GoDriveLog \
  --v3 \
  --harness \
  --config CONFIG \
  --vehicle VEHICLE_ID \
  --repo-root REPO_ROOT \
  --pattern sweep \
  --interval 100ms
```

Useful cadence values:

```text
50ms
100ms
```

Pattern names are explicit. Unknown pattern names are rejected.

## Patterns

### `sweep`

The sweep pattern uses an 11 second cycle:

1. start at sensor `min`
2. move from `min` to `max` over 5 seconds
3. pause at `max` for 1 second
4. move from `max` back to `min` over 5 seconds

This is intended for gauge, bar, and digit movement testing.

### `heartbeat`

The heartbeat pattern uses a 10 second cycle:

1. start slightly above `min` as the baseline
2. rise to a first smaller peak
3. return toward baseline
4. dip below baseline for the negative part of the cycle
5. rise to a larger second peak at `max`
6. return to baseline for the remainder of the cycle

This is intended for peak/response testing and for spotting update or redraw behaviour around quick changes.

### `fixed`

The fixed pattern holds each sensor at the midpoint between `min` and `max`.

This is mainly useful as a stable visual sanity check.

## Boundary

The harness path is:

```text
fake sensor value
-> sensors.SensorEvent
-> v3dashboard.Runtime.ApplyEvent
-> []v3dashboard.Scene
-> Fyne display adapter
```

It deliberately avoids:

```text
fake renderer
special widget playground
OBD polling changes
old UI retirement
```
