# v3.4 example layout moves

## Purpose

This cleanup rehomes the generated v3.4 example dashboards so each dashboard is self-contained under `examples/<dashboard_name>/`.

The goal is to remove the split between `examples/dashboards/`, `examples/assets/v3.4/`, and `assets/gauges/v3.4/` so future example slices do not keep rediscovering the same layout pothole with a lantern and a stern face.

## Old layout summary

Before this cleanup, the generated example dashboards were split across:

- `examples/dashboards/*.yaml`
- `examples/assets/v3.4/<theme>/...`
- `assets/gauges/v3.4/ornate-timber/...`

## New layout summary

After this cleanup, each generated example dashboard is self-contained under:

```text
examples/<dashboard_name>/
  dashboard.yaml
  assets/
    panel/
      background.png
      foreground.png
    gauges/
      <gauge_name>/
        gauge.yaml
        <gauge assets>
```

## Moved paths

| Status | Old path | New path | Notes |
|---|---|---|---|
| moved | `examples/dashboards/framework-smoke.yaml` | `examples/framework-smoke/dashboard.yaml` | Framework smoke dashboard config now lives beside the example assets. |
| moved | `examples/assets/v3.4/framework-smoke/panel/*.png` | `examples/framework-smoke/assets/panel/*.png` | Dashboard-local panel artwork. |
| moved | `examples/assets/v3.4/framework-smoke/digits/*.png` | `examples/framework-smoke/assets/digits/*.png` | Dashboard-local digit artwork. |
| moved | `examples/assets/v3.4/framework-smoke/indicator/*.png` | `examples/framework-smoke/assets/indicator/*.png` | Dashboard-local indicator artwork. |
| moved | `examples/dashboards/ornate-timber.yaml` | `examples/ornate-timber/dashboard.yaml` | Ornate timber dashboard config now lives beside the example assets. |
| moved | `examples/assets/v3.4/ornate-timber/panel/*.png` | `examples/ornate-timber/assets/panel/*.png` | Dashboard-local panel artwork. |
| moved | `examples/assets/v3.4/ornate-timber/gauges/speed_numeric/*.png` | `examples/ornate-timber/assets/gauges/speed_numeric/*.png` | Numeric gauge artwork moved into the example dashboard tree. |
| moved | `examples/assets/v3.4/ornate-timber/gauges/radial_rpm/*.png` | `examples/ornate-timber/assets/gauges/radial_rpm/*.png` | Radial gauge artwork moved into the example dashboard tree. |
| moved | `examples/assets/v3.4/ornate-timber/gauges/trip_odometer/*.png` | `examples/ornate-timber/assets/gauges/trip_odometer/*.png` | Odometer gauge artwork moved into the example dashboard tree. |
| moved | `examples/assets/v3.4/ornate-timber/gauges/check_engine_indicator/*.png` | `examples/ornate-timber/assets/gauges/check_engine_indicator/*.png` | Indicator gauge artwork moved into the example dashboard tree. |
| moved | `examples/assets/v3.4/ornate-timber/gauges/fuel_bar/*.png` | `examples/ornate-timber/assets/gauges/fuel_bar/*.png` | Bar gauge artwork moved into the example dashboard tree. |
| moved | `examples/assets/v3.4/ornate-timber/gauges/rpm_segmented/*.png` | `examples/ornate-timber/assets/gauges/rpm_segmented/*.png` | Segmented gauge artwork moved into the example dashboard tree. |
| moved | `assets/gauges/v3.4/ornate-timber/speed_numeric/gauge.yaml` | `examples/ornate-timber/assets/gauges/speed_numeric/gauge.yaml` | Gauge package YAML now sits beside the gauge assets. |
| moved | `assets/gauges/v3.4/ornate-timber/radial_rpm/gauge.yaml` | `examples/ornate-timber/assets/gauges/radial_rpm/gauge.yaml` | Gauge package YAML now sits beside the gauge assets. |
| moved | `assets/gauges/v3.4/ornate-timber/trip_odometer/gauge.yaml` | `examples/ornate-timber/assets/gauges/trip_odometer/gauge.yaml` | Gauge package YAML now sits beside the gauge assets. |
| moved | `assets/gauges/v3.4/ornate-timber/check_engine_indicator/gauge.yaml` | `examples/ornate-timber/assets/gauges/check_engine_indicator/gauge.yaml` | Gauge package YAML now sits beside the gauge assets. |
| moved | `assets/gauges/v3.4/ornate-timber/fuel_bar/gauge.yaml` | `examples/ornate-timber/assets/gauges/fuel_bar/gauge.yaml` | Gauge package YAML now sits beside the gauge assets. |
| moved | `assets/gauges/v3.4/ornate-timber/rpm_segmented/gauge.yaml` | `examples/ornate-timber/assets/gauges/rpm_segmented/gauge.yaml` | Gauge package YAML now sits beside the gauge assets. |
| removed | `examples/assets/v3.4/README.md` | n/a | Old layout note removed because the new layout supersedes it. |

## Old paths intentionally left behind

None.

## Old paths that should no longer be used

- `examples/dashboards/framework-smoke.yaml`
- `examples/dashboards/ornate-timber.yaml`
- `examples/assets/v3.4/framework-smoke/`
- `examples/assets/v3.4/ornate-timber/`
- `assets/gauges/v3.4/ornate-timber/`
- `examples/assets/v3.4/README.md`

## Canonical paths after this PR

- `examples/framework-smoke/dashboard.yaml`
- `examples/framework-smoke/assets/panel/background.png`
- `examples/framework-smoke/assets/panel/foreground.png`
- `examples/framework-smoke/assets/digits/*.png`
- `examples/framework-smoke/assets/indicator/*.png`
- `examples/ornate-timber/dashboard.yaml`
- `examples/ornate-timber/assets/panel/background.png`
- `examples/ornate-timber/assets/panel/foreground.png`
- `examples/ornate-timber/assets/gauges/speed_numeric/gauge.yaml`
- `examples/ornate-timber/assets/gauges/radial_rpm/gauge.yaml`
- `examples/ornate-timber/assets/gauges/trip_odometer/gauge.yaml`
- `examples/ornate-timber/assets/gauges/check_engine_indicator/gauge.yaml`
- `examples/ornate-timber/assets/gauges/fuel_bar/gauge.yaml`
- `examples/ornate-timber/assets/gauges/rpm_segmented/gauge.yaml`
