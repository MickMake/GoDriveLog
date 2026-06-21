# CHANGES

## Unreleased

- Updated the v3 Fyne display scene sink path to use non-blocking latest-only submissions for live dashboard and harness updates, preventing display rendering from throttling sensor/harness event cadence while preserving latest-frame coalescing and render error visibility.
- Updated `v3.2.6` Fyne radial rendering to prepare 1-degree radial needle frame sets outside normal live update sweeps, keeping live updates to keyed image resource swaps and preserving keyed canvas object reuse.
- Added `v3.2.6` Fyne radial gauge rendering, including ordered radial layer rendering, image-space needle rotation, normalised pivot placement, rotated-needle resource caching, and adapter coverage.
- Added `v3.2.5` radial gauge scene model support, including dashboard runtime routing, package-owned pivots, value-map angle calculation, needle scene part data, non-ok needle suppression, and radial scene signatures.
- Added `v3.2.4` Fyne seven-segment rendering hardening: stable keyed canvas object reuse, glass overlay ordering coverage, and a deterministic adapter benchmark for repeated digit updates.
- Added `v3.2.3` seven-segment gauge support through the dashboard scene path, including `type: gauge` package loading, package-owned sensor state, static layers, digit positions, package-owned formatting, non-ok suppression, Fyne adapter positioning, and scene signatures.
- Removed the redundant post-`Validate` gauge widget ownership pass now that ownership validation lives inside `Validate`.
- Added `v3.2.2` dashboard config support for `type: gauge` widgets with package-owned gauge paths, placement, and scale.
- Added validation tests proving gauge widgets do not own sensors and reject widget-level `sensor` fields.
- Added `v3.2.1` gauge package loader support for self-contained packages under `assets/gauges/**/gauge.yaml`, including `seven_segment` and `radial` package parsing.
- Added gauge package loader tests for valid seven-segment and radial packages, arbitrary package directory names, shared relative image paths, missing `gauge.yaml`, unsupported types, and asset-tree traversal rejection.
- Added `v3.2.0` planning baseline docs for self-contained gauge packages under `docs/v3.2/`, including `ImplementationState.md`, `OpenDecisions.md`, `ReleasePlan.md`, `CarryForward.md`, and per-slice prompts.
- Documented the v3.2 gauge package architecture: dashboard `type: gauge` widgets place gauge packages, gauge packages own sensor binding and radial gauge definition, and `assets/gauges/**/gauge.yaml` is the only required gauge package filename.
- Normalised active v3 root config, examples, and simple dashboard assets.
- Archived v3.0 baseline docs under `docs/archive/v3.0/`.
- Marked v3.1.7 dashboard event efficiency and v3.1.8 retirement readiness as deferred, not cancelled, until the gauge package direction is established.
- Added `v3.1.6` explicit sensor status semantics for `missing`, `unsupported`, `timeout`, `parse_error`, generic `error`, `unknown`, `ok`, and `stale`.
- Added `v3.1.5` typed sensor values for v3 sensor state/events, making `kind` mandatory and rejecting empty or unknown value kinds.
- Added parser/config value-kind contract checks so configured sensor `value_kind` must match the selected parser output kind; OBD sensors currently derive `numeric` from the parser contract when omitted.
- Updated JSONL event records to write a typed `value` object instead of a bare numeric value, including explicit error values for bad or unavailable readings.
- Added `v3.1.4` daily JSONL rotation for v3 event logs, deriving concrete files such as `logs/vw_caddy-2026-06-18.jsonl` from configured base paths such as `logs/vw_caddy.jsonl`.
- Added tests for daily JSONL path generation and writer rollover across a date boundary.
- Documented `v3.1.4` scope: daily rotation only, with no configurable rotation modes, retention policy, compression, upload, or logging architecture rewrite.
- Documented `v3.1.3` dashboard performance design decisions, including Raspberry Pi 4 2GB memory-churn rationale, Fyne image reuse, scene coalescing, render error propagation, config-derived startup window sizing, and shutdown constraints.
- Added `v3.1.3` dashboard scene update coalescing so v3 display rendering keeps only the latest pending scene instead of queueing stale frames.
- Wired both `--v3` and `--v3 --harness` display paths through the coalescing scene sink, preserving sensor polling and logging priority over dashboard freshness.
- Added tests proving scene sink submission returns while rendering is busy and stale pending frames are replaced by the latest scene.
- Added `v3.1.2` dashboard harness behind `--v3 --harness`, feeding fake sensor events through the real v3 dashboard scene path and Fyne display adapter.
- Removed the temporary `--repo-root` flag and added ordered relative asset search paths: config directory + vehicle ID, current working directory + vehicle ID, config directory, then current working directory.
- Added explicit harness patterns: `sweep`, `heartbeat`, and `fixed`, including tuned sweep and heartbeat timing tests.
- Documented dashboard harness command usage, asset search path order, pattern semantics, and 50ms/100ms cadence options.
- Added `v3.1.1` Fyne display adapter for v3 dashboard scene output, keeping display code below the dashboard runtime boundary.
- Wired the `--v3` command path to show selected v3 dashboard scenes in a Fyne window while retaining the existing runtime as the default path.
- Added adapter tests covering repo-relative asset rendering and rejection of escaping asset paths.
- Added `v3.1.0` runnable v3 command path behind `--v3`, including selected vehicle resolution, endpoint connection, sensor polling runtime startup, selected JSONL event subscribers, dashboard scene boundary logging, and signal-based clean shutdown.
- Added `internal/runtime/v3runtime` orchestration tests covering v3 config load, selected vehicle connection, selected JSONL event output, and reader cleanup.
- Split v3.1 implementation prompts into per-slice files under `docs/v3.1/`.
- Reframed PR 51 as docs-only planning setup, with implementation starting at `v3.1.0`.
- Added blocking and impact metadata to v3.1 open decisions.
- Clarified v3.1 slice docs-update rules for `MigrationState.md` and `OpenDecisions.md`.
- Revised the v3.1 release plan around the runnable command path, display adapter, dashboard/gauge test harness, and dashboard update cadence targets.
- Expanded v3.1 carry-forward and release plan docs with v3.0 implementation details and per-slice checkpoints.
- Added v3.1 release planning stubs under `docs/v3.1/` for the next implementation phase.
- Added v3 inverse implementation audit documentation for old/current behaviours not yet fully rebuilt as v3.
- Added v3 retirement audit documentation for old/current paths that may be reviewed for later removal or archiving.
- Added v3 richer dashboard widget rendering for `bar_display` and `frame_gauge`.
- Added dashboard tests for bar fill mapping, reverse fill direction, zones, frame clamping, sensor status handling, and unchanged frame output handling.
- Added v3 richer asset registry support for bar and frame asset families.
- Added reusable decoded bar cell and frame asset structs for later dashboard widgets.
- Added v3 smallest selected-dashboard scene runtime for image, digit display, and indicator widgets.
- Added selected-dashboard scene tests for RuntimePlan dashboard selection, digit formatting, decimal point overlays, indicator status mapping, and unchanged formatted output handling.
- Added v3 minimal asset registry for image, digit, and indicator asset families.
- Added reusable decoded image asset structs so future widgets can avoid hot-path asset loading.
- Added tests for repository-root asset path resolution, missing asset errors, decoded digit assets, and required indicator states.
- Updated v3 migration state for the v3.0.11 retirement audit slice.
- Updated v3 migration state for the v3.0.10 implementation slice.
- Updated v3 migration state for the v3.0.9 implementation slice.
- Updated v3 migration state for the v3.0.7 implementation slice.

## 0.1 - 2026-06-08

- Created PR-tail package for GoDriveLog dashboard v2.7 throttle fixture completion.
- Restored throttle frame count from 3 to 11 in `config.example.yaml`.
- Added placeholder SVG throttle frames 003 through 010 for 30% through 100%.
- Kept changes limited to the v2.7 example/dashboard fixture assets.
