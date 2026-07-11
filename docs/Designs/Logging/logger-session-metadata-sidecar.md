# Session Metadata Sidecar

Index: 7

Status: desired

Area: logging, replay metadata, config provenance

Effort: 4-7 Codex hours

Add a session metadata sidecar next to each GoDriveLog event log.

The sidecar should capture enough context to replay, validate, or interpret a log later, even if the active dashboard config has changed.

## Proposed shape

```json
{
  "schema": "godrivelog.session.v1",
  "vehicle_id": "caddy",
  "vehicle_name": "VW Caddy 2019 SWB",
  "started_at": "...",
  "ended_at": "...",
  "config_path": "...",
  "config_sha": "...",
  "sensors": {
    "rpm": {
      "pid": "010C",
      "unit": "rpm",
      "min": 0,
      "max": 8000
    }
  }
}
```

## Rules

- Keep high-volume sensor events in the `.gdl.jsonl` file.
- Keep session-level context in `.gdl.meta.json`.
- Do not duplicate full session metadata onto every event line.
- Capture enough sensor mapping information to support replay and conversion audits.
- Treat the sidecar as optional for reading older logs but preferred for new logs.

## Possible future slice

```text
v3.x session metadata sidecar
```
