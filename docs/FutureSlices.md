# Future Slices

This file is a parking lot for approved or desired follow-on ideas that are not part of the current implementation slice.

Use this to capture "oh, also implement this later" notes without making the active slice ambiguous. Future prompts may reference this file, but items here are not current scope unless a later prompt explicitly promotes them.

## Indexed FutureSlices table

Effort is rough **Codex hours**, assuming the v3 dashboard/gauge code is already loaded in context and tests exist. Not human hours. Codex hours are weird little mushroom hours.

| # | State | Description | Area | Effort |
|---:|---|---|---|---:|
| 1 | Near / implementation-ready | Bar gauge overshoot follow-up | `gauge/bar`, `realism.overshoot`, animation | 2–4 |
| 2 | Near / needs spec tightening | Radial movement options | `gauge/radial`, movement policy, runtime animation | 3–5 |
| 3 | Later / visual feature | Radial needle trail | `gauge/radial`, renderer, animation history | 4–7 |
| 4 | Later / visual feature | Radial peak hold | `gauge/radial`, display marker, state tracking | 3–6 |
| 5 | Medium / useful soon | Value zones / warning-danger assets | `gauge/assets`, renderer, config validation | 4–7 |
| 6 | Medium / foundational logging | Canonical GoDriveLog Event Log | logging, sensor events, schema/versioning | 5–9 |
| 7 | Medium / pairs with event log | Session metadata sidecar | logging, replay metadata, config provenance | 4–7 |
| 8 | Medium / high value dev tool | JSONL dashboard replay | dashboard runtime, logs, replay CLI | 6–10 |
| 9 | Near / bounded utility | JSONL log validation | logs, CLI, schema validation | 3–5 |
| 10 | Later / architecture boundary | External converter boundary | `tools/converters`, import/export architecture | 3–6 |
| 11 | Near / performance polish | Needle Animation Performance | `gauge/radial`, animation loop, renderer, Pi4 performance | 3–6 |

A few notes from the table:

- **#1, #9, and #11** are the most directly sliceable.
- **#6–#8** are chunky but important because they form the logging/replay backbone.
- **#3 and #4** are fun, but they’re polish. Nice polish, but still polish.
- **#11 should not change the look.** It should protect the look from dropped/undersampled frames. That distinction is the whole ballgame.

## 0. Guidelines

- Keep entries small and slice-shaped.
- Mark ideas as `deferred`, `desired`, `exploratory`, or `rejected`.
- Do not treat this file as an implementation checklist.
- Do not let vague mentions here expand the current slice.
- Prefer a later dedicated prompt/spec before implementation.

## 1. Bar gauge overshoot follow-up

Status: deferred

`v3.5.10` is currently being treated as the radial overshoot slice because the active prompt/spec only defines radial overshoot behaviour. Bar gauge overshoot remains approved as a follow-up idea, but it must not be pulled into the radial overshoot implementation by inference from older `radial/bar overshoot` wording.

Bar gauges should eventually support `realism.overshoot`, but this was intentionally left out of the radial overshoot slice to avoid ambiguous behaviour and accidental scope creep.

Notes:

- Display-only.
- Bounded pass-and-settle movement.
- Should compose cleanly with bar damping/smoothing.
- Do not copy radial behaviour blindly; bar movement has its own visual semantics.
- A bar overshoot should affect the displayed fill/level extent, not mutate source sensor values.
- Clamp final settled display to the real target/range after the overshoot tail completes.
- Consider vertical and horizontal bars, plus different origins, when defining the later prompt.
- Keep radial overshoot behaviour unchanged when this is implemented.

Possible future slice:

```text
v3.5.x bar overshoot
```

## 2. Radial movement options

Status: desired

Radial gauges should eventually support the scalar `movement` options that already exist for gauge movement selection, while preserving current behaviour as the compatibility default.

Proposed movement meanings for radial gauges:

- `instant`: current radial behaviour; immediately render the needle at the target angle with no interpolation.
- `linear`: interpolate the displayed needle angle from the previous displayed angle to the target angle at constant progress.
- `bell`: interpolate with a slow start, faster middle, and slow end.

Rules:

- `instant` must preserve existing radial semantics.
- Movement must be display-only.
- Movement must animate displayed angle/position only; it must not mutate source values, logs, exported values, configured ranges, or input data.
- Do not pre-render or cache unbounded intermediate needle images.
- Prefer small per-gauge transition state such as previous angle, target angle, elapsed time, duration, movement mode, and active/inactive state.
- Keep needle geometry and image assets reusable; rotate or transform at render time rather than generating a frame cache.
- Do not combine this with damping, stiction, overshoot, peg bounce, needle trail, or peak hold unless a later slice explicitly defines composition.

Possible future slice:

```text
v3.5.x radial movement options
```

## 3. Radial needle trail

Status: desired

Add optional radial-only `realism.needle_trail` support.

Needle trail renders a bounded history of previous displayed needle positions as fading ghost needles. It is a visual afterimage effect, not a movement curve.

Proposed config shape:

```yaml
realism:
  needle_trail:
    length: 12
    decay_ms: 500
```

Options:

- `length`: maximum number of historical displayed needle positions retained. Default: `12`.
- `decay_ms`: time in milliseconds for retained trail samples to fade out. Default: `500`.

Rules:

- Radial-only.
- Disabled by default.
- Display-only.
- Must not mutate source values, logs, exported values, configured ranges, or input data.
- Store only a bounded history of displayed needle angles/positions and timestamps.
- Trail samples should fade and be discarded deterministically.
- Do not store an unbounded render history.
- Do not place this under `movement`; `movement` selects the travel curve, while `needle_trail` is a render-history effect.

Possible future slice:

```text
v3.5.19 radial needle trail
```

## 4. Radial peak hold

Status: desired

Add optional radial-only `realism.peak_hold` support.

Peak hold displays a secondary marker or needle at the highest displayed value reached. It is an instrument display feature, not a source value change.

Proposed config shape:

```yaml
realism:
  peak_hold:
    hold_ms: 0
    decay_ms: 1000
```

Options:

- `hold_ms`: how long to hold the peak after the displayed needle stops increasing. `0` means hold indefinitely.
- `decay_ms`: optional time for the peak marker to release/return after the hold expires.

Rules:

- Radial-only.
- Disabled by default.
- Display-only.
- Must not mutate source values, logs, exported values, configured ranges, or input data.
- Peak tracking should use displayed value/angle semantics defined by the later implementation prompt.
- If decay is enabled, release should be bounded and deterministic.
- Do not place this under `movement`; `movement` selects the travel curve, while `peak_hold` is a display marker/history feature.

Possible future slice:

```text
v3.5.20 radial peak hold
```

## 5. Value zones / warning-danger assets

Status: desired

Support optional value zones that select warning/danger variants of gauge assets when the source value reaches a configured range.

This should be a separate gauge-display feature, not part of `realism.overshoot`.

Proposed config shape:

```yaml
zones:
  warning:
    min: 6000
    max: 7000
  danger:
    min: 7000
    max: 8000
```

Proposed asset convention:

```text
needle.png
needle_warning.png
needle_danger.png
face.png
face_warning.png
face_danger.png
bar.png
bar_warning.png
bar_danger.png
```

Rules:

- Zone selection should follow the real/source target value, not any temporary animated display value.
- If a zone-specific asset exists for a layer, use it.
- If a zone-specific asset does not exist, fall back to the normal asset.
- Overshoot may visually pass a threshold, but should not change the zone state unless the real/source value is in that zone.
- Avoid surprising behaviour where a temporary animation makes the gauge appear to enter warning or danger falsely.

Possible future slice:

```text
v3.5.x value zones / warning-danger assets
```

## 6. Canonical GoDriveLog Event Log

Status: desired

Promote the current JSONL event logger output into a formal, versioned GoDriveLog-owned event log format.

The native working log format should be newline-delimited JSON with one event per line. This is the format GoDriveLog core writes, validates, and replays. Other formats should be converted into this format rather than being supported directly inside the runtime.

Proposed file naming:

```text
*.gdl.jsonl
*.gdl.meta.json
```

Proposed event schema marker:

```json
{"schema":"godrivelog.event.v1"}
```

Rules:

- Treat GoDriveLog JSONL as the canonical event log, not as incidental logger output.
- Add a schema marker or schema version to every event record.
- Keep one complete event per line.
- Preserve the existing event-oriented shape: kind, sensor id, timestamps, status, typed value, previous status, and error.
- JSONL events should represent sensor events, not rendered dashboard state.
- Do not use CSV, MDF4, BLF, Parquet, ROS bag, or any other external format as the native runtime log format.
- Industry or third-party formats may be supported by converters, importers, or exporters outside the core runtime.
- Keep logs inspectable, appendable, streamable, and replayable.
- Maintain backwards compatibility or provide a clear migration path if the existing JSONL shape changes.

Possible future slice:

```text
v3.x canonical event log v1
```

## 7. Session metadata sidecar

Status: desired

Add a session metadata sidecar next to each GoDriveLog event log.

The sidecar should capture enough context to replay, validate, or interpret a log later, even if the active dashboard config has changed.

Proposed shape:

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

Rules:

- Keep high-volume sensor events in the `.gdl.jsonl` file.
- Keep session-level context in `.gdl.meta.json`.
- Do not duplicate full session metadata onto every event line.
- Capture enough sensor mapping information to support replay and conversion audits.
- Treat the sidecar as optional for reading older logs but preferred for new logs.

Possible future slice:

```text
v3.x session metadata sidecar
```

## 8. JSONL dashboard replay

Status: desired

Add a replay mode that consumes GoDriveLog event logs and feeds recorded events back into the dashboard runtime.

This is a core development and validation feature. It allows a real OBD session to be captured once and replayed repeatedly without the vehicle attached.

Proposed command shape:

```text
godrivelog dashboard replay --config dashboard.yaml --log drive.gdl.jsonl
```

Replay path:

```text
.gdl.jsonl -> SensorEvent stream -> dashboard runtime -> renderer
```

Rules:

- Replay recorded sensor events directly; do not pretend JSONL is an OBD adapter.
- Do not mutate the live OBD polling path.
- Preserve `event_at` and `read_at` semantics when rebuilding events.
- Replay should feed the same dashboard boundary used by live runtime rendering.
- Replay should not write new source sensor values unless explicitly configured to log replay output.
- Replay should be deterministic for the same input log, dashboard config, and replay options.
- Preview mode remains separate: preview is manual one-gauge testing; replay is recorded event-stream playback.

Useful options:

```text
--speed 1.0      # original timing
--speed 2.0      # double speed
--speed 0        # no sleeps / fastest possible
--from <time>    # optional later slice
--to <time>      # optional later slice
--loop           # optional later slice
```

Possible future slice:

```text
v3.x JSONL dashboard replay
```

## 9. JSONL log validation

Status: desired

Add a validator for GoDriveLog event logs before replay or conversion.

Proposed command shape:

```text
godrivelog logs validate drive.gdl.jsonl
```

Rules:

- Validate that every line is valid JSON.
- Validate known schema markers.
- Validate required fields.
- Validate timestamps are parseable.
- Validate typed value objects.
- Validate status/error semantics.
- Warn, rather than fail, on non-monotonic timestamps unless a later spec requires strict ordering.
- Produce useful line-numbered errors for converter/debugging work.

Possible future slice:

```text
v3.x GoDriveLog log validator
```

## 10. External converter boundary

Status: desired

Keep foreign-format conversion outside GoDriveLog core runtime.

Converters should live under `tools/converters` and convert external telemetry/log formats into canonical GoDriveLog event logs.

Proposed layout:

```text
tools/
  converters/
    README.md
    csv-to-gdl-jsonl/
    racechrono-to-gdl-jsonl/
    decoded-can-csv-to-gdl-jsonl/
```

Rules:

- GoDriveLog core should understand GoDriveLog event logs, not every external telemetry format.
- Foreign formats convert into `.gdl.jsonl` plus optional `.gdl.meta.json`.
- Converters may understand CSV, RaceChrono, Torque Pro, decoded CAN CSV, racing datasets, or other third-party formats.
- Converter-specific mapping files are allowed and encouraged.
- Do not add converter dependencies to the dashboard runtime.
- Do not let a one-off converter become a production runtime dependency.
- Import mapping should be explicit enough to preserve sensor ids, units, timestamps, and source provenance.

Possible future slice:

```text
v3.x tools/converters boundary
```

## 11. Needle Animation Performance

Status: desired

Improve the reliability of short radial needle animations on lower-powered display targets such as Raspberry Pi 4.

The current radial peg-bounce and overshoot effects can look visually correct when rendered smoothly, but very short animation tails may occasionally appear to skip or vanish under load. This is most noticeable for subtle effects where the important motion happens over only a few rendered frames.

This is not a request to make the needle movement stronger, larger, or more dramatic. The visual result should remain essentially unchanged. The goal is to make existing subtle radial needle movement render more consistently.

Areas to investigate:

- animation tick cadence;
- render/update scheduling;
- minimum visible animation sample count;
- short movement tail duration;
- Raspberry Pi 4 GPU/CPU load;
- asset size and rotation cost;
- whether radial needle animation needs fixed-timestep sampling independent of frame delivery.

Rules:

- Preserve existing visual semantics for radial overshoot, peg bounce, damping, and stiction.
- Do not increase rebound or overshoot amplitude as a performance fix.
- Do not change source sensor values, logs, exported values, or configured ranges.
- Prefer timing/sampling/render-loop improvements over changing gauge behaviour.
- Keep this radial-focused unless profiling proves another gauge type is affected.
- Any fix should be validated on Raspberry Pi 4 or equivalent low-power target, not only on desktop.
- Treat occasional dropped frames as a rendering/performance problem, not as a realism tuning problem.

Possible future slice:

```text
v3.x radial needle animation performance


