---
gauge_id: E13
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# E13 — Thermal bimetal automotive gauge

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `E13` |
| Section | `E` |
| Name | Thermal bimetal automotive gauge |
| Representative names or models | King-Seeley and Smiths fuel/temperature gauges; pulsed instrument-voltage regulator systems |
| Measured or indicated | Fuel level, coolant temperature, oil pressure by sender resistance |
| Era | 1920s-c.1980s common; restoration market remains |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 10 |
| Further-reading links | 1 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Response speed, lag and delay](<../quirks/response_speed_lag_and_delay.md>) | — | “Intentionally sluggish”<br>“cold start produces a slow crawl from rest” |
| [Cool-down, heat soak and continued motion](<../quirks/cool_down_heat_soak_and_continued_motion.md>) | — | “heat soak after shutdown” |
| [Thermal behaviour and temperature effects](<../quirks/thermal_behaviour_and_temperature_effects.md>) | — | “heat soak after shutdown”<br>“voltage and ambient temperature matter”<br>“cold start produces a slow crawl from rest” |
| [Hysteresis](<../quirks/hysteresis.md>) | — | “bimetal hysteresis” |
| [Drive power and regulation](<../quirks/drive_power_and_regulation.md>) | — | “pulsed regulator can cause a faint periodic quiver” |
| [Flutter, jitter, tremor and quiver](<../quirks/flutter_jitter_tremor_and_quiver.md>) | — | “pulsed regulator can cause a faint periodic quiver” |
| [Supply-voltage sensitivity](<../quirks/supply_voltage_sensitivity.md>) | — | “pulsed regulator can cause a faint periodic quiver”<br>“voltage and ambient temperature matter” |
| [Averaging and integration](<../quirks/averaging_and_integration.md>) | — | “sender slosh is averaged” |
| [Electrical or mechanical remote coupling](<../quirks/electrical_or_mechanical_remote_coupling.md>) | — | “sender slosh is averaged” |
| [Fluid surge, slosh and foam](<../quirks/fluid_surge_slosh_and_foam.md>) | — | “sender slosh is averaged” |

## Image references

- [Wikimedia Commons images: vintage bimetal automotive fuel temperature gauge](<https://commons.wikimedia.org/w/index.php?search=vintage+bimetal+automotive+fuel+temperature+gauge&title=Special:MediaSearch&type=image>)

## Further reading

- [Fuel gauge](<https://en.wikipedia.org/wiki/Fuel_gauge>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=E13]`

[Back to Radial pointer](../README.md)
