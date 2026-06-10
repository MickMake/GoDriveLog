# Dashboard v2 Overview

## Current direction

GoDriveLog now uses a fast fixed 1920x480 instrument dashboard as the normal runtime display.

Runtime path:

```text
sensor polling -> sensors.StateStore -> internal/ui/instrument_dashboard.go -> direct Fyne canvas object updates
```

This replaces the older configurable scene renderer path:

```text
sensor state -> decoders -> scene evaluation -> generic renderer -> Fyne canvas updates
```

That older path was useful for experiments, but it is no longer the preferred or normal runtime direction.

## Active implementation

The fast instrument dashboard:

- constructs its Fyne object graph once;
- reads sensor values from `sensors.StateStore`;
- updates existing `canvas.Text` and `canvas.Rectangle` objects directly;
- opens at 1920x480;
- is wired from `cmd/GoDriveLog/main.go` as the app display path;
- is designed for the RaceDemoReader and live OBD state boundary.

## Legacy baseline

The old config-scene dashboard baseline is preserved in Git history at:

```text
legacy-config-scene-dashboard
```

Use that ref for rollback comparison or archaeology. Do not reintroduce it as an old/new runtime preference.

## Version map

| Version | Status | Result |
|---:|---|---|
| v2.0.0 | merged | Added fast 1920x480 instrument dashboard skeleton as the normal display path. |
| v2.0.1 | merged | Improved visual layout, warning/failure display, and race-demo readability. |
| v2.0.2 | active | Removes retired config-scene renderer runtime code and stale docs/config examples. |
| v2.0.3 | future | Optional polish, profiling, or asset refinement after the fast path is stable. |

## Guardrails

Do not:

- reintroduce the old config-scene renderer as a runtime preference;
- route the fast dashboard through `scene.Evaluate`;
- change live OBD behaviour during dashboard cleanup;
- change RaceDemoScenario behaviour during dashboard cleanup;
- rewrite sensor polling or logging;
- keep dead renderer code merely because it exists.

Do:

- keep the fast display path simple and explicit;
- update only files needed for the cleanup stage;
- keep fallback through Git history, not production runtime code;
- document any remaining legacy config/schema pieces honestly.

## Remaining legacy pieces

`internal/config/dashboard.go` may remain temporarily because config loading still validates a small `dashboard` block. That block is now a schema placeholder, not the active renderer configuration.

The active dashboard itself is implemented in:

```text
internal/ui/instrument_dashboard.go
```
