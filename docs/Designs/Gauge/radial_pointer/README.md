---
gauge_group: radial_pointer
catalogue_version: "0.2"
primary_gauge_count: 56
supporting_quirk_count: 109
---

# Radial pointer

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A pointer or hand rotates about a pivot and indicates a value against a fixed circular or arc-shaped scale. The pointer may be directly driven, linked through gears, or positioned by an electrical movement or servo.

**Catalogue definition:** One or more hands sweep an arc or circle against a fixed scale.

## How the group encodes a value

Value is encoded primarily as angular position. Multiple hands, coloured arcs, counter windows, tell-tales and contact markers may add range, state or history.

## Classification boundary

Use this group when the moving index is the principal readout. A rotating card or scale behind a fixed index belongs under rotating_scale_or_scene instead.

## Simulation baseline

Treat pointer angle, acceleration, damping, stops, backlash and reading geometry as separate concerns. A perfectly smooth interpolation is usually the least convincing option.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 56 |
| Share of catalogue | 41.18% |
| Alternate members | 5 |
| Canonical quirks represented | 109 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `E01` | [D’Arsonval / Weston moving-coil meter](<gauges/e01_darsonval_weston_moving_coil_meter.md>) | Weston Model 60; Simpson 260 movement; panel ammeters and voltmeters | DC current or any transduced quantity | 1880s-present |
| `E02` | [Taut-band meter movement](<gauges/e02_taut_band_meter_movement.md>) | Suspended moving coil without pivots; high-grade laboratory and aerospace panel meters | Current, voltage or transduced variables | mid-1900s-present; now specialist |
| `E03` | [Moving-iron meter](<gauges/e03_moving_iron_meter.md>) | Attraction and repulsion types; AC switchboard ammeters and voltmeters | AC or DC current/voltage | 1880s-present |
| `E04` | [Electrodynamometer instrument](<gauges/e04_electrodynamometer_instrument.md>) | Siemens electrodynamometer; dynamometer wattmeter | AC/DC current, voltage or true power | 1880-c.1930s common; precision descendants remain |
| `E05` | [Hot-wire ammeter](<gauges/e05_hot_wire_ammeter.md>) | Thermal-expansion wire meters in early radio and switchboards | RMS-like AC/RF current | 1880s-c.1940s common; niche later |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `E01` | [D’Arsonval / Weston moving-coil meter](<gauges/e01_darsonval_weston_moving_coil_meter.md>) | DC current or any transduced quantity | 1880s-present | 10 | 2 | 1 |
| `E02` | [Taut-band meter movement](<gauges/e02_taut_band_meter_movement.md>) | Current, voltage or transduced variables | mid-1900s-present; now specialist | 7 | 1 | 1 |
| `E03` | [Moving-iron meter](<gauges/e03_moving_iron_meter.md>) | AC or DC current/voltage | 1880s-present | 7 | 1 | 1 |
| `E04` | [Electrodynamometer instrument](<gauges/e04_electrodynamometer_instrument.md>) | AC/DC current, voltage or true power | 1880-c.1930s common; precision descendants remain | 7 | 2 | 1 |
| `E05` | [Hot-wire ammeter](<gauges/e05_hot_wire_ammeter.md>) | RMS-like AC/RF current | 1880s-c.1940s common; niche later | 6 | 1 | 1 |
| `E06` | [Thermocouple RF meter](<gauges/e06_thermocouple_rf_meter.md>) | RF current or power via heating | 1910s-present in specialist RF metrology | 9 | 2 | 1 |
| `E07` | [Electrostatic voltmeter](<gauges/e07_electrostatic_voltmeter.md>) | High voltage or electrostatic potential | 1870s-present; mainly laboratory/high-voltage niche | 9 | 1 | 1 |
| `E12` | [Cross-coil / ratiometer gauge](<gauges/e12_cross_coil_ratiometer_gauge.md>) | Resistance ratio, fuel level, temperature or pressure sender output | 1910s-present; dominant automotive form c.1930s-1970s | 7 | 1 | 1 |
| `E13` | [Thermal bimetal automotive gauge](<gauges/e13_thermal_bimetal_automotive_gauge.md>) | Fuel level, coolant temperature, oil pressure by sender resistance | 1920s-c.1980s common; restoration market remains | 10 | 1 | 1 |
| `E14` | [Air-core gauge](<gauges/e14_air_core_gauge.md>) | Vehicle speed, fuel, temperature or arbitrary angular command | 1960s-present; common c.1970s-1990s | 6 | 1 | 1 |
| `E15` | [Stepper-motor gauge](<gauges/e15_stepper_motor_gauge.md>) | Digitally commanded speed, RPM, fuel, temperature and other values | c.1980s-present | 9 | 1 | 1 |
| `E16` | [Synchro / selsyn remote indicator](<gauges/e16_synchro_selsyn_remote_indicator.md>) | Remote angular position such as heading, valve position or antenna bearing | 1910s-present; peak use c.1930s-1970s | 7 | 1 | 1 |
| `E17` | [Servo and torque-motor indicator](<gauges/e17_servo_and_torque_motor_indicator.md>) | Remote position or electrically computed variable | 1930s-present | 10 | 1 | 1 |
| `E19` | [Contact meter relay / limit meter](<gauges/e19_contact_meter_relay_limit_meter.md>) | Analog value plus one or more control thresholds | 1930s-present | 7 | 1 | 1 |
| `E25` | [Concentric multi-pointer instrument](<gauges/e25_concentric_multi_pointer_instrument.md>) | Several related variables in one aperture | 1910s-present | 5 | 1 | 1 |
| `P03` | [Aneroid barometer and pressure capsule mechanism](<gauges/p03_aneroid_barometer_and_pressure_capsule_mechanism.md>) | Atmospheric pressure or altitude via pressure | 1844-present | 10 | 2 | 1 |
| `P04` | [Bourdon-tube pressure gauge, including C, spiral and helical tubes](<gauges/p04_bourdon_tube_pressure_gauge_including_c_spiral_and_helical_tubes.md>) | Gauge, absolute, vacuum or compound pressure | 1849-present | 13 | 2 | 1 |
| `P05` | [Diaphragm pressure gauge](<gauges/p05_diaphragm_pressure_gauge.md>) | Low pressure, differential pressure, viscous or corrosive process pressure | late 1800s-present | 10 | 2 | 1 |
| `P06` | [Capsule and bellows pressure gauges](<gauges/p06_capsule_and_bellows_pressure_gauges.md>) | Very low gas pressure, differential pressure, altitude or displacement | late 1800s-present | 7 | 2 | 1 |
| `P07` | [Magnetic-coupled differential gauge](<gauges/p07_magnetic_coupled_differential_gauge.md>) | Low differential pressure, filter loading, clean-room pressure and air velocity | 1953-present | 11 | 1 | 1 |
| `P08` | [Compound, duplex and twin-pointer pressure gauges](<gauges/p08_compound_duplex_and_twin_pointer_pressure_gauges.md>) | Two pressures, pressure plus vacuum, or one value over a signed range | late 1800s-present | 9 | 2 | 1 |
| `P09` | [Liquid-filled, snubbed and damped pressure gauge](<gauges/p09_liquid_filled_snubbed_and_damped_pressure_gauge.md>) | Pressure in vibrating or pulsating machinery | c.1930s-present | 10 | 1 | 1 |
| `P10` | [Contact, alarm and meter-relay gauge](<gauges/p10_contact_alarm_and_meter_relay_gauge.md>) | Measured value plus alarm or control threshold | early 1900s-present | 8 | 1 | 1 |
| `P11` | [Drag, tell-tale and witness-pointer gauge](<gauges/p11_drag_tell_tale_and_witness_pointer_gauge.md>) | Maximum, minimum or both extrema of pressure, speed, temperature or load | late 1800s-present | 5 | 1 | 1 |
| `P13` | [Thermal-conductivity vacuum gauges: Pirani and thermocouple](<gauges/p13_thermal_conductivity_vacuum_gauges_pirani_and_thermocouple.md>) | Rough and medium vacuum inferred from heat loss | 1906-present | 10 | 2 | 1 |
| `P14` | [Cold-cathode ionisation gauge](<gauges/p14_cold_cathode_ionisation_gauge.md>) | High vacuum | 1937-present | 7 | 2 | 1 |
| `P15` | [Hot-cathode ionisation gauge](<gauges/p15_hot_cathode_ionisation_gauge.md>) | High and ultra-high vacuum | 1950-present | 8 | 2 | 1 |
| `P16` | [Capacitance diaphragm manometer](<gauges/p16_capacitance_diaphragm_manometer.md>) | Absolute or differential pressure, especially process vacuum | c.1960s-present | 9 | 1 | 1 |
| `P23` | [Bimetal dial thermometer](<gauges/p23_bimetal_dial_thermometer.md>) | Temperature | 1800s-present | 8 | 2 | 1 |
| `P24` | [Filled-system capillary thermometer](<gauges/p24_filled_system_capillary_thermometer.md>) | Remote temperature | late 1800s-present | 9 | 1 | 1 |
| `P26` | [Hair hygrometer and hygrograph](<gauges/p26_hair_hygrometer_and_hygrograph.md>) | Relative humidity | 1783-present; largely meteorological/historical today | 11 | 2 | 1 |
| `P31` | [Cup, propeller and vane anemometers](<gauges/p31_cup_propeller_and_vane_anemometers.md>) | Wind speed and direction | 1846-present | 10 | 1 | 1 |
| `P34` | [Diving depth gauge: capillary and Bourdon](<gauges/p34_diving_depth_gauge_capillary_and_bourdon.md>) | Water depth | c.1940s-present; capillary types now niche/back-up | 6 | 2 | 1 |
| `P35` | [Mercury and aneroid sphygmomanometer](<gauges/p35_mercury_and_aneroid_sphygmomanometer.md>) | Cuff pressure and blood pressure during auscultation | 1896-present; mercury use declining | 9 | 1 | 1 |
| `P37` | [Spring balance, dial force gauge and dynamometer](<gauges/p37_spring_balance_dial_force_gauge_and_dynamometer.md>) | Force, weight or tension | 1700s-present | 11 | 1 | 1 |
| `P38` | [Dial indicator, test indicator and dial torque wrench](<gauges/p38_dial_indicator_test_indicator_and_dial_torque_wrench.md>) | Small displacement, runout or applied torque | late 1800s-present | 9 | 2 | 1 |
| `P39` | [Mechanical tachometer and eddy-current speedometer](<gauges/p39_mechanical_tachometer_and_eddy_current_speedometer.md>) | Rotational speed or road speed | 1817-present; mechanical automotive versions dominant c.1900s-1980s | 11 | 2 | 1 |
| `P40` | [Chronometric tachometer](<gauges/p40_chronometric_tachometer.md>) | Rotational speed averaged over timed samples | c.1920s-1960s common; still restored and reproduced | 7 | 2 | 1 |
| `X01` | [Three-pointer altimeter](<gauges/x01_three_pointer_altimeter.md>) | Pressure altitude | c.1920s-c.1970s common; some remain in service | 6 | 2 | 1 |
| `X02` | [Counter-pointer altimeter](<gauges/x02_counter_pointer_altimeter.md>) | Pressure altitude | c.1940s-present; common in transport and military aircraft mid-century | 7 | 2 | 1 |
| `X03` | [Drum-pointer and counter-drum altimeter](<gauges/x03_drum_pointer_and_counter_drum_altimeter.md>) | Pressure altitude | c.1950s-present | 9 | 2 | 1 |
| `X04` | [Airspeed indicator with coloured operating arcs](<gauges/x04_airspeed_indicator_with_coloured_operating_arcs.md>) | Indicated airspeed derived from dynamic pressure | c.1910-present | 8 | 2 | 1 |
| `X05` | [Vertical-speed indicator](<gauges/x05_vertical_speed_indicator.md>) | Rate of climb or descent | c.1920-present | 8 | 2 | 1 |
| `X07` | [Turn-and-bank, turn coordinator and slip-skid ball](<gauges/x07_turn_and_bank_turn_coordinator_and_slip_skid_ball.md>) | Turn rate and coordination/lateral acceleration | 1910s-present | 8 | 2 | 1 |
| `X08` | [Radar altimeter dial](<gauges/x08_radar_altimeter_dial.md>) | Height above terrain | 1940s-present | 8 | 1 | 1 |
| `X10` | [Aircraft dual- and triple-pointer engine gauge](<gauges/x10_aircraft_dual_and_triple_pointer_engine_gauge.md>) | Two or three engine or propeller variables | 1930s-present | 7 | 1 | 1 |
| `X11` | [Railway duplex brake and vacuum gauge](<gauges/x11_railway_duplex_brake_and_vacuum_gauge.md>) | Brake-pipe, reservoir, cylinder or vacuum pressures | 1870s-present on heritage and some service stock | 6 | 2 | 1 |
| `X13` | [Engine-order telegraph / Chadburn](<gauges/x13_engine_order_telegraph_chadburn.md>) | Commanded engine direction/speed and acknowledgement | 1870s-c.1950s dominant; electronic descendants remain | 5 | 2 | 1 |
| `X16` | [VU meter / standard volume indicator](<gauges/x16_vu_meter_standard_volume_indicator.md>) | Perceived programme level, not instantaneous peak | 1939-present | 7 | 2 | 1 |
| `X17` | [Peak programme meter](<gauges/x17_peak_programme_meter.md>) | Quasi-peak audio programme level | 1930s-present | 7 | 2 | 1 |
| `X18` | [Radio S-meter](<gauges/x18_radio_s_meter.md>) | Relative received radio signal strength | 1930s-present | 6 | 2 | 1 |
| `X19` | [Geiger-counter analog ratemeter](<gauges/x19_geiger_counter_analog_ratemeter.md>) | Radiation count rate or inferred dose rate | 1928-present | 12 | 2 | 1 |
| `X21` | [Automotive econometer / manifold-vacuum gauge](<gauges/x21_automotive_econometer_manifold_vacuum_gauge.md>) | Engine manifold vacuum as a rough load/economy proxy | c.1930s-c.1980s common; aftermarket niche remains | 7 | 1 | 1 |
| `X25` | [Follow-the-pointer photographic exposure meter](<gauges/x25_follow_the_pointer_photographic_exposure_meter.md>) | Scene luminance, exposure value, aperture/shutter combinations | 1930s-c.1990s common; analog models remain | 9 | 1 | 1 |
| `X33` | [Synchroscope](<gauges/x33_synchroscope.md>) | Frequency difference and phase angle between AC sources | c.1900-present | 3 | 2 | 1 |
| `X36` | [Mechanical G-meter with min/max witness needles](<gauges/x36_mechanical_g_meter_with_min_max_witness_needles.md>) | Normal acceleration / load factor | 1930s-present | 6 | 1 | 1 |

## Alternate members

These gauges have a different primary group but also use this action or display form. They are not included in this group’s counts.

| ID | Gauge | Primary group |
|---|---|---|
| `P02` | [Mercury barometer: cistern, Fortin and wheel forms](<../liquid_column/gauges/p02_mercury_barometer_cistern_fortin_and_wheel_forms.md>) | [Liquid column](<../liquid_column/README.md>) |
| `P18` | [Float-and-tape tank gauge](<../moving_tape_ribbon_or_map/gauges/p18_float_and_tape_tank_gauge.md>) | [Moving tape, ribbon or map](<../moving_tape_ribbon_or_map/README.md>) |
| `P41` | [Mechanical odometer, cyclometer and utility-meter register](<../rolling_drum_or_counter/gauges/p41_mechanical_odometer_cyclometer_and_utility_meter_register.md>) | [Rolling drum or counter](<../rolling_drum_or_counter/README.md>) |
| `X09` | [RMI and HSI multi-pointer navigation indicator](<../rotating_scale_or_scene/gauges/x09_rmi_and_hsi_multi_pointer_navigation_indicator.md>) | [Rotating scale or scene](<../rotating_scale_or_scene/README.md>) |
| `X32` | [Railway speed indicator and recorder](<../chart_or_trace_recorder/gauges/x32_railway_speed_indicator_and_recorder.md>) | [Chart or trace recorder](<../chart_or_trace_recorder/README.md>) |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Response speed, lag and delay](<quirks/response_speed_lag_and_delay.md>) | Finite response time between a change at the input and the corresponding visible indication. | 19 | 33.93% | 21 |
| [Scale linearity and nonlinearity](<quirks/scale_linearity_and_nonlinearity.md>) | Variation in how equal input increments map to equal or unequal distances on the displayed scale. | 15 | 26.79% | 16 |
| [Thermal behaviour and temperature effects](<quirks/thermal_behaviour_and_temperature_effects.md>) | Changes in reading, appearance, sensitivity or dynamics caused by instrument or ambient temperature. | 14 | 25.00% | 19 |
| [Flutter, jitter, tremor and quiver](<quirks/flutter_jitter_tremor_and_quiver.md>) | Small rapid random or periodic movements around the nominal indication. | 13 | 23.21% | 13 |
| [Hysteresis](<quirks/hysteresis.md>) | Different indicated values for the same input depending on whether the input approached from above or below. | 13 | 23.21% | 13 |
| [Overload, saturation and damage](<quirks/overload_saturation_and_damage.md>) | Behaviour when the input exceeds the useful range, including pegging, clipping, recovery changes or permanent harm. | 12 | 21.43% | 12 |
| [Friction and drag](<quirks/friction_and_drag.md>) | Motion resistance that slows, biases or distorts the indication. | 11 | 19.64% | 11 |
| [Damping](<quirks/damping.md>) | Deliberate or inherent suppression of rapid movement and oscillation in the indication. | 10 | 17.86% | 10 |
| [Drift and long-term stability](<quirks/drift_and_long_term_stability.md>) | Slow change in indicated value or behaviour over time despite an unchanged input. | 10 | 17.86% | 10 |
| [Shock and vibration effects](<quirks/shock_and_vibration_effects.md>) | Temporary or permanent indication changes caused by mechanical shock or sustained vibration. | 10 | 17.86% | 10 |
| [Zero drift and offset](<quirks/zero_drift_and_offset.md>) | A non-zero indication at the true zero point, including offsets that change with time or conditions. | 10 | 17.86% | 11 |
| [Compressed or expanded scale](<quirks/compressed_or_expanded_scale.md>) | A deliberately or accidentally nonuniform scale that gives some ranges more display space than others. | 8 | 14.29% | 8 |
| [Pressure pulsation and pneumatic-line dynamics](<quirks/pressure_pulsation_and_pneumatic_line_dynamics.md>) | Indication effects caused by pressure waves, tubing volume, restrictions and compressible-fluid behaviour. | 8 | 14.29% | 9 |
| [Scale markings, zones and legends](<quirks/scale_markings_zones_and_legends.md>) | Visual information carried by ticks, numerals, colour bands, labels and operating zones. | 8 | 14.29% | 9 |
| [Stiction and sticking](<quirks/stiction_and_sticking.md>) | Static friction or adhesion that prevents motion until enough force accumulates, often followed by a jump. | 8 | 14.29% | 8 |
| [Zero adjustment and checking](<quirks/zero_adjustment_and_checking.md>) | Procedures or controls used to establish, verify or restore the zero reference. | 8 | 14.29% | 8 |
| [Calibration, correction and compensation](<quirks/calibration_correction_and_compensation.md>) | Adjustments or correction factors required to relate the raw indication to the intended quantity. | 7 | 12.50% | 7 |
| [Creep](<quirks/creep.md>) | Very slow movement under a constant input, load or retained state. | 7 | 12.50% | 7 |
| [End stops, pegging and overflow](<quirks/end_stops_pegging_and_overflow.md>) | Behaviour at physical or representational range limits, including contact with stops and beyond-range indications. | 7 | 12.50% | 7 |
| [Settling and return behaviour](<quirks/settling_and_return_behaviour.md>) | How the indication approaches a stable value or returns after a transient, disturbance or release. | 7 | 12.50% | 7 |
| [Wraparound and multi-turn indication](<quirks/wraparound_and_multi_turn_indication.md>) | Behaviour when a value crosses a cyclic boundary or requires multiple revolutions to represent its full range. | 7 | 12.50% | 8 |
| [Averaging and integration](<quirks/averaging_and_integration.md>) | Intentional or inherent smoothing that indicates an average, accumulated or integrated quantity rather than instantaneous input. | 6 | 10.71% | 6 |
| [Contamination, dirt and fouling](<quirks/contamination_dirt_and_fouling.md>) | Reading or appearance changes caused by deposits, dust, oxidation, residue or biological growth. | 6 | 10.71% | 6 |
| [Human-factor ambiguity and misreading](<quirks/human_factor_ambiguity_and_misreading.md>) | Display features that make an otherwise functioning instrument easy to interpret incorrectly. | 6 | 10.71% | 8 |
| [Multi-pointer and multi-channel interaction](<quirks/multi_pointer_and_multi_channel_interaction.md>) | Visual or mechanical interaction among multiple pointers, scales or measurement channels. | 6 | 10.71% | 11 |
| [Orientation sensitivity](<quirks/orientation_sensitivity.md>) | Dependence on instrument mounting angle or orientation relative to gravity or fields. | 6 | 10.71% | 6 |
| [Backlash and lash](<quirks/backlash_and_lash.md>) | Lost motion when direction reverses because of clearance in gears, linkages or drive parts. | 5 | 8.93% | 5 |
| [Channel mismatch and unequal dynamics](<quirks/channel_mismatch_and_unequal_dynamics.md>) | Different calibration, response or motion between channels intended to behave alike. | 5 | 8.93% | 7 |
| [Chatter](<quirks/chatter.md>) | Rapid repeated switching or movement near a threshold or unstable equilibrium. | 5 | 8.93% | 5 |
| [Electrical or mechanical remote coupling](<quirks/electrical_or_mechanical_remote_coupling.md>) | Errors and dynamics introduced when a sensor transmits position or value to a remote indicator. | 5 | 8.93% | 5 |
| [Frequency, waveform and source dependence](<quirks/frequency_waveform_and_source_dependence.md>) | Changes in reading or response caused by signal frequency, waveform shape or source characteristics. | 5 | 8.93% | 9 |
| [Logarithmic scale](<quirks/logarithmic_scale.md>) | A scale where equal distances represent equal ratios rather than equal numeric increments. | 5 | 8.93% | 5 |
| [Power-off and power-loss behaviour](<quirks/power_off_and_power_loss_behaviour.md>) | What the indication does when drive power is removed, including retained, blank, parked or misleading states. | 5 | 8.93% | 6 |
| [Thresholds, deadband and switching points](<quirks/thresholds_deadband_and_switching_points.md>) | Regions or levels where no change occurs, or where a discrete state changes with defined or variable thresholds. | 5 | 8.93% | 6 |
| [Witness, peak-hold and retained extrema](<quirks/witness_peak_hold_and_retained_extrema.md>) | Markers or memory mechanisms that preserve minimum, maximum or peak values after the live indication moves away. | 5 | 8.93% | 8 |
| [Bounce](<quirks/bounce.md>) | Repeated rebounds or reversals after a mechanical impact, contact change or rapid movement. | 4 | 7.14% | 4 |
| [Bubble, void and separated-column behaviour](<quirks/bubble_void_and_separated_column_behaviour.md>) | Errors or discontinuities caused by trapped gas, empty spaces or broken fluid columns. | 4 | 7.14% | 5 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 4 | 7.14% | 4 |
| [Electrical contacts, grounds and wiring](<quirks/electrical_contacts_grounds_and_wiring.md>) | Display faults caused by contact resistance, grounding, broken conductors, polarity or connection layout. | 4 | 7.14% | 5 |
| [Gas or medium dependence](<quirks/gas_or_medium_dependence.md>) | Dependence on the type, pressure, composition or condition of the working gas or surrounding medium. | 4 | 7.14% | 4 |
| [Hunting and following error](<quirks/hunting_and_following_error.md>) | Repeated corrective motion or persistent difference between commanded and indicated position in a following system. | 4 | 7.14% | 5 |
| [Invalid, out-of-range and warning flags](<quirks/invalid_out_of_range_and_warning_flags.md>) | Explicit indications that a reading is unavailable, unreliable, unsafe or beyond the valid range. | 4 | 7.14% | 5 |
| [Magnetic-field, deviation and remanence effects](<quirks/magnetic_field_deviation_and_remanence_effects.md>) | Influence of external or retained magnetism on indication, zero and calibration. | 4 | 7.14% | 4 |
| [Mechanical noise and cadence](<quirks/mechanical_noise_and_cadence.md>) | Audible clicks, hums, impacts or rhythms produced by the display mechanism. | 4 | 7.14% | 6 |
| [Mechanical wear, stretch and permanent set](<quirks/mechanical_wear_stretch_and_permanent_set.md>) | Long-term dimensional or elastic change in moving parts, springs, fibres or linkages. | 4 | 7.14% | 4 |
| [Movement torque and sensitivity](<quirks/movement_torque_and_sensitivity.md>) | The relationship between applied drive, restoring force and resulting visible movement. | 4 | 7.14% | 4 |
| [Overshoot](<quirks/overshoot.md>) | Temporary travel beyond the final steady indication after an input change. | 4 | 7.14% | 4 |
| [Reference bugs and set markers](<quirks/reference_bugs_and_set_markers.md>) | Movable or fixed markers used to record targets, limits, headings or comparison values. | 4 | 7.14% | 4 |
| [Shadows, depth and occlusion](<quirks/shadows_depth_and_occlusion.md>) | Visual effects caused by layered parts blocking, shading or appearing at different depths. | 4 | 7.14% | 6 |
| [Cool-down, heat soak and continued motion](<quirks/cool_down_heat_soak_and_continued_motion.md>) | Behaviour after drive or heat is removed while stored thermal or mechanical energy remains. | 3 | 5.36% | 4 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 3 | 5.36% | 3 |
| [Fluid leakage](<quirks/fluid_leakage.md>) | Escape of the working fluid, causing bias, loss of range, contamination or failure. | 3 | 5.36% | 3 |
| [Manual adjustment and setup](<quirks/manual_adjustment_and_setup.md>) | User controls and preparation steps needed to configure the instrument before use. | 3 | 5.36% | 3 |
| [Manual reset, tare and caging](<quirks/manual_reset_tare_and_caging.md>) | User-operated mechanisms that clear, zero, restrain or protect the indication. | 3 | 5.36% | 3 |
| [Parallax](<quirks/parallax.md>) | Apparent reading error caused by viewing the indicator and scale from the wrong angle or depth relationship. | 3 | 5.36% | 3 |
| [Revolution counting and coarse/fine readout](<quirks/revolution_counting_and_coarse_fine_readout.md>) | Methods that combine turns, counters or multiple indicators to represent both large range and fine resolution. | 3 | 5.36% | 3 |
| [Valid operating region](<quirks/valid_operating_region.md>) | The part of the scale or operating envelope in which the indication is considered reliable. | 3 | 5.36% | 3 |
| [Warm-up behaviour](<quirks/warm_up_behaviour.md>) | Transient changes after startup while temperature, discharge, illumination or mechanics stabilise. | 3 | 5.36% | 3 |
| [Ageing and material degradation](<quirks/ageing_and_material_degradation.md>) | Progressive change due to wear, fatigue, chemical decay, phosphor loss, embrittlement or similar ageing processes. | 2 | 3.57% | 2 |
| [Ambient-pressure and altitude sensitivity](<quirks/ambient_pressure_and_altitude_sensitivity.md>) | Changes in indication caused by surrounding atmospheric pressure or elevation. | 2 | 3.57% | 2 |
| [Attack and release ballistics](<quirks/attack_and_release_ballistics.md>) | Different rise and fall time constants used to emphasise peaks or produce a specified dynamic response. | 2 | 3.57% | 2 |
| [Balance, imbalance and pointer balance](<quirks/balance_imbalance_and_pointer_balance.md>) | Errors or dynamics caused by unequal mass distribution around a pivot or rotating assembly. | 2 | 3.57% | 2 |
| [Capillary, wetting and surface-tension effects](<quirks/capillary_wetting_and_surface_tension_effects.md>) | Fluid effects arising from narrow passages, adhesion and surface forces at boundaries. | 2 | 3.57% | 2 |
| [Carry, rollover and digit transition](<quirks/carry_rollover_and_digit_transition.md>) | The mechanical or visual sequence as counters advance between digits and propagate carries. | 2 | 3.57% | 3 |
| [Centre-zero, bidirectional and asymmetric range](<quirks/centre_zero_bidirectional_and_asymmetric_range.md>) | Scales that represent positive and negative values around a centre, possibly with unequal ranges or behaviour. | 2 | 3.57% | 4 |
| [Colour variation and colour shift](<quirks/colour_variation_and_colour_shift.md>) | Changes or inconsistencies in displayed colour across level, age, temperature, angle or individual units. | 2 | 3.57% | 3 |
| [Construction tolerances and unit variation](<quirks/construction_tolerances_and_unit_variation.md>) | Differences between nominally identical instruments caused by manufacturing and assembly tolerances. | 2 | 3.57% | 4 |
| [Cross-axis and acceleration sensitivity](<quirks/cross_axis_and_acceleration_sensitivity.md>) | Response to acceleration or forces along axes other than the intended measurement axis. | 2 | 3.57% | 2 |
| [Dead time, missed events and duplicate counts](<quirks/dead_time_missed_events_and_duplicate_counts.md>) | Intervals or mechanisms that cause events to be ignored, delayed or counted more than once. | 2 | 3.57% | 2 |
| [Drive power and regulation](<quirks/drive_power_and_regulation.md>) | Dependence on the available drive energy and the quality of its regulation. | 2 | 3.57% | 2 |
| [Electrical hum, buzz and whine](<quirks/electrical_hum_buzz_and_whine.md>) | Audible vibration or tone produced by electrical drive, magnetic parts or switching. | 2 | 3.57% | 2 |
| [Filament behaviour and failure](<quirks/filament_behaviour_and_failure.md>) | Warm-up, sag, resistance change, brightness variation and breakage of incandescent filaments. | 2 | 3.57% | 3 |
| [Fluid surge, slosh and foam](<quirks/fluid_surge_slosh_and_foam.md>) | Transient level or pressure effects caused by fluid motion, aeration or froth. | 2 | 3.57% | 2 |
| [Homing, parking and startup sweep](<quirks/homing_parking_and_startup_sweep.md>) | Controlled movement to a reference, rest position or full-scale test during startup or shutdown. | 2 | 3.57% | 3 |
| [Immersion, stem-conduction and contact error](<quirks/immersion_stem_conduction_and_contact_error.md>) | Temperature-reading errors caused by installation depth, heat flow along a probe or imperfect thermal contact. | 2 | 3.57% | 2 |
| [Inertia and coasting](<quirks/inertia_and_coasting.md>) | Continued or delayed motion due to mass and stored kinetic energy. | 2 | 3.57% | 2 |
| [Input loading and ratio dependence](<quirks/input_loading_and_ratio_dependence.md>) | Reading changes caused by the instrument altering the measured system or depending on a ratio of inputs. | 2 | 3.57% | 2 |
| [Latching and state-retention behaviour](<quirks/latching_and_state_retention_behaviour.md>) | A displayed state that remains mechanically, magnetically, electrically or optically retained until reset or rewritten. | 2 | 3.57% | 2 |
| [Meniscus behaviour](<quirks/meniscus_behaviour.md>) | The curved liquid surface and its reading conventions, wetting shape and movement. | 2 | 3.57% | 2 |
| [Operator procedure and ritual](<quirks/operator_procedure_and_ritual.md>) | Required handling or reading practices that materially affect the result. | 2 | 3.57% | 2 |
| [Pen, stylus and trace artefacts](<quirks/pen_stylus_and_trace_artefacts.md>) | Line-width, drag, skipping, smear, lift-off or other defects introduced by a recording or tracing element. | 2 | 3.57% | 2 |
| [Power-up and self-test behaviour](<quirks/power_up_and_self_test_behaviour.md>) | Visible startup sequences, checks, sweeps or initial states used when power is applied. | 2 | 3.57% | 2 |
| [Range switching and multiple ranges](<quirks/range_switching_and_multiple_ranges.md>) | Selection or combination of alternate scales, sensitivities or measuring ranges. | 2 | 3.57% | 3 |
| [Rate limiting and motion limits](<quirks/rate_limiting_and_motion_limits.md>) | Restrictions on how quickly or how far an indication may move. | 2 | 3.57% | 3 |
| [Readout masking and narrow-window transitions](<quirks/readout_masking_and_narrow_window_transitions.md>) | Partial or ambiguous readings caused by apertures that reveal only a small portion of a moving scale or drum. | 2 | 3.57% | 2 |
| [Ringing and oscillation](<quirks/ringing_and_oscillation.md>) | Repeated decaying or sustained motion around a target following excitation or disturbance. | 2 | 3.57% | 2 |
| [Scale direction and interpretation](<quirks/scale_direction_and_interpretation.md>) | Whether increasing values move clockwise, anticlockwise, upward or otherwise, and how that direction is understood. | 2 | 3.57% | 2 |
| [Sensor plumbing faults](<quirks/sensor_plumbing_faults.md>) | Errors caused by incorrectly connected, leaking, reversed or contaminated pressure and fluid lines. | 2 | 3.57% | 2 |
| [Snap action](<quirks/snap_action.md>) | A rapid transition between stable positions once a threshold is crossed. | 2 | 3.57% | 2 |
| [Supply-voltage sensitivity](<quirks/supply_voltage_sensitivity.md>) | Changes in reading, brightness or dynamics caused by variations in supply voltage. | 2 | 3.57% | 3 |
| [Blockage and restricted passages](<quirks/blockage_and_restricted_passages.md>) | Errors or slow response caused by obstructed tubing, ports, capillaries or vents. | 1 | 1.79% | 1 |
| [Cable, belt, tape and roller errors](<quirks/cable_belt_tape_and_roller_errors.md>) | Slip, stretch, tracking, tension or geometry errors in flexible mechanical transmission and transport parts. | 1 | 1.79% | 1 |
| [Case, lens and enclosure deformation](<quirks/case_lens_and_enclosure_deformation.md>) | Reading or visual changes caused by pressure, heat, stress or damage deforming the enclosure or window. | 1 | 1.79% | 1 |
| [Chart and paper transport](<quirks/chart_and_paper_transport.md>) | Timing, speed, alignment and mechanical behaviour of the recording medium. | 1 | 1.79% | 1 |
| [Clock drift, timing and seams](<quirks/clock_drift_timing_and_seams.md>) | Errors caused by imperfect timebase speed, chart joins, scan boundaries or repeating mechanical cycles. | 1 | 1.79% | 1 |
| [Detents and mechanical indexing](<quirks/detents_and_mechanical_indexing.md>) | Discrete stable positions established by notches, pawls, stops or indexing mechanisms. | 1 | 1.79% | 1 |
| [Electrostatic leakage and charge retention](<quirks/electrostatic_leakage_and_charge_retention.md>) | Loss or persistence of electric charge affecting electrostatic instruments and displays. | 1 | 1.79% | 1 |
| [Gas-discharge instability and dropout](<quirks/gas_discharge_instability_and_dropout.md>) | Flicker, extinction or irregular conduction after a gas discharge has started. | 1 | 1.79% | 2 |
| [Gas-discharge striking and ignition](<quirks/gas_discharge_striking_and_ignition.md>) | The voltage, timing and transient behaviour required to initiate a gas discharge. | 1 | 1.79% | 2 |
| [Humidity and moisture sensitivity](<quirks/humidity_and_moisture_sensitivity.md>) | Changes caused by water vapour, condensation, absorption or damp contamination. | 1 | 1.79% | 1 |
| [Illumination and backlighting](<quirks/illumination_and_backlighting.md>) | Lighting systems that make scales, legends or display elements visible and their associated unevenness or ageing. | 1 | 1.79% | 1 |
| [Operator signalling and acknowledgement](<quirks/operator_signalling_and_acknowledgement.md>) | Human commands, confirmations or acknowledgements that form part of the indication system. | 1 | 1.79% | 2 |
| [Phase, synchronism and rotation direction](<quirks/phase_synchronism_and_rotation_direction.md>) | Dependence on relative phase, locked timing or the direction of a rotating field or mechanism. | 1 | 1.79% | 4 |
| [Qualitative or non-precision indication](<quirks/qualitative_or_non_precision_indication.md>) | An indication intended for trend, state or rough comparison rather than accurate numeric measurement. | 1 | 1.79% | 1 |
| [Random, statistical and batch variation](<quirks/random_statistical_and_batch_variation.md>) | Unpredictable events or unit-to-unit differences that are part of the observed behaviour. | 1 | 1.79% | 1 |
| [Reversal and direction-change error](<quirks/reversal_and_direction_change_error.md>) | Error or transient behaviour specifically introduced when motion or input direction reverses. | 1 | 1.79% | 1 |
| [Safety, guarding and fail-safe behaviour](<quirks/safety_guarding_and_fail_safe_behaviour.md>) | Features or failure modes intended to protect the user, equipment or validity of the indication. | 1 | 1.79% | 1 |
| [Tapping](<quirks/tapping.md>) | Intentional tapping used to release friction, settle the mechanism or obtain a representative reading. | 1 | 1.79% | 1 |
| [Viewing-angle dependence](<quirks/viewing_angle_dependence.md>) | Changes in readability, colour, contrast or apparent value with observer angle. | 1 | 1.79% | 1 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
