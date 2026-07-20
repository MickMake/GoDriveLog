# Gauge Imperfections

Index: 14

Status: desired implementation backlog

Area: `gauge/radial`, `gauge/bar`, indicator, numeric, display artefacts, mechanical wear, electrical artefacts, gauge realism

Effort: 7-12 Codex hours

## Implementation goal

Add optional gauge-level `realism.imperfections` support for controlled, deterministic, display-only gauge ageing, wear, vibration, electrical noise, and display artefacts.

This is not a bug mode. It is a realism layer for gauges that should look like real physical instruments: old, worn, cable-driven, electrically noisy, cheaply multiplexed, or slightly temperamental.

The goal is character without corrupting the data.

## Proposed config direction

Use this as a direction, not a locked schema:

```yaml
gauges:
  speed:
    realism:
      imperfections:
        preset: worn_cable_speedo

        mechanical:
          idle_needle_vibration:
            enabled: true
            amplitude: 0.8
            frequency_hz: 18

          eddy_speedo_wobble:
            enabled: true
            low_speed_wobble: 14
            high_speed_wobble: 2
            asymmetry: 70/30
            bias: under_read
            damping: 0.6

        electrical:
          brownout_dip:
            enabled: false

        display:
          led_multiplex_flicker:
            enabled: false
          gas_discharge_jitter:
            enabled: false

        intermittent:
          flicker:
            enabled: false
            seed: gauge
```

## Implementation planning notes

- Imperfections are downstream of source values and value mapping.
- Some imperfections are position perturbations; some are visual-only. Keep that distinction explicit.
- Position imperfections may alter rendered needle angle, bar extent, or displayed marker position.
- Visual-only imperfections may alter opacity, brightness, asset selection, segment visibility, or glow.
- Imperfections must not alter source values, logs, exports, configured ranges, sensor state, odometer source values, or persisted data.
- All pseudo-random imperfections must be deterministic under replay.
- Use stable seeds such as gauge id, session id, or an explicit configured seed.
- Do not make imperfections active by default.
- Do not apply speedometer-specific faults to all radial gauges without explicit configuration.

## Composition rules

Suggested order:

```text
source value
value mapping
stable display value
power / lighting lifecycle
movement realism
imperfection layer
render
```

Imperfections should compose with:

```text
power lifecycle
lighting mode
pointer markers
movement realism
asset selection
replay
```

## Pointer marker rule

By default, pointer markers should ignore imperfection perturbations.

Pointer markers should track the stable displayed value before imperfection effects unless a later feature explicitly adds an option to include imperfections.

Reason: a wobbly needle is a presentation defect, not a real source value.

## Suggested future implementation tickets

- Specify the `realism.imperfections` config model.
- Implement radial `idle_needle_vibration`.
- Implement eddy-current speedometer wobble.
- Implement brownout dip display effects.
- Implement LED multiplex flicker.
- Implement gas-discharge jitter.
- Implement intermittent deterministic flicker.
- Implement imperfection presets.

## Do not

- Do not mutate source values, logs, exports, configured ranges, pointer marker samples, odometer source values, or sensor state.
- Do not make everything randomly twitch.
- Do not let random flicker run unbounded.
- Do not make imperfections active by default.
- Do not let effects hide the actual useful value unless the user deliberately configures a strong fault.
