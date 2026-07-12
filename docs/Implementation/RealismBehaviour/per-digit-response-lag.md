# `per_digit_response_lag` — Implementation

## Purpose
Audits whether per-digit response lag exists in current code.

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
No per-slot timing or lag state was found for numeric or segmented displays.

## Configuration
No `realism.per_digit_response_lag` key was found in current code.

## Behaviour
Current multi-digit displays update as one current state.

## Rendering
No per-digit lag render path was found.

## Tests
No feature-specific tests found.

## Limitations
This record only covers current repository code.

## Deviations from Design
The design calls for slot-level response lag. Current code does not implement it.

## Remaining Work
Add slot-level timing only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/dashboard/gauges/scene.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `NumericScene`
- `SegmentedScene`

Searches performed:
- `per_digit_response_lag`
- `realism.per_digit_response_lag`
