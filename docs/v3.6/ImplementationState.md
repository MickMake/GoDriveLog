# v3.6 Implementation State

Status: v3.6 planning scaffold prepared; implementation slices pending

Current target: v3.6.0 pointer marker planning docs

Current branch: `docs/v3.6-pointer-marker-spec-fix`

## Scope

v3.6 is the pointer marker gauge enhancement pass after the completed v3.5 gauge realism pass.

Pointer markers are instrument-realism features, not statistical overlays.

Pointer markers observe the final rendered indicator position for the gauge family:

- radial gauges observe the final rendered needle position;
- bar gauges observe the final rendered bar fill/indicator position.

If no realism effects are enabled, the rendered indicator path is equivalent to the mapped source value, so pointer marker behaviour naturally reflects true input data. If realism effects are enabled, pointer markers follow the realistic rendered behaviour. For example, a radial max marker may capture overshoot if the live pointer visibly overshoots.

v3.6 must stay focused on pointer markers only. Broader gauge realism audit/backlog material belongs in v3.7 or later.

## Current decisions

- Use `realism.pointer_markers` as the config key.
- Support simple boolean marker enables: `max: true`, `min: true`, `average: true`.
- Support optional `window: <duration>` for rolling min/max history.
- Do not support long-form marker objects such as `max.enabled` in v3.6.
- Do not support top-level `pointer_markers: true`.
- Unknown keys under `pointer_markers` must fail config loading clearly.
- Do not add a separate `source: value` / `source: pointer` switch in v3.6.
- Pointer markers always sample the final rendered indicator position.
- Pointer markers must not sample source values, logs, exports, clean mapped values, or pre-realism values.
- Support radial and bar gauges in v3.6.
- Use a shared marker engine where practical; keep rendering family-specific.
- Min/max markers use daily local reset mode when `window` is absent.
- Min/max markers use rolling-window history when `window` is present.
- `average` is an old-style highly damped pointer marker, not a mathematical average.
- `average` uses a fixed 10 second time constant in v3.6.
- Pointer markers render above the live needle/bar and below overlay/glass/bezel/frame layers.
- Pointer markers use explicit marker PNG assets where provided.
- Pointer marker state is runtime-only; do not add database persistence in v3.6.
- Keep future gauge realism audit/backlog material out of v3.6.

## Config key ownership

| Key | Gauge families | Notes |
| --- | --- | --- |
| `pointer_markers.max` | radial, bar | Tracks the furthest/highest final rendered indicator position in the active min/max history mode. |
| `pointer_markers.min` | radial, bar | Tracks the lowest/least final rendered indicator position in the active min/max history mode. |
| `pointer_markers.average` | radial, bar | 10 second highly damped average-style pointer marker; not a statistical average. |
| `pointer_markers.window` | radial, bar | Optional positive finite rolling duration for min/max history only. |

## Asset ownership

| Gauge family | Marker assets | Notes |
| --- | --- | --- |
| radial | `needle_min.png`, `needle_max.png`, `needle_average.png` | Same pivot/rotation model as live radial needle. Render above live needle and below overlay/glass/bezel. |
| bar | `marker_min.png`, `marker_max.png`, `marker_average.png` | Placed along the bar axis. Respect horizontal/vertical/reversed/origin direction. Render above live bar/fill and below overlay/frame/glass. |

## Scope boundaries

Allowed in v3.6:

- docs and prompts for v3.6 pointer marker planning;
- shared pointer marker config/state model;
- shared min/max marker engine;
- radial min/max marker rendering;
- bar min/max marker rendering;
- daily local reset mode for min/max when `window` is absent;
- rolling-window min/max history when `window` is present;
- a shared 10 second damped average marker engine;
- radial and bar average marker rendering;
- explicit marker PNG asset support;
- final pointer-marker tests, previews, and docs/checkpoint work.

Not allowed in v3.6:

- database persistence for marker state;
- mathematical average/statistical overlays;
- source/log/export mutation;
- hidden config defaults that alter existing gauges;
- long-form pointer marker object config;
- procedural replacement of explicit marker PNG assets unless the gauge renderer already has a documented marker convention;
- odometer backlash cleanup;
- v3.5 realism implementation audits;
- broad future gauge realism matrices;
- implementing multiple later slices while doing an earlier one.

## Checklist

- [ ] v3.6.0 pointer marker planning docs
- [ ] v3.6.1 shared pointer marker config/state
- [ ] v3.6.2 shared min/max marker engine
- [ ] v3.6.3 radial pointer marker rendering
- [ ] v3.6.4 bar pointer marker rendering
- [ ] v3.6.5 average pointer marker engine
- [ ] v3.6.6 average pointer marker rendering
- [ ] v3.6.7 tests, previews, docs checkpoint

## Next-slice workflow

When asked to do the next v3.6 slice:

1. Read this file.
2. Use `Current target` only if it is explicitly set and its checklist item is still unchecked; otherwise find the first unchecked allowed slice.
3. Read `docs/v3.6/ReleasePlan.md`.
4. Read the matching prompt in `docs/v3.6/prompts/`.
5. Make only that slice's changes.
6. Update this checklist and any relevant docs.
7. Advance `Current target` to the next intended unchecked slice, or set it to `none` when v3.6 is complete.
8. Do not implement later slices early.
9. Run the relevant local tests/checks.
10. Commit the completed slice with a clear message.
11. Push the branch to GitHub.
12. Raise a pull request against main.

Then enter the review-fix loop:

13. Wait for Codex GitHub review feedback and CI results.
14. If CI fails, review requests changes, or unresolved review comments require code changes:
    * inspect the feedback;
    * make the smallest safe fixes only;
    * do not refactor unrelated code;
    * rerun relevant tests/checks;
    * commit and push the fixes.
15. Repeat the review-fix loop at most 3 times.

Stop when either:

- CI passes and review is clear enough to merge; or
- the review-fix loop hits the limit and needs human direction.
