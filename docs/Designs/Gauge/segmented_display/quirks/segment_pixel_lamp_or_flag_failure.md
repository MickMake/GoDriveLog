---
gauge_group: segmented_display
quirk: "segment, pixel, lamp or flag failure"
catalogue_version: "0.2"
gauge_count_in_group: 3
gauge_count_global: 10
---

# Segment, pixel, lamp or flag failure

**Gauge group:** [Segmented display](<../README.md>)

## Definition

Individual display elements that fail open, fail active, stick, weaken or respond intermittently.

This is a canonical umbrella label. The exact per-gauge wording and any qualifiers below remain the authoritative detail; gauges grouped here may reach the same visible symptom through different physics.

## Frequency

| Scope | Gauges | Share | Source statements |
|---|---:|---:|---:|
| Segmented display | 3 of 13 | 23.08% | 3 |
| Entire catalogue | 10 of 136 | 7.35% | 10 |

## Gauges and preserved evidence

### [D13 — LED seven-segment display](<../gauges/d13_led_seven_segment_display.md>)

- **Measured or indicated:** Numeric readout
- **Era:** late 1960s-present
- **Preserved source phrases:**
  - “dead or weak segments”
- **Image references:**
  - [Wikimedia Commons images: vintage red LED seven segment display](<https://commons.wikimedia.org/w/index.php?search=vintage+red+LED+seven+segment+display&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [Texas Instruments: LED-display brightness uniformity and ghosting](<https://www.ti.com/lit/pdf/sbva057>)
  - [Analog Devices: MAX7219 multiplexed LED driver](<https://www.analog.com/en/products/max7219.html>)
  - [Seven-segment display](<https://en.wikipedia.org/wiki/Seven-segment_display>)

### [D14 — LED 14- and 16-segment starburst display](<../gauges/d14_led_14_and_16_segment_starburst_display.md>)

- **Measured or indicated:** Letters, numerals and symbols
- **Era:** 1970s-present
- **Preserved source phrases:**
  - “failed segment changes one glyph into another rather than merely losing a bar”
- **Image references:**
  - [Wikimedia Commons images: vintage LED 16 segment alphanumeric display](<https://commons.wikimedia.org/w/index.php?search=vintage+LED+16+segment+alphanumeric+display&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [Fourteen-segment display](<https://en.wikipedia.org/wiki/Fourteen-segment_display>)
  - [Sixteen-segment display](<https://en.wikipedia.org/wiki/Sixteen-segment_display>)

### [D20 — Twisted-nematic segmented LCD](<../gauges/d20_twisted_nematic_segmented_lcd.md>)

- **Measured or indicated:** Numeric, icons, bars and fixed legends
- **Era:** 1971-present
- **Preserved source phrases:**
  - “weak elastomer contacts cause missing rows”
- **Image references:**
  - [Wikimedia Commons images: segmented twisted nematic LCD meter](<https://commons.wikimedia.org/w/index.php?search=segmented+twisted+nematic+LCD+meter&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [IEEE Spectrum: RCA and the first liquid-crystal displays](<https://spectrum.ieee.org/the-father-of-the-lcd>)
  - [Twisted nematic field effect](<https://en.wikipedia.org/wiki/Twisted_nematic_field_effect>)

## Other gauge groups with this quirk

| Gauge group | Gauges with quirk | Group share |
|---|---:|---:|
| [Dot-matrix or cell array](<../../dot_matrix_or_cell_array/quirks/segment_pixel_lamp_or_flag_failure.md>) | 1 | 20.00% |
| [Bar, column, wedge or moving-dot display](<../../bar_or_wedge_display/quirks/segment_pixel_lamp_or_flag_failure.md>) | 1 | 33.33% |
| [Indicator lamp or illuminated legend](<../../indicator_lamp/quirks/segment_pixel_lamp_or_flag_failure.md>) | 3 | 100.00% |
| [Mechanical flag, shutter or semaphore](<../../mechanical_flag_or_shutter/quirks/segment_pixel_lamp_or_flag_failure.md>) | 1 | 33.33% |
| [Flip-element array](<../../flip_element_array/quirks/segment_pixel_lamp_or_flag_failure.md>) | 1 | 50.00% |

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Global quirk index: [gauge_quirk_index_v0.2.json](<../../_data/gauge_quirk_index_v0.2.json>)
- Group quirk index: [gauge_group_quirk_index_v0.2.json](<../../_data/gauge_group_quirk_index_v0.2.json>)

[Back to Segmented display quirks](../README.md) · [Back to canonical quirk index](../../QUIRKS.md)
