# v3.5 Implementation State

Status: v3.5.0 docs complete

Current target: v3.5.1 Gauge Preview Mode

Current branch: `v3.5-gauge-realism-docs-2`

## Scope

v3.5 is the gauge realism pass.

It adds believable gauge behaviour without changing the v3.4 gauge family model. The intent is to make gauges look more like real mechanisms when values change, while avoiding perpetual ambient effects.

The final v3.5 tail also includes two small radial-only display refinements that need renderer support: an optional needle drop shadow and an optional display-only calibration offset.

## Current decisions

- All v3.5 realism options live under the `realism` key.
- Keep realism config collapsed where possible.
- `realism.movement` is a scalar, not a nested object.
- `realism.movement` supports `click` and `smooth`.
- If `realism.movement` is omitted, it defaults to `click`.
- `click` means the display updates directly to the next value unless another enabled realism option adds visible movement.
- Existing top-level `movement` may remain supported for backwards compatibility, but new v3.5 config should use `realism.movement`.
- Unknown realism options must fail config loading.
- Known realism options used on unsupported gauge types must fail config loading.
- `realism.order` may optionally control the order of enabled realism behaviours.
- Do not rely on YAML key order to control behaviour order.
- Gauge Preview Mode is the simple visual viewer for one gauge at a time.
- Gauge Preview Mode CLI is `godrivelog dashboard preview <file>`.
- `<file>` is mandatory and positional.
- Preview files are normal YAML configs, not a special metadata system.
- Each single-feature preview file should enable one realism feature only.
- Each gauge type may also have one deliberate `99-all-options` preview file.
- Radial needle shadow is a static renderer feature, not dynamic parallax or lighting.
- Radial calibration offset is display-only and must not change input values.
- Hysteresis applies only to radial and bar gauges in v3.5.

## Approved v3.5 realism options

| Option | Applies to |
|---|---|
| `movement` | relevant gauge types |
| `wraparound` | odometer |
| `drum_slop` | odometer |
| `carry_drag` | odometer |
| `snap_settle` | odometer |
| `backlash` | odometer |
| `hysteresis` | radial, bar |
| `stiction` | radial |
| `damping` | radial, bar |
| `overshoot` | radial, bar |
| `peg_bounce` | radial |
| `needle_shadow` | radial |
| `calibration_offset` | radial |

## Scope boundaries

Allowed in v3.5:

- static imperfection;
- finite value-change movement;
- Gauge Preview Mode;
- deterministic, bounded behaviour;
- display-only realism options;
- small radial-only display refinements that need renderer support.

Not allowed in v3.5:

- idle needle vibration;
- random flicker or shimmer;
- gas-discharge jitter;
- LED multiplex flicker;
- power-on sweep;
- brownout dip;
- lazy power-off;
- dynamic parallax or gyro/light-driven shadow movement;
- general physics engine;
- generated artwork, videos, screenshot reports, or visual diff machinery.

## Slice checklist

- [x] v3.5.0 movement realism docs
- [ ] v3.5.1 Gauge Preview Mode
- [ ] v3.5.2 odometer wraparound
- [ ] v3.5.3 odometer drum slop
- [ ] v3.5.4 finite movement lifecycle
- [ ] v3.5.5 shared movement policy
- [ ] v3.5.6 odometer eased roll
- [ ] v3.5.7 odometer carry-drag / 9-drag
- [ ] v3.5.8 radial damping
- [ ] v3.5.9 radial stiction
- [ ] v3.5.10 radial/bar overshoot
- [ ] v3.5.11 radial peg bounce
- [ ] v3.5.12 indicator thermal fade
- [ ] v3.5.13 bar smoothing
- [ ] v3.5.14 odometer snap / settle
- [ ] v3.5.15 odometer backlash
- [ ] v3.5.16 display-only hysteresis
- [ ] v3.5.17 radial needle drop shadow
- [ ] v3.5.18 radial calibration offset

## Next-slice workflow

When asked to do the next slice:

1. Read this file.
2. Find the first unchecked slice.
3. Read `docs/v3.5/ReleasePlan.md`.
4. Read the matching prompt in `docs/v3.5/prompts/`.
5. Make only that slice's changes.
6. Update this checklist and any relevant docs.
7. Do not implement later slices early.
