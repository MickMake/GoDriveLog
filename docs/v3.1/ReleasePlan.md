# GoDriveLog v3.1 release plan

Status: planning
Owner: migration implementor

## Purpose

This file gives the initial v3.1 roadmap.

v3.1 starts after the v3.0 foundation and focuses on the remaining work needed to run the app through the v3 path.

## Release goal

Make the v3 path runnable and visible before old runtime, UI, and logging paths are retired.

## Release principles

- Use existing v3 foundation packages before adding new foundations.
- Keep every slice small enough to review in isolation.
- Keep config as data.
- Keep dashboard code below the sensor/event boundary.
- Retire old paths only after replacement behaviour is verified.

## Planned slices

| Version | Slice | Goal |
|---|---|---|
| v3.1.0 | release planning stubs | Create and bulk out v3.1 planning directory. |
| v3.1.1 | runnable command path | Wire the selected vehicle runtime path. |
| v3.1.2 | display adapter | Show v3 dashboard scene output. |
| v3.1.3 | JSONL rotation decision | Decide whether daily rotation survives. |
| v3.1.4 | typed sensor value decision | Decide whether numeric sensor values remain enough. |
| v3.1.5 | unsupported and missing semantics | Decide how unavailable sensors are represented. |
| v3.1.6 | dashboard event efficiency | Reduce avoidable scene rebuild work if needed. |
| v3.1.7 | retirement readiness review | Re-check old paths before removal or archive slices. |

## v3.1.0 acceptance checkpoints

- `docs/v3.1/` exists.
- v3.1 process docs exist.
- Carry-forward work from v3.0 is represented.
- Open design decisions are tracked.
- Completed v3.0 history is summarised, not bulk-copied.
- No code, test, schema, runtime, archive, or deletion changes.

## v3.1.1 runnable command path checkpoints

- Loads v3 config.
- Selects one vehicle.
- Resolves RuntimePlan.
- Connects the configured endpoint.
- Starts sensor polling runtime.
- Wires selected JSONL logs to sensor events.
- Exposes a dashboard output boundary, even if the first visible adapter lands later.
- Shuts down cleanly.
- Leaves old command/runtime paths available until the new path is verified.

## v3.1.2 display adapter checkpoints

- Consumes v3 dashboard scene output.
- Does not read sensors directly.
- Does not access OBD endpoints.
- Keeps display concerns below the dashboard runtime boundary.
- Provides enough visible output to prove selected dashboard rendering works.
- Documents whether old Fyne renderer caching/resource behaviour is ported or rejected.

## v3.1.3 JSONL rotation checkpoints

- Chooses exact configured path, explicit daily log type, or explicit rotation option.
- Documents the choice in `OpenDecisions.md` or closes the decision there.
- Updates config/docs only if the chosen option needs config representation.
- Does not silently inherit old daily rotation.

## v3.1.4 typed sensor value checkpoints

- Decides whether `float64` sensor values remain enough for v3.1.
- Documents any boolean/status convention used by indicators or logs.
- Avoids broad type rewrites unless a concrete runtime need exists.

## v3.1.5 unsupported and missing semantics checkpoints

- Decides whether unavailable sensors need explicit runtime events.
- Documents how unsupported, missing, stale, error, and recovery states should appear in logs and dashboards.
- Keeps logging and dashboard interpretation consistent.

## v3.1.6 dashboard event efficiency checkpoints

- Optimises only after a real display path exists.
- Avoids dashboard polling.
- Avoids YAML formulas.
- Avoids widget-owned sensor reads.
- Avoids endpoint access from dashboard code.

## v3.1.7 retirement readiness checkpoints

- Compares old/current paths against v3.1 replacements.
- Confirms old path tests have been ported, rejected, or retained intentionally.
- Confirms daily JSONL behaviour is decided.
- Confirms display path is usable.
- Confirms old runtime/UI/logger paths are safe to remove or archive.
- Produces a retirement-ready checklist before any removal slice.

## First implementation target

The first real implementation slice should be `v3.1.1-runnable-command-path`.

That slice should prove that the v3 foundation can run together under one selected vehicle.

## Non-goals for v3.1.0

- No code changes.
- No test changes.
- No schema changes.
- No runtime behaviour changes.
- No old-code removal.
- No file archiving.
