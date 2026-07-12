# Value Zones / Warning-Danger Assets — Implementation

## Purpose
Audits whether gauge-package assets switch by value zone in current code.

## Implementation Status
Not implemented.

Verified current code does not provide the designed feature in the audited scope.

## Packages and Files
- `internal/config/v3config/config.go`
- `internal/config/v3config/validate.go`
- `internal/dashboard/v3dashboard/dashboard.go`

## Types
- `WidgetConfig`
- `ZoneConfig`

## Functions and Methods
- `validateZones`
- `barCellNameForValue`

## Runtime Flow
Current zone handling is limited to `bar_display` widgets in the dashboard asset system, not gauge packages.

## Configuration
`WidgetConfig` declares `Zones []ZoneConfig` for dashboard `bar_display` widgets. No gauge package `zones` field or gauge-package warning/danger asset switch was found.

## Behaviour
Current code can choose a bar-set cell name by zone for `bar_display` widgets. It does not switch gauge-package assets by value zone.

## Rendering
Zone-based output is handled in `barParts` through `barCellNameForValue` for `bar_display` widgets, not through gauge package scene rendering.

## Tests
- `TestBarDisplayUsesZonesBySensorValue`

## Limitations
The repository contains a related but narrower zone feature than the design describes.

## Deviations from Design
The design describes gauge-package asset switching. Current code implements widget-level bar cell selection only.

## Remaining Work
Add gauge-package zone config and asset selection only if this design is scheduled.

## Verification Notes

Files inspected:
- `internal/config/v3config/config.go`
- `internal/config/v3config/validate.go`
- `internal/dashboard/v3dashboard/dashboard.go`

Symbols verified:
- `WidgetConfig`
- `ZoneConfig`
- `validateZones`
- `barCellNameForValue`

Configuration verified:
- `zones`
- `up_to`
- `cell`

Tests inspected:
- `TestBarDisplayUsesZonesBySensorValue`

Searches performed:
- `zones`
- `warning`
- `danger`
- `barCellNameForValue`
