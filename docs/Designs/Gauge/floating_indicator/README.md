---
gauge_group: floating_indicator
catalogue_version: "0.2"
primary_gauge_count: 3
supporting_quirk_count: 22
---

# Floating indicator

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A buoyant float, bob, bulb or body moves within a fluid and its position directly represents the value.

**Catalogue definition:** A buoyant float, bob or bulb directly indicates the value by its position.

## How the group encodes a value

Value is encoded by the equilibrium position of the floating element, sometimes against a fixed scale and sometimes by identifying which calibrated float settles where.

## Classification boundary

Use this group when buoyancy positions the visible indicator. A liquid surface without a separate float belongs under liquid_column.

## Simulation baseline

Model buoyancy, density, drag, wall contact, bobbing and settling rather than treating the indicator as a massless cursor.

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
| `P20` | [Variable-area flowmeter / rotameter](<gauges/p20_variable_area_flowmeter_rotameter.md>) | KROHNE Rotameter; Brooks Sho-Rate; glass or armoured metal tubes | Volumetric or mass flow after calibration | 1908-present |
| `P28` | [Hydrometer and floating density gauges](<gauges/p28_hydrometer_and_floating_density_gauges.md>) | Lactometer, alcoholmeter, API hydrometer, battery hydrometer | Liquid density or concentration | antiquity-present; precision glass forms from 1700s |
| `P29` | [Galileo thermometer](<gauges/p29_galileo_thermometer.md>) | Sealed column with weighted glass bulbs and temperature tags | Approximate ambient temperature | 1600s concept; decorative use from 1900s-present |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `P20` | [Variable-area flowmeter / rotameter](<gauges/p20_variable_area_flowmeter_rotameter.md>) | Volumetric or mass flow after calibration | 1908-present | 13 | 2 | 1 |
| `P28` | [Hydrometer and floating density gauges](<gauges/p28_hydrometer_and_floating_density_gauges.md>) | Liquid density or concentration | antiquity-present; precision glass forms from 1700s | 10 | 2 | 1 |
| `P29` | [Galileo thermometer](<gauges/p29_galileo_thermometer.md>) | Approximate ambient temperature | 1600s concept; decorative use from 1900s-present | 9 | 1 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `X30` | [Mechanical float tide gauge and stilling-well recorder](<../chart_or_trace_recorder/gauges/x30_mechanical_float_tide_gauge_and_stilling_well_recorder.md>) | [Chart or trace recorder](<../chart_or_trace_recorder/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Float and buoyancy behaviour](<quirks/float_and_buoyancy_behaviour.md>) | Motion and equilibrium effects arising from buoyant indicators in fluids. | 3 | 100.00% | 5 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 3 | 100.00% | 3 |
| [Bounce](<quirks/bounce.md>) | Repeated rebounds or reversals after a mechanical impact, contact change or rapid movement. | 2 | 66.67% | 2 |
| [Bubble, void and separated-column behaviour](<quirks/bubble_void_and_separated_column_behaviour.md>) | Errors or discontinuities caused by trapped gas, empty spaces or broken fluid columns. | 2 | 66.67% | 2 |
| [Calibration, correction and compensation](<quirks/calibration_correction_and_compensation.md>) | Adjustments or correction factors required to relate the raw indication to the intended quantity. | 2 | 66.67% | 2 |
| [Density and viscosity dependence](<quirks/density_and_viscosity_dependence.md>) | Changes in indication caused by fluid density, viscosity or their variation with conditions. | 2 | 66.67% | 2 |
| [Scale direction and interpretation](<quirks/scale_direction_and_interpretation.md>) | Whether increasing values move clockwise, anticlockwise, upward or otherwise, and how that direction is understood. | 2 | 66.67% | 3 |
| [Stiction and sticking](<quirks/stiction_and_sticking.md>) | Static friction or adhesion that prevents motion until enough force accumulates, often followed by a jump. | 2 | 66.67% | 2 |
| [Ambient-pressure and altitude sensitivity](<quirks/ambient_pressure_and_altitude_sensitivity.md>) | Changes in indication caused by surrounding atmospheric pressure or elevation. | 1 | 33.33% | 1 |
| [Chatter](<quirks/chatter.md>) | Rapid repeated switching or movement near a threshold or unstable equilibrium. | 1 | 33.33% | 1 |
| [Construction tolerances and unit variation](<quirks/construction_tolerances_and_unit_variation.md>) | Differences between nominally identical instruments caused by manufacturing and assembly tolerances. | 1 | 33.33% | 1 |
| [Contamination, dirt and fouling](<quirks/contamination_dirt_and_fouling.md>) | Reading or appearance changes caused by deposits, dust, oxidation, residue or biological growth. | 1 | 33.33% | 1 |
| [Damping](<quirks/damping.md>) | Deliberate or inherent suppression of rapid movement and oscillation in the indication. | 1 | 33.33% | 1 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 1 | 33.33% | 2 |
| [Human-factor ambiguity and misreading](<quirks/human_factor_ambiguity_and_misreading.md>) | Display features that make an otherwise functioning instrument easy to interpret incorrectly. | 1 | 33.33% | 1 |
| [Meniscus behaviour](<quirks/meniscus_behaviour.md>) | The curved liquid surface and its reading conventions, wetting shape and movement. | 1 | 33.33% | 1 |
| [Orientation sensitivity](<quirks/orientation_sensitivity.md>) | Dependence on instrument mounting angle or orientation relative to gravity or fields. | 1 | 33.33% | 1 |
| [Pressure pulsation and pneumatic-line dynamics](<quirks/pressure_pulsation_and_pneumatic_line_dynamics.md>) | Indication effects caused by pressure waves, tubing volume, restrictions and compressible-fluid behaviour. | 1 | 33.33% | 1 |
| [Qualitative or non-precision indication](<quirks/qualitative_or_non_precision_indication.md>) | An indication intended for trend, state or rough comparison rather than accurate numeric measurement. | 1 | 33.33% | 1 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 1 | 33.33% | 1 |
| [Scale linearity and nonlinearity](<quirks/scale_linearity_and_nonlinearity.md>) | Variation in how equal input increments map to equal or unequal distances on the displayed scale. | 1 | 33.33% | 1 |
| [Settling and return behaviour](<quirks/settling_and_return_behaviour.md>) | How the indication approaches a stable value or returns after a transient, disturbance or release. | 1 | 33.33% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
