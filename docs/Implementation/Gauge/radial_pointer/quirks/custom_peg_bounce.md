# Custom radial peg bounce quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.peg_bounce` |
| New Gauge group | `radial_pointer` |
| Paired custom quirk design | `docs/Designs/Gauge/radial_pointer/quirks/custom_peg_bounce.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/radial_pointer/gauges/custom_radial.md` |
| Runtime code impact | None |

## Current implementation model

Current code treats `realism.peg_bounce` as an implemented realism option for radial and bar gauges.

For the current `radial` renderer, the quirk is applied to the displayed state of the needle. It is display-only and must not mutate the underlying sensor reading, configured range, logged value, or exported value.

## Current rendering effect

the needle can rebound after contacting the minimum or maximum stop.

## Configuration boundary

The current user-facing realism key remains:

```yaml
realism:
  peg_bounce: ...
```

This documentation does not rename the runtime key and does not introduce a Gauge-directory-specific configuration name.

## Current limitations and boundaries

This implementation record is intentionally narrow. It documents the existing `radial` support for `realism.peg_bounce` only.

It does not claim:

- bounce away from non-stop positions;
- random jitter;
- unbounded oscillation;
- changing configured min/max values;


## Documentation boundary

This file records current GoDriveLog implementation behaviour only.

It does not:

- record implementation status;
- rename runtime gauge types;
- replace existing `RealismBehaviour` records;
- describe unimplemented or future realism behaviour as current behaviour.

Implementation status belongs only in `docs/Status.md`.


## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/Designs/RealismBehaviour/realism-behaviour-guide.md`
- `docs/Status.md`
- Existing `docs/Designs/RealismBehaviour/*` and `docs/Implementation/RealismBehaviour/*` records where present

