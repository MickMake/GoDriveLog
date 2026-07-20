# Custom radial pointer markers quirk implementation

## Identity

| Field | Value |
|---|---|
| Old GoDriveLog type | `radial` |
| Old realism key | `realism.pointer_markers` |
| New Gauge group | `radial_pointer` |
| Paired custom quirk design | `docs/Designs/Gauge/radial_pointer/quirks/custom_pointer_markers.md` |
| Paired custom gauge implementation | `docs/Implementation/Gauge/radial_pointer/gauges/custom_radial.md` |
| Runtime code impact | None |

## Naming note

This documentation uses `pointer_markers` as the current GoDriveLog realism key.

The same behaviour is also referred to as **witness markers** in older realism/design notes. Within this custom Gauge documentation set, **pointer markers** and **witness markers** are interchangeable names for the same current behaviour unless a document explicitly says otherwise.

## Current implementation model

Current code treats `realism.pointer_markers` as an implemented current-code realism key.

For the current GoDriveLog `radial` gauge, the behaviour applies to rendered marker state associated with the displayed pointer position or angle. It must not change source sensor values, persisted log output, exported values, or configured range semantics.

## Configuration boundary

The old GoDriveLog realism key remains `realism.pointer_markers`.

This document does not rename that key and does not introduce a new Gauge-tree runtime configuration name.

## Current limitations and exclusions

Pointer markers are not `stat_markers`. Current status records `pointer_markers` as implemented and `stat_markers` as not implemented. Do not merge those behaviours without a later audited code change.

## Documentation boundary

This file records current GoDriveLog custom quirk implementation behaviour only.

It does not:
- record implementation status;
- describe future gauge work as implemented;
- rename runtime package types;
- replace or migrate existing documentation.

Implementation status belongs only in `docs/Status.md`.


## Historical source basis

- `docs/v3.5/ImplementationState.md`
- `docs/v3.5/RealismBehaviourGuide.md`
- `docs/Status.md`
