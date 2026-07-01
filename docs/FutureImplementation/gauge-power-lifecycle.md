# Gauge Power Lifecycle

Index: 12

Status: desired

Area: `gauge/radial`, `gauge/bar`, indicator, numeric, odometer, dashboard runtime, power-state events, gauge realism

Effort: 6-10 Codex hours

Add gauge-level power-on and power-off realism driven by a dashboard power-state signal, such as ACC on/off.

The dashboard runtime detects the external power state. Each gauge owns how it reacts to that state.

This keeps gauges self-contained: a realistic gauge should know how it wakes, settles, powers down, blanks, drops, sweeps, or holds when vehicle accessory power changes.

## Boundary

Separate the feature into two layers:

```text
ACC / GPIO input -> dashboard power signal provider -> runtime power events -> gauge realism.power behaviour
```

The dashboard/runtime layer is responsible for detecting and debouncing a power-state signal.

The gauge layer is responsible for presentation behaviour when the runtime broadcasts power events.

Do not wire GPIO details directly into gauge rendering.

## Electrical note

A Raspberry Pi GPIO input should not be treated as a raw automotive input.

The ACC signal should be electrically isolated and level-shifted before the Pi sees it. The intended software design should allow a GPIO-backed provider, but the implementation docs should still warn that vehicle 12V ACC is noisy and must not be connected directly to a Pi GPIO.

## Proposed dashboard input config

```yaml
dashboard:
  power:
    input:
      type: gpio
      pin: 17
      active: high
      debounce_ms: 100
      stable_ms: 250
```

Supported input provider types should eventually include:

```text
gpio
manual
mock
none
```

`manual` and `mock` are useful for desktop preview, tests, and CI where GPIO is not available.

## Proposed gauge config

Gauge behaviour should live under gauge realism:

```yaml
gauges:
  rpm:
    realism:
      power:
        on:
          behaviour: needle_sweep
          sweep_to: max
          duration_ms: 1200
          settle: live_value

        off:
          behaviour: slow_drop
          target: min
          duration_ms: 1600
          final: rest
```

Global defaults may exist, but per-gauge config should win.

```yaml
dashboard:
  power:
    defaults:
      on:
        behaviour: none
      off:
        behaviour: none
```

Use dashboard defaults only as a convenience layer. The actual visual response belongs to the gauge.

## Runtime power states

The runtime should expose a small power lifecycle state machine:

```text
off
powering_on
live
powering_off
```

Expected transitions:

```text
off + ACC on -> powering_on
powering_on complete -> live
live + ACC off -> powering_off
powering_off complete -> off
powering_on + ACC off -> powering_off
powering_off + ACC on -> powering_on
```

The last two transitions matter because ACC can change quickly during key-position changes or vehicle startup.

## Gauge-family examples

### Radial gauges

Power on:

```yaml
realism:
  power:
    on:
      behaviour: needle_sweep
      sweep_to: max
      duration_ms: 1000
      settle: live_value
```

Power off:

```yaml
realism:
  power:
    off:
      behaviour: slow_drop
      target: min
      duration_ms: 1200
      final: rest
```

Radial power behaviour may temporarily own the displayed needle angle during power-on or power-off.

### Bar gauges

```yaml
realism:
  power:
    on:
      behaviour: fill_sweep
      sweep_to: max
      duration_ms: 800
      settle: live_value

    off:
      behaviour: drain
      target: min
      duration_ms: 1000
```

Bar power behaviour may temporarily own the displayed fill/reveal extent.

### Indicator gauges

```yaml
realism:
  power:
    on:
      behaviour: self_test
      duration_ms: 500

    off:
      behaviour: fade_out
      duration_ms: 300
```

Indicator power behaviour may briefly light the indicator for self-test before settling to the live state.

### Numeric gauges

```yaml
realism:
  power:
    on:
      behaviour: reveal
      duration_ms: 300

    off:
      behaviour: blank
```

Numeric power behaviour may blank, reveal, or hold values depending on the gauge.

### Odometer gauges

Odometer power behaviour should be conservative.

Possible behaviours:

```text
hold
blank
reveal
```

Do not animate odometer source values during power-on or power-off.

## Initial behaviour set

Start small:

```text
none
hold
blank
needle_sweep
slow_drop
fill_sweep
drain
self_test
fade_in
fade_out
reveal
```

Not every behaviour applies to every gauge family. Unsupported combinations should either be rejected during config validation or ignored with a clear warning, depending on existing config style.

## Composition rules

Power lifecycle is display-only.

It must not mutate:

```text
source values
logs
exports
configured ranges
stat marker samples
odometer source values
sensor state
```

During power-on or power-off, a gauge may temporarily own its displayed position or visibility. Once the lifecycle animation completes, normal gauge rendering and existing realism behaviour resume.

Startup sweeps, shutdown drops, indicator self-tests, fades, reveals, and blanks are presentation choreography. They are not sensor readings.

## Sampling rules

Other display features should not treat power lifecycle choreography as live signal data.

In particular:

- stat markers should ignore startup sweep samples;
- stat markers should ignore shutdown drop/drain samples;
- logs and exports should continue to reflect real source values, not lifecycle positions;
- replay should be able to simulate power events deterministically.

## Preview and test support

Add a mock/manual power signal provider so power-on and power-off behaviour can be previewed without vehicle hardware.

Useful preview controls:

```text
power on
power off
toggle power
replay scripted power events
```

Possible mock script shape:

```yaml
dashboard:
  power:
    input:
      type: mock
      script:
        - at_ms: 0
          state: off
        - at_ms: 1000
          state: on
        - at_ms: 8000
          state: off
```

## Do not

- Do not connect vehicle ACC wiring assumptions directly to gauge code.
- Do not require GPIO for desktop preview or tests.
- Do not feed lifecycle choreography into source values, logs, exports, stat markers, configured ranges, or odometer values.
- Do not make all gauges behave the same by dashboard fiat; each gauge should own its realistic power response.
- Do not add endless ambient animation. Power behaviour should be bounded and event-driven.

## Possible future slices

```text
power signal provider and runtime power events
gauge realism.power config model
radial gauge power sweep/drop
bar gauge power fill/drain
indicator gauge power self-test/fade
preview support for mock/manual power events
```
