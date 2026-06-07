# Dashboard Implementation Prompts - v2.x.x series

## How to use this file

Each prompt is intended to start a separate implementation chat or branch.

Before using any implementation prompt, follow the project coding rules:

1. Check latest `main`.
2. Check whether any previous branches or PRs are unmerged.
3. If unmerged work exists, stop and ask what to do.
4. Branch from latest `main`.
5. Keep changes small.
6. Do not support legacy dashboard config unless explicitly requested.
7. Do not touch directories listed in the prompt guardrails.
8. Update docs/tests for the stage being implemented.

The prompts are intentionally condensed. Do not invite the architecture goblin to brunch.

---

## Global guardrails for all v2 dashboard work

```text
Project: MickMake/GoDriveLog

Hard rule:
Do not implement compatibility with the old per-PID display.widget dashboard model.
Read and obey docs/dashboard/v2/repo-structure-guardrails.md

Architecture rule:
Sensors produce state. The dashboard scene consumes state. The dashboard scene must not own polling, OBD reading, or logging.

Do not touch unless required:
- internal/sensors real OBD reader behaviour
- internal/logger JSONL behaviour
- mock reader behaviour, except where test support requires small changes
- unrelated CLI behaviour
- asset image files unless the stage is specifically about assets

Keep:
- app starts from config
- mock mode works
- logging works
- tests pass

Prefer:
- small packages
- table-driven tests
- clear validation errors
- boring code
```

---

## v2.0.x prompt - New config schema

```text
Implement stage v2.0.x of the GoDriveLog dashboard rewrite.

Goal:
Separate sensor configuration from dashboard visual configuration.

Required outcome:
- Add a top-level sensors section.
- Add a top-level dashboard section with at least canvas width/height.
- Remove dashboard/display ownership from sensor/PID config.
- Keep OBD polling, mock mode, and JSONL logging working.
- Active sensor extraction should use sensors, not vehicle.pids.*.display.
- Add/update tests for config loading and active sensor extraction.
- Update config.example.yaml to the new v2.0.x shape.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not implement rendering.
- Do not implement assets, decoders, blocks, or layers beyond placeholder structs if required.
- Do not support old vehicle.pids.*.display.
- Do not touch internal/sensors OBD decoding unless absolutely necessary.
```

---

## v2.1.x prompt - Dashboard config validation only

```text
Implement stage v2.1.x of the GoDriveLog dashboard rewrite.

Goal:
Load and validate the dashboard schema, but do not render it yet.

Required schema areas:
- dashboard.canvas
- assets
- decoders
- blocks
- layers

Required outcome:
- Config structs exist for dashboard scene config.
- Validation catches duplicate IDs, missing IDs, invalid types, missing references, invalid canvas size, and invalid geometry.
- Good config loads.
- Bad config returns clear errors.
- Tests cover valid and invalid config examples.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not load images yet.
- Do not implement Fyne rendering.
- Do not implement decoder execution.
- Do not support old display.widget config.
- Do not touch polling/logging behaviour.
```

---

## v2.2.x prompt - Sensor state boundary

```text
Implement stage v2.2.x of the GoDriveLog dashboard rewrite.

Goal:
Introduce a neutral runtime sensor state boundary.

Required outcome:
- Add a SensorState model containing ID, value, unit, min, max, updated time, and status/error.
- Add a StateStore or equivalent latest-value store.
- Polling writes readings and errors into the state store.
- Logging continues to work.
- Dashboard-facing code reads from state, not PID config.
- Add tests for state updates, error/stale handling, and retrieval.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not implement decoders.
- Do not implement asset loading.
- Do not render the new dashboard yet.
- Do not make dashboard code call the OBD reader.
- Do not change real OBD PID decoding unless absolutely necessary.
```

---

## v2.3.x prompt - Decoder engine

```text
Implement stage v2.3.x of the GoDriveLog dashboard rewrite.

Goal:
Add a reusable decoder engine that converts sensor state into useful visual-ready outputs.

Required decoder types:
- normalize
- threshold
- frame_index
- format_number
- digits
- boolean

Required outcome:
- Decoders are configured by ID.
- Decoders can reference sensor values.
- Decoders can reference previous decoder outputs if that is simple and safe.
- Decoder execution is independent of Fyne.
- Add table-driven tests for each decoder type.
- Errors are clear for unknown inputs, invalid thresholds, bad frame counts, and invalid formats.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not implement arbitrary eval.
- Do not implement a scripting language.
- Do not render images.
- Do not bind decoders to specific sensors like rpm in code.
- Do not touch old widget packages.
```

---

## v2.4.x prompt - Asset registry

```text
Implement stage v2.4.x of the GoDriveLog dashboard rewrite.

Goal:
Add an asset registry for dashboard visuals.

Required asset types:
- image
- frame_set
- charset

Required outcome:
- Assets are configured by ID.
- Asset paths resolve relative to the config file or configured asset root.
- Missing image files fail validation/load with clear errors.
- Frame sets validate frame count and generated paths.
- Charsets validate character-to-image mappings.
- Assets are cached after load.
- Add tests using small fixture assets.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not implement scene rendering yet.
- Do not silently substitute missing images.
- Do not load assets every frame.
- Do not add remote asset fetching.
- Do not touch OBD/logger code.
```

---

## v2.5.x prompt - Scene primitives

```text
Implement stage v2.5.x of the GoDriveLog dashboard rewrite.

Goal:
Add scene primitive models and runtime evaluation without depending on old widgets.

Required primitives:
- image
- sprite_frame
- sprite_text
- group
- condition
- z-order

Required outcome:
- Scene elements can be sorted by z-order.
- Conditions can show/hide elements based on sensor or decoder values.
- sprite_frame can resolve a frame index to a frame_set asset.
- sprite_text can resolve formatted text to charset glyphs.
- group can contain child elements.
- Add tests for z-order, condition evaluation, frame resolution, and text/glyph resolution.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not implement full Fyne renderer if that becomes large; keep this stage model/evaluation focused.
- Do not add gauges.
- Do not hardcode RPM, speed, throttle, or coolant.
- Do not support old display.widget config.
```

---

## v2.6.x prompt - Fyne scene renderer

```text
Implement stage v2.6.x of the GoDriveLog dashboard rewrite.

Goal:
Render configured scene primitives in Fyne.

Required outcome:
- Add a Fyne scene renderer that draws layers/elements in z-order.
- Render image assets.
- Render sprite_frame elements.
- Render sprite_text elements from charsets.
- Apply visibility conditions on update.
- Scene updates when StateStore values change.
- The renderer does not perform decoder logic directly; it consumes resolved scene/decoder state.
- Add tests where practical and keep manual mock-mode verification simple.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not keep old panel renderer wired in.
- Do not make Fyne renderer responsible for polling or logging.
- Do not implement visual editor.
- Do not add remote config sync.
- Do not implement old widget compatibility.
```

---

## v2.7.x prompt - First real dashboard

```text
Implement stage v2.7.x of the GoDriveLog dashboard rewrite.

Goal:
Deliver the first real asset-driven dashboard as a vertical slice.

Required dashboard:
- static background
- RPM seven-segment sprite text
- throttle sprite-frame bar
- redline glow overlay
- status/error/stale indicator

Required outcome:
- Add a complete example dashboard config.
- Add minimal fixture/example assets or documented placeholders.
- Mock mode demonstrates the dashboard changing.
- Redline glow appears when threshold condition is met.
- README explains how to run the example.
- Tests remain green.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not try to implement every dashboard style.
- Do not add editor tooling.
- Do not add Google Drive sync.
- Do not revive old display.widget.
```

---

## v2.8.x prompt - Remove old widgets

```text
Implement stage v2.8.x of the GoDriveLog dashboard rewrite.

Goal:
Delete the old dashboard widget architecture.

Remove or retire:
- vehicle.pids.*.display
- DisplayConfig.Widget
- validDisplayWidget()
- old panel dashboard
- hardcoded widget factory if unused
- dead radial/bar/speedhud/ramped packages if no longer referenced

Required outcome:
- Build has no dependency on old display.widget.
- Config examples use only the new scene dashboard.
- README describes v2 dashboard config.
- Tests pass.
- Dead code is removed, not left as a decorative fossil.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not keep compatibility shims.
- Do not leave old widget names in docs.
- Do not delete useful non-dashboard code.
- Do not touch OBD/logger behaviour unless a compile break requires a minimal fix.
```

---

## v2.9.x prompt - Reusable block library

```text
Implement stage v2.9.x of the GoDriveLog dashboard rewrite.

Goal:
Add reusable visual blocks so dashboards can be built from config and images.

Initial blocks:
- seven_segment_number
- percent_frame_bar
- state_lamp
- glowing_number_box
- labelled_sensor_value
- warning_overlay
- stale_overlay
- static_panel

Required outcome:
- Blocks have named inputs.
- Blocks are sensor-agnostic.
- Blocks expand into scene primitives.
- Block input validation catches missing/wrong inputs.
- Example dashboard uses blocks instead of repeated raw primitives.
- Documentation lists block names, inputs, and examples.

Guardrails:
- Read and obey docs/dashboard/v2/repo-structure-guardrails.md
- Do not make blocks hardcode rpm/speed/throttle unless they are example instances.
- Do not add a plugin system.
- Do not add a scripting language.
- Do not add visual editor tooling.
```
