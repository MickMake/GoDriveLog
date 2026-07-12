# Numeric Display — Implementation

## Purpose

Records the current implementation of the numeric gauge display, including value formatting, digit-slot allocation, reusable digit assets, decimal-point overlays, layer ordering, validation and current limitations.

## Implementation Status

Implemented.

Current code formats numeric sensor values, maps the formatted characters into configured digit slots, renders reusable character assets, and represents a decimal point as a separate overlay attached to the preceding digit slot.

The optional realism behaviours documented elsewhere, including ghosting, leading-zero policy, load sag, bleed, uneven brightness and per-digit response lag, are not part of this implementation status.

## Packages and Files

- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/gauges/scene_test.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/v3dashboard/gauge_widget_test.go`

## Types

### `Package`

Relevant fields:

- `Type`
- `Sensor`
- `Format`
- `Size`
- `Layers`
- `DigitSet`
- `Digits`
- `Path`
- `AssetRoot`

### `DigitSet`

Defines reusable assets for a numeric display:

- `Background`
- `Characters`
- `DecimalPoint`
- `Foreground`
- `Spacing`

`Characters` maps formatted characters such as `0` through `9` and `-` to asset paths.

`DecimalPoint` is a separate asset path. It is not embedded into each character image.

### `Digits`

Defines:

- `Count`
- `Positions`

`Positions` must contain one configured position for every digit slot.

### `Scene`

Relevant fields:

- `Text`
- `DigitPositions`
- `Parts`
- `Status`
- `Error`

### `ScenePart`

Numeric scenes use these relevant fields:

- `Kind`
- `AssetPath`
- `Slot`
- `Character`
- `Position`

Relevant part kinds are:

- `ScenePartKindLayer`
- `ScenePartKindBackground`
- `ScenePartKindCharacter`
- `ScenePartKindDecimalPoint`
- `ScenePartKindForeground`

## Functions and Methods

### `NumericScene`

Builds a numeric gauge scene from a gauge `Package`, `Placement`, and `sensors.SensorState`.

It verifies the package type, placement scale, digit count and configured digit-position count before creating the scene.

### `formatValue`

Formats the live sensor value using `Package.Format`.

### `splitTextIntoSlots`

Separates formatted output into visible character slots and decimal-point ownership.

Decimal points do not consume a digit slot. Each decimal point is associated with the preceding character slot.

### `digitPosition`

Returns the configured position for a numeric digit slot.

### `stateForPackage`

Normalises the incoming sensor state for the package sensor before scene construction.

## Runtime Flow

1. The dashboard obtains the current sensor state for the numeric gauge.
2. `NumericScene` verifies that the package type is `numeric`.
3. It verifies:
   - placement scale is greater than zero;
   - `digits.count` is greater than zero;
   - `digits.positions` contains exactly one entry per digit slot.
4. Static underlay layers are added to the scene.
5. For a non-OK sensor state:
   - no live digit, background, decimal-point or foreground parts are created;
   - static overlay layers are retained;
   - the scene returns with its status and error information.
6. For an OK sensor state:
   - `formatValue` produces the display text;
   - `splitTextIntoSlots` returns the visible characters and the slots owning decimal points;
   - each configured slot is processed in order.
7. For each slot, scene parts are appended in this order:
   - optional digit background;
   - character asset, unless the slot contains a space;
   - decimal-point overlay, when required;
   - optional digit foreground.
8. Package overlay layers are appended after all digit-slot parts.
9. The dashboard renderer consumes the resulting scene parts and draws their referenced assets.

## Configuration

Numeric gauge packages use these fields:

```yaml
type: numeric
sensor: <sensor-id>
format: <printf-style format>

digit_set:
  background: <optional asset>
  characters:
    "0": <asset>
    "1": <asset>
    # ...
    "9": <asset>
    "-": <asset>
  decimal_point: <optional unless required by format>
  foreground: <optional asset>
  spacing: <integer>

digits:
  count: <positive integer>
  positions:
    - [x, y]
```

Verified configuration behaviour:

- `type` must be `numeric` for `NumericScene`.
- `digits.count` must be greater than zero.
- `digits.positions` must contain exactly `digits.count` entries.
- `format` controls the text produced from the sensor value.
- `digit_set.characters` supplies assets for rendered non-space characters.
- `digit_set.decimal_point` is required only when the formatted output contains a decimal point.
- A missing character asset causes scene construction to fail.
- A missing decimal-point asset causes scene construction to fail when the current formatted output requires it.
- A decimal point does not consume an additional configured digit position.

## Behaviour

The implementation supports:

- reusable assets for digits and minus signs;
- configured displays with two, three, four and five digit slots;
- leading zeroes produced by the configured format string;
- negative values where a `-` asset exists;
- separate per-slot background and foreground assets;
- static package underlay and overlay layers;
- decimal points rendered as overlays on their owning digit slots;
- no live numeric parts when sensor status is not OK;
- scene signatures that change when the formatted visible output changes.

Leading-zero treatment is currently determined by the format string. There is no separate `leading_zero_behaviour` implementation.

## Rendering

`NumericScene` emits a deterministic ordered list of scene parts.

For each digit slot, the order is:

1. background;
2. character;
3. decimal point;
4. foreground.

Package underlay layers appear before the digit-slot parts. Package overlay layers appear after them.

The decimal point uses the same configured position as its owning digit slot and is identified by `ScenePartKindDecimalPoint`.

The implementation does not create duplicate digit glyphs containing decimal points.

## Tests

Verified tests in `internal/dashboard/gauges/scene_test.go` include:

- `TestNumericSceneUsesPackageOwnedFormatPositionsAndStaticLayers`
- `TestNumericSceneEmitsPanelUnderDigitsAndGlassOverDigits`
- `TestNumericSceneSupportsTwoThreeFourAndFiveDigitShapes`
- `TestNumericSceneDoesNotRenderLiveDigitsForNonOKStates`
- `TestNumericSceneSignatureChangesWithFormattedOutput`
- `TestNumericSceneSignatureIncludesDigitPositionsForNonOKState`
- `TestNumericSceneRejectsMissingDecimalPointWhenFormatNeedsIt`

These tests cover formatting, configured positions, layer order, supported slot counts, non-OK states, scene signatures and required decimal-point assets.

## Limitations

- The numeric display has no independent leading-zero behaviour; leading zeroes come from the format string.
- Decimal-point position is tied to the owning digit slot. No separately configurable decimal-point offset was verified.
- No decimal-point-specific fade, brightness, bleed, ghosting or timing behaviour was found.
- No per-digit response lag was found.
- No load-sag or uneven-brightness behaviour was found.
- Asset availability is validated during scene construction when the relevant formatted character or decimal point is required.
- The audited tests verify scene construction. Hardware-specific visual appearance was not verified.

## Deviations from Design

- The core numeric-display and decimal-point-overlay architecture is implemented.
- The Design allows future realism behaviours to affect the decimal point independently. No such decimal-point-specific realism controls currently exist.
- `DigitSet.Spacing` is present in the package type, but numeric placement in `NumericScene` uses explicit `digits.positions`; no spacing-driven slot-layout path was verified.
- No verified deviation was found in the rule that a decimal point does not consume a digit slot.

## Remaining Work

The following work remains only if the corresponding separate behaviour designs are pursued:

- leading-zero behaviour beyond printf formatting;
- ghosting;
- segment or digit bleed;
- uneven brightness;
- load sag;
- per-digit response lag;
- decimal-point-specific visual timing or brightness effects.

No additional work is required for the current base decimal-point overlay architecture.

## Verification Notes

Files inspected:

- `internal/dashboard/gauges/package.go`
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/gauges/scene_test.go`
- `internal/dashboard/v3dashboard/dashboard.go`
- `internal/dashboard/v3dashboard/gauge_widget_test.go`

Symbols verified:

- `Package`
- `DigitSet`
- `Digits`
- `Scene`
- `ScenePart`
- `NumericScene`
- `formatValue`
- `splitTextIntoSlots`
- `digitPosition`
- `stateForPackage`
- `ScenePartKindBackground`
- `ScenePartKindCharacter`
- `ScenePartKindDecimalPoint`
- `ScenePartKindForeground`

Configuration verified:

- `type`
- `sensor`
- `format`
- `digit_set.background`
- `digit_set.characters`
- `digit_set.decimal_point`
- `digit_set.foreground`
- `digit_set.spacing`
- `digits.count`
- `digits.positions`

Tests inspected:

- `TestNumericSceneUsesPackageOwnedFormatPositionsAndStaticLayers`
- `TestNumericSceneEmitsPanelUnderDigitsAndGlassOverDigits`
- `TestNumericSceneSupportsTwoThreeFourAndFiveDigitShapes`
- `TestNumericSceneDoesNotRenderLiveDigitsForNonOKStates`
- `TestNumericSceneSignatureChangesWithFormattedOutput`
- `TestNumericSceneSignatureIncludesDigitPositionsForNonOKState`
- `TestNumericSceneRejectsMissingDecimalPointWhenFormatNeedsIt`

Unable to verify:

- final visual appearance across all supplied gauge artwork;
- behaviour on Raspberry Pi hardware;
- whether every example numeric gauge supplies all assets required by every possible formatted value.

