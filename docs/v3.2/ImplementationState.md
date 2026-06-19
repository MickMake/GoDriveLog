# GoDriveLog v3.2 implementation state

Status: planning baseline in progress
Current target: v3.2.0 planning baseline
Current branch: v3.2.0-planning-baseline

## Purpose

This file records the current implementation state for v3.2. Update it in every v3.2 slice PR.

## Current direction

v3.2 adds self-contained gauge packages while keeping the existing dashboard and widget model.

A gauge package lives at:

```text
assets/gauges/**/gauge.yaml
```

The directory names under `assets/gauges/` are arbitrary. They do not imply gauge type, renderer type, sensor type, or style semantics.

The only required filename is `gauge.yaml`.

## Architecture summary

```text
dashboard
  widgets[]
    type: gauge
    gauge: gauges/.../<package-dir>
    position: [x, y]
    scale: n

gauge package
  gauge.yaml
  image files
  optional shared image files via relative paths
```

For v3.2:

- the dashboard widget places the gauge;
- the gauge package owns the sensor binding;
- the gauge package owns value mapping;
- the gauge package owns visual layers;
- the gauge package owns pivots;
- no widget-level sensor override is planned;
- no code inheritance is planned;
- no cluster layer is planned.

## Desired minimal example

Dashboard config:

```yaml
widgets:
  - id: rpm
    type: gauge
    gauge: gauges/classic/rpm
    position: [100, 20]
    scale: 0.75
```

Gauge package:

```text
assets/
  gauges/
    classic/
      images/
        bezel.png
        face_dark.png
        needle_red.png
        glass_overlay.png
      rpm/
        gauge.yaml
        ticks.png
```

`assets/gauges/classic/rpm/gauge.yaml`:

```yaml
id: rpm
type: radial
sensor: engine_rpm

size:
  width: 512
  height: 512

layers:
  background: ../images/bezel.png
  face: ../images/face_dark.png
  ticks: ticks.png
  needle: ../images/needle_red.png
  overlay: ../images/glass_overlay.png

pivot:
  face: { x: 0.5, y: 0.55 }
  needle: { x: 0.5, y: 0.9 }

value_map:
  min: 0
  max: 8000
  start_angle: -135
  end_angle: 135
  clamp: true
```

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.2.0 | in progress | Planning docs and prompts. |

## Pending slices

| Version | Status | Next action |
|---|---|---|
| v3.2.1 | not started | Implement gauge package loader. |
| v3.2.2 | not started | Add gauge widget config support. |
| v3.2.3 | not started | Add radial gauge scene model. |
| v3.2.4 | not started | Add Fyne radial renderer. |
| v3.2.5 | not started | Add example gauge package. |
| v3.2.6 | not started | Add harness verification. |
| v3.2.7 | not started | Checkpoint next direction. |

## Deferred v3.1 work

- v3.1.7 dashboard event efficiency is deferred.
- v3.1.8 retirement readiness is deferred.

They are not cancelled. They should be reconsidered at the v3.2.7 checkpoint.

## Update rule

Every v3.2 implementation PR must update this file with:

- completed version;
- current branch;
- current state;
- next target;
- any important implementation notes;
- any deferred items.
