---
gauge_group: indicator_lamp
catalogue_version: "0.2"
primary_gauge_count: 3
supporting_quirk_count: 18
---

# Indicator lamp or illuminated legend

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A lamp, lens, illuminated legend or light pipe indicates a discrete state, warning or coarse condition.

**Catalogue definition:** A lamp, lens or illuminated legend indicates a state or coarse condition.

## How the group encodes a value

Information is encoded by on/off state, colour, brightness, flash cadence or which legend is illuminated.

## Classification boundary

Use this group for individual lamps and legends. Arrays that form arbitrary glyphs belong under dot_matrix_or_cell_array or segmented_display.

## Simulation baseline

Include warm-up, filament fade, lens diffusion, ambient-light washout, supply variation, flash timing and failed lamps.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 3 |
| Share of catalogue | 2.21% |
| Alternate members | 0 |
| Canonical quirks represented | 18 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `D01` | [Incandescent pilot lamp and annunciator](<gauges/d01_incandescent_pilot_lamp_and_annunciator.md>) | Miniature bayonet bulbs; aircraft post lights; jewel lamps; warning panels | Binary status, alarm, backlight or legend illumination | 1880s-present |
| `D02` | [Edge-lit engraved acrylic and light-pipe display](<gauges/d02_edge_lit_engraved_acrylic_and_light_pipe_display.md>) | Aircraft edge-lit panels; engraved Perspex legends; automotive light pipes | Fixed legends, scales, pointers and annunciators | 1930s-present |
| `X29` | [Automotive tell-tale / “idiot light” cluster](<gauges/x29_automotive_tell_tale_idiot_light_cluster.md>) | Oil-pressure, charge, coolant and brake warning lamps; bulb-check circuits | Binary fault or threshold state | 1930s-present |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `D01` | [Incandescent pilot lamp and annunciator](<gauges/d01_incandescent_pilot_lamp_and_annunciator.md>) | Binary status, alarm, backlight or legend illumination | 1880s-present | 13 | 1 | 1 |
| `D02` | [Edge-lit engraved acrylic and light-pipe display](<gauges/d02_edge_lit_engraved_acrylic_and_light_pipe_display.md>) | Fixed legends, scales, pointers and annunciators | 1930s-present | 5 | 1 | 1 |
| `X29` | [Automotive tell-tale / “idiot light” cluster](<gauges/x29_automotive_tell_tale_idiot_light_cluster.md>) | Binary fault or threshold state | 1930s-present | 10 | 1 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Ageing and material degradation](<quirks/ageing_and_material_degradation.md>) | Progressive change due to wear, fatigue, chemical decay, phosphor loss, embrittlement or similar ageing processes. | 3 | 100.00% | 3 |
| [Segment, pixel, lamp or flag failure](<quirks/segment_pixel_lamp_or_flag_failure.md>) | Individual display elements that fail open, fail active, stick, weaken or respond intermittently. | 3 | 100.00% | 3 |
| [Brightness and contrast variation](<quirks/brightness_and_contrast_variation.md>) | Changes in luminance or visual separation between active and inactive parts of the display. | 2 | 66.67% | 2 |
| [Brightness nonuniformity and gradients](<quirks/brightness_nonuniformity_and_gradients.md>) | Unequal luminance across a display, character, segment, tube or field. | 2 | 66.67% | 3 |
| [Colour variation and colour shift](<quirks/colour_variation_and_colour_shift.md>) | Changes or inconsistencies in displayed colour across level, age, temperature, angle or individual units. | 2 | 66.67% | 2 |
| [Ghosting, crosstalk and light leakage](<quirks/ghosting_crosstalk_and_light_leakage.md>) | Unwanted partial activation or illumination of neighbouring, previous or nominally inactive display elements. | 2 | 66.67% | 2 |
| [Supply-voltage sensitivity](<quirks/supply_voltage_sensitivity.md>) | Changes in reading, brightness or dynamics caused by variations in supply voltage. | 2 | 66.67% | 2 |
| [Warm-up behaviour](<quirks/warm_up_behaviour.md>) | Transient changes after startup while temperature, discharge, illumination or mechanics stabilise. | 2 | 66.67% | 2 |
| [Bloom, halo and penumbra](<quirks/bloom_halo_and_penumbra.md>) | Spreading or soft-edged light around a spot, trace, segment or projected image. | 1 | 33.33% | 1 |
| [Cool-down, heat soak and continued motion](<quirks/cool_down_heat_soak_and_continued_motion.md>) | Behaviour after drive or heat is removed while stored thermal or mechanical energy remains. | 1 | 33.33% | 1 |
| [Electrical contacts, grounds and wiring](<quirks/electrical_contacts_grounds_and_wiring.md>) | Display faults caused by contact resistance, grounding, broken conductors, polarity or connection layout. | 1 | 33.33% | 1 |
| [Filament behaviour and failure](<quirks/filament_behaviour_and_failure.md>) | Warm-up, sag, resistance change, brightness variation and breakage of incandescent filaments. | 1 | 33.33% | 1 |
| [Flicker, scan and PWM artefacts](<quirks/flicker_scan_and_pwm_artefacts.md>) | Visible modulation caused by multiplexing, scanning, pulse-width control or interaction with cameras and eye motion. | 1 | 33.33% | 1 |
| [Illumination and backlighting](<quirks/illumination_and_backlighting.md>) | Lighting systems that make scales, legends or display elements visible and their associated unevenness or ageing. | 1 | 33.33% | 2 |
| [Power-up and self-test behaviour](<quirks/power_up_and_self_test_behaviour.md>) | Visible startup sequences, checks, sweeps or initial states used when power is applied. | 1 | 33.33% | 1 |
| [Random, statistical and batch variation](<quirks/random_statistical_and_batch_variation.md>) | Unpredictable events or unit-to-unit differences that are part of the observed behaviour. | 1 | 33.33% | 1 |
| [Shared-current or shared-supply dimming](<quirks/shared_current_or_shared_supply_dimming.md>) | Brightness reduction or interaction when multiple elements draw from a common limited current or supply. | 1 | 33.33% | 1 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 1 | 33.33% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
