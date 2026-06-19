# GoDriveLog v3.2 config examples

Version: 0.1

Draft v3.2 gauge-package configuration examples.

These files are intentionally YAML-first. Artwork is represented by expected file paths only, so the real PNG/SVG assets can be added or replaced later.

## Contents

- `assets/gauges/seven_segment/amber/**/gauge.yaml` - amber seven-segment packages covering 2, 3, 4, and 5 digits.
- `assets/gauges/seven_segment/green/**/gauge.yaml` - green seven-segment packages covering 2, 3, 4, and 5 digits.
- `assets/gauges/radial/simple_rpm/gauge.yaml` - radial RPM gauge package.

## Intended v3.2 model

Dashboard widgets place gauges:

```yaml
- id: rpm
  type: gauge
  gauge: assets/gauges/seven_segment/amber/4_digit_rpm
  position: [780, 40]
  scale: 1.0
```

Gauge packages own sensor binding, formatting, value mapping, layout geometry, and local asset references:

```yaml
id: amber_4_digit_rpm
type: seven_segment
sensor: rpm
format: "%04.0f"
```

## Notes

- Directory names are examples only; v3.2 should care about `gauge.yaml`, not infer type from path names.
- Asset paths inside `gauge.yaml` are relative to that package file.
- Shared artwork paths use `../../../shared/...` from the package directories.
- These examples are design fixtures, not guaranteed runnable config until the v3.2 loader/schema slices exist.
