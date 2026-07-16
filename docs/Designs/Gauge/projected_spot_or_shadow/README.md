---
gauge_group: projected_spot_or_shadow
catalogue_version: "0.2"
primary_gauge_count: 2
supporting_quirk_count: 15
---

# Projected spot or shadow

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A moving light spot, reflected beam, vane shadow or projected optical image traverses a scale or target.

**Catalogue definition:** A projected light spot, vane shadow or optical image moves across a scale.

## How the group encodes a value

A small mechanical movement is optically magnified into a larger visible displacement.

## Classification boundary

Use this group for an optical pointer or shadow cast onto a local scale. Head-up overlays belong under projected_symbology.

## Simulation baseline

Include focus, spot size, halo, optical alignment, ambient light and the geometry between mirror, lamp, scale and observer.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 2 |
| Share of catalogue | 1.47% |
| Alternate members | 0 |
| Canonical quirks represented | 15 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `D03` | [Shadow meter / moving-vane tuning indicator](<gauges/d03_shadow_meter_moving_vane_tuning_indicator.md>) | Philips shadow meter; pre-war radio tuning indicators | Radio tuning strength, balance or null | c.1932-c.1950s common |
| `E09` | [Mirror / light-spot galvanometer](<gauges/e09_mirror_light_spot_galvanometer.md>) | Thomson mirror galvanometer; ballistic galvanometer with projected light spot | Very small current, charge or magnetic-flux change | 1850s-c.1950s common; laboratory demonstrations remain |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `D03` | [Shadow meter / moving-vane tuning indicator](<gauges/d03_shadow_meter_moving_vane_tuning_indicator.md>) | Radio tuning strength, balance or null | c.1932-c.1950s common | 8 | 1 | 1 |
| `E09` | [Mirror / light-spot galvanometer](<gauges/e09_mirror_light_spot_galvanometer.md>) | Very small current, charge or magnetic-flux change | 1850s-c.1950s common; laboratory demonstrations remain | 9 | 2 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Bloom, halo and penumbra](<quirks/bloom_halo_and_penumbra.md>) | Spreading or soft-edged light around a spot, trace, segment or projected image. | 2 | 100.00% | 2 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 2 | 100.00% | 2 |
| [Ageing and material degradation](<quirks/ageing_and_material_degradation.md>) | Progressive change due to wear, fatigue, chemical decay, phosphor loss, embrittlement or similar ageing processes. | 1 | 50.00% | 1 |
| [Ambient-light readability](<quirks/ambient_light_readability.md>) | Dependence of legibility on sunlight, darkness, glare or surrounding illumination. | 1 | 50.00% | 1 |
| [Brightness and contrast variation](<quirks/brightness_and_contrast_variation.md>) | Changes in luminance or visual separation between active and inactive parts of the display. | 1 | 50.00% | 1 |
| [Damping](<quirks/damping.md>) | Deliberate or inherent suppression of rapid movement and oscillation in the indication. | 1 | 50.00% | 1 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 1 | 50.00% | 1 |
| [Drift and long-term stability](<quirks/drift_and_long_term_stability.md>) | Slow change in indicated value or behaviour over time despite an unchanged input. | 1 | 50.00% | 1 |
| [Focus and optical alignment](<quirks/focus_and_optical_alignment.md>) | Sharpness and registration effects caused by optical focus and component alignment. | 1 | 50.00% | 2 |
| [Ringing and oscillation](<quirks/ringing_and_oscillation.md>) | Repeated decaying or sustained motion around a target following excitation or disturbance. | 1 | 50.00% | 1 |
| [Scale linearity and nonlinearity](<quirks/scale_linearity_and_nonlinearity.md>) | Variation in how equal input increments map to equal or unequal distances on the displayed scale. | 1 | 50.00% | 1 |
| [Shadows, depth and occlusion](<quirks/shadows_depth_and_occlusion.md>) | Visual effects caused by layered parts blocking, shading or appearing at different depths. | 1 | 50.00% | 2 |
| [Shock and vibration effects](<quirks/shock_and_vibration_effects.md>) | Temporary or permanent indication changes caused by mechanical shock or sustained vibration. | 1 | 50.00% | 1 |
| [Warm-up behaviour](<quirks/warm_up_behaviour.md>) | Transient changes after startup while temperature, discharge, illumination or mechanics stabilise. | 1 | 50.00% | 1 |
| [Zero drift and offset](<quirks/zero_drift_and_offset.md>) | A non-zero indication at the true zero point, including offsets that change with time or conditions. | 1 | 50.00% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
