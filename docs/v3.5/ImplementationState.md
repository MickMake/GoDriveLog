# v3.5 Implementation State

Status: v3.5.0 docs branch in progress

Current target: v3.5.0 movement realism docs

Current branch: `v3.5-gauge-realism-docs`

## Scope

v3.5 is the gauge realism pass.

It adds believable gauge behaviour without changing the v3.4 gauge family model. The intent is to make gauges look more like real mechanisms when values change, while avoiding perpetual ambient effects.

The final v3.5 tail also includes two small radial-only display refinements that need renderer support: an optional needle drop shadow and an optional display-only calibration offset.

## Current decisions

- The existing odometer movement modes `smooth` and `click` remain valid base modes.
- Odometer realism options layer on top of `smooth`/`click`; they do not replace those modes.
- The manual inspection harness is deliberately simple and visual-first.
- Harness showcase files are normal YAML configs, not a special metadata system.
- Harness start value is the midpoint of the configured value range.
- Each single-feature harness case should enable one realism feature only.
- Each gauge type may also have one deliberate `99-all-options` case.
- Radial needle shadow is a static renderer feature, not dynamic parallax or lighting.
- Radial calibration offset is display-only and must not change input values.

## Scope boundaries

Allowed in v3.5:

- static imperfection;
- finite value-change movement;
- manual visual inspection cases;
- deterministic, bounded behaviour;
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
- general physics engine.

## Slice checklist

- [ ] v3.5.0 movement realism docs
- [ ] v3.5.1 manual gauge inspection harness
- [ ] v3.5.2 odometer wraparound
- [ ] v3.5.3 odometer drum slop
- [ ] v3.5.4 finite movement lifecycle
- [ ] v3.5.5 shared movement policy
- [ ] v3.5.6 odometer eased roll
- [ ] v3.5.7 odometer carry-drag / 9-drag
- [ ] v3.5.8 radial damping
- [ ] v3.5.9 radial stiction
- [ ] v3.5.10 radial overshoot
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
