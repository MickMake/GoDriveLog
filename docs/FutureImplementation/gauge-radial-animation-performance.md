# Radial Animation Performance

Index: 11

Status: desired

Area: `gauge/radial`, animation loop, renderer, low-power display performance

Effort: 3-6 Codex hours

Improve the reliability of short radial gauge animations on lower-powered display targets.

The current radial bounce and overshoot effects can look visually correct when rendered smoothly, but very short animation tails may occasionally appear to skip or vanish under load. This is most noticeable for subtle effects where the important motion happens over only a few rendered frames.

This is not a request to make the movement stronger, larger, or more dramatic. The visual result should remain essentially unchanged. The goal is to make existing subtle radial movement render more consistently.

## Areas to investigate

- animation tick cadence;
- render/update scheduling;
- minimum visible animation sample count;
- short movement tail duration;
- low-power display target load;
- asset size and rotation cost;
- whether radial animation needs fixed-timestep sampling independent of frame delivery.

## Rules

- Preserve existing visual semantics for radial overshoot, peg bounce, damping, and stiction.
- Do not increase rebound or overshoot amplitude as a performance fix.
- Do not change source sensor values, logs, exported values, or configured ranges.
- Prefer timing/sampling/render-loop improvements over changing gauge behaviour.
- Keep this radial-focused unless profiling proves another gauge type is affected.
- Any fix should be validated on a low-power display target, not only on desktop.
- Treat occasional dropped frames as a rendering/performance problem, not as a realism tuning problem.

## Possible future slice

```text
v3.x radial animation performance
```
