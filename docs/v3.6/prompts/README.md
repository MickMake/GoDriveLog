# v3.6 Prompt Index

These prompt files define the v3.6 pointer marker implementation slices.

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
- `v3.6.1-shared-pointer-marker-config-state.md`
- `v3.6.2-shared-min-max-marker-engine.md`
- `v3.6.3-radial-pointer-marker-rendering.md`
- `v3.6.4-bar-pointer-marker-rendering.md`
- `v3.6.5-average-pointer-marker-engine.md`
- `v3.6.6-average-pointer-marker-rendering.md`
- `v3.6.7-pointer-marker-tests-previews-docs.md`

## Shared rules

- Keep slices small.
- Prefer deterministic behaviour.
- Preserve existing gauge behaviour when new config is absent or disabled.
- Keep pointer markers display-only.
- Never mutate source values, logs, exports, configured ranges, or input data.
- Pointer markers sample the final rendered indicator position for the gauge family.
- Do not sample source values, logs, exports, clean mapped values, or pre-realism values.
- Do not add a `source: value` / `source: pointer` switch in v3.6.
- Support simple boolean marker config only: `max`, `min`, and `average`.
- Reject long-form marker object config and unknown marker keys.
- Keep `average` as a highly damped visual pointer, not a statistical average.
- Use the fixed 10 second average pointer time constant in v3.6.
- Render pointer markers above the live needle/bar and below overlay/glass/bezel/frame layers.
- Use explicit marker PNG assets where provided.
- Do not implement persistence in v3.6.
- Do not implement odometer, numeric, segmented, indicator, or broad realism-audit work in v3.6.
