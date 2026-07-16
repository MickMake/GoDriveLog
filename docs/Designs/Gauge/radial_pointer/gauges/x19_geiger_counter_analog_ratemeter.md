---
gauge_id: X19
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# X19 — Geiger-counter analog ratemeter

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `X19` |
| Section | `X` |
| Name | Geiger-counter analog ratemeter |
| Representative names or models | Victoreen, Ludlum and civil-defence survey meters with GM tube |
| Measured or indicated | Radiation count rate or inferred dose rate |
| Era | 1928-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 12 |
| Further-reading links | 2 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Flutter, jitter, tremor and quiver](<../quirks/flutter_jitter_tremor_and_quiver.md>) | — | “Individual random pulses create authentic Poisson twitch and clicks” |
| [Mechanical noise and cadence](<../quirks/mechanical_noise_and_cadence.md>) | — | “Individual random pulses create authentic Poisson twitch and clicks” |
| [Random, statistical and batch variation](<../quirks/random_statistical_and_batch_variation.md>) | — | “Individual random pulses create authentic Poisson twitch and clicks” |
| [Averaging and integration](<../quirks/averaging_and_integration.md>) | — | “RC averaging makes fast rise/slow settling” |
| [Settling and return behaviour](<../quirks/settling_and_return_behaviour.md>) | — | “RC averaging makes fast rise/slow settling” |
| [Range switching and multiple ranges](<../quirks/range_switching_and_multiple_ranges.md>) | — | “range switch changes scale” |
| [Dead time, missed events and duplicate counts](<../quirks/dead_time_missed_events_and_duplicate_counts.md>) | — | “detector dead time causes high-rate nonlinearity or saturation” |
| [Overload, saturation and damage](<../quirks/overload_saturation_and_damage.md>) | — | “detector dead time causes high-rate nonlinearity or saturation” |
| [Scale linearity and nonlinearity](<../quirks/scale_linearity_and_nonlinearity.md>) | — | “detector dead time causes high-rate nonlinearity or saturation” |
| [Scale markings, zones and legends](<../quirks/scale_markings_zones_and_legends.md>) | — | “zero adjust and battery-check zones” |
| [Zero adjustment and checking](<../quirks/zero_adjustment_and_checking.md>) | — | “zero adjust and battery-check zones” |
| [Logarithmic scale](<../quirks/logarithmic_scale.md>) | — | “logarithmic faces common” |

## Image references

- [Wikimedia Commons images: vintage Geiger counter analog meter](<https://commons.wikimedia.org/w/index.php?search=vintage+Geiger+counter+analog+meter&title=Special:MediaSearch&type=image>)

## Further reading

- [US NRC: Geiger-Mueller counter test procedure](<https://www.nrc.gov/docs/ML0037/ML003739460.pdf>)
- [Geiger counter](<https://en.wikipedia.org/wiki/Geiger_counter>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=X19]`

[Back to Radial pointer](../README.md)
