---
gauge_id: X08
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# X08 — Radar altimeter dial

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `X08` |
| Section | `X` |
| Name | Radar altimeter dial |
| Representative names or models | Low-range radio/radar altimeter with decision-height bug and warning flag |
| Measured or indicated | Height above terrain |
| Era | 1940s-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 8 |
| Further-reading links | 1 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Valid operating region](<../quirks/valid_operating_region.md>) | — | “Needle is meaningful only at low altitude” |
| [Flutter, jitter, tremor and quiver](<../quirks/flutter_jitter_tremor_and_quiver.md>) | — | “terrain and bank cause lively movement” |
| [Response speed, lag and delay](<../quirks/response_speed_lag_and_delay.md>) | — | “terrain and bank cause lively movement” |
| [Invalid, out-of-range and warning flags](<../quirks/invalid_out_of_range_and_warning_flags.md>) | — | “OFF/invalid flag” |
| [Reference bugs and set markers](<../quirks/reference_bugs_and_set_markers.md>) | — | “adjustable decision-height bug triggers lamp/tone” |
| [Compressed or expanded scale](<../quirks/compressed_or_expanded_scale.md>) | — | “scale often expands near the ground and compresses high values” |
| [Homing, parking and startup sweep](<../quirks/homing_parking_and_startup_sweep.md>) | — | “power-up test slews pointer” |
| [Power-up and self-test behaviour](<../quirks/power_up_and_self_test_behaviour.md>) | — | “power-up test slews pointer” |

## Image references

- [Wikimedia Commons images: vintage radar altimeter dial decision height bug](<https://commons.wikimedia.org/w/index.php?search=vintage+radar+altimeter+dial+decision+height+bug&title=Special:MediaSearch&type=image>)

## Further reading

- [Radar altimeter](<https://en.wikipedia.org/wiki/Radar_altimeter>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=X08]`

[Back to Radial pointer](../README.md)
