# v3.6 Release Plan

v3.6 follows the completed v3.5 gauge realism pass.

v3.6 starts with pointer markers for radial and bar gauges, then keeps room for small follow-on gauge enhancements. The goal is to add useful instrument-style behaviours without turning the dashboard into a physics engine with opinions.

## Theme

v3.6 pointer markers are instrument-realism features, not statistical overlays.

Pointer markers track what the gauge indicator visibly does:

- radial markers follow the rendered needle angle;
- bar markers follow the rendered bar position/fill;
- if no realism is enabled, the rendered indicator path represents true mapped data;
- if realism is enabled, the marker follows the realistic rendered movement.

This means min/max markers can capture overshoot, bounce, damping lag, or other rendered indicator behaviour when those effects are enabled.

## Config shape

Use the `realism.pointer_markers` key.

Example radial or bar gauge config:

```yaml
realism:
  pointer_markers:
    max: true
    min: true
```

Expanded future-friendly shape:

```yaml
realism:
  pointer_markers:
    max:
      enabled: true
    min:
      enabled: true
    damped:
      enabled: true
```

Implementations may accept the short boolean form if the existing config parser supports it cleanly, but the documented long form should remain available for later options.

## Pointer marker semantics

### Max marker

The max marker records the highest rendered indicator position seen since the marker was reset.

For radial gauges, this means the highest rendered needle angle or equivalent mapped position, depending on existing radial geometry conventions.

For bar gauges, this means the furthest/highest rendered bar fill or bar indicator position, depending on the existing bar orientation model.

### Min marker

The min marker records the lowest rendered indicator position seen since the marker was reset.

### Damped secondary marker

The damped marker is a slow secondary indicator. It is inspired by heavily damped mechanical gauges that visually steady a fluctuating reading.

It is not a mathematical average and must not be documented as an arithmetic mean.

A later slice may choose a simple deterministic damping model that follows the same input target or rendered indicator target with heavier damping than the live indicator.

## Odometer backlash tail slice

Odometer `backlash` is required as a v3.6 tail cleanup slice.

It is already documented as an odometer realism behaviour in v3.5 planning, but current code support was not found. Existing odometer realism can create general mechanical feel, but cannot fully create direction-change slack.

`backlash` should mean:

> when an odometer-style value reverses direction, wheel movement shows a small bounded amount of mechanical slack before following the new direction and settling exactly on the target.

The implementation must be:

- odometer-only;
- display-only;
- deterministic;
- bounded and subtle;
- disabled unless configured;
- non-mutating for source values, logs, exports, configured ranges, or input data.

Do not treat `movement: smooth` or `movement: click` as required future movement implementations while doing this cleanup. `linear`, `ease_out`, and `bell` are already smooth movement modes. `click` remains reserved unless a later prompt defines a genuinely distinct stepped-wheel behaviour.

## Reset/session behaviour

Pointer markers are runtime/session state in v3.6.

Default reset events:

- gauge instance creation;
- dashboard/runtime reload;
- relevant gauge config change;
- data source identity change.

Do not add database persistence in v3.6 unless a later prompt explicitly promotes it.

## Rendering model

Pointer markers should render as separate visual indicators appropriate to the gauge family.

Radial marker assets may use names such as:

```text
needle_max.png
needle_min.png
needle_damped.png
```

Bar marker assets or render primitives should be family-specific. Do not force radial needle asset conventions onto bar gauges.

## Slice plan

| Slice | Name | Intent |
| --- | --- | --- |
| v3.6.0 | Pointer marker planning docs | Establish docs, prompts, scope, and terminology. |
| v3.6.1 | Radial pointer marker max | Add radial max marker following rendered needle path. |
| v3.6.2 | Radial pointer marker min | Add radial min marker using the same semantics. |
| v3.6.3 | Pointer marker reset/session behaviour | Define and implement reset/session rules. |
| v3.6.4 | Radial damped secondary pointer marker | Add slow secondary radial pointer, not mathematical average. |
| v3.6.5 | Bar pointer marker max | Add max marker for rendered bar path/fill. |
| v3.6.6 | Bar pointer marker min | Add min marker for rendered bar path/fill. |
| v3.6.7 | Bar damped secondary pointer marker | Add slow secondary bar marker. |
| v3.6.8 | Enhancement backlog triage | Review and promote later v3.6 candidates. |
| v3.6.9 | Odometer backlash cleanup | Implement missing odometer `backlash` and document `smooth`/`click` cleanup decisions. |

## Backlog candidates

These are possible v3.6+ enhancements, not first-tranche requirements:

- explicit manual marker reset control;
- optional marker labels/tooltips;
- marker preview fixtures for all gauge families;
- indicator marker equivalents;
- numeric/segmented baseline preview improvements;
- mathematical statistics overlays, if separately named and scoped;
- persistent marker state;
- marker styling overrides beyond default assets.

## Non-goals

- Do not mutate source values, logs, exports, configured ranges, or input data.
- Do not add true statistical average under `pointer_markers`.
- Do not implement persistence as part of min/max.
- Do not bundle radial and bar renderer changes in the same slice.
- Do not let v3.6 become a random wishlist drawer. The drawer already contains three screws, a dead battery, and a suspiciously important bracket.
