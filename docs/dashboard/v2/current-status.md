# Dashboard v2 Current Status

## Active stage

v2.0.2 - Legacy config-scene renderer cleanup

## Current runtime display

GoDriveLog now starts the fast fixed 1920x480 instrument dashboard from `cmd/GoDriveLog/main.go`.

Runtime path:

```text
sensor polling -> sensors.StateStore -> internal/ui/instrument_dashboard.go -> direct Fyne canvas object updates
```

The old config-driven scene renderer is no longer a normal user-facing runtime path.

## Completed

- v2.0.0 - Fast instrument renderer skeleton merged.
- v2.0.1 - Visual polish for the 1920x480 instrument dashboard merged.

## Current branch

- feature/v2-dashboard-renderer-cleanup

## Decisions

- Fast instrument dashboard is the primary display direction.
- Old config-scene dashboard fallback is via Git history/ref, not production runtime code.
- Legacy baseline ref exists as `legacy-config-scene-dashboard`.
- Live OBD behaviour remains protected.
- RaceDemoScenario behaviour remains protected.
- Sensors produce state. The fast instrument dashboard consumes state directly.

## Known risks

- Full visual verification remains manual because the connector cannot launch the Fyne window.
- Some config structs for dashboard YAML may remain temporarily if config loading/tests still require the schema shape, but they are not a normal runtime display path.
