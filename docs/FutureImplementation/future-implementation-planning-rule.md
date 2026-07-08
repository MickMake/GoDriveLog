# Planning Rule

Origin: `docs/v3.7/PlannedFeatures.md`

Do not implement anything from this file as part of an unrelated slice.

Promotion should be explicit:

1. choose one small candidate;
2. define its user-facing config;
3. define which gauge families support it;
4. add docs and prompt slice(s);
5. then implement it in a dedicated branch/PR.
