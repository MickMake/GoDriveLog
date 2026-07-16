---
gauge_id: X29
gauge_group: indicator_lamp
catalogue_version: "0.2"
date_confidence: M
---

# X29 — Automotive tell-tale / “idiot light” cluster

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Indicator lamp or illuminated legend](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `X29` |
| Section | `X` |
| Name | Automotive tell-tale / “idiot light” cluster |
| Representative names or models | Oil-pressure, charge, coolant and brake warning lamps; bulb-check circuits |
| Measured or indicated | Binary fault or threshold state |
| Era | 1930s-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 10 |
| Further-reading links | 1 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Warm-up behaviour](<../quirks/warm_up_behaviour.md>) | — | “Lamp warm-up and fade” |
| [Power-up and self-test behaviour](<../quirks/power_up_and_self_test_behaviour.md>) | — | “bulb-check sweep at key-on” |
| [Flicker, scan and PWM artefacts](<../quirks/flicker_scan_and_pwm_artefacts.md>) | — | “pressure switches flicker near idle” |
| [Brightness and contrast variation](<../quirks/brightness_and_contrast_variation.md>) | — | “alternator lamp glows dimly under small voltage difference” |
| [Supply-voltage sensitivity](<../quirks/supply_voltage_sensitivity.md>) | — | “alternator lamp glows dimly under small voltage difference” |
| [Electrical contacts, grounds and wiring](<../quirks/electrical_contacts_grounds_and_wiring.md>) | — | “shared grounds make unrelated lamps ghost” |
| [Ghosting, crosstalk and light leakage](<../quirks/ghosting_crosstalk_and_light_leakage.md>) | — | “shared grounds make unrelated lamps ghost” |
| [Ageing and material degradation](<../quirks/ageing_and_material_degradation.md>) | — | “coloured filters fade” |
| [Colour variation and colour shift](<../quirks/colour_variation_and_colour_shift.md>) | — | “coloured filters fade” |
| [Segment, pixel, lamp or flag failure](<../quirks/segment_pixel_lamp_or_flag_failure.md>) | — | “failed bulb silently removes the warning” |

## Image references

- [Wikimedia Commons images: vintage automotive warning light telltale cluster](<https://commons.wikimedia.org/w/index.php?search=vintage+automotive+warning+light+telltale+cluster&title=Special:MediaSearch&type=image>)

## Further reading

- [Telltale (automotive)](<https://en.wikipedia.org/wiki/Telltale_(automotive)>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=X29]`

[Back to Indicator lamp or illuminated legend](../README.md)
