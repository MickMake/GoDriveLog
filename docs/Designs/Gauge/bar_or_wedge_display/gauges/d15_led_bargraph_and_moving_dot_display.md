---
gauge_id: D15
gauge_group: bar_or_wedge_display
catalogue_version: "0.2"
date_confidence: M
---

# D15 — LED bargraph and moving-dot display

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Bar, column, wedge or moving-dot display](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `D15` |
| Section | `D` |
| Name | LED bargraph and moving-dot display |
| Representative names or models | LM3914/LM3915 meters; automotive bar gauges; audio peak bars |
| Measured or indicated | Level, percentage, frequency bands or threshold state |
| Era | 1970s-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 10 |
| Further-reading links | 2 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Display geometry and motion mode](<../quirks/display_geometry_and_motion_mode.md>) | — | “Dot versus bar mode” |
| [Chatter](<../quirks/chatter.md>) | — | “discrete thresholds chatter with noisy input” |
| [Discrete, quantised or stepwise motion](<../quirks/discrete_quantised_or_stepwise_motion.md>) | — | “discrete thresholds chatter with noisy input” |
| [Logarithmic scale](<../quirks/logarithmic_scale.md>) | — | “logarithmic and linear driver variants” |
| [Scale linearity and nonlinearity](<../quirks/scale_linearity_and_nonlinearity.md>) | — | “logarithmic and linear driver variants” |
| [Witness, peak-hold and retained extrema](<../quirks/witness_peak_hold_and_retained_extrema.md>) | — | “a peak dot may decay separately” |
| [Shared-current or shared-supply dimming](<../quirks/shared_current_or_shared_supply_dimming.md>) | — | “current budget makes full bars dimmer in poor designs” |
| [Colour variation and colour shift](<../quirks/colour_variation_and_colour_shift.md>) | — | “thermal colour shift” |
| [Thermal behaviour and temperature effects](<../quirks/thermal_behaviour_and_temperature_effects.md>) | — | “thermal colour shift” |
| [Segment, pixel, lamp or flag failure](<../quirks/segment_pixel_lamp_or_flag_failure.md>) | — | “one dark segment is conspicuous” |

## Image references

- [Wikimedia Commons images: LED bargraph LM3914 display](<https://commons.wikimedia.org/w/index.php?search=LED+bargraph+LM3914+display&title=Special:MediaSearch&type=image>)

## Further reading

- [Texas Instruments LM3914 datasheet](<https://www.ti.com/lit/ds/symlink/lm3914.pdf>)
- [LM3914](<https://en.wikipedia.org/wiki/LM3914>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=D15]`

[Back to Bar, column, wedge or moving-dot display](../README.md)
