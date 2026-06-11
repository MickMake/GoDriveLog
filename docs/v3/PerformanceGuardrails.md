# GoDriveLog v3 performance guardrails

Status: performance guidance  
Applies to: current display stabilisation and future v3 renderer work  
References: `MigrationGuardrails.md`, `ImplementationGuardrails.md`, `config.full.yaml`

## 1. Purpose

This document captures performance rules for the current display path and the future v3 renderer.

The current display path may need tactical speed work before the full v3 migration is complete. That work is allowed. It should not distort the v3 schema or architecture.

Performance fixes are welcome. Performance-shaped architecture accidents are not.

## 2. Core rule

```text
Optimise rendering locally.
Keep the v3 model clean.
```

The v3 target remains:

```text
vehicle endpoint
-> sensor polling runtime
-> sensor events
-> logs and dashboards as subscribers
```

Renderer performance should improve how dashboards consume state and draw pixels. It should not make dashboards poll sensors, own sensor timing, or require schema tricks.

## 3. Current display performance lane

Current display performance work is a tactical lane.

Allowed:

- profile current display/rendering hot spots
- cache decoded image assets
- cache scaled or prepared render resources where safe
- avoid re-decoding images every frame
- avoid full scene rebuilds when one value changes
- avoid full window invalidation when one widget changes
- reduce layout recalculation churn
- reduce excessive allocation in hot paths
- batch UI updates where the toolkit benefits from batching
- keep current-runtime performance fixes isolated

Not allowed:

- changing v3 config shape to suit current renderer limitations
- adding dashboard-level polling as a speed workaround
- pushing current renderer internals into v3 asset schema
- adding temporary performance flags with no removal condition
- turning old renderer concepts into required v3 concepts by accident

## 4. Performance work classification

Every performance change should be classified as one of these:

| Type | Meaning | Where it belongs |
|---|---|---|
| Tactical current-runtime fix | Makes current display usable now | current renderer/runtime path |
| Reusable renderer lesson | Applies cleanly to v3 widget renderer | document here, then implement in v3 |
| Schema requirement | Requires config shape change | stop and review docs first |
| Measurement-only change | Adds profiling/timing/logging | local instrumentation, removable |

If a speed fix starts wanting config surface area, stop and ask why.

## 5. Measurement guardrails

Do not optimise blind if the slowdown is not obvious.

Useful measurements:

- frame/render duration
- number of images decoded per frame
- number of widgets rebuilt per update
- number of allocations per update path
- number of UI invalidations per sensor event
- event rate from sensor runtime
- actual dashboard update rate
- time spent formatting values
- time spent loading/scaling assets

Prefer simple instrumentation that can be removed or left behind as low-noise debug logging.

## 6. Asset performance guardrails

Image-backed dashboards need asset discipline.

Rules:

- asset paths are repository-root relative
- load assets once where practical
- decode assets once where practical
- validate asset packs at startup
- do not discover missing files in the render hot path
- do not rescale native-size assets every update unless unavoidable
- reuse prepared images/resources where the UI toolkit allows it
- validate dimensions before rendering starts

Digit, bar, frame, indicator, and image assets should be registered and ready before the dashboard goes live.

## 7. Widget update guardrails

Widgets should update from changed state, not from whole-scene panic.

Rules:

- track the last rendered value/state per widget where useful
- skip redraw when rendered output is unchanged
- update only affected widgets where practical
- avoid reformatting values if neither value nor status changed
- indicator widgets should update on status changes even if boolean value did not change
- digit widgets should compare formatted output, not raw float noise, when deciding if visual output changed
- formatted decimal separators do not consume digit character slots
- decimal-capable digit formats require a `decimal_point` asset

For example, if speed changes from `42.1` to `42.2` and the widget format is `%03.0f`, the rendered string remains `042`; the widget does not need to redraw.

## 8. Event rate and coalescing guardrails

The sensor runtime may emit frequent events. The dashboard does not need to waste work on invisible changes.

Rules:

- sensors own polling rate
- event state updates should be cheap
- dashboard widgets decide whether an event changes rendered output
- renderer should avoid doing work for unchanged visual state
- do not solve event volume by making dashboards poll sensors
- state/change coalescing must not hide status transitions
- state/change coalescing must not hide stale transitions
- state/change coalescing must not hide error transitions
- state/change coalescing must not hide recovery transitions

If event volume becomes a problem, add state/change coalescing at the event/dashboard boundary, not OBD reads inside the renderer.

Performance fixes that hide error states are not fixes. They are lies with better frame times.

## 9. Renderer guardrails

The renderer should be boring and predictable.

Rules:

- initialise assets before live rendering
- separate state update from draw/update work
- keep UI toolkit calls on the correct thread
- avoid large allocations in per-frame/per-event code
- prefer dirty-widget updates over full-scene rebuilds
- keep native-size asset rendering where possible
- isolate toolkit-specific performance hacks inside renderer packages

Do not put Fyne-specific compromises into config schema.

## 10. Dashboard schema guardrails

Performance must not invent schema unless the requirement is real and durable.

Do not add config fields for:

- dashboard refresh cadence
- per-widget refresh cadence
- renderer throttle knobs
- cache toggles
- debug redraw switches
- backend-specific tuning
- stale timeout tuning
- asset root tuning

A future config field may be justified, but the bar is high. First prove the need with measurement and explain why it belongs in user config instead of renderer internals.

## 11. Current runtime stabilisation checklist

Before deeper v3 dashboard work, current display stabilisation should check:

- are images decoded repeatedly?
- are assets loaded repeatedly?
- is the whole scene rebuilt for every update?
- are unchanged formatted values redrawn?
- are layout calculations repeated unnecessarily?
- are UI updates happening too often?
- are sensor updates faster than visible dashboard updates need to be?
- are stale/error/status changes still visible after optimisation?

If a fix makes errors or stale states invisible, it is not a fix.

## 12. V3 renderer first useful target

The first v3 renderer target should be deliberately small:

```text
image + digit_display + indicator
```

Performance expectations for that slice:

- assets load once
- missing assets fail before live display
- digit formatted output drives redraw decisions
- indicators update on value or status changes
- stale/error/recovery transitions remain visible
- unchanged visual state does not redraw
- dashboard does not poll sensors

Do not start with frame gauges, curved bars, glass overlays, and glorious retro excess. Those are dessert. Delicious dessert, but still dessert.

## 13. Regression guardrails

Add performance-sensitive tests or checks where practical:

- digit display does not redraw if formatted output is unchanged
- digit display does redraw when status changes
- indicator redraws when status changes
- stale transitions are not hidden by coalescing
- recovery transitions are not hidden by coalescing
- asset registry does not reload the same file repeatedly
- missing asset fails during validation/load
- dashboard receives sensor state without calling reader/endpoint code

Full performance benchmarks can come later. Boundary tests should come early.

## 14. When to create issues

Create a specific issue/task when:

- a hotspot is measured and needs code work
- a current-renderer optimisation may affect v3 design
- a proposed config field exists only for speed
- a workaround needs a removal condition
- current display speed blocks dashboard iteration

Keep issues phrased around behaviour and measurement, not vibes.

Bad:

```text
Dashboard is slow.
```

Better:

```text
Current Fyne dashboard rebuilds full scene on every RPM update; avoid full rebuild when formatted output is unchanged.
```

## 15. Performance definition of done

A performance fix is done when:

- the bottleneck is identified or the fix is obviously local
- the change is isolated to the current runtime or renderer boundary
- the v3 schema remains unchanged unless explicitly reviewed
- stale/error/unknown/recovery display behaviour still works
- the improvement can be explained in one or two sentences
- any temporary workaround has a removal condition

Fast is good. Fast and understandable is better. Fast and mysterious is just a future ambush wearing running shoes.
