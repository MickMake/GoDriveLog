---
gauge_group: liquid_column
catalogue_version: "0.2"
primary_gauge_count: 9
supporting_quirk_count: 42
---

# Liquid column

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A liquid height, meniscus, interface or differential column directly indicates the measured quantity in a tube, well, bulb or gauge glass.

**Catalogue definition:** A liquid height, meniscus or liquid/gas boundary moves in a tube, well or gauge glass.

## How the group encodes a value

Value is encoded by one or more fluid boundaries. Tube geometry, liquid properties, gravity and viewing position can all affect the apparent reading.

## Classification boundary

Use this group when the observed liquid column itself is the readout. A free floating object within a liquid belongs under floating_indicator.

## Simulation baseline

The meniscus and fluid dynamics matter as much as the nominal level: include wetting, slosh, bubbles, settling and temperature-dependent properties where relevant.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 9 |
| Share of catalogue | 6.62% |
| Alternate members | 2 |
| Canonical quirks represented | 42 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `P01` | [Liquid-column manometer: U-tube, well-type and inclined](<gauges/p01_liquid_column_manometer_u_tube_well_type_and_inclined.md>) | Open/closed U-tubes; cistern/well manometers; inclined draft gauges | Gauge, differential or absolute pressure | 1640s-present; inclined industrial forms common from late 1800s |
| `P02` | [Mercury barometer: cistern, Fortin and wheel forms](<gauges/p02_mercury_barometer_cistern_fortin_and_wheel_forms.md>) | Torricellian barometer; Fortin adjustable cistern; wheel barometer with float and pulley | Atmospheric pressure | 1643-present; mercury versions now restricted but still historical/reference instruments |
| `P12` | [McLeod compression vacuum gauge](<gauges/p12_mcleod_compression_vacuum_gauge.md>) | Herbert McLeod gauge; laboratory mercury compression gauges | Low absolute pressure by compressing a known gas volume | 1874-c.1970s common; still used as a reference in some laboratories |
| `P17` | [Tubular, reflex and transparent sight-level gauges](<gauges/p17_tubular_reflex_and_transparent_sight_level_gauges.md>) | Boiler gauge glass; reflex prism glass; transparent double-plate level gauge | Liquid level, phase boundary or boiler water level | 1700s-present; industrial reflex forms from late 1800s |
| `P21` | [Liquid-in-glass thermometer](<gauges/p21_liquid_in_glass_thermometer.md>) | Mercury, alcohol, kerosene and spirit thermometers; clinical and industrial stem types | Temperature | 1600s-present; mercury increasingly restricted |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `P01` | [Liquid-column manometer: U-tube, well-type and inclined](<gauges/p01_liquid_column_manometer_u_tube_well_type_and_inclined.md>) | Gauge, differential or absolute pressure | 1640s-present; inclined industrial forms common from late 1800s | 14 | 2 | 1 |
| `P02` | [Mercury barometer: cistern, Fortin and wheel forms](<gauges/p02_mercury_barometer_cistern_fortin_and_wheel_forms.md>) | Atmospheric pressure | 1643-present; mercury versions now restricted but still historical/reference instruments | 9 | 2 | 1 |
| `P12` | [McLeod compression vacuum gauge](<gauges/p12_mcleod_compression_vacuum_gauge.md>) | Low absolute pressure by compressing a known gas volume | 1874-c.1970s common; still used as a reference in some laboratories | 7 | 2 | 1 |
| `P17` | [Tubular, reflex and transparent sight-level gauges](<gauges/p17_tubular_reflex_and_transparent_sight_level_gauges.md>) | Liquid level, phase boundary or boiler water level | 1700s-present; industrial reflex forms from late 1800s | 9 | 1 | 1 |
| `P21` | [Liquid-in-glass thermometer](<gauges/p21_liquid_in_glass_thermometer.md>) | Temperature | 1600s-present; mercury increasingly restricted | 11 | 1 | 1 |
| `P22` | [Six’s maximum-minimum thermometer](<gauges/p22_sixs_maximum_minimum_thermometer.md>) | Maximum and minimum temperature since reset | 1780s-present | 8 | 2 | 1 |
| `P27` | [Wet- and dry-bulb psychrometer](<gauges/p27_wet_and_dry_bulb_psychrometer.md>) | Humidity, dew point or wet-bulb temperature | 1800s-present | 6 | 2 | 1 |
| `X12` | [Steam-boiler water gauge glass](<gauges/x12_steam_boiler_water_gauge_glass.md>) | Boiler water level | 1800s-present | 6 | 1 | 1 |
| `X26` | [Goethe weather glass / water barometer](<gauges/x26_goethe_weather_glass_water_barometer.md>) | Atmospheric pressure, confounded by temperature | 1600s-present as household curiosity | 9 | 1 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `P34` | [Diving depth gauge: capillary and Bourdon](<../radial_pointer/gauges/p34_diving_depth_gauge_capillary_and_bourdon.md>) | [Radial pointer](<../radial_pointer/README.md>) |
| `P35` | [Mercury and aneroid sphygmomanometer](<../radial_pointer/gauges/p35_mercury_and_aneroid_sphygmomanometer.md>) | [Radial pointer](<../radial_pointer/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Capillary, wetting and surface-tension effects](<quirks/capillary_wetting_and_surface_tension_effects.md>) | Fluid effects arising from narrow passages, adhesion and surface forces at boundaries. | 6 | 66.67% | 6 |
| [Meniscus behaviour](<quirks/meniscus_behaviour.md>) | The curved liquid surface and its reading conventions, wetting shape and movement. | 6 | 66.67% | 6 |
| [Bubble, void and separated-column behaviour](<quirks/bubble_void_and_separated_column_behaviour.md>) | Errors or discontinuities caused by trapped gas, empty spaces or broken fluid columns. | 5 | 55.56% | 5 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 5 | 55.56% | 7 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 4 | 44.44% | 4 |
| [Calibration, correction and compensation](<quirks/calibration_correction_and_compensation.md>) | Adjustments or correction factors required to relate the raw indication to the intended quantity. | 3 | 33.33% | 4 |
| [Contamination, dirt and fouling](<quirks/contamination_dirt_and_fouling.md>) | Reading or appearance changes caused by deposits, dust, oxidation, residue or biological growth. | 3 | 33.33% | 3 |
| [Fluid surge, slosh and foam](<quirks/fluid_surge_slosh_and_foam.md>) | Transient level or pressure effects caused by fluid motion, aeration or froth. | 3 | 33.33% | 6 |
| [Parallax](<quirks/parallax.md>) | Apparent reading error caused by viewing the indicator and scale from the wrong angle or depth relationship. | 3 | 33.33% | 3 |
| [Safety, guarding and fail-safe behaviour](<quirks/safety_guarding_and_fail_safe_behaviour.md>) | Features or failure modes intended to protect the user, equipment or validity of the indication. | 3 | 33.33% | 5 |
| [Construction tolerances and unit variation](<quirks/construction_tolerances_and_unit_variation.md>) | Differences between nominally identical instruments caused by manufacturing and assembly tolerances. | 2 | 22.22% | 2 |
| [Human-factor ambiguity and misreading](<quirks/human_factor_ambiguity_and_misreading.md>) | Display features that make an otherwise functioning instrument easy to interpret incorrectly. | 2 | 22.22% | 2 |
| [Operator procedure and ritual](<quirks/operator_procedure_and_ritual.md>) | Required handling or reading practices that materially affect the result. | 2 | 22.22% | 3 |
| [Optical distortion and refraction](<quirks/optical_distortion_and_refraction.md>) | Apparent displacement or shape changes caused by lenses, glass, liquid, curved windows or refractive interfaces. | 2 | 22.22% | 2 |
| [Scale linearity and nonlinearity](<quirks/scale_linearity_and_nonlinearity.md>) | Variation in how equal input increments map to equal or unequal distances on the displayed scale. | 2 | 22.22% | 2 |
| [Witness, peak-hold and retained extrema](<quirks/witness_peak_hold_and_retained_extrema.md>) | Markers or memory mechanisms that preserve minimum, maximum or peak values after the live indication moves away. | 2 | 22.22% | 2 |
| [Airflow and ventilation dependence](<quirks/airflow_and_ventilation_dependence.md>) | Changes caused by cooling airflow, convection or ventilation through and around the instrument. | 1 | 11.11% | 2 |
| [Ambient-pressure and altitude sensitivity](<quirks/ambient_pressure_and_altitude_sensitivity.md>) | Changes in indication caused by surrounding atmospheric pressure or elevation. | 1 | 11.11% | 1 |
| [Backlash and lash](<quirks/backlash_and_lash.md>) | Lost motion when direction reverses because of clearance in gears, linkages or drive parts. | 1 | 11.11% | 1 |
| [Blockage and restricted passages](<quirks/blockage_and_restricted_passages.md>) | Errors or slow response caused by obstructed tubing, ports, capillaries or vents. | 1 | 11.11% | 1 |
| [Channel mismatch and unequal dynamics](<quirks/channel_mismatch_and_unequal_dynamics.md>) | Different calibration, response or motion between channels intended to behave alike. | 1 | 11.11% | 1 |
| [Compressed or expanded scale](<quirks/compressed_or_expanded_scale.md>) | A deliberately or accidentally nonuniform scale that gives some ranges more display space than others. | 1 | 11.11% | 1 |
| [Density and viscosity dependence](<quirks/density_and_viscosity_dependence.md>) | Changes in indication caused by fluid density, viscosity or their variation with conditions. | 1 | 11.11% | 1 |
| [End stops, pegging and overflow](<quirks/end_stops_pegging_and_overflow.md>) | Behaviour at physical or representational range limits, including contact with stops and beyond-range indications. | 1 | 11.11% | 1 |
| [Evaporation and drying](<quirks/evaporation_and_drying.md>) | Loss of liquid or solvent that changes level, concentration, response or appearance. | 1 | 11.11% | 1 |
| [Gas or medium dependence](<quirks/gas_or_medium_dependence.md>) | Dependence on the type, pressure, composition or condition of the working gas or surrounding medium. | 1 | 11.11% | 1 |
| [Gravity-related behaviour](<quirks/gravity_related_behaviour.md>) | Dependence on local gravity magnitude or direction. | 1 | 11.11% | 1 |
| [Illumination and backlighting](<quirks/illumination_and_backlighting.md>) | Lighting systems that make scales, legends or display elements visible and their associated unevenness or ageing. | 1 | 11.11% | 1 |
| [Immersion, stem-conduction and contact error](<quirks/immersion_stem_conduction_and_contact_error.md>) | Temperature-reading errors caused by installation depth, heat flow along a probe or imperfect thermal contact. | 1 | 11.11% | 1 |
| [Manual adjustment and setup](<quirks/manual_adjustment_and_setup.md>) | User controls and preparation steps needed to configure the instrument before use. | 1 | 11.11% | 1 |
| [Manual reset, tare and caging](<quirks/manual_reset_tare_and_caging.md>) | User-operated mechanisms that clear, zero, restrain or protect the indication. | 1 | 11.11% | 1 |
| [Orientation sensitivity](<quirks/orientation_sensitivity.md>) | Dependence on instrument mounting angle or orientation relative to gravity or fields. | 1 | 11.11% | 1 |
| [Qualitative or non-precision indication](<quirks/qualitative_or_non_precision_indication.md>) | An indication intended for trend, state or rough comparison rather than accurate numeric measurement. | 1 | 11.11% | 1 |
| [Readout masking and narrow-window transitions](<quirks/readout_masking_and_narrow_window_transitions.md>) | Partial or ambiguous readings caused by apertures that reveal only a small portion of a moving scale or drum. | 1 | 11.11% | 1 |
| [Reference bugs and set markers](<quirks/reference_bugs_and_set_markers.md>) | Movable or fixed markers used to record targets, limits, headings or comparison values. | 1 | 11.11% | 1 |
| [Ringing and oscillation](<quirks/ringing_and_oscillation.md>) | Repeated decaying or sustained motion around a target following excitation or disturbance. | 1 | 11.11% | 1 |
| [Scale direction and interpretation](<quirks/scale_direction_and_interpretation.md>) | Whether increasing values move clockwise, anticlockwise, upward or otherwise, and how that direction is understood. | 1 | 11.11% | 2 |
| [Sensor plumbing faults](<quirks/sensor_plumbing_faults.md>) | Errors caused by incorrectly connected, leaking, reversed or contaminated pressure and fluid lines. | 1 | 11.11% | 1 |
| [Settling and return behaviour](<quirks/settling_and_return_behaviour.md>) | How the indication approaches a stable value or returns after a transient, disturbance or release. | 1 | 11.11% | 1 |
| [Shadows, depth and occlusion](<quirks/shadows_depth_and_occlusion.md>) | Visual effects caused by layered parts blocking, shading or appearing at different depths. | 1 | 11.11% | 1 |
| [Stiction and sticking](<quirks/stiction_and_sticking.md>) | Static friction or adhesion that prevents motion until enough force accumulates, often followed by a jump. | 1 | 11.11% | 1 |
| [Tapping](<quirks/tapping.md>) | Intentional tapping used to release friction, settle the mechanism or obtain a representative reading. | 1 | 11.11% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
