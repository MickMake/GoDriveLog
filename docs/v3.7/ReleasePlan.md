# v3.7 Release Plan

v3.7 is the gauge realism follow-up area created after v3.6 was narrowed to pointer markers only.

v3.7 starts with odometer `backlash` cleanup. The slice must still audit current code before changing it, but implementation is in scope because the feature is required.

## Theme

Keep v3.6 clean by parking unrelated gauge realism ideas here, then promote only one small, verified item at a time.

Candidate v3.7 work may include:

- odometer `backlash` cleanup;
- broad realism implementation audits;
- numeric and segmented display realism design;
- bar realism follow-up review;
- indicator realism scope review;
- pointer-marker follow-ups such as persistence, labels, styling, statistics overlays, or reset controls.

## First slice: odometer backlash cleanup

Odometer `backlash` was removed from v3.6 because it is not pointer-marker work.

v3.7.0 must audit and then implement missing/incomplete support for:

```yaml
realism:
  backlash: true
```

The slice should confirm:

- whether `backlash` is documented as complete in v3.5;
- whether config parsing recognises `realism.backlash`;
- whether validation gates `backlash` to odometer gauges;
- whether runtime code implements direction-change slack;
- whether tests cover parsing, validation, reversal behaviour, and final settling;
- whether preview fixtures/examples exist.

`backlash` should mean:

> when an odometer-style value reverses direction, wheel movement shows a small bounded amount of mechanical slack before following the new direction and settling exactly on the target.

The implementation must be:

- odometer-only;
- display-only;
- deterministic;
- bounded and subtle;
- disabled unless configured;
- non-mutating for source values, logs, exports, configured ranges, or input data.

## Planning rule

Do not implement anything from v3.7 during a different slice.

Promotion should be explicit:

1. choose one small candidate;
2. audit the current code/docs state;
3. define its user-facing config;
4. define which gauge families support it;
5. add or update docs and prompt slice(s);
6. then implement it in a dedicated branch/PR.

## Non-goals

- Do not reopen v3.6 pointer-marker scope from here.
- Do not use v3.7 as proof that v3.5 was wrong without a focused audit.
- Do not bundle unrelated gauge families into one implementation slice.
- Do not create a general physics engine. The dashboard is allowed to be charming, not sentient.
