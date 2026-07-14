# Odometer Backlash

Applies to: odometer.

Design status: approved  
Estimated effort: 3–5 Codex hours

## Purpose

Backlash adds a small amount of visible mechanical slack when an odometer reverses direction.

It simulates play in gears, shafts, or couplings where clearance must be taken up before movement continues in the opposite direction.

## Expected visual behaviour

When an odometer-style value reverses direction, wheel movement shows a small bounded amount of mechanical slack before following the new direction and settling exactly on the correct rendered target.

The effect must remain subtle and must never make the displayed value confusing or inaccurate.

## Candidate configuration

If implementation is required, add odometer-only support for:

```yaml
realism:
  backlash: true
```

Keep the first implementation boolean-only unless a small amount or timing shape is already clearly supported by existing parser conventions.

Do not add hidden default behaviour. Existing gauges without `backlash`, or with `backlash: false`, must render exactly as they did before.

## Required behaviour

The implementation must be:

- odometer-only;
- display-only;
- deterministic;
- bounded and subtle;
- active only when configured;
- safe for forward and reverse value changes;
- settled exactly at the target when movement completes.

It must not mutate source values, logs, exports, configured ranges, or input data.

## Direction behaviour

### Same-direction movement

Same-direction odometer movement must not be visibly affected by backlash.

### Direction reversal

When direction reverses:

1. retain a small bounded wheel offset in the previous direction;
2. reduce the temporary backlash offset smoothly to zero while movement continues toward the new target;
3. apply the effect only to wheels that are moving;
4. do not restart or accumulate the effect during repeated updates;
5. settle all wheel offsets exactly on the target when movement completes.

Backlash is a display offset, not a delay to the source value or a separate physics simulation.

## Implementation requirements

If confirmed missing:

- Add config parsing support for `realism.backlash`.
- Validate that `backlash` is accepted only for odometer gauges.
- Add runtime behaviour for direction-change slack during odometer wheel movement.
- Preserve existing behaviour when `backlash` is absent or false.
- Ensure same-direction odometer movement is not visibly affected.
- Ensure direction reversal produces a small bounded temporary offset that is fully removed before settling exactly on the target position.
- Ensure final wheel offsets settle exactly on the target offsets.
- Keep the behaviour deterministic.
- Do not introduce random jitter.
- Keep feature state local unless an existing small helper already fits cleanly.

## Interaction with other odometer behaviour

Backlash must remain distinct from:

- `drum_slop` — static wheel alignment imperfection;
- `carry_drag` — rollover coupling between adjacent wheels;
- `snap_settle` — landing or detent behaviour;
- movement curves — timing and interpolation of ordinary movement;
- `wraparound` — route selection across digit-strip boundaries.

Backlash may compose with these behaviours, but it must not duplicate or replace them.

## Tests

Add or update tests for:

- YAML parsing accepts `realism.backlash: true` for odometers.
- YAML validation rejects `backlash` for non-odometer gauge types.
- Existing odometer movement remains unchanged when `backlash` is absent.
- Existing odometer movement remains unchanged when `backlash: false`.
- Same-direction movement is unaffected or effectively equivalent when backlash is enabled.
- Forward-to-reverse movement applies bounded slack and then catches up.
- Reverse-to-forward movement applies bounded slack and then catches up.
- Final offsets settle exactly on target offsets.
- Repeated updates do not accumulate drift.
- Behaviour remains deterministic for the same inputs and elapsed time.

## Preview and documentation

Add a small preview fixture or documented example showing:

- movement in one direction;
- a clear reversal;
- the brief backlash take-up;
- exact final settling.

Update relevant documentation so backlash is represented consistently.

## Constraints

- No hidden default behaviour.
- No randomness.
- No source-value mutation.
- No persistent numeric error.
- No unbounded delay.
- No broad physics engine.
- No renderer redesign.
- No generic realism framework.
- No support for non-odometer gauges.
- No change to ordinary movement when disabled.

## Good result

A direction reversal has a brief mechanical feel, then the odometer follows the new direction and lands exactly on the correct value.

## Bad result

The wheel delays too long, jitters, drifts, loses numeric correctness, changes ordinary same-direction movement, or behaves as though a small elephant has become lodged in the gearbox.

## Non-goals

This design does not define:

- general movement curves;
- drum alignment error;
- carry coupling;
- detent or snap behaviour;
- wraparound direction;
- random wear or vibration;
- backlash for non-odometer gauges.
