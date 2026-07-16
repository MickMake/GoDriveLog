---
gauge_id: X18
gauge_group: radial_pointer
catalogue_version: "0.2"
date_confidence: M
---

# X18 — Radio S-meter

> Generated from the v0.2 catalogue. The JSON entry is authoritative.

**Primary group:** [Radial pointer](<../README.md>)

## Catalogue entry

| Field | Value |
|---|---|
| ID | `X18` |
| Section | `X` |
| Name | Radio S-meter |
| Representative names or models | Receiver signal-strength meter, usually S1-S9 then dB above S9 |
| Measured or indicated | Relative received radio signal strength |
| Era | 1930s-present |
| Date confidence | `M` — Approximate adoption/common-use range from reliable overviews |
| Canonical quirks | 6 |
| Further-reading links | 2 |
| Image-reference links | 1 |

## Quirks to simulate

Canonical labels make behaviours searchable; qualifiers and exact source phrases preserve the gauge-specific meaning.

| Canonical quirk | Qualifiers | Preserved source phrases |
|---|---|---|
| [Frequency, waveform and source dependence](<../quirks/frequency_waveform_and_source_dependence.md>) | — | “Driven by AGC rather than a precision detector”<br>“needle rises fast and falls with AGC decay” |
| [Compressed or expanded scale](<../quirks/compressed_or_expanded_scale.md>) | — | “low end is often compressed and receiver-specific” |
| [Construction tolerances and unit variation](<../quirks/construction_tolerances_and_unit_variation.md>) | — | “low end is often compressed and receiver-specific”<br>“S9 conventionally corresponds to 50 microvolts on HF but real sets vary” |
| [Calibration, correction and compensation](<../quirks/calibration_correction_and_compensation.md>) | — | “S9 conventionally corresponds to 50 microvolts on HF but real sets vary” |
| [Attack and release ballistics](<../quirks/attack_and_release_ballistics.md>) | — | “needle rises fast and falls with AGC decay” |
| [Drift and long-term stability](<../quirks/drift_and_long_term_stability.md>) | — | “zero wanders with gain” |

## Image references

- [Wikimedia Commons images: vintage radio S meter analog](<https://commons.wikimedia.org/w/index.php?search=vintage+radio+S+meter+analog&title=Special:MediaSearch&type=image>)

## Further reading

- [ARRL: Practical S-meter calibration and nonlinearity](<https://www.arrl.org/test-procedures-to-measure-bpl-interference>)
- [S meter](<https://en.wikipedia.org/wiki/S_meter>)

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Entry key: `entries[id=X18]`

[Back to Radial pointer](../README.md)
