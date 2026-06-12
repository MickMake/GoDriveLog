# GoDriveLog v3 working-code inventory and seam plan

Status: implementation slice output  
Target version: `v3.0.0`  
Branch: `v3.0.0-working-code-inventory`

## 1. Purpose

This document records the first v3.0.0 implementation slice: inventory the current working code, map it to v3 roles, and choose seams before replacing or refactoring anything.

This is intentionally docs-only. Runtime code comes later. The point is to avoid turning useful current behaviour into accidental v3 architecture just because it already exists and looks at us with hopeful little eyes.

## 2. v3 target model

The target runtime model remains:

```text
selected vehicle
-> OBD endpoint
-> sensor polling runtime
-> sensor events
-> selected logs and dashboards as subscribers
```

The selected vehicle is the runtime profile. Sensors and assets are global catalogues. Logs and dashboards are global subscriber/display definitions selected by the vehicle.

## 3. Files inspected

Planning and guardrail docs:

- `docs/v3/MigrationState.md`
- `docs/v3/ChatPrompts.md`
- `docs/v3/MigrationGuardrails.md`
- `docs/v3/ImplementationGuardrails.md`
- `docs/v3/DirectoryStructure.md`

Current implementation files:

- `cmd/GoDriveLog/main.go`
- `internal/config/config.go`
- `internal/config/runtime.go`
- `internal/config/dashboard.go`
- `internal/sensors/reader.go`
- `internal/sensors/elmobd_reader.go`
- `internal/sensors/state.go`
- `internal/sensors/state_store.go`
- `internal/logger/jsonl.go`
- `internal/ui/dashboard.go`
- `internal/dashboard/assets/registry.go`
- `internal/dashboard/decoders/registry.go`
- `internal/dashboard/scene/scene.go`
- `internal/dashboard/renderer/fyne/renderer.go`

Tests inspected:

- `internal/config/config_test.go`
- `internal/sensors/state_store_test.go`
- discovered dashboard asset, decoder, scene, renderer, and UI tests through repository search

## 4. Inventory summary

| Area | Current shape | v3 role | Decision |
|---|---|---|---|
| Config root | `obd`, `log`, `vehicle`, `sensors`, `dashboard` | strict v3 root: `vehicles`, `sensors`, `assets`, `logs`, `dashboards` | replace shape, reuse strict loading habit |
| Runtime startup | `main.go` wires config, logger, reader, polling goroutines, Fyne window, dashboard | selected vehicle -> endpoint -> sensor runtime -> subscribers | replace orchestration, salvage behaviour |
| OBD/ELM327 | `sensors.Reader`, `MockReader`, `ELMOBDReader` | endpoint/reader seam behind selected vehicle endpoint | wrap/refactor |
| Sensor polling | polling loops live in `main.go` | shared sensor polling runtime | refactor |
| Sensor state | `SensorState`, `StateStore`, stale snapshots | latest-state store plus event source | reuse/refactor |
| Logging | daily JSONL writer consumes `sensors.Reading` | log subscriber consuming sensor events | wrap/refactor |
| Dashboard config | assets/decoders/blocks/layers mini-scene model | dashboard widgets selected by vehicle | archive/replace shape, salvage rendering ideas |
| Dashboard runtime | UI refresh loop polls state store on cadence | dashboard subscriber consuming latest state/events | refactor |
| Assets | dashboard-local asset list and `asset_root` | global repo-root-relative asset catalogue | replace shape, reuse loader pieces |
| Renderer | scene -> Fyne cached objects | v3 widget renderer | reuse/refactor locally |
| Tests | good validation/state/renderer behaviour tests but tied to old shapes | boundary tests for v3 config/runtime/assets/widgets | keep/rewrite selectively |

## 5. Config loading and structs

Current code:

- `internal/config/config.go` defines a root `Config` with `OBD`, `Log`, `Vehicle`, `Sensors`, and one `Dashboard`.
- Loading uses `yaml.Decoder.KnownFields(true)`, which is worth preserving as a strict-loading habit.
- Current defaults include top-level OBD, log directory, dashboard refresh, and dashboard render minimums.
- Sensor config uses `refresh` and per-sensor `log`; v3 wants `poll`, with logs selecting sensors separately.
- Current validation already rejects some legacy shapes, which is good, but it validates the wrong target shape.

Decision: **replace config shape, reuse strict-loading approach**.

v3 role:

- `Config` should become the documented v3 root shape: `vehicles`, `sensors`, `assets`, `logs`, `dashboards`.
- Vehicles select logs and dashboards by ID.
- Sensors and assets remain global.

Seam/boundary:

- Add v3 config structs and loader behind a clear v3 config boundary before wiring runtime.
- Keep any old config support out of the v3 loader.
- Prefer a clean v3 package/file split matching `docs/v3/DirectoryStructure.md`.

Must not leak into v3:

- top-level `obd`
- top-level singular `log`
- singular `vehicle`
- singular `dashboard`
- `dashboard.refresh_ms` or `dashboard.render_min_ms`
- dashboard-local `asset_root`
- sensor-owned `log`
- `refresh` as the v3 sensor cadence name
- asset/decoder/block/layer config as the v3 dashboard schema

Tests to preserve or rewrite:

- Preserve the strict unknown-field intent.
- Rewrite validation tests around v3 root and nested unknown fields.
- Keep useful reference-validation patterns, but point them at v3 vehicle/log/dashboard/sensor/asset references.

## 6. Runtime startup and command flow

Current code:

- `cmd/GoDriveLog/main.go` loads config, derives active sensors, creates a state store, opens JSONL logging, chooses mock or ELMOBD reader, creates the Fyne app/window, starts dashboard refresh, and launches one goroutine per active sensor.
- It currently decides active sensors from `sensor.log == true`, so logging selection drives polling.
- Dashboard and logger are wired directly in `main.go`.

Decision: **replace orchestration, salvage behaviour**.

v3 role:

- `main.go` should become a thin command entry point.
- Runtime resolution should produce a plan from loaded config plus selected vehicle.
- Sensor polling should be owned by a sensor runtime, not by `main.go`.
- Logs and dashboards should be subscribers selected by the vehicle.

Seam/boundary:

- Introduce a RuntimePlan-style boundary after strict v3 config loading.
- Keep endpoint selection, sensor polling, logging, and dashboard startup behind runtime package seams.

Must not leak into v3:

- logger deciding active sensors
- dashboard refresh settings as config cadence
- direct Fyne setup as the core runtime shape
- mock mode as an OBD config boolean that spreads through runtime
- one-off goroutine orchestration inside `main.go`

Tests to preserve or rewrite:

- Add later tests proving selected vehicle controls endpoint, logs, and dashboards.
- Add tests proving sensors are polled independently of log/dashboard consumers.

## 7. OBD, ELM327, vehicle, endpoint, and adapter code

Current code:

- `internal/sensors/reader.go` defines a simple `Reader` interface: `Read(ctx, pid) (float64, string, error)`.
- `MockReader` provides simulated values for common PIDs.
- `internal/sensors/elmobd_reader.go` wraps `github.com/rzetterberg/elmobd` and maps known PIDs to typed elmobd commands.
- Unsupported PIDs return an error.
- Current endpoint selection is `if cfg.OBD.MockMode { mock } else { ELMOBDReader(address, debug) }`.

Decision: **wrap/refactor**.

v3 role:

- Existing ELMOBD command knowledge is valuable behind the v3 endpoint/reader seam.
- The mock reader is useful as an early simulator, but the v3 docs prefer endpoint addresses such as `serial://...` and `tcp://...` over endpoint-type branching.

Seam/boundary:

- Create an endpoint connector abstraction later that can support serial and TCP simulator endpoints.
- Wrap existing ELMOBD code behind that seam where practical.
- Keep PID command mapping isolated from config ownership.

Must not leak into v3:

- `mock_mode` as a core runtime branch
- direct dependency on elmobd types outside the adapter boundary
- returning fake numeric values for unsupported/error states as if they were valid readings
- endpoint-specific concepts inside sensors, logs, dashboards, or assets

Tests to preserve or rewrite:

- Add adapter tests around supported and unsupported PIDs where mocking is practical.
- Add endpoint address validation tests in v3 config validation.

## 8. Sensor polling, state, cache, and status logic

Current code:

- Polling loops are in `main.go`, one goroutine per active sensor.
- `internal/sensors/state.go` defines `SensorState` with value, unit, range, status, error, updated timestamp, and stale threshold.
- `internal/sensors/state_store.go` provides a thread-safe latest-state store.
- Stale status is computed from `StaleAfter` during snapshots.
- Current stale threshold derives from `refresh * 2` in `internal/config/runtime.go`.

Decision: **reuse/refactor**.

v3 role:

- The latest-state store is close to the v3 need and should be salvaged.
- v3 needs sensor events with value, unit, status, original read timestamp, and sequence/version.
- v3 stale rule should be `max(sensor.poll * 3, 1000ms)`, not current `refresh * 2`.

Seam/boundary:

- Split current state-store behaviour from current config/runtime assumptions.
- Add event emission later around state transitions: first reading, value change, status change, recovery, stale, error, unsupported.

Must not leak into v3:

- `refresh` naming
- `log` flag determining whether a sensor is active
- stale threshold from old config/runtime helper
- float-only values if boolean/status sensors are needed
- treating `0` as an error, unsupported, or missing value

Tests to preserve or rewrite:

- Preserve StateStore initialization, value update, error preservation, sorted snapshots, and stale handling tests.
- Add event tests later for first reading, unchanged value suppression, value change, status change, stale transition, and recovery.

## 9. Logging and JSONL writer behaviour

Current code:

- `internal/logger/jsonl.go` opens a daily JSONL file under a directory.
- `Write` serializes `sensors.Reading` as one JSON object per line.
- The writer rotates by current wall-clock date.
- No dedicated logger tests were found.

Decision: **wrap/refactor**.

v3 role:

- Keep daily JSONL writing behaviour if still wanted.
- Change input from `sensors.Reading` to v3 sensor events.
- Logs should be global subscriber definitions selected by vehicle.
- Logs should write first readings, value changes, and status changes, without spamming unchanged duplicate readings.

Seam/boundary:

- Wrap the existing writer behind a v3 log subscriber interface.
- Keep file writing separate from event filtering.

Must not leak into v3:

- logger as hidden sensor scheduler
- `sensor.log` as polling selector
- writer timestamp replacing sensor read timestamp
- writing every poll tick by default

Tests to preserve or add:

- Add JSONL tests for file creation, one-object-per-line encoding, rotation behaviour if retained, and close idempotence.
- Add subscriber tests for first reading/value-change/status-change filtering and read timestamp preservation.

## 10. Dashboard, renderer, Fyne, display, and widget code

Current code:

- `internal/ui/dashboard.go` builds a dashboard from current dashboard config, asset registry, decoder execution, scene evaluation, and Fyne renderer.
- It refreshes on a ticker and pulls a stale-aware snapshot from the state store.
- `internal/config/dashboard.go` defines the old dashboard schema: canvas, asset root, assets, decoders, blocks, layers, conditions, geometry.
- `internal/dashboard/decoders` implements a decoder pipeline.
- `internal/dashboard/scene/scene.go` resolves old blocks into primitive scene elements.
- `internal/dashboard/renderer/fyne/renderer.go` renders scene elements using cached Fyne objects and resources.

Decision: **archive/replace schema, reuse/refactor renderer ideas**.

v3 role:

- v3 dashboards are global definitions selected by vehicles.
- v3 widgets reference global sensors and global assets.
- v3 starts with `image`, `digit_display`, and `indicator`, later adding `bar_display` and `frame_gauge`.
- Rendering should consume sensor state/events, not poll OBD directly.

Seam/boundary:

- Treat the old decoder/block/layer model as an archive/current-dashboard implementation, not as v3 schema.
- Reuse Fyne caching/resource-update ideas where they fit native-size image-backed v3 widgets.
- Keep any compatibility layer temporary and outside the v3 core model.

Must not leak into v3:

- dashboard-local `assets`, `decoders`, `blocks`, and `layers`
- expression-like conditions as v3 widget logic
- dashboard refresh/cadence fields
- Fyne renderer concepts defining config schema
- geometry-first `rect` model when v3 widgets use `position` and asset-native dimensions

Tests to preserve or rewrite:

- Preserve renderer caching/performance behaviour tests where practical.
- Preserve scene/condition tests only as current-dashboard archive tests unless they map cleanly to v3 widgets.
- Add v3 widget tests later for image, digit display, indicator, bar display, and frame gauge semantics.

## 11. Asset loading and image handling

Current code:

- `internal/dashboard/assets/registry.go` loads image, frame-set, and charset assets from dashboard-local config.
- It resolves paths relative to the config file plus optional `asset_root`.
- It rejects remote paths.
- It supports generated frame paths using `{index}` or zero-padded `{index:0N}` markers.

Decision: **replace shape, reuse loader pieces**.

v3 role:

- Assets become global catalogues: `digit_sets`, `bar_sets`, `frame_sets`, `indicator_sets`, and `image_sets`.
- Asset paths should be repository-root relative.
- Widgets map sensor state to asset family behaviour.

Seam/boundary:

- Build a v3 asset registry with global families.
- Reuse path reading, frame expansion, and useful error wording where it still fits.

Must not leak into v3:

- `asset_root`
- config-file-relative active examples
- `charset` as the v3 digit model
- dashboard-owned assets
- remote path support
- asset rules/conditions/scripts

Tests to preserve or rewrite:

- Preserve missing-file and bad-frame-pattern behaviours where useful.
- Rewrite tests around repository-root-relative asset paths and v3 asset families.
- Add tests for digit set required characters, indicator `off/on/unknown`, bar set `off/on`, frame range validation, and useful dimension mismatch errors.

## 12. Test inventory

Useful current tests:

- Config tests prove strict rejection of several obsolete config shapes and validate dashboard references.
- State store tests prove latest-state behaviour, metadata preservation, sorted snapshots, stale derivation, and per-sensor stale thresholds.
- Dashboard asset, decoder, scene, renderer, and UI tests likely preserve useful display behaviour.

Decision: **keep tests that prove behaviour, rewrite tests that protect obsolete shape**.

Preserve:

- strict unknown-field attitude
- reference validation style
- state-store behaviour
- stale/error visibility principles
- renderer cache/update behaviour
- asset loading error clarity

Rewrite/archive:

- tests that assert old root config shape
- tests that assert dashboard decoder/block/layer config as active target
- tests that assert dashboard refresh cadence in config
- tests that treat sensor `log` as active runtime selection

## 13. First implementation seam

The first code seam after this document should be:

```text
strict v3 config loading and validation, isolated from current runtime wiring
```

That matches the queued `v3.0.2` implementation direction, but should only start after `v3.0.1` freezes/tightens the schema docs.

Recommended near-term order:

1. Finish this v3.0.0 inventory/seam-plan PR.
2. Verify and merge it.
3. Run v3.0.1 to freeze/tighten docs and schema examples.
4. Start v3.0.2 strict v3 config loading and validation.

## 14. Risks and trade-offs

- The current code has useful behaviour, but much of it is attached to the old dashboard/config shape.
- Reusing too much too early risks dragging old assumptions into v3.
- Rewriting everything would also be wasteful; the sensor state store, ELMOBD PID knowledge, JSONL writer, asset loader pieces, and Fyne renderer caching are all useful.
- The clean path is boring: strict config first, runtime plan second, endpoint/sensor/log/dashboard seams after that.

## 15. Explicit non-goals for this slice

This slice does not:

- implement runtime code
- change the v3 config schema
- add config aliases
- convert old config to v3 config
- delete current code
- start v3.0.1 or v3.0.2 work
- make the current dashboard model the v3 dashboard model

## 16. Summary decision

Salvage behaviour and proven implementation knowledge. Do not preserve the old machine.

The current code should be mined carefully, not crowned king. Crowning old code is how architecture gets a sceptre, a tax policy, and a dungeon full of TODOs.
