# CHANGES

## 2.0.2 - 2026-06-10

- Removed the legacy config-scene dashboard runtime bridge from `internal/ui/dashboard.go`.
- Removed the old config-scene asset registry, decoder engine, scene evaluator, and generic Fyne scene renderer packages and their old-renderer tests.
- Kept the fast 1920x480 instrument dashboard as the only normal app display path.
- Updated `config.example.yaml` so the dashboard block is a minimal schema placeholder rather than a scene-renderer configuration.
- Updated dashboard docs to point to the fast instrument renderer first and to treat the config-scene renderer as retired legacy/history.
- Confirmed the legacy baseline ref remains `legacy-config-scene-dashboard`.

## 2.0.1 - 2026-06-10

- Improved the fast 1920x480 instrument dashboard visuals for RaceDemoReader.
- Added larger RPM/speed/gear areas, right-side oil/coolant/battery/warning/failure/reset readouts, bottom status text, and colour-coded warning/failure states.
- Increased instrument dashboard refresh cadence to 50ms.
- Added small refresh guards for text colours, rectangle colours, and bar sizes.

## 2.0.0 - 2026-06-10

- Added the fixed 1920x480 fast instrument dashboard as the normal app display path.
- Created fixed Fyne canvas objects once and updated them directly from `sensors.StateStore` snapshots.
- Stopped routing normal app display through the old config-scene renderer.
- Preserved the old config-scene dashboard baseline in Git history at `legacy-config-scene-dashboard`.
