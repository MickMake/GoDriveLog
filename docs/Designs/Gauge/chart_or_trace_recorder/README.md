---
gauge_group: chart_or_trace_recorder
catalogue_version: "0.2"
primary_gauge_count: 6
supporting_quirk_count: 35
---

# Chart or trace recorder

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A pen, stylus, scratch point, light beam or mechanical marker records a changing value on moving paper, film, smoked surface or another physical medium.

**Catalogue definition:** A pen, stylus, light beam or scratch records a changing value on a moving medium.

## How the group encodes a value

The instantaneous value is represented by trace position while time is represented by medium motion, drum rotation or chart angle.

## Classification boundary

Use this group when a persistent physical record is a primary output. A non-recording CRT trace belongs under vector_or_storage_trace.

## Simulation baseline

Model both axes: sensor or pen dynamics and chart transport. Include line width, drag, gaps, clock error, paper seams and channel alignment.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 6 |
| Share of catalogue | 4.41% |
| Alternate members | 2 |
| Canonical quirks represented | 35 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `P30` | [Rain gauge and tipping-bucket recorder](<gauges/p30_rain_gauge_and_tipping_bucket_recorder.md>) | Standard cylindrical gauge; siphon gauge; tipping bucket with reed switch | Rainfall depth and rate | 1400s-present; tipping buckets common from late 1800s/1900s |
| `P42` | [Strip-chart and circular-chart recorder](<gauges/p42_strip_chart_and_circular_chart_recorder.md>) | Bristol, Honeywell and Foxboro recorders; clockwork thermographs and barographs | One or more variables over time | late 1800s-present; now specialist/legacy |
| `X30` | [Mechanical float tide gauge and stilling-well recorder](<gauges/x30_mechanical_float_tide_gauge_and_stilling_well_recorder.md>) | Standard tide gauge with float, pulley, clock and pen drum; tide staff | Sea level versus time | c.1830s-c.1980s dominant; historical and some backup use |
| `X31` | [Mechanical seismograph and smoked-paper drum](<gauges/x31_mechanical_seismograph_and_smoked_paper_drum.md>) | Wiechert, Milne and Wood-Anderson instruments; helicorders | Ground displacement, velocity or acceleration over time | 1880s-present; mechanical photographic forms mainly historical |
| `X32` | [Railway speed indicator and recorder](<gauges/x32_railway_speed_indicator_and_recorder.md>) | Hasler/Teloc and Hasler-Hausshalter instruments; Flaman recorders | Locomotive speed plus a tamper-resistant time history | 1891-present; electromechanical RT/A forms from 1930s |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `P30` | [Rain gauge and tipping-bucket recorder](<gauges/p30_rain_gauge_and_tipping_bucket_recorder.md>) | Rainfall depth and rate | 1400s-present; tipping buckets common from late 1800s/1900s | 6 | 2 | 1 |
| `P42` | [Strip-chart and circular-chart recorder](<gauges/p42_strip_chart_and_circular_chart_recorder.md>) | One or more variables over time | late 1800s-present; now specialist/legacy | 12 | 1 | 1 |
| `X30` | [Mechanical float tide gauge and stilling-well recorder](<gauges/x30_mechanical_float_tide_gauge_and_stilling_well_recorder.md>) | Sea level versus time | c.1830s-c.1980s dominant; historical and some backup use | 15 | 2 | 1 |
| `X31` | [Mechanical seismograph and smoked-paper drum](<gauges/x31_mechanical_seismograph_and_smoked_paper_drum.md>) | Ground displacement, velocity or acceleration over time | 1880s-present; mechanical photographic forms mainly historical | 8 | 1 | 1 |
| `X32` | [Railway speed indicator and recorder](<gauges/x32_railway_speed_indicator_and_recorder.md>) | Locomotive speed plus a tamper-resistant time history | 1891-present; electromechanical RT/A forms from 1930s | 9 | 2 | 1 |
| `X35` | [Tachograph / recording speedometer](<gauges/x35_tachograph_recording_speedometer.md>) | Vehicle speed, distance, time and driver activity history | 1920s-present; paper-disc peak c.1950s-2000s | 9 | 1 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `E10` | [String galvanometer](<../resonant_or_oscillating_element/gauges/e10_string_galvanometer.md>) | [Resonant or oscillating element](<../resonant_or_oscillating_element/README.md>) |
| `P26` | [Hair hygrometer and hygrograph](<../radial_pointer/gauges/p26_hair_hygrometer_and_hygrograph.md>) | [Radial pointer](<../radial_pointer/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Chart and paper transport](<quirks/chart_and_paper_transport.md>) | Timing, speed, alignment and mechanical behaviour of the recording medium. | 5 | 83.33% | 6 |
| [Pen, stylus and trace artefacts](<quirks/pen_stylus_and_trace_artefacts.md>) | Line-width, drag, skipping, smear, lift-off or other defects introduced by a recording or tracing element. | 5 | 83.33% | 11 |
| [Clock drift, timing and seams](<quirks/clock_drift_timing_and_seams.md>) | Errors caused by imperfect timebase speed, chart joins, scan boundaries or repeating mechanical cycles. | 4 | 66.67% | 6 |
| [Recording alignment and channel offset](<quirks/recording_alignment_and_channel_offset.md>) | Misregistration between recorded channels, time axes, pens or reference lines. | 4 | 66.67% | 6 |
| [Drift and long-term stability](<quirks/drift_and_long_term_stability.md>) | Slow change in indicated value or behaviour over time despite an unchanged input. | 3 | 50.00% | 4 |
| [Channel mismatch and unequal dynamics](<quirks/channel_mismatch_and_unequal_dynamics.md>) | Different calibration, response or motion between channels intended to behave alike. | 2 | 33.33% | 2 |
| [Damping](<quirks/damping.md>) | Deliberate or inherent suppression of rapid movement and oscillation in the indication. | 2 | 33.33% | 2 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 2 | 33.33% | 3 |
| [Evaporation and drying](<quirks/evaporation_and_drying.md>) | Loss of liquid or solvent that changes level, concentration, response or appearance. | 2 | 33.33% | 2 |
| [Friction and drag](<quirks/friction_and_drag.md>) | Motion resistance that slows, biases or distorts the indication. | 2 | 33.33% | 2 |
| [Operator procedure and ritual](<quirks/operator_procedure_and_ritual.md>) | Required handling or reading practices that materially affect the result. | 2 | 33.33% | 3 |
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 2 | 33.33% | 2 |
| [Tamper evidence and official seals](<quirks/tamper_evidence_and_official_seals.md>) | Marks, seals or mechanisms that reveal unauthorised adjustment or preserve legal metrology status. | 2 | 33.33% | 2 |
| [Backlash and lash](<quirks/backlash_and_lash.md>) | Lost motion when direction reverses because of clearance in gears, linkages or drive parts. | 1 | 16.67% | 1 |
| [Blockage and restricted passages](<quirks/blockage_and_restricted_passages.md>) | Errors or slow response caused by obstructed tubing, ports, capillaries or vents. | 1 | 16.67% | 1 |
| [Bounce](<quirks/bounce.md>) | Repeated rebounds or reversals after a mechanical impact, contact change or rapid movement. | 1 | 16.67% | 1 |
| [Cable, belt, tape and roller errors](<quirks/cable_belt_tape_and_roller_errors.md>) | Slip, stretch, tracking, tension or geometry errors in flexible mechanical transmission and transport parts. | 1 | 16.67% | 1 |
| [Calibration, correction and compensation](<quirks/calibration_correction_and_compensation.md>) | Adjustments or correction factors required to relate the raw indication to the intended quantity. | 1 | 16.67% | 1 |
| [Chatter](<quirks/chatter.md>) | Rapid repeated switching or movement near a threshold or unstable equilibrium. | 1 | 16.67% | 1 |
| [Contamination, dirt and fouling](<quirks/contamination_dirt_and_fouling.md>) | Reading or appearance changes caused by deposits, dust, oxidation, residue or biological growth. | 1 | 16.67% | 1 |
| [Creep](<quirks/creep.md>) | Very slow movement under a constant input, load or retained state. | 1 | 16.67% | 1 |
| [Dead time, missed events and duplicate counts](<quirks/dead_time_missed_events_and_duplicate_counts.md>) | Intervals or mechanisms that cause events to be ignored, delayed or counted more than once. | 1 | 16.67% | 3 |
| [Electrical or mechanical remote coupling](<quirks/electrical_or_mechanical_remote_coupling.md>) | Errors and dynamics introduced when a sensor transmits position or value to a remote indicator. | 1 | 16.67% | 1 |
| [End stops, pegging and overflow](<quirks/end_stops_pegging_and_overflow.md>) | Behaviour at physical or representational range limits, including contact with stops and beyond-range indications. | 1 | 16.67% | 1 |
| [Float and buoyancy behaviour](<quirks/float_and_buoyancy_behaviour.md>) | Motion and equilibrium effects arising from buoyant indicators in fluids. | 1 | 16.67% | 1 |
| [Fluid surge, slosh and foam](<quirks/fluid_surge_slosh_and_foam.md>) | Transient level or pressure effects caused by fluid motion, aeration or froth. | 1 | 16.67% | 1 |
| [Inertia and coasting](<quirks/inertia_and_coasting.md>) | Continued or delayed motion due to mass and stored kinetic energy. | 1 | 16.67% | 1 |
| [Input slip and wheel-slip error](<quirks/input_slip_and_wheel_slip_error.md>) | Mismatch between actual movement and sensed movement due to slipping at a wheel, roller or drive interface. | 1 | 16.67% | 1 |
| [Mechanical noise and cadence](<quirks/mechanical_noise_and_cadence.md>) | Audible clicks, hums, impacts or rhythms produced by the display mechanism. | 1 | 16.67% | 1 |
| [Range switching and multiple ranges](<quirks/range_switching_and_multiple_ranges.md>) | Selection or combination of alternate scales, sensitivities or measuring ranges. | 1 | 16.67% | 1 |
| [Ringing and oscillation](<quirks/ringing_and_oscillation.md>) | Repeated decaying or sustained motion around a target following excitation or disturbance. | 1 | 16.67% | 1 |
| [Scale markings, zones and legends](<quirks/scale_markings_zones_and_legends.md>) | Visual information carried by ticks, numerals, colour bands, labels and operating zones. | 1 | 16.67% | 1 |
| [Thresholds, deadband and switching points](<quirks/thresholds_deadband_and_switching_points.md>) | Regions or levels where no change occurs, or where a discrete state changes with defined or variable thresholds. | 1 | 16.67% | 1 |
| [Wind, splash and collection error](<quirks/wind_splash_and_collection_error.md>) | Measurement error caused by airflow, splashing, wetting or imperfect collection of precipitation or fluid. | 1 | 16.67% | 2 |
| [Zero drift and offset](<quirks/zero_drift_and_offset.md>) | A non-zero indication at the true zero point, including offsets that change with time or conditions. | 1 | 16.67% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
