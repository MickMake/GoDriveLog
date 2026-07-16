---
gauge_group: flip_element_array
catalogue_version: "0.2"
primary_gauge_count: 2
supporting_quirk_count: 15
---

# Flip-element array

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

An array of bistable discs, dots, tiles or flags physically flips between contrasting faces to form symbols, levels or fields.

**Catalogue definition:** Bistable discs, dots, tiles or flags physically flip to form symbols or levels.

## How the group encodes a value

Each element retains one of two visible states, allowing an image to persist without continuous drive power.

## Classification boundary

Use this group for arrays of physical bistable elements. A single flag belongs under mechanical_flag_or_shutter.

## Simulation baseline

Update elements individually, including travel time, sound, unequal timing, stuck elements and retained state after power loss.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 2 |
| Share of catalogue | 1.47% |
| Alternate members | 0 |
| Canonical quirks represented | 15 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `E20` | [Flip-dot / flip-disc display](<gauges/e20_flip_dot_flip_disc_display.md>) | Ferranti-Packard and Luminator transit signs; AlfaZeta modules | Text, numerals, symbols, bar graphs or mimic diagrams | 1960s-present |
| `P19` | [Magnetic level indicator with flip flags](<gauges/p19_magnetic_level_indicator_with_flip_flags.md>) | External chamber, magnetic float and red/white or black/yellow rotating flags | Liquid level or interface level | mid-1900s-present |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `E20` | [Flip-dot / flip-disc display](<gauges/e20_flip_dot_flip_disc_display.md>) | Text, numerals, symbols, bar graphs or mimic diagrams | 1960s-present | 11 | 2 | 1 |
| `P19` | [Magnetic level indicator with flip flags](<gauges/p19_magnetic_level_indicator_with_flip_flags.md>) | Liquid level or interface level | mid-1900s-present | 7 | 1 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 2 | 100.00% | 3 |
| [Stiction and sticking](<quirks/stiction_and_sticking.md>) | Static friction or adhesion that prevents motion until enough force accumulates, often followed by a jump. | 2 | 100.00% | 3 |
| [Viewing-angle dependence](<quirks/viewing_angle_dependence.md>) | Changes in readability, colour, contrast or apparent value with observer angle. | 2 | 100.00% | 2 |
| [Ambient-light readability](<quirks/ambient_light_readability.md>) | Dependence of legibility on sunlight, darkness, glare or surrounding illumination. | 1 | 50.00% | 1 |
| [Bistable display memory and update](<quirks/bistable_display_memory_and_update.md>) | Displays whose cells retain state without continuous power and require explicit transitions to update. | 1 | 50.00% | 1 |
| [Brightness and contrast variation](<quirks/brightness_and_contrast_variation.md>) | Changes in luminance or visual separation between active and inactive parts of the display. | 1 | 50.00% | 1 |
| [Flicker, scan and PWM artefacts](<quirks/flicker_scan_and_pwm_artefacts.md>) | Visible modulation caused by multiplexing, scanning, pulse-width control or interaction with cameras and eye motion. | 1 | 50.00% | 1 |
| [Float and buoyancy behaviour](<quirks/float_and_buoyancy_behaviour.md>) | Motion and equilibrium effects arising from buoyant indicators in fluids. | 1 | 50.00% | 2 |
| [Latching and state-retention behaviour](<quirks/latching_and_state_retention_behaviour.md>) | A displayed state that remains mechanically, magnetically, electrically or optically retained until reset or rewritten. | 1 | 50.00% | 1 |
| [Magnetic-field, deviation and remanence effects](<quirks/magnetic_field_deviation_and_remanence_effects.md>) | Influence of external or retained magnetism on indication, zero and calibration. | 1 | 50.00% | 1 |
| [Mechanical noise and cadence](<quirks/mechanical_noise_and_cadence.md>) | Audible clicks, hums, impacts or rhythms produced by the display mechanism. | 1 | 50.00% | 1 |
| [Power-off and power-loss behaviour](<quirks/power_off_and_power_loss_behaviour.md>) | What the indication does when drive power is removed, including retained, blank, parked or misleading states. | 1 | 50.00% | 1 |
| [Refresh, erase and update artefacts](<quirks/refresh_erase_and_update_artefacts.md>) | Visible effects produced while changing, clearing or rewriting a display. | 1 | 50.00% | 1 |
| [Segment, pixel, lamp or flag failure](<quirks/segment_pixel_lamp_or_flag_failure.md>) | Individual display elements that fail open, fail active, stick, weaken or respond intermittently. | 1 | 50.00% | 1 |
| [Shadows, depth and occlusion](<quirks/shadows_depth_and_occlusion.md>) | Visual effects caused by layered parts blocking, shading or appearing at different depths. | 1 | 50.00% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
