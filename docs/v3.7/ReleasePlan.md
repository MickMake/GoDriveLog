# v3.7 Release Plan

v3.7 is the future gauge realism follow-up area created by moving non-pointer-marker material out of v3.6.

v3.7 is not yet an active implementation release. It is a planning shelf for work that needs sharper scope before code changes begin.

## Theme

Keep v3.6 clean by parking unrelated gauge realism ideas here.

Candidate v3.7 work may include:

- odometer `backlash` cleanup;
- broad realism implementation audits;
- numeric and segmented display realism design;
- bar realism follow-up review;
- indicator realism scope review;
- pointer-marker follow-ups such as persistence, labels, styling, statistics overlays, or reset controls.

## First candidate: odometer backlash cleanup

Odometer `backlash` was removed from v3.6 because it is not pointer-marker work.

Before implementation, v3.7 should confirm:

- whether `backlash` is documented as complete in v3.5;
- whether code support is actually missing;
- whether the fix belongs in v3.7, a v3.5 cleanup, or a dedicated bugfix release;
- whether the intended config and runtime behaviour are still desirable.

If promoted, `backlash` should mean:

> when an odometer-style value reverses direction, wheel movement shows a small bounded amount of mechanical slack before following the new direction and settling exactly on the target.

The implementation must be:

- odometer-only;
- display-only;
- deterministic;
- bounded and subtle;
- disabled unless configured;
- non-mutating for source values, logs, exports, configured ranges, or input data.

## Planning rule

Do not implement anything from v3.7 during v3.6 work.

Promotion should be explicit:

1. choose one small candidate;
2. define its user-facing config;
3. define which gauge families support it;
4. add docs and prompt slice(s);
5. then implement it in a dedicated branch/PR.

## Non-goals

- Do not reopen v3.6 pointer-marker scope from here.
- Do not use v3.7 as proof that v3.5 was wrong without a focused audit.
- Do not bundle unrelated gauge families into one implementation slice.
- Do not create a general physics engine. The dashboard is allowed to be charming, not sentient.
