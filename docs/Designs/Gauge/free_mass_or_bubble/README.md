---
gauge_group: free_mass_or_bubble
catalogue_version: "0.2"
primary_gauge_count: 2
supporting_quirk_count: 14
---

# Free mass, leaf, ball or bubble

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A freely deflecting mass, leaf, ball, pendulum or bubble is itself the indicator and responds directly to gravity, acceleration, charge or another applied influence.

**Catalogue definition:** A freely deflecting mass, leaf, ball, pendulum or bubble is itself the indicator.

## How the group encodes a value

Value or state is inferred from displacement, deflection, separation or motion of the free element.

## Classification boundary

Use this group when the visible body is not constrained to a conventional pointer linkage or buoyant level indication.

## Simulation baseline

Include inertia, oscillation, cross-axis response, contact with boundaries and orientation-dependent equilibrium.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 2 |
| Share of catalogue | 1.47% |
| Alternate members | 1 |
| Canonical quirks represented | 14 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `P32` | [Bubble level, ball inclinometer and pendulum clinometer](<gauges/p32_bubble_level_ball_inclinometer_and_pendulum_clinometer.md>) | Spirit level; aircraft slip-skid ball; marine clinometer; pendulum angle gauge | Tilt, bank, slope or lateral acceleration | 1600s-present |
| `X34` | [Gold-leaf or aluminium-leaf electroscope](<gauges/x34_gold_leaf_or_aluminium_leaf_electroscope.md>) | Bennet gold-leaf electroscope; Braun electroscope | Presence, sign or rough magnitude of electric charge | 1787-present; education and early physics |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `P32` | [Bubble level, ball inclinometer and pendulum clinometer](<gauges/p32_bubble_level_ball_inclinometer_and_pendulum_clinometer.md>) | Tilt, bank, slope or lateral acceleration | 1600s-present | 8 | 2 | 1 |
| `X34` | [Gold-leaf or aluminium-leaf electroscope](<gauges/x34_gold_leaf_or_aluminium_leaf_electroscope.md>) | Presence, sign or rough magnitude of electric charge | 1787-present; education and early physics | 7 | 1 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `X07` | [Turn-and-bank, turn coordinator and slip-skid ball](<../radial_pointer/gauges/x07_turn_and_bank_turn_coordinator_and_slip_skid_ball.md>) | [Radial pointer](<../radial_pointer/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Scale linearity and nonlinearity](<quirks/scale_linearity_and_nonlinearity.md>) | Variation in how equal input increments map to equal or unequal distances on the displayed scale. | 2 | 100.00% | 2 |
| [Bubble, void and separated-column behaviour](<quirks/bubble_void_and_separated_column_behaviour.md>) | Errors or discontinuities caused by trapped gas, empty spaces or broken fluid columns. | 1 | 50.00% | 2 |
| [Contamination, dirt and fouling](<quirks/contamination_dirt_and_fouling.md>) | Reading or appearance changes caused by deposits, dust, oxidation, residue or biological growth. | 1 | 50.00% | 1 |
| [Cross-axis and acceleration sensitivity](<quirks/cross_axis_and_acceleration_sensitivity.md>) | Response to acceleration or forces along axes other than the intended measurement axis. | 1 | 50.00% | 1 |
| [Damping](<quirks/damping.md>) | Deliberate or inherent suppression of rapid movement and oscillation in the indication. | 1 | 50.00% | 1 |
| [Density and viscosity dependence](<quirks/density_and_viscosity_dependence.md>) | Changes in indication caused by fluid density, viscosity or their variation with conditions. | 1 | 50.00% | 1 |
| [Electrostatic leakage and charge retention](<quirks/electrostatic_leakage_and_charge_retention.md>) | Loss or persistence of electric charge affecting electrostatic instruments and displays. | 1 | 50.00% | 3 |
| [Friction and drag](<quirks/friction_and_drag.md>) | Motion resistance that slows, biases or distorts the indication. | 1 | 50.00% | 1 |
| [Gravity-related behaviour](<quirks/gravity_related_behaviour.md>) | Dependence on local gravity magnitude or direction. | 1 | 50.00% | 1 |
| [Humidity and moisture sensitivity](<quirks/humidity_and_moisture_sensitivity.md>) | Changes caused by water vapour, condensation, absorption or damp contamination. | 1 | 50.00% | 1 |
| [Latching and state-retention behaviour](<quirks/latching_and_state_retention_behaviour.md>) | A displayed state that remains mechanically, magnetically, electrically or optically retained until reset or rewritten. | 1 | 50.00% | 1 |
| [Qualitative or non-precision indication](<quirks/qualitative_or_non_precision_indication.md>) | An indication intended for trend, state or rough comparison rather than accurate numeric measurement. | 1 | 50.00% | 1 |
| [Shock and vibration effects](<quirks/shock_and_vibration_effects.md>) | Temporary or permanent indication changes caused by mechanical shock or sustained vibration. | 1 | 50.00% | 1 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 1 | 50.00% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
