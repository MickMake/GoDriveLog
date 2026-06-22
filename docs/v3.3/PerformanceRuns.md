# GoDriveLog v3.3 performance runs

Status: active v3.3 renderer evidence

## Purpose

This file records renderer and baseline-dashboard performance evidence for the v3.3 renderer decision.

The goal is not to chase synthetic benchmark trophies. The goal is to answer a practical question: can the selected renderer make the real GoDriveLog dashboard usable on the target class of hardware?

## Renderer decision evidence

Fyne was not merely visually jerky on the Raspberry Pi baseline dashboard workload. It was effectively unusable as a live vehicle dashboard, with visible updates arriving roughly every several seconds.

Ebiten is the default and only active v3.3 renderer implementation because it makes the same scene path usable on the target hardware while preserving the renderer boundary.

The required path remains:

```text
OBD / harness source
-> prepared vehicle/sensor data
-> runtime event path
-> dashboard scene generation
-> display sink / latest submission
-> renderer adapter
-> screen
```

## Recorded runs

| Date | Source | Renderer | Hardware | Pattern | Interval | Duration | Events | Submitted | Rendered | Superseded | Last render | Notes |
|---|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| 2026-06-22 | PR74/local | ebiten | carpi / Raspberry Pi | sweep | 50ms | 60s | 3543 | 1086 | 1056 | 30 | 225.479µs | Smooth and responsive on target hardware. |

## Derived notes for 2026-06-22 run

- Harness emitted about 59 events/second over the 60 second run.
- Display submitted about 18 scene batches/second.
- Display rendered about 17.6 scene batches/second.
- Superseded scene batches were 30 of 1086, about 2.8%.
- Last render was about 0.225ms.

The important result is not just the low last-render value. It is that the real baseline scene path remained usable and responsive on the Raspberry Pi.

## Recording rule

For future performance runs, record:

- branch or PR;
- renderer;
- hardware;
- command or pattern details;
- duration and interval;
- `events`;
- `display_submitted`;
- `display_rendered`;
- `display_superseded`;
- `display_last_render`;
- subjective smoothness notes.
