# Future Implementation

This directory is a parking lot for approved or desired follow-on implementation ideas that are not part of the current implementation slice.

Use this to capture "oh, also implement this later" notes without making the active slice ambiguous. Future prompts may reference this directory, but items here are not current scope unless a later prompt explicitly promotes them.

## Boundary with RealismBehaviourGuide

Gauge realism behaviour definitions live in [`../RealismBehaviourGuide/`](../RealismBehaviourGuide/).

Use this directory for implementation tickets only:

- what should be implemented later;
- why it matters;
- dependencies and risks;
- rough effort;
- suggested implementation order;
- links to the canonical behaviour definition.

Do not redefine realism behaviour here. If a future implementation item needs to explain what a realism option means, link to the matching Realism Behaviour Guide page instead.

In short:

```text
RealismBehaviourGuide = definition / behaviour / real-world simulation
FutureImplementation = implementation ticket / backlog / build plan
```

## Indexed FutureSlices table

Effort is rough **Codex hours**, assuming the v3 dashboard/gauge code is already loaded in context and tests exist. Not human hours. Codex hours are weird little mushroom hours.

| # | State | Description | Area | Effort | File |
|---:|---|---|---|---:|---|
| 1 | Implemented in v3.5.19 | Bar gauge overshoot follow-up | `gauge/bar`, `realism.overshoot`, animation | IMPLEMENTED | [`gauge-bar-overshoot-follow-up.md`](gauge-bar-overshoot-follow-up.md) |
| 2 | Near / needs spec tightening | Radial movement options | `gauge/radial`, movement policy, runtime animation | 3-5 | [`gauge-radial-movement-options.md`](gauge-radial-movement-options.md) |
| 3 | Later / visual feature | Radial needle trail | `gauge/radial`, renderer, animation history | 4-7 | [`gauge-radial-needle-trail.md`](gauge-radial-needle-trail.md) |
| 4 | Implemented as v3.6 pointer markers | Gauge stat markers | `gauge/radial`, `gauge/bar`, renderer, rolling-window statistics, marker assets | IMPLEMENTED | [`gauge-stat-markers.md`](gauge-stat-markers.md) |
| 5 | Medium / useful soon | Value zones / warning-danger assets | `gauge/assets`, renderer, config validation | 4-7 | [`gauge-assets-value-zones-warning-danger-assets.md`](gauge-assets-value-zones-warning-danger-assets.md) |
| 6 | Medium / foundational logging | Canonical GoDriveLog Event Log | logging, sensor events, schema/versioning | 5-9 | [`logger-canonical-event-log.md`](logger-canonical-event-log.md) |
| 7 | Medium / pairs with event log | Session metadata sidecar | logging, replay metadata, config provenance | 4-7 | [`logger-session-metadata-sidecar.md`](logger-session-metadata-sidecar.md) |
| 8 | Medium / high value dev tool | JSONL dashboard replay | dashboard runtime, logs, replay CLI | 6-10 | [`dashboard-jsonl-replay.md`](dashboard-jsonl-replay.md) |
| 9 | Near / bounded utility | JSONL log validation | logs, CLI, schema validation | 3-5 | [`logger-jsonl-log-validation.md`](logger-jsonl-log-validation.md) |
| 10 | Later / architecture boundary | External converter boundary | `tools/converters`, import/export architecture | 3-6 | [`tools-converters-external-converter-boundary.md`](tools-converters-external-converter-boundary.md) |
| 11 | Near / performance polish | Needle Animation Performance | `gauge/radial`, animation loop, renderer, Pi4 performance | 3-6 | [`gauge-radial-animation-performance.md`](gauge-radial-animation-performance.md) |
| 12 | Later / gauge realism | Gauge power lifecycle | `gauge/radial`, `gauge/bar`, indicator, numeric, odometer, dashboard runtime, ACC/power events | 6-10 | [`gauge-power-lifecycle.md`](gauge-power-lifecycle.md) |
| 13 | Later / gauge realism | Gauge lighting mode | `gauge/radial`, `gauge/bar`, indicator, numeric, odometer, dashboard runtime, lights-state events, alternate asset sets | 4-7 | [`gauge-lighting-mode.md`](gauge-lighting-mode.md) |
| 14 | Later / gauge realism | Gauge imperfections | `gauge/radial`, `gauge/bar`, indicator, numeric, display artefacts, mechanical wear, electrical artefacts | 7-12 | [`gauge-imperfections.md`](gauge-imperfections.md) |
| 15 | Later / config reuse | Gauge presets | `gauge/config`, `gauge/assets`, `gauge/realism`, config loading, validation, reusable gauge profiles | 5-9 | [`gauge-presets.md`](gauge-presets.md) |

## Extracted historical/planning notes

| Description | File |
|---|---|
| Status legend | [`v37-status-legend.md`](v37-status-legend.md) |
| Gauge realism map | [`v37-gauge-realism-map.md`](v37-gauge-realism-map.md) |
| Candidate: odometer backlash | [`gauge-odometer-backlash.md`](gauge-odometer-backlash.md) |
| Odometer movement cleanup candidates | [`gauge-odometer-movement-cleanup-candidates.md`](gauge-odometer-movement-cleanup-candidates.md) |
| Indicator realism scope | [`gauge-indicator-realism-scope.md`](gauge-indicator-realism-scope.md) |
| Numeric and segmented display realism candidates | [`gauge-numeric-segmented-display-realism-candidates.md`](gauge-numeric-segmented-display-realism-candidates.md) |
| Bar realism scope | [`gauge-bar-realism-scope.md`](gauge-bar-realism-scope.md) |
| Planning rule | [`future-implementation-planning-rule.md`](future-implementation-planning-rule.md) |

These files may contain historical planning context. Current realism behaviour definitions should be checked in [`../RealismBehaviourGuide/`](../RealismBehaviourGuide/) before using them for implementation work.

A few notes from the table:

- **#1 is implemented.** It stays indexed here only as historical follow-up context.
- **#9 and #11** are the most directly sliceable unimplemented items.
- **#6-#8** are chunky but important because they form the logging/replay backbone.
- **#3 is fun, but it’s polish. Nice polish, but still polish.**
- **#11 should not change the look.** It should protect the look from dropped/undersampled frames. That distinction is the whole ballgame.

## Guidelines

- Keep entries small and slice-shaped.
- Mark ideas as `deferred`, `desired`, `exploratory`, or `rejected`.
- Do not treat this directory as an implementation checklist.
- Do not let vague mentions here expand the current slice.
- Prefer a later dedicated prompt/spec before implementation.
- Link to `../RealismBehaviourGuide/` for realism behaviour definitions instead of redefining them here.
