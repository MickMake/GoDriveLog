# v3.6 Implementation State

Status: v3.6 planning scaffold prepared; implementation slices pending

Current target: v3.6.0 pointer marker planning docs

Current branch: `docs/v3.6-planning`

## Scope

v3.6 is the pointer marker gauge enhancement pass after the completed v3.5 gauge realism pass.

Pointer markers are instrument-realism features, not statistical overlays.

Pointer markers observe the rendered indicator path for the gauge family:

- radial gauges observe the rendered needle angle;
- bar gauges observe the rendered bar fill/position.

If no realism effects are enabled, the rendered indicator path is equivalent to the mapped source value, so pointer marker behaviour naturally reflects true input data. If realism effects are enabled, pointer markers follow the realistic rendered behaviour. For example, a radial max marker may capture overshoot if the live pointer actually overshoots.

v3.6 must stay focused on pointer markers only. Broader gauge realism audit/backlog material belongs in v3.7 or later.

## Current decisions

- Use `realism.pointer_markers` as the config key.
- Do not add a separate `source: value` / `source: pointer` switch in v3.6.
- Min/max pointer markers always reference what the rendered indicator does.
- Radial pointer markers come before bar pointer markers.
- Bar pointer markers should reuse the same semantics after radial behaviour is proven.
- The damped secondary marker is a mechanical/visual feature, not a mathematical average.
- Do not promise a true arithmetic average in v3.6.
- Pointer marker state is runtime/session state only unless a later release explicitly adds persistence.
- Keep future gauge realism audit/backlog material out of v3.6.

## Config key ownership

| Key | Gauge families | Notes |
| --- | --- | --- |
| `pointer_markers.max` | radial, bar | Tracks highest rendered indicator position seen by the marker. |
| `pointer_markers.min` | radial, bar | Tracks lowest rendered indicator position seen by the marker. |
| `pointer_markers.damped` | radial, bar | Slow secondary indicator; not a mathematical average. |

## Scope boundaries

Allowed in v3.6:

- docs and prompts for v3.6 pointer marker planning;
- radial pointer max/min markers;
- bar pointer max/min markers;
- explicit marker reset/session behaviour;
- a damped secondary marker that behaves like a slow mechanical indicator;
- final pointer-marker docs/checkpoint work.

Not allowed in v3.6:

- database persistence for marker state;
- mathematical average/statistical overlays;
- source/log/export mutation;
- hidden config defaults that alter existing gauges;
- applying radial marker assets to bar gauges without a family-specific rendering plan;
- odometer backlash cleanup;
- v3.5 realism implementation audits;
- broad future gauge realism matrices;
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
- [ ] v3.6.8 pointer marker docs/checkpoint

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
