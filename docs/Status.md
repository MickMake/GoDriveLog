# GoDriveLog Documentation Status

Pillar 2 - The current project truth.

This file records the current state of documented GoDriveLog features against current code.
It exists as a map of where the project is now, so the current state can be recovered quickly without replaying old release notes, prompts, branches, or review history.

This file is not a release gate and not a wishlist. Planned work may appear while a release is active. As work progresses, each row should be updated to reflect the current truth.

The `Version` column records the release or slice version the item belongs to. It does not always mean the item was successfully implemented in that version.

## Status values

| Status | Meaning |
|—|—|
| Planned | Accepted into a release or slice plan, but implementation has not landed yet. |
| In progress | Work has started but is not complete or not yet audited. |
| Implemented | Current code supports the documented behaviour in scope. |
| Partially implemented | Some current code support exists, but the implemented behaviour is incomplete, compatibility-only, or differs from the documented scope. |
| Not implemented | Design or slice exists, but current code does not support it. |
| Deferred | Deliberately moved out of the current release or active scope. |
| Superseded | Replaced by another design, config key, implementation path, or naming model. |
| Unable to verify | Current evidence is insufficient to verify the feature end-to-end. |

## Gauge status summary

| Status | Count |
|—|—:|
| Implemented | 24 |
| Partially implemented | 4 |
| Planned | 7 |
| In progress | 0 |
| Not implemented | 0 |
| Deferred | 0 |
| Superseded | 0 |
| Unable to verify | 0 |

## Gauge implementation status

### `radial_pointer`

| Name | Status | Version | Current config key | Quirk/Gauge doc | Notes |
|—|—|—|—|—|—|
| `custom_radial` | Implemented | v3.4 | `type: radial` | [Design](Designs/Gauge/radial_pointer/gauges/custom_radial.md) / [Implementation](Implementation/Gauge/radial_pointer/gauges/custom_radial.md) | Current GoDriveLog radial renderer, mapped to the mechanism-first `radial_pointer` group. |
| `damping` | Implemented | v3.5.8 | `realism.damping` | [Design](Designs/Gauge/radial_pointer/quirks/custom_damping.md) / [Implementation](Implementation/Gauge/radial_pointer/quirks/custom_damping.md) | Radial value-change smoothing. |
| `stiction` | Implemented | v3.5.9 | `realism.stiction` | [Design](Designs/Gauge/radial_pointer/quirks/custom_stiction.md) / [Implementation](Implementation/Gauge/radial_pointer/quirks/custom_stiction.md) | Holds small display changes until the configured threshold is exceeded. |
| `overshoot` | Implemented | v3.5.10 | `realism.overshoot` | [Design](Designs/Gauge/radial_pointer/quirks/custom_overshoot.md) / [Implementation](Implementation/Gauge/radial_pointer/quirks/custom_overshoot.md) | Finite overshoot and settle behaviour for radial pointer motion. |
| `peg_bounce` | Implemented | v3.5.11 | `realism.peg_bounce` | [Design](Designs/Gauge/radial_pointer/quirks/custom_peg_bounce.md) / [Implementation](Implementation/Gauge/radial_pointer/quirks/custom_peg_bounce.md) | End-stop bounce when the displayed pointer hits a configured range limit. |
| `hysteresis` | Implemented | v3.5.16 | `realism.hysteresis` | [Design](Designs/Gauge/radial_pointer/quirks/custom_hysteresis.md) / [Implementation](Implementation/Gauge/radial_pointer/quirks/custom_hysteresis.md) | Display-side hysteresis/deadband behaviour. |
| `needle_shadow` | Implemented | v3.5.17 | `realism.needle_shadow` | [Design](Designs/Gauge/radial_pointer/quirks/custom_needle_shadow.md) / [Implementation](Implementation/Gauge/radial_pointer/quirks/custom_needle_shadow.md) | Static renderer depth effect; not dynamic parallax or lighting. |
| `calibration_offset` | Implemented | v3.5.18 | `realism.calibration_offset` | [Design](Designs/Gauge/radial_pointer/quirks/custom_calibration_offset.md) / [Implementation](Implementation/Gauge/radial_pointer/quirks/custom_calibration_offset.md) | Display-only offset; does not change the source sensor value. |
| `pointer_markers` / witness markers | Implemented | v3.5 | `realism.pointer_markers` | [Design](Designs/Gauge/radial_pointer/quirks/custom_pointer_markers.md) / [Implementation](Implementation/Gauge/radial_pointer/quirks/custom_pointer_markers.md) | Current config key is `pointer_markers`; older notes may call these witness markers. |
| `movement_policy` | Partially implemented | v3.5.5 | `realism.movement_policy` | [Design](Designs/Gauge/radial_pointer/quirks/custom_movement_policy.md) / [Implementation](Implementation/Gauge/radial_pointer/quirks/custom_movement_policy.md) | Accepted and used as a runtime transition policy. It is not a standalone physical quirk; visible effect depends on another timed movement behaviour such as damping, overshoot, or peg bounce. |

### `bar_or_wedge_display`

| Name | Status | Version | Current config key | Quirk/Gauge doc | Notes |
|—|—|—|—|—|—|
| `custom_bar` | Implemented | v3.4 | `type: bar` | [Design](Designs/Gauge/bar_or_wedge_display/gauges/custom_bar.md) / [Implementation](Implementation/Gauge/bar_or_wedge_display/gauges/custom_bar.md) | Current GoDriveLog continuous bar renderer, mapped to the mechanism-first `bar_or_wedge_display` group. |
| `custom_segmented` | Implemented | v3.4 | `type: segmented` | [Design](Designs/Gauge/bar_or_wedge_display/gauges/custom_segmented.md) / [Implementation](Implementation/Gauge/bar_or_wedge_display/gauges/custom_segmented.md) | Current GoDriveLog sparse percent-threshold image-selection gauge; intentionally mapped here rather than `segmented_display`. |
| `damping` | Implemented | v3.5.13 | `realism.damping` | [Design](Designs/Gauge/bar_or_wedge_display/quirks/custom_damping.md) / [Implementation](Implementation/Gauge/bar_or_wedge_display/quirks/custom_damping.md) | Bar fill/reveal smoothing. |
| `overshoot` | Implemented | v3.5.19 | `realism.overshoot` | [Design](Designs/Gauge/bar_or_wedge_display/quirks/custom_overshoot.md) / [Implementation](Implementation/Gauge/bar_or_wedge_display/quirks/custom_overshoot.md) | Finite overshoot and settle behaviour for displayed fill/reveal extent. |
| `hysteresis` | Implemented | v3.5.20 | `realism.hysteresis` | [Design](Designs/Gauge/bar_or_wedge_display/quirks/custom_hysteresis.md) / [Implementation](Implementation/Gauge/bar_or_wedge_display/quirks/custom_hysteresis.md) | Display-side hysteresis/deadband behaviour for bar fill/reveal extent. |
| `stiction` | Implemented | v3.5.21 | `realism.stiction` | [Design](Designs/Gauge/bar_or_wedge_display/quirks/custom_stiction.md) / [Implementation](Implementation/Gauge/bar_or_wedge_display/quirks/custom_stiction.md) | Holds small display changes until the configured threshold is exceeded. |
| `peg_bounce` | Implemented | v3.5.22 | `realism.peg_bounce` | [Design](Designs/Gauge/bar_or_wedge_display/quirks/custom_peg_bounce.md) / [Implementation](Implementation/Gauge/bar_or_wedge_display/quirks/custom_peg_bounce.md) | End-stop bounce on displayed fill/reveal extent. |
| `pointer_markers` / witness markers | Implemented | v3.5 | `realism.pointer_markers` | [Design](Designs/Gauge/bar_or_wedge_display/quirks/custom_pointer_markers.md) / [Implementation](Implementation/Gauge/bar_or_wedge_display/quirks/custom_pointer_markers.md) | Current config key is `pointer_markers`; older notes may call these witness markers. |
| `movement_policy` | Partially implemented | v3.5.5 | `realism.movement_policy` | [Design](Designs/Gauge/bar_or_wedge_display/quirks/custom_movement_policy.md) / [Implementation](Implementation/Gauge/bar_or_wedge_display/quirks/custom_movement_policy.md) | Accepted and used as a runtime transition policy. It is not a standalone physical quirk; visible effect depends on another timed movement behaviour such as damping, overshoot, or peg bounce. |

### `rolling_drum_or_counter`

| Name | Status | Version | Current config key | Quirk/Gauge doc | Notes |
|—|—|—|—|—|—|
| `custom_odometer` | Implemented | v3.4 | `type: odometer` | [Design](Designs/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md) / [Implementation](Implementation/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md) | Current GoDriveLog odometer renderer, mapped to the mechanism-first `rolling_drum_or_counter` group. |
| `movement` | Partially implemented | v3.5.6b | `odometer.movement` | [Design](Designs/Gauge/rolling_drum_or_counter/quirks/custom_movement.md) / [Implementation](Implementation/Gauge/rolling_drum_or_counter/quirks/custom_movement.md) | `instant`, `linear`, `ease_out`, and `bell` have concrete behaviour. `smooth` and `click` are recognised but fall back to `instant`. |
| `wraparound` | Partially implemented | v3.5.2 | `realism.wraparound` | [Design](Designs/Gauge/rolling_drum_or_counter/quirks/custom_wraparound.md) / [Implementation](Implementation/Gauge/rolling_drum_or_counter/quirks/custom_wraparound.md) | Config is accepted for odometer gauges, but current scene rendering does not branch on it; current wheel circularity is effectively compatibility behaviour. |
| `drum_slop` | Implemented | v3.5.3 | `realism.drum_slop` | [Design](Designs/Gauge/rolling_drum_or_counter/quirks/custom_drum_slop.md) / [Implementation](Implementation/Gauge/rolling_drum_or_counter/quirks/custom_drum_slop.md) | Per-wheel visual offset/slop for odometer drums. |
| `carry_drag` | Implemented | v3.5.7 | `realism.carry_drag` | [Design](Designs/Gauge/rolling_drum_or_counter/quirks/custom_carry_drag.md) / [Implementation](Implementation/Gauge/rolling_drum_or_counter/quirks/custom_carry_drag.md) | Visual carry-drag behaviour during digit carry. |
| `snap_settle` | Implemented | v3.5.14 | `realism.snap_settle` | [Design](Designs/Gauge/rolling_drum_or_counter/quirks/custom_snap_settle.md) / [Implementation](Implementation/Gauge/rolling_drum_or_counter/quirks/custom_snap_settle.md) | Finite settle/snap behaviour after odometer wheel movement. |
| `backlash` | Planned | v3.7.1 | `realism.backlash` not accepted | [Design](Designs/RealismBehaviour/odometer-backlash.md) / [Implementation](Implementation/RealismBehaviour/odometer-backlash.md) | v3.5.15 was corrected as not implemented on `main`; v3.7.1 plans odometer backlash as future work. |

### `indicator_lamp`

| Name | Status | Version | Current config key | Quirk/Gauge doc | Notes |
|—|—|—|—|—|—|
| `custom_indicator` | Implemented | v3.4 | `type: indicator` | [Design](Designs/Gauge/indicator_lamp/gauges/custom_indicator.md) / [Implementation](Implementation/Gauge/indicator_lamp/gauges/custom_indicator.md) | Current GoDriveLog two-state indicator renderer, mapped to the mechanism-first `indicator_lamp` group. |
| `thermal_fade` | Implemented | v3.5.12 | `realism.thermal_fade` | [Design](Designs/Gauge/indicator_lamp/quirks/custom_thermal_fade.md) / [Implementation](Implementation/Gauge/indicator_lamp/quirks/custom_thermal_fade.md) | Indicator on/off alpha transition to simulate finite lamp response. |

### `segmented_display`

| Name | Status | Version | Current config key | Quirk/Gauge doc | Notes |
|—|—|—|—|—|—|
| `custom_numeric` | Implemented | v3.4 | `type: numeric` | [Design](Designs/Gauge/segmented_display/gauges/custom_numeric.md) / [Implementation](Implementation/Gauge/segmented_display/gauges/custom_numeric.md) | Current GoDriveLog formatted-value image-slot renderer; old `seven_segment` naming was replaced by `numeric`. |
| `per_digit_response_lag` | Planned | v3.7.2 | `realism.per_digit_response_lag` not accepted | [Design](Designs/RealismBehaviour/per-digit-response-lag.md) / [Implementation](Implementation/RealismBehaviour/per-digit-response-lag.md) | Planned v3.7 numeric/segmented realism slice; no current code support yet. |
| `leading_zero_behaviour` | Planned | v3.7.3 | `realism.leading_zero_behaviour` not accepted | [Design](Designs/RealismBehaviour/numeric-leading-zero-behaviour.md) / [Implementation](Implementation/RealismBehaviour/numeric-leading-zero-behaviour.md) | Planned v3.7 numeric display behaviour; no current code support yet. |
| `segment_bleed` / `digit_bleed` | Planned | v3.7.4 | `realism.segment_bleed` / `realism.digit_bleed` not accepted | [Design](Designs/RealismBehaviour/segment-bleed-digit-bleed.md) / [Implementation](Implementation/RealismBehaviour/segment-bleed-digit-bleed.md) | Planned v3.7 numeric/segmented display realism slice; no current code support yet. |
| `ghosting` | Planned | v3.7.5 | `realism.ghosting` not accepted | [Design](Designs/RealismBehaviour/numeric-ghosting.md) / [Implementation](Implementation/RealismBehaviour/numeric-ghosting.md) | Planned v3.7 numeric/segmented display realism slice; no current code support yet. |
| `uneven_brightness` | Planned | v3.7.6 | `realism.uneven_brightness` not accepted | [Design](Designs/RealismBehaviour/uneven-brightness.md) / [Implementation](Implementation/RealismBehaviour/uneven-brightness.md) | Planned v3.7 numeric/segmented display realism slice; no current code support yet. |
| `load_sag` | Planned | v3.7.7 | `realism.load_sag` not accepted | [Design](Designs/RealismBehaviour/numeric-load-sag.md) / [Implementation](Implementation/RealismBehaviour/numeric-load-sag.md) | Planned v3.7 numeric/segmented display realism slice; no current code support yet. |

## Release planning checkpoints

### v3.7

| Name | Status | Version | Current config key | Quirk/Gauge doc | Notes |
|—|—|—|—|—|—|
| v3.7 release planning docs | Planned | v3.7.0 | — | [Implementation state](v3.7/ImplementationState.md) | Current v3.7 state is not started. |
| Odometer backlash | Planned | v3.7.1 | `realism.backlash` not accepted | [Design](Designs/RealismBehaviour/odometer-backlash.md) / [Implementation](Implementation/RealismBehaviour/odometer-backlash.md) | Also listed under `rolling_drum_or_counter`. |
| Per-digit response lag | Planned | v3.7.2 | `realism.per_digit_response_lag` not accepted | [Design](Designs/RealismBehaviour/per-digit-response-lag.md) / [Implementation](Implementation/RealismBehaviour/per-digit-response-lag.md) | Also listed under `segmented_display`. |
| Leading-zero behaviour | Planned | v3.7.3 | `realism.leading_zero_behaviour` not accepted | [Design](Designs/RealismBehaviour/numeric-leading-zero-behaviour.md) / [Implementation](Implementation/RealismBehaviour/numeric-leading-zero-behaviour.md) | Also listed under `segmented_display`. |
| Segment and digit bleed | Planned | v3.7.4 | `realism.segment_bleed` / `realism.digit_bleed` not accepted | [Design](Designs/RealismBehaviour/segment-bleed-digit-bleed.md) / [Implementation](Implementation/RealismBehaviour/segment-bleed-digit-bleed.md) | Also listed under `segmented_display`. |
| Ghosting | Planned | v3.7.5 | `realism.ghosting` not accepted | [Design](Designs/RealismBehaviour/numeric-ghosting.md) / [Implementation](Implementation/RealismBehaviour/numeric-ghosting.md) | Also listed under `segmented_display`. |
| Uneven brightness | Planned | v3.7.6 | `realism.uneven_brightness` not accepted | [Design](Designs/RealismBehaviour/uneven-brightness.md) / [Implementation](Implementation/RealismBehaviour/uneven-brightness.md) | Also listed under `segmented_display`. |
| Load sag | Planned | v3.7.7 | `realism.load_sag` not accepted | [Design](Designs/RealismBehaviour/numeric-load-sag.md) / [Implementation](Implementation/RealismBehaviour/numeric-load-sag.md) | Also listed under `segmented_display`. |
| Tests, previews and docs checkpoint | Planned | v3.7.8 | — | [Implementation state](v3.7/ImplementationState.md) | Planned release checkpoint slice. |

## Notes for future updates

- This file should be updated whenever code truth changes.
- Do not mark a row `Implemented` unless current code supports the documented behaviour in scope.
- Parser support alone is not enough for `Implemented` when the design requires runtime/rendering behaviour.
- `Planned` rows are allowed while a release is active.
- At release close, planned rows should be updated to the current truth: `Implemented`, `Partially implemented`, `Not implemented`, `Deferred`, `Superseded`, or left `Planned` only if the release remains open.
- Do not add every historical idea here. Track accepted release/slice scope and current code truth.
