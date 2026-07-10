# JSONL Log Validation

Index: 9

Status: desired

Area: logs, CLI, schema validation

Effort: 3-5 Codex hours

Add a validator for GoDriveLog event logs before replay or conversion.

## Proposed command shape

```text
godrivelog logs validate drive.gdl.jsonl
```

## Rules

- Validate that every line is valid JSON.
- Validate known schema markers.
- Validate required fields.
- Validate timestamps are parseable.
- Validate typed value objects.
- Validate status/error semantics.
- Warn, rather than fail, on non-monotonic timestamps unless a later spec requires strict ordering.
- Produce useful line-numbered errors for converter/debugging work.

## Possible future slice

```text
v3.x GoDriveLog log validator
```
