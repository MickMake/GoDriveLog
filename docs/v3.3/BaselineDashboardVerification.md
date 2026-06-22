# GoDriveLog v3.3 baseline dashboard verification

Status: v3.3 renderer decision complete; Ebiten is active default

## Purpose

This baseline verifies renderer behaviour through the real GoDriveLog dashboard path.

It is intentionally not a demo renderer benchmark. The active renderer must be driven by the same upstream runtime/harness path used by real dashboards.

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

The dashboard size is `1920x720`. The current example widget placement is visually tuned and fits inside that logical working window using declared package sizes and widget scales.

| Widget | Package size | Scale | Position | Effective bounds |
|---|---:|---:|---:|---:|
| `temp_3_digit` | `316x150` | `0.5` | `[40, 40]` | `x=40..198`, `y=40..115` |
| `speed_3_digit` | `316x150` | `0.5` | `[430, 40]` | `x=430..588`, `y=40..115` |
| `rpm_4_digit` | `398x150` | `0.5` | `[820, 40]` | `x=820..1019`, `y=40..115` |
| `radial_rpm` | `512x512` | `0.9` | `[1220, 80]` | `x=1220..1681`, `y=80..541` |

Seven-segment digit positions are package-local source-artwork alignment coordinates. They may exceed the declared logical package width because the source artwork and dashboard fit box are not always the same coordinate system. Do not change those values simply to make the YAML look numerically tidy; change them only if rendered output is wrong.

## Required full path

The active renderer must use:

```text
OBD / harness source
-> prepared vehicle/sensor data
-> runtime event path
-> dashboard scene generation
-> display sink / latest submission
-> renderer adapter
   -> Ebiten implementation
-> screen
```

Do not compare against renderer-local fake values. A cardboard tachometer can be very fast; that does not make it useful.

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

`--renderer ebiten` is explicit for readability. Ebiten is already the default renderer in the active v3.3 command path.

## Fyne historical note

Fyne was used through the v3.2.x dashboard line. It is not an active v3.3 renderer.

On Raspberry Pi hardware, Fyne was observed to be effectively unusable for the real baseline dashboard workload, with visible updates arriving roughly every several seconds rather than behaving like a live dashboard. The v3.3 renderer decision moved the active path to Ebiten while preserving the same upstream scene model.

For supported Fyne dashboard behaviour, use a v3.2.x tag or branch.

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

## Recorded comparison table

| Renderer | Hardware | Pattern | Interval | Duration | Events | Submitted | Rendered | Superseded | Last render | Smoothness notes |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| Fyne | Raspberry Pi | sweep | 50ms | 60s | | | | | | Effectively unusable; visible updates roughly every several seconds. |
| Ebiten | carpi / Raspberry Pi | sweep | 50ms | 60s | 3543 | 1086 | 1056 | 30 | 225.479µs | Smooth and responsive on target hardware. |

## Decision result

Ebiten is the default and only active v3.3 renderer implementation.

This decision is functional, not cosmetic. The dashboard must be usable as a live vehicle display on Raspberry Pi-class hardware. Fyne remains historical in the v3.2.x line, while v3.3 keeps the runtime, dashboard scene model, and display-sink boundaries renderer-independent.

## Related docs

- `docs/v3.3/PerformanceRuns.md`
- `docs/v3.3/ImplementationState.md`
- `docs/v3.2/ImplementationState.md`
