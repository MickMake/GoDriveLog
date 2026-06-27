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

Applies to: relevant gauge types.

### `click`

**Visual intent:**

The gauge display updates directly to the next displayed value unless another enabled realism option adds visible movement.

**Good result:**

The display feels crisp and stable. It should suit gauges that step between positions or values.

**Bad result:**

The display unexpectedly drifts, eases, or animates when no other realism option asks it to.

### `smooth`

**Visual intent:**

The gauge display may interpolate continuously where the gauge type supports continuous movement.

**Good result:**

The display visibly moves between values in a clean, controlled way.

**Bad result:**

The display keeps moving after it should have settled, or creates movement on gauge types that cannot sensibly show it.

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
