# v3.5 Preview Files

v3.5 does not need generated artwork.

The examples for this release are Gauge Preview Mode files: small, normal YAML files that render one gauge at a time so realism options can be judged by eye.

## Preview intent

Preview files should make behaviour easy to see, not produce polished dashboards.

Use normal dashboard/gauge YAML. Do not add a separate metadata layer unless a later slice proves there is no sane alternative.

## Suggested location

```text
examples/gauge-realism/
  odometer/
  radial/
  bar/
  numeric/
  segmented/
  indicator/
```

## File rules

For each relevant gauge type:

- `00-baseline.yaml` uses no realism options.
- Each numbered feature file enables one realism option only.
- `99-all-options.yaml` deliberately enables all supported options for that gauge type.

Do not add arbitrary combinations. The all-options file is the only combination case.

Numeric, segmented, and indicator gauges should only get baseline preview files in v3.5 unless a specific realism option applies to them.

## Odometer notes

Odometers support `realism.movement: click` and `realism.movement: smooth`.

If `realism.movement` is omitted, it defaults to `click`.

Prefer feature preview files against the default mode unless the feature specifically needs both smooth and click coverage.

Acceptable examples:

```text
odometer/00-baseline-click.yaml
odometer/00-baseline-smooth.yaml
odometer/01-wraparound.yaml
odometer/05-snap-settle-click.yaml
odometer/99-all-options-click.yaml
odometer/99-all-options-smooth.yaml
```

Do not double every preview file automatically. Only add smooth/click-specific coverage where it helps visual judgement.

## Gauge Preview CLI

Gauge Preview Mode is launched with:

```text
godrivelog dashboard preview <file>
```

`<file>` is mandatory and must point to a normal dashboard/gauge YAML file.

Optional flags may select the gauge inside the file, override the starting value, or tune preview step sizes. Realism options stay in YAML.

Suggested optional flags:

```text
--gauge <name>
--value <number>
--step <number>
--fine-step <number>
--coarse-step <number>
```

## Preview controls

The preview should infer the start value from the gauge value range unless `--value` is supplied:

```text
start = (min + max) / 2
```

Use keyboard input for manual preview:

- Left arrow: min.
- Right arrow: max.
- Up arrow: increment.
- Down arrow: decrement.
- Shift + Up/Down: coarse increment/decrement.
- Ctrl/Cmd + Up/Down: fine increment/decrement.
- R: reset to midpoint.
- Space: replay last transition.
- Esc/Q: quit.
- Mouse wheel may also increment/decrement.

## What not to add

Do not add generated art, videos, screenshot reports, visual diff outputs, test metadata, or a second CLI-based realism configuration system in v3.5.
