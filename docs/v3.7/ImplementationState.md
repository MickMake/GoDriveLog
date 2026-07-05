# v3.7 Implementation State

Status: active planning; v3.7.0 audit slice prepared

Current target: v3.7.0 odometer backlash audit

Current branch: `v37audit`

## Scope

v3.7 is the gauge realism follow-up release area created after v3.6 was narrowed to pointer markers only.

The first v3.7 task is an audit, not an implementation. It must determine whether odometer `backlash` is already implemented, partially implemented, missing, or wrongly documented before any code change is planned.

## Current decisions

- v3.6 remains pointer-marker-only.
- v3.7 may own broader gauge realism audit/backlog work.
- Odometer `backlash` cleanup is the first candidate, but v3.7.0 is audit-only.
- Do not implement `backlash` until the audit produces a clear finding.
- Numeric and segmented display realism candidates need separate design before implementation.
- Bar realism audit work must not contradict completed v3.5 docs without a focused code audit.
- Do not bundle unrelated gauge families into one slice.

## Candidate checklist

- [ ] v3.7.0 odometer backlash audit
- [ ] v3.7.1 odometer backlash implementation decision / prompt, if audit confirms missing support
- [ ] numeric/segmented display realism design
- [ ] bar realism implementation audit
- [ ] indicator realism scope review
- [ ] marker persistence/statistics/styling follow-up review

## v3.7.0 audit questions

The v3.7.0 audit must answer:

1. Is `backlash` documented as implemented or completed anywhere in v3.5 docs?
2. Does config parsing currently recognise `realism.backlash`?
3. Does validation allow or reject `backlash`, and for which gauge families?
4. Does odometer runtime code apply direction-change slack?
5. Are tests present for parsing, validation, runtime reversal behaviour, and final settling?
6. Are preview fixtures or examples present?
7. If support is missing, should the fix be a v3.7 slice or a smaller bugfix/cleanup PR?

## Next-slice workflow

When v3.7 is explicitly activated:

1. Read this file.
2. Read `docs/v3.7/ReleasePlan.md`.
3. Read `docs/v3.7/PlannedFeatures.md`.
4. Read the matching prompt under `docs/v3.7/prompts/`.
5. Make only that slice's changes.
6. For v3.7.0, produce an audit finding before proposing implementation.
7. Do not implement broad backlog/audit material in one slice.
8. Update docs and checklist after the slice.
