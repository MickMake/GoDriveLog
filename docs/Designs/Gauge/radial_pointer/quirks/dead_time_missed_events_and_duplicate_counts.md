---
gauge_group: radial_pointer
quirk: "dead time, missed events and duplicate counts"
catalogue_version: "0.2"
gauge_count_in_group: 2
gauge_count_global: 4
---

# Dead time, missed events and duplicate counts

**Gauge group:** [Radial pointer](<../README.md>)

## Definition

Intervals or mechanisms that cause events to be ignored, delayed or counted more than once.

This is a canonical umbrella label. The exact per-gauge wording and any qualifiers below remain the authoritative detail; gauges grouped here may reach the same visible symptom through different physics.

## Frequency

| Scope | Gauges | Share | Source statements |
|---|---:|---:|---:|
| Radial pointer | 2 of 56 | 3.57% | 2 |
| Entire catalogue | 4 of 136 | 2.94% | 6 |

## Gauges and preserved evidence

### [E15 — Stepper-motor gauge](<../gauges/e15_stepper_motor_gauge.md>)

- **Measured or indicated:** Digitally commanded speed, RPM, fuel, temperature and other values
- **Era:** c.1980s-present
- **Preserved source phrases:**
  - “missed steps and lost zero”
- **Image references:**
  - [Wikimedia Commons images: Switec X27.168 gauge stepper motor](<https://commons.wikimedia.org/w/index.php?search=Switec+X27.168+gauge+stepper+motor&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [Adafruit: Switec/Juken X27.168 gauge stepper specifications](<https://www.adafruit.com/product/2424>)

### [X19 — Geiger-counter analog ratemeter](<../gauges/x19_geiger_counter_analog_ratemeter.md>)

- **Measured or indicated:** Radiation count rate or inferred dose rate
- **Era:** 1928-present
- **Preserved source phrases:**
  - “detector dead time causes high-rate nonlinearity or saturation”
- **Image references:**
  - [Wikimedia Commons images: vintage Geiger counter analog meter](<https://commons.wikimedia.org/w/index.php?search=vintage+Geiger+counter+analog+meter&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [US NRC: Geiger-Mueller counter test procedure](<https://www.nrc.gov/docs/ML0037/ML003739460.pdf>)
  - [Geiger counter](<https://en.wikipedia.org/wiki/Geiger_counter>)

## Other gauge groups with this quirk

| Gauge group | Gauges with quirk | Group share |
|---|---:|---:|
| [Chart or trace recorder](<../../chart_or_trace_recorder/quirks/dead_time_missed_events_and_duplicate_counts.md>) | 1 | 16.67% |
| [Rotating scan or strobe](<../../rotating_scan_or_strobe/quirks/dead_time_missed_events_and_duplicate_counts.md>) | 1 | 25.00% |

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Global quirk index: [gauge_quirk_index_v0.2.json](<../../_data/gauge_quirk_index_v0.2.json>)
- Group quirk index: [gauge_group_quirk_index_v0.2.json](<../../_data/gauge_group_quirk_index_v0.2.json>)

[Back to Radial pointer quirks](../README.md) · [Back to canonical quirk index](../../QUIRKS.md)
