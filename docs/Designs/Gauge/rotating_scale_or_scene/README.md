---
gauge_group: rotating_scale_or_scene
catalogue_version: "0.2"
primary_gauge_count: 3
supporting_quirk_count: 22
---

# Rotating scale or scene

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A compass card, scale, drum, horizon or scene rotates behind one or more fixed references or overlaid pointers.

**Catalogue definition:** A compass card, scale, drum, horizon or scene rotates behind a fixed reference.

## How the group encodes a value

Value is encoded by the angular relationship between the moving background and a fixed index, aircraft symbol or pointer system.

## Classification boundary

Use this group when the scale or scene moves. If the scale is fixed and the hand moves, use radial_pointer.

## Simulation baseline

Separate scene rotation, pointer overlays, card inertia, gimbal limits, precession and visual masking.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 3 |
| Share of catalogue | 2.21% |
| Alternate members | 1 |
| Canonical quirks represented | 22 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `P33` | [Liquid magnetic compass card](<gauges/p33_liquid_magnetic_compass_card.md>) | Marine bowl compass; aircraft whiskey compass; floating card with lubber line | Magnetic heading | 1200s-present |
| `X06` | [Mechanical attitude indicator / artificial horizon](<gauges/x06_mechanical_attitude_indicator_artificial_horizon.md>) | Vacuum- or electrically driven gyro horizon; caged aerobatic variants | Aircraft pitch and bank attitude | 1920s-present |
| `X09` | [RMI and HSI multi-pointer navigation indicator](<gauges/x09_rmi_and_hsi_multi_pointer_navigation_indicator.md>) | Bendix radio magnetic indicator; horizontal situation indicator | Magnetic heading, selected course, bearing and deviation | 1950s-present |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `P33` | [Liquid magnetic compass card](<gauges/p33_liquid_magnetic_compass_card.md>) | Magnetic heading | 1200s-present | 10 | 2 | 1 |
| `X06` | [Mechanical attitude indicator / artificial horizon](<gauges/x06_mechanical_attitude_indicator_artificial_horizon.md>) | Aircraft pitch and bank attitude | 1920s-present | 8 | 2 | 1 |
| `X09` | [RMI and HSI multi-pointer navigation indicator](<gauges/x09_rmi_and_hsi_multi_pointer_navigation_indicator.md>) | Magnetic heading, selected course, bearing and deviation | 1950s-present | 7 | 2 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `E11` | [Ferraris induction meter](<../rolling_drum_or_counter/gauges/e11_ferraris_induction_meter.md>) | [Rolling drum or counter](<../rolling_drum_or_counter/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Invalid, out-of-range and warning flags](<quirks/invalid_out_of_range_and_warning_flags.md>) | Explicit indications that a reading is unavailable, unreliable, unsafe or beyond the valid range. | 2 | 66.67% | 2 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 2 | 66.67% | 2 |
| [Wraparound and multi-turn indication](<quirks/wraparound_and_multi_turn_indication.md>) | Behaviour when a value crosses a cyclic boundary or requires multiple revolutions to represent its full range. | 2 | 66.67% | 2 |
| [Bubble, void and separated-column behaviour](<quirks/bubble_void_and_separated_column_behaviour.md>) | Errors or discontinuities caused by trapped gas, empty spaces or broken fluid columns. | 1 | 33.33% | 1 |
| [Cross-axis and acceleration sensitivity](<quirks/cross_axis_and_acceleration_sensitivity.md>) | Response to acceleration or forces along axes other than the intended measurement axis. | 1 | 33.33% | 1 |
| [Damping](<quirks/damping.md>) | Deliberate or inherent suppression of rapid movement and oscillation in the indication. | 1 | 33.33% | 1 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 1 | 33.33% | 1 |
| [Drift and long-term stability](<quirks/drift_and_long_term_stability.md>) | Slow change in indicated value or behaviour over time despite an unchanged input. | 1 | 33.33% | 1 |
| [Gyro erection, precession and tumble](<quirks/gyro_erection_precession_and_tumble.md>) | Gyroscopic recovery, drift and loss-of-attitude behaviours associated with spinning-mass instruments. | 1 | 33.33% | 3 |
| [Hunting and following error](<quirks/hunting_and_following_error.md>) | Repeated corrective motion or persistent difference between commanded and indicated position in a following system. | 1 | 33.33% | 1 |
| [Magnetic-field, deviation and remanence effects](<quirks/magnetic_field_deviation_and_remanence_effects.md>) | Influence of external or retained magnetism on indication, zero and calibration. | 1 | 33.33% | 2 |
| [Manual reset, tare and caging](<quirks/manual_reset_tare_and_caging.md>) | User-operated mechanisms that clear, zero, restrain or protect the indication. | 1 | 33.33% | 2 |
| [Multi-pointer and multi-channel interaction](<quirks/multi_pointer_and_multi_channel_interaction.md>) | Visual or mechanical interaction among multiple pointers, scales or measurement channels. | 1 | 33.33% | 4 |
| [Optical distortion and refraction](<quirks/optical_distortion_and_refraction.md>) | Apparent displacement or shape changes caused by lenses, glass, liquid, curved windows or refractive interfaces. | 1 | 33.33% | 1 |
| [Orientation sensitivity](<quirks/orientation_sensitivity.md>) | Dependence on instrument mounting angle or orientation relative to gravity or fields. | 1 | 33.33% | 1 |
| [Overload, saturation and damage](<quirks/overload_saturation_and_damage.md>) | Behaviour when the input exceeds the useful range, including pegging, clipping, recovery changes or permanent harm. | 1 | 33.33% | 1 |
| [Overshoot](<quirks/overshoot.md>) | Temporary travel beyond the final steady indication after an input change. | 1 | 33.33% | 1 |
| [Parallax](<quirks/parallax.md>) | Apparent reading error caused by viewing the indicator and scale from the wrong angle or depth relationship. | 1 | 33.33% | 1 |
| [Power-off and power-loss behaviour](<quirks/power_off_and_power_loss_behaviour.md>) | What the indication does when drive power is removed, including retained, blank, parked or misleading states. | 1 | 33.33% | 1 |
| [Settling and return behaviour](<quirks/settling_and_return_behaviour.md>) | How the indication approaches a stable value or returns after a transient, disturbance or release. | 1 | 33.33% | 1 |
| [Shadows, depth and occlusion](<quirks/shadows_depth_and_occlusion.md>) | Visual effects caused by layered parts blocking, shading or appearing at different depths. | 1 | 33.33% | 1 |
| [Warm-up behaviour](<quirks/warm_up_behaviour.md>) | Transient changes after startup while temperature, discharge, illumination or mechanics stabilise. | 1 | 33.33% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
