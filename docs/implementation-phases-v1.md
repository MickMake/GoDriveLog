# Version 1.x Implementation Phases

This document defines the staged implementation plan for GoDriveLog version 1.x.

The guiding rule is simple: implement the documented config contract first, keep the runtime behaviour small, and avoid adding future features until the basic app works on a Pi with the OBD adapter.

Config source of truth: [`docs/config.md`](config.md)

Reference config: [`config.example.yaml`](../config.example.yaml)

## Version 1.1 - YAML config loader

Status: IMPLEMENTED

Implement YAML configuration loading that matches the documented config shape one-for-one.

### Scope

- Add YAML parsing using a small, standard Go YAML dependency.
- Replace the current config structs with structs that mirror the config contract:
  - `Config`
  - `LogConfig`
  - `VehicleConfig`
  - `PIDConfig`
  - `DisplayConfig`
  - `PositionConfig`
- Load the default config file path as `config.example.yaml`.
- Keep only simple defaults:
  - `log.rotate: daily`
  - `log.directory: ./log`
- Validate obvious configuration errors at startup.

### Validation

Fail fast if:

- `vehicle.name` is empty.
- `vehicle.pids` is empty.
- `log.rotate` is not `daily`.
- `log.directory` is empty.
- A PID has an unsupported `type`.
- A PID has `type: obd` but no raw `pid` value.
- A displayed PID has no `display.style` or no valid `display.position`.

### Out of scope

- Engine start/stop detection.
- Virtual PID calculations.
- JSON backwards compatibility.
- Complex config migration.

## Version 1.2 - Runtime PID list

Status: IMPLEMENTED

Build a simple runtime PID list from `vehicle.pids`.

### Scope

- Convert the PID map into a runtime slice/list.
- Preserve both:
  - the config key, such as `rpm`
  - the raw OBD PID, such as `010C`
- Include only active OBD PIDs in the polling list.

### Active polling rule

Poll when:

```text
type == obd AND (log == true OR display.enabled == true)
```

### Virtual handling

- Accept `type: virtual` as valid config.
- Do not poll it.
- Do not calculate it.
- Do not display it unless/until virtual support exists in a future release.

### Out of scope

- Generic PID decoding.
- Engine condition dependencies.
- Virtual PID dependency graphs.

## Version 1.3 - Poll, log, and display split

Status: IMPLEMENTED

Separate polling from logging and display behaviour.

### Scope

- Poll every active OBD PID using its configured `refresh` interval in milliseconds.
- If `log: true`, write the reading to JSONL.
- If `display.enabled: true`, update the dashboard.
- Allow these combinations:
  - log and display
  - log only
  - display only
  - known but inactive

### Log record shape

Include both the sensor key and raw PID in each JSONL record:

```json
{
  "time": "2026-06-04T10:15:30Z",
  "key": "rpm",
  "pid": "010C",
  "name": "rpm",
  "value": 812.0,
  "unit": "rpm",
  "source": "mock"
}
```

### Out of scope

- Session metadata files.
- Batch logging.
- Database storage.
- Engine-state-based logging.

## Version 1.4 - Dashboard from display config

Status: IMPLEMENTED

Build the dashboard from `display.enabled` PID entries only.

### Scope

- Create widgets only for PIDs where `display.enabled: true`.
- Use:
  - `display.style`
  - `display.position.x`
  - `display.position.y`
  - `display.position.width`
  - `display.position.height`
  - `min`
  - `max`
  - `unit`
- Keep the existing simple display behaviours for:
  - `gauge`
  - `bar`
  - `graph`
- Continue to show per-sensor status:
  - waiting
  - ok
  - error
  - stale

### Out of scope

- Pretty dashboards.
- Themes.
- Layout managers.
- Drag/drop editing.
- Advanced graphing.

## Version 1.5 - Daily log rotation

Replace engine-start log rotation with simple daily rotation.

### Scope

Use the documented config:

```yaml
log:
  rotate: daily
  directory: ./log
```

Implement:

- Open a date-based JSONL log file inside `log.directory`.
- Before each write, check the current date.
- If the date changed, close the old file and open a new one.
- Keep logging simple and synchronous.

### Suggested filename shape

```text
YYYYMMDD.jsonl
```

or, if multiple runs per day need separate files later:

```text
YYYYMMDD-HHMMSS.jsonl
```

Start with the simplest workable option.

### Out of scope

- Engine-start rotation.
- Engine-stop rotation.
- Compression.
- Retention policy.
- Session metadata files.

## Version 1.6 - OBD PID mappings

Extend the current `elmobd` reader to support the PIDs in the reference config.

### Scope

Support these OBD PIDs:

| Key | PID | Unit | Meaning |
|---|---:|---|---|
| `engine_load` | `0104` | `%` | Calculated engine load |
| `coolant_temp` | `0105` | `C` | Engine coolant temperature |
| `rpm` | `010C` | `rpm` | Engine RPM |
| `speed` | `010D` | `km/h` | Vehicle speed |
| `intake_air_temp` | `010F` | `C` | Intake air temperature |
| `throttle_position` | `0111` | `%` | Throttle position |
| `fuel_level` | `012F` | `%` | Fuel tank level |
| `control_module_voltage` | `0142` | `V` | Control module voltage |

### Implementation style

- Keep this as a small switch or map in the `elmobd` adapter.
- Return numeric values and units.
- Keep config unit as the display/log unit of record.

### Out of scope

- Generic raw PID decoder.
- Manufacturer-specific PIDs.
- PID auto-discovery.
- Unit conversion.

## Version 1.7 - Documentation cleanup

Update user-facing docs after the code phases are complete.

### Scope

- Update `README.md` to reference:
  - `docs/config.md`
  - `config.example.yaml`
- Update run command examples:

```bash
./GoDriveLog -config config.example.yaml
```

- Mark old JSON configs as legacy or remove them if no longer useful.
- Document the current limitations clearly:
  - `virtual` accepted but not implemented.
  - engine start/stop detection deferred.
  - daily log rotation only.

### Out of scope

- Long tutorials.
- Hardware troubleshooting guides.
- Packaging/release automation.

## Deferred future releases

These are intentionally not part of version 1.x implementation unless explicitly promoted later:

- Virtual PID calculations.
- Engine start detection.
- Engine stop detection.
- Engine-state-based log rotation.
- Complex condition parsing.
- PID auto-discovery.
- Advanced dashboard layouts.
- Log compression/retention.
