# Dashboard v2 Decisions

## 001 - No legacy display.widget support

Decision:
The v2 dashboard rewrite will not support the old per-PID display.widget model.

Reason:
The project is not in production. Legacy support would slow the rewrite and preserve the wrong abstraction.

Consequence:
Existing config examples must be rewritten, not migrated at runtime.