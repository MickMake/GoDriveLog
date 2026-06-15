# GoDriveLog v3.1 open decisions

Status: planning
Owner: migration implementor

## Purpose

This file tracks design decisions that remain open for v3.1.

Each decision records what it blocks, what it affects, and how long it can safely remain open.

## Active decisions

### 1. Dashboard and gauge test harness shape

Status: open

Question: should the v3.1 test harness target whole dashboards, individual widgets, individual gauges, or all of these?

Default: support the smallest useful path first, but do not block later dashboard/widget selection.

Blocks: final v3.1.0 harness shape.

Impacts: display adapter, performance testing, visual regression workflow.

Can defer until: end of v3.1.0.

Decision options:

- Whole dashboard only.
- Widget/gauge only.
- Both, with smallest useful path first.

Decision: pending.

Required dummy data patterns:

- `sweep`: min to max to min over 10 seconds.
- `heartbeat`: pulse or rhythm pattern for peak/response testing.

### 2. Dashboard update cadence

Status: open

Question: can v3.1 support 50ms updates on Raspberry Pi 4, or should 100ms be the initial accepted target?

Default: design for 50ms, accept 100ms if 50ms is not realistic without unsafe complexity.

Blocks: final v3.1.3 performance target.

Impacts: v3.1.0 harness cadence options, display adapter design, renderer update strategy.

Can defer until: v3.1.3.

Decision options:

- 50ms preferred and required.
- 50ms preferred, 100ms acceptable.
- 100ms initial target with later optimisation.

Decision: pending.

### 3. JSONL rotation

Status: open

Question: should v3.1 keep daily JSONL rotation, use exact configured paths only, or add an explicit rotation option?

Default: exact configured path only.

Blocks: v3.1.4 and old logger retirement.

Impacts: log config docs, logger tests, archive/removal readiness for `internal/logger/jsonl.go`.

Can defer until: v3.1.4.

Decision options:

- Exact configured path only.
- Explicit `daily_jsonl` log type.
- Rotation option under the v3 log definition.

Decision: pending.

### 4. Sensor value typing

Status: open

Question: are numeric sensor values enough for v3.1, or does v3.1 need typed values for boolean and status sensors?

Default: keep numeric values unless a concrete display or logging need proves otherwise.

Blocks: v3.1.5 only if runtime/display work proves numeric values are not enough.

Impacts: indicators, logs, unsupported/missing semantics, possible future non-OBD signals.

Can defer until: v3.1.5.

Decision options:

- Keep `float64` for v3.1.
- Add documented boolean/status conventions without changing the core value type.
- Add typed sensor values only if a small concrete need is proven.

Decision: pending.

### 5. Unsupported and missing sensors

Status: open

Question: should unavailable sensors produce explicit runtime events, or remain represented through error and missing state handling?

Default: keep current status handling and document any mapping clearly.

Blocks: v3.1.6.

Impacts: dashboard status display, JSONL status logging, diagnostics, retirement readiness.

Can defer until: v3.1.6.

Decision options:

- Use existing status/error/missing handling.
- Add explicit unsupported/unavailable runtime events.
- Document a hybrid mapping if needed.

Decision: pending.

### 6. Display adapter target

Status: open

Question: should the first v3.1 display adapter be Fyne, headless, or both?

Default: prefer the smallest visible adapter that proves v3 dashboard output can be displayed.

Blocks: v3.1.2.

Impacts: old Fyne renderer retirement, v3.1.0 harness output mode, v3.1.3 performance measurements.

Can defer until: start of v3.1.2.

Decision options:

- Fyne first.
- Headless first.
- Both, if still small.

Decision: pending.

### 7. Minimum runnable path

Status: open

Question: what is the smallest acceptable runnable v3.1 app path?

Default: selected vehicle, endpoint connector, sensor polling runtime, selected log output, and selected dashboard output boundary.

Blocks: v3.1.1.

Impacts: command wiring, manual verification, old command/runtime retirement readiness.

Can defer until: start of v3.1.1.

Decision options:

- Selected vehicle + endpoint + sensors + JSONL only.
- Selected vehicle + endpoint + sensors + JSONL + dashboard output boundary.
- Temporary v3 command until the main command can safely switch.

Decision: pending.