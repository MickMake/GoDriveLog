# Realism Behaviour Guide

This directory records what each gauge realism option does, what visual behaviour it should produce, and what kind of real-world gauge or display behaviour it is trying to simulate.

Use these notes when:

- implementing a realism option;
- writing preview YAML;
- judging Gauge Preview Mode output by eye;
- deciding whether a new realism idea belongs in code, artwork, or future planning.

## General visual doctrine

- Keep each behaviour subtle.
- Keep movement finite: it should visibly settle.
- Prefer simple, readable behaviour over clever behaviour.
- Do not let a realism option change source values, logs, exported values, configured ranges, or input data.
- When a combined `99-all-options` preview looks wrong, inspect the single-option previews first.

## Realism options

| Option | Gauge families | Simulates | Status |
|---|---|---|---|
| [`movement`](movement.md) | odometer; other gauges parse/fallback as documented | deliberate mechanical/display movement rather than instant redraw | partial/family-specific |
| [`wraparound`](wraparound.md) | odometer | continuous rolling number drums crossing digit-strip boundaries | implemented |
| [`drum_slop`](drum-slop.md) | odometer | imperfectly aligned mechanical odometer drums | implemented |
| [`carry_drag`](carry-drag.md) | odometer | rollover coupling where a lower drum drags the next drum | implemented |
| [`snap_settle`](snap-settle.md) | odometer | wheel landing into detent/click position | implemented |
| [`backlash`](backlash.md) | odometer | direction-change slack in worn gears/drives | not implemented |
| [`hysteresis`](hysteresis.md) | radial, bar | direction-dependent mechanical/friction offset | implemented |
| [`stiction`](stiction.md) | radial, bar | static friction before visible movement releases | implemented |
| [`damping`](damping.md) | radial, bar | mass/fluid/electrical smoothing and lag | implemented |
| [`overshoot`](overshoot.md) | radial, bar | momentum carrying past target then settling | implemented |
| [`peg_bounce`](peg-bounce.md) | radial, bar | needle/end-stop contact and rebound | implemented |
| [`pointer_markers`](pointer-markers.md) | radial, bar | min/max witness markers and damped follower pointer markers | implemented |
| [`thermal_fade`](thermal-fade.md) | indicator | incandescent bulb warm-up/cool-down | implemented |
| [`needle_shadow`](needle-shadow.md) | radial | physical needle depth casting a static shadow | implemented |
| [`calibration_offset`](calibration-offset.md) | radial | slightly misaligned physical needle calibration | implemented |

## Notes on real-world simulation

These behaviours are not meant to be a physics engine. They are small visual cues that make a digital dashboard feel like it is imitating physical instruments:

- **Odometer options** usually simulate mechanical number drums, gear coupling, detents, and alignment errors.
- **Radial options** usually simulate analogue needle movement, friction, damping, stop pegs, witness markers, and physical dial imperfections.
- **Bar options** simulate a physical or electronically damped level indicator using the bar's displayed fill/reveal extent, plus optional witness markers along the bar travel.
- **Indicator options** simulate lamp behaviour, especially incandescent response, while static lamp appearance still belongs in supplied artwork.
