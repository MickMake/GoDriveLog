# Custom bar peg bounce quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `bar` |
| Old realism key | `realism.peg_bounce` |
| New Gauge group | `bar_or_wedge_display` |
| Paired custom quirk design | `docs/Designs/Gauge/bar_or_wedge_display/quirks/custom_peg_bounce.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/bar_or_wedge_display/gauges/custom_bar.md` |
| Runtime code impact | None |

## Current implementation model

Current code treats `realism.peg_bounce` as an implemented realism option for radial and bar gauges.

For the current `bar` renderer, the quirk is applied to the displayed state of the fill or reveal extent. It is display-only and must not mutate the underlying sensor reading, configured range, logged value, or exported value.

## Current rendering effect

the fill or reveal extent can rebound after contacting its minimum or maximum end of travel.

## Configuration boundary

The current user-facing realism key remains:

```yaml
realism:
  peg_bounce: ...
```

This documentation does not rename the runtime key and does not introduce a Gauge-directory-specific configuration name.

## Current limitations and boundaries

This implementation record is intentionally narrow. It documents the existing `bar` support for `realism.peg_bounce` only.

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
- describe unimplemented or future realism behaviour as current behaviour.

Implementation status belongs only in `docs/Status.md`.


## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/Status.md`
