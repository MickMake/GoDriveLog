# GoDriveLog v3.1.3 dashboard update performance

Status: v3.1.3 implementation

## Target

- Preferred visible dashboard cadence: `50ms`, about 20Hz.
- Minimum acceptable visible dashboard cadence: `100ms`, about 10Hz.
- Reference target remains Raspberry Pi 4 class hardware, including the Raspberry Pi 4 2GB memory limit.
- Sensor polling and JSONL logging correctness take priority over dashboard freshness.

## Problem being solved

This slice is a narrow dashboard performance fix for the visible Fyne display path.

The Raspberry Pi 4 2GB target was running out of memory because the Fyne adapter rebuilt the full canvas/image object tree on every dashboard update. That behaviour created avoidable allocation churn, RSS growth, and slower visible updates.

The goal of this slice is not to redesign the dashboard runtime. The goal is to stop the visible display path from doing wasteful work on every frame.

## Design decisions

### Reuse Fyne image objects when possible

When the rendered dashboard part count is stable, the Fyne adapter reuses existing `canvas.Image` objects and updates their resource, position, and size instead of recreating the entire canvas object tree.

```text
same rendered part count
-> update existing canvas.Image objects
-> avoid full object-tree rebuild
```

If the rendered part count changes, rebuilding the visible object list is still allowed. Stable dashboards should avoid that hot-path churn.

### Coalesce visible display frames

The visible v3 display paths use a coalescing scene sink.

```text
v3 dashboard scenes
-> latest-scene sink
-> Fyne display adapter
```

The sink keeps only the latest pending scene while rendering is still processing an earlier scene. If several updates arrive before display rendering catches up, older pending display frames are replaced. Runtime sensor events and JSONL log events are not dropped by this mechanism; the coalescing boundary is only for visible dashboard scenes.

The latest dashboard scene wins. Stale visual frames do not build up behind a slow renderer.

### Preserve render error propagation

The sink must not hide display adapter errors. `Submit()` returns once its submitted frame has either:

- rendered successfully;
- been superseded by a newer pending frame; or
- hit a render error.

That preserves prompt/runtime error visibility while still allowing stale visual frames to be replaced by newer scenes.

### Keep Fyne work on the Fyne thread

All Fyne UI updates stay inside `fyne.DoAndWait`. The dashboard/runtime code submits scenes to the display boundary; it does not manipulate Fyne objects directly.

### Size the v3 window from selected dashboard config

The v3 window is sized from the selected v3 dashboard configuration before the window is created. The earlier hard-coded `800x480` startup size is no longer the source of truth for v3 dashboards.

The old later resize safety net was removed because window manipulation during event-drain/shutdown caused Fyne thread warnings during Ctrl-C shutdown.

### Keep shutdown boring

Ctrl-C should stop the harness/runtime cleanly and print the summary without Fyne thread warnings. The shutdown rule is:

```text
stop runtime/harness
close scene sink
avoid window manipulation during drain/shutdown
```

## Scope

This slice does not change the v3 schema, add dashboard polling, or let widgets read sensors directly. It is a local visible-display-path optimisation below the dashboard scene boundary.

This slice also does not complete deeper dirty-widget rendering, per-widget invalidation, lower-allocation scene generation, or full dashboard event-efficiency work.

## Known limitation

This slice reduces display memory churn and stale-frame buildup in the visible Fyne path. It does not fully eliminate every possible dashboard backpressure path under sustained slow rendering.

The current design still favours prompt render-error propagation and simple shutdown semantics over a deeper asynchronous dashboard event pipeline. If sustained slow rendering still affects runtime behaviour, that belongs in the later dashboard event efficiency slice.

## Follow-up

`v3.1.7` remains the deeper dashboard event efficiency slice for dirty-widget updates, lower allocation paths, more granular redraw behaviour, and any remaining sustained-render-backpressure work.

## Manual check

```bash
go run ./cmd/GoDriveLog \
  --v3 \
  --harness \
  --config ./docs/v3/config.example.yaml \
  --vehicle vw_caddy \
  --pattern sweep \
  --interval 50ms
```

Expected checks:

- RSS should not grow rapidly from repeated Fyne canvas/image object-tree rebuilds.
- `50ms` is the preferred visual cadence.
- If `50ms` is visually unreliable, retry with `--interval 100ms`.
- Ctrl-C should stop cleanly, print the summary, and avoid Fyne thread warnings.
- If `100ms` is still unreliable on Raspberry Pi 4 class hardware, document the bottleneck and keep deeper renderer work queued for `v3.1.7`.
