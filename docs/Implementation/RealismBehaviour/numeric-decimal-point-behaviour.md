# `decimal_point_behaviour` — Implementation

## Purpose
Audits current decimal-point handling against the decimal-point behaviour design.

## Implementation Status
Partially implemented.

Verified current code implements part of the design, but the audited scope also has missing or different behaviour.

## Packages and Files
- `internal/dashboard/gauges/scene.go`
- `internal/config/v3config/validate.go`
- `internal/dashboard/v3dashboard/dashboard_test.go`

## Types
- `DigitSet`

## Functions and Methods
- `NumericScene`
- `splitTextIntoSlots`
- `formatUsesDecimalPoint`

## Runtime Flow
No decimal-point-specific runtime state was found. Decimal point handling is part of numeric formatting and scene generation.

## Configuration
Current code verifies `digit_set.decimal_point` when the configured numeric format requires a decimal point. No `realism.decimal_point_behaviour` key was found.

## Behaviour
Numeric scenes can render decimal point assets in specific slots, and validation rejects formats that need a decimal point when the digit set lacks one. No fade, bleed, or ghosting behaviour specific to decimal points was found.

## Rendering
`NumericScene` renders a `decimal_point` part after the character part for a slot and before the foreground layer.

## Tests
- `TestNumericSceneRejectsMissingDecimalPointWhenFormatNeedsIt`
- `TestDigitDecimalPointDoesNotConsumeSlot`
- `TestDigitDefaultFormatDoesNotRequireDecimalPoint`
- `TestDigitDecimalPointRendersBeforeForegroundForSlot`

## Limitations
Current code covers basic decimal-point rendering and validation only.

## Deviations from Design
The design calls for decimal-point-specific realism behaviour. Current code only implements base decimal-point layout and asset requirements.

## Remaining Work
Add a dedicated decimal-point realism contract only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/scene.go`
- `internal/config/v3config/validate.go`
- `internal/dashboard/v3dashboard/dashboard_test.go`

Symbols verified:
- `DigitSet`
- `NumericScene`
- `splitTextIntoSlots`
- `formatUsesDecimalPoint`

Configuration verified:
- `digit_set.decimal_point`
- `format`

Tests inspected:
- `TestNumericSceneRejectsMissingDecimalPointWhenFormatNeedsIt`
- `TestDigitDecimalPointDoesNotConsumeSlot`
- `TestDigitDefaultFormatDoesNotRequireDecimalPoint`
- `TestDigitDecimalPointRendersBeforeForegroundForSlot`

Searches performed:
- `decimal_point_behaviour`
- `decimal_point`
- `splitTextIntoSlots`
