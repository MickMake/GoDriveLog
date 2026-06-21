# GoDriveLog v3.3 baseline dashboard verification

Status: v3.3 renderer comparison planning

## Purpose

This baseline verifies renderer behaviour through the real GoDriveLog dashboard path.

It is intentionally not a demo renderer benchmark. Fyne and Ebiten must both be driven by the same upstream runtime/harness path.

Canonical config:

```text
examples/baseline-dashboard.yaml
```

## Workload

| Gauge | Type | Sensor | Notes |
|---|---|---|---|
| Temperature | 3 digit seven-segment | `coolant_temperature` | Range `-10..40`; exercises minus-symbol rendering. |
| Speed | 3 digit seven-segment | `speed` | Normal changing numeric display. |
| RPM numeric | 4 digit seven-segment | `rpm` | High-frequency digit changes. |
| RPM radial | radial gauge | `rpm` | Same RPM sensor as numeric RPM. |

## Required full path

Both renderers must use:

```text
OBD / harness source
-> prepared vehicle/sensor data
-> runtime event path
-> dashboard scene generation
-> display sink / latest submission
-> selected renderer adapter
   -> Fyne OR Ebiten
-> screen
```

Do not compare against renderer-local fake values. A cardboard tachometer can be very fast; that does not make it useful.

## Primary Fyne command

```bash
go run ./cmd/GoDriveLog \
  --harness \
  --config ./examples/baseline-dashboard.yaml \
  --vehicle vw_caddy \
  --pattern sweep \
  --interval 50ms \
  --duration 60s \
  --renderer fyne
```

## Primary Ebiten command

```bash
go run ./cmd/GoDriveLog \
  --harness \
  --config ./examples/baseline-dashboard.yaml \
  --vehicle vw_caddy \
  --pattern sweep \
  --interval 50ms \
  --duration 60s \
  --renderer ebiten
```

## Expected summary fields

Record these fields for every run:

| Field | Meaning |
|---|---|
| `events` | Total fake sensor events emitted by the harness. |
| `display_submitted` | Scene batches accepted by the latest-only display sink. |
| `display_rendered` | Scene batches actually rendered by the display adapter. |
| `display_superseded` | Scene batches replaced before rendering. |
| `display_last_render` | Duration of the last display render. |
| average render duration | Useful if available. |
| worst render duration | Useful if available. |
| CPU use | Prefer Pi measurement. |
| subjective smoothness | Note stutter, lag, or visible tearing. |

## Comparison table

| Renderer | Pattern | Interval | Duration | Events | Submitted | Rendered | Superseded | Last render | CPU | Smoothness notes |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| Fyne | sweep | 50ms | 60s | | | | | | | |
| Ebiten | sweep | 50ms | 60s | | | | | | | |

## Decision guidance

| Result | Likely decision |
|---|---|
| Ebiten wins by 2x+ and feels smoother | Continue Ebiten path. |
| Ebiten wins by 30-80% | Inspect complexity hard before continuing. |
| Ebiten wins by 10-20% | Probably not worth switching. |
| Ebiten is similar or slower | Keep Fyne primary. |
| Ebiten is unstable on Pi | Pause or abandon Ebiten. |

## Notes

- The v3.3.1 spike starts with Ebiten runtime radial needle rotation.
- If runtime rotation is too expensive, the next candidate is prepared/cached needle frames.
- Fyne remains the default renderer until the measurements justify changing it.
