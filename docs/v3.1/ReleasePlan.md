# GoDriveLog v3.1 release plan

Status: planning
Owner: migration implementor

## Purpose

This file gives the v3.1 implementation roadmap.

v3.1 starts after the v3.0 foundation and focuses on the remaining work needed to run, view, test, and eventually retire the old app path.

PR 51 is a docs-only planning baseline. It is not a numbered implementation slice.

## Release goal

Make the v3 path runnable, independently testable, visible, performant, and retirement-ready.

## Release principles

- Use existing v3 foundation packages before adding new foundations.
- Keep every slice small enough to review in isolation.
- Keep config as data.
- Keep dashboard code below the sensor/event boundary.
- Make visual output testable before spending days on runtime wiring.
- Target fast dashboard updates on Raspberry Pi 4 hardware.
- Retire old paths only after replacement behaviour is verified.

## Branch-chat workflow

The planning docs define the work. Each implementation chat should:

1. Read this file.
2. Read `docs/v3.1/prompts/README.md`.
3. Read the prompt file for the target slice under `docs/v3.1/prompts/`.
4. Confirm the previous relevant PR is merged into `main`.
5. Confirm there are no blocking open PRs.
6. Create a branch from latest `main` using the target version prefix.
7. Implement only that version slice.
8. Update `CHANGES.md` and `docs/v3.1/MigrationState.md`.
9. Update `docs/v3.1/OpenDecisions.md` only when a decision is resolved, changed, added, or explicitly deferred.
10. Open a PR.
11. Stop.

Do not redesign the release plan inside a slice chat.

## Planned implementation slices

| Version | Slice | Goal |
|---|---|---|
| v3.1.0 | dashboard and gauge test harness | Test gauges/widgets independently with dummy data before full runtime wiring. |
| v3.1.1 | runnable command path | Wire the selected vehicle runtime path. |
| v3.1.2 | display adapter | Show v3 dashboard scene output through a practical adapter. |
| v3.1.3 | dashboard update performance target | Support 50ms target updates or at least 100ms updates on Pi 4. |
| v3.1.4 | JSONL rotation decision | Decide whether daily rotation survives. |
| v3.1.5 | typed sensor values | Decide whether numeric sensor values remain enough. |
| v3.1.6 | unsupported and missing semantics | Decide how unavailable sensors are represented. |
| v3.1.7 | dashboard event efficiency | Reduce avoidable scene rebuild work if needed. |
| v3.1.8 | retirement readiness | Re-check old paths before one-swoop removal or archive slices. |

## PR 51 planning baseline checkpoints

- `docs/v3.1/` exists.
- v3.1 process docs exist without treating planning as an implementation version.
- Carry-forward work from v3.0 is represented.
- Open design decisions are tracked.
- Completed v3.0 history is summarised, not bulk-copied.
- Per-slice prompt files exist under `docs/v3.1/prompts/`.
- No code, test, schema, runtime, archive, or deletion changes.

## v3.1.0 dashboard and gauge test harness checkpoints

- Provides a way to test a dashboard, gauge, widget, or display element independently.
- Runs without OBD.
- Runs without full runtime startup.
- Feeds dummy data through the v3 dashboard path where practical.
- Supports a `sweep` pattern from min to max to min over 10 seconds.
- Supports a `heartbeat` rhythm pattern for peak/response testing.
- Allows update cadence selection, including 50ms and 100ms where practical.
- Keeps this tooling separate from production OBD polling.

## v3.1.1 runnable command path checkpoints

- Loads v3 config.
- Selects one vehicle.
- Resolves RuntimePlan.
- Connects the configured endpoint.
- Starts sensor polling runtime.
- Wires selected JSONL logs to sensor events.
- Exposes a dashboard output boundary, even if the visible adapter lands separately.
- Shuts down cleanly.
- Leaves old command/runtime paths available until the new path is verified.

## v3.1.2 display adapter checkpoints

- Consumes v3 dashboard scene output.
- Does not read sensors directly.
- Does not access OBD endpoints.
- Keeps display concerns below the dashboard runtime boundary.
- Provides enough visible output to prove selected dashboard rendering works.
- Documents whether old Fyne renderer caching/resource behaviour is ported or rejected.

## v3.1.3 dashboard update performance checkpoints

- Defines dashboard update targets: 50ms preferred, 100ms minimum acceptable.
- Treats Raspberry Pi 4 as the reference hardware target.
- Measures or structures the path so display rendering does not block OBD polling or logging.
- Uses the dashboard test harness to exercise visual update cadence where possible.
- Avoids premature micro-optimisation before the display/test path is measurable.

## v3.1.4 JSONL rotation checkpoints

- Chooses exact configured path, explicit daily log type, or explicit rotation option.
- Documents the choice in `OpenDecisions.md` or closes the decision there.
- Updates config/docs only if the chosen option needs config representation.
- Does not silently inherit old daily rotation.

## v3.1.5 typed sensor value checkpoints

- Decides whether `float64` sensor values remain enough for v3.1.
- Documents any boolean/status convention used by indicators or logs.
- Avoids broad type rewrites unless a concrete runtime need exists.

## v3.1.6 unsupported and missing semantics checkpoints

- Decides whether unavailable sensors need explicit runtime events.
- Documents how unsupported, missing, stale, error, and recovery states should appear in logs and dashboards.
- Keeps logging and dashboard interpretation consistent.

## v3.1.7 dashboard event efficiency checkpoints

- Optimises only after a real display path exists.
- Avoids dashboard polling.
- Avoids YAML formulas.
- Avoids widget-owned sensor reads.
- Avoids endpoint access from dashboard code.

## v3.1.8 retirement readiness checkpoints

- Compares old/current paths against v3.1 replacements.
- Confirms old path tests have been ported, rejected, or retained intentionally.
- Confirms daily JSONL behaviour is decided.
- Confirms display path is usable.
- Confirms old runtime/UI/logger paths are safe to remove or archive.
- Produces a retirement-ready checklist before any removal slice.
- Supports the final goal of a boring one-swoop deletion PR.

## First implementation target

The first real implementation slice is `v3.1.0-dashboard-gauge-test-harness`.

That slice should prove dashboard and gauge output can be inspected independently before full runtime wiring grows around it.

## Non-goals for PR 51

- No numbered implementation slice.
- No code changes.
- No test changes.
- No schema changes.
- No runtime behaviour changes.
- No old-code removal.
- No file archiving.