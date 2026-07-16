---
gauge_group: projected_symbology
catalogue_version: "0.2"
primary_gauge_count: 1
supporting_quirk_count: 8
---

# Projected symbology

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

Collimated or head-up symbols are optically projected into the operator’s field of view and registered with an external scene or sight line.

**Catalogue definition:** Collimated or head-up symbols are projected into the operator’s field of view.

## How the group encodes a value

Information is encoded by projected symbols whose apparent distance and alignment are controlled optically.

## Classification boundary

Use this group for scene-registered projected overlays. A local projected spot on a scale belongs under projected_spot_or_shadow.

## Simulation baseline

Include collimation, parallax, eyebox limits, clipping, focus, combiner reflections, brightness adaptation and alignment error.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 1 |
| Share of catalogue | 0.74% |
| Alternate members | 0 |
| Canonical quirks represented | 8 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `D28` | [Head-up and collimated projected symbology](<gauges/d28_head_up_and_collimated_projected_symbology.md>) | Reflector gunsights; aircraft HUDs; automotive combiner HUDs | Speed, attitude, aiming, navigation and warnings at optical infinity | 1940s-present |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `D28` | [Head-up and collimated projected symbology](<gauges/d28_head_up_and_collimated_projected_symbology.md>) | Speed, attitude, aiming, navigation and warnings at optical infinity | 1940s-present | 8 | 1 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Ambient-light readability](<quirks/ambient_light_readability.md>) | Dependence of legibility on sunlight, darkness, glare or surrounding illumination. | 1 | 100.00% | 1 |
| [Bloom, halo and penumbra](<quirks/bloom_halo_and_penumbra.md>) | Spreading or soft-edged light around a spot, trace, segment or projected image. | 1 | 100.00% | 1 |
| [Contamination, dirt and fouling](<quirks/contamination_dirt_and_fouling.md>) | Reading or appearance changes caused by deposits, dust, oxidation, residue or biological growth. | 1 | 100.00% | 1 |
| [Ghosting, crosstalk and light leakage](<quirks/ghosting_crosstalk_and_light_leakage.md>) | Unwanted partial activation or illumination of neighbouring, previous or nominally inactive display elements. | 1 | 100.00% | 1 |
| [Optical distortion and refraction](<quirks/optical_distortion_and_refraction.md>) | Apparent displacement or shape changes caused by lenses, glass, liquid, curved windows or refractive interfaces. | 1 | 100.00% | 1 |
| [Parallax](<quirks/parallax.md>) | Apparent reading error caused by viewing the indicator and scale from the wrong angle or depth relationship. | 1 | 100.00% | 1 |
| [Phosphor behaviour and ageing](<quirks/phosphor_behaviour_and_ageing.md>) | Brightness, colour, persistence and degradation characteristics of phosphor-based displays. | 1 | 100.00% | 1 |
| [Projected-display eyebox and clipping](<quirks/projected_display_eyebox_and_clipping.md>) | Loss or truncation of projected symbols when the observer moves outside the usable viewing volume. | 1 | 100.00% | 2 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
