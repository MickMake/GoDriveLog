---
gauge_group: mechanical_flag_or_shutter
catalogue_version: "0.2"
primary_gauge_count: 3
supporting_quirk_count: 17
---

# Mechanical flag, shutter or semaphore

> The JSON files are authoritative. This README is the human-readable group view.

## Definition

A physical flag, shutter, ball, striped drum, vane or semaphore changes position to reveal a state or warning.

**Catalogue definition:** A flag, shutter, vane, striped drum or semaphore changes between visible states.

## How the group encodes a value

Information is encoded by the visible face, colour, pattern or position of a discrete mechanical element.

## Classification boundary

Use this group for a small number of stateful mechanical indicators. Pixel-like arrays of flipping elements belong under flip_element_array.

## Simulation baseline

Model latching, snap action, partial travel, bounce, mechanical cadence, reset method and the power-loss state.

## Catalogue coverage

| Metric | Value |
|---|---:|
| Primary gauges | 3 |
| Share of catalogue | 2.21% |
| Alternate members | 0 |
| Canonical quirks represented | 17 |

Primary counts use the mutually exclusive `gauge_group` field. Alternate memberships are cross-references only.

## Example gauges

| ID | Gauge | Representative names or models | Measured or indicated | Era |
|---|---|---|---|---|
| `E18` | [Magnetic latching flag, ball and shutter annunciator](<gauges/e18_magnetic_latching_flag_ball_and_shutter_annunciator.md>) | Aircraft OFF flags; red/white magnetic balls; relay semaphore indicators | Binary or small-state warning/status indication | 1920s-present |
| `E26` | [Barber-pole warning indicator](<gauges/e26_barber_pole_warning_indicator.md>) | Striped rotating drum or sliding flag for overspeed, unsafe gear, cabin altitude and control limits | Binary warning or moving limit boundary | 1930s-present |
| `X23` | [Railway block, semaphore and route repeater](<gauges/x23_railway_block_semaphore_and_route_repeater.md>) | Miniature signal-arm repeaters; block telegraph needles; route-indicator shutters | Track occupancy, signal state, route or permission | 1850s-present in heritage and legacy signalling |

## All primary gauges

| ID | Gauge | Measured or indicated | Era | Quirks | Sources | Images |
|---|---|---|---|---:|---:|---:|
| `E18` | [Magnetic latching flag, ball and shutter annunciator](<gauges/e18_magnetic_latching_flag_ball_and_shutter_annunciator.md>) | Binary or small-state warning/status indication | 1920s-present | 8 | 1 | 1 |
| `E26` | [Barber-pole warning indicator](<gauges/e26_barber_pole_warning_indicator.md>) | Binary warning or moving limit boundary | 1930s-present | 5 | 1 | 1 |
| `X23` | [Railway block, semaphore and route repeater](<gauges/x23_railway_block_semaphore_and_route_repeater.md>) | Track occupancy, signal state, route or permission | 1850s-present in heritage and legacy signalling | 9 | 2 | 1 |

## Supporting quirks

| Quirk | Definition | Gauges in group | Group share | Source statements |
|---|---|---:|---:|---:|
| [Snap action](<quirks/snap_action.md>) | A rapid transition between stable positions once a threshold is crossed. | 3 | 100.00% | 3 |
| [Bounce](<quirks/bounce.md>) | Repeated rebounds or reversals after a mechanical impact, contact change or rapid movement. | 2 | 66.67% | 2 |
| [Contamination, dirt and fouling](<quirks/contamination_dirt_and_fouling.md>) | Reading or appearance changes caused by deposits, dust, oxidation, residue or biological growth. | 2 | 66.67% | 2 |
| [Display geometry and motion mode](<quirks/display_geometry_and_motion_mode.md>) | The physical path, layout, orientation or mode by which the visible indication moves or changes. | 2 | 66.67% | 3 |
| [Chatter](<quirks/chatter.md>) | Rapid repeated switching or movement near a threshold or unstable equilibrium. | 1 | 33.33% | 1 |
| [Discrete, quantised or stepwise motion](<quirks/discrete_quantised_or_stepwise_motion.md>) | Indication that changes in finite increments rather than continuously. | 1 | 33.33% | 1 |
| [Electrical contacts, grounds and wiring](<quirks/electrical_contacts_grounds_and_wiring.md>) | Display faults caused by contact resistance, grounding, broken conductors, polarity or connection layout. | 1 | 33.33% | 1 |
| [Gravity-related behaviour](<quirks/gravity_related_behaviour.md>) | Dependence on local gravity magnitude or direction. | 1 | 33.33% | 1 |
| [Invalid, out-of-range and warning flags](<quirks/invalid_out_of_range_and_warning_flags.md>) | Explicit indications that a reading is unavailable, unreliable, unsafe or beyond the valid range. | 1 | 33.33% | 1 |
| [Latching and state-retention behaviour](<quirks/latching_and_state_retention_behaviour.md>) | A displayed state that remains mechanically, magnetically, electrically or optically retained until reset or rewritten. | 1 | 33.33% | 1 |
| [Mechanical noise and cadence](<quirks/mechanical_noise_and_cadence.md>) | Audible clicks, hums, impacts or rhythms produced by the display mechanism. | 1 | 33.33% | 1 |
| [Power-off and power-loss behaviour](<quirks/power_off_and_power_loss_behaviour.md>) | What the indication does when drive power is removed, including retained, blank, parked or misleading states. | 1 | 33.33% | 1 |
| [Safety, guarding and fail-safe behaviour](<quirks/safety_guarding_and_fail_safe_behaviour.md>) | Features or failure modes intended to protect the user, equipment or validity of the indication. | 1 | 33.33% | 1 |
| [Scale markings, zones and legends](<quirks/scale_markings_zones_and_legends.md>) | Visual information carried by ticks, numerals, colour bands, labels and operating zones. | 1 | 33.33% | 2 |
| [Segment, pixel, lamp or flag failure](<quirks/segment_pixel_lamp_or_flag_failure.md>) | Individual display elements that fail open, fail active, stick, weaken or respond intermittently. | 1 | 33.33% | 1 |
| [Settling and return behaviour](<quirks/settling_and_return_behaviour.md>) | How the indication approaches a stable value or returns after a transient, disturbance or release. | 1 | 33.33% | 1 |
| [Thresholds, deadband and switching points](<quirks/thresholds_deadband_and_switching_points.md>) | Regions or levels where no change occurs, or where a discrete state changes with defined or variable thresholds. | 1 | 33.33% | 2 |

## Reference files

- [gauge_display_research_catalog_v0.2.json](<../_data/gauge_display_research_catalog_v0.2.json>)
- [gauge_group_index_v0.2.json](<../_data/gauge_group_index_v0.2.json>)
- [gauge_group_quirk_index_v0.2.json](<../_data/gauge_group_quirk_index_v0.2.json>)

[Back to all gauge groups](../GAUGE_GROUPS.md)
