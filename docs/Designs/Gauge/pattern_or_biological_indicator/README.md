---
gauge_group: pattern_or_biological_indicator
catalogue_version: "0.2"
primary_gauge_count: 2
supporting_quirk_count: 10
---

# Pattern-forming or biological indicator

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

Crystals, suspended matter, living organisms or another emergent system forms patterns or behaviours interpreted as an indication.

**Catalogue definition:** Crystals, living organisms or other emergent patterns provide the indication.

## How the group encodes a value

Information is qualitative and arises from distributed physical or biological responses rather than a calibrated pointer or numeral.

## Classification boundary

Use this group for emergent pattern-forming or biological displays whose interpretation is indirect.

## Simulation baseline

Include stochastic variation, environmental dependence, slow response, history, ambiguous interpretation and specimen-to-specimen differences.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 2 |
| Share of catalogue | 1.47% |
| Alternate members | 0 |
| Canonical quirks represented | 10 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `X22` | [Storm glass / FitzRoy weather glass](<gauges/x22_storm_glass_fitzroy_weather_glass.md>) | Camphor, salts, alcohol and water in a sealed glass | Claimed weather forecast; actually mainly temperature-driven crystallisation | c.1850s-present as curiosity/decor |
| `X27` | [Tempest Prognosticator / leech barometer](<gauges/x27_tempest_prognosticator_leech_barometer.md>) | George Merryweather’s 1851 twelve-bottle leech instrument | Claimed storm warning from animal behaviour | 1851; exhibition novelty and modern replicas |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `X22` | [Storm glass / FitzRoy weather glass](<gauges/x22_storm_glass_fitzroy_weather_glass.md>) | Claimed weather forecast; actually mainly temperature-driven crystallisation | c.1850s-present as curiosity/decor | 5 | 2 | 1 |
| `X27` | [Tempest Prognosticator / leech barometer](<gauges/x27_tempest_prognosticator_leech_barometer.md>) | Claimed storm warning from animal behaviour | 1851; exhibition novelty and modern replicas | 7 | 1 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Biological and pattern-forming variability](<quirks/biological_and_pattern_forming_variability.md>) | Variation inherent in living organisms, crystallisation or other emergent pattern-forming processes. | 2 | 100.00% | 4 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 2 | 100.00% | 2 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 1 | 50.00% | 1 |
| [Hysteresis](<quirks/hysteresis.md>) | Different indicated values for the same input depending on whether the input approached from above or below. | 1 | 50.00% | 1 |
| [Mechanical noise and cadence](<quirks/mechanical_noise_and_cadence.md>) | Audible clicks, hums, impacts or rhythms produced by the display mechanism. | 1 | 50.00% | 1 |
| [Operator procedure and ritual](<quirks/operator_procedure_and_ritual.md>) | Required handling or reading practices that materially affect the result. | 1 | 50.00% | 1 |
| [Operator signalling and acknowledgement](<quirks/operator_signalling_and_acknowledgement.md>) | Human commands, confirmations or acknowledgements that form part of the indication system. | 1 | 50.00% | 1 |
| [Qualitative or non-precision indication](<quirks/qualitative_or_non_precision_indication.md>) | An indication intended for trend, state or rough comparison rather than accurate numeric measurement. | 1 | 50.00% | 2 |
| [Random, statistical and batch variation](<quirks/random_statistical_and_batch_variation.md>) | Unpredictable events or unit-to-unit differences that are part of the observed behaviour. | 1 | 50.00% | 1 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 1 | 50.00% | 2 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
