# v3.6 Implementation State

Status: v3.6 planning scaffold prepared; implementation slices pending

Current target: v3.6.0 pointer marker planning docs

Current branch: `docs/v3.6-planning`

## Scope

v3.6 is the gauge enhancement pass after v3.5.

The first v3.6 theme is pointer markers for radial and bar gauges. Pointer markers are instrument-realism features, not statistical overlays.

Pointer markers observe the rendered indicator path for the gauge family:

- radial gauges observe the rendered needle angle;
- bar gauges observe the rendered bar fill/position.

If no realism effects are enabled, the rendered indicator path is equivalent to the mapped source value, so pointer marker behaviour naturally reflects true input data. If realism effects are enabled, pointer markers follow the realistic rendered behaviour. For example, a radial max marker may capture overshoot if the live pointer actually overshoots.

v3.6 should remain open to other small gauge enhancements, but pointer markers are the initial implementation tranche.

## Current decisions

- Use `realism.pointer_markers` as the config key.
- Do not add a separate `source: value` / `source: pointer` switch in v3.6.
- Min/max pointer markers always reference what the rendered indicator does.
- Radial pointer markers come before bar pointer markers.
- Bar pointer markers should reuse the same semantics after radial behaviour is proven.
- The damped secondary marker is a mechanical/visual feature, not a mathematical average.
- Do not promise a true arithmetic average in v3.6 unless a later prompt explicitly defines statistical averaging.
- Pointer marker state is runtime/session state only unless a later slice explicitly adds persistence.
- v3.6 docs should keep an enhancement backlog so future work can be promoted into slices without bloating the first marker implementation.
- Odometer `backlash` is required as a v3.6 tail implementation slice because existing odometer realism cannot fully create direction-change slack.
- Do not implement `movement: smooth` as a separate odometer mode; `linear`, `ease_out`, and `bell` are already the smooth movement modes.
- Do not implement `movement: click` as a separate odometer mode unless a later prompt defines distinct stepped-wheel behaviour.

## Config key ownership

| Key | Gauge families | Notes |
| --- | --- | --- |
| `pointer_markers.max` | radial, bar | Tracks highest rendered indicator position seen by the marker. |
| `pointer_markers.min` | radial, bar | Tracks lowest rendered indicator position seen by the marker. |
| `pointer_markers.damped` | radial, bar | Slow secondary indicator; not a mathematical average. |
| `backlash` | odometer | Direction-change slack for odometer wheels; v3.6.9 required missing implementation. |

## Scope boundaries

Allowed in v3.6:

- docs and prompts for v3.6 planning;
- radial pointer max/min markers;
- bar pointer max/min markers;
- explicit marker reset/session behaviour;
- a damped secondary marker that behaves like a slow mechanical indicator;
- odometer `backlash` as the promoted v3.6 tail realism cleanup slice;
- small future enhancement candidates documented as backlog items.

Not allowed in the first v3.6 marker tranche:

- database persistence for marker state;
- mathematical average unless explicitly promoted later;
- source/log/export mutation;
- hidden config defaults that alter existing gauges;
- applying radial marker assets to bar gauges without a family-specific rendering plan;
- implementing multiple later slices while doing an earlier one.

## Checklist

- [ ] v3.6.0 pointer marker planning docs
- [ ] v3.6.1 radial pointer marker max
- [ ] v3.6.2 radial pointer marker min
- [ ] v3.6.3 pointer marker reset/session behaviour
- [ ] v3.6.4 radial damped secondary pointer marker
- [ ] v3.6.5 bar pointer marker max
- [ ] v3.6.6 bar pointer marker min
- [ ] v3.6.7 bar damped secondary pointer marker
- [ ] v3.6.8 enhancement backlog triage
- [ ] v3.6.9 implement odometer backlash cleanup

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
