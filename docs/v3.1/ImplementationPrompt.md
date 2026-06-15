# GoDriveLog v3.1 implementation prompt

Status: canonical planning prompt
Owner: migration implementor

## Purpose

Use this file as the single implementation prompt source for v3.1 branch chats.

This file should be read by future implementation chats. It should not be split into per-version prompt files.

## Branch-chat workflow

The planning chat owns this file.

A slice implementation chat should:

1. Confirm the target version.
2. Read the matching version section in this file.
3. Confirm the previous v3.1 PR is merged into `main`, unless the planning chat explicitly says the current branch is still open for planning updates.
4. Confirm there are no blocking open PRs.
5. Create a branch from latest `main` using the target version prefix.
6. Implement only the named slice.
7. Update `CHANGES.md`.
8. Update `docs/v3.1/MigrationState.md`.
9. Update `docs/v3.1/OpenDecisions.md` only when the slice resolves, changes, adds, or explicitly defers a decision.
10. Open a PR.
11. Stop.

Do not redesign the release plan inside a slice chat.

## Required docs update per slice

Every implementation slice must update:

- `CHANGES.md`
- `docs/v3.1/MigrationState.md`

A slice must also update `docs/v3.1/OpenDecisions.md` when it:

- resolves a decision
- changes a default position
- adds a new decision
- explicitly defers an existing decision with a reason

Do not edit `OpenDecisions.md` just to say there was no change.

## Global implementation rules

- Keep the implementation small.
- Use existing v3 foundation packages before adding new foundations.
- Keep config as data.
- Do not add YAML formulas.
- Do not let dashboard code access OBD endpoints.
- Do not let widgets own sensor reads.
- Do not remove old code unless the target slice explicitly allows retirement work.
- Add or update tests when runtime behaviour changes.
- Keep branch names prefixed with the target version.

## Standard branch format

```text
v3.1.x-short-purpose
```

## v3.1.1 dashboard and gauge test harness

Instruction to slice chat: implement only v3.1.1.

Branch prefix:

```text
v3.1.1
```

Goal:

Create a way to test dashboards, gauges, widgets, or display elements independently without OBD and without starting the full runtime.

Why this is first:

Visual output needs to be testable early. A runnable runtime is not enough if the dashboard looks bad, updates poorly, or cannot be inspected in isolation.

Expected behaviour:

- Run without OBD.
- Run without full runtime startup.
- Feed dummy values into the v3 dashboard path where practical.
- Support a `sweep` pattern from min to max to min over 10 seconds.
- Support a `heartbeat` rhythm pattern for peak/response testing.
- Support or prepare for cadence options including 50ms and 100ms.
- Keep this test tool separate from production polling.

OpenDecisions.md:

- Update only if the slice resolves or changes the harness shape or cadence decision.

Do not:

- Implement the full display adapter slice.
- Wire the full command runtime.
- Change OBD polling.
- Retire old UI paths.
- Change v3 schema unless the slice explicitly proves it is required.

Acceptance checks:

- A developer can exercise at least one selected dashboard/widget path using dummy data.
- Sweep and heartbeat patterns are implemented or clearly scaffolded with tests/docs.
- The path is suitable for later visual/performance testing.
- `CHANGES.md` and `docs/v3.1/MigrationState.md` are updated.

## v3.1.2 runnable command path

Instruction to slice chat: implement only v3.1.2.

Branch prefix:

```text
v3.1.2
```

Goal:

Wire the existing v3 foundation into a runnable app path for one selected vehicle.

Expected behaviour:

- Load v3 config.
- Select one vehicle.
- Resolve RuntimePlan.
- Connect the configured endpoint.
- Start the sensor polling runtime.
- Wire selected JSONL logs to sensor events.
- Expose a dashboard output boundary, even if the visible adapter lands separately.
- Shut down cleanly.

OpenDecisions.md:

- Update only if the slice resolves or changes the minimum runnable path decision.

Do not:

- Build a new config system.
- Rebuild RuntimePlan.
- Rebuild the sensor runtime.
- Implement display adapter details beyond the boundary needed for wiring.
- Retire old command/runtime paths.

Acceptance checks:

- The v3 path can run through the selected vehicle pipeline.
- Old runtime path remains available.
- Tests or documented manual verification cover the new path.
- `CHANGES.md` and `docs/v3.1/MigrationState.md` are updated.

## v3.1.3 display adapter

Instruction to slice chat: implement only v3.1.3.

Branch prefix:

```text
v3.1.3
```

Goal:

Show v3 dashboard scene output through a practical display adapter.

Expected behaviour:

- Consume v3 dashboard scene output.
- Keep display concerns below the dashboard runtime boundary.
- Provide enough visible output to prove selected dashboard rendering works.
- Reuse or deliberately reject useful old Fyne renderer behaviour.

OpenDecisions.md:

- Update when the display adapter target is resolved, changed, or deferred.

Do not:

- Let dashboard code read sensors directly.
- Let display code access OBD endpoints.
- Rebuild the dashboard runtime.
- Retire old Fyne renderer paths yet.

Acceptance checks:

- A selected v3 dashboard can be shown through the adapter or a minimal visible path.
- Sensor and endpoint boundaries are preserved.
- Renderer caching/resource decisions are documented if relevant.
- `CHANGES.md` and `docs/v3.1/MigrationState.md` are updated.

## v3.1.4 dashboard update performance target

Instruction to slice chat: implement only v3.1.4.

Branch prefix:

```text
v3.1.4
```

Goal:

Support fast dashboard update cadence on Raspberry Pi 4 class hardware.

Performance target:

- Preferred dashboard update cadence: 50ms, about 20Hz.
- Minimum acceptable cadence: 100ms, about 10Hz.
- Rendering must not block OBD polling or logging.

Expected behaviour:

- Measure, structure, or test the dashboard update path so the target is realistic.
- Use the dashboard/gauge test harness where useful.
- Avoid broad rewrites until there is measurable evidence.

OpenDecisions.md:

- Update when the 50ms or 100ms cadence target is resolved, changed, or explicitly deferred.

Do not:

- Add dashboard polling as the solution.
- Add YAML formulas.
- Let widgets read sensors directly.
- Let dashboard code access OBD endpoints.

Acceptance checks:

- The target cadence is represented in code, tests, docs, or benchmark-style harness work.
- Any limitations are documented clearly.
- `CHANGES.md` and `docs/v3.1/MigrationState.md` are updated.

## v3.1.5 JSONL rotation decision

Instruction to slice chat: implement only v3.1.5.

Branch prefix:

```text
v3.1.5
```

Goal:

Decide and document whether v3.1 keeps daily JSONL rotation.

Allowed outcomes:

1. Keep exact configured JSONL path only.
2. Add an explicit log type such as `daily_jsonl`.
3. Add a documented rotation option under the v3 log definition.

OpenDecisions.md:

- Update this file because this slice exists to resolve or explicitly defer the JSONL rotation decision.

Do not:

- Silently inherit old daily rotation.
- Change config schema casually.
- Mix old logger behaviour into v3 without naming the choice.

Acceptance checks:

- The decision is documented.
- `OpenDecisions.md` is updated or the decision is closed there.
- Tests/docs cover any runtime behaviour change.
- `CHANGES.md` and `docs/v3.1/MigrationState.md` are updated.

## v3.1.6 typed sensor values

Instruction to slice chat: implement only v3.1.6.

Branch prefix:

```text
v3.1.6
```

Goal:

Decide whether numeric `float64` sensor values remain enough for v3.1.

Expected behaviour:

- Document current numeric convention.
- Identify whether boolean or status values need stronger typing.
- Avoid broad type rewrites unless a concrete runtime, display, or logging need exists.

OpenDecisions.md:

- Update this file because this slice exists to resolve or explicitly defer the typed sensor value decision.

Do not:

- Turn config into a programming language.
- Break existing numeric sensors.
- Rewrite the sensor model without a small acceptance target.

Acceptance checks:

- Decision is documented.
- Indicator/logging conventions are clear.
- Tests/docs cover any model change.
- `CHANGES.md` and `docs/v3.1/MigrationState.md` are updated.

## v3.1.7 unsupported and missing semantics

Instruction to slice chat: implement only v3.1.7.

Branch prefix:

```text
v3.1.7
```

Goal:

Decide how unsupported, missing, stale, error, and recovery states should appear in events, logs, and dashboards.

Expected behaviour:

- Decide whether unavailable sensors need explicit runtime events.
- Keep dashboard and logging interpretation consistent.
- Document mappings clearly.

OpenDecisions.md:

- Update this file because this slice exists to resolve or explicitly defer unsupported and missing semantics.

Do not:

- Invent multiple incompatible meanings for missing/unsupported.
- Hide unsupported sensors in a way that makes diagnostics useless.
- Break existing error/recovery behaviour.

Acceptance checks:

- Semantics are documented.
- Logs and dashboards use compatible status meanings.
- Tests/docs cover any runtime behaviour change.
- `CHANGES.md` and `docs/v3.1/MigrationState.md` are updated.

## v3.1.8 dashboard event efficiency

Instruction to slice chat: implement only v3.1.8.

Branch prefix:

```text
v3.1.8
```

Goal:

Reduce avoidable dashboard scene rebuild work if the display path shows it is needed.

Expected behaviour:

- Optimise after the display path and test harness make cost visible.
- Avoid rebuilding unaffected widgets or dashboards when practical.
- Preserve event-driven dashboard behaviour.

OpenDecisions.md:

- Update only if the slice discovers or resolves a decision about event efficiency.

Do not:

- Add dashboard polling.
- Add YAML formulas.
- Let widgets own sensor reads.
- Let dashboard code access endpoints.

Acceptance checks:

- Any optimisation has a clear before/after reason.
- Correctness is preserved.
- Tests cover unchanged output or reduced rebuild behaviour where practical.
- `CHANGES.md` and `docs/v3.1/MigrationState.md` are updated.

## v3.1.9 retirement readiness

Instruction to slice chat: implement only v3.1.9.

Branch prefix:

```text
v3.1.9
```

Goal:

Prepare for a later one-swoop deletion or archive PR.

Expected behaviour:

- Compare old/current paths against v3.1 replacements.
- Confirm old tests are ported, rejected, or intentionally retained.
- Confirm display path is usable or intentionally deferred.
- Confirm JSONL behaviour is decided.
- Produce a retirement-ready checklist.

OpenDecisions.md:

- Update if this slice closes, defers, or discovers decisions that affect deletion readiness.

Do not:

- Delete old paths in this slice unless explicitly approved.
- Archive files casually.
- Remove tests that still capture wanted behaviour.

Acceptance checks:

- Retire-last paths are reviewed.
- Deletion readiness is explicit.
- Remaining blockers are named.
- `CHANGES.md` and `docs/v3.1/MigrationState.md` are updated.

## Standard completion summary

- Target version.
- Branch created from latest `main`.
- Files changed.
- Behaviour changed or not changed.
- Tests run or not run.
- Docs updated.
- OpenDecisions changed or not applicable.
- PR opened.
