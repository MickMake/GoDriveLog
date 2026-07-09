# Future Implementation

This directory is an intent register for things we may want to implement later.

It is deliberately lightweight. It should not contain architecture, behaviour definitions, full specs, or implementation checklists.

## Rule

```text
RealismBehaviourGuide = definition / behaviour / real-world simulation
FutureImplementation = intent to implement later
```

Use this file to say:

- what we might implement;
- why we care;
- where the canonical behaviour/design definition lives;
- whether it is near, later, implemented, or only a vague candidate.

Do not use this directory to define detailed architecture. When an item is promoted, write a dedicated release plan, prompt, issue, or PR spec at that time.

## Intent register

| Intent | Status | Canonical / detail source |
|---|---|---|
| Bar gauge overshoot follow-up | implemented in v3.5.19 | [`../RealismBehaviourGuide/overshoot.md`](../RealismBehaviourGuide/overshoot.md) |
| Radial movement options | near / needs spec tightening | [`../RealismBehaviourGuide/movement.md`](../RealismBehaviourGuide/movement.md) |
| Radial needle trail | later / visual polish | future dedicated spec needed |
| Gauge pointer/stat markers | implemented in v3.6 | [`../RealismBehaviourGuide/pointer-markers.md`](../RealismBehaviourGuide/pointer-markers.md) |
| Value zones / warning-danger assets | useful soon | future dedicated spec needed |
| Canonical GoDriveLog event log | medium / foundational | future dedicated spec needed |
| Session metadata sidecar | medium / pairs with event log | future dedicated spec needed |
| JSONL dashboard replay | medium / high-value dev tool | future dedicated spec needed |
| JSONL log validation | near / bounded utility | future dedicated spec needed |
| External converter boundary | later / architecture boundary | future dedicated spec needed |
| Needle animation performance | near / performance polish | future dedicated spec needed |
| Gauge power lifecycle | later / gauge realism | future dedicated spec needed |
| Gauge lighting mode | later / gauge realism | future dedicated spec needed |
| Gauge imperfections | later / gauge realism | [`../RealismBehaviourGuide/imperfections.md`](../RealismBehaviourGuide/imperfections.md) |
| Gauge presets | later / config reuse | future dedicated spec needed |
| Odometer backlash | candidate / not implemented | [`../RealismBehaviourGuide/backlash.md`](../RealismBehaviourGuide/backlash.md) |
| Numeric/segmented display realism | candidate / not implemented | [`../RealismBehaviourGuide/`](../RealismBehaviourGuide/) |
| Bar stepped/quantized fill | candidate / not implemented | [`../RealismBehaviourGuide/stepped-fill.md`](../RealismBehaviourGuide/stepped-fill.md), [`../RealismBehaviourGuide/quantized-fill.md`](../RealismBehaviourGuide/quantized-fill.md) |

## Promotion rule

Before implementing anything from this file:

1. choose one small intent;
2. check whether a canonical definition already exists;
3. create or update the canonical definition if needed;
4. write a dedicated implementation prompt/spec;
5. implement only that promoted slice.
