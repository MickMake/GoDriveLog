# Future Slices

This file is a parking lot for approved or desired follow-on ideas that are not part of the current implementation slice.

Use this to capture "oh, also implement this later" notes without making the active slice ambiguous. Future prompts may reference this file, but items here are not current scope unless a later prompt explicitly promotes them.

## Guidelines

- Keep entries small and slice-shaped.
- Mark ideas as `deferred`, `desired`, `exploratory`, or `rejected`.
- Do not treat this file as an implementation checklist.
- Do not let vague mentions here expand the current slice.
- Prefer a later dedicated prompt/spec before implementation.

## Bar gauge overshoot

Status: deferred

Bar gauges should eventually support `realism.overshoot`, but this was intentionally left out of the v3.5.10 radial overshoot slice.

Notes:

- Display-only.
- Bounded pass-and-settle movement.
- Should compose cleanly with bar damping/smoothing.
- Do not copy radial behaviour blindly; bar movement has its own visual semantics.
- Keep radial overshoot behaviour unchanged when this is implemented.

Possible future slice:

```text
v3.5.x bar overshoot
```

## Value zones / warning-danger assets

Status: desired

Support optional value zones that select warning/danger variants of gauge assets when the source value reaches a configured range.

This should be a separate gauge-display feature, not part of `realism.overshoot`.

Proposed config shape:

```yaml
zones:
  warning:
    min: 6000
    max: 7000
  danger:
    min: 7000
    max: 8000
```

Proposed asset convention:

```text
needle.png
needle_warning.png
needle_danger.png
face.png
face_warning.png
face_danger.png
bar.png
bar_warning.png
bar_danger.png
```

Rules:

- Zone selection should follow the real/source target value, not any temporary animated display value.
- If a zone-specific asset exists for a layer, use it.
- If a zone-specific asset does not exist, fall back to the normal asset.
- Overshoot may visually pass a threshold, but should not change the zone state unless the real/source value is in that zone.
- Avoid surprising behaviour where a temporary animation makes the gauge appear to enter warning or danger falsely.

Possible future slice:

```text
v3.5.x value zones / warning-danger assets
```
