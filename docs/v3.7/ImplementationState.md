# v3.7 Implementation State

Status: backlog scaffold prepared; implementation slices not started

Current target: none

Current branch: `docs/v3.6-planning`

## Scope

v3.7 is the holding area for gauge realism follow-up work that was intentionally removed from the v3.6 pointer marker release.

v3.7 is not active implementation scope yet. Treat these docs as planning/audit material until a future release decision promotes one small slice.

## Current decisions

- v3.6 remains pointer-marker-only.
- v3.7 may own broader gauge realism audit/backlog work.
- Odometer `backlash` cleanup is a candidate v3.7 slice, not a v3.6 tail slice.
- Numeric and segmented display realism candidates need separate design before implementation.
- Bar realism audit work must not contradict completed v3.5 docs without a focused code audit.
- Do not implement anything from v3.7 as part of a v3.6 slice.

## Candidate checklist

- [ ] v3.7.0 odometer backlash cleanup audit/implementation decision
- [ ] numeric/segmented display realism design
- [ ] bar realism implementation audit
- [ ] indicator realism scope review
- [ ] marker persistence/statistics/styling follow-up review

## Next-slice workflow

When v3.7 is explicitly activated:

1. Read this file.
2. Read `docs/v3.7/ReleasePlan.md`.
3. Read `docs/v3.7/PlannedFeatures.md`.
4. Pick one small candidate only.
5. Create or use a matching prompt under `docs/v3.7/prompts/`.
6. Do not implement broad backlog/audit material in one slice.
7. Update docs and checklist after the slice.
