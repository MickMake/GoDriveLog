# GoDriveLog gauge package examples

Version: 0.1

Gauge-package configuration examples.

These files are intentionally YAML-first. Artwork is represented by expected file paths only, so the real PNG/SVG assets can be added or replaced later.

## Contents

- `assets/gauges/7Seg/amber/**/gauge.yaml` - amber numeric packages using seven-segment artwork, covering 2, 3, 4, and 5 digits.
- `assets/gauges/7Seg/green/**/gauge.yaml` - green numeric packages using seven-segment artwork, covering 2, 3, 4, and 5 digits.
- `assets/gauges/radial/simple_rpm/gauge.yaml` - radial RPM gauge package.

## Gauge model

Dashboard widgets place gauges:

```yaml
- id: rpm
  type: gauge
  gauge: assets/gauges/7Seg/amber/4_digit_rpm
  position: [780, 40]
  scale: 1.0
```

Gauge packages own sensor binding, formatting, value mapping, layout geometry, and local asset references:

```yaml
id: amber_4_digit_rpm
type: numeric
sensor: rpm
format: "%04.0f"
```

## Notes

- Directory names are examples only; the loader cares about `gauge.yaml`, not inferred type from path names.
- `7Seg` is the current example directory name for numeric packages that use seven-segment artwork.
- Asset paths inside `gauge.yaml` are relative to that package file.
- Relative paths such as `../` and `../../` are acceptable when they stay inside the asset tree and do not wander up and back down through unrelated folders.
- These examples are runnable gauge package fixtures for the active v3 dashboard path.
