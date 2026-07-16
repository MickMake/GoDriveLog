---
gauge_group: vector_or_storage_trace
catalogue_version: "0.2"
primary_gauge_count: 2
supporting_quirk_count: 14
---

# Vector or storage trace

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A beam or equivalent trace draws waveforms, radar returns, vectors or stored imagery by steering directly through display space.

**Catalogue definition:** A beam or equivalent trace draws waveforms, radar returns or stored vector imagery.

## How the group encodes a value

Information is encoded as luminous paths, spots and intensity rather than a raster of independently addressed pixels.

## Classification boundary

Use this group for vector or storage traces. Physical paper recording belongs under chart_or_trace_recorder.

## Simulation baseline

Model sweep timing, phosphor persistence, bloom, focus, intensity modulation, retrace, erase behaviour and burn-in.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 2 |
| Share of catalogue | 1.47% |
| Alternate members | 0 |
| Canonical quirks represented | 14 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `D26` | [CRT oscilloscope, vector and radar display](<gauges/d26_crt_oscilloscope_vector_and_radar_display.md>) | Round radar PPI; rectangular oscilloscope; vector monitors and engine analysers | Waveforms, range/bearing, vectors or computed traces | 1930s-c.2000s mainstream; specialist/collector use remains |
| `D27` | [Storage CRT](<gauges/d27_storage_crt.md>) | Tektronix direct-view bistable storage tubes; radar and analytical instruments | Waveforms or traces that must remain without refresh | 1950s-c.1980s common |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `D26` | [CRT oscilloscope, vector and radar display](<gauges/d26_crt_oscilloscope_vector_and_radar_display.md>) | Waveforms, range/bearing, vectors or computed traces | 1930s-c.2000s mainstream; specialist/collector use remains | 9 | 2 | 1 |
| `D27` | [Storage CRT](<gauges/d27_storage_crt.md>) | Waveforms or traces that must remain without refresh | 1950s-c.1980s common | 8 | 1 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Bloom, halo and penumbra](<quirks/bloom_halo_and_penumbra.md>) | Spreading or soft-edged light around a spot, trace, segment or projected image. | 2 | 100.00% | 2 |
| [Burn-in and permanent image wear](<quirks/burn_in_and_permanent_image_wear.md>) | Lasting visible damage or nonuniform ageing from repeatedly displaying the same pattern. | 2 | 100.00% | 2 |
| [Persistence and afterglow](<quirks/persistence_and_afterglow.md>) | Continued visibility after excitation is reduced or removed. | 2 | 100.00% | 3 |
| [Brightness and contrast variation](<quirks/brightness_and_contrast_variation.md>) | Changes in luminance or visual separation between active and inactive parts of the display. | 1 | 50.00% | 1 |
| [Brightness nonuniformity and gradients](<quirks/brightness_nonuniformity_and_gradients.md>) | Unequal luminance across a display, character, segment, tube or field. | 1 | 50.00% | 1 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 1 | 50.00% | 1 |
| [Flicker, scan and PWM artefacts](<quirks/flicker_scan_and_pwm_artefacts.md>) | Visible modulation caused by multiplexing, scanning, pulse-width control or interaction with cameras and eye motion. | 1 | 50.00% | 2 |
| [Focus and optical alignment](<quirks/focus_and_optical_alignment.md>) | Sharpness and registration effects caused by optical focus and component alignment. | 1 | 50.00% | 1 |
| [Latching and state-retention behaviour](<quirks/latching_and_state_retention_behaviour.md>) | A displayed state that remains mechanically, magnetically, electrically or optically retained until reset or rewritten. | 1 | 50.00% | 1 |
| [Optical distortion and refraction](<quirks/optical_distortion_and_refraction.md>) | Apparent displacement or shape changes caused by lenses, glass, liquid, curved windows or refractive interfaces. | 1 | 50.00% | 3 |
| [Parallax](<quirks/parallax.md>) | Apparent reading error caused by viewing the indicator and scale from the wrong angle or depth relationship. | 1 | 50.00% | 1 |
| [Pen, stylus and trace artefacts](<quirks/pen_stylus_and_trace_artefacts.md>) | Line-width, drag, skipping, smear, lift-off or other defects introduced by a recording or tracing element. | 1 | 50.00% | 2 |
| [Phosphor behaviour and ageing](<quirks/phosphor_behaviour_and_ageing.md>) | Brightness, colour, persistence and degradation characteristics of phosphor-based displays. | 1 | 50.00% | 1 |
| [Refresh, erase and update artefacts](<quirks/refresh_erase_and_update_artefacts.md>) | Visible effects produced while changing, clearing or rewriting a display. | 1 | 50.00% | 2 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
