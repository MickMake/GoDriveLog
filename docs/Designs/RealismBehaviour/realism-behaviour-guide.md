# Realism Behaviour Guide

This guide is the canonical definition of gauge realism behaviour in GoDriveLog.

It describes what each behaviour represents, why it exists, the real-world mechanism it imitates, the expected visual result, applicable gauge families, constraints, and interactions with other behaviours.

Implementation progress does not belong in this guide. Current support is recorded in [`../../Status.md`](../../Status.md).

## General visual doctrine

- Keep each behaviour subtle.
- Keep movement finite: it should visibly settle.
- Prefer simple, readable behaviour over clever behaviour.
- Do not let a realism option change source values, logs, exported values, configured ranges, or input data.
- Apply realism to displayed state only.
- Keep behaviour deterministic for the same configuration, seed, input sequence, and elapsed time.
- When a combined `99-all-options` preview looks wrong, inspect the single-option previews first.

## Capability matrix

This matrix describes conceptual applicability. It does not state whether a behaviour is implemented.

| Behaviour | Gauge families | Simulates |
|---|---|---|
| [`movement`](movement.md) | odometer, radial, bar, numeric, segmented | Deliberate mechanical or display movement rather than instant redraw. |
| [`wraparound`](wraparound.md) | odometer | Continuous rolling number drums crossing digit-strip boundaries. |
| [`odometer-drum-slop`](odometer-drum-slop.md) | odometer | Imperfectly aligned mechanical odometer drums. |
| [`odometer-carry-drag`](odometer-carry-drag.md) | odometer | Rollover coupling where a lower drum drags the next drum. |
| [`snap-settle`](snap-settle.md) | odometer, radial, bar | A moving element landing into a stable detent or final position. |
| [`hysteresis`](hysteresis.md) | radial, bar | Direction-dependent mechanical or friction offset. |
| [`stiction`](stiction.md) | radial, bar | Static friction before visible movement releases. |
| [`damping`](damping.md) | radial, bar | Mass, fluid, or electrical smoothing and lag. |
| [`overshoot`](overshoot.md) | radial, bar | Momentum carrying past the target before settling. |
| [`peg-bounce`](peg-bounce.md) | radial, bar | Contact with an end stop followed by bounded rebound. |
| [`witness-markers`](witness-markers.md) | radial, bar | Retained minimum, maximum, or follower markers. |
| [`stat-markers`](stat-markers.md) | radial, bar | Markers derived from displayed statistics such as minimum or maximum. |
| [`thermal-fade`](thermal-fade.md) | indicator | Incandescent lamp warm-up and cool-down. |
| [`needle-shadow`](needle-shadow.md) | radial | Physical needle depth casting a visible shadow. |
| [`calibration-offset`](calibration-offset.md) | radial | Slight physical misalignment between the true reading and displayed needle position. |
| [`needle-trail`](needle-trail.md) | radial | A bounded fading history of previous displayed needle positions. |
| [`backlash`](backlash.md) | odometer, radial, bar | Direction-change slack in worn gears, linkages, or drives. |
| [`per-digit-response-lag`](per-digit-response-lag.md) | numeric, segmented | Slot-level display update lag or stagger. |
| [`numeric-leading-zero-behaviour`](numeric-leading-zero-behaviour.md) | numeric, segmented | Display-specific treatment of leading zeroes. |
| [`segment-bleed-digit-bleed`](segment-bleed-digit-bleed.md) | numeric, segmented | Faint visibility from inactive segments or digit masks. |
| [`numeric-ghosting`](numeric-ghosting.md) | numeric, segmented | Persistence of the previous character or image. |
| [`uneven-brightness`](uneven-brightness.md) | numeric, segmented, indicator | Stable brightness variation between slots, segments, or lamps. |
| [`numeric-load-sag`](numeric-load-sag.md) | numeric, segmented | Display dimming caused by electrical load. |
| [`stepped-fill`](stepped-fill.md) | bar, segmented | Block-style or discrete-step fill behaviour. |
| [`quantized-fill`](quantized-fill.md) | bar, segmented | Limited display resolution and thresholded visible changes. |
| [`imperfections`](imperfections.md) | radial, bar, indicator, numeric, segmented | Deterministic ageing, wear, vibration, noise, and display artefacts. |
| [`lighting-mode`](lighting-mode.md) | radial, bar, indicator, numeric, segmented | Behaviour and appearance changes caused by illumination state. |
| [`gauge-power-lifecycle`](gauge-power-lifecycle.md) | all powered gauges | Startup, shutdown, power interruption, warm-up, and retained display state. |
| [`gauge-presets`](gauge-presets.md) | all gauge families | Reusable combinations of realism behaviours representing a gauge technology or condition. |
| [`value-zones-warning-danger-assets`](value-zones-warning-danger-assets.md) | radial, bar, indicator | Visual treatment of normal, warning, and danger ranges. |

## Behaviour document requirements

Each behaviour document should describe:

1. **Purpose** — why the behaviour exists.
2. **Physical mechanism** — what real-world effect it imitates.
3. **Applicable gauge families** — where it conceptually belongs.
4. **Expected visual behaviour** — what should be visible.
5. **Configuration model** — the intended user-facing controls.
6. **Constraints** — bounds, determinism, settling, and failure behaviour.
7. **Interactions** — how it combines with other behaviours.
8. **Non-goals** — what the behaviour must not attempt to simulate.
9. **Preview guidance** — how to judge the behaviour independently.

## Family guidance

### Odometer

Odometer behaviours usually simulate mechanical number drums, gear coupling, detents, backlash, alignment error, carry interaction, and finite rolling movement.

### Radial

Radial behaviours usually simulate analogue needle movement, friction, damping, inertia, stop pegs, calibration error, witness markers, physical depth, and dial illumination.

### Bar

Bar behaviours apply to the displayed fill or reveal extent. They may simulate damping, friction, overshoot, discrete resolution, retained markers, or illumination effects.

### Numeric and segmented

Numeric and segmented behaviours simulate display technology quirks such as slot lag, ghosting, inactive segment bleed, uneven brightness, load sag, decimal-point handling, and leading-zero handling.

Whole-image segmented gauges and per-digit numeric rendering are separate rendering models. A behaviour document must state which model it addresses when the distinction matters.

### Indicator

Indicator behaviours simulate lamp response, power state, thermal fade, uneven brightness, and illumination effects. Static lamp artwork remains the responsibility of supplied assets.

## Interaction principles

- Movement behaviours should operate on displayed state rather than source state.
- Multiple movement behaviours must have a defined order of application.
- Static visual imperfections should not alter value mapping.
- Power lifecycle behaviour may gate or modulate other visual behaviours.
- Brightness effects should combine predictably and remain bounded.
- Random-looking variation must be deterministic.
- A behaviour should degrade safely when required assets or configuration are absent.

## Documentation boundaries

Architecture documents that are not themselves realism behaviours belong elsewhere in the design tree.

Examples:

- MQTT architecture belongs under `docs/Designs/Runtime/`.
- JSONL replay belongs under `docs/Designs/Logging/`.
- Rendering and dashboard composition belong under `docs/Designs/Dashboard/`.

## Status and implementation

This guide intentionally contains no implementation-status claims.

Use:

- [`../../Status.md`](../../Status.md) for current progress;
- [`../../Implementation.md`](../../Implementation.md) for implementation records;
- matching files under `docs/Implementation/RealismBehaviour/` for feature-specific implementation details.
