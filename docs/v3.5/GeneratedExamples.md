# v3.5 Generated Examples

v3.5 does not need generated artwork.

The examples for this release are visual inspection fixtures: small YAML files that render one gauge at a time in the manual inspection harness.

## Fixture intent

Fixtures should make behaviour easy to see, not produce polished dashboards.

Use normal dashboard/gauge YAML where possible. Do not add a separate metadata layer unless a later slice proves it is necessary.

## Suggested fixture location

```text
examples/harness/gauge-realism/
  odometer/
  radial/
  indicator/
  bar/
  segmented/
  numeric/
```

## Fixture rules

For each relevant gauge type:

- `00-baseline.yaml` uses no realism options.
- Each numbered feature file enables one realism option only.
- `99-all-options.yaml` deliberately enables all supported options for that gauge type.

Do not add arbitrary combinations. The all-options file is the only combination case.

## Odometer notes

Odometers already support `movement: smooth` and `movement: click`. Keep those base modes.

Prefer feature fixtures against the default mode unless the feature specifically needs both smooth and click coverage.

Acceptable examples:

```text
odometer/00-baseline-smooth.yaml
odometer/00-baseline-click.yaml
odometer/01-wraparound.yaml
odometer/05-snap-settle-click.yaml
odometer/99-all-options-smooth.yaml
odometer/99-all-options-click.yaml
```

Do not double every fixture automatically. Only add click-specific coverage where it helps visual inspection.

## Harness controls

The manual harness should infer the start value from the gauge value range:

```text
start = (min + max) / 2
```

Use keyboard input for manual inspection:

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

Do not add generated art, videos, screenshot reports, or visual diff outputs in v3.5.
