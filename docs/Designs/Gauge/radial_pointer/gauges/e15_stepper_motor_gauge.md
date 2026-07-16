---
gauge_id: E15
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: H
---

# E15 — Stepper-motor gauge

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `E15` |
| Section | `E` |
| Name | Stepper-motor gauge |
| Representative names or models | Switec/Juken X25 and X27.168; VID29; modern instrument-cluster motors |
| Measured or indicated | Digitally commanded speed, RPM, fuel, temperature and other values |
| Era | c.1980s-present |
| Date confidence | `H` — Dated primary, official or museum support |
| Canonical quirks | 9 |
| Further-reading links | 1 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Discrete, quantised or stepwise motion](<../quirks/discrete_quantised_or_stepwise_motion.md>) | — | “Discrete angular steps” |
| [End stops, pegging and overflow](<../quirks/end_stops_pegging_and_overflow.md>) | — | “startup homing sweep into a mechanical stop” |
| [Homing, parking and startup sweep](<../quirks/homing_parking_and_startup_sweep.md>) | — | “startup homing sweep into a mechanical stop”<br>“key-off parking behaviour” |
| [Chatter](<../quirks/chatter.md>) | — | “stop chatter” |
| [Dead time, missed events and duplicate counts](<../quirks/dead_time_missed_events_and_duplicate_counts.md>) | — | “missed steps and lost zero” |
| [Zero drift and offset](<../quirks/zero_drift_and_offset.md>) | — | “missed steps and lost zero” |
| [Rate limiting and motion limits](<../quirks/rate_limiting_and_motion_limits.md>) | — | “rate limiting”<br>“some motors allow only about 315 degrees, not a full circle” |
| [Backlash and lash](<../quirks/backlash_and_lash.md>) | — | “gear lash” |
| [Wraparound and multi-turn indication](<../quirks/wraparound_and_multi_turn_indication.md>) | — | “some motors allow only about 315 degrees, not a full circle” |

## Image references

- [Wikimedia Commons images: Switec X27.168 gauge stepper motor](<https://commons.wikimedia.org/w/index.php?search=Switec+X27.168+gauge+stepper+motor&title=Special:MediaSearch&type=image>)

## Further reading

- [Adafruit: Switec/Juken X27.168 gauge stepper specifications](<https://www.adafruit.com/product/2424>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=E15]`

[Back to Radial pointer](../README.md)
