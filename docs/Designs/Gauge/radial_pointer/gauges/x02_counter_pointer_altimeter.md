---
gauge_id: X02
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# X02 — Counter-pointer altimeter

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `X02` |
| Section | `X` |
| Name | Counter-pointer altimeter |
| Representative names or models | Drum counter for thousands plus one fine pointer for hundreds |
| Measured or indicated | Pressure altitude |
| Era | c.1940s-present; common in transport and military aircraft mid-century |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 7 |
| Further-reading links | 2 |
| Image-reference links | 1 |

## Alternate group memberships

- [Rolling drum or counter](<../../rolling_drum_or_counter/README.md>)

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Revolution counting and coarse/fine readout](<../quirks/revolution_counting_and_coarse_fine_readout.md>) | — | “Pointer supplies fine altitude while digits supply coarse place value” |
| [Carry, rollover and digit transition](<../quirks/carry_rollover_and_digit_transition.md>) | — | “counter carry can occur while pointer crosses zero”<br>“half-visible digits” |
| [Readout masking and narrow-window transitions](<../quirks/readout_masking_and_narrow_window_transitions.md>) | — | “half-visible digits” |
| [Invalid, out-of-range and warning flags](<../quirks/invalid_out_of_range_and_warning_flags.md>) | — | “mechanical flag for negative/out-of-range” |
| [Manual adjustment and setup](<../quirks/manual_adjustment_and_setup.md>) | — | “barometric drum” |
| [Flutter, jitter, tremor and quiver](<../quirks/flutter_jitter_tremor_and_quiver.md>) | — | “vibration can make the counter shimmer” |
| [Shock and vibration effects](<../quirks/shock_and_vibration_effects.md>) | — | “vibration can make the counter shimmer” |

## Image references

- [Wikimedia Commons images: counter pointer aircraft altimeter](<https://commons.wikimedia.org/w/index.php?search=counter+pointer+aircraft+altimeter&title=Special:MediaSearch&type=image>)

## Further reading

- [NASA: Altimeter display types and reading performance](<https://ntrs.nasa.gov/citations/19660017807>)
- [FAA: Pilot’s Handbook of Aeronautical Knowledge, Flight Instruments](<https://www.faa.gov/regulations_policies/handbooks_manuals/aviation/phak>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=X02]`

[Back to Radial pointer](../README.md)
