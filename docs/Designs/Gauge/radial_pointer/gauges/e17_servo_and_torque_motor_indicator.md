---
gauge_id: E17
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# E17 — Servo and torque-motor indicator

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `E17` |
| Section | `E` |
| Name | Servo and torque-motor indicator |
| Representative names or models | Autopilot repeaters; radar bearing indicators; industrial servo meters |
| Measured or indicated | Remote position or electrically computed variable |
| Era | 1930s-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 10 |
| Further-reading links | 1 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Hunting and following error](<../quirks/hunting_and_following_error.md>) | — | “Closed-loop hunting”<br>“following error” |
| [Rate limiting and motion limits](<../quirks/rate_limiting_and_motion_limits.md>) | — | “velocity and acceleration limits” |
| [Backlash and lash](<../quirks/backlash_and_lash.md>) | — | “gear backlash” |
| [Electrical hum, buzz and whine](<../quirks/electrical_hum_buzz_and_whine.md>) | — | “motor buzz” |
| [End stops, pegging and overflow](<../quirks/end_stops_pegging_and_overflow.md>) | — | “end-stop limit switches” |
| [Thresholds, deadband and switching points](<../quirks/thresholds_deadband_and_switching_points.md>) | — | “end-stop limit switches” |
| [Drift and long-term stability](<../quirks/drift_and_long_term_stability.md>) | — | “power loss may drop a warning flag or let the pointer drift” |
| [Invalid, out-of-range and warning flags](<../quirks/invalid_out_of_range_and_warning_flags.md>) | — | “power loss may drop a warning flag or let the pointer drift” |
| [Power-off and power-loss behaviour](<../quirks/power_off_and_power_loss_behaviour.md>) | — | “power loss may drop a warning flag or let the pointer drift” |
| [Power-up and self-test behaviour](<../quirks/power_up_and_self_test_behaviour.md>) | — | “test mode may slew through full scale” |

## Image references

- [Wikimedia Commons images: servo motor analog indicator gauge](<https://commons.wikimedia.org/w/index.php?search=servo+motor+analog+indicator+gauge&title=Special:MediaSearch&type=image>)

## Further reading

- [Servomechanism](<https://en.wikipedia.org/wiki/Servomechanism>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=E17]`

[Back to Radial pointer](../README.md)
