---
gauge_id: X03
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# X03 — Drum-pointer and counter-drum altimeter

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `X03` |
| Section | `X` |
| Name | Drum-pointer and counter-drum altimeter |
| Representative names or models | Kollsman-style drum-pointer; counter-drum-pointer displays studied by NASA |
| Measured or indicated | Pressure altitude |
| Era | c.1950s-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 9 |
| Further-reading links | 2 |
| Image-reference links | 1 |

## Alternate group memberships

- [Rolling drum or counter](<../../rolling_drum_or_counter/README.md>)

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Human-factor ambiguity and misreading](<../quirks/human_factor_ambiguity_and_misreading.md>) | — | “Rolling drums reduce hand-counting but introduce transition ambiguity” |
| [Invalid, out-of-range and warning flags](<../quirks/invalid_out_of_range_and_warning_flags.md>) | — | “drum can be masked until valid”<br>“baro setting and low-altitude warning sectors” |
| [Readout masking and narrow-window transitions](<../quirks/readout_masking_and_narrow_window_transitions.md>) | — | “drum can be masked until valid” |
| [Revolution counting and coarse/fine readout](<../quirks/revolution_counting_and_coarse_fine_readout.md>) | — | “pointer often makes one revolution per 1,000 feet” |
| [Wraparound and multi-turn indication](<../quirks/wraparound_and_multi_turn_indication.md>) | — | “pointer often makes one revolution per 1,000 feet” |
| [Carry, rollover and digit transition](<../quirks/carry_rollover_and_digit_transition.md>) | — | “carry is intentionally delayed or accelerated” |
| [Response speed, lag and delay](<../quirks/response_speed_lag_and_delay.md>) | — | “carry is intentionally delayed or accelerated” |
| [Manual adjustment and setup](<../quirks/manual_adjustment_and_setup.md>) | — | “baro setting and low-altitude warning sectors” |
| [Scale markings, zones and legends](<../quirks/scale_markings_zones_and_legends.md>) | — | “baro setting and low-altitude warning sectors” |

## Image references

- [Wikimedia Commons images: counter drum pointer altimeter](<https://commons.wikimedia.org/w/index.php?search=counter+drum+pointer+altimeter&title=Special:MediaSearch&type=image>)

## Further reading

- [NASA: Counter-drum-pointer altimeter study](<https://ntrs.nasa.gov/citations/19730011592>)
- [NASA: Altimeter display types and reading performance](<https://ntrs.nasa.gov/citations/19660017807>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=X03]`

[Back to Radial pointer](../README.md)
