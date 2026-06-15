# GoDriveLog v3.1 open decisions

Status: planning
Owner: migration implementor

## Purpose

This file tracks design decisions that remain open for v3.1.

## Active decisions

### 1. Dashboard and gauge test harness shape

Question: should the v3.1 test harness target whole dashboards, individual widgets, individual gauges, or all of these?

Default until decided: support the smallest useful path first, but do not block later dashboard/widget selection.

Required dummy data patterns:

- `sweep`: min to max to min over 10 seconds.
- `heartbeat`: pulse or rhythm pattern for peak/response testing.

### 2. Dashboard update cadence

Question: can v3.1 support 50ms updates on Raspberry Pi 4, or should 100ms be the initial accepted target?

Default until measured: design for 50ms, accept 100ms if 50ms is not realistic without unsafe complexity.

### 3. JSONL rotation

Question: should v3.1 keep daily JSONL rotation, use exact configured paths only, or add an explicit rotation option?

Default until decided: exact configured path only.

### 4. Sensor value typing

Question: are numeric sensor values enough for v3.1, or does v3.1 need typed values for boolean and status sensors?

Default until decided: keep numeric values unless a concrete display or logging need proves otherwise.

### 5. Unsupported and missing sensors

Question: should unavailable sensors produce explicit runtime events, or remain represented through error and missing state handling?

Default until decided: keep current status handling and document any mapping clearly.

### 6. Display adapter target

Question: should the first v3.1 display adapter be Fyne, headless, or both?

Default until decided: prefer the smallest visible adapter that proves v3 dashboard output can be displayed.

### 7. Minimum runnable path

Question: what is the smallest acceptable runnable v3.1 app path?

Default until decided: selected vehicle, endpoint connector, sensor polling runtime, selected log output, and selected dashboard output boundary.
