# Needle Trail

Applies to: radial.

Config key: `realism.needle_trail`.

## Purpose

`needle_trail` renders a bounded history of previous displayed needle positions as fading ghost needles.

It is a render-history effect, not a movement curve. `movement` controls how the live needle travels from one displayed position to another; `needle_trail` controls how recently displayed positions remain visible after the live needle has moved on.

## Real-world mechanism

Needle trail simulates persistence or afterimage around a moving pointer.

Depending on the visual treatment, it can resemble:

- optical persistence;
- a long-exposure effect;
- phosphor-like persistence;
- a stylised analogue trace.

It is deliberate visual theatre rather than a source-data effect.

## Proposed configuration

```yaml
realism:
  needle_trail:
    length: 12
    decay_ms: 500
```

| Option | Default | Meaning |
|---|---:|---|
| `length` | `12` | Maximum number of historical displayed needle positions retained. |
| `decay_ms` | `500` | Time in milliseconds for retained trail samples to fade out. |

## Source of truth

Trail samples must be derived from the final displayed radial needle position.

They must not alter or feed back into:

- source values;
- logs or exported values;
- configured ranges;
- input data;
- movement targets;
- other realism state.

## Behaviour rules

- Radial-only.
- Disabled by default.
- Display-only.
- Deterministic for the same displayed-position and timing sequence.
- Retain only a bounded history of displayed angles or positions and timestamps.
- Discard expired samples after their fade duration.
- Never retain an unbounded render history.
- Do not pre-render or cache an unbounded set of intermediate needle images.
- Reuse the normal needle geometry or asset and render transformed historical instances.
- Keep this option under `realism`, not under `movement`.

## Composition expectations

A future implementation must define exact composition with radial movement and other realism options.

Until then:

- trail history should observe the final displayed needle path;
- movement should remain responsible only for the live needle travel curve;
- trail should not modify damping, stiction, hysteresis, overshoot, peg bounce, or marker state.

## Constraints

- The live needle must remain easy to read.
- The effect must remain subtle.
- Memory use must remain bounded.
- Trail samples must expire cleanly.
- The effect must not change the underlying reading.
- Missing or invalid configuration should disable the effect safely.

## Good result

The live needle remains clear while a short sequence of fading historical positions follows behind it. The effect ends cleanly, consumes bounded memory, and never changes the underlying reading.

## Bad result

The trail obscures the dial, grows without limit, persists forever, changes the target angle, samples raw source values instead of displayed position, or becomes a pre-rendered frame factory with delusions of grandeur.
