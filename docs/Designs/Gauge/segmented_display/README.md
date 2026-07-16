---
gauge_group: segmented_display
catalogue_version: "0.2"
primary_gauge_count: 13
supporting_quirk_count: 38
---

# Segmented display

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

Fixed luminous, gas-discharge, liquid-crystal, electrochromic or emissive segments combine to form numerals and limited glyphs.

**Catalogue definition:** Fixed luminous, gas-discharge, liquid-crystal or material segments form numbers or glyphs.

## How the group encodes a value

A predefined set of segments is activated in patterns; the geometry limits which characters can be represented.

## Classification boundary

Use this group when character strokes are fixed. An addressable grid belongs under dot_matrix_or_cell_array.

## Simulation baseline

Represent segment geometry, drive method, multiplexing, brightness interactions, viewing angle, persistence, ageing and partial failures.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 13 |
| Share of catalogue | 9.56% |
| Alternate members | 1 |
| Canonical quirks represented | 38 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `D05` | [Nixie numeric cold-cathode tube](<gauges/d05_nixie_numeric_cold_cathode_tube.md>) | Burroughs NIXIE; IN-12, IN-14, ZM1040 and NL-5441 | Numerals, decimal points and occasional symbols | c.1955-c.1980 mainstream; boutique revival today |
| `D06` | [Alphanumeric Nixie tube](<gauges/d06_alphanumeric_nixie_tube.md>) | Burroughs B-7971 and B-8971; segmented or shaped-cathode alphabet tubes | Letters, numerals and symbols | 1960s-c.1980s; collector use today |
| `D07` | [Pixie / top-view cold-cathode display](<gauges/d07_pixie_top_view_cold_cathode_display.md>) | Philips ZM1050 and ZM1051; Burroughs Pixie variants | Numerals and symbols | c.1960s-c.1970s |
| `D08` | [Panaplex planar gas-discharge display](<gauges/d08_panaplex_planar_gas_discharge_display.md>) | Burroughs Panaplex II; Sperry planar numeric panels | Multi-digit numeric and segmented readouts | 1969-c.1980s common |
| `D11` | [Numitron / Minitron incandescent segment display](<gauges/d11_numitron_minitron_incandescent_segment_display.md>) | RCA Numitron; IEE and Apollo-era Minitron; DR2000/DR2110 | Numeric or limited alphanumeric display | c.1968-c.1980s common; collector use today |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `D05` | [Nixie numeric cold-cathode tube](<gauges/d05_nixie_numeric_cold_cathode_tube.md>) | Numerals, decimal points and occasional symbols | c.1955-c.1980 mainstream; boutique revival today | 9 | 2 | 1 |
| `D06` | [Alphanumeric Nixie tube](<gauges/d06_alphanumeric_nixie_tube.md>) | Letters, numerals and symbols | 1960s-c.1980s; collector use today | 7 | 1 | 1 |
| `D07` | [Pixie / top-view cold-cathode display](<gauges/d07_pixie_top_view_cold_cathode_display.md>) | Numerals and symbols | c.1960s-c.1970s | 7 | 2 | 1 |
| `D08` | [Panaplex planar gas-discharge display](<gauges/d08_panaplex_planar_gas_discharge_display.md>) | Multi-digit numeric and segmented readouts | 1969-c.1980s common | 8 | 1 | 1 |
| `D11` | [Numitron / Minitron incandescent segment display](<gauges/d11_numitron_minitron_incandescent_segment_display.md>) | Numeric or limited alphanumeric display | c.1968-c.1980s common; collector use today | 9 | 2 | 1 |
| `D12` | [Vacuum fluorescent display](<gauges/d12_vacuum_fluorescent_display.md>) | Numeric, alphanumeric, bargraph and custom icons | 1967-present | 10 | 2 | 1 |
| `D13` | [LED seven-segment display](<gauges/d13_led_seven_segment_display.md>) | Numeric readout | late 1960s-present | 8 | 3 | 1 |
| `D14` | [LED 14- and 16-segment starburst display](<gauges/d14_led_14_and_16_segment_starburst_display.md>) | Letters, numerals and symbols | 1970s-present | 8 | 2 | 1 |
| `D19` | [Dynamic-scattering-mode LCD](<gauges/d19_dynamic_scattering_mode_lcd.md>) | Numeric and simple segmented readout | 1968-c.1975 commercially | 8 | 2 | 1 |
| `D20` | [Twisted-nematic segmented LCD](<gauges/d20_twisted_nematic_segmented_lcd.md>) | Numeric, icons, bars and fixed legends | 1971-present | 10 | 2 | 1 |
| `D23` | [Segmented electrophoretic / e-paper display](<gauges/d23_segmented_electrophoretic_e_paper_display.md>) | Persistent numbers, icons, bars and labels | late 1990s-present | 8 | 2 | 1 |
| `D24` | [Electrochromic segmented display](<gauges/d24_electrochromic_segmented_display.md>) | Persistent or semi-persistent symbols, numbers and tint levels | 1970s-present; niche | 9 | 1 | 1 |
| `D25` | [Segmented OLED and passive-matrix OLED](<gauges/d25_segmented_oled_and_passive_matrix_oled.md>) | Numbers, icons, text and bar graphs | late 1990s-present | 7 | 2 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `P16` | [Capacitance diaphragm manometer](<../radial_pointer/gauges/p16_capacitance_diaphragm_manometer.md>) | [Radial pointer](<../radial_pointer/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Ghosting, crosstalk and light leakage](<quirks/ghosting_crosstalk_and_light_leakage.md>) | Unwanted partial activation or illumination of neighbouring, previous or nominally inactive display elements. | 8 | 61.54% | 8 |
| [Brightness and contrast variation](<quirks/brightness_and_contrast_variation.md>) | Changes in luminance or visual separation between active and inactive parts of the display. | 7 | 53.85% | 7 |
| [Colour variation and colour shift](<quirks/colour_variation_and_colour_shift.md>) | Changes or inconsistencies in displayed colour across level, age, temperature, angle or individual units. | 7 | 53.85% | 10 |
| [Ageing and material degradation](<quirks/ageing_and_material_degradation.md>) | Progressive change due to wear, fatigue, chemical decay, phosphor loss, embrittlement or similar ageing processes. | 6 | 46.15% | 8 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 6 | 46.15% | 7 |
| [Glyph, segment and aperture geometry](<quirks/glyph_segment_and_aperture_geometry.md>) | The shapes and proportions of characters, segments, masks and viewing windows. | 5 | 38.46% | 10 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 5 | 38.46% | 5 |
| [Brightness nonuniformity and gradients](<quirks/brightness_nonuniformity_and_gradients.md>) | Unequal luminance across a display, character, segment, tube or field. | 4 | 30.77% | 6 |
| [Flicker, scan and PWM artefacts](<quirks/flicker_scan_and_pwm_artefacts.md>) | Visible modulation caused by multiplexing, scanning, pulse-width control or interaction with cameras and eye motion. | 4 | 30.77% | 4 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 4 | 30.77% | 4 |
| [Shadows, depth and occlusion](<quirks/shadows_depth_and_occlusion.md>) | Visual effects caused by layered parts blocking, shading or appearing at different depths. | 4 | 30.77% | 6 |
| [Bloom, halo and penumbra](<quirks/bloom_halo_and_penumbra.md>) | Spreading or soft-edged light around a spot, trace, segment or projected image. | 3 | 23.08% | 3 |
| [Gas-discharge striking and ignition](<quirks/gas_discharge_striking_and_ignition.md>) | The voltage, timing and transient behaviour required to initiate a gas discharge. | 3 | 23.08% | 3 |
| [Segment, pixel, lamp or flag failure](<quirks/segment_pixel_lamp_or_flag_failure.md>) | Individual display elements that fail open, fail active, stick, weaken or respond intermittently. | 3 | 23.08% | 3 |
| [Shared-current or shared-supply dimming](<quirks/shared_current_or_shared_supply_dimming.md>) | Brightness reduction or interaction when multiple elements draw from a common limited current or supply. | 3 | 23.08% | 3 |
| [Viewing-angle dependence](<quirks/viewing_angle_dependence.md>) | Changes in readability, colour, contrast or apparent value with observer angle. | 3 | 23.08% | 3 |
| [Warm-up behaviour](<quirks/warm_up_behaviour.md>) | Transient changes after startup while temperature, discharge, illumination or mechanics stabilise. | 3 | 23.08% | 3 |
| [Ambient-light readability](<quirks/ambient_light_readability.md>) | Dependence of legibility on sunlight, darkness, glare or surrounding illumination. | 2 | 15.38% | 2 |
| [Burn-in and permanent image wear](<quirks/burn_in_and_permanent_image_wear.md>) | Lasting visible damage or nonuniform ageing from repeatedly displaying the same pattern. | 2 | 15.38% | 3 |
| [Cathode poisoning](<quirks/cathode_poisoning.md>) | Loss of emission or uneven operation in cathodes that are underused or operated under unsuitable conditions. | 2 | 15.38% | 2 |
| [Filament behaviour and failure](<quirks/filament_behaviour_and_failure.md>) | Warm-up, sag, resistance change, brightness variation and breakage of incandescent filaments. | 2 | 15.38% | 6 |
| [Latching and state-retention behaviour](<quirks/latching_and_state_retention_behaviour.md>) | A displayed state that remains mechanically, magnetically, electrically or optically retained until reset or rewritten. | 2 | 15.38% | 2 |
| [Parallax](<quirks/parallax.md>) | Apparent reading error caused by viewing the indicator and scale from the wrong angle or depth relationship. | 2 | 15.38% | 2 |
| [Phosphor behaviour and ageing](<quirks/phosphor_behaviour_and_ageing.md>) | Brightness, colour, persistence and degradation characteristics of phosphor-based displays. | 2 | 15.38% | 3 |
| [Power-off and power-loss behaviour](<quirks/power_off_and_power_loss_behaviour.md>) | What the indication does when drive power is removed, including retained, blank, parked or misleading states. | 2 | 15.38% | 2 |
| [Sputter, haze and deposits](<quirks/sputter_haze_and_deposits.md>) | Material deposition or internal contamination that clouds a window or changes discharge behaviour. | 2 | 15.38% | 2 |
| [Bistable display memory and update](<quirks/bistable_display_memory_and_update.md>) | Displays whose cells retain state without continuous power and require explicit transitions to update. | 1 | 7.69% | 1 |
| [Channel mismatch and unequal dynamics](<quirks/channel_mismatch_and_unequal_dynamics.md>) | Different calibration, response or motion between channels intended to behave alike. | 1 | 7.69% | 1 |
| [Chemical diffusion and reaction front](<quirks/chemical_diffusion_and_reaction_front.md>) | Movement or spread of a chemical change through a material used as the indication. | 1 | 7.69% | 1 |
| [Cool-down, heat soak and continued motion](<quirks/cool_down_heat_soak_and_continued_motion.md>) | Behaviour after drive or heat is removed while stored thermal or mechanical energy remains. | 1 | 7.69% | 1 |
| [Display overlap and adjacent indication](<quirks/display_overlap_and_adjacent_indication.md>) | Visual interference where neighbouring pointers, digits, traces or display regions overlap. | 1 | 7.69% | 1 |
| [Drift and long-term stability](<quirks/drift_and_long_term_stability.md>) | Slow change in indicated value or behaviour over time despite an unchanged input. | 1 | 7.69% | 1 |
| [Drive power and regulation](<quirks/drive_power_and_regulation.md>) | Dependence on the available drive energy and the quality of its regulation. | 1 | 7.69% | 1 |
| [Electrical contacts, grounds and wiring](<quirks/electrical_contacts_grounds_and_wiring.md>) | Display faults caused by contact resistance, grounding, broken conductors, polarity or connection layout. | 1 | 7.69% | 1 |
| [Flutter, jitter, tremor and quiver](<quirks/flutter_jitter_tremor_and_quiver.md>) | Small rapid random or periodic movements around the nominal indication. | 1 | 7.69% | 1 |
| [Frequency, waveform and source dependence](<quirks/frequency_waveform_and_source_dependence.md>) | Changes in reading or response caused by signal frequency, waveform shape or source characteristics. | 1 | 7.69% | 1 |
| [Illumination and backlighting](<quirks/illumination_and_backlighting.md>) | Lighting systems that make scales, legends or display elements visible and their associated unevenness or ageing. | 1 | 7.69% | 1 |
| [Refresh, erase and update artefacts](<quirks/refresh_erase_and_update_artefacts.md>) | Visible effects produced while changing, clearing or rewriting a display. | 1 | 7.69% | 2 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
