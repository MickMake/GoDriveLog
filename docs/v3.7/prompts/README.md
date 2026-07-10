# v3.7 Prompt Index

These prompt files define the v3.7 odometer and image-based numeric display realism implementation slices.

If the user says any of the following:

- "do the next v3.7 slice"
- "continue v3.7"
- "start the next v3.7 slice"
- "implement v3.7.x"

then the agent must:

1. Read `docs/v3.7/ImplementationState.md`.
2. Use `Current target` only if it is explicitly set and still unchecked; otherwise find the first unchecked allowed slice.
3. Read `docs/v3.7/ReleasePlan.md`.
4. Read the matching prompt file under `docs/v3.7/prompts/`.
5. Inspect the current implementation before proposing code changes.
6. Make only that slice's changes.
7. Update `docs/v3.7/ImplementationState.md` and relevant realism-guide docs.
8. Do not implement later slices early.
9. After the slice is complete, follow the finalisation and PR cycle in `docs/v3.7/ImplementationState.md`.

When the user names a specific version such as `implement v3.7.5`, use the matching prompt file and confirm it is allowed by `docs/v3.7/ImplementationState.md` before making code changes.

## Prompt files

- `v3.7.0-release-planning-docs.md`
- `v3.7.1-odometer-backlash.md`
- `v3.7.2-per-digit-response-lag.md`
- `v3.7.3-leading-zero-behaviour.md`
- `v3.7.4-segment-digit-bleed.md`
- `v3.7.5-ghosting.md`
- `v3.7.6-uneven-brightness.md`
- `v3.7.7-load-sag.md`
- `v3.7.8-tests-previews-docs.md`

## Shared rules

- KISS.
- Keep slices small and independently reviewable.
- Extend the existing renderer directly.
- Do not redesign or replace the renderer.
- Keep each realism option locally owned unless a tiny shared helper is obviously simpler.
- Do not create a generic realism engine.
- Do not create a shared numeric-display runtime framework.
- Preserve existing gauge behaviour when the new option is absent or disabled.
- Keep new realism display-only, deterministic, bounded, and subtle.
- Never mutate source values, logs, exports, configured ranges, or input data.
- Keep the whole-image `segmented` renderer separate from the image-based numeric/seven-segment renderer.
- Do not move existing segmented hysteresis into a new shared display abstraction.
- Use supplied image assets and existing scene composition where practical.
- Do not inspect image pixels to infer segment geometry or display load.
- Define and document the exact config shape in the active slice before implementation.
- Add focused tests for enabled and disabled behaviour.
- Add or update a dedicated preview for each realism option.
- Do not refactor unrelated code.
- Stop after the active slice is complete.

## Slice boundaries

### v3.7.0 release planning docs

Activate the release and prepare the release plan, implementation state, prompt index, and per-slice prompt files.

No runtime implementation belongs in this slice.

### v3.7.1 odometer backlash

Add odometer-only bounded slack after direction reversal.

Do not change ordinary forward movement or unrelated odometer realism options.

### v3.7.2 per-digit response lag

Add small deterministic update delays to numeric digit slots.

Any runtime state introduced must remain feature-local unless later reuse is already concrete and simpler.

### v3.7.3 leading-zero behaviour

Add presentation-only handling for leading zero slots.

Do not alter formatting, numeric meaning, source values, or slot alignment.

### v3.7.4 segment and digit bleed

Add subtle inactive segment/digit imagery using explicit image assets or existing documented layers.

Do not build a procedural segment renderer.

### v3.7.5 ghosting

Add finite previous-character fade behaviour for changed numeric slots.

Do not retain unbounded history or make mixed values long-lived.

### v3.7.6 uneven brightness

Add stable deterministic per-slot brightness variation.

Do not introduce per-frame random flicker.

### v3.7.7 load sag

Add subtle brightness reduction based on known displayed segment load.

Use a known character load table or explicit configuration, not pixel analysis.

### v3.7.8 tests, previews, docs checkpoint

Verify config validation, runtime behaviour, finite settling, disabled behaviour, existing gauge regression coverage, preview packages, and realism-guide status.

Do not introduce new realism options in the checkpoint slice. The checkpoint is a broom, not another room.
