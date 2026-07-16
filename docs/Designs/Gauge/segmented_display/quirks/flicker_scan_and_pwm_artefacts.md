---
gauge_group: segmented_display
quirk: "flicker, scan and PWM artefacts"
catalogue_version: "0.2"
gauge_count_in_group: 4
gauge_count_global: 9
---

# Flicker, scan and PWM artefacts

**Gauge group:** [Segmented display](<../README.md>)

## Definition

Visible modulation caused by multiplexing, scanning, pulse-width control or interaction with cameras and eye motion.

This is a canonical umbrella label. The exact per-gauge wording and any qualifiers below remain the authoritative detail; gauges grouped here may reach the same visible symptom through different physics.

## Frequency

| Scope | Gauges | Share | Source statements |
|---|---:|---:|---:|
| Segmented display | 4 of 13 | 30.77% | 4 |
| Entire catalogue | 9 of 136 | 6.62% | 11 |

## Gauges and preserved evidence

### [D05 — Nixie numeric cold-cathode tube](<../gauges/d05_nixie_numeric_cold_cathode_tube.md>)

- **Measured or indicated:** Numerals, decimal points and occasional symbols
- **Era:** c.1955-c.1980 mainstream; boutique revival today
- **Preserved source phrases:**
  - “high-voltage striking flicker”
- **Image references:**
  - [Wikimedia Commons images: Nixie tube digits close up](<https://commons.wikimedia.org/w/index.php?search=Nixie+tube+digits+close+up&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [Computer History Museum: Nixie tube](<https://www.computerhistory.org/collections/catalog/X1589.99C>)
  - [Nixie tube](<https://en.wikipedia.org/wiki/Nixie_tube>)

### [D13 — LED seven-segment display](<../gauges/d13_led_seven_segment_display.md>)

- **Measured or indicated:** Numeric readout
- **Era:** late 1960s-present
- **Preserved source phrases:**
  - “Multiplex flicker and scan order”
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
  - “multiplex shimmer”
- **Image references:**
  - [Wikimedia Commons images: vintage LED 16 segment alphanumeric display](<https://commons.wikimedia.org/w/index.php?search=vintage+LED+16+segment+alphanumeric+display&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [Fourteen-segment display](<https://en.wikipedia.org/wiki/Fourteen-segment_display>)
  - [Sixteen-segment display](<https://en.wikipedia.org/wiki/Sixteen-segment_display>)

### [D25 — Segmented OLED and passive-matrix OLED](<../gauges/d25_segmented_oled_and_passive_matrix_oled.md>)

- **Measured or indicated:** Numbers, icons, text and bar graphs
- **Era:** late 1990s-present
- **Preserved source phrases:**
  - “extremely fast pixels but driver may use PWM”
- **Image references:**
  - [Wikimedia Commons images: segmented OLED instrument display](<https://commons.wikimedia.org/w/index.php?search=segmented+OLED+instrument+display&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [American Chemical Society: OLED milestone](<https://www.acs.org/education/whatischemistry/landmarks/organic-light-emitting-diodes.html>)
  - [OLED](<https://en.wikipedia.org/wiki/OLED>)

## Other gauge groups with this quirk

| Gauge group | Gauges with quirk | Group share |
|---|---:|---:|
| [Dot-matrix or cell array](<../../dot_matrix_or_cell_array/quirks/flicker_scan_and_pwm_artefacts.md>) | 1 | 20.00% |
| [Rotating scan or strobe](<../../rotating_scan_or_strobe/quirks/flicker_scan_and_pwm_artefacts.md>) | 1 | 25.00% |
| [Indicator lamp or illuminated legend](<../../indicator_lamp/quirks/flicker_scan_and_pwm_artefacts.md>) | 1 | 33.33% |
| [Flip-element array](<../../flip_element_array/quirks/flicker_scan_and_pwm_artefacts.md>) | 1 | 50.00% |
| [Vector or storage trace](<../../vector_or_storage_trace/quirks/flicker_scan_and_pwm_artefacts.md>) | 1 | 50.00% |

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Global quirk index: [gauge_quirk_index_v0.2.json](<../../_data/gauge_quirk_index_v0.2.json>)
- Group quirk index: [gauge_group_quirk_index_v0.2.json](<../../_data/gauge_group_quirk_index_v0.2.json>)

[Back to Segmented display quirks](../README.md) · [Back to canonical quirk index](../../QUIRKS.md)
