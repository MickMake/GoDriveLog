# Dashboard v2 Decisions

## 001 - No legacy visual ownership support

Decision:
The v2 dashboard rewrite will not support the old model where each sensor/PID directly owns its visual presentation.

Reason:
The project is not in production. Legacy support would slow the rewrite and preserve the wrong abstraction.

Consequence:
Existing config examples must be rewritten to the dashboard v2 scene model, not migrated at runtime.

## 002 - Remove old widget package tree

Decision:
The old standalone widget package tree is removed as part of v2.8.0.

Reason:
The application now boots through the dashboard v2 scene configuration, and keeping unused widget code would invite accidental compatibility work.

Consequence:
Future visual work should live under the dashboard v2 scene, decoder, asset, and renderer packages.
