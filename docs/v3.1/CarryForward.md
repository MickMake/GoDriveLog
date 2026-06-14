# GoDriveLog v3.1 carry-forward list

Status: planning
Owner: migration implementor

## Purpose

This file records unfinished work carried forward from the v3.0 docs.

The original v3.0 docs remain the source history. This file keeps the active v3.1 reminders short enough to use during implementation review.

## Source docs

- `docs/v3/WorkingCodeInventory.md`
- `docs/v3/MigrationState.md`
- `docs/v3/RetirementAudit.md`
- `docs/v3/InverseImplementationAudit.md`

## Confirmed v3.0 foundation

Do not rebuild these from scratch in v3.1:

- Strict v3 config loading and validation.
- RuntimePlan resolution.
- Vehicle endpoint abstraction.
- Sensor polling runtime.
- Sensor event spine and latest-state store.
- Selected JSONL event logging.
- Asset registry for image, digit, indicator, bar, and frame assets.
- Selected dashboard scene runtime for image, digit display, indicator, bar display, and frame gauge widgets.

v3.1 should mostly wire these pieces together into a practical app path.

## Exact old/current paths to review

These old/current paths still matter because they contain runnable behaviour, display behaviour, logging behaviour, or useful tests:

- `cmd/GoDriveLog/main.go`
- `internal/config/config.go`
- `internal/config/runtime.go`
- `internal/logger/jsonl.go`
- `internal/ui/dashboard.go`
- `internal/dashboard/renderer/fyne/`
- `internal/dashboard/assets/`
- `internal/dashboard/decoders/`
- `internal/dashboard/scene/`

Do not remove or archive these paths only because v3 foundation packages exist. They can be retired only after the relevant v3.1 replacement behaviour is verified.

## Exact v3 paths to use

These are the v3 foundation paths v3.1 should use first:

- `internal/config/v3config/`
- `internal/vehicle/endpoint.go`
- `internal/sensors/runtime.go`
- `internal/logger/event_jsonl.go`
- `internal/assets/registry.go`
- `internal/dashboard/v3dashboard/`

If a v3.1 slice starts by building a parallel replacement for one of these, stop and justify it first. Prefer wiring existing seams.

## Carried work items

### Runnable app path

The active app path still needs to be wired through v3 config, RuntimePlan, endpoint connection, sensor polling runtime, selected logging, and dashboard output.

Minimum useful path:

1. Load v3 config.
2. Select one vehicle.
3. Resolve RuntimePlan.
4. Connect endpoint.
5. Start sensor polling runtime.
6. Subscribe selected JSONL logs.
7. Connect selected dashboard output boundary.
8. Shut down cleanly.

Do not retire `cmd/GoDriveLog/main.go`, `internal/config/runtime.go`, or the old UI/runtime path until this works for at least one selected vehicle.

### Display adapter

The v3 dashboard scene runtime exists, but v3.1 still needs a practical display adapter before old UI paths can be retired.

The first adapter should prove that v3 dashboard scene output can be shown without letting dashboard code read sensors or endpoints directly.

Retirement warning:

- Do not retire `internal/ui/dashboard.go` until a v3 display path exists.
- Do not retire `internal/dashboard/renderer/fyne/` until its useful rendering and caching lessons are ported or deliberately rejected.

### JSONL rotation

The old logger supported daily rotation. The v3 logger currently writes to the configured path.

v3.1 must decide whether rotation survives and how it is represented.

Acceptable choices:

1. Keep exact configured path only.
2. Add an explicit log type such as `daily_jsonl`.
3. Add a documented rotation option under the v3 log definition.

Do not silently inherit daily rotation.

### Sensor value typing

Current v3 sensor state uses numeric values.

v3.1 must decide whether boolean or status values need stronger typing before richer status widgets, non-OBD derived signals, or user-facing boolean sensor config are added.

Do not block the runnable command path on this unless the implementation exposes a concrete need.

### Unsupported and missing sensors

v3.1 must decide whether unavailable sensors need explicit runtime events or whether current status/error handling is enough.

If unavailable sensors remain error-path details, document the mapping so logging and dashboard behaviour do not drift apart.

### Dashboard event efficiency

Current dashboard event handling may rebuild more scene state than necessary.

This is not the first problem to solve. Optimise only after a real display path exists and the cost is visible.

Do not solve this with dashboard polling, YAML formulas, widget-owned sensor reads, or endpoint access from dashboard code.

### Asset loader lessons

Before archiving old asset loader paths, compare old test coverage and error cases.

Carry forward or deliberately reject:

- Remote path rejection.
- Generated frame paths using `{index}`.
- Zero-padded generated frame paths using `{index:0N}`.
- Missing-file diagnostics.
- Bad frame-pattern diagnostics.
- Repository-root path resolution expectations.

### Renderer lessons

Before retiring the old Fyne renderer, decide what to do with old resource caching and image/canvas update behaviour.

v3 scene output is renderer-neutral. The display adapter should stay below the dashboard runtime boundary.

## Retire-last list

Retire these last, or only after their replacement behaviour is actively wired and verified:

1. `cmd/GoDriveLog/main.go`
2. `internal/ui/dashboard.go`
3. `internal/dashboard/renderer/fyne/`
4. `internal/logger/jsonl.go`, if daily rotation remains wanted
5. Old dashboard asset, decoder, scene, and renderer tests that still capture wanted behaviour

## Suggested v3.1 order

1. Add a runnable v3 command path behind the existing command name or a temporary v3 command.
2. Wire selected vehicle to endpoint connector and sensor polling runtime.
3. Wire selected JSONL subscribers to sensor events.
4. Build a display adapter for v3 dashboard scene output.
5. Decide whether daily JSONL rotation survives in v3.
6. Review typed sensor values.
7. Review unsupported and missing sensor semantics.
8. Port or deliberately reject old asset and renderer lessons.
9. Only then start archive or removal slices.
