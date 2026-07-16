---
gauge_id: P30
gauge_group: chart_or_trace_recorder
catalogue_version: "0.2"
date_confidence: M
---

# P30 — Rain gauge and tipping-bucket recorder

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Chart or trace recorder](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `P30` |
| Section | `P` |
| Name | Rain gauge and tipping-bucket recorder |
| Representative names or models | Standard cylindrical gauge; siphon gauge; tipping bucket with reed switch |
| Measured or indicated | Rainfall depth and rate |
| Era | 1400s-present; tipping buckets common from late 1800s/1900s |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 6 |
| Further-reading links | 2 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Discrete, quantised or stepwise motion](<../quirks/discrete_quantised_or_stepwise_motion.md>) | — | “Discrete bucket tips”<br>“chart traces form staircase ramps” |
| [Dead time, missed events and duplicate counts](<../quirks/dead_time_missed_events_and_duplicate_counts.md>) | — | “dead time while tipping”<br>“high-rate under-read”<br>“mechanical bounce may generate duplicate electrical pulses” |
| [Wind, splash and collection error](<../quirks/wind_splash_and_collection_error.md>) | — | “undercatch in wind”<br>“splash-in/out” |
| [Blockage and restricted passages](<../quirks/blockage_and_restricted_passages.md>) | — | “blockage” |
| [Evaporation and drying](<../quirks/evaporation_and_drying.md>) | — | “evaporation” |
| [Bounce](<../quirks/bounce.md>) | — | “mechanical bounce may generate duplicate electrical pulses” |

## Image references

- [Wikimedia Commons images: tipping bucket rain gauge mechanism](<https://commons.wikimedia.org/w/index.php?search=tipping+bucket+rain+gauge+mechanism&title=Special:MediaSearch&type=image>)

## Further reading

- [Met Office: Rainfall observations](<https://www.metoffice.gov.uk/weather/guides/observations/how-we-measure-rainfall>)
- [Rain gauge](<https://en.wikipedia.org/wiki/Rain_gauge>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=P30]`

[Back to Chart or trace recorder](../README.md)
