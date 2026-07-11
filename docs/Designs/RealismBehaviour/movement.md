# `movement`

Applies to: all gauge types for parsing; concrete behaviour is defined per gauge type.

Status: **partial / family-specific**.

## What it does

`movement` is the single movement knob. It should remain a collapsed scalar, not a nested policy object.

Gauge types that do not yet have concrete movement behaviour may accept `movement` and keep their current immediate display behaviour until a later slice defines more. Do not invent broad physics or idle animation to make unsupported gauge types look busy.

## What it simulates in real life

`movement` simulates the difference between a display that redraws instantly and a physical mechanism that takes a finite amount of time to reach a new reading.

For odometers, this is the most literal: it simulates number drums rolling between digit positions. For radial gauges, it describes how the displayed needle travels from its previous displayed angle to the target angle.

## Radial movement values

Radial gauges should eventually support:

```yaml
movement: instant | linear | bell
```

Current radial behaviour is equivalent to `movement: instant`.

### Radial `instant`

The needle is rendered immediately at the target angle with no interpolation.

This must preserve existing radial behaviour and must not be overloaded to mean visible ticking or stepped motion.

### Radial `linear`

The displayed needle angle interpolates from the previous displayed angle to the target angle at constant progress.

The movement should be bounded, deterministic, and settle exactly on the target angle.

### Radial `bell`

The displayed needle starts slowly, moves faster through the middle, then slows into the target.

This is a simple smooth movement curve, not overshoot, damping, stiction, hysteresis, peg bounce, or needle trail.

### Radial implementation rules

- Movement is display-only.
- Animate displayed angle or position only.
- Do not mutate source values, logs, exported values, configured ranges, or input data.
- Keep only small bounded transition state such as previous angle, target angle, elapsed time, duration, movement mode, and active state.
- Reuse needle geometry and image assets; rotate or transform at render time.
- Do not pre-render or cache an unbounded set of intermediate needle images.
- Do not combine radial movement with other realism behaviours unless a later slice explicitly defines composition.

## Odometer movement values

Odometers use:

```yaml
odometer:
  movement: instant | linear | ease_out | bell | smooth | click
```

## `instant`

### Visual intent

The digit display updates directly to the target value with no visible animation.

### Real-world analogue

A direct electronic update, or a mechanical instrument being represented without showing its movement phase.

### Good result

The reading changes immediately and settles exactly on the target value.

### Bad result

The wheel drifts, eases, leaves a fractional resting offset, or continues ticking after the value has changed.

## `linear`

### Visual intent

The wheel rolls from the previous digit position to the target digit position at constant speed.

### Real-world analogue

A simple motor-driven or mechanically coupled odometer drum moving at steady speed.

### Good result

The movement is plain, deterministic, readable, and settles exactly on the target digit position.

### Bad result

The movement changes speed unexpectedly, rolls the long way without a path rule asking for it, or settles between digits.

## `ease_out`

### Visual intent

The wheel starts moving quickly, then slows into the target.

### Real-world analogue

A driven mechanism with friction or damping that loses speed as it lands into position.

### Good result

The wheel feels mechanically eased without becoming theatrical. At completion it lands exactly on the target digit position.

### Bad result

The wheel eases so slowly that it feels sluggish, fails to settle, or leaves a fractional display state behind.

## `bell`

### Visual intent

The wheel starts slow, speeds up through the middle, then slows into the target. This is a bell-curve velocity / smoothstep-style movement.

### Real-world analogue

A smooth servo-like or well-damped mechanism accelerating and decelerating cleanly.

### Good result

The roll feels smooth and mechanical while still being simple and bounded.

### Bad result

The movement feels springy, bouncy, or like a physics simulation. Bell movement is not snap/settle, backlash, or carry-drag.

## `smooth`

### Visual intent

Recognised only for odometers in this slice. Reserved for future enhancement.

### Real-world analogue

Not defined yet. It might eventually describe a deliberately smoother wheel response, but it currently has no distinct behaviour.

### Good result

The configuration emits a clear warning and falls back to `instant`.

### Bad result

The system silently treats `smooth` as a real odometer movement mode, or invents generic smoothing without a documented slice.

## `click`

### Visual intent

Recognised only for odometers in this slice. Reserved for future stepped-click enhancement.

### Real-world analogue

Potential future simulation of a ratcheting odometer drum or stepped detent movement.

### Good result

The configuration emits a clear warning and falls back to `instant`.

### Bad result

The system silently treats `click` as a real odometer movement mode, or implements stepped clicking early.
