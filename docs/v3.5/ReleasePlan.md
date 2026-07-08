# v3.5 Gauge Realism Release Plan

Status: implemented except odometer backlash; stale backlash implementation claim corrected

v3.5 adds gauge realism in small, inspectable slices. It builds on the v3.4 gauge package work and keeps the existing asset-driven dashboard model.

## Purpose

v3.5 improves how gauges behave when values change. It focuses on:

- static imperfection that does not require a frame tick;
- finite movement responses after value changes;
- small deterministic display effects that need renderer support;
- simple preview support so each behaviour can be judged by eye.

Use `docs/v3.5/RealismBehaviourGuide.md` for the intended visual feel of each implemented realism option.

It does not add a general physics engine, idle animation, ambient flicker, dashboard power lifecycle effects, or asset-generation work.

## Existing v3.4 behaviour to preserve

Do not break the v3.4 gauge type model:

- `numeric` remains an image-character formatted value display.
- `radial` remains a value-to-angle transform gauge.
- `odometer` remains a rolling digit/wheel-strip transform gauge.
- `indicator` remains an on/off image-selection gauge.
- `bar` remains a continuous transform/reveal gauge.
- `segmented` remains a stepped threshold image-selection gauge.

## Realism configuration doctrine

Most v3.5 realism options live under the `realism` key.

`movement` is the exception. It is the single scalar movement knob and should be accepted by any gauge type for now. Gauge types that do not yet have concrete movement behaviour may parse `movement` and use their current immediate behaviour until a later slice defines more.

Keep config collapsed where possible. Simple options may be booleans or scalars. Expand an option only when it needs settings.

Odometer movement shape:

```yaml
odometer:
  movement: bell

realism:
  wraparound: true
  drum_slop:
    offsets: [0.0, 0.03, -0.02]
```

Odometer `movement` controls the main odometer wheel movement phase.

Allowed odometer movement values are:

- `instant`
- `linear`
- `ease_out`
- `bell`
- `smooth`
- `click`

Odometer `instant` means the digit display jumps immediately to the target value with no animation.

Odometer `linear` means the wheel rolls from the old digit position to the target digit position at constant speed.

Odometer `ease_out` means the wheel starts fast, then slows into the target.

Odometer `bell` means the wheel starts slow, speeds up through the middle, then slows into the target.

Odometer `smooth` is recognised only, reserved for future enhancement, and should warn then fall back to `instant`.

Odometer `click` is recognised only, reserved for future stepped-click enhancement, and should warn then fall back to `instant`.

`realism.movement_policy` is obsolete for odometer movement. Do not use or recommend it for odometers.

Existing top-level `movement` may remain supported for backwards compatibility where already present.

Unknown realism options must fail configuration loading with a clear error. Known realism options used on unsupported gauge types must also fail. Unknown movement values must fail configuration loading unless a gauge type explicitly documents a recognised fallback. Invalid `realism.order` entries must fail.

Realism options affect display behaviour only. They must not mutate source values, logs, exported values, configured ranges, or input data.

## Approved v3.5 realism options

These options are implemented or supported in v3.5:

| Option | Applies to | Summary |
|---|---|---|
| `movement` | all gauge types for parsing; concrete behaviour defined per gauge type | Single scalar movement knob. Odometer defines `instant`, `linear`, `ease_out`, `bell`, recognised `smooth`, and recognised `click`. |
| `wraparound` | odometer | Roll cleanly through digit strip boundaries, especially `9 -> 0`. |
| `drum_slop` | odometer | Static per-wheel alignment imperfection. |
| `carry_drag` | odometer | Higher digit creeps during lower digit rollover. |
| `snap_settle` | odometer | Mechanical snap into final digit position with a small settle. |
| `hysteresis` | radial, bar | Direction-dependent displayed offset without changing the source value. |
| `stiction` | radial, bar | Sticky threshold behaviour before visible movement releases. |
| `damping` | radial, bar | Lagged/smoothed response to value changes. |
| `overshoot` | radial, bar | Bounded pass-and-settle movement. |
| `peg_bounce` | radial, bar | Tiny bounded bounce at configured min/max stops. For bars this is end-stop bounce on fill/reveal extent. |
| `thermal_fade` | indicator | Soft incandescent-style on/off response. |
| `needle_shadow` | radial | Optional offset/tinted copy of the rotating needle for visual depth. |
| `calibration_offset` | radial | Optional display-only degree offset for imperfect needle alignment. |

Not implemented in v3.5:

| Option | Applies to | Status |
|---|---|---|
| `backlash` | odometer | Planned in old v3.5 docs, but not implemented on `main`; future work only. |

Numeric and segmented gauges do not get extra realism behaviour in v3.5 unless a later slice explicitly adds it. Indicator gauges support `thermal_fade`. These gauge types may still have baseline preview files.

## Odometer movement phase model

Odometer realism should be composable:

```text
route -> lead_in -> travel -> settle -> rest
```

The phase model is internal implementation structure, not the public YAML shape. Public config remains feature-oriented. Do not expose `movement.pre`, `movement.primary`, or `movement.post` unless a later docs slice explicitly changes the public config model.

For v3.5.6:

```text
default route -> none -> instant / linear / ease_out / bell -> none -> existing static offsets
```

Implemented odometer slices fit around this:

| Slice | Feature | Phase |
|---|---|---|
| v3.5.2 | `wraparound` | `route` path rule |
| v3.5.3 | `drum_slop` | `rest` static offset |
| v3.5.6 | `movement` | `travel` curve |
| v3.5.7 | `carry_drag` / 9-drag | `lead_in` / overlap movement on neighbouring wheels |
| v3.5.14 | `snap_settle` | `settle` tail |

Future odometer slices may add:

| Candidate | Phase |
|---|---|
| `backlash` | `lead_in` / `settle` direction-change slack |

The main odometer movement phase must not render by permanently feeding fractional numeric odometer values back into the source display value.

Bad model:

```text
displayValue = 5998.5
```

Good model:

```text
from digit position -> movement phase -> target digit position
```

At the end of the movement phase, the handover position must be exactly the target digit position. Later settle effects such as `snap_settle` may start from that handover position.

## Bar movement model

Bar realism should work in displayed fill/reveal extent, not by mutating source values or configured ranges.

For bar gauges:

```text
source value -> target fill/reveal extent -> finite movement behaviour -> rendered fill/reveal extent
```

Radial implementations may be used as references for timing, lifecycle, transition shape, and clamping where sensible. They must not be copied as angle/needle rendering logic. Bar render output is an extent, clip, reveal, or transform appropriate to the existing bar renderer.

## Realism ordering

`realism.order` may optionally define the order in which enabled realism behaviours are applied.

Example:

```yaml
realism:
  order:
    - hysteresis
    - stiction
    - damping
    - overshoot
    - peg_bounce
```

Do not rely on YAML key order to control behaviour order.

If `realism.order` is omitted, use the documented default order for the gauge type.

Default radial order:

```text
hysteresis
stiction
damping
overshoot
peg_bounce
calibration_offset
needle_shadow
```

Default bar order:

```text
hysteresis
stiction
damping
overshoot
peg_bounce
```

Default odometer order:

```text
route: wraparound
lead_in: carry_drag
travel: movement
settle: snap_settle
rest: drum_slop
```

Default indicator order:

```text
thermal_fade
```

The default order is a starting point. `99-all-options` preview files exist so the combined feel can be judged visually and tuned later.

## Scope doctrine

Do:

- add one behaviour per slice;
- keep every change deterministic and bounded;
- tick only while a finite transition is active;
- keep preview YAML files simple and normal-looking;
- use Gauge Preview Mode to judge whether behaviour looks right.

Do not:

- add idle vibration, random flicker, shimmer, multiplex flicker, or gas-discharge jitter;
- add dashboard startup sweep, brownout, or lazy power-off in v3.5;
- add a general physics engine;
- add arbitrary combinations of preview cases;
- add preview metadata unless there is no sane alternative.

## Gauge Preview Mode

Gauge Preview Mode is the simple visual way to see what a gauge realism option does.

It is not a test harness. It is a viewer.

It loads one normal dashboard/gauge YAML file, renders one gauge, and lets the user manually change the displayed value.

CLI:

```text
godrivelog dashboard preview <file>
```

`<file>` is mandatory and must point to a normal dashboard/gauge YAML file.

Optional CLI flags may select a gauge inside the file, override the starting value, or tune preview step sizes. The YAML file remains the source of realism configuration.

Suggested optional flags:

```text
--gauge <name>
--value <number>
--step <number>
--fine-step <number>
--coarse-step <number>
```

The preview starts each gauge at the midpoint of its configured value range unless `--value` is supplied.

Suggested controls:

- Left arrow: jump to min.
- Right arrow: jump to max.
- Up arrow: increment.
- Down arrow: decrement.
- Shift + Up/Down: coarse increment/decrement.
- Ctrl/Cmd + Up/Down: fine increment/decrement.
- R: reset to midpoint.
- Space: replay last transition.
- Esc/Q: quit.
- Mouse wheel may also increment/decrement.

Do not add screenshots, videos, visual diff reports, or generated art in v3.5.

## Preview files

Each gauge type should have normal YAML preview files under a dedicated examples directory.

Suggested location:

```text
examples/gauge-realism/
  odometer/
  radial/
  bar/
  numeric/
  segmented/
  indicator/
```

Each relevant gauge type gets:

- one baseline case with no realism option enabled;
- one file per single realism option;
- one deliberate `99-all-options` case to see whether the full stack looks good or turns into brass soup.

The all-options case is a taste test, not a debugging case. If it looks wrong, inspect the single-option previews first, then adjust order, defaults, or composition rules.

Numeric and segmented gauges should only get baseline previews in v3.5 unless a specific realism option applies to them. Indicator gauges may also get a `thermal_fade` preview.

## Version slices

| Version | Slice | Summary |
|---|---|---|
| v3.5.0 | Movement realism docs | Add v3.5 planning docs and prompt structure. |
| v3.5.1 | Gauge Preview Mode | Add single-gauge interactive preview mode. |
| v3.5.2 | Odometer wraparound | Roll through digit strip boundaries cleanly. |
| v3.5.3 | Odometer drum slop | Add fixed per-wheel alignment offsets. |
| v3.5.4 | Finite movement lifecycle | Add static -> changed -> moving -> settled lifecycle. |
| v3.5.5 | Shared movement groundwork | Add reusable lifecycle helpers for finite movement; `realism.movement_policy` is not used for odometer movement. |
| v3.5.6 | Odometer movement | Define and apply the main odometer wheel movement phase. |
| v3.5.7 | Odometer carry-drag | Make higher digits creep during lower digit rollover. |
| v3.5.8 | Radial damping | Add lagged needle response. |
| v3.5.9 | Radial stiction | Add thresholded movement release for sticky needles. |
| v3.5.10 | Radial overshoot | Add bounded radial pass-and-settle movement. |
| v3.5.11 | Radial peg bounce | Add tiny bounce on radial min/max physical stops. |
| v3.5.12 | Indicator thermal fade | Add asymmetric incandescent-style on/off response. |
| v3.5.13 | Bar smoothing | Add smoothed bar movement and optional rise/fall timing. |
| v3.5.14 | Odometer snap/settle | Add mechanical snap into digit position. |
| v3.5.15 | Odometer backlash | Not implemented on `main`; moved to future work. |
| v3.5.16 | Radial display-only hysteresis | Add direction-dependent displayed offset for radial gauges. |
| v3.5.17 | Radial needle drop shadow | Draw an optional offset/tinted copy of the rotating needle behind the real needle. |
| v3.5.18 | Radial calibration offset | Add an optional display-only degree offset for imperfect needle alignment. |
| v3.5.19 | Bar overshoot | Add the missing bounded pass-and-settle movement for bar fill/reveal extent. |
| v3.5.20 | Bar hysteresis | Add direction-dependent displayed offset for bar fill/reveal extent. |
| v3.5.21 | Bar stiction | Add thresholded movement release for sticky bar fill/reveal extent. |
| v3.5.22 | Bar peg bounce | Add end-stop bounce for bar fill/reveal extent using `realism.peg_bounce`. |

## Parked for later

These are good ideas, but not v3.5:

- odometer `backlash`;
- idle needle vibration;
- gas/nixie flicker;
- LED multiplex flicker;
- power-on sweep/gauge dance;
- brownout dip;
- lazy power-off/capacitive bleed-down;
- dynamic parallax or gyro/light-driven visual effects;
- asset-only presentation work that does not need code changes.
