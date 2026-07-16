---
gauge_group: chart_or_trace_recorder
quirk: "dead time, missed events and duplicate counts"
catalogue_version: "0.2"
gauge_count_in_group: 1
gauge_count_global: 4
---

# Dead time, missed events and duplicate counts

**Gauge group:** [Chart or trace recorder](<../README.md>)

## Definition

Intervals or mechanisms that cause events to be ignored, delayed or counted more than once.

This is a canonical umbrella label. The exact per-gauge wording and any qualifiers below remain the authoritative detail; gauges grouped here may reach the same visible symptom through different physics.

## Frequency

| Scope | Gauges | Share | Source statements |
|---|---:|---:|---:|
| Chart or trace recorder | 1 of 6 | 16.67% | 3 |
| Entire catalogue | 4 of 136 | 2.94% | 6 |

## Gauges and preserved evidence

### [P30 — Rain gauge and tipping-bucket recorder](<../gauges/p30_rain_gauge_and_tipping_bucket_recorder.md>)

- **Measured or indicated:** Rainfall depth and rate
- **Era:** 1400s-present; tipping buckets common from late 1800s/1900s
- **Preserved source phrases:**
  - “dead time while tipping”
  - “high-rate under-read”
  - “mechanical bounce may generate duplicate electrical pulses”
- **Image references:**
  - [Wikimedia Commons images: tipping bucket rain gauge mechanism](<https://commons.wikimedia.org/w/index.php?search=tipping+bucket+rain+gauge+mechanism&title=Special:MediaSearch&type=image>)
- **Further reading:**
  - [Met Office: Rainfall observations](<https://www.metoffice.gov.uk/weather/guides/observations/how-we-measure-rainfall>)
  - [Rain gauge](<https://en.wikipedia.org/wiki/Rain_gauge>)

## Other gauge groups with this quirk

| Gauge group | Gauges with quirk | Group share |
|---|---:|---:|
| [Radial pointer](<../../radial_pointer/quirks/dead_time_missed_events_and_duplicate_counts.md>) | 2 | 3.57% |
| [Rotating scan or strobe](<../../rotating_scan_or_strobe/quirks/dead_time_missed_events_and_duplicate_counts.md>) | 1 | 25.00% |

## Reference

- Catalogue: [gauge_display_research_catalog_v0.2.json](<../../_data/gauge_display_research_catalog_v0.2.json>)
- Global quirk index: [gauge_quirk_index_v0.2.json](<../../_data/gauge_quirk_index_v0.2.json>)
- Group quirk index: [gauge_group_quirk_index_v0.2.json](<../../_data/gauge_group_quirk_index_v0.2.json>)

[Back to Chart or trace recorder quirks](../README.md) · [Back to canonical quirk index](../../QUIRKS.md)
