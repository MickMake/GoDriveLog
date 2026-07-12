# `segment_bleed` / `digit_bleed` — Implementation

## Purpose
Audits whether segment or digit bleed exists in current code.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/scene.go`

## Types
None found in current code.

## Functions and Methods
- `NumericScene`
- `SegmentedScene`

## Runtime Flow
No inactive-mask or bleed state was found.

## Configuration
No `realism.segment_bleed` or `realism.digit_bleed` key was found in current code.

## Behaviour
Current numeric and segmented scenes render active state only.

## Rendering
No inactive-segment or inactive-digit overlay path was found.

## Tests
No feature-specific tests found.

## Limitations
This record only covers current repository code.

## Deviations from Design
The design describes segment and digit bleed. Current code does not implement them.

## Remaining Work
Add the feature only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/scene.go`

Symbols verified:
- `NumericScene`
- `SegmentedScene`

Searches performed:
- `segment_bleed`
- `digit_bleed`
