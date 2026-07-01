# Gauge Imperfections

Index: 14

Status: desired

Area: `gauge/radial`, `gauge/bar`, indicator, numeric, display artefacts, mechanical wear, electrical artefacts, gauge realism

Effort: 7-12 Codex hours

Add optional gauge-level `realism.imperfections` support for controlled, deterministic, display-only gauge ageing, wear, vibration, electrical noise, and display artefacts.

This is not a bug mode. It is a realism layer for gauges that should look like real physical instruments: old, worn, cable-driven, electrically noisy, cheaply multiplexed, or slightly temperamental.

The goal is character without corrupting the data.

## Concept

A clean gauge shows the mapped source value as accurately as the configured gauge allows.

An imperfect gauge still uses the real source value, but may render the visible presentation with bounded mechanical or visual artefacts.

Examples:

```text
idle needle vibration
worn cable speedo wobble
brownout dip during engine start
LED multiplex flicker
gas-discharge jitter
intermittent display flicker
```

All imperfections are display-only unless a later implementation explicitly says otherwise.

## Proposed config shape

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

Use this as a direction, not a locked schema.

## Imperfection categories

### Mechanical imperfections

Mechanical imperfections affect displayed position, motion, or pointer stability.

Examples:

```text
idle needle vibration
worn cable speedo wobble
sticky pivot
loose needle
end-stop rattle
```

Mechanical imperfections mostly apply to radial gauges, though some effects may later map to bar gauges.

### Electrical imperfections

Electrical imperfections simulate voltage and grounding artefacts.

Examples:

```text
brownout dip during engine start
voltage sag
needle twitch on power transition
backlight shimmer
bad-ground flicker
```

These may compose with dashboard-provided power lifecycle events, but the visible response belongs to the gauge.

### Display technology imperfections

Display technology imperfections simulate quirks of specific visual technologies.

Examples:

```text
LED multiplex flicker
gas-discharge jitter
LCD ghosting
seven-segment uneven brightness
backlight PWM shimmer
```

These are usually visual-only and should not change displayed gauge position.

### Intermittent imperfections

Intermittent imperfections are rare, bounded, deterministic faults.

Examples:

```text
brief random-looking flicker
momentary dimming
brief missing segment
short needle twitch
```

Avoid unbounded randomness. Use seeded pseudo-random timing so replay and screenshots can be reproduced.

## Position vs visual imperfections

Separate imperfections into two implementation classes.

### Position imperfections

Position imperfections change the rendered gauge position but not the source value.

Examples:

```text
idle needle vibration
eddy-current speedo wobble
brownout needle dip
sticky needle
```

These effects may alter the rendered needle angle, bar extent, or displayed marker position.

They must not alter source values, logs, exports, configured ranges, or stat marker samples.

### Visual-only imperfections

Visual-only imperfections change appearance without changing displayed value position.

Examples:

```text
LED multiplex flicker
gas-discharge jitter
backlight shimmer
brief display flicker
uneven segment brightness
```

These effects may alter opacity, brightness, asset selection, segment visibility, or glow.

They must not alter the displayed value, source values, logs, exports, configured ranges, or stat marker samples.

## Eddy-current speedo wobble

The Torana-style speedo fault should be treated as a named imperfection, not generic random jitter.

Old cable-driven speedometers often use an eddy-current mechanism. The vehicle turns a speedo cable. The cable spins a magnet or magnetic disc inside the gauge. The spinning magnetic field applies drag torque to a nearby cup or disc connected to the needle. A return spring resists that torque. The faster the cable spins, the more torque is applied and the higher the needle reads.

If a magnet is missing, weak, or uneven, the magnetic drag is no longer smooth across each rotation. The needle sees a lumpy pull pattern instead of a smooth average force.

That produces a wobble that is still related to speed, but not cleanly centred around the true speed.

Important behaviour:

```text
slower speed -> larger visible oscillation
faster speed -> smaller visible oscillation
asymmetric wobble, not equal +/- movement
usually biased toward under-reading or over-reading
bounded around the real speed
```

This should model an uneven magnetic force curve, not arbitrary random noise.

Suggested config direction:

```yaml
realism:
  imperfections:
    mechanical:
      eddy_speedo_wobble:
        enabled: true
        low_speed_wobble: 14
        high_speed_wobble: 2
        asymmetry: 70/30
        bias: under_read
        damping: 0.6
        speed_sensitive: true
```

Meaning:

- `low_speed_wobble` controls the large slow-speed visible swing.
- `high_speed_wobble` controls the smaller high-speed visible swing.
- `asymmetry` describes the uneven force shape, such as mostly normal pull with a weaker phase.
- `bias` controls whether the wobble tends to dip below or surge above the real value.
- `damping` controls how much the needle smooths the uneven force pulses.
- `speed_sensitive` means oscillation frequency and visible amplitude are driven by displayed speed.

This effect should initially be radial-only and most appropriate for speedometer-style gauges.

## Idle needle vibration

Idle needle vibration simulates small engine/chassis vibration affecting a pointer at idle or low speed.

This is different from speedo wobble.

```text
idle vibration = small high-frequency shake
speedo wobble = speed-linked uneven magnetic/cable-drive oscillation
```

Suggested config direction:

```yaml
realism:
  imperfections:
    mechanical:
      idle_needle_vibration:
        enabled: true
        amplitude: 0.8
        frequency_hz: 18
        only_when_idle: true
```

Keep this subtle. If the viewer notices the effect before noticing the gauge, it is too strong.

## Brownout dip

Brownout dip simulates voltage sag while starting the engine or during an electrical load event.

It may affect:

```text
needle position
bar brightness
numeric display brightness
indicator brightness
backlight brightness
```

Suggested config direction:

```yaml
realism:
  imperfections:
    electrical:
      brownout_dip:
        enabled: true
        trigger: engine_start
        duration_ms: 450
        brightness_drop: 0.35
        needle_drop: 3
```

Brownout should compose with power lifecycle, but should not be identical to power on/off behaviour.

## LED multiplex flicker

LED multiplex flicker simulates the slight shimmer or scan artefact from multiplexed LED or seven-segment displays.

Suggested config direction:

```yaml
realism:
  imperfections:
    display:
      led_multiplex_flicker:
        enabled: true
        intensity: subtle
        frequency_hz: 60
```

This should usually be visual-only.

## Gas-discharge jitter

Gas-discharge jitter simulates unstable glow, uneven brightness, or flickering discharge-style displays.

Suggested config direction:

```yaml
realism:
  imperfections:
    display:
      gas_discharge_jitter:
        enabled: true
        intensity: low
        seed: gauge
```

This should usually be visual-only.

## Intermittent flicker

Intermittent flicker simulates occasional minor faults such as bad contacts, dirty connectors, or ageing electronics.

This is the easiest one to overdo.

Suggested config direction:

```yaml
realism:
  imperfections:
    intermittent:
      flicker:
        enabled: true
        probability_per_minute: 0.2
        duration_ms: [40, 120]
        intensity: subtle
        seed: gauge
```

Rules:

- rare by default;
- bounded duration;
- deterministic seed;
- no unbounded random flashing;
- no effect on source values.

## Presets

Presets may provide convenient imperfection bundles.

Possible presets:

```text
clean
old_analogue
worn_cable_speedo
cheap_led
weak_battery
bad_ground
gas_discharge
race_car_vibration
```

A preset should expand to ordinary config values. Manual config should still be able to override individual effects.

Example:

```yaml
realism:
  imperfections:
    preset: worn_cable_speedo
```

Could expand to something like:

```yaml
realism:
  imperfections:
    mechanical:
      eddy_speedo_wobble:
        enabled: true
        low_speed_wobble: 14
        high_speed_wobble: 2
        asymmetry: 70/30
        bias: under_read
        damping: 0.6
      idle_needle_vibration:
        enabled: true
        amplitude: 0.8
```

## Composition rules

Imperfections are downstream of source values and value mapping.

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

Some imperfections are position perturbations. Some are visual-only. Keep that distinction explicit.

Imperfections should compose with:

```text
power lifecycle
lighting mode
stat markers
movement realism
asset selection
replay
```

## Stat marker rules

By default, stat markers should ignore imperfection perturbations.

Stat markers should track the stable displayed value before imperfection effects unless a later feature explicitly adds an option to include imperfections.

Reason: a wobbly needle is a presentation defect, not a real source value.

## Replay and determinism

All pseudo-random imperfections must be deterministic under replay.

Use stable seeds, such as:

```text
gauge id
session id
explicit configured seed
```

Do not use uncontrolled runtime randomness for visible gauge behaviour.

## Do not

- Do not mutate source values, logs, exports, configured ranges, stat marker samples, odometer source values, or sensor state.
- Do not make everything randomly twitch.
- Do not let random flicker run unbounded.
- Do not make imperfections active by default.
- Do not apply speedometer-specific faults to all radial gauges without explicit configuration.
- Do not let effects hide the actual useful value unless the user deliberately configures a strong fault.

## Possible future slices

```text
gauge imperfections config model
radial idle needle vibration
eddy-current speedo wobble
brownout dip display effect
LED multiplex flicker
gas-discharge jitter
intermittent deterministic flicker
imperfection presets
```
