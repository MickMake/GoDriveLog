# Gauges

## Purpose

This document explains the **gauge side** of the architecture.

Here, a **gauge** is not a one-off handcrafted implementation of a specific real-world instrument. Instead, a gauge is a combination of:

- a **gauge type**
- a set of **assets**
- a set of **quirks**

The gauge type provides the physical mechanism being simulated. Assets provide the visible look. Quirks provide the behaviour.

## What a gauge type is

A **gauge type** represents a physical mechanism or display topology.

It answers questions like:

- Is this a radial pointer?
- Is this a rolling drum counter?
- Is this a bar gauge?
- Is this a segmented numeric display?
- Is this an indicator lamp?

A gauge type is therefore the runtime representation of **how the gauge is physically organised**, not what brand or exact model it is.

## What a gauge type is not

A gauge type is **not**:

- a brand
- a manufacturer
- a single researched real-world part
- a visual theme
- a finished rendered gauge

Those things may influence assets or presets, but they should not multiply the number of gauge types.

## Gauge type responsibilities

A gauge type should define:

1. **Topology**  
   What visible elements exist and how they are arranged.

2. **Supported asset concepts**  
   Which visual parts may be presented to it.

3. **Supported quirk vocabulary**  
   Which quirk names or families can be applied meaningfully.

4. **Behaviour application rules**  
   How quirks apply to that topology.

## Topology

Topology is the most important part of a gauge type.

Examples:

### Radial pointer

Topology may include:

- pivot
- needle
- scale arc
- optional markers
- peg limits

### Rolling drum counter

Topology may include:

- drum wheels
- digit windows
- carry relationship between wheels
- wrap behaviour

### Segmented / numeric display

Topology may include:

- digits or slots
- segments within digits
- neighbour relationships
- leading-zero behaviour rules

### Indicator

Topology may include:

- lamp body
- on/off state
- intensity state
- optional warm-up / cool-down behaviour

## Assets

A gauge type can accept different assets via config.

Some asset ideas are common across many gauge types, such as:

- face
- bezel
- overlay

Some asset ideas are type-specific, such as:

- needle
- drum digits
- segment masks
- glow sprite
- lamp lens

The runtime should therefore support both:

- **common asset concepts** shared across types
- **gauge-type-specific asset concepts** for special cases

## Asset rule

Assets must remain about **appearance**, not behaviour.

Examples of asset concerns:

- printed face style
- bezel finish
- overlay markings
- drum artwork
- segment shape
- lens colour
- dirt, scratches, tint, shadow, glow appearance

Examples of things that should stay out of assets:

- lag behaviour
- backlash
- carry drag
- ghost persistence rules
- load sag logic
- thermal response timing

If assets start carrying behaviour semantics, the architecture becomes muddled and difficult to test.

## Quirk vocabulary per gauge type

Not every gauge type should support every quirk.

That is not a weakness; it is part of keeping the system honest.

Examples:

- `carry_drag` makes sense for a rolling counter or similar discrete carry topology
- `peg_bounce` makes sense for a pointer gauge with a physical end stop
- `thermal_fade` makes sense for indicator-like or emissive displays
- `leading_zero_behaviour` makes sense for numeric or discrete digit displays

So each gauge type should declare its own **quirk vocabulary**.

## Plug-and-play behaviour

The goal is for a gauge to be assembled in a plug-and-play style:

```text
Gauge Type + Assets + Quirks -> Final Gauge Behaviour
```

This does **not** mean every gauge type accepts every option.

It means:

- a gauge type exposes the options that make sense for it
- assets can be swapped without rewriting behaviour code
- quirks can be adjusted to tune the feel
- presets can package common real-world behaviours

## Example

### Example: rolling drum counter

```yaml
 type: rolling_drum_counter
 assets:
   face: old_cream_counter_face
   bezel: worn_black_bezel
   overlay: dusty_glass_overlay
   drums: cream_drum_digits
 quirks:
   lag: 0.15
   drum_slop: 0.08
   carry_drag: 0.22
   wraparound: shortest_path
```

Interpretation:

- the type defines the rolling-drum mechanism
- the assets define the visual look
- the quirks define the feel

The same gauge type could then be made to feel tighter, older, lazier, or cleaner by changing quirk values without changing the type itself.

## Presets

A gauge type should support both:

- direct quirk values
- named presets

Example:

```yaml
 type: radial
 preset: worn_mechanical_pointer
```

A preset is not a new gauge type. It is simply a packaged quirk configuration for a gauge type.

## Design constraint

To minimise churn, the current runtime should evolve by **refining gauge types**, not by inventing a new architectural layer every time a new real-world example appears.

The intended direction is:

- stable set of gauge types
- flexible assets
- expandable quirk vocabulary
- reusable runtime behaviour

That allows the platform to simulate a wide variety of real gauges without becoming a sprawling collection of special-case implementations.
