# Dashboard v2 Decisions

## 001 - No legacy visual ownership support

Decision:
The v2 dashboard rewrite will not support the old model where each sensor/PID directly owns its visual presentation.

Reason:
The project is not in production. Legacy support would slow the rewrite and preserve the wrong abstraction.

Consequence:
Existing config examples must not be migrated at runtime.

## 002 - Fast instrument dashboard is primary

Decision:
The normal runtime display is the fixed 1920x480 fast instrument dashboard in `internal/ui/instrument_dashboard.go`.

Reason:
The config-scene renderer was flexible but too generic for the live vehicle dashboard path. The fast dashboard updates fixed Fyne objects directly from `sensors.StateStore` snapshots.

Consequence:
Do not add an old/new display preference. The legacy config-scene dashboard is available only through Git history at `legacy-config-scene-dashboard`.

## 003 - Remove config-scene runtime stack

Decision:
v2.0.2 removes the old config-scene runtime bridge and generic renderer packages.

Reason:
Keeping an unused renderer stack invites accidental compatibility work and makes the active runtime path harder to understand.

Consequence:
Future display work should extend the fast instrument dashboard unless a later prompt explicitly introduces a new renderer direction.
