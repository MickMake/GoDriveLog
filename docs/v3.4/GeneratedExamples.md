# GoDriveLog v3.4 generated example dashboards

## Purpose

v3.4.6 through v3.4.9 add generated example dashboards after the gauge behaviour slices are complete.

These examples should prove the final v3.4 gauge set using deterministic, local, reproducible assets. They are not a new renderer model and must not add runtime visual-style fields.

## Core rule

Runtime config continues to describe behaviour:

```text
numeric
radial
odometer
indicator
bar
segmented
```

Visual identity belongs to generated PNG assets, dashboard layout, and gauge package artwork.

Generator-internal theme options are allowed. Runtime `style` fields are not.

## Reproducibility requirements

Generated dashboard assets must be reproducible from committed scripts and stable inputs:

- Use local procedural drawing only.
- Use stable seed/config values where noise is used.
- Commit the source script/config/docs that explain how to regenerate outputs.
- Do not use remote image generation.
- Do not download stock art.
- Do not hand-edit generated PNGs as the source of truth.
- Keep outputs deterministic enough for review.

## Dimension rules

Source asset dimensions are authoritative.

- Do not reinterpret, normalize, or infer one digit set's source dimensions from another digit set.
- If a dashboard needs a smaller or larger rendered display, use dashboard/widget config `scale`.
- Generated example assets may choose their own dimensions, but each generated digit set must stay internally consistent across its slot-positioned assets.
- Decimal points are overlays on the current or preceding digit cell; they do not consume their own digit slot.
- Decimal-point artwork may keep a small visible dot in the lower-right area, but the canvas must align with that digit set's own cell dimensions.

## Planned themes

| Version | Theme | Directory name | Summary |
|---|---|---|---|
| v3.4.6 | Framework | `framework-smoke` or equivalent | Minimal smoke-test theme proving the generation pipeline. |
| v3.4.7 | Ornate timber | `ornate-timber` | Master-carpenter dashboard: multiple timber tones, timber needles, timber ticks, routed/carved/inlaid look. |
| v3.4.8 | Neon grid | `neon-grid` | Dark retro-tech dashboard with Tron-like neon blue glow and grid/circuit accents. |
| v3.4.9 | Steam scrap | `steam-scrap` | Steampunk/scrapyard dashboard with brass/copper/iron plates, pipes, wires, rivets, lamps, and overbuilt decoration. |

## Expected output shape

Exact layout may change during implementation, but generated assets should live under a predictable v3.4 example tree, for example:

```text
examples/assets/v3.4/<theme>/
examples/dashboards/<theme>.yaml
```

The framework may introduce helper modules such as:

```text
scripts/example_assets/
  canvas.py
  wood.py
  metal.py
  glow.py
  gauges.py
  numeric.py
  layout.py
```

Keep the final shape simple. If a single script is enough for v3.4.6, use one script first and extract helpers only when the themed dashboards need them.

## v3.4.6 framework output

The framework slice now uses:

```bash
go run ./scripts/generate-example-assets -theme framework-smoke
```

Current committed smoke-test output:

```text
examples/assets/v3.4/framework-smoke/
examples/dashboards/framework-smoke.yaml
```

The smoke dashboard proves the deterministic asset pipeline, path conventions, and active harness/runtime loading path. It is intentionally small and does not claim the final ornate timber, neon-grid, or steam-scrap art direction yet.

## Gauge coverage target

Each complete themed dashboard should cover as much of the v3.4 gauge model as practical:

- `numeric` displays for values such as speed, voltage, temperature, or diagnostics.
- `radial` gauges for RPM/speed-style values.
- `odometer` for distance/trip-style values.
- `indicator` warning lamps or status tiles.
- `bar` for continuous fuel/temperature-style levels.
- `segmented` for stepped level/alert visuals.

The examples do not need to prove every edge case. They should prove that the completed gauge families can coexist in one coherent dashboard.

## Theme direction

### Ornate timber

The ornate timber dashboard should look like it was made by someone who owns sharp chisels, knows how to use them, and may have opinions about end grain.

Use:

- contrasting timber species/treatments;
- routed or carved panel edges;
- timber needles;
- timber tick marks or inlays;
- darker recessed gauge faces for contrast;
- subtle screw heads, plugs, or joinery details;
- amber/green numeric displays behind smoked acrylic if useful.

Avoid brown soup. The design needs contrast and readable instrumentation.

### Neon grid

The neon-grid dashboard should feel like dark retro tech with a Tron-style influence.

Use:

- near-black dashboard background;
- neon blue as the primary colour;
- cyan/white secondary highlights;
- subtle glow halos around ticks, needles, bars, and indicators;
- circuit/grid accents as background decoration;
- luminous blue numeric assets.

Avoid over-glowing everything. If all pixels shout, none of them say anything useful.

### Steam scrap

The steam-scrap dashboard should look deliberately overbuilt and made from reused mechanical/electrical leftovers.

Use:

- tarnished brass, copper, iron, and aged cream faces;
- unnecessary pipes and fittings as decoration;
- rivets, bolts, mismatched plates, wire loops, and small lamps;
- brass or oxidised metal bezels;
- nixie-ish orange numeric displays where useful;
- visible patched/salvaged construction.

Keep pipes and wires decorative. Do not add renderer behaviour just to make the dashboard look like it needs a pressure certificate.

## Slice boundaries

- v3.4.6 establishes the generation framework and smoke-test path only.
- v3.4.7 adds only the ornate timber dashboard.
- v3.4.8 adds only the neon-grid dashboard.
- v3.4.9 adds only the steam-scrap dashboard.

Do not combine the three themed dashboards in one PR unless explicitly re-scoped later.
