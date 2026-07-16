---
gauge_group: optical_null_or_match
catalogue_version: "0.2"
primary_gauge_count: 1
supporting_quirk_count: 7
---

# Optical null or match

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

The operator adjusts or compares the instrument until a visual difference disappears or two optical appearances match.

**Catalogue definition:** The reading is found by visually matching brightness, colour or disappearance at a null.

## How the group encodes a value

The result is found at a null or match rather than read directly from a continuously displaced indicator.

## Classification boundary

Use this group when human visual matching is part of the measurement method, not merely a way of reading a scale.

## Simulation baseline

Model observer sensitivity, emissivity or colour mismatch, optical alignment, adjustment cadence and uncertainty around the null.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 1 |
| Share of catalogue | 0.74% |
| Alternate members | 0 |
| Canonical quirks represented | 7 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `P44` | [Disappearing-filament optical pyrometer](<gauges/p44_disappearing_filament_optical_pyrometer.md>) | Holborn-Kurlbaum style; Leeds & Northrup optical pyrometers | High temperature by matching filament brightness | c.1901-1980s common; still niche |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `P44` | [Disappearing-filament optical pyrometer](<gauges/p44_disappearing_filament_optical_pyrometer.md>) | High temperature by matching filament brightness | c.1901-1980s common; still niche | 7 | 1 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Ambient-light readability](<quirks/ambient_light_readability.md>) | Dependence of legibility on sunlight, darkness, glare or surrounding illumination. | 1 | 100.00% | 1 |
| [Calibration, correction and compensation](<quirks/calibration_correction_and_compensation.md>) | Adjustments or correction factors required to relate the raw indication to the intended quantity. | 1 | 100.00% | 1 |
| [Colour variation and colour shift](<quirks/colour_variation_and_colour_shift.md>) | Changes or inconsistencies in displayed colour across level, age, temperature, angle or individual units. | 1 | 100.00% | 2 |
| [Filament behaviour and failure](<quirks/filament_behaviour_and_failure.md>) | Warm-up, sag, resistance change, brightness variation and breakage of incandescent filaments. | 1 | 100.00% | 2 |
| [Operator procedure and ritual](<quirks/operator_procedure_and_ritual.md>) | Required handling or reading practices that materially affect the result. | 1 | 100.00% | 1 |
| [Optical matching, emissivity and observer effects](<quirks/optical_matching_emissivity_and_observer_effects.md>) | Measurement variation caused by visual matching, target emissivity and differences in human observation. | 1 | 100.00% | 4 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 1 | 100.00% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
