# Radial Animation Performance

Index: 11

Status: partially implemented / mostly superseded by v3.2.6 and v3.3 renderer work

Area: dashboard renderer, radial gauge rendering, animation loop, scene delivery, low-power display performance

Effort: historical audit / future targeted fix only if the issue reappears

Improve the reliability of short radial gauge animations on lower-powered display targets.

The original issue was that small radial movements could look jerky, skipped, or absent under load. This was most noticeable for subtle radial effects where the important motion happens over only a few rendered frames.

This was not a request to make radial movement stronger, larger, or more dramatic. The intended visual result should remain essentially unchanged. The goal was to make existing subtle radial movement render more consistently.

## Current implementation state

This concern has been **partially implemented** and **mostly superseded** by later renderer work.

Repository evidence shows that the project implemented general renderer/delivery improvements rather than a dedicated radial-only animation sampler:

- v3.2.6 added deterministic one-degree prepared radial needle frames.
- v3.2.6 reused keyed rendered image resources during normal live radial updates.
- v3.2.6 added latest-only scene coalescing via `LatestSink.SubmitLatest`.
- v3.2.6 added display sink stats for submitted, rendered, superseded, and render-duration values.
- v3.3 moved the active dashboard renderer path to Ebiten after Fyne proved unsuitable on Raspberry Pi baseline hardware.
- v3.3 recorded target-hardware evidence that the Ebiten path was smooth and responsive on the baseline dashboard workload.

These changes are relevant and useful, but they do not exactly implement the original design as written.

## What was implemented

Implemented or evidenced:

- prepared radial needle frame sets;
- renderer object/resource reuse for radial updates;
- latest-frame dashboard scene coalescing;
- non-blocking display submission for runtime/harness paths;
- display rendering statistics;
- Ebiten renderer path as the active v3 dashboard renderer;
- Raspberry Pi baseline performance evidence for the active Ebiten path.

## What was not proven implemented

Not clearly implemented as a dedicated feature:

- radial-only fixed timestep animation sampling;
- minimum visible radial animation sample count;
- guaranteed visibility of every short radial movement tail;
- radial-specific animation timing/tail extension;
- radial-animation-specific config or tuning surface.

## Design interpretation

Treat this as a dashboard renderer/performance concern, not as a new radial gauge mechanism.

The original radial-specific symptom was real, but the practical fix appears to have been addressed mostly through renderer architecture and scene delivery changes rather than by adding a radial-only animation subsystem.

## Areas originally investigated

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
- Prefer timing, sampling, renderer, or scene-delivery improvements over changing gauge behaviour.
- Keep this radial-focused unless profiling proves another gauge type is affected.
- Any future fix should be validated on a low-power display target, not only on desktop.
- Treat occasional dropped frames as a rendering/performance problem, not as realism tuning.
- Do not introduce a radial-only animation subsystem unless the current renderer path still shows the problem on target hardware.

## Future work trigger

Do not reopen this as active work unless subtle radial movement is still visibly jerky on the current Ebiten dashboard path.

If it does reappear, the next investigation should verify:

```text
OBD / harness source
-> prepared vehicle/sensor data
-> runtime event path
-> dashboard scene generation
-> display sink / latest submission
-> renderer adapter
-> screen
```

The fix should be based on target-hardware profiling and should preserve existing radial realism behaviour.

## Historical source basis

- `docs/Designs/RealismBehaviour/radial-animation-performance.md`
- `docs/Implementation/RealismBehaviour/radial-animation-performance.md`
- `docs/v3.2/ImplementationState.md`
- `docs/v3.3/ImplementationState.md`
- `docs/v3.3/PerformanceRuns.md`
- `docs/v3.5/ImplementationState.md`
