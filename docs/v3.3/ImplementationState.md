# GoDriveLog v3.3 implementation state

Status: v3.3.4 post-Fyne renderer housekeeping in progress
Current target: v3.3.4 post-Fyne renderer housekeeping
Current branch: v334-housekeeping

## Purpose

This file records the current implementation state for v3.3. Update it in every v3.3 slice PR.

## Renderer decision

The renderer decision is made:

```text
v3.2.x = final supported Fyne dashboard line
v3.3.x = Ebiten-first active renderer line
v3.4.x = future feature/platform preparation line, not part of v3.3.4
v4.0.x = later product/platform hardening and shipping work
```

Ebiten is now the default and only active v3.3 dashboard renderer implementation. Fyne is no longer an active v3.3 dashboard runtime target, and its legacy code/dependencies have been removed from the active branch.

The renderer remains a boundary, not the owner of the app. Runtime, harness, sensors, logs, and dashboard scene generation must remain renderer-independent.

The active renderer boundary is:

```text
runtime / harness
  -> prepared vehicle/sensor data
  -> dashboard scene model
  -> display sink
  -> renderer adapter
     -> Ebiten implementation
```

## v3.3.4 housekeeping audit

v3.3.4 exists to make the v3.3 renderer transition clean, boring, and explicit before opening the v3.4 series.

Audit result expected for this slice:

- No active Go package imports `fyne.io`.
- No Fyne module dependency remains in `go.mod` or `go.sum`.
- Remaining Fyne references are historical documentation, v3.2 documentation, changelog history, or the `fyne_legacy` command notice only.
- Ebiten remains the default active renderer implementation.
- The renderer flag remains readable for command examples, but `ebiten` is the only active v3.3 renderer value.
- Gauge examples keep their visually tuned positions; comments document that seven-segment digit positions are source-artwork/package-local alignment coordinates.
- v3.2.9 is closed as superseded by the v3.3 renderer decision, not implemented as a separate branch.

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

The current example gauge positions are visually tuned. For seven-segment packages, declared `size` is the logical instrument size used for dashboard fit checks, while `digits.positions` are package-local source-artwork alignment coordinates. Those coordinates may exceed the declared logical width when the source art requires it.

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

## Raspberry Pi evidence

The renderer decision is based on target-hardware behaviour, not cosmetic preference.

Observed Raspberry Pi result before the Ebiten decision: Fyne was effectively unusable on the real baseline dashboard workload, with visible updates arriving roughly every several seconds rather than behaving like a live dashboard.

Recorded Ebiten baseline run on 2026-06-22, using `carpi` / Raspberry Pi, `sweep`, `50ms`, `60s`:

```text
events=3543
sensors=3
dashboards=1
display_submitted=1086
display_rendered=1056
display_superseded=30
display_last_render=225.479µs
```

This evidence is also recorded in `docs/v3.3/PerformanceRuns.md`.

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.3.0 | complete | Planning docs, prompts, reusable examples path, and renderer-spike checkpoint setup. |
| v3.3.1 | complete | Added the first Ebiten adapter through the real v3 dashboard scene path, fixed Linux GLFW symbol collisions by separating renderer builds, accepted `--duration`, and retained comparable display-sink stats. |
| v3.3.2 | complete | Promoted Ebiten to the primary v3.3 renderer, retired Fyne from the active v3 dashboard runtime path, documented v3.2.x as the final supported Fyne line, and kept mobile packaging as a v3.4 strategic follow-up. |
| v3.3.3 | complete | Removed legacy Fyne dashboard code packages and Fyne module dependencies from the active v3.3 branch, leaving only the `fyne_legacy` notice. |
| v3.3.4 | in progress | Post-Fyne renderer housekeeping: audit Fyne removal, preserve renderer boundary wording, document example gauge geometry, record Pi performance evidence, and close v3.2.9 as superseded. |

## Pending slices

| Version | Status | Next action |
|---|---|---|
| v3.4.0 | not started | Start in a separate v3.4 branch/chat after v3.3.4 closes; do not add v3.4 feature work to this slice. |

## Update rule

Every v3.3 implementation PR must update this file with:

- completed version;
- current branch;
- next target;
- any changed decisions or deferrals.
