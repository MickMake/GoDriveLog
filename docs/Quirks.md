# Quirks

## Purpose

This document explains the **quirk side** of the architecture.

In this project, **quirks** are the behaviour parameters attached to a gauge type.

They are the vocabulary used to describe how a gauge feels when driven.

That vocabulary should remain stable. There is no need to replace it with new architecture jargon.

## What a quirk is

A quirk is a configurable behaviour characteristic.

Examples:

- lag
- backlash
- drum_slop
- carry_drag
- thermal_fade
- ghosting
- uneven_brightness
- load_sag
- leading_zero_behaviour
- wraparound

Each quirk should describe a display-side behaviour, not a source-data mutation.

## What quirks are not

Quirks are not:

- asset names
- manufacturer names
- visual themes
- raw gauge research examples

Quirks are the behaviour knobs that help simulate real physical instrument feel.

## Main rule

> **Assets give the gauge its visible look. Code gives it its behaviour or feel.**

Quirks therefore belong on the **code / runtime** side of the system.

Some quirks may require asset support for how an effect is drawn, but the decision that the effect exists should remain in runtime logic.

Example:

- `ghosting` may need a visual mask or glow treatment from assets
- but the logic for how much ghosting occurs and when it appears is still a quirk/runtime concern

## Why quirks need commonality

The original one-slice-per-quirk approach is easy to document, but it does not scale well if the aim is to simulate a very wide range of real gauges.

To keep the system practical, quirks should be grouped by shared runtime behaviour ideas.

This does **not** mean removing the word `quirks`.

It means:

- keep quirks as the user-facing vocabulary
- organise implementation around common behavioural ideas underneath

## Internal commonality

Internally, quirks can be understood as belonging to broad behaviour groups.

These are internal implementation categories, not necessarily public terminology.

### 1. Response dynamics

Quirks that affect how a display responds over time.

Examples:

- lag
- damping
- overshoot
- warm-up / cool-down style response
- per-digit response lag
- snap-settle style behaviour

### 2. Memory / hysteresis / slack

Quirks that depend on previous state or direction changes.

Examples:

- hysteresis
- stiction
- backlash
- drum_slop

### 3. Coupling rules

Quirks where one displayed element influences another.

Examples:

- carry_drag
- neighbour influence
- shared power effects
- follower-style marker behaviour

### 4. Optical artefact behaviour

Quirks that create visual residue or contamination effects.

Examples:

- ghosting
- digit_bleed
- persistence-like effects

### 5. Intensity / energy behaviour

Quirks that affect luminous output or visible intensity.

Examples:

- uneven_brightness
- load_sag
- thermal_fade when expressed visually

### 6. Boundary and policy rules

Quirks that define behaviour at limits, rollover points, or display policy choices.

Examples:

- wraparound
- peg_bounce
- leading_zero_behaviour
- calibration offset style behaviour

## Why this helps

This allows the project to keep the word **quirks** while still reducing duplicated implementation work.

In practice it means:

- different quirks may share some runtime mechanics
- the config still talks about quirks
- the user still tunes quirks
- the code avoids unnecessary duplication

## Same quirk, different gauge

The same quirk name can behave slightly differently on different gauge types.

That is expected.

For example, `backlash` on:

- a radial pointer gauge
- a rolling drum counter

belongs to the same broad behaviour idea, but will not behave identically.

That does **not** mean the project needs two unrelated implementations from scratch.

It means the runtime can reuse common mechanics while applying them differently for different gauge topologies.

## How quirks should work together

Quirks must be able to combine.

That is essential if the goal is to let the user explore behaviour space by adjusting values.

This means the runtime needs:

- a defined order of application
- bounded ranges where sensible
- clear handling of unsupported combinations
- tolerance for extreme settings that may create strange but deterministic results

The project does not need to prevent every weird combination.

It does need to make the system understandable and predictable enough that tuning remains useful.

## Order of application

A practical order is:

```text
source value
 -> boundary / policy rules
 -> response behaviour
 -> memory / slack behaviour
 -> coupling behaviour
 -> optical / intensity effects
 -> final draw
```

The exact order can evolve, but it should be documented and kept stable.

## Quirks as options

Quirks should become a stable option vocabulary per gauge type.

This means a gauge type should expose:

- which quirks are supported
- what values or modes are allowed
- what defaults exist
- which presets package common combinations

Example:

```yaml
 type: segmented_numeric
 quirks:
   per_digit_response_lag: 0.10
   ghosting: 0.15
   uneven_brightness: 0.08
   load_sag: 0.05
   leading_zero_behaviour: suppress
```

## Presets and manual tuning

The system should support both:

### Presets

Useful for quickly approximating a real behaviour.

```yaml
 type: rolling_drum_counter
 preset: worn_mechanical_odometer
```

### Manual tuning

Useful for exploration and fine control.

```yaml
 type: rolling_drum_counter
 quirks:
   lag: 0.15
   drum_slop: 0.08
   carry_drag: 0.22
```

## Research vs runtime

The research work should remain separate from runtime quirk design.

The research describes:

- real observed gauges
- known physical mechanisms
- observed quirk examples
- iconic examples

The runtime design describes:

- what quirk vocabulary is supported
- how quirks combine
- which gauge types support which quirks
- how the platform simulates behaviour

That separation prevents the runtime docs from becoming a giant research archive.

## Direction from here

The next practical step is not to invent a whole new terminology set.

It is to:

1. keep the word **quirks**
2. define a stable quirk vocabulary per gauge type
3. document supported ranges and modes
4. document how quirks combine
5. allow presets plus manual tuning

That keeps the system flexible without making it bloated or obscure.
