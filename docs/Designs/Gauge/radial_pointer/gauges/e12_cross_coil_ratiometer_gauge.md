---
gauge_id: E12
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# E12 — Cross-coil / ratiometer gauge

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `E12` |
| Section | `E` |
| Name | Cross-coil / ratiometer gauge |
| Representative names or models | Two-coil fuel and temperature gauges; Wheatstone-bridge ratiometer indicators |
| Measured or indicated | Resistance ratio, fuel level, temperature or pressure sender output |
| Era | 1910s-present; dominant automotive form c.1930s-1970s |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 7 |
| Further-reading links | 1 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Input loading and ratio dependence](<../quirks/input_loading_and_ratio_dependence.md>) | — | “Angle depends on current ratio, reducing supply-voltage sensitivity” |
| [Supply-voltage sensitivity](<../quirks/supply_voltage_sensitivity.md>) | — | “Angle depends on current ratio, reducing supply-voltage sensitivity” |
| [Settling and return behaviour](<../quirks/settling_and_return_behaviour.md>) | — | “weak or absent return spring” |
| [Electrical or mechanical remote coupling](<../quirks/electrical_or_mechanical_remote_coupling.md>) | — | “sender noise makes small tremors” |
| [Flutter, jitter, tremor and quiver](<../quirks/flutter_jitter_tremor_and_quiver.md>) | — | “sender noise makes small tremors” |
| [Scale linearity and nonlinearity](<../quirks/scale_linearity_and_nonlinearity.md>) | — | “nonlinear sender/coil geometry” |
| [Electrical contacts, grounds and wiring](<../quirks/electrical_contacts_grounds_and_wiring.md>) | — | “open-circuit faults can drive to an extreme” |

## Image references

- [Wikimedia Commons images: cross coil ratiometer fuel gauge movement](<https://commons.wikimedia.org/w/index.php?search=cross+coil+ratiometer+fuel+gauge+movement&title=Special:MediaSearch&type=image>)

## Further reading

- [Ratiometer](<https://en.wikipedia.org/wiki/Ratiometer>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=E12]`

[Back to Radial pointer](../README.md)
