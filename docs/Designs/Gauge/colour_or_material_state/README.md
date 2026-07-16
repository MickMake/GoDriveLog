---
gauge_group: colour_or_material_state
catalogue_version: "0.2"
primary_gauge_count: 3
supporting_quirk_count: 20
---

# Colour or material-state indicator

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

Colour, opacity, scattering, phase, chemical state or another material property changes to indicate value, exposure or history.

**Catalogue definition:** Colour, opacity, scattering, phase or chemical state changes to indicate value or history.

## How the group encodes a value

The material itself becomes the display, sometimes reversibly and sometimes irreversibly.

## Classification boundary

Use this group when material-state change is the primary indication rather than merely the light source behind conventional segments.

## Simulation baseline

Model transition thresholds, diffusion, hysteresis, temperature dependence, ageing, nonuniformity and irreversible response where applicable.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 3 |
| Share of catalogue | 2.21% |
| Alternate members | 0 |
| Canonical quirks represented | 20 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `D18` | [Powder and thin-film electroluminescent display](<gauges/d18_powder_and_thin_film_electroluminescent_display.md>) | EL panels; Lumineq/Beneq rugged displays; automotive backlights | Fixed legends, segmented readouts or monochrome matrix graphics | 1950s-present; matrix peak c.1980s-2000s |
| `D30` | [Irreversible chemical and time-temperature indicator](<gauges/d30_irreversible_chemical_and_time_temperature_indicator.md>) | Time-temperature integrators; freeze indicators; sterilisation tape; gas exposure badges | Cumulative thermal, chemical or environmental exposure | mid-1900s-present |
| `P25` | [Thermochromic liquid-crystal strip and leuco-dye indicator](<gauges/p25_thermochromic_liquid_crystal_strip_and_leuco_dye_indicator.md>) | Forehead strips; aquarium strips; battery/test labels; reversible and irreversible temperature labels | Surface temperature or threshold exposure | 1960s-present |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `D18` | [Powder and thin-film electroluminescent display](<gauges/d18_powder_and_thin_film_electroluminescent_display.md>) | Fixed legends, segmented readouts or monochrome matrix graphics | 1950s-present; matrix peak c.1980s-2000s | 7 | 2 | 1 |
| `D30` | [Irreversible chemical and time-temperature indicator](<gauges/d30_irreversible_chemical_and_time_temperature_indicator.md>) | Cumulative thermal, chemical or environmental exposure | mid-1900s-present | 9 | 1 | 1 |
| `P25` | [Thermochromic liquid-crystal strip and leuco-dye indicator](<gauges/p25_thermochromic_liquid_crystal_strip_and_leuco_dye_indicator.md>) | Surface temperature or threshold exposure | 1960s-present | 12 | 2 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Ambient-light readability](<quirks/ambient_light_readability.md>) | Dependence of legibility on sunlight, darkness, glare or surrounding illumination. | 3 | 100.00% | 3 |
| [Colour variation and colour shift](<quirks/colour_variation_and_colour_shift.md>) | Changes or inconsistencies in displayed colour across level, age, temperature, angle or individual units. | 3 | 100.00% | 4 |
| [Irreversible or one-way response](<quirks/irreversible_or_one_way_response.md>) | An indication that cannot return to its prior state without replacement or a separate reset process. | 2 | 66.67% | 4 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 2 | 66.67% | 2 |
| [Viewing-angle dependence](<quirks/viewing_angle_dependence.md>) | Changes in readability, colour, contrast or apparent value with observer angle. | 2 | 66.67% | 2 |
| [Witness, peak-hold and retained extrema](<quirks/witness_peak_hold_and_retained_extrema.md>) | Markers or memory mechanisms that preserve minimum, maximum or peak values after the live indication moves away. | 2 | 66.67% | 2 |
| [Ageing and material degradation](<quirks/ageing_and_material_degradation.md>) | Progressive change due to wear, fatigue, chemical decay, phosphor loss, embrittlement or similar ageing processes. | 1 | 33.33% | 1 |
| [Brightness and contrast variation](<quirks/brightness_and_contrast_variation.md>) | Changes in luminance or visual separation between active and inactive parts of the display. | 1 | 33.33% | 3 |
| [Chemical diffusion and reaction front](<quirks/chemical_diffusion_and_reaction_front.md>) | Movement or spread of a chemical change through a material used as the indication. | 1 | 33.33% | 2 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 1 | 33.33% | 1 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 1 | 33.33% | 1 |
| [Display overlap and adjacent indication](<quirks/display_overlap_and_adjacent_indication.md>) | Visual interference where neighbouring pointers, digits, traces or display regions overlap. | 1 | 33.33% | 1 |
| [Electrical hum, buzz and whine](<quirks/electrical_hum_buzz_and_whine.md>) | Audible vibration or tone produced by electrical drive, magnetic parts or switching. | 1 | 33.33% | 1 |
| [Humidity and moisture sensitivity](<quirks/humidity_and_moisture_sensitivity.md>) | Changes caused by water vapour, condensation, absorption or damp contamination. | 1 | 33.33% | 1 |
| [Hysteresis](<quirks/hysteresis.md>) | Different indicated values for the same input depending on whether the input approached from above or below. | 1 | 33.33% | 1 |
| [Immersion, stem-conduction and contact error](<quirks/immersion_stem_conduction_and_contact_error.md>) | Temperature-reading errors caused by installation depth, heat flow along a probe or imperfect thermal contact. | 1 | 33.33% | 1 |
| [Manual reset, tare and caging](<quirks/manual_reset_tare_and_caging.md>) | User-operated mechanisms that clear, zero, restrain or protect the indication. | 1 | 33.33% | 1 |
| [Random, statistical and batch variation](<quirks/random_statistical_and_batch_variation.md>) | Unpredictable events or unit-to-unit differences that are part of the observed behaviour. | 1 | 33.33% | 1 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 1 | 33.33% | 1 |
| [Thresholds, deadband and switching points](<quirks/thresholds_deadband_and_switching_points.md>) | Regions or levels where no change occurs, or where a discrete state changes with defined or variable thresholds. | 1 | 33.33% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
