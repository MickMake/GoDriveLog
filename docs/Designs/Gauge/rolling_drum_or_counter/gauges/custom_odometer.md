# Custom odometer gauge

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `odometer` |
| New Gauge group | `rolling_drum_or_counter` |
| Documentation role | Custom current GoDriveLog gauge design |
| Runtime code impact | None |

## Design intent

The GoDriveLog `odometer` gauge is a transform gauge. It displays a numeric value using rolling digit or wheel-strip assets clipped through fixed windows.

The type models a rolling counter display rather than a plain formatted text readout.

## Behaviour model

An `odometer` gauge:

- receives a numeric value;
- decomposes that value into wheel positions;
- maps each wheel to a strip offset;
- clips the visible portion of each strip through a configured window;
- draws the result as a set of rolling digit wheels.

## Movement model

The current design distinguishes two public movement modes:

| Movement | Meaning |
|---|---|
| `smooth` | Continuous strip offset between digit positions. |
| `click` | Stepped movement that snaps to digit positions. |

## Asset model

Each wheel is backed by a strip asset. The artwork, digit shape, ageing, tint, and decorative casing belong to assets rather than gauge type names.

## Current design boundaries

The GoDriveLog odometer model is a flat strip/window renderer. It is not a full physical gear train, curved drum, or mechanical counter simulator.

## Not current behaviour

Do not treat backlash, gear lash, advanced easing, inertia, curved depth, or rear-wheel wraparound as part of this current gauge design unless separately documented and verified.

## Documentation boundary

This file documents the current GoDriveLog custom gauge design only.

It does not:
- rename the runtime gauge type;
- change package YAML;
- claim generic catalogue coverage;
- record implementation status;
- describe future renderer work as current behaviour.

Implementation status belongs only in `docs/Status.md`.

## Historical source basis

- `docs/v3.4/ReleasePlan.md`
- `docs/v3.4/ImplementationState.md`
- `docs/v3.4/prompts/v3.4.0-gauge-type-docs.md`

