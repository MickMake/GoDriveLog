# `per_digit_response_lag`

Applies to: numeric, segmented.

Status: **candidate / not implemented**.

## Purpose

Simulate the behaviour of older multiplexed multi-digit displays where successive digit positions receive updated values slightly later than preceding positions.

The effect is purely visual. It changes only when each digit slot updates, never the underlying source value.

## Expected visual behaviour

When a displayed value changes:

- changed digit slots update sequentially across the display;
- each slot retains its previous rendered value until its update point;
- the delay between slots is small and deterministic;
- the complete display settles quickly on the current source value.

The effect should be subtle and should never make the display difficult to read.

## What it simulates

Older multiplexed LED and segmented displays often refreshed individual digit positions sequentially rather than updating every position simultaneously.

This feature recreates that display-driver behaviour.

It does **not** simulate mechanical movement, ageing, or wear.

## Candidate configuration

Initial implementation should support:

```yaml
realism:
  per_digit_response_lag:
    enabled: true
    direction: left_to_right
```

Supported directions:

- `left_to_right`
- `right_to_left`

If omitted, the default direction is:

```text
left_to_right
```

No user-configurable timing values are required for the initial implementation.

## Required behaviour

The implementation must be:

- display-only;
- deterministic;
- bounded;
- active only when configured;
- independent of source values;
- settled exactly on the current value when complete.

It must not mutate:

- source values;
- logs;
- exports;
- configured ranges;
- input data.

## Slot behaviour

When a value changes:

- only digit slots whose displayed value changes participate;
- unchanged digit slots remain unchanged;
- participating slots update in configured display order;
- each participating slot updates after a small fixed offset from the previous slot;
- all slots settle on the newest source value within a bounded period.

If a newer source value arrives before settling completes:

- pending updates must target the newest value;
- previous pending values must not be queued or replayed;
- settling remains bounded.

## Behaviour and appearance boundary

Code is responsible for:

- slot update order;
- slot update timing;
- settling behaviour.

Images and visual assets are responsible for:

- digit appearance;
- segment appearance;
- colours;
- glow;
- styling;
- ageing effects.

This feature must not introduce code-driven visual styling.

## Implementation requirements

If confirmed missing:

- add configuration parsing;
- validate configuration values;
- implement per-slot update timing;
- support both update directions;
- preserve existing behaviour when disabled;
- keep behaviour deterministic;
- keep implementation local unless an existing helper cleanly fits.

## Tests

Add or update tests for:

- configuration parsing;
- configuration validation;
- left-to-right updates;
- right-to-left updates;
- unchanged slots remain unchanged;
- only changed slots participate;
- newer values replace pending targets;
- deterministic behaviour;
- disabled behaviour matches current rendering;
- final displayed value exactly matches the source value.

## Preview and documentation

Add a preview showing:

- simultaneous value change;
- left-to-right update;
- right-to-left update;
- rapid successive updates;
- exact final settling.

Update relevant documentation accordingly.

## Constraints

- No hidden default behaviour.
- No randomness.
- No source-value mutation.
- No queued historical values.
- No configurable timing in the initial implementation.
- No renderer redesign.
- No styling changes.
- No image inspection.

## Good result

The display behaves like an older multiplexed segmented display whose digit positions update sequentially, while remaining readable and quickly settling on the correct value.

## Bad result

Digits update randomly, queue old values, remain visibly stale, become difficult to read, or attempt to simulate mechanical behaviour.

## Non-goals

This design does not define:

- mechanical realism;
- odometer behaviour;
- glyph animation;
- image styling;
- glow or ageing;
- segment rendering;
- configurable timing.
