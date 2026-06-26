# v3.5 Gauge Realism Release Plan

Status: planned

v3.5 adds gauge realism in small, inspectable slices. It builds on the v3.4 gauge package work and keeps the existing asset-driven dashboard model.

## Purpose

v3.5 improves how gauges behave when values change. It focuses on:

- static imperfection that does not require a frame tick;
- finite movement responses after value changes;
- small deterministic radial display effects that need renderer support;
- visual inspection support so each behaviour can be judged by eye.

It does not add a general physics engine, idle animation, ambient flicker, dashboard power lifecycle effects, or asset-generation work.

## Existing v3.4 behaviour to preserve

Do not break the v3.4 gauge type model:

- `numeric` remains an image-character formatted value display.
- `radial` remains a value-to-angle transform gauge.
- `odometer` remains a rolling digit/wheel-strip transform gauge.
- `indicator` remains an on/off image-selection gauge.
- `bar` remains a continuous transform/reveal gauge.
- `segmented` remains a stepped threshold image-selection gauge.

For odometers, preserve the existing base movement modes:

- `smooth` means wheel position follows fractional values continuously.
- `click` means wheel position steps/snaps to digit positions.

v3.5 realism options layer around those modes. Do not rename or replace them.

## Scope doctrine

Do:

- add one behaviour per slice;
- keep every change deterministic and testable;
- tick only while a finite transition is active;
- keep visual showcase YAML files simple and normal-looking;
- use the manual inspection harness to judge whether behaviour looks right.

Do not:

- add idle vibration, random flicker, shimmer, multiplex flicker, or gas-discharge jitter;
- add dashboard startup sweep, brownout, or lazy power-off in v3.5;
- add a general physics engine;
- add arbitrary combinations of visual inspection cases;
- add harness metadata unless there is no sane alternative.

## Visual inspection harness rule

Each gauge type should have visual inspection YAML cases under a dedicated harness examples directory.

Each gauge type gets:

- one baseline case with no realism option enabled;
- one file per single realism option;
- one deliberate `99-all-options` case to see whether the full stack looks good or turns into brass soup.

The all-options case is a taste test, not a debugging case. If it looks wrong, debug from the single-feature cases.

The harness starts each gauge at the midpoint of its configured value range.

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

## Version slices

| Version | Slice | Summary |
|---|---|---|
| v3.5.0 | Movement realism docs | Add v3.5 planning docs and prompt structure. |
| v3.5.1 | Manual gauge inspection harness | Add single-gauge interactive harness mode. |
| v3.5.2 | Odometer wraparound | Roll through digit strip boundaries cleanly. |
| v3.5.3 | Odometer drum slop | Add fixed per-wheel alignment offsets. |
| v3.5.4 | Finite movement lifecycle | Add static -> changed -> moving -> settled lifecycle. |
| v3.5.5 | Shared movement policy | Add simple policies such as immediate, linear, ease_out. |
| v3.5.6 | Odometer eased roll | Apply finite easing to odometer wheel movement. |
| v3.5.7 | Odometer carry-drag | Make higher digits creep during lower digit rollover. |
| v3.5.8 | Radial damping | Add lagged needle response. |
| v3.5.9 | Radial stiction | Add thresholded movement release for sticky needles. |
| v3.5.10 | Radial overshoot | Add bounded pass-and-settle movement. |
| v3.5.11 | Radial peg bounce | Add tiny bounce on min/max physical stops. |
| v3.5.12 | Indicator thermal fade | Add asymmetric incandescent-style on/off response. |
| v3.5.13 | Bar smoothing | Add smoothed bar movement and optional rise/fall timing. |
| v3.5.14 | Odometer snap/settle | Add mechanical snap into digit position. |
| v3.5.15 | Odometer backlash | Add direction-change slack/settle. |
| v3.5.16 | Display-only hysteresis | Add direction-dependent displayed offset without changing source value. |
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
