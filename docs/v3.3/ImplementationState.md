# GoDriveLog v3.3 implementation state

Status: v3.3.3 Fyne dependencies removed from active branch
Current target: v3.4 mobile/platform packaging preparation
Current branch: v3.3.3-remove-fyne-dependencies

## Purpose

This file records the current implementation state for v3.3. Update it in every v3.3 slice PR.

## Renderer decision

The renderer decision is made:

```text
v3.2.x = final supported Fyne dashboard line
v3.3.x = Ebiten migration and primary renderer line
v3.4.x = Ebiten-first platform packaging line
```

Ebiten is now the primary v3 dashboard renderer. Fyne is no longer an active v3.3 dashboard runtime target, and its legacy code/dependencies have been removed from the active branch.

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

The v3.3 command path keeps only a `fyne_legacy` build-tag notice that points users back to v3.2.x. That notice does not import or depend on Fyne.

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

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.3.0 | complete | Planning docs, prompts, reusable examples path, and renderer-spike checkpoint setup. |
| v3.3.1 | complete | Added the first Ebiten adapter through the real v3 dashboard scene path, fixed Linux GLFW symbol collisions by separating renderer builds, accepted `--duration`, and retained comparable display-sink stats. |
| v3.3.2 | complete | Promoted Ebiten to the primary v3.3 renderer, retired Fyne from the active v3 dashboard runtime path, documented v3.2.x as the final supported Fyne line, and kept mobile packaging as a v3.4 strategic follow-up. |
| v3.3.3 | implemented for review | Removed legacy Fyne dashboard code packages and Fyne module dependencies from the active v3.3 branch, leaving only the `fyne_legacy` notice. |

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
