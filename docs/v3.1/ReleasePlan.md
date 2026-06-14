# GoDriveLog v3.1 release plan

Status: planning stub
Owner: migration implementor

## Purpose

This file gives the initial v3.1 roadmap.

v3.1 starts after the v3.0 foundation and focuses on the remaining work needed to run the app through the v3 path.

## Release goal

Make the v3 path runnable and visible before old runtime, UI, and logging paths are retired.

## Planned slices

| Version | Slice | Goal |
|---|---|---|
| v3.1.0 | release planning stubs | Create v3.1 planning directory. |
| v3.1.1 | runnable command path | Wire the selected vehicle runtime path. |
| v3.1.2 | display adapter | Show v3 dashboard scene output. |
| v3.1.3 | JSONL rotation decision | Decide whether daily rotation survives. |
| v3.1.4 | typed sensor value decision | Decide whether numeric sensor values remain enough. |
| v3.1.5 | unsupported and missing semantics | Decide how unavailable sensors are represented. |
| v3.1.6 | dashboard event efficiency | Reduce avoidable scene rebuild work if needed. |
| v3.1.7 | retirement readiness review | Re-check old paths before removal or archive slices. |

## First implementation target

The first real implementation slice should be `v3.1.1-runnable-command-path`.

## Non-goals for v3.1.0

- No code changes.
- No test changes.
- No schema changes.
- No old-code removal.
- No file archiving.
