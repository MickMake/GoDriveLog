---
gauge_id: X15
gauge_group: rotating_scan_or_strobe
catalogue_version: "0.2"
date_confidence: H
---

# X15 — Mechanical-wheel flasher sonar

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Rotating scan or strobe](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `X15` |
| Section | `X` |
| Name | Mechanical-wheel flasher sonar |
| Representative names or models | Vexilar FL/FLX series and earlier ice-fishing flashers |
| Measured or indicated | Current sonar return versus depth |
| Era | 1960s-present |
| Date confidence | `H` — Dated primary, official or museum support |
| Canonical quirks | 9 |
| Further-reading links | 1 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Periodic marker and rotating-scan behaviour](<../quirks/periodic_marker_and_rotating_scan_behaviour.md>) | — | “A high-speed wheel carries fixed lights that flash at timed positions” |
| [Latching and state-retention behaviour](<../quirks/latching_and_state_retention_behaviour.md>) | — | “no scrolling history” |
| [Colour variation and colour shift](<../quirks/colour_variation_and_colour_shift.md>) | — | “colours encode return strength” |
| [Scale direction and interpretation](<../quirks/scale_direction_and_interpretation.md>) | — | “colours encode return strength” |
| [Flicker, scan and PWM artefacts](<../quirks/flicker_scan_and_pwm_artefacts.md>) | — | “visible spoke/scan persistence” |
| [Persistence and afterglow](<../quirks/persistence_and_afterglow.md>) | — | “visible spoke/scan persistence” |
| [Mechanical noise and cadence](<../quirks/mechanical_noise_and_cadence.md>) | — | “motor whir” |
| [Bloom, halo and penumbra](<../quirks/bloom_halo_and_penumbra.md>) | — | “gain blooms targets” |
| [Multiple echoes, interference and remapped scale](<../quirks/multiple_echoes_interference_and_remapped_scale.md>) | — | “gain blooms targets”<br>“interference creates moving false marks”<br>“bottom lock and zoom remap the scale” |

## Image references

- [Wikimedia Commons images: Vexilar mechanical flasher sonar display](<https://commons.wikimedia.org/w/index.php?search=Vexilar+mechanical+flasher+sonar+display&title=Special:MediaSearch&type=image>)

## Further reading

- [Vexilar FLX-28 manual: how flashers work](<https://vexilar.com/pages/vexilar-flx-28-flasher-instruction-manual>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=X15]`

[Back to Rotating scan or strobe](../README.md)
