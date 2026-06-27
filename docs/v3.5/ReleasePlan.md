# v3.5 Gauge Realism Release Plan

Status: planned

v3.5 adds gauge realism in small, inspectable slices. It builds on the v3.4 gauge package work and keeps the existing asset-driven dashboard model.

## Purpose

v3.5 improves how gauges behave when values change. It focuses on:

- static imperfection that does not require a frame tick;
- finite movement responses after value changes;
- small deterministic display effects that need renderer support;
- simple preview support so each behaviour can be judged by eye.

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

All v3.5 realism options live under the `realism` key.

Keep the config collapsed where possible. Simple options may be booleans or scalars. Expand an option only when it needs settings.

Preferred shape:

```yaml
realism:
  movement: click
  wraparound: true
  damping:
    duration_ms: 300
```

`realism.movement` is the base movement mode. It is a collapsed scalar, not a nested object.

Allowed movement values are:

- `click`
- `smooth`

If omitted, `movement` defaults to `click`.

`click` means the display updates directly to the next value with no visible transition unless another enabled realism option adds one.

`smooth` means the display may interpolate continuously where that gauge type supports it.

The existing top-level `movement` field may remain supported for backwards compatibility, but new v3.5 configuration should use `realism.movement`.

Unknown realism options must fail configuration loading with a clear error. Known realism options used on unsupported gauge types must also fail. Invalid `movement` values and invalid `realism.order` entries must fail.

Realism options affect display behaviour only. They must not mutate source values, logs, exported values, configured ranges, or input data.

## Approved v3.5 realism options

These options are in scope for v3.5 and should be defined, not parked:

| Option | Applies to | Summary |
|---|---|---|
| `movement` | relevant gauge types | Base display movement mode, default `click`, optional `smooth`. |
| `wraparound` | odometer | Roll cleanly through digit strip boundaries, especially `9 -> 0`. |
| `drum_slop` | odometer | Static per-wheel alignment imperfection. |
| `carry_drag` | odometer | Higher digit creeps during lower digit rollover. |
| `snap_settle` | odometer | Mechanical snap into final digit position with a small settle. |
| `backlash` | odometer | Direction-change slack/settle. |
| `hysteresis` | radial, bar | Direction-dependent displayed offset without changing the source value. |
| `stiction` | radial | Sticky needle threshold behaviour. |
| `damping` | radial, bar | Lagged/smoothed response to value changes. |
| `overshoot` | radial, bar | Bounded pass-and-settle movement. |
| `peg_bounce` | radial | Tiny bounded bounce at configured min/max physical stops. |
| `needle_shadow` | radial | Optional offset/tinted copy of the rotating needle for visual depth. |
| `calibration_offset` | radial | Optional display-only degree offset for imperfect needle alignment. |

Numeric, indicator, and segmented gauges do not get extra realism behaviour in v3.5 unless a later slice explicitly adds it. They may still have baseline preview files.

## Realism ordering

`realism.order` may optionally define the order in which enabled realism behaviours are applied.

Example:

```yaml
realism:
  movement: smooth
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
damping
overshoot
```

Default odometer order:

```text
wraparound
carry_drag
backlash
snap_settle
drum_slop
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

Numeric, segmented, and indicator gauges should only get baseline previews in v3.5 unless a specific realism option applies to them.

## Version slices

| Version | Slice | Summary |
|---|---|---|
| v3.5.0 | Movement realism docs | Add v3.5 planning docs and prompt structure. |
| v3.5.1 | Gauge Preview Mode | Add single-gauge interactive preview mode. |
| v3.5.2 | Odometer wraparound | Roll through digit strip boundaries cleanly. |
| v3.5.3 | Odometer drum slop | Add fixed per-wheel alignment offsets. |
| v3.5.4 | Finite movement lifecycle | Add static -> changed -> moving -> settled lifecycle. |
| v3.5.5 | Shared movement policy | Add simple policies such as immediate, linear, ease_out if still needed after `realism.movement`. |
| v3.5.6 | Odometer eased roll | Apply finite easing to odometer wheel movement. |
| v3.5.7 | Odometer carry-drag | Make higher digits creep during lower digit rollover. |
| v3.5.8 | Radial damping | Add lagged needle response. |
| v3.5.9 | Radial stiction | Add thresholded movement release for sticky needles. |
| v3.5.10 | Radial/bar overshoot | Add bounded pass-and-settle movement. |
| v3.5.11 | Radial peg bounce | Add tiny bounce on min/max physical stops. |
| v3.5.12 | Indicator thermal fade | Add asymmetric incandescent-style on/off response. |
| v3.5.13 | Bar smoothing | Add smoothed bar movement and optional rise/fall timing. |
| v3.5.14 | Odometer snap/settle | Add mechanical snap into digit position. |
| v3.5.15 | Odometer backlash | Add direction-change slack/settle. |
| v3.5.16 | Display-only hysteresis | Add direction-dependent displayed offset for radial and bar gauges only. |
| v3.5.17 | Radial needle drop shadow | Draw an optional offset/tinted copy of the rotating needle behind the real needle. |
| v3.5.18 | Radial calibration offset | Add an optional display-only degree offset for imperfect needle alignment. |

## Parked for later

These are good ideas, but not v3.5:

- idle needle vibration;
- gas/nixie flicker;
- LED multiplex flicker;
- power-on sweep/gauge dance;
- brownout dip;
- lazy power-off/capacitive bleed-down;
- dynamic parallax or gyro/light-driven visual effects;
- asset-only presentation work that does not need code changes.
