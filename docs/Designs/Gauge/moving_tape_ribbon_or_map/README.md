---
gauge_group: moving_tape_ribbon_or_map
catalogue_version: "0.2"
primary_gauge_count: 5
supporting_quirk_count: 29
---

# Moving tape, ribbon or map

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A strip, band, tape, ribbon or printed map translates or rolls past a fixed aperture or index.

**Catalogue definition:** A strip, band, tape or printed map travels through a window or past a fixed index.

## How the group encodes a value

Value or position is encoded by the portion of a continuous medium visible in the window.

## Classification boundary

Use this group for a travelling continuous medium. Discrete numeral wheels belong under rolling_drum_or_counter.

## Simulation baseline

Model transport inertia, roller geometry, splice or seam behaviour, tension, slip, edge tracking and the narrow viewing window.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 5 |
| Share of catalogue | 3.68% |
| Alternate members | 0 |
| Canonical quirks represented | 29 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `E23` | [Ribbon and tape meter](<gauges/e23_ribbon_and_tape_meter.md>) | Moving-colour ribbon indicators; aircraft vertical tape instruments; process tape meters | Speed, altitude, temperature, level or arbitrary trend variable | 1930s-present; niche |
| `P18` | [Float-and-tape tank gauge](<gauges/p18_float_and_tape_tank_gauge.md>) | Mechanical tank gauging tape; Whessoe and Varec-style systems | Large storage-tank liquid level | late 1800s-present |
| `X20` | [Ribbon, rolling-drum and “thermometer” automotive speedometer](<gauges/x20_ribbon_rolling_drum_and_thermometer_automotive_speedometer.md>) | 1950s-1970s American ribbon speedometers; Citroën rotating drum; linear colour band displays | Road speed and sometimes odometer | 1930s-c.1980s common |
| `X24` | [Mechanical/electromechanical vertical-tape flight instrument](<gauges/x24_mechanical_electromechanical_vertical_tape_flight_instrument.md>) | Vertical-scale airspeed, altitude and engine instruments; moving printed belts | Flight or engine variable with a moving linear scale | 1940s-present; mechanical versions mainly mid-century |
| `X28` | [Mechanical moving-map and roller-map display](<gauges/x28_mechanical_moving_map_and_roller_map_display.md>) | Aircraft dead-reckoning moving maps; scrolling road-map navigators; roller chart displays | Estimated position, route progress or terrain beneath vehicle | 1930s-c.1980s common in specialist navigation |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `E23` | [Ribbon and tape meter](<gauges/e23_ribbon_and_tape_meter.md>) | Speed, altitude, temperature, level or arbitrary trend variable | 1930s-present; niche | 8 | 1 | 1 |
| `P18` | [Float-and-tape tank gauge](<gauges/p18_float_and_tape_tank_gauge.md>) | Large storage-tank liquid level | late 1800s-present | 10 | 1 | 1 |
| `X20` | [Ribbon, rolling-drum and “thermometer” automotive speedometer](<gauges/x20_ribbon_rolling_drum_and_thermometer_automotive_speedometer.md>) | Road speed and sometimes odometer | 1930s-c.1980s common | 11 | 1 | 1 |
| `X24` | [Mechanical/electromechanical vertical-tape flight instrument](<gauges/x24_mechanical_electromechanical_vertical_tape_flight_instrument.md>) | Flight or engine variable with a moving linear scale | 1940s-present; mechanical versions mainly mid-century | 7 | 1 | 1 |
| `X28` | [Mechanical moving-map and roller-map display](<gauges/x28_mechanical_moving_map_and_roller_map_display.md>) | Estimated position, route progress or terrain beneath vehicle | 1930s-c.1980s common in specialist navigation | 8 | 1 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Cable, belt, tape and roller errors](<quirks/cable_belt_tape_and_roller_errors.md>) | Slip, stretch, tracking, tension or geometry errors in flexible mechanical transmission and transport parts. | 5 | 100.00% | 10 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 3 | 60.00% | 4 |
| [Readout masking and narrow-window transitions](<quirks/readout_masking_and_narrow_window_transitions.md>) | Partial or ambiguous readings caused by apertures that reveal only a small portion of a moving scale or drum. | 3 | 60.00% | 4 |
| [Wraparound and multi-turn indication](<quirks/wraparound_and_multi_turn_indication.md>) | Behaviour when a value crosses a cyclic boundary or requires multiple revolutions to represent its full range. | 3 | 60.00% | 3 |
| [End stops, pegging and overflow](<quirks/end_stops_pegging_and_overflow.md>) | Behaviour at physical or representational range limits, including contact with stops and beyond-range indications. | 2 | 40.00% | 2 |
| [Illumination and backlighting](<quirks/illumination_and_backlighting.md>) | Lighting systems that make scales, legends or display elements visible and their associated unevenness or ageing. | 2 | 40.00% | 2 |
| [Mechanical wear, stretch and permanent set](<quirks/mechanical_wear_stretch_and_permanent_set.md>) | Long-term dimensional or elastic change in moving parts, springs, fibres or linkages. | 2 | 40.00% | 2 |
| [Operator procedure and ritual](<quirks/operator_procedure_and_ritual.md>) | Required handling or reading practices that materially affect the result. | 2 | 40.00% | 2 |
| [Stiction and sticking](<quirks/stiction_and_sticking.md>) | Static friction or adhesion that prevents motion until enough force accumulates, often followed by a jump. | 2 | 40.00% | 2 |
| [Bounce](<quirks/bounce.md>) | Repeated rebounds or reversals after a mechanical impact, contact change or rapid movement. | 1 | 20.00% | 1 |
| [Brightness nonuniformity and gradients](<quirks/brightness_nonuniformity_and_gradients.md>) | Unequal luminance across a display, character, segment, tube or field. | 1 | 20.00% | 1 |
| [Calibration, correction and compensation](<quirks/calibration_correction_and_compensation.md>) | Adjustments or correction factors required to relate the raw indication to the intended quantity. | 1 | 20.00% | 1 |
| [Carry, rollover and digit transition](<quirks/carry_rollover_and_digit_transition.md>) | The mechanical or visual sequence as counters advance between digits and propagate carries. | 1 | 20.00% | 2 |
| [Colour variation and colour shift](<quirks/colour_variation_and_colour_shift.md>) | Changes or inconsistencies in displayed colour across level, age, temperature, angle or individual units. | 1 | 20.00% | 1 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 1 | 20.00% | 1 |
| [Drift and long-term stability](<quirks/drift_and_long_term_stability.md>) | Slow change in indicated value or behaviour over time despite an unchanged input. | 1 | 20.00% | 1 |
| [Electrical or mechanical remote coupling](<quirks/electrical_or_mechanical_remote_coupling.md>) | Errors and dynamics introduced when a sensor transmits position or value to a remote indicator. | 1 | 20.00% | 1 |
| [Float and buoyancy behaviour](<quirks/float_and_buoyancy_behaviour.md>) | Motion and equilibrium effects arising from buoyant indicators in fluids. | 1 | 20.00% | 1 |
| [Fluid surge, slosh and foam](<quirks/fluid_surge_slosh_and_foam.md>) | Transient level or pressure effects caused by fluid motion, aeration or froth. | 1 | 20.00% | 1 |
| [Flutter, jitter, tremor and quiver](<quirks/flutter_jitter_tremor_and_quiver.md>) | Small rapid random or periodic movements around the nominal indication. | 1 | 20.00% | 1 |
| [Friction and drag](<quirks/friction_and_drag.md>) | Motion resistance that slows, biases or distorts the indication. | 1 | 20.00% | 1 |
| [Human-factor ambiguity and misreading](<quirks/human_factor_ambiguity_and_misreading.md>) | Display features that make an otherwise functioning instrument easy to interpret incorrectly. | 1 | 20.00% | 2 |
| [Manual reset, tare and caging](<quirks/manual_reset_tare_and_caging.md>) | User-operated mechanisms that clear, zero, restrain or protect the indication. | 1 | 20.00% | 1 |
| [Map, dead-reckoning and position drift](<quirks/map_dead_reckoning_and_position_drift.md>) | Accumulated position error in moving-map or dead-reckoning displays. | 1 | 20.00% | 5 |
| [Optical distortion and refraction](<quirks/optical_distortion_and_refraction.md>) | Apparent displacement or shape changes caused by lenses, glass, liquid, curved windows or refractive interfaces. | 1 | 20.00% | 1 |
| [Overshoot](<quirks/overshoot.md>) | Temporary travel beyond the final steady indication after an input change. | 1 | 20.00% | 1 |
| [Parallax](<quirks/parallax.md>) | Apparent reading error caused by viewing the indicator and scale from the wrong angle or depth relationship. | 1 | 20.00% | 1 |
| [Rate limiting and motion limits](<quirks/rate_limiting_and_motion_limits.md>) | Restrictions on how quickly or how far an indication may move. | 1 | 20.00% | 1 |
| [Shadows, depth and occlusion](<quirks/shadows_depth_and_occlusion.md>) | Visual effects caused by layered parts blocking, shading or appearing at different depths. | 1 | 20.00% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
