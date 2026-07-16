# Gauge Catalogue Documentation v0.2

> Generated from the v0.2 JSON catalogue and indexes. The JSON files remain authoritative; these Markdown files provide a browsable GitHub view.

## Catalogue summary

- **Gauges:** 136
- **Primary gauge groups:** 24
- **Canonical quirks:** 149
- **Primary group-to-quirk pages:** 562
- **Source quirk statements retained:** 885

## Navigation

- [Gauge group index](<GAUGE_GROUPS.md>)
- [Gauge index](<GAUGES.md>)
- [Canonical quirk index](<QUIRKS.md>)
- [Data and interpretation notes](<DATA_NOTES.md>)
- [Build and validation report](<BUILD_REPORT.md>)

## Gauge groups

| Gauge group | Definition | Primary gauges | Alternate members |
|---|---|---:|---:|
| [Radial pointer](<radial_pointer/README.md>) | One or more hands sweep an arc or circle against a fixed scale. | 56 | 5 |
| [Segmented display](<segmented_display/README.md>) | Fixed luminous, gas-discharge, liquid-crystal or material segments form numbers or glyphs. | 13 | 1 |
| [Liquid column](<liquid_column/README.md>) | A liquid height, meniscus or liquid/gas boundary moves in a tube, well or gauge glass. | 9 | 2 |
| [Chart or trace recorder](<chart_or_trace_recorder/README.md>) | A pen, stylus, light beam or scratch records a changing value on a moving medium. | 6 | 2 |
| [Dot-matrix or cell array](<dot_matrix_or_cell_array/README.md>) | An addressable matrix of dots or cells forms characters, symbols or fields. | 5 | 1 |
| [Moving tape, ribbon or map](<moving_tape_ribbon_or_map/README.md>) | A strip, band, tape or printed map travels through a window or past a fixed index. | 5 | 0 |
| [Rotating scan or strobe](<rotating_scan_or_strobe/README.md>) | Timing is converted to apparent position by a rotating scan, flash, glow point or strobe. | 4 | 0 |
| [Bar, column, wedge or moving-dot display](<bar_or_wedge_display/README.md>) | Magnitude is shown as a continuous or discrete bar, wedge, column or advancing dot. | 3 | 0 |
| [Colour or material-state indicator](<colour_or_material_state/README.md>) | Colour, opacity, scattering, phase or chemical state changes to indicate value or history. | 3 | 0 |
| [Floating indicator](<floating_indicator/README.md>) | A buoyant float, bob or bulb directly indicates the value by its position. | 3 | 1 |
| [Indicator lamp or illuminated legend](<indicator_lamp/README.md>) | A lamp, lens or illuminated legend indicates a state or coarse condition. | 3 | 0 |
| [Mechanical flag, shutter or semaphore](<mechanical_flag_or_shutter/README.md>) | A flag, shutter, vane, striped drum or semaphore changes between visible states. | 3 | 0 |
| [Rolling drum or counter](<rolling_drum_or_counter/README.md>) | Numeral drums, cyclometer wheels or register wheels rotate to form a count or value. | 3 | 4 |
| [Rotating scale or scene](<rotating_scale_or_scene/README.md>) | A compass card, scale, drum, horizon or scene rotates behind a fixed reference. | 3 | 1 |
| [Flip-element array](<flip_element_array/README.md>) | Bistable discs, dots, tiles or flags physically flip to form symbols or levels. | 2 | 0 |
| [Free mass, leaf, ball or bubble](<free_mass_or_bubble/README.md>) | A freely deflecting mass, leaf, ball, pendulum or bubble is itself the indicator. | 2 | 1 |
| [Linear pointer](<linear_pointer/README.md>) | A pointer, blade, hairline or marker moves along a straight or edgewise scale. | 2 | 1 |
| [Pattern-forming or biological indicator](<pattern_or_biological_indicator/README.md>) | Crystals, living organisms or other emergent patterns provide the indication. | 2 | 0 |
| [Projected spot or shadow](<projected_spot_or_shadow/README.md>) | A projected light spot, vane shadow or optical image moves across a scale. | 2 | 0 |
| [Resonant or oscillating element](<resonant_or_oscillating_element/README.md>) | A vibrating reed, string or other oscillating element is observed directly. | 2 | 0 |
| [Vector or storage trace](<vector_or_storage_trace/README.md>) | A beam or equivalent trace draws waveforms, radar returns or stored vector imagery. | 2 | 0 |
| [Optical null or match](<optical_null_or_match/README.md>) | The reading is found by visually matching brightness, colour or disappearance at a null. | 1 | 0 |
| [Projected symbology](<projected_symbology/README.md>) | Collimated or head-up symbols are projected into the operator’s field of view. | 1 | 0 |
| [Split-flap display](<split_flap/README.md>) | Physical character flaps cycle through intermediate positions to reach the target. | 1 | 0 |

## Directory layout

```text
<gauge_group>/
  README.md                 # Type definition, gauges and supporting quirks
  gauges/
    <id>_<name>.md          # Full human-readable gauge entry
  quirks/
    <canonical_quirk>.md    # Group-specific quirk detail and evidence
_data/
  *.json                    # Exact v0.2 reference files
```

Image references are links supplied by the catalogue; they are not bundled image files. Many are Wikimedia Commons searches rather than a claim that one particular photograph is the definitive model.
