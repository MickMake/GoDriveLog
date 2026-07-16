---
gauge_group: rolling_drum_or_counter
catalogue_version: "0.2"
primary_gauge_count: 3
supporting_quirk_count: 20
---

# Rolling drum or counter

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

One or more numeral drums, cyclometer wheels or register wheels rotate to form a cumulative or direct numeric reading.

**Catalogue definition:** Numeral drums, cyclometer wheels or register wheels rotate to form a count or value.

## How the group encodes a value

Value is encoded by aligned characters on rotating wheels, usually with mechanical carry between adjacent decades or units.

## Classification boundary

Use this group for rotating numeral registers. A full moving tape belongs under moving_tape_ribbon_or_map; a luminous segmented numeral belongs under segmented_display.

## Simulation baseline

Carry timing, partial digits, wheel alignment, backlash and rollover are essential. Perfectly simultaneous digit changes look suspiciously modern.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 3 |
| Share of catalogue | 2.21% |
| Alternate members | 4 |
| Canonical quirks represented | 20 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `E11` | [Ferraris induction meter](<gauges/e11_ferraris_induction_meter.md>) | Rotating-disc electricity meter; Sangamo and Landis & Gyr watt-hour meters | Integrated AC electrical energy | 1880s-present; electronic replacement accelerating since 1990s |
| `E24` | [Roller, drum and cyclometer indicator](<gauges/e24_roller_drum_and_cyclometer_indicator.md>) | Rolling-number frequency meters; drum speedometers; counter drums paired with a pointer | Discrete or quasi-continuous numeric value | 1800s-present |
| `P41` | [Mechanical odometer, cyclometer and utility-meter register](<gauges/p41_mechanical_odometer_cyclometer_and_utility_meter_register.md>) | Number drums, Geneva carries, pointer sub-dials; gas, water and electricity registers | Distance, volume, energy, cycles or accumulated count | 1600s-present |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `E11` | [Ferraris induction meter](<gauges/e11_ferraris_induction_meter.md>) | Integrated AC electrical energy | 1880s-present; electronic replacement accelerating since 1990s | 10 | 2 | 1 |
| `E24` | [Roller, drum and cyclometer indicator](<gauges/e24_roller_drum_and_cyclometer_indicator.md>) | Discrete or quasi-continuous numeric value | 1800s-present | 5 | 1 | 1 |
| `P41` | [Mechanical odometer, cyclometer and utility-meter register](<gauges/p41_mechanical_odometer_cyclometer_and_utility_meter_register.md>) | Distance, volume, energy, cycles or accumulated count | 1600s-present | 9 | 2 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `P39` | [Mechanical tachometer and eddy-current speedometer](<../radial_pointer/gauges/p39_mechanical_tachometer_and_eddy_current_speedometer.md>) | [Radial pointer](<../radial_pointer/README.md>) |
| `X02` | [Counter-pointer altimeter](<../radial_pointer/gauges/x02_counter_pointer_altimeter.md>) | [Radial pointer](<../radial_pointer/README.md>) |
| `X03` | [Drum-pointer and counter-drum altimeter](<../radial_pointer/gauges/x03_drum_pointer_and_counter_drum_altimeter.md>) | [Radial pointer](<../radial_pointer/README.md>) |
| `X20` | [Ribbon, rolling-drum and “thermometer” automotive speedometer](<../moving_tape_ribbon_or_map/gauges/x20_ribbon_rolling_drum_and_thermometer_automotive_speedometer.md>) | [Moving tape, ribbon or map](<../moving_tape_ribbon_or_map/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Backlash and lash](<quirks/backlash_and_lash.md>) | Lost motion when direction reverses because of clearance in gears, linkages or drive parts. | 3 | 100.00% | 3 |
| [Carry, rollover and digit transition](<quirks/carry_rollover_and_digit_transition.md>) | The mechanical or visual sequence as counters advance between digits and propagate carries. | 2 | 66.67% | 6 |
| [Inertia and coasting](<quirks/inertia_and_coasting.md>) | Continued or delayed motion due to mass and stored kinetic energy. | 2 | 66.67% | 3 |
| [Creep](<quirks/creep.md>) | Very slow movement under a constant input, load or retained state. | 1 | 33.33% | 1 |
| [Damping](<quirks/damping.md>) | Deliberate or inherent suppression of rapid movement and oscillation in the indication. | 1 | 33.33% | 1 |
| [Electrical hum, buzz and whine](<quirks/electrical_hum_buzz_and_whine.md>) | Audible vibration or tone produced by electrical drive, magnetic parts or switching. | 1 | 33.33% | 1 |
| [Flutter, jitter, tremor and quiver](<quirks/flutter_jitter_tremor_and_quiver.md>) | Small rapid random or periodic movements around the nominal indication. | 1 | 33.33% | 1 |
| [Manual reset, tare and caging](<quirks/manual_reset_tare_and_caging.md>) | User-operated mechanisms that clear, zero, restrain or protect the indication. | 1 | 33.33% | 1 |
| [Mechanical wear, stretch and permanent set](<quirks/mechanical_wear_stretch_and_permanent_set.md>) | Long-term dimensional or elastic change in moving parts, springs, fibres or linkages. | 1 | 33.33% | 1 |
| [Movement torque and sensitivity](<quirks/movement_torque_and_sensitivity.md>) | The relationship between applied drive, restoring force and resulting visible movement. | 1 | 33.33% | 1 |
| [Periodic marker and rotating-scan behaviour](<quirks/periodic_marker_and_rotating_scan_behaviour.md>) | Repeated markers or scan positions tied to rotation, timing or a cyclic reference. | 1 | 33.33% | 1 |
| [Power-off and power-loss behaviour](<quirks/power_off_and_power_loss_behaviour.md>) | What the indication does when drive power is removed, including retained, blank, parked or misleading states. | 1 | 33.33% | 1 |
| [Readout masking and narrow-window transitions](<quirks/readout_masking_and_narrow_window_transitions.md>) | Partial or ambiguous readings caused by apertures that reveal only a small portion of a moving scale or drum. | 1 | 33.33% | 1 |
| [Reversal and direction-change error](<quirks/reversal_and_direction_change_error.md>) | Error or transient behaviour specifically introduced when motion or input direction reverses. | 1 | 33.33% | 1 |
| [Revolution counting and coarse/fine readout](<quirks/revolution_counting_and_coarse_fine_readout.md>) | Methods that combine turns, counters or multiple indicators to represent both large range and fine resolution. | 1 | 33.33% | 1 |
| [Scale direction and interpretation](<quirks/scale_direction_and_interpretation.md>) | Whether increasing values move clockwise, anticlockwise, upward or otherwise, and how that direction is understood. | 1 | 33.33% | 1 |
| [Settling and return behaviour](<quirks/settling_and_return_behaviour.md>) | How the indication approaches a stable value or returns after a transient, disturbance or release. | 1 | 33.33% | 1 |
| [Snap action](<quirks/snap_action.md>) | A rapid transition between stable positions once a threshold is crossed. | 1 | 33.33% | 1 |
| [Stiction and sticking](<quirks/stiction_and_sticking.md>) | Static friction or adhesion that prevents motion until enough force accumulates, often followed by a jump. | 1 | 33.33% | 1 |
| [Tamper evidence and official seals](<quirks/tamper_evidence_and_official_seals.md>) | Marks, seals or mechanisms that reveal unauthorised adjustment or preserve legal metrology status. | 1 | 33.33% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
