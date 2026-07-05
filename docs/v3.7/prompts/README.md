# v3.7 Prompt Index

These prompt files define v3.7 follow-up slices. They are not part of v3.6.

Do not run a v3.7 prompt unless the user explicitly activates v3.7 or names the specific v3.7 slice.

## Prompt files

- `v3.7.0-implement-odometer-backlash-cleanup.md` — implement odometer `backlash` cleanup after confirming the current code state.

## Shared rules

- Keep slices small.
- Prefer deterministic display behaviour.
- Preserve existing gauge behaviour when new config is absent or disabled.
- Never mutate source values, logs, exports, configured ranges, or input data.
- Audit current code before assuming an old checklist is wrong.
- Do not bundle unrelated gauge families into one slice.
