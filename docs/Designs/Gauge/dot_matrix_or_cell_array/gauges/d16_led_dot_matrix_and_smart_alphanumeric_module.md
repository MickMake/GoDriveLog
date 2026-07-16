---
gauge_id: D16
gauge_group: dot_matrix_or_cell_array
catalogue_version: "0.2"
date_confidence: M
---

# D16 — LED dot-matrix and smart alphanumeric module

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Dot-matrix or cell array](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `D16` |
| Section | `D` |
| Name | LED dot-matrix and smart alphanumeric module |
| Representative names or models | 5x7 matrices; HP intelligent displays; scrolling message modules |
| Measured or indicated | Text, symbols, plots and coarse gauges |
| Era | 1970s-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 7 |
| Further-reading links | 2 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Flicker, scan and PWM artefacts](<../quirks/flicker_scan_and_pwm_artefacts.md>) | — | “Row/column scan lines”<br>“rolling-shutter banding on cameras” |
| [Glyph, segment and aperture geometry](<../quirks/glyph_segment_and_aperture_geometry.md>) | — | “fixed internal fonts” |
| [Segment, pixel, lamp or flag failure](<../quirks/segment_pixel_lamp_or_flag_failure.md>) | — | “dead pixels and column drivers” |
| [Brightness and contrast variation](<../quirks/brightness_and_contrast_variation.md>) | — | “brightness varies with number of active dots unless current is compensated” |
| [Shared-current or shared-supply dimming](<../quirks/shared_current_or_shared_supply_dimming.md>) | — | “brightness varies with number of active dots unless current is compensated” |
| [Power-up and self-test behaviour](<../quirks/power_up_and_self_test_behaviour.md>) | — | “power-up test patterns, random garbage or a brief blank are period-authentic” |
| [Random, statistical and batch variation](<../quirks/random_statistical_and_batch_variation.md>) | — | “power-up test patterns, random garbage or a brief blank are period-authentic” |

## Image references

- [Wikimedia Commons images: vintage LED dot matrix smart display](<https://commons.wikimedia.org/w/index.php?search=vintage+LED+dot+matrix+smart+display&title=Special:MediaSearch&type=image>)

## Further reading

- [Texas Instruments: LED-display brightness uniformity and ghosting](<https://www.ti.com/lit/pdf/sbva057>)
- [Dot-matrix display](<https://en.wikipedia.org/wiki/Dot-matrix_display>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=D16]`

[Back to Dot-matrix or cell array](../README.md)
