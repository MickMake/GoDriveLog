# Custom bar pointer markers quirk

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| Old realism key | `realism.pointer_markers` |
| New Gauge group | `bar_or_wedge_display` |
| Paired custom gauge design | `docs/Designs/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Documentation role | Custom current GoDriveLog quirk design |
| Runtime code impact | None |

## Naming note

This documentation uses `pointer_markers` as the current GoDriveLog realism key.

The same behaviour is also referred to as **witness markers** in older realism/design notes. Within this custom Gauge documentation set, **pointer markers** and **witness markers** are interchangeable names for the same current behaviour unless a document explicitly says otherwise.

## Design intent

This quirk displays retained marker state associated with the gauge reading, such as a minimum, maximum, follower, or tell-tale position.

For the current GoDriveLog `bar` gauge, the behaviour applies to displayed displayed fill or reveal extent only. It must not alter the input sensor value, configured ranges, exported values, or logs.

## Expected visible behaviour

The expected visible effect is one or more marker elements showing remembered positions alongside the live displayed value.

## Gauge-family boundary

This custom quirk belongs to the current GoDriveLog `bar` renderer and is documented under the `bar_or_wedge_display` Gauge group.

It is not a generic definition of every mechanical witness pointer, tell-tale, min/max register, or statistical marker. Generic physical gauge catalogue quirks remain separate from current GoDriveLog custom behaviour.

## Constraints

Pointer markers should remain deterministic and should operate on displayed state. They must not mutate source readings or replace the main value mapping.

## Non-goals

This is not `stat_markers`, automatic statistical analysis, logging summary output, or a future generic marker subsystem.

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
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Status.md`
