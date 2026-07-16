---
gauge_group: resonant_or_oscillating_element
catalogue_version: "0.2"
primary_gauge_count: 2
supporting_quirk_count: 15
---

# Resonant or oscillating element

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A vibrating reed, string or other tuned element is observed directly, usually to identify resonance, frequency or a balance condition.

**Catalogue definition:** A vibrating reed, string or other oscillating element is observed directly.

## How the group encodes a value

Value is encoded by which element vibrates most strongly, by oscillation amplitude, or by a visible beat pattern.

## Classification boundary

Use this group when the oscillating physical element is the display rather than merely an internal sensor.

## Simulation baseline

Represent resonance bandwidth, beating, amplitude build-up and decay, neighbouring-element response and source waveform dependence.

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
| `E08` | [Vibrating-reed frequency meter](<gauges/e08_vibrating_reed_frequency_meter.md>) | Frahm reed meter; switchboard 45–65 Hz indicators | AC frequency | c.1900-present; legacy power panels and generators |
| `E10` | [String galvanometer](<gauges/e10_string_galvanometer.md>) | Einthoven electrocardiograph; photographic string oscillograph | Rapid electrical waveform, especially ECG | 1901-c.1930s common; historically important |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `E08` | [Vibrating-reed frequency meter](<gauges/e08_vibrating_reed_frequency_meter.md>) | AC frequency | c.1900-present; legacy power panels and generators | 6 | 1 | 1 |
| `E10` | [String galvanometer](<gauges/e10_string_galvanometer.md>) | Rapid electrical waveform, especially ECG | 1901-c.1930s common; historically important | 10 | 1 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Resonance, beating and tuned elements](<quirks/resonance_beating_and_tuned_elements.md>) | Selective response, beat patterns and amplitude behaviour of resonant components. | 2 | 100.00% | 3 |
| [Aliasing, harmonics and false solutions](<quirks/aliasing_harmonics_and_false_solutions.md>) | Incorrect apparent readings caused by sampling, strobing, harmonics or multiple possible synchronisation points. | 1 | 50.00% | 1 |
| [Damping](<quirks/damping.md>) | Deliberate or inherent suppression of rapid movement and oscillation in the indication. | 1 | 50.00% | 1 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 1 | 50.00% | 1 |
| [Drift and long-term stability](<quirks/drift_and_long_term_stability.md>) | Slow change in indicated value or behaviour over time despite an unchanged input. | 1 | 50.00% | 1 |
| [Human-factor ambiguity and misreading](<quirks/human_factor_ambiguity_and_misreading.md>) | Display features that make an otherwise functioning instrument easy to interpret incorrectly. | 1 | 50.00% | 1 |
| [Magnetic-field, deviation and remanence effects](<quirks/magnetic_field_deviation_and_remanence_effects.md>) | Influence of external or retained magnetism on indication, zero and calibration. | 1 | 50.00% | 1 |
| [Movement torque and sensitivity](<quirks/movement_torque_and_sensitivity.md>) | The relationship between applied drive, restoring force and resulting visible movement. | 1 | 50.00% | 1 |
| [Operator procedure and ritual](<quirks/operator_procedure_and_ritual.md>) | Required handling or reading practices that materially affect the result. | 1 | 50.00% | 1 |
| [Overload, saturation and damage](<quirks/overload_saturation_and_damage.md>) | Behaviour when the input exceeds the useful range, including pegging, clipping, recovery changes or permanent harm. | 1 | 50.00% | 1 |
| [Pen, stylus and trace artefacts](<quirks/pen_stylus_and_trace_artefacts.md>) | Line-width, drag, skipping, smear, lift-off or other defects introduced by a recording or tracing element. | 1 | 50.00% | 1 |
| [Ringing and oscillation](<quirks/ringing_and_oscillation.md>) | Repeated decaying or sustained motion around a target following excitation or disturbance. | 1 | 50.00% | 1 |
| [Shock and vibration effects](<quirks/shock_and_vibration_effects.md>) | Temporary or permanent indication changes caused by mechanical shock or sustained vibration. | 1 | 50.00% | 1 |
| [Stiction and sticking](<quirks/stiction_and_sticking.md>) | Static friction or adhesion that prevents motion until enough force accumulates, often followed by a jump. | 1 | 50.00% | 1 |
| [Viewing-angle dependence](<quirks/viewing_angle_dependence.md>) | Changes in readability, colour, contrast or apparent value with observer angle. | 1 | 50.00% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
