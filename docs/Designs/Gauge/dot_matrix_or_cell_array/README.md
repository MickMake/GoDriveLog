---
gauge_group: dot_matrix_or_cell_array
catalogue_version: "0.2"
primary_gauge_count: 5
supporting_quirk_count: 27
---

# Dot-matrix or cell array

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

An addressable two-dimensional array of dots, cells or pixels forms characters, symbols or graphic fields without relying on fixed character strokes.

**Catalogue definition:** An addressable matrix of dots or cells forms characters, symbols or fields.

## How the group encodes a value

Information is encoded by the state of independently or row/column-addressed cells.

## Classification boundary

Use this group for matrix-addressed fields. Fixed seven-, fourteen- or sixteen-segment characters belong under segmented_display.

## Simulation baseline

Include addressing cadence, row/column artefacts, response time, crosstalk, viewing angle, retained state and defective cells.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 5 |
| Share of catalogue | 3.68% |
| Alternate members | 1 |
| Canonical quirks represented | 27 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `D16` | [LED dot-matrix and smart alphanumeric module](<gauges/d16_led_dot_matrix_and_smart_alphanumeric_module.md>) | 5x7 matrices; HP intelligent displays; scrolling message modules | Text, symbols, plots and coarse gauges | 1970s-present |
| `D17` | [Monochrome storage plasma panel](<gauges/d17_monochrome_storage_plasma_panel.md>) | University of Illinois PLATO panels; Owens-Illinois and IBM orange plasma displays | Text, graphics, terminals and instrument pages | 1964-c.1990s; specialist descendants later |
| `D21` | [STN/FSTN passive-matrix LCD](<gauges/d21_stn_fstn_passive_matrix_lcd.md>) | Industrial HMIs, telephones, test instruments and vehicle displays | Text, graphics and segmented or dot-matrix gauges | 1980s-present; peak c.1990s-2000s |
| `D22` | [Cholesteric and other bistable LCD](<gauges/d22_cholesteric_and_other_bistable_lcd.md>) | Kent Displays writing panels; reflective shelf labels and status indicators | Persistent text, symbols or status with near-zero holding power | 1990s-present |
| `D29` | [Electrowetting and interferometric MEMS display](<gauges/d29_electrowetting_and_interferometric_mems_display.md>) | Liquavista electrowetting; Qualcomm Mirasol interferometric modulator | Reflective low-power text, colour graphics or segmented indicators | 2000s-present; limited commercial adoption |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `D16` | [LED dot-matrix and smart alphanumeric module](<gauges/d16_led_dot_matrix_and_smart_alphanumeric_module.md>) | Text, symbols, plots and coarse gauges | 1970s-present | 7 | 2 | 1 |
| `D17` | [Monochrome storage plasma panel](<gauges/d17_monochrome_storage_plasma_panel.md>) | Text, graphics, terminals and instrument pages | 1964-c.1990s; specialist descendants later | 9 | 2 | 1 |
| `D21` | [STN/FSTN passive-matrix LCD](<gauges/d21_stn_fstn_passive_matrix_lcd.md>) | Text, graphics and segmented or dot-matrix gauges | 1980s-present; peak c.1990s-2000s | 7 | 1 | 1 |
| `D22` | [Cholesteric and other bistable LCD](<gauges/d22_cholesteric_and_other_bistable_lcd.md>) | Persistent text, symbols or status with near-zero holding power | 1990s-present | 6 | 1 | 1 |
| `D29` | [Electrowetting and interferometric MEMS display](<gauges/d29_electrowetting_and_interferometric_mems_display.md>) | Reflective low-power text, colour graphics or segmented indicators | 2000s-present; limited commercial adoption | 5 | 2 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `D18` | [Powder and thin-film electroluminescent display](<../colour_or_material_state/gauges/d18_powder_and_thin_film_electroluminescent_display.md>) | [Colour or material-state indicator](<../colour_or_material_state/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Colour variation and colour shift](<quirks/colour_variation_and_colour_shift.md>) | Changes or inconsistencies in displayed colour across level, age, temperature, angle or individual units. | 3 | 60.00% | 4 |
| [Viewing-angle dependence](<quirks/viewing_angle_dependence.md>) | Changes in readability, colour, contrast or apparent value with observer angle. | 3 | 60.00% | 5 |
| [Bistable display memory and update](<quirks/bistable_display_memory_and_update.md>) | Displays whose cells retain state without continuous power and require explicit transitions to update. | 2 | 40.00% | 3 |
| [Latching and state-retention behaviour](<quirks/latching_and_state_retention_behaviour.md>) | A displayed state that remains mechanically, magnetically, electrically or optically retained until reset or rewritten. | 2 | 40.00% | 2 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 2 | 40.00% | 2 |
| [Bloom, halo and penumbra](<quirks/bloom_halo_and_penumbra.md>) | Spreading or soft-edged light around a spot, trace, segment or projected image. | 1 | 20.00% | 1 |
| [Brightness and contrast variation](<quirks/brightness_and_contrast_variation.md>) | Changes in luminance or visual separation between active and inactive parts of the display. | 1 | 20.00% | 1 |
| [Brightness nonuniformity and gradients](<quirks/brightness_nonuniformity_and_gradients.md>) | Unequal luminance across a display, character, segment, tube or field. | 1 | 20.00% | 1 |
| [Burn-in and permanent image wear](<quirks/burn_in_and_permanent_image_wear.md>) | Lasting visible damage or nonuniform ageing from repeatedly displaying the same pattern. | 1 | 20.00% | 1 |
| [Calibration, correction and compensation](<quirks/calibration_correction_and_compensation.md>) | Adjustments or correction factors required to relate the raw indication to the intended quantity. | 1 | 20.00% | 1 |
| [Case, lens and enclosure deformation](<quirks/case_lens_and_enclosure_deformation.md>) | Reading or visual changes caused by pressure, heat, stress or damage deforming the enclosure or window. | 1 | 20.00% | 1 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 1 | 20.00% | 1 |
| [Electrical hum, buzz and whine](<quirks/electrical_hum_buzz_and_whine.md>) | Audible vibration or tone produced by electrical drive, magnetic parts or switching. | 1 | 20.00% | 1 |
| [Flicker, scan and PWM artefacts](<quirks/flicker_scan_and_pwm_artefacts.md>) | Visible modulation caused by multiplexing, scanning, pulse-width control or interaction with cameras and eye motion. | 1 | 20.00% | 2 |
| [Ghosting, crosstalk and light leakage](<quirks/ghosting_crosstalk_and_light_leakage.md>) | Unwanted partial activation or illumination of neighbouring, previous or nominally inactive display elements. | 1 | 20.00% | 1 |
| [Glyph, segment and aperture geometry](<quirks/glyph_segment_and_aperture_geometry.md>) | The shapes and proportions of characters, segments, masks and viewing windows. | 1 | 20.00% | 1 |
| [Meniscus behaviour](<quirks/meniscus_behaviour.md>) | The curved liquid surface and its reading conventions, wetting shape and movement. | 1 | 20.00% | 1 |
| [Persistence and afterglow](<quirks/persistence_and_afterglow.md>) | Continued visibility after excitation is reduced or removed. | 1 | 20.00% | 1 |
| [Phosphor behaviour and ageing](<quirks/phosphor_behaviour_and_ageing.md>) | Brightness, colour, persistence and degradation characteristics of phosphor-based displays. | 1 | 20.00% | 1 |
| [Power-up and self-test behaviour](<quirks/power_up_and_self_test_behaviour.md>) | Visible startup sequences, checks, sweeps or initial states used when power is applied. | 1 | 20.00% | 1 |
| [Random, statistical and batch variation](<quirks/random_statistical_and_batch_variation.md>) | Unpredictable events or unit-to-unit differences that are part of the observed behaviour. | 1 | 20.00% | 1 |
| [Refresh, erase and update artefacts](<quirks/refresh_erase_and_update_artefacts.md>) | Visible effects produced while changing, clearing or rewriting a display. | 1 | 20.00% | 1 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 1 | 20.00% | 1 |
| [Scope or implementation boundary](<quirks/scope_or_implementation_boundary.md>) | A catalogue note that distinguishes the represented mechanism from excluded or more generic implementations. | 1 | 20.00% | 1 |
| [Segment, pixel, lamp or flag failure](<quirks/segment_pixel_lamp_or_flag_failure.md>) | Individual display elements that fail open, fail active, stick, weaken or respond intermittently. | 1 | 20.00% | 1 |
| [Shared-current or shared-supply dimming](<quirks/shared_current_or_shared_supply_dimming.md>) | Brightness reduction or interaction when multiple elements draw from a common limited current or supply. | 1 | 20.00% | 1 |
| [Stiction and sticking](<quirks/stiction_and_sticking.md>) | Static friction or adhesion that prevents motion until enough force accumulates, often followed by a jump. | 1 | 20.00% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
