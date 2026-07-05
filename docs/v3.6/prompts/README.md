# v3.6 Prompt Index

These prompt files define the v3.6 implementation slices.

If the user says any of the following:

- "do the next v3.6 slice"
- "continue v3.6"
- "start the next v3.6 slice"
- "implement v3.6.x"

then the agent must:

1. Read `docs/v3.6/ImplementationState.md`.
2. Use `Current target` only if it is explicitly set and still unchecked; otherwise find the first unchecked allowed slice.
3. Read `docs/v3.6/ReleasePlan.md`.
4. Read the matching prompt file under `docs/v3.6/prompts/`.
5. Make only that slice's changes.
6. Update `docs/v3.6/ImplementationState.md` and any relevant docs.
7. Do not implement later slices early.
8. After the slice is complete, follow the finalisation / PR cycle in `docs/v3.6/ImplementationState.md`.

When the user names a specific version such as `implement v3.6.5`, use the matching prompt file below and confirm it is allowed by `docs/v3.6/ImplementationState.md` before making code changes.

## Prompt files

- `v3.6.0-pointer-marker-docs.md`
- `v3.6.1-radial-pointer-marker-max.md`
- `v3.6.2-radial-pointer-marker-min.md`
- `v3.6.3-pointer-marker-reset-session.md`
- `v3.6.4-radial-damped-secondary-pointer.md`
- `v3.6.5-bar-pointer-marker-max.md`
- `v3.6.6-bar-pointer-marker-min.md`
- `v3.6.7-bar-damped-secondary-pointer.md`
- `v3.6.8-enhancement-backlog.md`

## Shared rules

- Keep slices small.
- Prefer deterministic behaviour.
- Preserve existing gauge behaviour when new config is absent or disabled.
- Keep pointer markers display-only.
- Never mutate source values, logs, exports, configured ranges, or input data.
- Pointer markers follow the rendered indicator path for the gauge family.
- Do not add a `source: value` / `source: pointer` switch in v3.6.
- Do not call the damped secondary marker a mathematical average.
- Do not implement persistence unless a prompt explicitly asks for it.
