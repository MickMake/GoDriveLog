# `ghosting` — Implementation

## Purpose
Audits whether numeric or segmented ghosting exists in current code.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
None found in current code.

## Functions and Methods
- `NumericScene`
- `SegmentedScene`

## Runtime Flow
No previous-glyph history or decay state was found in numeric or segmented runtime paths.

## Configuration
No `realism.ghosting` key was found in current code.

## Behaviour
Numeric and segmented displays render the current state only.

## Rendering
No ghost-image render path was found.

## Tests
No feature-specific tests found.

## Limitations
This record only covers current repository code.

## Deviations from Design
The design describes bounded display ghosting. Current code has no matching config or render path.

## Remaining Work
Add history state and rendering only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `NumericScene`
- `SegmentedScene`

Searches performed:
- `ghosting`
- `realism.ghosting`
