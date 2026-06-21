# GoDriveLog v3.2 baseline dashboard verification

Status: v3.2.8 baseline verification

## Purpose

This baseline exercises the realistic v3.2 gauge workload through the existing dashboard harness and Fyne display adapter. It is verification-only: it does not redesign renderers, introduce Ebiten, add sensor overrides, or change dashboard runtime ownership rules.

The baseline config is:

```text
docs/v3.2/baseline-dashboard.yaml
```

The workload is:

- 3 digit seven-segment temperature display using `coolant_temperature`, range `-10..40`, including `-` glyph rendering via `format: "%03.0f"`.
- 3 digit seven-segment speed display using `speed`.
- 4 digit seven-segment RPM display using `rpm`.
- Radial RPM gauge using the same `rpm` sensor as the RPM digit display.

## Primary command

```bash
go run ./cmd/GoDriveLog --harness --config ./docs/v3.2/baseline-dashboard.yaml --vehicle vw_caddy --pattern sweep --interval 50ms
```

Close the harness window or send `Ctrl-C` to stop the run and emit the summary line.

## Expected summary fields

Each successful harness run should finish with a log line matching this shape:

```text
v3 dashboard harness summary: vehicle=vw_caddy sensors=3 dashboards=1 pattern=<pattern> interval=50ms events=<n> display_submitted=<n> display_rendered=<n> display_superseded=<n> display_last_render=<duration>
```

Record these fields for each verification run:

| Field | Expected meaning |
|---|---|
| `events` | Total fake sensor events emitted by the harness. With three sensors, this grows by three per harness tick plus the initial emit. |
| `display_submitted` | Number of scene batches accepted by the latest-only display sink. |
| `display_rendered` | Number of scene batches actually rendered by the display adapter. |
| `display_superseded` | Number of submitted scene batches replaced by a newer scene before rendering. This may be non-zero on slower hardware. |
| `display_last_render` | Duration of the last display render. |

## Verification runs

### Fixed pattern

```bash
go run ./cmd/GoDriveLog --harness --config ./docs/v3.2/baseline-dashboard.yaml --vehicle vw_caddy --pattern fixed --interval 50ms
```

Expected visual result:

- temperature display is stable near midpoint of `-10..40`;
- speed display is stable near midpoint of `0..220`;
- RPM digit display and radial RPM gauge show the same stable RPM-derived midpoint;
- after the initial render, `display_submitted` should usually grow slowly or stay low because formatted values stop changing.

### Sweep pattern

```bash
go run ./cmd/GoDriveLog --harness --config ./docs/v3.2/baseline-dashboard.yaml --vehicle vw_caddy --pattern sweep --interval 50ms
```

Expected visual result:

- temperature sweeps across `-10..40`, proving the minus glyph appears around negative values;
- speed sweeps across `0..220`;
- RPM digit display and radial RPM gauge move together from low to high and back;
- `display_submitted` should increase through the run;
- `display_rendered` may be lower than `display_submitted` when the latest-only sink coalesces frames;
- `display_superseded` may be non-zero, especially on slower hardware or with the radial gauge visible.

### Heartbeat pattern

```bash
go run ./cmd/GoDriveLog --harness --config ./docs/v3.2/baseline-dashboard.yaml --vehicle vw_caddy --pattern heartbeat --interval 50ms
```

Expected visual result:

- values pulse through quick peaks and dips;
- temperature should again enter the negative range and exercise the minus glyph;
- RPM digit display and radial RPM gauge should pulse together;
- this run is useful for spotting redraw lag or stale-frame behaviour.

### Non-ok / missing state

The current CLI harness patterns emit `ok` sensor states only. Non-ok and missing sensor-state coverage already exists at lower dashboard/runtime test boundaries, but this baseline does not add a CLI switch or renderer redesign for non-ok simulation.

## Notes

- The baseline deliberately reuses existing v3.2 gauge package examples under `docs/v3.2/assets/gauges/`.
- The only added package is `assets/gauges/7Seg/green/3_digit_temp`, which reuses existing three-digit seven-segment artwork so `-10..40` can be verified without changing renderer code.
- The RPM digit display and radial RPM gauge intentionally share the `rpm` sensor so mismatched update behaviour is obvious.
