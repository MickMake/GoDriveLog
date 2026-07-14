# `leading_zero_behaviour`

Applies to: numeric.

Status: **candidate / not implemented**.

## Purpose

Control how unused high-order integer zero digits are displayed.

The feature affects only the rendered display. It never changes the underlying numeric value, formatting, logging, exports, or calculations.

## Expected visual behaviour

Leading integer zeroes may be deliberately:

- shown;
- blanked;
- dimmed.

The behaviour should feel intentional and consistent across the entire display.

## What it simulates

Many real numeric displays intentionally suppress or de-emphasise unused leading digits to improve readability.

This feature reproduces that display behaviour.

It does **not** simulate mechanical behaviour or alter numeric formatting.

## Candidate configuration

Initial implementation should support:

```yaml
realism:
  leading_zero_behaviour: show
```

Supported values:

```text
show
blank
dim
```

If omitted, existing display behaviour is preserved.

## Required behaviour

The implementation must be:

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
- formatting outside the rendered display.

## Slot behaviour

Only high-order integer zero digits are candidates.

The implementation must:

- preserve at least one visible zero when the value is zero;
- never suppress significant digits;
- never suppress fractional digits;
- never suppress digits following a decimal point;
- ignore signs, decimal separators and units.

## Behaviour modes

### `show`

Leading zeroes are rendered normally.

### `blank`

Leading zero slots are not rendered.

### `dim`

Leading zeroes remain visible but are rendered using reduced alpha.

The digit image itself is unchanged.

## Behaviour and appearance boundary

Code is responsible for:

- determining which digit slots are leading zeroes;
- selecting the configured behaviour;
- applying simple presentation transforms such as alpha reduction.

Images are responsible for:

- digit appearance;
- colours;
- styling;
- typography;
- visual identity.

This feature must not introduce new glyph rendering, styling, or replacement digit assets.

## Implementation requirements

If confirmed missing:

- add configuration parsing;
- validate configuration values;
- implement the three behaviour modes;
- preserve existing behaviour when the configuration is absent;
- keep behaviour deterministic;
- keep implementation local unless an existing helper cleanly fits.

## Tests

Add or update tests for:

- configuration parsing;
- configuration validation;
- show mode;
- blank mode;
- dim mode;
- zero value handling;
- decimal values;
- negative values;
- significant zeroes;
- deterministic behaviour;
- disabled behaviour.

## Preview and documentation

Add previews demonstrating:

- integer values;
- zero;
- decimal values;
- negative values;
- each behaviour mode.

Update relevant documentation accordingly.

## Constraints

- No source-value mutation.
- No numeric formatting changes.
- No placeholder mode.
- No randomness.
- No new digit assets.
- No renderer redesign.

## Good result

Leading zero handling is deliberate, readable and visually consistent while preserving the numeric meaning of the display.

## Bad result

Significant digits disappear, decimal values are altered, numeric meaning changes, or behaviour depends on formatting side effects.

## Non-goals

This design does not define:

- numeric formatting;
- mechanical realism;
- odometer behaviour;
- placeholder characters;
- image styling;
- renderer architecture.