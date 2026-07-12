# Implementation Index

This index links each design document to its matching implementation record and the current verified implementation status.

## Logging

| Document | Status | Purpose |
|---|---|---|
| [JSONL Dashboard Replay](Logging/dashboard-jsonl-replay.md) | Not implemented | Tracks the missing replay mode for feeding recorded event logs back through the dashboard runtime. |
| [Canonical GoDriveLog Event Log](Logging/logger-canonical-event-log.md) | Partially implemented | Records how far the current JSONL event stream has progressed toward a formal GoDriveLog-owned log format. |
| [JSONL Log Validation](Logging/logger-jsonl-log-validation.md) | Not implemented | Tracks the absent validator for canonical GoDriveLog event logs. |
| [Session Metadata Sidecar](Logging/logger-session-metadata-sidecar.md) | Not implemented | Tracks the proposed sidecar file that would capture replay and provenance metadata next to each event log. |
| [External Converter Boundary](Logging/tools-converters-external-converter-boundary.md) | Not implemented | Tracks the planned boundary that keeps foreign-format conversion out of GoDriveLog core runtime. |

## Realism Behaviour

| Document | Status | Purpose |
|---|---|---|
| [`backlash`](RealismBehaviour/backlash.md) | Not implemented | Tracks the planned odometer backlash effect for direction-change slack. |
| [Bar Gauge Overshoot Follow-Up](RealismBehaviour/bar-overshoot-follow-up.md) | Implemented | Records the follow-up idea for bar-gauge overshoot and how it now exists on `main`. |
| [Bar Realism Scope](RealismBehaviour/bar-realism-scope.md) | Partially implemented | Tracks which planned realism behaviours for bar gauges have landed and which remain backlog items. |
| [`calibration_offset`](RealismBehaviour/calibration-offset.md) | Implemented | Tracks the fixed angular display offset for radial needles. |
| [`damping`](RealismBehaviour/damping.md) | Implemented | Tracks the lag-and-catch-up behaviour for radial and bar gauges. |
| [Gauge Imperfections](RealismBehaviour/gauge-imperfections.md) | Partially implemented | Tracks the broader backlog for visible gauge imperfections across multiple gauge families. |
| [Gauge Power Lifecycle](RealismBehaviour/gauge-power-lifecycle.md) | Not implemented | Tracks the planned gauge-level power-on and power-off realism driven by a dashboard power signal. |
| [Gauge Presets](RealismBehaviour/gauge-presets.md) | Not implemented | Tracks the planned reusable preset/profile layer for gauge visuals and realism. |
| [`hysteresis`](RealismBehaviour/hysteresis.md) | Implemented | Tracks direction-dependent displayed offsets for radial and bar gauges. |
| [`realism.imperfections`](RealismBehaviour/imperfections.md) | Not implemented | Tracks the proposed umbrella `realism.imperfections` config layer. |
| [Indicator Realism Scope](RealismBehaviour/indicator-realism-scope.md) | Partially implemented | Tracks the limited current realism support for indicator gauges. |
| [Gauge Lighting Mode](RealismBehaviour/lighting-mode.md) | Not implemented | Tracks the planned per-gauge reaction to dashboard lights-state changes. |
| [`movement`](RealismBehaviour/movement.md) | Partially implemented | Tracks the single movement knob and its family-specific behaviour across gauges. |
| [`needle_shadow`](RealismBehaviour/needle-shadow.md) | Implemented | Tracks the static shadow/depth cue for radial needles. |
| [Needle Trail](RealismBehaviour/needle-trail.md) | Not implemented | Tracks the planned fading history of previous radial needle positions. |
| [`decimal_point_behaviour`](RealismBehaviour/numeric-decimal-point-behaviour.md) | Not implemented | Tracks the planned independent behaviour rules for decimal points in numeric and segmented displays. |
| [`ghosting`](RealismBehaviour/numeric-ghosting.md) | Not implemented | Tracks the planned residual afterimage effect for numeric and segmented displays. |
| [`leading_zero_behaviour`](RealismBehaviour/numeric-leading-zero-behaviour.md) | Not implemented | Tracks the planned deliberate handling of leading zero slots for numeric and segmented displays. |
| [`load_sag`](RealismBehaviour/numeric-load-sag.md) | Not implemented | Tracks the planned brightness sag effect for high-load numeric and segmented values. |
| [Numeric and Segmented Display Realism Candidates](RealismBehaviour/numeric-segmented-display-realism-candidates.md) | Not implemented | Tracks the backlog of candidate realism behaviours for numeric and segmented displays. |
| [Candidate: Odometer Backlash](RealismBehaviour/odometer-backlash.md) | Not implemented | Tracks the backlog note confirming that odometer backlash is not implemented on `main`. |
| [`carry_drag`](RealismBehaviour/odometer-carry-drag.md) | Implemented | Tracks the early-coupling movement of higher odometer digits near rollover. |
| [`drum_slop`](RealismBehaviour/odometer-drum-slop.md) | Implemented | Tracks the fixed per-wheel alignment imperfection for odometers. |
| [Odometer Movement Cleanup Candidates](RealismBehaviour/odometer-movement-cleanup-candidates.md) | Partially implemented | Tracks the cleanup note for reserved odometer movement values such as `smooth` and `click`. |
| [`overshoot`](RealismBehaviour/overshoot.md) | Implemented | Tracks bounded pass-and-settle movement for radial and bar gauges. |
| [`peg_bounce`](RealismBehaviour/peg-bounce.md) | Implemented | Tracks the tap-rebound-settle behaviour when radial or bar gauges hit display limits. |
| [`per_digit_response_lag`](RealismBehaviour/per-digit-response-lag.md) | Not implemented | Tracks the planned slot-by-slot update lag for numeric and segmented displays. |
| [`quantized_fill`](RealismBehaviour/quantized-fill.md) | Not implemented | Tracks the planned discrete-resolution fill behaviour for bar and segmented displays. |
| [Radial Animation Performance](RealismBehaviour/radial-animation-performance.md) | Partially implemented | Tracks the reliability of subtle radial animations on slower render targets. |
| [Radial Movement Options](RealismBehaviour/radial-movement-options.md) | Partially implemented | Tracks the planned scalar movement options for radial gauges. |
| [Realism Behaviour Guide](RealismBehaviour/realism-behaviour-guide.md) | Partially implemented | Tracks how much of the canonical realism guide currently exists in code. |
| [`segment_bleed` / `digit_bleed`](RealismBehaviour/segment-bleed-digit-bleed.md) | Not implemented | Tracks the planned faint inactive-segment visibility for numeric and segmented displays. |
| [`snap_settle`](RealismBehaviour/snap-settle.md) | Implemented | Tracks the short landing snap for odometer wheels. |
| [Gauge Stat Markers](RealismBehaviour/stat-markers.md) | Not implemented | Preserves the historical note for an older statistical marker concept that should not be implemented as written. |
| [`stepped_fill`](RealismBehaviour/stepped-fill.md) | Not implemented | Tracks the planned block-style fill behaviour for bar and segmented displays. |
| [`stiction`](RealismBehaviour/stiction.md) | Implemented | Tracks thresholded release behaviour for small radial and bar changes. |
| [`thermal_fade`](RealismBehaviour/thermal-fade.md) | Implemented | Tracks incandescent-style warm-up and cool-down for indicators. |
| [`uneven_brightness`](RealismBehaviour/uneven-brightness.md) | Not implemented | Tracks the planned stable slot/region brightness variation for numeric and segmented displays. |
| [Value Zones / Warning-Danger Assets](RealismBehaviour/value-zones-warning-danger-assets.md) | Not implemented | Tracks the planned asset-switching feature that selects warning or danger variants when values enter configured zones. |
| [`pointer_markers`](RealismBehaviour/witness-markers.md) | Implemented | Tracks the pointer-marker feature that records and renders reference positions for final displayed indicator geometry. |
| [`wraparound`](RealismBehaviour/wraparound.md) | Implemented | Tracks continuous odometer wheel routing through digit-strip boundaries. |

## Runtime

| Document | Status | Purpose |
|---|---|---|
| [MQTT Architecture Notes](Runtime/mqtt-architecture.md) | Not implemented | Tracks the future architectural direction for decoupling telemetry producers and consumers through MQTT. |
| [GoDriveLog Pi4 Fyne Kiosk Setup](Runtime/setup-kiosk-mode-on-pi4.md) | Not implemented | Tracks the Raspberry Pi kiosk setup note for a Fyne-based display stack. |
