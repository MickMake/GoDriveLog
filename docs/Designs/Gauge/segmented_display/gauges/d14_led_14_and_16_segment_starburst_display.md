---
gauge_id: D14
gauge_group: segmented_display
catalogue_version: "0.2"
date_confidence: M
---

# D14 — LED 14- and 16-segment starburst display

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Segmented display](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `D14` |
| Section | `D` |
| Name | LED 14- and 16-segment starburst display |
| Representative names or models | DL-1414 smart display; HDSP alphanumeric modules; avionics and test gear |
| Measured or indicated | Letters, numerals and symbols |
| Era | 1970s-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 8 |
| Further-reading links | 2 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Brightness nonuniformity and gradients](<../quirks/brightness_nonuniformity_and_gradients.md>) | — | “Diagonal joints create bright knots”<br>“segment-current mismatch” |
| [Glyph, segment and aperture geometry](<../quirks/glyph_segment_and_aperture_geometry.md>) | — | “Diagonal joints create bright knots”<br>“many glyphs are compromises”<br>“internal-font variants have idiosyncratic lowercase”<br>“failed segment changes one glyph into another rather than merely losing a bar” |
| [Channel mismatch and unequal dynamics](<../quirks/channel_mismatch_and_unequal_dynamics.md>) | — | “segment-current mismatch” |
| [Flicker, scan and PWM artefacts](<../quirks/flicker_scan_and_pwm_artefacts.md>) | — | “multiplex shimmer” |
| [Flutter, jitter, tremor and quiver](<../quirks/flutter_jitter_tremor_and_quiver.md>) | — | “multiplex shimmer” |
| [Shared-current or shared-supply dimming](<../quirks/shared_current_or_shared_supply_dimming.md>) | — | “shared current limits can dim dense letters” |
| [Display overlap and adjacent indication](<../quirks/display_overlap_and_adjacent_indication.md>) | — | “failed segment changes one glyph into another rather than merely losing a bar” |
| [Segment, pixel, lamp or flag failure](<../quirks/segment_pixel_lamp_or_flag_failure.md>) | — | “failed segment changes one glyph into another rather than merely losing a bar” |

## Image references

- [Wikimedia Commons images: vintage LED 16 segment alphanumeric display](<https://commons.wikimedia.org/w/index.php?search=vintage+LED+16+segment+alphanumeric+display&title=Special:MediaSearch&type=image>)

## Further reading

- [Fourteen-segment display](<https://en.wikipedia.org/wiki/Fourteen-segment_display>)
- [Sixteen-segment display](<https://en.wikipedia.org/wiki/Sixteen-segment_display>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=D14]`

[Back to Segmented display](../README.md)
