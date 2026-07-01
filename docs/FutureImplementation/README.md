# Future Implementation

This directory is a parking lot for approved or desired follow-on ideas that are not part of the current implementation slice.

Use this to capture "oh, also implement this later" notes without making the active slice ambiguous. Future prompts may reference this directory, but items here are not current scope unless a later prompt explicitly promotes them.

## Indexed FutureSlices table

Effort is rough **Codex hours**, assuming the v3 dashboard/gauge code is already loaded in context and tests exist. Not human hours. Codex hours are weird little mushroom hours.

| # | State | Description | Area | Effort | File |
|---:|---|---|---|---:|---|
| 1 | Near / implementation-ready | Bar gauge overshoot follow-up | `gauge/bar`, `realism.overshoot`, animation | 2-4 | [`gauge-bar-overshoot-follow-up.md`](gauge-bar-overshoot-follow-up.md) |
| 2 | Near / needs spec tightening | Radial movement options | `gauge/radial`, movement policy, runtime animation | 3-5 | [`gauge-radial-movement-options.md`](gauge-radial-movement-options.md) |
| 3 | Later / visual feature | Radial needle trail | `gauge/radial`, renderer, animation history | 4-7 | [`gauge-radial-needle-trail.md`](gauge-radial-needle-trail.md) |
| 4 | Later / visual feature | Radial peak hold | `gauge/radial`, display marker, state tracking | 3-6 | [`gauge-radial-peak-hold.md`](gauge-radial-peak-hold.md) |
| 5 | Medium / useful soon | Value zones / warning-danger assets | `gauge/assets`, renderer, config validation | 4-7 | [`gauge-assets-value-zones-warning-danger-assets.md`](gauge-assets-value-zones-warning-danger-assets.md) |
| 6 | Medium / foundational logging | Canonical GoDriveLog Event Log | logging, sensor events, schema/versioning | 5-9 | [`logger-canonical-event-log.md`](logger-canonical-event-log.md) |
| 7 | Medium / pairs with event log | Session metadata sidecar | logging, replay metadata, config provenance | 4-7 | [`logger-session-metadata-sidecar.md`](logger-session-metadata-sidecar.md) |
| 8 | Medium / high value dev tool | JSONL dashboard replay | dashboard runtime, logs, replay CLI | 6-10 | [`dashboard-jsonl-replay.md`](dashboard-jsonl-replay.md) |
| 9 | Near / bounded utility | JSONL log validation | logs, CLI, schema validation | 3-5 | [`logger-jsonl-log-validation.md`](logger-jsonl-log-validation.md) |
| 10 | Later / architecture boundary | External converter boundary | `tools/converters`, import/export architecture | 3-6 | [`tools-converters-external-converter-boundary.md`](tools-converters-external-converter-boundary.md) |
| 11 | Near / performance polish | Needle Animation Performance | `gauge/radial`, animation loop, renderer, Pi4 performance | 3-6 | [`gauge-radial-animation-performance.md`](gauge-radial-animation-performance.md) |

A few notes from the table:

- **#1, #9, and #11** are the most directly sliceable.
- **#6-#8** are chunky but important because they form the logging/replay backbone.
- **#3 and #4** are fun, but they’re polish. Nice polish, but still polish.
- **#11 should not change the look.** It should protect the look from dropped/undersampled frames. That distinction is the whole ballgame.

## Guidelines

- Keep entries small and slice-shaped.
- Mark ideas as `deferred`, `desired`, `exploratory`, or `rejected`.
- Do not treat this directory as an implementation checklist.
- Do not let vague mentions here expand the current slice.
- Prefer a later dedicated prompt/spec before implementation.
