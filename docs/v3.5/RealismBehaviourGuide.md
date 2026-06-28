# v3.5 Realism Behaviour Guide

This guide defines the intended visual feel of each v3.5 realism option.

Use it when implementing a slice, writing preview YAML, or judging Gauge Preview Mode output by eye.

The goal is believable gauge behaviour: small, clear details that make gauges feel like real mechanisms when values change.

## General visual rules

- Keep each behaviour subtle.
- Keep movement finite: it should visibly settle.
- Prefer simple, readable behaviour over clever behaviour.
- Do not let a realism option change source values, logs, exported values, configured ranges, or input data.
- When a combined `99-all-options` preview looks wrong, inspect the single-option previews first.

## `movement`

Applies to: all gauge types for parsing; concrete behaviour is defined per gauge type.

`movement` is the single movement knob. It should remain a collapsed scalar, not a nested policy object.

Gauge types that do not yet have concrete movement behaviour may accept `movement` and keep their current immediate display behaviour until a later slice defines more. Do not invent broad physics or idle animation to make unsupported gauge types look busy.

### Odometer movement values

Odometers use:

```yaml
odometer:
  movement: instant | linear | ease_out | bell | smooth | click
```

#### `instant`

**Visual intent:**

The digit display updates directly to the target value with no visible animation.

**Good result:**

The reading changes immediately and settles exactly on the target value.

**Bad result:**

The wheel drifts, eases, leaves a fractional resting offset, or continues ticking after the value has changed.

#### `linear`

**Visual intent:**

The wheel rolls from the previous digit position to the target digit position at constant speed.

**Good result:**

The movement is plain, deterministic, readable, and settles exactly on the target digit position.

**Bad result:**

The movement changes speed unexpectedly, rolls the long way without a path rule asking for it, or settles between digits.

#### `ease_out`

**Visual intent:**

The wheel starts moving quickly, then slows into the target.

**Good result:**

The wheel feels mechanically eased without becoming theatrical. At completion it lands exactly on the target digit position.

**Bad result:**

The wheel eases so slowly that it feels sluggish, fails to settle, or leaves a fractional display state behind.

#### `bell`

**Visual intent:**

The wheel starts slow, speeds up through the middle, then slows into the target. This is a bell-curve velocity / smoothstep-style movement.

**Good result:**

The roll feels smooth and mechanical while still being simple and bounded.

**Bad result:**

The movement feels springy, bouncy, or like a physics simulation. Bell movement is not snap/settle, backlash, or carry-drag.

#### `smooth`

**Visual intent:**

Recognised only for odometers in this slice. Reserved for future enhancement.

**Good result:**

The configuration emits a clear warning and falls back to `instant`.

**Bad result:**

The system silently treats `smooth` as a real odometer movement mode, or invents generic smoothing without a documented slice.

#### `click`

**Visual intent:**

Recognised only for odometers in this slice. Reserved for future stepped-click enhancement.

**Good result:**

The configuration emits a clear warning and falls back to `instant`.

**Bad result:**

The system silently treats `click` as a real odometer movement mode, or implements stepped clicking early.

## `wraparound`

Applies to: odometer.

**Visual intent:**

Odometer wheels roll cleanly through digit strip boundaries, especially `9 -> 0` and `0 -> 9`.

**Good result:**

A rollover looks like one continuous drum motion through the nearest boundary.

**Bad result:**

The wheel jumps, reverses unexpectedly, rolls the long way around, or briefly shows an impossible digit position.

## `drum_slop`

Applies to: odometer.

**Visual intent:**

Each odometer wheel may sit with a small fixed alignment imperfection, as if the drums are not perfectly centred in the window.

**Good result:**

Digits look slightly mechanical and imperfect while still being readable.

**Bad result:**

Digits become hard to read, offsets change between runs, or wheels drift while idle.

## `carry_drag`

Applies to: odometer.

**Visual intent:**

A higher digit begins to creep as the lower digit approaches rollover.

**Good result:**

The next wheel looks lightly dragged by the rolling lower wheel, then lands in the correct final digit.

**Bad result:**

The higher digit moves too early, moves too far, or appears to change value before the rollover is visually justified.

## `snap_settle`

Applies to: odometer.

**Visual intent:**

An odometer wheel lands into its final digit position with a small mechanical snap and quick settle.

**Good result:**

The wheel feels like it has clicked into place.

**Bad result:**

The wheel bounces repeatedly, overshoots so far it becomes distracting, or keeps moving after it has settled.

## `backlash`

Applies to: odometer.

**Visual intent:**

When direction changes, the wheel shows a tiny amount of slack before movement takes up in the new direction.

**Good result:**

The change of direction feels slightly mechanical without making the reading confusing.

**Bad result:**

The wheel loses numeric correctness, delays too long, or appears disconnected from the value change.

## `hysteresis`

Applies to: radial, bar.

**Visual intent:**

The displayed position may rest slightly differently depending on whether the value approached from below or above.

**Good result:**

A rising value and falling value can settle with a tiny direction-dependent display difference, while still clearly representing the same source value.

**Bad result:**

The displayed offset is large enough to look wrong, changes the source value, or accumulates error over time.

## `stiction`

Applies to: radial.

**Visual intent:**

A needle resists very small changes, then releases once the change is large enough to notice.

**Good result:**

Tiny value changes may not move the needle immediately. When it does move, it makes a small catch-up movement and settles.

**Bad result:**

The needle sticks during large changes, jumps violently, or behaves unpredictably.

## `damping`

Applies to: radial, bar.

**Visual intent:**

The displayed position lags behind the source value and catches up smoothly.

**Good result:**

Needles and bars feel weighty and calm during value changes.

**Bad result:**

The display feels sluggish, never reaches the target, or keeps drifting after the value has settled.

## `overshoot`

Applies to: radial, bar.

**Visual intent:**

The displayed position may pass the target slightly, then return and settle.

**Good result:**

The movement gives a small sense of momentum without stealing attention.

**Bad result:**

The display swings too far, oscillates repeatedly, or overshoots during tiny changes where it looks silly.

## `peg_bounce`

Applies to: radial.

**Visual intent:**

When a radial needle reaches the configured minimum or maximum stop, it should appear to tap the stop, rebound slightly, then settle quickly.

**Good result:**

The bounce is small, quick, and only visible at the physical stop.

**Bad result:**

The needle bounces during ordinary in-range movement, passes through the stop, or keeps bouncing.

## `thermal_fade`

Applies to: indicator.

**Visual intent:**

An indicator turns on and off with a soft incandescent-style response.

**Good result:**

On-state appears to warm in rather than appearing instantly. Off-state fades away softly and then settles fully off.

**Bad result:**

The indicator flickers, pulses, randomly changes brightness, or remains partly on after it should be off.

## `needle_shadow`

Applies to: radial.

**Visual intent:**

A radial needle may have a fixed shadow or depth cue behind it.

**Good result:**

The needle feels visually raised from the dial face while still reading clearly.

**Bad result:**

The shadow distracts from the needle position, appears as a second needle, or makes the gauge harder to read.

## `calibration_offset`

Applies to: radial.

**Visual intent:**

A radial needle may display with a small fixed angular offset, as if the needle is not perfectly aligned.

**Good result:**

The gauge looks slightly imperfect while still clearly representing the configured value range.

**Bad result:**

The offset changes source values, pushes the needle outside sensible visual bounds, or makes the gauge look broken.
