# v3.7 Implementation State

Status: active planning; v3.7.0 implementation slice prepared

Current target: v3.7.0 implement odometer backlash cleanup

Current branch: `v37audit`

## Scope

v3.7 is the gauge realism follow-up release area created after v3.6 was narrowed to pointer markers only.

The first v3.7 task is odometer `backlash`: audit the current code/docs state, then implement missing support if the audit confirms it is absent or incomplete.

## Current decisions

- v3.6 remains pointer-marker-only.
- v3.7 may own broader gauge realism audit/backlog work.
- Odometer `backlash` cleanup is required as the first v3.7 implementation slice.
- v3.7.0 must still inspect current code before changing it, but implementation is in scope.
- Numeric and segmented display realism candidates need separate design before implementation.
- Bar realism audit work must not contradict completed v3.5 docs without a focused code audit.
- Do not bundle unrelated gauge families into one slice.

## Candidate checklist

- [ ] v3.7.0 implement odometer backlash cleanup
- [ ] numeric/segmented display realism design
- [ ] bar realism implementation audit
- [ ] indicator realism scope review
- [ ] marker persistence/statistics/styling follow-up review

## v3.7.0 requirements

The v3.7.0 slice must:

1. Confirm the current `backlash` docs/code state.
2. Add config parsing support for `realism.backlash` if missing.
3. Restrict `backlash` to odometer gauges.
4. Add deterministic odometer direction-change slack.
5. Preserve existing odometer behaviour when `backlash` is absent or false.
6. Keep source values, logs, exports, configured ranges, and input data unchanged.
7. Add parsing, validation, runtime, and final-settling tests.
8. Add or update preview/docs examples.

## Next-slice workflow

When v3.7 is explicitly activated:

1. Read this file.
2. Read `docs/v3.7/ReleasePlan.md`.
3. Read `docs/v3.7/PlannedFeatures.md`.
4. Read the matching prompt under `docs/v3.7/prompts/`.
5. Make only that slice's changes.
6. For v3.7.0, audit first, then implement only the confirmed missing/incomplete `backlash` support.
7. Do not implement broad backlog/audit material in one slice.
8. Update docs and checklist after the slice.
