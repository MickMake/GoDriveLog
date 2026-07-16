---
gauge_group: radial_pointer
quirk: "invalid, out-of-range and warning flags"
catalogue_version: "0.2"
gauge_count_in_group: 4
gauge_count_global: 7
---

# Invalid, out-of-range and warning flags

**Gauge group:** [Radial pointer](<../README.md>)

## Definition

Explicit indications that a reading is unavailable, unreliable, unsafe or beyond the valid range.

This is a canonical umbrella label. The exact per-gauge wording and any qualifiers below remain the authoritative detail; gauges grouped here may reach the same visible symptom through different physics.

## Frequency

| Scope | Gauges | Share | Source statements |
|---|---:|---:|---:|
| Radial pointer | 4 of 56 | 7.14% | 5 |
| Entire catalogue | 7 of 136 | 5.15% | 8 |

## Gauges and preserved evidence

### [E17 — Servo and torque-motor indicator](<../gauges/e17_servo_and_torque_motor_indicator.md>)

- **Measured or indicated:** Remote position or electrically computed variable
- **Era:** 1930s-present
- **Preserved source phrases:**
  - “power loss may drop a warning flag or let the pointer drift”
- **Image references:**
  - [Wikimedia Commons images: servo motor analog indicator gauge](<https://commons.wikimedia.org/w/index.php?search=servo+motor+analog+indicator+gauge&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [Servomechanism](<https://en.wikipedia.org/wiki/Servomechanism>)

### [X02 — Counter-pointer altimeter](<../gauges/x02_counter_pointer_altimeter.md>)

- **Measured or indicated:** Pressure altitude
- **Era:** c.1940s-present; common in transport and military aircraft mid-century
- **Preserved source phrases:**
  - “mechanical flag for negative/out-of-range”
- **Image references:**
  - [Wikimedia Commons images: counter pointer aircraft altimeter](<https://commons.wikimedia.org/w/index.php?search=counter+pointer+aircraft+altimeter&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [NASA: Altimeter display types and reading performance](<https://ntrs.nasa.gov/citations/19660017807>)
  - [FAA: Pilot’s Handbook of Aeronautical Knowledge, Flight Instruments](<https://www.faa.gov/regulations_policies/handbooks_manuals/aviation/phak>)

### [X03 — Drum-pointer and counter-drum altimeter](<../gauges/x03_drum_pointer_and_counter_drum_altimeter.md>)

- **Measured or indicated:** Pressure altitude
- **Era:** c.1950s-present
- **Preserved source phrases:**
  - “drum can be masked until valid”
  - “baro setting and low-altitude warning sectors”
- **Image references:**
  - [Wikimedia Commons images: counter drum pointer altimeter](<https://commons.wikimedia.org/w/index.php?search=counter+drum+pointer+altimeter&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [NASA: Counter-drum-pointer altimeter study](<https://ntrs.nasa.gov/citations/19730011592>)
  - [NASA: Altimeter display types and reading performance](<https://ntrs.nasa.gov/citations/19660017807>)

### [X08 — Radar altimeter dial](<../gauges/x08_radar_altimeter_dial.md>)

- **Measured or indicated:** Height above terrain
- **Era:** 1940s-present
- **Preserved source phrases:**
  - “OFF/invalid flag”
- **Image references:**
  - [Wikimedia Commons images: vintage radar altimeter dial decision height bug](<https://commons.wikimedia.org/w/index.php?search=vintage+radar+altimeter+dial+decision+height+bug&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [Radar altimeter](<https://en.wikipedia.org/wiki/Radar_altimeter>)

## Other gauge groups with this quirk

| Gauge group | Gauges with quirk | Group share |
|---|---:|---:|
| [Mechanical flag, shutter or semaphore](<../../mechanical_flag_or_shutter/quirks/invalid_out_of_range_and_warning_flags.md>) | 1 | 33.33% |
| [Rotating scale or scene](<../../rotating_scale_or_scene/quirks/invalid_out_of_range_and_warning_flags.md>) | 2 | 66.67% |

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Global quirk index: [gauge_quirk_index_v0.2.json](<../../_data/gauge_quirk_index_v0.2.json>)
- Group quirk index: [gauge_group_quirk_index_v0.2.json](<../../_data/gauge_group_quirk_index_v0.2.json>)

[Back to Radial pointer quirks](../README.md) · [Back to canonical quirk index](../../QUIRKS.md)
