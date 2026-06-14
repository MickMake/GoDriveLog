# GoDriveLog v3.1 carry-forward list

Status: planning stub
Owner: migration implementor

## Purpose

This file records unfinished work carried forward from the v3.0 docs.

The original v3.0 docs remain the source history. This file keeps the active v3.1 reminders short and reviewable.

## Carried from v3.0 audits

### Runnable app path

The active app path still needs to be wired through v3 config, RuntimePlan, endpoint connection, sensor polling runtime, selected logging, and dashboard output.

### Display adapter

The v3 dashboard scene runtime exists, but v3.1 still needs a practical display adapter before old UI paths can be retired.

### JSONL rotation

The old logger supported daily rotation. The v3 logger currently writes to the configured path. v3.1 must decide whether rotation survives and how it is represented.

### Sensor value typing

Current v3 sensor state uses numeric values. v3.1 must decide whether boolean or status values need stronger typing.

### Unsupported and missing sensors

v3.1 must decide whether unavailable sensors need explicit runtime events or whether current status/error handling is enough.

### Dashboard event efficiency

Current dashboard event handling may rebuild more scene state than necessary. Optimise later only after the display path is real.

### Retirement readiness

Do not retire old runtime, UI, renderer, or logging paths until their replacement behaviour is verified in v3.1.

## Source docs

- `docs/v3/WorkingCodeInventory.md`
- `docs/v3/MigrationState.md`
- `docs/v3/RetirementAudit.md`
- `docs/v3/InverseImplementationAudit.md`
