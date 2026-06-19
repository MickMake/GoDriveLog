# GoDriveLog v3 inverse implementation audit

Status: implementation slice output  
Target version: `v3.0.12`  
Branch: `v3.0.12-inverse-implementation-audit`

## 1. Purpose

This document inverts `docs/v3/RetirementAudit.md`.

The retirement audit asks which old or current paths may be removable later. This inverse audit asks which old or current behaviours have not yet been fully rebuilt as v3 behaviour.

The goal is to stop useful behaviour being removed before v3 can actually replace it. This is a review map only. It does not remove, move, archive, or rewrite code.

## 2. Inverse audit rule

Treat old or current code as still needed when any of these are true:

1. The v3 replacement package exists but is not wired into a runnable app path.
2. The v3 replacement covers only part of the old behaviour.
3. The v3 design deliberately deferred a decision that the old code already answers.
4. Tests preserve useful behaviour only on the old path.
5. Removing the old path would make the app less runnable, less visible, or harder to verify.

A path can be a retirement candidate and still contain behaviour worth porting first. That is not a contradiction. It is the normal awkward middle of a migration, where the old machine is both wrong and useful, like a ladder made of bad decisions.

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

Planning and audit docs:

- `docs/v3/WorkingCodeInventory.md`
- `docs/v3/MigrationState.md`
- `docs/v3/ImplementationGuardrails.md`
- `docs/v3/RetirementAudit.md`

Old/current runtime paths:

- `cmd/GoDriveLog/main.go`
- `internal/config/config.go`
- `internal/config/runtime.go`
- `internal/logger/jsonl.go`
- `internal/ui/dashboard.go`
- `internal/dashboard/renderer/fyne/`
- `internal/dashboard/assets/`
- `internal/dashboard/decoders/`
- `internal/dashboard/scene/`

v3 implementation paths:

- `internal/config/v3config/`
- `internal/vehicle/endpoint.go`
- `internal/sensors/runtime.go`
- `internal/logger/event_jsonl.go`
- `internal/assets/registry.go`
- `internal/dashboard/v3dashboard/`

## 5. Summary table

| Area | Old/current behaviour | v3 state | Gap | Priority | Recommendation |
|---|---|---|---|---|---|
| Runnable app path | `cmd/GoDriveLog/main.go` starts config, reader, logger, polling, state store, and Fyne UI | v3 packages exist separately | No runnable v3 command path wires selected vehicle -> endpoint -> polling runtime -> subscribers | Critical | Implement before retiring old command/runtime wiring |
| v3 UI/display adapter | Old `internal/ui/dashboard.go` plus old Fyne renderer displays dashboard output | `v3dashboard` produces scenes/parts | No practical Fyne/display adapter consumes v3 scenes | Critical | Build v3 display adapter before retiring old UI/renderer |
| Daily JSONL rotation | Old logger writes one file per date under a directory | v3 event writer writes exact configured path | Rotation decision not carried into v3 | Medium | Decide explicitly; do not inherit silently |
| Typed sensor values | Old and current sensor state are numeric | v3 state still stores `float64` values | Boolean/status sensors are represented by numeric convention | Medium | Keep acceptable short term, but review before broader indicators/status widgets |
| Unsupported/missing events | Old reader errors on unsupported PIDs; dashboard can show missing state | v3 has `missing/unsupported` status constant | Runtime has no distinct unsupported event kind | Medium | Add only if useful for diagnostics/display; otherwise document error mapping |
| Dashboard event efficiency | Old UI ticks and snapshots state | v3 `ApplyEvent()` rebuilds scenes before detecting unchanged output | Known performance follow-up remains | Low-medium | Optimise after v3 display path exists |
| Asset loader details | Old loader has path/frame expansion and error handling | v3 asset registry exists with richer asset families | Most important shape is implemented; some old error wording/tests may still be useful | Low | Mine tests/error cases before archiving old asset loader |
| Renderer caching | Old Fyne renderer caches image/canvas resources | v3 scene runtime is renderer-neutral | Useful Fyne cache/resource-update behaviour not ported | Medium | Port or deliberately reject before retiring old Fyne renderer |

## 6. Confirmed v3 foundations that do exist

These are not gaps by themselves:

- Strict v3 config root exists under `internal/config/v3config/` with `vehicles`, `sensors`, `assets`, `logs`, and `dashboards`.
- RuntimePlan resolution exists.
- Endpoint abstraction exists for selected vehicle OBD addresses.
- Sensor polling runtime exists and emits sensor events.
- Selected JSONL subscriber exists for v3 sensor events.
- Global asset registry exists for image, digit, indicator, bar, and frame assets.
- Selected dashboard scene runtime exists for image, digit display, indicator, bar display, and frame gauge widgets.

These pieces are good v3 foundation. The gap is mostly integration and display, not lack of bricks.

## 7. Critical gap: runnable v3 command path

Current old behaviour:

- Load the old config shape.
- Derive active sensors from old sensor/runtime rules.
- Create a shared state store.
- Open the old JSONL logger.
- Select mock or ELMOBD reader from old OBD config.
- Create the Fyne app/window.
- Start the old dashboard refresh loop.
- Start polling goroutines directly inside `main.go`.
- Write log records from the polling loop.

v3 state:

- v3 config loading and validation exists.
- RuntimePlan resolution exists.
- Endpoint connector exists.
- Polling runtime exists.
- JSONL event subscriber exists.
- v3 dashboard runtime exists.

Gap:

There is no active v3 command path that wires those pieces into the documented runtime pipeline.

Removal warning:

Do not retire `cmd/GoDriveLog/main.go`, `internal/config/runtime.go`, or the old UI/runtime path until a v3 command can run at least one selected vehicle from config through endpoint, polling runtime, selected log subscriber, and selected dashboard scene/display path.

Recommended implementation slice:

```text
v3.1.0-runnable-command-path
```

Keep it thin. The command should mostly glue existing v3 seams together rather than inventing a new runtime kingdom with banners, taxes, and a surprisingly large YAML court.

## 8. Critical gap: v3 UI/display adapter

Current old behaviour:

- `internal/ui/dashboard.go` builds a visible Fyne dashboard.
- The old dashboard path loads dashboard-local assets, decoders, scene elements, and Fyne renderer output.
- The old Fyne renderer likely contains useful resource caching and update behaviour.

v3 state:

- `internal/dashboard/v3dashboard/` can render v3 dashboards into scene/widget/part output.
- It consumes sensor state/events rather than reading OBD directly.
- It is intentionally renderer-neutral.

Gap:

There is no practical v3 adapter that turns v3 scenes/parts into a visible Fyne display.

Removal warning:

Do not retire `internal/ui/dashboard.go` or `internal/dashboard/renderer/fyne/` until one of these is true:

1. A v3 Fyne/display adapter exists and is wired into a runnable command.
2. Fyne is deliberately rejected and a different display adapter exists.
3. The project chooses headless logging/runtime first and archives display work intentionally.

Recommended implementation slice:

```text
v3.1.1-v3-display-adapter
```

## 9. Medium gap: daily JSONL rotation decision

Current old behaviour:

- `internal/logger/jsonl.go` writes daily JSONL files named by date under a directory.
- It rotates when the wall-clock date changes.
- The visible UI shows the active log path.

v3 state:

- `internal/logger/event_jsonl.go` writes event records to the exact path configured by the selected log definition.
- It preserves event/read time separately from logger write time.
- It filters duplicate unchanged state.

Gap:

The v3 path intentionally does not preserve daily rotation. That may be correct, but it is still a decision point.

Recommendation:

Do not silently reintroduce old daily rotation. Pick one of these later:

1. Keep exact configured JSONL path only.
2. Add an explicit v3 log type such as `daily_jsonl`.
3. Add a documented rotation option under the v3 log definition.

The safest near-term position is exact configured path only, because implicit rotation is how logs become a filing cabinet designed by a ferret.

## 10. Medium gap: typed sensor values

Current state:

- `sensors.SensorState` stores `Value float64`.
- Indicator widgets currently treat non-zero values as `on` and zero as `off`.
- Missing, stale, and error states are handled through explicit status.

Gap:

Boolean/status sensors are represented by numeric convention rather than a typed value model.

Recommendation:

Do not block the runnable v3 path on this. Numeric values are good enough for RPM, speed, bars, frames, and early indicator work.

Revisit before adding richer status widgets, non-OBD derived signals, or user-facing boolean sensor config. A future typed value model should stay small and should not turn config into a programming language wearing a false moustache.

## 11. Medium gap: unsupported and missing runtime events

Current state:

- v3 defines `missing/unsupported` as a sensor status.
- The dashboard runtime can render missing widget sensor state.
- The polling runtime emits first read, value change, status change, stale, error, and recovery events.

Gap:

Unsupported PIDs or unavailable sensors do not appear to have a distinct runtime event kind yet. They may currently be represented as errors or dashboard-side missing state.

Recommendation:

Keep this as a design question, not an immediate blocker. Add a distinct unsupported/missing event only if it improves diagnostics, dashboard behaviour, or log clarity.

If unsupported remains an error-path detail, document that explicitly so nobody later invents three half-compatible meanings for the same thing.

## 12. Low-medium gap: dashboard event efficiency

Current state:

`docs/v3/MigrationState.md` already carries a follow-up that `ApplyEvent()` rebuilds selected dashboard scenes via `Snapshot()` before detecting unchanged rendered output by scene signature.

Gap:

This is acceptable for the current richer-widget slice, but it will become wasteful once the v3 display adapter exists and scenes get heavier.

Recommendation:

Do not optimise before the visible v3 path exists. First make it correct and visible. Then avoid rebuilding unaffected widgets or dashboards on every sensor event.

Do not solve this by adding dashboard polling, YAML formulas, widget-owned sensor reads, or endpoint access from dashboard code.

## 13. Low gap: old asset loader details

Current old behaviour worth checking before archive:

- Rejects remote paths.
- Handles generated frame paths using `{index}` or zero-padded `{index:0N}` markers.
- Has missing-file and bad-frame-pattern test value.
- Contains useful error wording and path-resolution lessons.

v3 state:

- Global asset registry exists.
- Asset families exist for image, digit, indicator, bar, and frame assets.
- v3 paths are repository-root relative.

Gap:

Most v3 asset behaviour exists. The remaining risk is losing useful tests or error cases before they are translated.

Recommendation:

Before archiving old `internal/dashboard/assets/`, compare test coverage and error cases. Port only what still fits the v3 shape.

## 14. Retire-last list from the inverse audit

Retire these last, or only after their v3 replacements are actively wired:

1. `cmd/GoDriveLog/main.go`
2. `internal/ui/dashboard.go`
3. `internal/dashboard/renderer/fyne/`
4. `internal/logger/jsonl.go`, if daily rotation remains wanted
5. old dashboard asset/scene/decoder tests that still capture wanted behaviour

## 15. Suggested implementation order

1. Add a runnable v3 command path behind the existing command name or a temporary v3 command.
2. Wire selected vehicle -> endpoint connector -> polling runtime.
3. Wire selected JSONL subscribers to sensor events.
4. Build a v3 display adapter for `v3dashboard.Scene` output.
5. Decide whether daily JSONL rotation survives in v3.
6. Review typed sensor values and unsupported/missing event semantics.
7. Port or deliberately reject old Fyne renderer caching lessons.
8. Only then start archive/removal slices from `RetirementAudit.md`.

## 16. Explicit non-goals

This audit does not:

- delete code
- move code
- archive files
- change runtime behaviour
- change tests
- change v3 schema
- add compatibility aliases
- implement the v3 command path
- implement the v3 display adapter
- decide final deletion dates

## 17. Summary recommendation

The v3 foundation is real, but the active runnable application is still mostly old-path wiring.

Do not begin removal work yet. The next safe implementation target is a runnable v3 path that proves the existing v3 packages can work together under one selected vehicle.

In short: the parts bin is labelled, the pieces look good, and the dragon is mostly asleep. Build the cart before throwing away the old wheels.
