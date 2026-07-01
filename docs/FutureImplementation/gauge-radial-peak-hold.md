# Radial Peak Hold

Index: 4

Status: desired

Area: `gauge/radial`, display marker, state tracking

Effort: 3-6 Codex hours

Add optional radial-only `realism.peak_hold` support.

Peak hold displays a secondary marker or needle at the highest displayed value reached. It is an instrument display feature, not a source value change.

## Proposed config shape

```yaml
realism:
  peak_hold:
    hold_ms: 0
    decay_ms: 1000
```

## Options

- `hold_ms`: how long to hold the peak after the displayed needle stops increasing. `0` means hold indefinitely.
- `decay_ms`: optional time for the peak marker to release/return after the hold expires.

## Rules

- Radial-only.
- Disabled by default.
- Display-only.
- Must not mutate source values, logs, exported values, configured ranges, or input data.
- Peak tracking should use displayed value/angle semantics defined by the later implementation prompt.
- If decay is enabled, release should be bounded and deterministic.
- Do not place this under `movement`; `movement` selects the travel curve, while `peak_hold` is a display marker/history feature.

## Possible future slice

```text
v3.5.20 radial peak hold
```
