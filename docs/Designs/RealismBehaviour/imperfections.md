# `realism.imperfections`

Applies to: radial, bar, indicator, numeric, segmented.

Status: **desired future layer / not implemented**.

## What it would do

`realism.imperfections` would be an umbrella layer for controlled, deterministic, display-only gauge ageing, wear, vibration, electrical noise, and display artefacts.

This is not a bug mode. It is a realism layer for gauges that should look like real physical instruments: old, worn, cable-driven, electrically noisy, cheaply multiplexed, or slightly temperamental.

The goal is character without corrupting the data.

## What it simulates in real life

Real instruments are rarely perfect:

- mechanical parts wear;
- cables wobble;
- pivots stick;
- electrical supplies sag;
- displays flicker or ghost;
- bulbs fade;
- cheap drivers multiplex imperfectly;
- old gauges develop habits.

This layer would collect those broader ageing/noise/wear behaviours without pretending they are all the same kind of effect.

## Candidate categories

### Mechanical imperfections

Mechanical imperfections affect displayed position, motion, or pointer stability.

Examples:

- idle needle vibration;
- worn cable speedo wobble;
- sticky pivot;
- loose needle;
- end-stop rattle.

These mostly apply to radial gauges, though some may later map to bar gauges.

### Electrical imperfections

Electrical imperfections simulate voltage and grounding artefacts.

Examples:

- brownout dip during engine start;
- voltage sag;
- needle twitch on power transition;
- backlight shimmer;
- bad-ground flicker.

These may compose with dashboard-provided power lifecycle events, but the visible response belongs to the gauge.

### Display technology imperfections

Display technology imperfections simulate quirks of specific visual technologies.

Examples:

- LED multiplex flicker;
- gas-discharge jitter;
- LCD ghosting;
- seven-segment uneven brightness;
- backlight PWM shimmer.

These are usually visual-only and should not change displayed gauge position.

### Intermittent imperfections

Intermittent imperfections are rare, bounded, deterministic faults.

Examples:

- brief random-looking flicker;
- momentary dimming;
- brief missing segment;
- short needle twitch.

## Good result

The gauge gains believable character while remaining readable, deterministic, replayable, and faithful to the source data.

## Bad result

The gauge becomes noisy for its own sake, hides the reading, mutates the source value, uses unbounded randomness, or turns every instrument into the same haunted arcade cabinet.

## Design notes

Avoid unbounded randomness. Use seeded pseudo-random timing so replay and screenshots can be reproduced.

Keep this as a future umbrella until each concrete imperfection has a clear name, gauge-family scope, config shape, and rendering model.
