# GoDriveLog v3.2

Planning documents for the v3.2 dashboard gauge/instrument package series.

v3.2 benches the remaining v3.1.7 dashboard event efficiency and v3.1.8 retirement readiness work temporarily so the dashboard instrument package direction can be advanced while the design is clear.

## Documents

- `ReleasePlan.md` - implementation roadmap and version queue.
- `ImplementationState.md` - current v3.2 state and next target.
- `OpenDecisions.md` - unresolved decisions with blocking and impact notes.
- `CarryForward.md` - v3.1 lessons and constraints that still matter.
- `prompts/README.md` - prompt index and common implementation guardrails.
- `prompts/v3.2.x-*.md` - one prompt per implementation slice.

## Direction

v3.2 introduces self-contained gauge packages under `assets/gauges/**/gauge.yaml`.

A gauge package is a complete dashboard instrument package. It owns its sensor binding, visual definition, formatting, value mapping, layout geometry, renderer type, and local asset references.

A dashboard `type: gauge` widget places a gauge package on the dashboard using a gauge path, position, and scale.

The first package target is a seven-segment display package because it builds on the existing digit display path while moving panel, glass, digit positions, digit count, and `format` into the package. Radial gauges come after that.

Directory names under `assets/gauges/` are user-defined and carry no renderer meaning. The only required filename is `gauge.yaml`; the gauge type is declared inside that file.

Tiny app, fewer steering wheels, fewer tiny steering-wheel accessories.
