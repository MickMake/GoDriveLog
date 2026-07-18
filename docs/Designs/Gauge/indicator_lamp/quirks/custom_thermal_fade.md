# Custom indicator thermal fade quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `indicator` |
| Old realism key | `realism.thermal_fade` |
| New Gauge group | `indicator_lamp` |
| Paired custom gauge design | `docs/Designs/Gauge/indicator_lamp/gauges/custom_indicator.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Design intent

This quirk makes an indicator lamp fade in or fade out as if the visible element has thermal inertia.

For the current GoDriveLog `indicator` gauge, the behaviour applies to displayed lamp brightness only. It must not alter the input sensor value, configured thresholds, exported values, or logs.

## Physical mechanism being imitated

This quirk imitates an incandescent lamp filament warming up after power is applied and cooling after power is removed.

The on transition and off transition may feel different, because a real filament does not appear and disappear instantly.

## Expected visible behaviour

The expected visible effect is a lamp that warms into the on state and cools away from it instead of switching instantly between fully off and fully on.

## Good result

The on-state appears to warm in rather than appearing instantly. The off-state fades away softly and then settles fully off.

## Bad result

The indicator flickers, pulses, randomly changes brightness, or remains partly on after it should be off.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `indicator` renderer and is documented under the `indicator_lamp` Gauge group.

It is not a generic definition of every illuminated display or electrical failure mode. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Thermal fade is indicator-only in the current model. It should remain bounded, finite, deterministic, and display-only.

## Non-goals

This is not random flicker, bloom, lens dirt, weak bulb tint, ageing, power brownout, lamp failure, PWM scan artefacts, or a full power lifecycle model.

## Documentation boundary

This file documents the current GoDriveLog custom quirk design only.

It does not:
- rename the runtime gauge type;
- change package YAML;
- claim generic catalogue coverage;
- record implementation status;
- describe future gauge behaviour as current behaviour.

Implementation status belongs only in `docs/Status.md`.

## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/v3.5/RealismBehaviourGuide.md`
- `docs/Designs/RealismBehaviour/indicator-thermal-fade.md`
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Status.md`
