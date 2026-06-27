# v3.5 Prompts

These prompts are written for any coding agent. They avoid tool-specific language and should work whether the work is started from a short instruction like:

```text
do the next slice
```

or from a more explicit implementation request.

## How to use

1. Read `docs/v3.5/ImplementationState.md`.
2. Find the first unchecked slice.
3. Read `docs/v3.5/ReleasePlan.md`.
4. Read the matching prompt below.
5. Implement only that slice.
6. Update implementation state and relevant docs.
7. Do not implement future slices early.

## Prompt files

- `v3.5.0-movement-realism-docs.md`
- `v3.5.1-manual-gauge-inspection-harness.md`
- `v3.5.2-odometer-wraparound.md`
- `v3.5.3-odometer-drum-slop.md`
- `v3.5.4-finite-movement-lifecycle.md`
- `v3.5.5-shared-movement-policy.md`
- `v3.5.6-odometer-eased-roll.md`
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
- Keep odometer `smooth` and `click` base modes intact.
- Add tests where behaviour can be asserted.
- Add or update visual harness YAML only where useful.
- Do not implement asset-only presentation work in these code slices.
