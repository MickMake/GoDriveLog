---
gauge_id: D13
gauge_group: segmented_display
catalogue_version: "0.2"
date_confidence: M
---

# D13 — LED seven-segment display

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Segmented display](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `D13` |
| Section | `D` |
| Name | LED seven-segment display |
| Representative names or models | Monsanto/HP early modules; HP 5082 series; Kingbright and Lite-On modules |
| Measured or indicated | Numeric readout |
| Era | late 1960s-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 8 |
| Further-reading links | 3 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Flicker, scan and PWM artefacts](<../quirks/flicker_scan_and_pwm_artefacts.md>) | — | “Multiplex flicker and scan order” |
| [Ghosting, crosstalk and light leakage](<../quirks/ghosting_crosstalk_and_light_leakage.md>) | — | “ghost segments during switching” |
| [Colour variation and colour shift](<../quirks/colour_variation_and_colour_shift.md>) | — | “lens and tinted-filter effects” |
| [Segment, pixel, lamp or flag failure](<../quirks/segment_pixel_lamp_or_flag_failure.md>) | — | “dead or weak segments” |
| [Brightness and contrast variation](<../quirks/brightness_and_contrast_variation.md>) | — | “thermal brightness drift” |
| [Drift and long-term stability](<../quirks/drift_and_long_term_stability.md>) | — | “thermal brightness drift” |
| [Thermal behaviour and temperature effects](<../quirks/thermal_behaviour_and_temperature_effects.md>) | — | “thermal brightness drift” |
| [Shared-current or shared-supply dimming](<../quirks/shared_current_or_shared_supply_dimming.md>) | — | “count-dependent dimming occurs with a weak supply, shared resistor or shared current budget—not with a correctly regulated constant-current design” |

## Image references

- [Wikimedia Commons images: vintage red LED seven segment display](<https://commons.wikimedia.org/w/index.php?search=vintage+red+LED+seven+segment+display&title=Special:MediaSearch&type=image>)

## Further reading

- [Texas Instruments: LED-display brightness uniformity and ghosting](<https://www.ti.com/lit/pdf/sbva057>)
- [Analog Devices: MAX7219 multiplexed LED driver](<https://www.analog.com/en/products/max7219.html>)
- [Seven-segment display](<https://en.wikipedia.org/wiki/Seven-segment_display>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=D13]`

[Back to Segmented display](../README.md)
