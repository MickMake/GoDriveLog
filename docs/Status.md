# GoDriveLog Documentation Status

Pillar 2 - The current project truth.

This file records the current state of documented GoDriveLog features against current code.
It exists as a map of where the project is now, so the current state can be recovered quickly without replaying old release notes, prompts, branches, or review history.

This file is not a release gate and not a wishlist. Planned work may appear while a release is active. As work progresses, each row should be updated to reflect the current truth.

The `Version` column records the release or slice version the item belongs to. It does not always mean the item was successfully implemented in that version.

## Status values

<table>
  <thead>
    <tr>
      <th>Status</th>
      <th>Meaning</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Planned</td>
      <td>Accepted into a release or slice plan, but implementation has not landed yet.</td>
    </tr>
    <tr>
      <td>In progress</td>
      <td>Work has started but is not complete or not yet audited.</td>
    </tr>
    <tr>
      <td>Implemented</td>
      <td>Current code supports the documented behaviour in scope.</td>
    </tr>
    <tr>
      <td>Partially implemented</td>
      <td>Some current code support exists, but the implemented behaviour is incomplete, compatibility-only, or differs from the documented scope.</td>
    </tr>
    <tr>
      <td>Not implemented</td>
      <td>Design or slice exists, but current code does not support it.</td>
    </tr>
    <tr>
      <td>Deferred</td>
      <td>Deliberately moved out of the current release or active scope.</td>
    </tr>
    <tr>
      <td>Superseded</td>
      <td>Replaced by another design, config key, implementation path, or naming model.</td>
    </tr>
    <tr>
      <td>Unable to verify</td>
      <td>Current evidence is insufficient to verify the feature end-to-end.</td>
    </tr>
  </tbody>
</table>

## Gauge implementation status

### `radial_pointer`

<table>
  <thead>
    <tr><th>Name</th><th>Gauge</th><th>Status</th><th>Version</th><th>Current config key</th><th>Quirk/Gauge doc</th><th>Notes</th></tr>
  </thead>
  <tbody>
    <tr><td><code>custom_radial</code></td><td><code>custom_radial</code></td><td>Implemented</td><td>v3.4</td><td><code>type: radial</code></td><td><a href="Designs/Gauge/radial_pointer/gauges/custom_radial.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/gauges/custom_radial.md">Implementation</a></td><td>Current GoDriveLog radial renderer, mapped to the mechanism-first <code>radial_pointer</code> group.</td></tr>
    <tr><td><code>damping</code></td><td><code>custom_radial</code></td><td>Implemented</td><td>v3.5.8</td><td><code>realism.damping</code></td><td><a href="Designs/Gauge/radial_pointer/quirks/custom_damping.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/quirks/custom_damping.md">Implementation</a></td><td>Radial value-change smoothing.</td></tr>
    <tr><td><code>stiction</code></td><td><code>custom_radial</code></td><td>Implemented</td><td>v3.5.9</td><td><code>realism.stiction</code></td><td><a href="Designs/Gauge/radial_pointer/quirks/custom_stiction.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/quirks/custom_stiction.md">Implementation</a></td><td>Holds small display changes until the configured threshold is exceeded.</td></tr>
    <tr><td><code>overshoot</code></td><td><code>custom_radial</code></td><td>Implemented</td><td>v3.5.10</td><td><code>realism.overshoot</code></td><td><a href="Designs/Gauge/radial_pointer/quirks/custom_overshoot.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/quirks/custom_overshoot.md">Implementation</a></td><td>Finite overshoot and settle behaviour for radial pointer motion.</td></tr>
    <tr><td><code>peg_bounce</code></td><td><code>custom_radial</code></td><td>Implemented</td><td>v3.5.11</td><td><code>realism.peg_bounce</code></td><td><a href="Designs/Gauge/radial_pointer/quirks/custom_peg_bounce.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/quirks/custom_peg_bounce.md">Implementation</a></td><td>End-stop bounce when the displayed pointer hits a configured range limit.</td></tr>
    <tr><td><code>hysteresis</code></td><td><code>custom_radial</code></td><td>Implemented</td><td>v3.5.16</td><td><code>realism.hysteresis</code></td><td><a href="Designs/Gauge/radial_pointer/quirks/custom_hysteresis.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/quirks/custom_hysteresis.md">Implementation</a></td><td>Display-side hysteresis/deadband behaviour.</td></tr>
    <tr><td><code>needle_shadow</code></td><td><code>custom_radial</code></td><td>Implemented</td><td>v3.5.17</td><td><code>realism.needle_shadow</code></td><td><a href="Designs/Gauge/radial_pointer/quirks/custom_needle_shadow.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/quirks/custom_needle_shadow.md">Implementation</a></td><td>Static renderer depth effect; not dynamic parallax or lighting.</td></tr>
    <tr><td><code>calibration_offset</code></td><td><code>custom_radial</code></td><td>Implemented</td><td>v3.5.18</td><td><code>realism.calibration_offset</code></td><td><a href="Designs/Gauge/radial_pointer/quirks/custom_calibration_offset.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/quirks/custom_calibration_offset.md">Implementation</a></td><td>Display-only offset; does not change the source sensor value.</td></tr>
    <tr><td><code>pointer_markers</code> (witness markers)</td><td><code>custom_radial</code></td><td>Implemented</td><td>v3.5</td><td><code>realism.pointer_markers</code></td><td><a href="Designs/Gauge/radial_pointer/quirks/custom_pointer_markers.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/quirks/custom_pointer_markers.md">Implementation</a></td><td>Current config key is <code>pointer_markers</code>; older notes may call these witness markers.</td></tr>
    <tr><td><code>movement_policy</code></td><td><code>custom_radial</code></td><td>Partially implemented</td><td>v3.5.5</td><td><code>realism.movement_policy</code></td><td><a href="Designs/Gauge/radial_pointer/quirks/custom_movement_policy.md">Design</a><br><a href="Implementation/Gauge/radial_pointer/quirks/custom_movement_policy.md">Implementation</a></td><td><code>immediate</code>, <code>linear</code>, and <code>ease_out</code> are accepted/currently implemented. <code>bell</code> is part of the desired radial contract but is not accepted or applied for radial <code>movement_policy</code> yet. Plain movement policy may still depend on another timed movement behaviour to produce visible travel.</td></tr>
  </tbody>
</table>

### `bar_or_wedge_display`

<table>
  <thead>
    <tr><th>Name</th><th>Gauge</th><th>Status</th><th>Version</th><th>Current config key</th><th>Quirk/Gauge doc</th><th>Notes</th></tr>
  </thead>
  <tbody>
    <tr><td><code>custom_bar</code></td><td><code>custom_bar</code></td><td>Implemented</td><td>v3.4</td><td><code>type: bar</code></td><td><a href="Designs/Gauge/bar_or_wedge_display/gauges/custom_bar.md">Design</a><br><a href="Implementation/Gauge/bar_or_wedge_display/gauges/custom_bar.md">Implementation</a></td><td>Current GoDriveLog continuous bar renderer, mapped to the mechanism-first <code>bar_or_wedge_display</code> group.</td></tr>
    <tr><td><code>custom_segmented</code></td><td><code>custom_segmented</code></td><td>Implemented</td><td>v3.4</td><td><code>type: segmented</code></td><td><a href="Designs/Gauge/bar_or_wedge_display/gauges/custom_segmented.md">Design</a><br><a href="Implementation/Gauge/bar_or_wedge_display/gauges/custom_segmented.md">Implementation</a></td><td>Current GoDriveLog sparse percent-threshold image-selection gauge; intentionally mapped here rather than <code>segmented_display</code>.</td></tr>
    <tr><td><code>damping</code></td><td><code>custom_bar</code></td><td>Implemented</td><td>v3.5.13</td><td><code>realism.damping</code></td><td><a href="Designs/Gauge/bar_or_wedge_display/quirks/custom_damping.md">Design</a><br><a href="Implementation/Gauge/bar_or_wedge_display/quirks/custom_damping.md">Implementation</a></td><td>Bar fill/reveal smoothing.</td></tr>
    <tr><td><code>overshoot</code></td><td><code>custom_bar</code></td><td>Implemented</td><td>v3.5.19</td><td><code>realism.overshoot</code></td><td><a href="Designs/Gauge/bar_or_wedge_display/quirks/custom_overshoot.md">Design</a><br><a href="Implementation/Gauge/bar_or_wedge_display/quirks/custom_overshoot.md">Implementation</a></td><td>Finite overshoot and settle behaviour for displayed fill/reveal extent.</td></tr>
    <tr><td><code>hysteresis</code></td><td><code>custom_bar</code></td><td>Implemented</td><td>v3.5.20</td><td><code>realism.hysteresis</code></td><td><a href="Designs/Gauge/bar_or_wedge_display/quirks/custom_hysteresis.md">Design</a><br><a href="Implementation/Gauge/bar_or_wedge_display/quirks/custom_hysteresis.md">Implementation</a></td><td>Display-side hysteresis/deadband behaviour for bar fill/reveal extent.</td></tr>
    <tr><td><code>stiction</code></td><td><code>custom_bar</code></td><td>Implemented</td><td>v3.5.21</td><td><code>realism.stiction</code></td><td><a href="Designs/Gauge/bar_or_wedge_display/quirks/custom_stiction.md">Design</a><br><a href="Implementation/Gauge/bar_or_wedge_display/quirks/custom_stiction.md">Implementation</a></td><td>Holds small display changes until the configured threshold is exceeded.</td></tr>
    <tr><td><code>peg_bounce</code></td><td><code>custom_bar</code></td><td>Implemented</td><td>v3.5.22</td><td><code>realism.peg_bounce</code></td><td><a href="Designs/Gauge/bar_or_wedge_display/quirks/custom_peg_bounce.md">Design</a><br><a href="Implementation/Gauge/bar_or_wedge_display/quirks/custom_peg_bounce.md">Implementation</a></td><td>End-stop bounce on displayed fill/reveal extent.</td></tr>
    <tr><td><code>pointer_markers</code> (witness markers)</td><td><code>custom_bar</code></td><td>Implemented</td><td>v3.5</td><td><code>realism.pointer_markers</code></td><td><a href="Designs/Gauge/bar_or_wedge_display/quirks/custom_pointer_markers.md">Design</a><br><a href="Implementation/Gauge/bar_or_wedge_display/quirks/custom_pointer_markers.md">Implementation</a></td><td>Current config key is <code>pointer_markers</code>; older notes may call these witness markers.</td></tr>
    <tr><td><code>movement_policy</code></td><td><code>custom_bar</code></td><td>Partially implemented</td><td>v3.5.5</td><td><code>realism.movement_policy</code></td><td><a href="Designs/Gauge/bar_or_wedge_display/quirks/custom_movement_policy.md">Design</a><br><a href="Implementation/Gauge/bar_or_wedge_display/quirks/custom_movement_policy.md">Implementation</a></td><td>Accepted and used as a runtime transition policy. It is not a standalone physical quirk; visible effect depends on another timed movement behaviour.</td></tr>
  </tbody>
</table>

### `rolling_drum_or_counter`

<table>
  <thead>
    <tr><th>Name</th><th>Gauge</th><th>Status</th><th>Version</th><th>Current config key</th><th>Quirk/Gauge doc</th><th>Notes</th></tr>
  </thead>
  <tbody>
    <tr><td><code>custom_odometer</code></td><td><code>custom_odometer</code></td><td>Implemented</td><td>v3.4</td><td><code>type: odometer</code></td><td><a href="Designs/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md">Design</a><br><a href="Implementation/Gauge/rolling_drum_or_counter/gauges/custom_odometer.md">Implementation</a></td><td>Current GoDriveLog odometer renderer, mapped to the mechanism-first <code>rolling_drum_or_counter</code> group.</td></tr>
    <tr><td><code>movement</code></td><td><code>custom_odometer</code></td><td>Partially implemented</td><td>v3.5.6b</td><td><code>odometer.movement</code></td><td><a href="Designs/Gauge/rolling_drum_or_counter/quirks/custom_movement.md">Design</a><br><a href="Implementation/Gauge/rolling_drum_or_counter/quirks/custom_movement.md">Implementation</a></td><td><code>instant</code>, <code>linear</code>, <code>ease_out</code>, and <code>bell</code> have concrete behaviour. <code>smooth</code> and <code>click</code> are recognised but fall back to <code>instant</code>.</td></tr>
    <tr><td><code>wraparound</code></td><td><code>custom_odometer</code></td><td>Partially implemented</td><td>v3.5.2</td><td><code>realism.wraparound</code></td><td><a href="Designs/Gauge/rolling_drum_or_counter/quirks/custom_wraparound.md">Design</a><br><a href="Implementation/Gauge/rolling_drum_or_counter/quirks/custom_wraparound.md">Implementation</a></td><td>Config is accepted for odometer gauges, but current scene rendering does not branch on it; current wheel circularity is effectively compatibility behaviour.</td></tr>
    <tr><td><code>drum_slop</code></td><td><code>custom_odometer</code></td><td>Implemented</td><td>v3.5.3</td><td><code>realism.drum_slop</code></td><td><a href="Designs/Gauge/rolling_drum_or_counter/quirks/custom_drum_slop.md">Design</a><br><a href="Implementation/Gauge/rolling_drum_or_counter/quirks/custom_drum_slop.md">Implementation</a></td><td>Per-wheel visual offset/slop for odometer drums.</td></tr>
    <tr><td><code>carry_drag</code></td><td><code>custom_odometer</code></td><td>Implemented</td><td>v3.5.7</td><td><code>realism.carry_drag</code></td><td><a href="Designs/Gauge/rolling_drum_or_counter/quirks/custom_carry_drag.md">Design</a><br><a href="Implementation/Gauge/rolling_drum_or_counter/quirks/custom_carry_drag.md">Implementation</a></td><td>Visual carry-drag behaviour during digit carry.</td></tr>
    <tr><td><code>snap_settle</code></td><td><code>custom_odometer</code></td><td>Implemented</td><td>v3.5.14</td><td><code>realism.snap_settle</code></td><td><a href="Designs/Gauge/rolling_drum_or_counter/quirks/custom_snap_settle.md">Design</a><br><a href="Implementation/Gauge/rolling_drum_or_counter/quirks/custom_snap_settle.md">Implementation</a></td><td>Finite settle/snap behaviour after odometer wheel movement.</td></tr>
    <tr><td><code>backlash</code></td><td><code>custom_odometer</code></td><td>Planned</td><td>v3.7.1</td><td><code>realism.backlash</code> not accepted</td><td><a href="Designs/Gauge/rolling_drum_or_counter/quirks/custom_backlash.md">Design</a><br><span>Implementation pending</span></td><td>v3.5.15 was corrected as not implemented on <code>main</code>; v3.7.1 plans odometer backlash as future work.</td></tr>
  </tbody>
</table>

### `indicator_lamp`

<table>
  <thead>
    <tr><th>Name</th><th>Gauge</th><th>Status</th><th>Version</th><th>Current config key</th><th>Quirk/Gauge doc</th><th>Notes</th></tr>
  </thead>
  <tbody>
    <tr><td><code>custom_indicator</code></td><td><code>custom_indicator</code></td><td>Implemented</td><td>v3.4</td><td><code>type: indicator</code></td><td><a href="Designs/Gauge/indicator_lamp/gauges/custom_indicator.md">Design</a><br><a href="Implementation/Gauge/indicator_lamp/gauges/custom_indicator.md">Implementation</a></td><td>Current GoDriveLog two-state indicator renderer, mapped to the mechanism-first <code>indicator_lamp</code> group.</td></tr>
    <tr><td><code>thermal_fade</code></td><td><code>custom_indicator</code></td><td>Implemented</td><td>v3.5.12</td><td><code>realism.thermal_fade</code></td><td><a href="Designs/Gauge/indicator_lamp/quirks/custom_thermal_fade.md">Design</a><br><a href="Implementation/Gauge/indicator_lamp/quirks/custom_thermal_fade.md">Implementation</a></td><td>Indicator on/off alpha transition to simulate finite lamp response.</td></tr>
  </tbody>
</table>

### `segmented_display`

<table>
  <thead>
    <tr><th>Name</th><th>Gauge</th><th>Status</th><th>Version</th><th>Current config key</th><th>Quirk/Gauge doc</th><th>Notes</th></tr>
  </thead>
  <tbody>
    <tr><td><code>custom_numeric</code></td><td><code>custom_numeric</code></td><td>Implemented</td><td>v3.4</td><td><code>type: numeric</code></td><td><a href="Designs/Gauge/segmented_display/gauges/custom_numeric.md">Design</a><br><a href="Implementation/Gauge/segmented_display/gauges/custom_numeric.md">Implementation</a></td><td>Current GoDriveLog formatted-value image-slot renderer; old <code>seven_segment</code> naming was replaced by <code>numeric</code>.</td></tr>
    <tr><td><code>per_digit_response_lag</code></td><td><code>custom_numeric</code></td><td>Planned</td><td>v3.7.2</td><td><code>realism.per_digit_response_lag</code> not accepted</td><td><a href="Designs/Gauge/segmented_display/quirks/custom_per_digit_response_lag.md">Design</a><br><span>Implementation pending</span></td><td>Planned v3.7 numeric/segmented realism slice; no current code support yet.</td></tr>
    <tr><td><code>leading_zero_behaviour</code></td><td><code>custom_numeric</code></td><td>Planned</td><td>v3.7.3</td><td><code>realism.leading_zero_behaviour</code> not accepted</td><td><a href="Designs/Gauge/segmented_display/quirks/custom_leading_zero_behaviour.md">Design</a><br><span>Implementation pending</span></td><td>Planned v3.7 numeric display behaviour; no current code support yet.</td></tr>
    <tr><td><code>digit_bleed</code></td><td><code>custom_numeric</code></td><td>Planned</td><td>v3.7.4</td><td><code>realism.digit_bleed</code> not accepted</td><td><a href="Designs/Gauge/segmented_display/quirks/custom_digit_bleed.md">Design</a><br><span>Implementation pending</span></td><td>Planned v3.7 numeric/segmented display realism slice; no current code support yet. No separate current <code>segment_bleed</code> design is linked here.</td></tr>
    <tr><td><code>ghosting</code></td><td><code>custom_numeric</code></td><td>Planned</td><td>v3.7.5</td><td><code>realism.ghosting</code> not accepted</td><td><a href="Designs/Gauge/segmented_display/quirks/custom_ghosting.md">Design</a><br><span>Implementation pending</span></td><td>Planned v3.7 numeric/segmented display realism slice; no current code support yet.</td></tr>
    <tr><td><code>uneven_brightness</code></td><td><code>custom_numeric</code></td><td>Planned</td><td>v3.7.6</td><td><code>realism.uneven_brightness</code> not accepted</td><td><a href="Designs/Gauge/segmented_display/quirks/custom_uneven_brightness.md">Design</a><br><span>Implementation pending</span></td><td>Planned v3.7 numeric/segmented display realism slice; no current code support yet.</td></tr>
    <tr><td><code>load_sag</code></td><td><code>custom_numeric</code></td><td>Planned</td><td>v3.7.7</td><td><code>realism.load_sag</code> not accepted</td><td><a href="Designs/Gauge/segmented_display/quirks/custom_load_sag.md">Design</a><br><span>Implementation pending</span></td><td>Planned v3.7 numeric/segmented display realism slice; no current code support yet.</td></tr>
  </tbody>
</table>

## Release planning checkpoints

### v3.7

<table>
  <thead>
    <tr>
      <th>Name</th>
      <th>Status</th>
      <th>Version</th>
      <th>Current config key</th>
      <th>Quirk/Gauge doc</th>
      <th>Notes</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>v3.7 release planning docs</td>
      <td>Planned</td>
      <td>v3.7.0</td>
      <td>-</td>
      <td><a href="v3.7/ReleasePlan.md">Design</a><br><a href="v3.7/ImplementationState.md">Implementation</a></td>
      <td>Current v3.7 state is not started.</td>
    </tr>
    <tr>
      <td>Tests, previews and docs checkpoint</td>
      <td>Planned</td>
      <td>v3.7.8</td>
      <td>-</td>
      <td><a href="v3.7/ImplementationState.md">Implementation</a></td>
      <td>Planned release checkpoint slice.</td>
    </tr>
  </tbody>
</table>

## Notes for future updates

- This file should be updated whenever code truth changes.
- Do not mark a row `Implemented` unless current code supports the documented behaviour in scope.
- Parser support alone is not enough for `Implemented` when the design requires runtime/rendering behaviour.
- `Planned` rows are allowed while a release is active.
- At release close, planned rows should be updated to the current truth: `Implemented`, `Partially implemented`, `Not implemented`, `Deferred`, `Superseded`, or left `Planned` only if the release remains open.
- Do not add every historical idea here. Track accepted release/slice scope and current code truth.
- Do not link to retired `Designs/RealismBehaviour` or `Implementation/RealismBehaviour` paths.
- Do not link to implementation documentation for planned work until the implementation doc exists.
- Planned rows may link to current design docs only; implementation links should be added when the slice lands.
