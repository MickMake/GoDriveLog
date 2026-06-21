# GoDriveLog v3.3 implementation state

Status: v3.3.0 planning setup in progress
Current target: v3.3.0 renderer checkpoint planning
Current branch: v3.3.0-docs-renderer-plan-and-examples

## Purpose

This file records the current implementation state for v3.3. Update it in every v3.3 slice PR.

## Current direction

v3.3 tests whether an Ebiten renderer should be added beside Fyne for the live dashboard path.

Fyne remains the default renderer unless Ebiten clearly wins on Raspberry Pi hardware.

The current renderer boundary is useful:

```text
runtime / harness
  -> prepared vehicle/sensor data
  -> dashboard scene model
  -> display sink
  -> renderer adapter
```

v3.3 should swap only the renderer adapter when using `--renderer fyne|ebiten`.

## Full-path requirement

Ebiten must not be tested through a renderer-private demo path.

The required comparison path is:

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

A fake renderer-local RPM counter is not an acceptable benchmark. Cardboard dashboards have excellent frame rates and terrible evidentiary value.

## Baseline dashboard

The canonical reusable baseline config is:

```text
examples/baseline-dashboard.yaml
```

The baseline workload is:

| Gauge | Type | Sensor | Notes |
|---|---|---|---|
| Temperature | 3 digit seven-segment | `coolant_temperature` | Range `-10..40`; exercises minus-symbol rendering. |
| Speed | 3 digit seven-segment | `speed` | Normal changing numeric display. |
| RPM numeric | 4 digit seven-segment | `rpm` | High-frequency digit changes. |
| RPM radial | radial gauge | `rpm` | Same RPM sensor as numeric RPM; exercises radial needle rendering. |

## Renderer spike rule

v3.3.1 should start with Ebiten runtime radial needle rotation.

If runtime rotation is too costly on the Pi, switch to prepared/cached radial needle frames, matching the successful Fyne prepared-frame strategy.

## Example command shape

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

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.3.0 | in progress | Planning docs, prompts, reusable examples path, and renderer-spike checkpoint setup. |

## Pending slices

| Version | Status | Next action |
|---|---|---|
| v3.3.1 | not started | Add experimental Ebiten renderer backend beside Fyne using the real dashboard path. |
| v3.3.2 | not started | Run fixed-duration A/B comparison between Fyne and Ebiten. |
| v3.3.3 | not started | Decide whether to continue, promote, pause, or abandon Ebiten. |
| v3.3.4 | conditional | Act only if the v3.3.3 decision identifies a clear follow-up. |

## Update rule

Every v3.3 implementation PR must update this file with:

- completed version;
- current branch;
- next target;
- any changed decisions or deferrals.
