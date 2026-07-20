# Custom bar stiction quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| Old realism key | `realism.stiction` |
| New Gauge group | `bar_or_wedge_display` |
| Paired custom quirk design | `docs/Designs/Gauge/bar_or_wedge_display/quirks/custom_stiction.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Runtime code impact | None |

## Current implementation model

Current code treats `realism.stiction` as an implemented realism option for radial and bar gauges.

For the current `bar` renderer, the quirk is applied to the displayed state of the fill or reveal extent. It is display-only and must not mutate the underlying sensor reading, configured range, logged value, or exported value.

## Current rendering effect

the fill or reveal extent may hold through small changes, then release to a new displayed level.

## Configuration boundary

The current user-facing realism key remains:

```yaml
realism:
  stiction: ...
```

This documentation does not rename the runtime key and does not introduce a Gauge-directory-specific configuration name.

## Current limitations and boundaries

This implementation record is intentionally narrow. It documents the existing `bar` support for `realism.stiction` only.

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
