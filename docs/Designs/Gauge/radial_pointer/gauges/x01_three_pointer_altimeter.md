---
gauge_id: X01
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# X01 — Three-pointer altimeter

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `X01` |
| Section | `X` |
| Name | Three-pointer altimeter |
| Representative names or models | Classic sensitive altimeter with 100-, 1,000- and 10,000-foot hands |
| Measured or indicated | Pressure altitude |
| Era | c.1920s-c.1970s common; some remain in service |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 6 |
| Further-reading links | 2 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Human-factor ambiguity and misreading](<../quirks/human_factor_ambiguity_and_misreading.md>) | — | “Clock-like place-value ambiguity”<br>“a small 10,000-foot pointer is easy to miss”<br>“NASA documented serious misreading risk, especially near pointer overlaps” |
| [Wraparound and multi-turn indication](<../quirks/wraparound_and_multi_turn_indication.md>) | — | “all hands wrap” |
| [Manual adjustment and setup](<../quirks/manual_adjustment_and_setup.md>) | — | “barometric setting window and knob” |
| [Pressure pulsation and pneumatic-line dynamics](<../quirks/pressure_pulsation_and_pneumatic_line_dynamics.md>) | — | “static-pressure lag” |
| [Response speed, lag and delay](<../quirks/response_speed_lag_and_delay.md>) | — | “static-pressure lag” |
| [Shock and vibration effects](<../quirks/shock_and_vibration_effects.md>) | — | “mechanical vibration” |

## Image references

- [Wikimedia Commons images: three pointer aircraft altimeter](<https://commons.wikimedia.org/w/index.php?search=three+pointer+aircraft+altimeter&title=Special:MediaSearch&type=image>)

## Further reading

- [NASA: Altimeter display types and reading performance](<https://ntrs.nasa.gov/citations/19660017807>)
- [FAA: Pilot’s Handbook of Aeronautical Knowledge, Flight Instruments](<https://www.faa.gov/regulations_policies/handbooks_manuals/aviation/phak>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=X01]`

[Back to Radial pointer](../README.md)
