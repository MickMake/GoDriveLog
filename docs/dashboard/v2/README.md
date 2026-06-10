# GoDriveLog Dashboard v2

## Current runtime display

The normal GoDriveLog dashboard runtime is the fast fixed 1920x480 instrument dashboard.

Runtime path:

```text
sensor polling -> sensors.StateStore -> internal/ui/instrument_dashboard.go -> direct Fyne canvas object updates
```

The old config-driven scene renderer is no longer a normal runtime path. It should not be restored as an old/new display preference.

## Legacy baseline

The old config-scene dashboard baseline is preserved in Git history at:

```text
legacy-config-scene-dashboard
```

Use that ref for archaeology or rollback comparison, not as production runtime fallback code.

## Read in this order

1. `renderer-v2-overview.md`
2. `repo-structure-guardrails.md`
3. `prompts/v2.0.0-fast-instrument-renderer.md`

## Retired material

Older docs that describe assets, decoders, scene blocks, layers, or the generic Fyne scene renderer are historical planning/reference material unless explicitly reintroduced in a future prompt.

Current rule:

> Sensors produce state. The fast instrument dashboard consumes state directly.
