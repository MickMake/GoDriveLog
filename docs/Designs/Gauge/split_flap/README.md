---
gauge_group: split_flap
catalogue_version: "0.2"
primary_gauge_count: 1
supporting_quirk_count: 9
---

# Split-flap display

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A character position consists of mechanically sequenced flaps that rotate through intermediate symbols until the requested character is reached.

**Catalogue definition:** Physical character flaps cycle through intermediate positions to reach the target.

## How the group encodes a value

Each module displays one character from a fixed ordered set; multiple modules form words or numbers.

## Classification boundary

Use this group for split-flap character modules, not generic flip-dot arrays or rotating numeral counters.

## Simulation baseline

The route to the target matters: include sequential intermediate characters, module timing differences, flap impact, bounce and acoustic cadence.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 1 |
| Share of catalogue | 0.74% |
| Alternate members | 0 |
| Canonical quirks represented | 9 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `E21` | [Split-flap display](<gauges/e21_split_flap_display.md>) | Solari Cifra 5; Solari Udine airport boards; Gino Valle designs | Text, time, destinations, prices or status | 1956-present; peak public use c.1960s-1990s |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `E21` | [Split-flap display](<gauges/e21_split_flap_display.md>) | Text, time, destinations, prices or status | 1956-present; peak public use c.1960s-1990s | 9 | 2 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Ageing and material degradation](<quirks/ageing_and_material_degradation.md>) | Progressive change due to wear, fatigue, chemical decay, phosphor loss, embrittlement or similar ageing processes. | 1 | 100.00% | 1 |
| [Bounce](<quirks/bounce.md>) | Repeated rebounds or reversals after a mechanical impact, contact change or rapid movement. | 1 | 100.00% | 1 |
| [Channel mismatch and unequal dynamics](<quirks/channel_mismatch_and_unequal_dynamics.md>) | Different calibration, response or motion between channels intended to behave alike. | 1 | 100.00% | 1 |
| [Latching and state-retention behaviour](<quirks/latching_and_state_retention_behaviour.md>) | A displayed state that remains mechanically, magnetically, electrically or optically retained until reset or rewritten. | 1 | 100.00% | 1 |
| [Mechanical noise and cadence](<quirks/mechanical_noise_and_cadence.md>) | Audible clicks, hums, impacts or rhythms produced by the display mechanism. | 1 | 100.00% | 1 |
| [Overshoot](<quirks/overshoot.md>) | Temporary travel beyond the final steady indication after an input change. | 1 | 100.00% | 1 |
| [Power-off and power-loss behaviour](<quirks/power_off_and_power_loss_behaviour.md>) | What the indication does when drive power is removed, including retained, blank, parked or misleading states. | 1 | 100.00% | 1 |
| [Refresh, erase and update artefacts](<quirks/refresh_erase_and_update_artefacts.md>) | Visible effects produced while changing, clearing or rewriting a display. | 1 | 100.00% | 2 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 1 | 100.00% | 2 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
