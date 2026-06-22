# GoDriveLog v3.3 implementation state

Status: v3.3.2 Ebiten promoted to primary renderer
Current target: v3.4 mobile/platform packaging preparation
Current branch: v3.3.2-ebiten-primary-renderer

## Purpose

This file records the current implementation state for v3.3. Update it in every v3.3 slice PR.

## Renderer decision

The renderer decision is now made:

```text
v3.2.x = final supported Fyne dashboard line
v3.3.x = Ebiten migration and primary renderer line
v3.4.x = Ebiten-first platform packaging line
```

Ebiten is now the primary v3 dashboard renderer. Fyne is no longer an active v3.3 dashboard runtime target.

The decision is based on real dashboard workload feedback: Fyne was visibly jerky/noisy, while Ebiten was smooth on the same path. Ebiten also gives the project a strategic path toward Android and iOS packaging.

The active renderer boundary remains:

```text
runtime / harness
  -> prepared vehicle/sensor data
  -> dashboard scene model
  -> display sink
  -> Ebiten renderer adapter
```

## Fyne support policy

Fyne support ends with the v3.2.x line.

Use a v3.2.x tag or branch for supported Fyne dashboard builds. v3.3.x and later should not add new Fyne dashboard features or maintain feature parity with Fyne.

The v3.3 command path keeps only a `fyne_legacy` build-tag notice so accidental Fyne runs fail loudly and point users back to v3.2.x.

## Full-path requirement

Ebiten must not be tested through a renderer-private demo path.

The required path is:

```text
OBD / harness source
-> prepared vehicle/sensor data
-> runtime event path
-> dashboard scene generation
-> display sink / latest submission
-> Ebiten renderer adapter
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

## Primary command shape

Ebiten is now the normal v3 dashboard command path:

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

The `--renderer ebiten` flag remains explicit for readability, but Ebiten is the default renderer in the active v3.3 command build.

Fyne is not an active v3.3 renderer. The only v3.3 Fyne command shape is a legacy notice:

```bash
go run -tags fyne_legacy ./cmd/GoDriveLog
```

For supported Fyne dashboard behaviour, use the v3.2.x line instead.

## Renderer implementation rule

Ebiten runtime radial needle rotation remains acceptable until Raspberry Pi measurements prove it too expensive.

If runtime rotation becomes too costly on the Pi, switch to prepared/cached radial needle frames, matching the successful Fyne prepared-frame strategy from v3.2.x.

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.3.0 | complete | Planning docs, prompts, reusable examples path, and renderer-spike checkpoint setup. |
| v3.3.1 | complete | Added the first Ebiten adapter through the real v3 dashboard scene path, fixed Linux GLFW symbol collisions by separating renderer builds, accepted `--duration`, and retained comparable display-sink stats. |
| v3.3.2 | implemented for review | Promoted Ebiten to the primary v3.3 renderer, retired Fyne from the active v3 dashboard runtime path, documented v3.2.x as the final supported Fyne line, and kept mobile packaging as a v3.4 strategic follow-up. |

## Pending slices

| Version | Status | Next action |
|---|---|---|
| v3.4.0 | not started | Prepare Ebiten-first desktop, Raspberry Pi, Android, and iOS platform packaging without changing the dashboard scene model. |

## Update rule

Every v3.3 implementation PR must update this file with:

- completed version;
- current branch;
- next target;
- any changed decisions or deferrals.
