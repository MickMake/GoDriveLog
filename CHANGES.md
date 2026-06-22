# CHANGES

## Unreleased

- Updated the `v3.3.1` experimental Ebiten renderer spike to use an `ebiten` build-tagged command path, avoiding Linux GLFW linker symbol collisions between Fyne and Ebiten.
- Added `v3.3.1` experimental Ebiten renderer spike behind `--renderer ebiten`, keeping `--renderer fyne` as the default and feeding both renderers from the same v3 runtime/harness dashboard scene path.
- Added `--duration` for v3 runtime and harness runs so baseline renderer comparison commands can stop automatically after fixed intervals.
- Added a narrow `internal/dashboard/adapter/ebiten` scene adapter that caches decoded image assets, renders static seven-segment/gauge layers, and rotates radial needles at draw time for measurement before deciding whether prepared frames are needed.
- Added `v3.3.0` renderer planning docs under `docs/v3.3/`, including release plan, implementation state, baseline verification notes, and intent-named prompts for `v3.3.0` through `v3.3.3`.
- Added reusable `examples/baseline-dashboard.yaml` for the renderer comparison workload, with Fyne and Ebiten command examples using the same full dashboard path and fixed-duration runs.
- Updated `examples/README.md` to identify reusable dashboard examples and the shared repository-root `assets/` tree for active assets.
- Added a `docs/v3.3/baseline-dashboard.yaml` pointer so v3.3 mirrors the v3.2 root filenames while keeping the runnable config canonical under `examples/`.
- Added a CI-visible `v3.2 baseline harness verification` workflow step that runs the baseline harness pattern test headlessly without launching Fyne.
- Added CI-safe harness coverage for the v3.2 baseline dashboard config across `fixed`, `sweep`, and `heartbeat` patterns, asserting the selected vehicle, three sensors, one dashboard, four baseline gauge widgets, and deterministic event completion.
- Added `v3.2.8` baseline dashboard harness configuration at `docs/v3.2/baseline-dashboard.yaml`, exercising three-digit temperature, three-digit speed, four-digit RPM, and radial RPM gauge widgets through the existing harness/display path.
- Added `docs/v3.2/BaselineDashboardVerification.md` documenting fixed, sweep, heartbeat, and current non-ok/missing-state harness verification expectations, including `events`, `display_submitted`, `display_rendered`, `display_superseded`, and `display_last_render` summary fields.
- Added a small `docs/v3.2/assets/gauges/7Seg/green/3_digit_temp` gauge package that reuses existing green seven-segment artwork to verify `-10..40` temperature output and minus-symbol rendering.
- Updated the v3.2 implementation state to mark `v3.2.6` completed, `v3.2.7` skipped/absorbed by existing examples, `v3.2.8` completed, and `v3.2.9` as the next renderer checkpoint.
- Updated the v3 Fyne display scene sink path to use non-blocking latest-only submissions for live dashboard and harness updates, preventing display rendering from throttling sensor/harness event cadence while preserving latest-frame coalescing and render error visibility.
- Updated `v3.2.6` Fyne radial rendering to prepare 1-degree radial needle frame sets outside normal live update sweeps, keeping live updates to keyed image resource swaps and preserving keyed canvas object reuse.
- Added `v3.2.6` Fyne radial gauge rendering, including ordered radial layer rendering, image-space needle rotation, normalised pivot placement, rotated-needle resource caching, and adapter coverage.
- Added `v3.2.5` radial gauge scene model support, including dashboard runtime routing, package-owned pivots, value-map angle calculation, needle scene part data, non-ok needle suppression, and radial scene signatures.
- Added `v3.2.4` Fyne seven-segment rendering hardening: stable keyed canvas object reuse, glass overlay ordering coverage, and a deterministic adapter benchmark for repeated digit updates.
- Added `v3.2.3` seven-segment gauge support through the dashboard scene path, including `type: gauge` package loading, package-owned sensor state, static layers, digit positions, package-owned formatting, non-ok suppression, Fyne adapter positioning, and scene signatures.
- Removed the redundant post-`Validate` gauge widget ownership pass now that ownership validation lives inside `Validate`.
- Added `v3.2.2` dashboard config support for `type: gauge` widgets with package-owned gauge paths, placement, and scale.
