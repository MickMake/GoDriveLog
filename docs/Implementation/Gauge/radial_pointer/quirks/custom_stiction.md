# Custom radial stiction quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.stiction` |
| New Gauge group | `radial_pointer` |
| Paired custom quirk design | `docs/Designs/Gauge/radial_pointer/quirks/custom_stiction.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/radial_pointer/gauges/custom_radial.md` |
| Runtime code impact | None |

## Current implementation model

Current code treats `realism.stiction` as an implemented realism option for radial and bar gauges.

For the current `radial` renderer, the quirk is applied to the displayed state of the needle. It is display-only and must not mutate the underlying sensor reading, configured range, logged value, or exported value.

## Current rendering effect

the needle may hold briefly through small changes, then release to a new displayed angle.

## Configuration boundary

The current user-facing realism key remains:

```yaml
realism:
  stiction: ...
```

This documentation does not rename the runtime key and does not introduce a Gauge-directory-specific configuration name.

## Current limitations and boundaries

This implementation record is intentionally narrow. It documents the existing `radial` support for `realism.stiction` only.

It does not claim:

- random sticking;
- permanent jam simulation;
- sensor fault simulation;
- changing the source sensor value;


## Documentation boundary

This file records current GoDriveLog implementation behaviour only.

It does not:

- record implementation status;
- rename runtime gauge types;
- describe unimplemented or future realism behaviour as current behaviour.

Implementation status belongs only in `docs/Status.md`.


## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/Status.md`
