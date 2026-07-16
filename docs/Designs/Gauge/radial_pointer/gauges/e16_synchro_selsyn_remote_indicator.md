---
gauge_id: E16
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# E16 — Synchro / selsyn remote indicator

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `E16` |
| Section | `E` |
| Name | Synchro / selsyn remote indicator |
| Representative names or models | Bendix and General Electric synchros; ship, aircraft and radar repeaters |
| Measured or indicated | Remote angular position such as heading, valve position or antenna bearing |
| Era | 1910s-present; peak use c.1930s-1970s |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 7 |
| Further-reading links | 1 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Electrical or mechanical remote coupling](<../quirks/electrical_or_mechanical_remote_coupling.md>) | — | “Pointer follows a transmitter electrically” |
| [Hunting and following error](<../quirks/hunting_and_following_error.md>) | — | “null hunting” |
| [Electrical contacts, grounds and wiring](<../quirks/electrical_contacts_grounds_and_wiring.md>) | — | “phase and wiring errors”<br>“loss of one phase creates wrong or weak positions” |
| [Movement torque and sensitivity](<../quirks/movement_torque_and_sensitivity.md>) | — | “smooth but slightly springy torque” |
| [Human-factor ambiguity and misreading](<../quirks/human_factor_ambiguity_and_misreading.md>) | — | “loss of one phase creates wrong or weak positions” |
| [Electrical hum, buzz and whine](<../quirks/electrical_hum_buzz_and_whine.md>) | — | “400 Hz aircraft systems may hum” |
| [Power-off and power-loss behaviour](<../quirks/power_off_and_power_loss_behaviour.md>) | — | “power-off pointer is free or parked” |

## Image references

- [Wikimedia Commons images: synchro selsyn remote position indicator](<https://commons.wikimedia.org/w/index.php?search=synchro+selsyn+remote+position+indicator&title=Special:MediaSearch&type=image>)

## Further reading

- [Synchro](<https://en.wikipedia.org/wiki/Synchro>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=E16]`

[Back to Radial pointer](../README.md)
