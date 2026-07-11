# JSONL Dashboard Replay

Index: 8

Status: desired

Area: dashboard runtime, logs, replay CLI

Effort: 6-10 Codex hours

Add a replay mode that consumes GoDriveLog event logs and feeds recorded events back into the dashboard runtime.

This is a core development and validation feature. It allows a real OBD session to be captured once and replayed repeatedly without the vehicle attached.

## Proposed command shape

```text
godrivelog dashboard replay --config dashboard.yaml --log drive.gdl.jsonl
```

## Replay path

```text
.gdl.jsonl -> SensorEvent stream -> dashboard runtime -> renderer
```

## Rules

- Replay recorded sensor events directly; do not pretend JSONL is an OBD adapter.
- Do not mutate the live OBD polling path.
- Preserve `event_at` and `read_at` semantics when rebuilding events.
- Replay should feed the same dashboard boundary used by live runtime rendering.
- Replay should not write new source sensor values unless explicitly configured to log replay output.
- Replay should be deterministic for the same input log, dashboard config, and replay options.
- Preview mode remains separate: preview is manual one-gauge testing; replay is recorded event-stream playback.

## Useful options

```text
--speed 1.0      # original timing
--speed 2.0      # double speed
--speed 0        # no sleeps / fastest possible
--from <time>    # optional later slice
--to <time>      # optional later slice
--loop           # optional later slice
```

## Possible future slice

```text
v3.x JSONL dashboard replay
```
