# `load_sag`

Applies to: numeric.

Status: **candidate / not implemented**.

## Purpose

Simulate the slight reduction in display brightness caused by increased electrical load as more display segments become illuminated.

The effect is display-only. It never alters the source value, formatting, calculations, logging, exports, or active digit selection.

## Expected visual behaviour

As more display segments become illuminated, the display gradually dims.

Examples:

```text
111     → brightest
777     → slightly dimmer
888     → dimmer
888.8   → dimmest
```

The effect should be smooth, deterministic, and always preserve readability.

## What it simulates

Some electronic displays consume more current as additional segments are illuminated.

Displays with limited power supplies, ageing electronics, long wiring runs, or marginal drivers may dim slightly under heavier electrical load.

This feature recreates that behaviour.

It does **not** simulate:

- random flicker;
- brightness ripple;
- ghosting;
- digit bleed;
- response lag;
- uneven brightness.

## Candidate configuration

Initial implementation should support:

```yaml
realism:
  load_sag:
    impact: 0.20
```

Where:

- `impact: 0.0` disables visible sag.
- `impact: 1.0` produces the maximum supported sag while remaining readable.

Values must be within:

```text
0.0 – 1.0
```

Values outside this range fail validation.

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

## Segment load model

The implementation uses a fixed segment-load table.

Recommended initial values:

| Character | Load |
|-----------|-----:|
| 0 | 6.0 |
| 1 | 2.0 |
| 2 | 5.0 |
| 3 | 5.0 |
| 4 | 4.0 |
| 5 | 5.0 |
| 6 | 6.0 |
| 7 | 3.0 |
| 8 | 7.0 |
| 9 | 6.0 |
| Decimal point | 0.5 |

The implementation must use a known load table.

It must not inspect image pixels.

## Load calculation

Display load is calculated using only the currently rendered active display.

The calculation includes:

- active digits;
- decimal points.

The calculation excludes:

- ghost glyphs;
- digit bleed;
- background layers;
- overlays unrelated to illuminated digits.

Ghost glyphs inherit the resulting display brightness but never contribute additional electrical load.

## Leading-zero behaviour

Only visible illumination contributes to display load.

Therefore:

- `show` contributes full load;
- `dim` contributes load scaled by its rendered brightness;
- `blank` contributes no load.

## Display behaviour

Load sag is applied uniformly across the illuminated display.

The brightness multiplier applies to:

- active digits;
- ghost glyphs;
- decimal-point overlays.

It does **not** apply to:

- digit bleed;
- housing;
- glass;
- bezel;
- unrelated overlays.

## Behaviour and appearance boundary

Code is responsible for:

- calculating display load;
- applying the configured impact;
- determining the resulting brightness multiplier;
- keeping behaviour deterministic.

Images are responsible for:

- digit appearance;
- colour;
- glow;
- styling;
- visual identity.

Code must not reinterpret artwork or inspect rendered pixels.

## Implementation requirements

If confirmed missing:

- add configuration parsing;
- validate numeric-only usage;
- validate impact range;
- implement the fixed segment-load table;
- calculate display-wide electrical load;
- include decimal-point load;
- exclude ghost contribution;
- support leading-zero load behaviour;
- preserve existing behaviour when absent;
- keep behaviour deterministic;
- keep implementation local unless an existing helper cleanly fits;
- do not redesign the renderer.

## Tests

Add or update tests for:

- configuration parsing;
- configuration validation;
- impact validation;
- numeric-only scope;
- deterministic load calculation;
- decimal-point contribution;
- ghost exclusion;
- leading-zero behaviour;
- maximum load;
- minimum load;
- absent configuration behaviour.

## Preview and documentation

Add previews demonstrating:

- low-load values;
- medium-load values;
- maximum-load values;
- different impact settings;
- interaction with leading-zero behaviour.

Update relevant documentation accordingly.

## Constraints

- No source-value mutation.
- No numeric formatting changes.
- No segmented-gauge support in this slice.
- No runtime randomness.
- No pixel inspection.
- No renderer redesign.

## Good result

Displays with heavier electrical load become noticeably dimmer in a smooth, believable way while remaining easy to read. Higher `impact` values can intentionally simulate ageing or faulty display electronics.

## Bad result

Brightness changes randomly, flickers, depends on ghost images, alters numeric meaning, becomes unreadable, or requires image analysis.

## Non-goals

This design does not define:

- segmented-gauge load sag;
- brightness ripple;
- flicker;
- ghosting;
- digit bleed;
- uneven brightness;
- response lag;
- renderer architecture.