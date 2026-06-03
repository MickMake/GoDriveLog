# Version 1.x Implementation Prompt Pack

Use these prompts to run each GoDriveLog implementation phase in a fresh chat.

Each prompt assumes:

- The chat has access to the GitHub connector.
- The chat should use the GitHub connector to inspect and edit files in `MickMake/GoDriveLog`.
- The chat must not assume it can execute shell commands, run `go build`, run tests, or access a local checkout.
- The chat should make focused commits directly to `main` unless instructed otherwise.
- The chat should keep the implementation simple and aligned with `docs/config.md`.
- The chat should avoid adding future-release features.

Reference docs:

- `docs/config.md`
- `docs/implementation-phases-v1.md`
- `config.example.yaml`

---

## Prompt for Version 1.1 - YAML config loader

```text
You are implementing GoDriveLog version 1.1 in GitHub repo MickMake/GoDriveLog.

Important constraints:
- You can access and edit files using the GitHub connector.
- Do not assume you can execute shell commands, run go build, run go test, or use a local checkout.
- Commit focused changes directly to main.
- Keep it SIMPLE.
- Do not implement virtual PID calculations.
- Do not implement engine start/stop detection.
- Do not implement daily log rotation yet unless unavoidable for config compatibility.

Read these files first using the GitHub connector:
- docs/config.md
- docs/implementation-phases-v1.md
- config.example.yaml
- internal/config/config.go
- cmd/GoDriveLog/main.go
- go.mod

Task:
Implement the YAML config loader phase.

Requirements:
1. Add YAML parsing dependency using gopkg.in/yaml.v3 in go.mod.
2. Replace the old config structs with structs that mirror docs/config.md one-for-one:
   - Config
   - LogConfig
   - VehicleConfig
   - PIDConfig
   - DisplayConfig
   - PositionConfig
3. The default config path in cmd/GoDriveLog/main.go should become config.example.yaml.
4. Load YAML config from disk.
5. Keep simple defaults only:
   - log.rotate defaults to daily if empty
   - log.directory defaults to ./log if empty
   - obd_address defaults to serial:///dev/ttyUSB0 if empty
6. Validate obvious configuration errors:
   - vehicle.name must not be empty
   - vehicle.pids must not be empty
   - log.rotate must be daily
   - log.directory must not be empty
   - PID type must be obd or virtual
   - type obd must have pid
   - refresh must be positive for active PIDs
   - max must be greater than min
   - if display.enabled is true, style must not be empty and position width/height must be positive
7. Do not update the rest of the app logic except what is necessary to compile against renamed config fields as far as can be reasoned statically.

After editing, summarize:
- files changed
- key structs added
- known follow-up work for phase 1.2
- any areas that need CI/build verification because you could not run commands
```

---

## Prompt for Version 1.2 - Runtime PID list

```text
You are implementing GoDriveLog version 1.2 in GitHub repo MickMake/GoDriveLog.

Important constraints:
- You can access and edit files using the GitHub connector.
- Do not assume you can execute shell commands, run go build, run go test, or use a local checkout.
- Commit focused changes directly to main.
- Keep it SIMPLE.
- Do not implement virtual PID calculations.
- Do not implement engine start/stop detection.
- Do not implement daily log rotation yet.

Read these files first using the GitHub connector:
- docs/config.md
- docs/implementation-phases-v1.md
- config.example.yaml
- internal/config/config.go
- cmd/GoDriveLog/main.go
- internal/sensors/reader.go
- internal/ui/dashboard.go

Task:
Implement the runtime PID list phase.

Requirements:
1. Add a small runtime representation for active PIDs. It should preserve:
   - config key, e.g. rpm
   - raw OBD PID, e.g. 010C
   - unit
   - refresh interval
   - log flag
   - display config
   - min and max
2. Build the runtime list from vehicle.pids.
3. Poll only PIDs matching:
   type == obd AND (log == true OR display.enabled == true)
4. Accept type: virtual but skip it for now.
5. Keep inactive PIDs in config but out of the runtime polling list.
6. Do not create a complicated registry or plugin system.
7. Update app code only enough so the polling loop can iterate the new runtime PID list instead of the old sensors slice.

After editing, summarize:
- files changed
- runtime PID shape
- polling selection rule
- known follow-up work for phase 1.3
- any areas that need CI/build verification because you could not run commands
```

---

## Prompt for Version 1.3 - Poll, log, and display split

```text
You are implementing GoDriveLog version 1.3 in GitHub repo MickMake/GoDriveLog.

Important constraints:
- You can access and edit files using the GitHub connector.
- Do not assume you can execute shell commands, run go build, run go test, or use a local checkout.
- Commit focused changes directly to main.
- Keep it SIMPLE.
- Do not implement virtual PID calculations.
- Do not implement engine start/stop detection.
- Do not implement daily log rotation yet unless the current logger requires a minimal compatibility tweak.

Read these files first using the GitHub connector:
- docs/config.md
- docs/implementation-phases-v1.md
- config.example.yaml
- cmd/GoDriveLog/main.go
- internal/logger/jsonl.go
- internal/sensors/reader.go
- internal/ui/dashboard.go
- internal/config/config.go

Task:
Implement the poll/log/display split phase.

Requirements:
1. Poll active runtime PIDs from phase 1.2.
2. If a PID has log: true, write the reading to JSONL.
3. If a PID has display.enabled: true, update the dashboard.
4. Allow log-only, display-only, log-and-display, and inactive known PIDs.
5. Update the reading/log record shape to include both:
   - key, e.g. rpm
   - raw pid, e.g. 010C
6. Keep unit from config as the unit of record unless this conflicts badly with the existing reader.
7. Do not add session metadata or database storage.

After editing, summarize:
- files changed
- final log record fields
- how display-only and log-only are handled
- known follow-up work for phase 1.4
- any areas that need CI/build verification because you could not run commands
```

---

## Prompt for Version 1.4 - Dashboard from display config

```text
You are implementing GoDriveLog version 1.4 in GitHub repo MickMake/GoDriveLog.

Important constraints:
- You can access and edit files using the GitHub connector.
- Do not assume you can execute shell commands, run go build, run go test, or use a local checkout.
- Commit focused changes directly to main.
- Keep it SIMPLE.
- Do not implement virtual PID calculations.
- Do not implement engine start/stop detection.
- Do not implement advanced dashboard layout features.

Read these files first using the GitHub connector:
- docs/config.md
- docs/implementation-phases-v1.md
- config.example.yaml
- internal/ui/dashboard.go
- cmd/GoDriveLog/main.go
- internal/config/config.go

Task:
Implement dashboard creation from display config.

Requirements:
1. Create dashboard widgets only for PIDs where display.enabled == true.
2. Use display.style for gauge, bar, or graph.
3. Use display.position x, y, width, and height for widget placement.
4. Use min, max, and unit from PID config for scaling/display.
5. Keep the current simple gauge/bar/graph behaviour. Do not make it pretty yet.
6. Keep per-sensor status display:
   - waiting
   - ok
   - error
   - stale
7. Ensure log-only PIDs do not create dashboard widgets.

After editing, summarize:
- files changed
- how widgets are selected
- how display position/style is mapped
- known follow-up work for phase 1.5
- any areas that need CI/build verification because you could not run commands
```

---

## Prompt for Version 1.5 - Daily log rotation

```text
You are implementing GoDriveLog version 1.5 in GitHub repo MickMake/GoDriveLog.

Important constraints:
- You can access and edit files using the GitHub connector.
- Do not assume you can execute shell commands, run go build, run go test, or use a local checkout.
- Commit focused changes directly to main.
- Keep it SIMPLE.
- Do not implement engine start/stop detection.
- Do not implement engine-state-based log rotation.
- Do not add compression, retention, or session metadata files.

Read these files first using the GitHub connector:
- docs/config.md
- docs/implementation-phases-v1.md
- config.example.yaml
- internal/logger/jsonl.go
- cmd/GoDriveLog/main.go
- internal/config/config.go

Task:
Implement simple daily log rotation.

Requirements:
1. Use config:
   log.rotate: daily
   log.directory: ./log
2. Open a date-based JSONL file inside log.directory.
3. Before each write, check whether the date changed.
4. If the date changed, close the current file and open a new one.
5. Use a simple filename. Prefer YYYY-MM-DD.jsonl unless that conflicts with current behaviour.
6. Remove or stop using engine-start-based rotation.
7. Keep logging synchronous and boring.

After editing, summarize:
- files changed
- log filename format
- rotation trigger
- what old rotation logic was removed or bypassed
- known follow-up work for phase 1.6
- any areas that need CI/build verification because you could not run commands
```

---

## Prompt for Version 1.6 - OBD PID mappings

```text
You are implementing GoDriveLog version 1.6 in GitHub repo MickMake/GoDriveLog.

Important constraints:
- You can access and edit files using the GitHub connector.
- Do not assume you can execute shell commands, run go build, run go test, or use a local checkout.
- Commit focused changes directly to main.
- Keep it SIMPLE.
- Do not implement generic PID decoding.
- Do not implement PID auto-discovery.
- Do not implement manufacturer-specific PIDs.

Read these files first using the GitHub connector:
- docs/config.md
- docs/implementation-phases-v1.md
- config.example.yaml
- internal/sensors/elmobd_reader.go
- internal/sensors/reader.go

Task:
Extend the elmobd reader to support the PIDs listed in config.example.yaml.

Requirements:
Support these mappings:
- engine_load: 0104, %, calculated engine load
- coolant_temp: 0105, C, engine coolant temperature
- rpm: 010C, rpm, engine RPM
- speed: 010D, km/h, vehicle speed
- intake_air_temp: 010F, C, intake air temperature
- throttle_position: 0111, %, throttle position
- fuel_level: 012F, %, fuel tank level
- control_module_voltage: 0142, V, control module voltage

Implementation style:
1. Keep the adapter as a small switch or map.
2. Return numeric values and unit strings.
3. Prefer config unit as the log/display unit where possible.
4. Do not add a generic OBD decoder.
5. Do not add virtual support.

After editing, summarize:
- files changed
- supported PID mappings
- any PID that could not be mapped cleanly to elmobd
- known follow-up work for phase 1.7
- any areas that need CI/build verification because you could not run commands
```

---

## Prompt for Version 1.7 - Documentation cleanup

```text
You are implementing GoDriveLog version 1.7 in GitHub repo MickMake/GoDriveLog.

Important constraints:
- You can access and edit files using the GitHub connector.
- Do not assume you can execute shell commands, run go build, run go test, or use a local checkout.
- Commit focused changes directly to main.
- Keep it SIMPLE.
- This phase is documentation cleanup only unless tiny code references must be updated for consistency.

Read these files first using the GitHub connector:
- README.md
- docs/config.md
- docs/implementation-phases-v1.md
- config.example.yaml
- any old JSON config files in the repo

Task:
Clean up docs after phases 1.1 through 1.6 are implemented.

Requirements:
1. Update README.md to reference:
   - docs/config.md
   - config.example.yaml
2. Update run command examples:
   ./GoDriveLog -config config.example.yaml
3. Remove old references to JSON config if they are obsolete.
4. If old JSON config files remain, either mark them as legacy or remove them if clearly safe.
5. Clearly document current limitations:
   - type: virtual is accepted but not implemented
   - engine start/stop detection is future release
   - daily log rotation only
   - no PID auto-discovery
6. Keep README concise.

After editing, summarize:
- files changed
- obsolete docs removed or updated
- remaining known limitations
- any areas that need CI/build verification because you could not run commands
```
