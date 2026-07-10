# External Converter Boundary

Index: 10

Status: desired

Area: `tools/converters`, import/export architecture

Effort: 3-6 Codex hours

Keep foreign-format conversion outside GoDriveLog core runtime.

Converters should live under `tools/converters` and convert external telemetry/log formats into canonical GoDriveLog event logs.

## Proposed layout

```text
tools/
  converters/
    README.md
    csv-to-gdl-jsonl/
    racechrono-to-gdl-jsonl/
    decoded-can-csv-to-gdl-jsonl/
```

## Rules

- GoDriveLog core should understand GoDriveLog event logs, not every external telemetry format.
- Foreign formats convert into `.gdl.jsonl` plus optional `.gdl.meta.json`.
- Converters may understand CSV, RaceChrono, Torque Pro, decoded CAN CSV, racing datasets, or other third-party formats.
- Converter-specific mapping files are allowed and encouraged.
- Do not add converter dependencies to the dashboard runtime.
- Do not let a one-off converter become a production runtime dependency.
- Import mapping should be explicit enough to preserve sensor ids, units, timestamps, and source provenance.

## Possible future slice

```text
v3.x tools/converters boundary
```
