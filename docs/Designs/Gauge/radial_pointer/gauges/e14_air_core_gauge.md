---
gauge_id: E14
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# E14 — Air-core gauge

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `E14` |
| Section | `E` |
| Name | Air-core gauge |
| Representative names or models | Two perpendicular coils and permanent-magnet rotor; General Motors clusters and industrial indicators |
| Measured or indicated | Vehicle speed, fuel, temperature or arbitrary angular command |
| Era | 1960s-present; common c.1970s-1990s |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 6 |
| Further-reading links | 1 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Power-off and power-loss behaviour](<../quirks/power_off_and_power_loss_behaviour.md>) | — | “No return spring, so power-off position may be arbitrary”<br>“supply loss may leave the pointer stranded” |
| [Settling and return behaviour](<../quirks/settling_and_return_behaviour.md>) | — | “No return spring, so power-off position may be arbitrary” |
| [Response speed, lag and delay](<../quirks/response_speed_lag_and_delay.md>) | — | “fast response and 360-degree capability” |
| [Wraparound and multi-turn indication](<../quirks/wraparound_and_multi_turn_indication.md>) | — | “fast response and 360-degree capability”<br>“rotor can flip to an equivalent angle” |
| [Human-factor ambiguity and misreading](<../quirks/human_factor_ambiguity_and_misreading.md>) | — | “rotor can flip to an equivalent angle” |
| [Scale linearity and nonlinearity](<../quirks/scale_linearity_and_nonlinearity.md>) | — | “coil imbalance produces nonlinear zones” |

## Image references

- [Wikimedia Commons images: air core gauge movement instrument cluster](<https://commons.wikimedia.org/w/index.php?search=air+core+gauge+movement+instrument+cluster&title=Special:MediaSearch&type=image>)

## Further reading

- [Air core gauge](<https://en.wikipedia.org/wiki/Air_core_gauge>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=E14]`

[Back to Radial pointer](../README.md)
