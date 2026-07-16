---
gauge_group: bar_or_wedge_display
catalogue_version: "0.2"
primary_gauge_count: 3
supporting_quirk_count: 22
---

# Bar, column, wedge or moving-dot display

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

Magnitude is shown as the length, width, closure, count or moving position of a bar, column, wedge or dot sequence.

**Catalogue definition:** Magnitude is shown as a continuous or discrete bar, wedge, column or advancing dot.

## How the group encodes a value

Value is encoded spatially rather than as a full numeral, often making trends and limits quick to perceive.

## Classification boundary

Use this group for one-dimensional magnitude displays. A general two-dimensional matrix belongs under dot_matrix_or_cell_array.

## Simulation baseline

Model discrete versus continuous motion, threshold spacing, overlap, nonuniform brightness, persistence and endpoint behaviour.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 3 |
| Share of catalogue | 2.21% |
| Alternate members | 0 |
| Canonical quirks represented | 22 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `D04` | [Electron-ray “magic eye” tube](<gauges/d04_electron_ray_magic_eye_tube.md>) | 6E5, EM34, EM84 and similar tuning indicators | Radio tuning, signal level, recording level or balance | 1935-c.1970s common; collector use today |
| `D10` | [Linear neon glow bar](<gauges/d10_linear_neon_glow_bar.md>) | Soviet IN-9 and IN-13; Burroughs linear indicators | Analog level, tuning, audio level or process value | c.1960s-c.1980s; hobby use today |
| `D15` | [LED bargraph and moving-dot display](<gauges/d15_led_bargraph_and_moving_dot_display.md>) | LM3914/LM3915 meters; automotive bar gauges; audio peak bars | Level, percentage, frequency bands or threshold state | 1970s-present |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `D04` | [Electron-ray “magic eye” tube](<gauges/d04_electron_ray_magic_eye_tube.md>) | Radio tuning, signal level, recording level or balance | 1935-c.1970s common; collector use today | 8 | 2 | 1 |
| `D10` | [Linear neon glow bar](<gauges/d10_linear_neon_glow_bar.md>) | Analog level, tuning, audio level or process value | c.1960s-c.1980s; hobby use today | 7 | 1 | 1 |
| `D15` | [LED bargraph and moving-dot display](<gauges/d15_led_bargraph_and_moving_dot_display.md>) | Level, percentage, frequency bands or threshold state | 1970s-present | 10 | 2 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Brightness and contrast variation](<quirks/brightness_and_contrast_variation.md>) | Changes in luminance or visual separation between active and inactive parts of the display. | 2 | 66.67% | 3 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 2 | 66.67% | 4 |
| [Scale linearity and nonlinearity](<quirks/scale_linearity_and_nonlinearity.md>) | Variation in how equal input increments map to equal or unequal distances on the displayed scale. | 2 | 66.67% | 2 |
| [Ambient-light readability](<quirks/ambient_light_readability.md>) | Dependence of legibility on sunlight, darkness, glare or surrounding illumination. | 1 | 33.33% | 1 |
| [Burn-in and permanent image wear](<quirks/burn_in_and_permanent_image_wear.md>) | Lasting visible damage or nonuniform ageing from repeatedly displaying the same pattern. | 1 | 33.33% | 1 |
| [Chatter](<quirks/chatter.md>) | Rapid repeated switching or movement near a threshold or unstable equilibrium. | 1 | 33.33% | 1 |
| [Colour variation and colour shift](<quirks/colour_variation_and_colour_shift.md>) | Changes or inconsistencies in displayed colour across level, age, temperature, angle or individual units. | 1 | 33.33% | 1 |
| [Construction tolerances and unit variation](<quirks/construction_tolerances_and_unit_variation.md>) | Differences between nominally identical instruments caused by manufacturing and assembly tolerances. | 1 | 33.33% | 1 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 1 | 33.33% | 1 |
| [Display overlap and adjacent indication](<quirks/display_overlap_and_adjacent_indication.md>) | Visual interference where neighbouring pointers, digits, traces or display regions overlap. | 1 | 33.33% | 1 |
| [Flutter, jitter, tremor and quiver](<quirks/flutter_jitter_tremor_and_quiver.md>) | Small rapid random or periodic movements around the nominal indication. | 1 | 33.33% | 1 |
| [Gas-discharge instability and dropout](<quirks/gas_discharge_instability_and_dropout.md>) | Flicker, extinction or irregular conduction after a gas discharge has started. | 1 | 33.33% | 1 |
| [Gas-discharge striking and ignition](<quirks/gas_discharge_striking_and_ignition.md>) | The voltage, timing and transient behaviour required to initiate a gas discharge. | 1 | 33.33% | 1 |
| [Hysteresis](<quirks/hysteresis.md>) | Different indicated values for the same input depending on whether the input approached from above or below. | 1 | 33.33% | 1 |
| [Logarithmic scale](<quirks/logarithmic_scale.md>) | A scale where equal distances represent equal ratios rather than equal numeric increments. | 1 | 33.33% | 1 |
| [Overload, saturation and damage](<quirks/overload_saturation_and_damage.md>) | Behaviour when the input exceeds the useful range, including pegging, clipping, recovery changes or permanent harm. | 1 | 33.33% | 1 |
| [Phosphor behaviour and ageing](<quirks/phosphor_behaviour_and_ageing.md>) | Brightness, colour, persistence and degradation characteristics of phosphor-based displays. | 1 | 33.33% | 2 |
| [Segment, pixel, lamp or flag failure](<quirks/segment_pixel_lamp_or_flag_failure.md>) | Individual display elements that fail open, fail active, stick, weaken or respond intermittently. | 1 | 33.33% | 1 |
| [Shared-current or shared-supply dimming](<quirks/shared_current_or_shared_supply_dimming.md>) | Brightness reduction or interaction when multiple elements draw from a common limited current or supply. | 1 | 33.33% | 1 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 1 | 33.33% | 1 |
| [Warm-up behaviour](<quirks/warm_up_behaviour.md>) | Transient changes after startup while temperature, discharge, illumination or mechanics stabilise. | 1 | 33.33% | 1 |
| [Witness, peak-hold and retained extrema](<quirks/witness_peak_hold_and_retained_extrema.md>) | Markers or memory mechanisms that preserve minimum, maximum or peak values after the live indication moves away. | 1 | 33.33% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
