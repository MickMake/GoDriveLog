---
gauge_group: rotating_scan_or_strobe
catalogue_version: "0.2"
primary_gauge_count: 4
supporting_quirk_count: 19
---

# Rotating scan or strobe

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A rotating scan, glow point, wheel or timed strobe converts pulse timing, frequency or phase into an apparent angular or spatial position.

**Catalogue definition:** Timing is converted to apparent position by a rotating scan, flash, glow point or strobe.

## How the group encodes a value

The observed value depends on synchronism between a rotating or sequential process and the input signal.

## Classification boundary

Use this group when scanning or stroboscopic timing creates the indication, rather than merely refreshing a conventional display.

## Simulation baseline

Represent phase, scan cadence, persistence, aliasing, false locks, rotation direction, missed pulses and startup synchronisation.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 4 |
| Share of catalogue | 2.94% |
| Alternate members | 0 |
| Canonical quirks represented | 19 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `D09` | [Dekatron / glow-transfer counting tube](<gauges/d09_dekatron_glow_transfer_counting_tube.md>) | Ericsson and Mullard GC10-series; ten-position and divide-by-ten tubes | Counting state, divider position or rotating indicator | 1949-c.1970s common; collector demonstrations today |
| `P43` | [Stroboscopic tachometer](<gauges/p43_stroboscopic_tachometer.md>) | Mechanical vibrating-reed strobe; electronic xenon/LED stroboscope | Rotational speed or repeated-motion frequency | c.1910s-present |
| `X14` | [Rotating-neon fathometer](<gauges/x14_rotating_neon_fathometer.md>) | Dorsey Fathometer; Submarine Signal Company Model 808 | Water depth from acoustic echo time | 1920s-c.1960s common |
| `X15` | [Mechanical-wheel flasher sonar](<gauges/x15_mechanical_wheel_flasher_sonar.md>) | Vexilar FL/FLX series and earlier ice-fishing flashers | Current sonar return versus depth | 1960s-present |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `D09` | [Dekatron / glow-transfer counting tube](<gauges/d09_dekatron_glow_transfer_counting_tube.md>) | Counting state, divider position or rotating indicator | 1949-c.1970s common; collector demonstrations today | 5 | 2 | 1 |
| `P43` | [Stroboscopic tachometer](<gauges/p43_stroboscopic_tachometer.md>) | Rotational speed or repeated-motion frequency | c.1910s-present | 6 | 1 | 1 |
| `X14` | [Rotating-neon fathometer](<gauges/x14_rotating_neon_fathometer.md>) | Water depth from acoustic echo time | 1920s-c.1960s common | 6 | 2 | 1 |
| `X15` | [Mechanical-wheel flasher sonar](<gauges/x15_mechanical_wheel_flasher_sonar.md>) | Current sonar return versus depth | 1960s-present | 9 | 1 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Persistence and afterglow](<quirks/persistence_and_afterglow.md>) | Continued visibility after excitation is reduced or removed. | 4 | 100.00% | 4 |
| [Latching and state-retention behaviour](<quirks/latching_and_state_retention_behaviour.md>) | A displayed state that remains mechanically, magnetically, electrically or optically retained until reset or rewritten. | 2 | 50.00% | 2 |
| [Mechanical noise and cadence](<quirks/mechanical_noise_and_cadence.md>) | Audible clicks, hums, impacts or rhythms produced by the display mechanism. | 2 | 50.00% | 2 |
| [Multiple echoes, interference and remapped scale](<quirks/multiple_echoes_interference_and_remapped_scale.md>) | Ambiguous or transformed indications caused by multiple returns, interference or nonlinear mapping of time to scale. | 2 | 50.00% | 4 |
| [Periodic marker and rotating-scan behaviour](<quirks/periodic_marker_and_rotating_scan_behaviour.md>) | Repeated markers or scan positions tied to rotation, timing or a cyclic reference. | 2 | 50.00% | 2 |
| [Aliasing, harmonics and false solutions](<quirks/aliasing_harmonics_and_false_solutions.md>) | Incorrect apparent readings caused by sampling, strobing, harmonics or multiple possible synchronisation points. | 1 | 25.00% | 2 |
| [Ambient-light readability](<quirks/ambient_light_readability.md>) | Dependence of legibility on sunlight, darkness, glare or surrounding illumination. | 1 | 25.00% | 1 |
| [Bloom, halo and penumbra](<quirks/bloom_halo_and_penumbra.md>) | Spreading or soft-edged light around a spot, trace, segment or projected image. | 1 | 25.00% | 1 |
| [Calibration, correction and compensation](<quirks/calibration_correction_and_compensation.md>) | Adjustments or correction factors required to relate the raw indication to the intended quantity. | 1 | 25.00% | 1 |
| [Colour variation and colour shift](<quirks/colour_variation_and_colour_shift.md>) | Changes or inconsistencies in displayed colour across level, age, temperature, angle or individual units. | 1 | 25.00% | 1 |
| [Dead time, missed events and duplicate counts](<quirks/dead_time_missed_events_and_duplicate_counts.md>) | Intervals or mechanisms that cause events to be ignored, delayed or counted more than once. | 1 | 25.00% | 1 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 1 | 25.00% | 2 |
| [Drift and long-term stability](<quirks/drift_and_long_term_stability.md>) | Slow change in indicated value or behaviour over time despite an unchanged input. | 1 | 25.00% | 1 |
| [Flicker, scan and PWM artefacts](<quirks/flicker_scan_and_pwm_artefacts.md>) | Visible modulation caused by multiplexing, scanning, pulse-width control or interaction with cameras and eye motion. | 1 | 25.00% | 1 |
| [Gas-discharge instability and dropout](<quirks/gas_discharge_instability_and_dropout.md>) | Flicker, extinction or irregular conduction after a gas discharge has started. | 1 | 25.00% | 2 |
| [Human-factor ambiguity and misreading](<quirks/human_factor_ambiguity_and_misreading.md>) | Display features that make an otherwise functioning instrument easy to interpret incorrectly. | 1 | 25.00% | 1 |
| [Phase, synchronism and rotation direction](<quirks/phase_synchronism_and_rotation_direction.md>) | Dependence on relative phase, locked timing or the direction of a rotating field or mechanism. | 1 | 25.00% | 1 |
| [Scale direction and interpretation](<quirks/scale_direction_and_interpretation.md>) | Whether increasing values move clockwise, anticlockwise, upward or otherwise, and how that direction is understood. | 1 | 25.00% | 1 |
| [Thresholds, deadband and switching points](<quirks/thresholds_deadband_and_switching_points.md>) | Regions or levels where no change occurs, or where a discrete state changes with defined or variable thresholds. | 1 | 25.00% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
