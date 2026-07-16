---
gauge_group: linear_pointer
catalogue_version: "0.2"
primary_gauge_count: 2
supporting_quirk_count: 15
---

# Linear pointer

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A pointer, blade, hairline or marker moves along a straight scale rather than around a pivoted dial.

**Catalogue definition:** A pointer, blade, hairline or marker moves along a straight or edgewise scale.

## How the group encodes a value

Value is encoded as linear displacement along one axis, often viewed edgewise through a narrow window.

## Classification boundary

Use this group for a discrete moving pointer. A continuously travelling tape or ribbon belongs under moving_tape_ribbon_or_map.

## Simulation baseline

Model travel limits, scale projection, parallax, friction and any linkage nonlinearity independently from the source value.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 2 |
| Share of catalogue | 1.47% |
| Alternate members | 1 |
| Canonical quirks represented | 15 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `E22` | [Edgewise meter](<gauges/e22_edgewise_meter.md>) | Weston and GE edgewise switchboard meters; rack-mounted process indicators | Electrical or process variables in dense panels | 1930s-present; peak c.1940s-1980s |
| `P36` | [Quartz-fibre self-reading pocket dosimeter](<gauges/p36_quartz_fibre_self_reading_pocket_dosimeter.md>) | Lauritsen electroscope; Victoreen Model 541; Landsverk pocket dosimeters | Accumulated ionising-radiation dose | 1937-c.1990s common; collectors and limited specialist use today |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `E22` | [Edgewise meter](<gauges/e22_edgewise_meter.md>) | Electrical or process variables in dense panels | 1930s-present; peak c.1940s-1980s | 4 | 1 | 1 |
| `P36` | [Quartz-fibre self-reading pocket dosimeter](<gauges/p36_quartz_fibre_self_reading_pocket_dosimeter.md>) | Accumulated ionising-radiation dose | 1937-c.1990s common; collectors and limited specialist use today | 11 | 2 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `P37` | [Spring balance, dial force gauge and dynamometer](<../radial_pointer/gauges/p37_spring_balance_dial_force_gauge_and_dynamometer.md>) | [Radial pointer](<../radial_pointer/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Damping](<quirks/damping.md>) | Deliberate or inherent suppression of rapid movement and oscillation in the indication. | 1 | 50.00% | 1 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 1 | 50.00% | 3 |
| [Drift and long-term stability](<quirks/drift_and_long_term_stability.md>) | Slow change in indicated value or behaviour over time despite an unchanged input. | 1 | 50.00% | 1 |
| [Electrostatic leakage and charge retention](<quirks/electrostatic_leakage_and_charge_retention.md>) | Loss or persistence of electric charge affecting electrostatic instruments and displays. | 1 | 50.00% | 2 |
| [Focus and optical alignment](<quirks/focus_and_optical_alignment.md>) | Sharpness and registration effects caused by optical focus and component alignment. | 1 | 50.00% | 1 |
| [Humidity and moisture sensitivity](<quirks/humidity_and_moisture_sensitivity.md>) | Changes caused by water vapour, condensation, absorption or damp contamination. | 1 | 50.00% | 1 |
| [Manual reset, tare and caging](<quirks/manual_reset_tare_and_caging.md>) | User-operated mechanisms that clear, zero, restrain or protect the indication. | 1 | 50.00% | 1 |
| [Operator procedure and ritual](<quirks/operator_procedure_and_ritual.md>) | Required handling or reading practices that materially affect the result. | 1 | 50.00% | 1 |
| [Overload, saturation and damage](<quirks/overload_saturation_and_damage.md>) | Behaviour when the input exceeds the useful range, including pegging, clipping, recovery changes or permanent harm. | 1 | 50.00% | 1 |
| [Parallax](<quirks/parallax.md>) | Apparent reading error caused by viewing the indicator and scale from the wrong angle or depth relationship. | 1 | 50.00% | 1 |
| [Readout masking and narrow-window transitions](<quirks/readout_masking_and_narrow_window_transitions.md>) | Partial or ambiguous readings caused by apertures that reveal only a small portion of a moving scale or drum. | 1 | 50.00% | 1 |
| [Shock and vibration effects](<quirks/shock_and_vibration_effects.md>) | Temporary or permanent indication changes caused by mechanical shock or sustained vibration. | 1 | 50.00% | 1 |
| [Valid operating region](<quirks/valid_operating_region.md>) | The part of the scale or operating envelope in which the indication is considered reliable. | 1 | 50.00% | 1 |
| [Viewing-angle dependence](<quirks/viewing_angle_dependence.md>) | Changes in readability, colour, contrast or apparent value with observer angle. | 1 | 50.00% | 1 |
| [Zero drift and offset](<quirks/zero_drift_and_offset.md>) | A non-zero indication at the true zero point, including offsets that change with time or conditions. | 1 | 50.00% | 2 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
