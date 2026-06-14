# GoDriveLog v3 retirement audit

Status: implementation slice output  
Target version: `v3.0.11`  
Branch: `v3.0.11-retirement-audit`

## 1. Purpose

This document records a docs-only retirement audit for the end of the v3.0.x foundation line.

It identifies old or current paths that may be removed, archived, or left alone later. It does not remove, move, or archive anything. This is a review map for Mick before any future cleanup PRs.

## 2. Retirement rule

Retire a path only after all of these are true:

1. The v3 replacement path exists.
2. The v3 replacement covers the wanted behaviour.
3. Tests cover the behaviour that should survive.
4. The old path is no longer the active runtime path, or keeping both active paths causes real confusion.
5. Mick has reviewed the candidate and approved the actual removal or archive step.

## 3. v3 target reminder

```text
selected vehicle
-> OBD endpoint
-> sensor polling runtime
-> sensor events
-> selected logs and dashboards as subscribers
```

The selected vehicle owns the runtime profile. Sensors and assets remain global catalogues. Logs and dashboards are global definitions selected by the vehicle.

## 4. Areas inspected

Planning docs:

- `docs/v3/MigrationState.md`
- `docs/v3/MigrationGuardrails.md`
- `docs/v3/ImplementationGuardrails.md`
- `docs/v3/WorkingCodeInventory.md`

Old/current paths:

- `cmd/GoDriveLog/main.go`
- `internal/config/config.go`
- `internal/config/runtime.go`
- `internal/config/dashboard.go`
- `internal/logger/jsonl.go`
- `internal/ui/dashboard.go`
- `internal/dashboard/assets/`
- `internal/dashboard/decoders/`
- `internal/dashboard/scene/`
- `internal/dashboard/renderer/fyne/`

v3 foundation paths:

- `internal/config/v3config/`
- `internal/vehicle/endpoint.go`
- `internal/sensors/runtime.go`
- `internal/logger/event_jsonl.go`
- `internal/assets/registry.go`
- `internal/dashboard/v3dashboard/`

## 5. Summary table

| Area | Old/current path | v3 path | Recommendation | Confidence | Removal condition |
|---|---|---|---|---|---|
| Config root | `internal/config/config.go` | `internal/config/v3config/` | Remove or archive later | High | Active CLI uses v3 config and old root keys are no longer supported |
| Runtime orchestration | `cmd/GoDriveLog/main.go` | v3 resolve + endpoint + sensors + subscribers | Replace later | High | A v3 entrypoint can run the selected vehicle pipeline |
| Runtime sensor selection | `internal/config/runtime.go` | `v3config.Resolve` + `sensors.NewPollingRuntime` | Remove later | High | No active code uses `sensor.log` to decide polling |
| Old dashboard schema | `internal/config/dashboard.go` | v3 dashboard widgets/assets | Archive or remove later | High | Old block/layer config is no longer loaded by the active app |
| Old daily JSONL writer | `internal/logger/jsonl.go` | `internal/logger/event_jsonl.go` | Keep temporarily | Medium | v3 event logging is active and the date-rotation decision is settled |
| Old UI dashboard shell | `internal/ui/dashboard.go` | `internal/dashboard/v3dashboard` plus future UI adapter | Keep temporarily | Medium | A v3 UI/display adapter exists |
| Old dashboard asset loader | `internal/dashboard/assets/` | `internal/assets/registry.go` | Archive/remove later | Medium | v3 asset registry covers active dashboard assets |
| Old decoder pipeline | `internal/dashboard/decoders/` | Widget behaviour in Go code | Archive/remove later | High | No active config uses decoder IDs |
| Old scene evaluator | `internal/dashboard/scene/` | `v3dashboard.Render` / `ApplyEvent` | Archive/remove later | Medium | v3 renderer consumes v3 dashboard output directly |
| Old Fyne renderer | `internal/dashboard/renderer/fyne/` | Future v3 renderer adapter | Keep for now | Low-medium | Useful Fyne caching/resource behaviour is ported or deliberately rejected |
| Shared sensor reader/state | `internal/sensors/reader.go`, `state.go`, `state_store.go`, `elmobd_reader.go` | Used by v3 endpoint/sensor runtime | Keep | High | Do not remove while v3 depends on these pieces |

## 6. Candidate notes

### `internal/config/config.go`

Current role: loads the old root shape: `obd`, `log`, `vehicle`, `sensors`, `dashboard`. It also applies old defaults such as dashboard refresh and top-level log/OBD defaults.

v3 replacement: `internal/config/v3config/` defines the v3 root allow-list: `vehicles`, `sensors`, `assets`, `logs`, `dashboards`.

Recommendation: remove or archive later, after the active command path loads v3 config directly.

Risk: removing this too early may break the current executable before the v3 entrypoint exists.

### `internal/config/runtime.go`

Current role: builds runtime sensors from old config, uses `sensor.log` to decide polling, and derives stale timing from old `refresh`.

v3 replacement: `v3config.Resolve` plus `sensors.NewPollingRuntime`.

Recommendation: remove later, after no active code uses `config.ActiveSensors` or `config.SensorStateDefinitions`.

Risk: removing this too early breaks the current old runtime path.

### `internal/config/dashboard.go`

Current role: defines the old dashboard schema with `refresh_ms`, `render_min_ms`, `canvas`, `asset_root`, `assets`, `decoders`, `blocks`, and `layers`.

v3 replacement: v3 dashboard widgets and global v3 asset families.

Recommendation: archive or remove later, after v3 UI/display work is useful and old examples have been archived or rewritten.

Risk: useful visual design knowledge may disappear before it is translated into v3 widgets or docs.

### `cmd/GoDriveLog/main.go`

Current role: wires the old app path: old config, old logger, old reader branching, old per-sensor goroutines, and old ticker-driven Fyne dashboard.

v3 replacement: a future command path should use v3 config resolution, endpoint connector, sensor polling runtime, JSONL event subscribers, and v3 dashboard subscribers.

Recommendation: replace later, not in this audit.

Risk: the project could have good v3 packages but no runnable app path.

### `internal/logger/jsonl.go`

Current role: writes old `sensors.Reading` objects to daily JSONL files under a directory.

v3 replacement: `internal/logger/event_jsonl.go` writes selected v3 sensor events.

Recommendation: keep temporarily. Decide whether daily date-based rotation should survive in v3 before removing it.

Risk: daily rotation behaviour could be lost accidentally.

### `internal/ui/dashboard.go`

Current role: builds the old dashboard UI from old config, dashboard-local assets, decoders, scene evaluation, and the old Fyne renderer.

v3 replacement: `internal/dashboard/v3dashboard/` provides selected dashboard scene/widget state, but a practical v3 UI adapter still appears to be future work.

Recommendation: keep temporarily.

Risk: removing this too early may leave v3 dashboard logic without a display path.

### `internal/dashboard/assets/`

Current role: loads old dashboard-local assets, using config-relative paths and optional `asset_root`.

v3 replacement: `internal/assets/registry.go` loads global repo-root-relative asset families.

Recommendation: archive/remove later after active dashboard examples use v3 assets only.

Risk: useful asset-loading tests and error handling may be lost before v3 coverage is complete.

### `internal/dashboard/decoders/`

Current role: implements the old decoder pipeline.

v3 replacement: widget behaviour lives in Go code; `digit_display`, `indicator`, `bar_display`, and `frame_gauge` now cover the near-term display behaviours.

Recommendation: archive/remove later with high confidence.

Risk: some visual behaviours may need translation notes before removal.

### `internal/dashboard/scene/`

Current role: evaluates old blocks and layers into scene elements.

v3 replacement: `v3dashboard.Render` and `ApplyEvent` map v3 widgets to scenes/parts.

Recommendation: archive/remove later.

Risk: a future v3 renderer may still want ideas from this intermediate scene model.

### `internal/dashboard/renderer/fyne/`

Current role: renders old scene elements with Fyne and likely contains useful caching/resource-update behaviour.

v3 replacement: not complete yet. A v3 Fyne/display adapter should be implemented before retiring this.

Recommendation: keep for now.

Risk: removing this too early may discard useful UI performance knowledge.

## 7. Paths that should stay

Keep these as v3 foundation or shared seam code:

- `internal/config/v3config/`
- `internal/vehicle/endpoint.go`
- `internal/sensors/runtime.go`
- `internal/sensors/reader.go`
- `internal/sensors/elmobd_reader.go`
- `internal/sensors/state.go`
- `internal/sensors/state_store.go`
- `internal/logger/event_jsonl.go`
- `internal/assets/registry.go`
- `internal/dashboard/v3dashboard/`

These are either v3 implementation paths or reused behaviour behind v3 seams.

## 8. Suggested cleanup order after review

1. Finish or confirm a runnable v3 entrypoint.
2. Move active config examples and docs to v3-only usage.
3. Retire old config root and old runtime sensor selection.
4. Retire old logger only after deciding daily rotation.
5. Retire old dashboard decoder/block/layer schema after v3 UI rendering is useful.
6. Retire old Fyne renderer last, after useful rendering lessons are ported or rejected.

## 9. Open questions for Mick

1. Should daily JSONL rotation survive in v3, or is a configured exact JSONL path enough?
2. Should old dashboard visual examples be archived under `docs/archive`, or deleted after v3 equivalents exist?
3. Should the old Fyne renderer caching strategy become the basis for the v3 renderer adapter?
4. Should the first v3.1.0 runtime keep the same `cmd/GoDriveLog` command name, or introduce a separate command until v3 is proven?
5. Are any old decoder behaviours still wanted as first-class v3 widgets?

## 10. Explicit non-goals

This audit does not:

- delete code
- move code
- archive files
- change runtime behaviour
- change tests
- change v3 schema
- add compatibility aliases
- decide final deletion dates

## 11. Summary recommendation

The v3.0.x line has enough foundation to identify likely retirement candidates, but not enough to safely remove the old app/display path blindly.

Use this document as the review map. The next safe step is manual review, then small cleanup PRs only after v3.1.0 or another runnable v3 path proves each replacement.
