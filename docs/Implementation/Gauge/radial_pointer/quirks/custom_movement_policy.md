# Custom radial movement policy quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.movement_policy` |
| New Gauge group | `radial_pointer` |
| Paired custom quirk design | `docs/Designs/Gauge/radial_pointer/quirks/custom_movement_policy.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/radial_pointer/gauges/custom_radial.md` |
| Runtime code impact | None |

## Implementation status

**Partially implemented.**

Current code supports part of the intended radial movement-policy contract, but not the full documented policy set.

## Naming note

This document uses `movement_policy` because current code uses `Realism.MovementPolicy` for `radial` movement behaviour.

Do not document this as odometer `movement`. Odometer movement has a separate current configuration surface and a separate wheel movement model.

## Current implementation model

Current code treats `radial` movement as a realism policy applied to the displayed pointer angle.

The current audited policy values are:

| Value | Current code support | Notes |
|---|---|---|
| `immediate` | Implemented | Normalized default. |
| `linear` | Implemented | Used for finite displayed pointer transitions. |
| `ease_out` | Implemented | Used as a non-linear transition policy. |
| `bell` | Not implemented for radial `movement_policy` | Desired extension. Odometer has a bell-style movement curve, but radial movement policy does not currently accept or apply it. |

The behaviour applies to rendered state only. It must not change source sensor values, persisted log output, exported values, or configured range semantics.

## Configuration boundary

The old GoDriveLog realism key remains:

```yaml
realism:
  movement_policy: immediate
```

This document does not rename that key and does not introduce a new Gauge-tree runtime configuration name.

`bell` should be added to this existing `realism.movement_policy` contract if implemented. It should not introduce a separate scalar radial `movement` key.

## Current limitations and exclusions

This is only the current movement-policy surface. It does not imply a nested movement phase model, physics simulation, continuous idle animation, or support for odometer-only movement values.

Current radial movement policy is also constrained by runtime movement gating: visible movement may depend on another timed movement behaviour being active. A future `bell` implementation must decide whether explicit non-immediate policies should create their own bounded default transition duration for radial gauges.

## Remaining work for `bell`

To implement `bell` for radial movement policy, code would need to:

- add a `bell` movement-policy constant or otherwise accept `bell` as a valid `realism.movement_policy` value;
- validate `bell` for radial gauges;
- apply a bell/smoothstep-style curve to radial displayed pointer movement;
- prevent explicit `bell` from being forced back to `immediate` when no other movement quirk is active;
- provide a bounded default transition duration, or explicitly document that `bell` only takes effect when another timed movement behaviour is active;
- add tests for config parsing, validation, curve behaviour, and final exact settling.

## Documentation boundary

This file records current GoDriveLog custom quirk implementation behaviour only.

It does not:
- describe `bell` as currently implemented for radial gauges;
- rename runtime package types;
- replace or migrate existing documentation;
- replace `docs/Status.md` as the authoritative implementation-status record.

Implementation status belongs in `docs/Status.md`.

## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/v3.5/RealismBehaviourGuide.md`
- `docs/Status.md`
