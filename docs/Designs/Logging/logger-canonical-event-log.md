# Canonical GoDriveLog Event Log

Index: 6

Status: desired

Area: logging, sensor events, schema/versioning

Effort: 5-9 Codex hours

Promote the current JSONL event logger output into a formal, versioned GoDriveLog-owned event log format.

The native working log format should be newline-delimited JSON with one event per line. This is the format GoDriveLog core writes, validates, and replays. Other formats should be converted into this format rather than being supported directly inside the runtime.

## Proposed file naming

```text
*.gdl.jsonl
*.gdl.meta.json
```

## Proposed event schema marker

```json
{"schema":"godrivelog.event.v1"}
```

## Rules

- Treat GoDriveLog JSONL as the canonical event log, not as incidental logger output.
- Add a schema marker or schema version to every event record.
- Keep one complete event per line.
- Preserve the existing event-oriented shape: kind, sensor id, timestamps, status, typed value, previous status, and error.
- JSONL events should represent sensor events, not rendered dashboard state.
- Do not use CSV, MDF4, BLF, Parquet, ROS bag, or any other external format as the native runtime log format.
- Industry or third-party formats may be supported by converters, importers, or exporters outside the core runtime.
- Keep logs inspectable, appendable, streamable, and replayable.
- Maintain backwards compatibility or provide a clear migration path if the existing JSONL shape changes.

## Possible future slice

```text
v3.x canonical event log v1
```
