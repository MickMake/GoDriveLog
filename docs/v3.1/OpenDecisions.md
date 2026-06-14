# GoDriveLog v3.1 open decisions

Status: planning stub
Owner: migration implementor

## Purpose

This file tracks design decisions that remain open for v3.1.

## Active decisions

### 1. JSONL rotation

Question: should v3.1 keep daily JSONL rotation, use exact configured paths only, or add an explicit rotation option?

Default until decided: exact configured path only.

### 2. Sensor value typing

Question: are numeric sensor values enough for v3.1, or does v3.1 need typed values for boolean and status sensors?

Default until decided: keep numeric values unless a concrete display or logging need proves otherwise.

### 3. Unsupported and missing sensors

Question: should unavailable sensors produce explicit runtime events, or remain represented through error and missing state handling?

Default until decided: keep current status handling and document any mapping clearly.

### 4. Display adapter target

Question: should the first v3.1 display adapter be Fyne, headless, or both?

Default until decided: prefer the smallest visible adapter that proves v3 dashboard output can be displayed.

### 5. Minimum runnable path

Question: what is the smallest acceptable runnable v3.1 app path?

Default until decided: selected vehicle, endpoint connector, sensor polling runtime, selected log output, and selected dashboard output boundary.
