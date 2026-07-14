# `ghosting`

Applies to: numeric.

## Purpose

Simulate the short-lived persistence of a previously displayed digit after a digit changes.

The effect is display-only. It does not alter the source value, formatting, calculations, logging, exports, or active digit selection.

## Expected visual behaviour

When a digit slot changes:

- the new digit appears immediately when that slot updates;
- the previous digit remains faintly visible;
- the previous digit fades smoothly to complete transparency;
- only the current digit remains once the fade completes.

The effect should be subtle and must never make the displayed value ambiguous.

## What it simulates

Many real electronic displays exhibit a small amount of visible persistence after a segment or digit changes.

Examples include:

- LCD response persistence;
- multiplexed LED persistence;
- VFD residual glow;
- aged display decay.

This feature recreates that short-lived persistence.

It does **not** simulate:

- digit bleed;
- inactive segment visibility;
- response lag;
- brightness ripple;
- mechanical behaviour.

## Candidate configuration

Initial implementation should support:

```yaml
realism:
  ghosting: true
```

If omitted or false, existing behaviour is preserved.

No user-configurable timing or alpha values are required for the initial implementation.

## Required behaviour

The implementation must be:

- numeric-only;
- display-only;
- deterministic;
- bounded;
- active only when configured;
- independent of source values.

It must never modify:

- source values;
- calculations;
- logging;
- exports;
- numeric formatting;
- active digit selection.

## Slot behaviour

Ghosting operates independently for each digit slot.

A ghost is created only when that slot's rendered digit changes.

When created:

- the previous rendered digit becomes the ghost glyph;
- the current digit renders normally;
- the ghost fades smoothly to zero opacity;
- the ghost is removed once fully faded.

Unchanged digit slots must not create ghost images.

## Rapid updates

Only one ghost glyph may exist per digit slot.

If another rendered digit change occurs before the previous ghost has completely faded:

- the existing ghost is immediately discarded;
- the newly replaced digit becomes the new ghost;
- fading restarts from the initial opacity.

Previous ghost history must never be queued or stacked.

## Interaction with response lag

Ghosting is independent of response lag.

If response lag is disabled:

- the new digit appears immediately;
- ghost fading begins immediately.

If response lag is enabled:

- ghost creation occurs when the slot actually updates;
- source-value changes alone must not create ghost images.

## Layer behaviour

The ghost layer must sit immediately beneath the active digit layer.

Expected order:

```text
background / housing
digit bleed layer
ghost layer
active digit layer
foreground / glass / bezel
```

The ghost layer:

- must not replace the active digit;
- must not replace the digit bleed layer;
- must not replace unrelated layers;
- must remain confined to the digit slot.

## Behaviour and appearance boundary

Code is responsible for:

- detecting digit changes;
- creating and removing ghost glyphs;
- fade timing;
- fade opacity;
- layer placement.

Images are responsible for:

- digit appearance;
- colours;
- glow;
- styling;
- visual identity.

Code must not recreate or reinterpret the artwork.

## Symbols

Ghosting applies only to digit slots.

Decimal points, signs, separators, units, and other symbols are excluded from the initial implementation unless they are already part of the digit asset.

## Implementation requirements

If confirmed missing:

- add configuration parsing;
- validate numeric-only usage;
- retain one previous glyph per changed digit slot;
- implement bounded fade-out;
- preserve existing behaviour when absent or false;
- keep behaviour deterministic;
- keep implementation local unless an existing helper cleanly fits;
- do not redesign the renderer.

## Tests

Add or update tests for:

- configuration parsing;
- configuration validation;
- numeric-only scope;
- rejection for unsupported gauge types;
- single ghost creation;
- unchanged slots creating no ghost;
- rapid successive updates replacing the previous ghost;
- interaction with response lag enabled;
- interaction with response lag disabled;
- deterministic behaviour;
- disabled behaviour.

## Preview and documentation

Add previews demonstrating:

- ghosting disabled;
- ghosting enabled;
- rapid digit changes;
- interaction with response lag;
- correct layer ordering.

Update relevant documentation accordingly.

## Constraints

- No source-value mutation.
- No numeric formatting changes.
- No segmented-gauge support in this slice.
- No ghost queues.
- No multiple ghosts per slot.
- No configurable timing in the initial implementation.
- No configurable alpha in the initial implementation.
- No renderer redesign.

## Good result

The previously displayed digit briefly fades away while the current digit remains clear and immediately readable.

## Bad result

Multiple ghosts accumulate, ghosting begins before the digit changes, ghost images never disappear, or the current value becomes difficult to read.

## Non-goals

This design does not define:

- segmented-gauge ghosting;
- digit bleed;
- response lag;
- brightness variation;
- persistence tuning;
- renderer architecture.