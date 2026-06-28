# v3.5 Prompts

These prompts define the v3.5 slice sequence.

## Codex Usage

These prompt files are intended for ChatGPT, Codex, or a human applying the same slice rules.

If the user says any of the following:

- "implement the next slice"
- "do the next slice"
- "continue v3.5"
- "start the next v3.5 slice"

the agent must:

1. Read `docs/v3.5/ImplementationState.md`.
2. Find the first unchecked slice.
3. Read `docs/v3.5/ReleasePlan.md`.
4. Read `docs/v3.5/RealismBehaviourGuide.md`.
5. Read the matching prompt file under `docs/v3.5/prompts/`.
6. Implement only that slice.
7. Update implementation state and relevant docs.
8. Do not implement later slices early.
9. After the slice is complete, follow the finalisation / PR cycle in `docs/v3.5/ImplementationState.md`.

## Prompt files

- `v3.5.0-movement-realism-docs.md`
- `v3.5.1-manual-gauge-inspection-harness.md` - legacy filename; this slice is Gauge Preview Mode
- `v3.5.2-odometer-wraparound.md`
- `v3.5.3-odometer-drum-slop.md`
- `v3.5.4-finite-movement-lifecycle.md`
- `v3.5.5-shared-movement-policy.md` - groundwork only; `realism.movement_policy` is not used for odometer movement
- `v3.5.6-odometer-eased-roll.md` - legacy filename; this slice is now odometer main movement
- `v3.5.6a-document-odometer-movement-goal.md`
- `v3.5.6b-implement-odometer-movement-model.md`
- `asset-root-rule.md` - separate follow-up prompt for the gauge asset-root rule
- `v3.5.7-odometer-carry-drag.md`
- `v3.5.8-radial-damping.md`
- `v3.5.9-radial-stiction.md`
- `v3.5.10-radial-overshoot.md`
- `v3.5.11-radial-peg-bounce.md`
- `v3.5.12-indicator-thermal-fade.md`
- `v3.5.13-bar-smoothing.md`
- `v3.5.14-odometer-snap-settle.md`
- `v3.5.15-odometer-backlash.md`
- `v3.5.16-display-only-hysteresis.md`
- `v3.5.17-radial-needle-drop-shadow.md`
- `v3.5.18-radial-calibration-offset.md`

## Shared rules

- Keep each slice small.
- Preserve v3.4 gauge semantics.
- Do not add idle animation or ambient effects.
- Do not add a general physics engine.
- Most new realism configuration belongs under `realism`.
- `movement` is the exception: it is a single scalar movement knob that should be accepted by any gauge type for now.
- Gauge types without concrete movement behaviour may parse `movement` and keep their current immediate behaviour until their movement slice defines more.
- Odometer movement is controlled by `odometer.movement`.
- Do not use or recommend `realism.movement_policy` for odometer movement.
- Unknown realism options must produce a clear configuration error.
- Known realism options used on unsupported gauge types must produce a clear configuration error.
- Unknown movement values must fail configuration loading unless that gauge type explicitly documents a recognised fallback.
- Add checks where behaviour can be asserted.
- Add or update Gauge Preview Mode YAML only where useful.
- Preview files are normal YAML; do not add a special metadata layer.
- Keep asset-root changes out of movement slices unless the slice explicitly targets asset loading.
