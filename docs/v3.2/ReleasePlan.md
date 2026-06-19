# GoDriveLog v3.2 release plan

Status: planning
Owner: migration implementor

## Purpose

This file gives the v3.2 implementation roadmap.

v3.2 temporarily benches the remaining v3.1.7 dashboard event efficiency and v3.1.8 retirement readiness work so the gauge package direction can be advanced while the architecture is clear.

## Release goal

Add a minimal, self-contained gauge package model that lets dashboards place complete gauge instruments without refactoring the existing dashboard/widget architecture.

## Release principles

- Keep the existing dashboard and widget model.
- Add `type: gauge` as an extension, not a rewrite.
- Gauge widgets place gauge packages using a gauge path, position, and scale.
- Gauge packages own their sensor binding, visual layers, value mapping, pivots, and asset references.
- Directory names under `assets/gauges/` are arbitrary and carry no renderer meaning.
- The only required gauge package filename is `gauge.yaml`.
- Gauge type is declared inside `gauge.yaml`.
- Image paths inside `gauge.yaml` are resolved relative to that `gauge.yaml` file.
- File-based reuse is allowed through relative image references, such as `../images/needle.png`.
- Do not add code inheritance, clusters, sensor overrides, themes, variants, or procedural ticks in the first pass.
- Keep dashboard code below the sensor/event boundary.
- Preserve existing widget types and behaviour.

## Benched v3.1 work

The following v3.1 slices are not cancelled. They are deferred until the gauge package direction is in place:

| Version | Slice | Deferred because |
|---|---|---|
| v3.1.7 | dashboard event efficiency | The visual scene model may change for gauge widgets. Optimisation should happen after the new gauge path exists. |
| v3.1.8 | retirement readiness | Retirement readiness should consider the new dashboard gauge direction before old paths are reviewed for deletion/archive. |

## Branch-chat workflow

Each implementation chat should:

1. Read this file.
2. Read `docs/v3.2/prompts/README.md`.
3. Read `docs/v3.2/ImplementationState.md`.
4. Read `docs/v3.2/OpenDecisions.md`.
5. Read the prompt file for the target slice under `docs/v3.2/prompts/`.
6. Confirm the previous relevant PR is merged into `main`.
7. Confirm there are no blocking open PRs.
8. Create a branch from latest `main` using the target version prefix.
9. Implement only that version slice.
10. Update `CHANGES.md` and `docs/v3.2/ImplementationState.md`.
11. Update `docs/v3.2/OpenDecisions.md` only when a decision is resolved, changed, added, or explicitly deferred.
12. Open a PR.
13. Stop.

Do not redesign the release plan inside a slice chat.

## Planned implementation slices

| Version | Slice | Goal |
|---|---|---|
| v3.2.0 | planning baseline | Create the v3.2 planning docs and prompts. |
| v3.2.1 | gauge package loader | Load `assets/gauges/**/gauge.yaml` packages. |
| v3.2.2 | gauge widget support | Add `type: gauge` widgets that place gauge packages. |
| v3.2.3 | radial gauge scene model | Convert radial gauge package + sensor state into dashboard scene parts. |
| v3.2.4 | Fyne radial rendering | Render layered PNG radial gauges and rotate/place the needle. |
| v3.2.5 | example gauge package | Add one working example gauge package. |
| v3.2.6 | harness verification | Exercise gauge widgets through the existing v3 dashboard harness. |
| v3.2.7 | checkpoint | Decide whether to resume v3.1.7/v3.1.8, continue gauge work, or add cluster support later. |

## v3.2.0 planning baseline checkpoints

- `docs/v3.2/` contains the planning documents.
- `docs/v3.2/prompts/` contains one prompt per planned implementation slice.
- v3.1.7 and v3.1.8 are clearly marked as deferred, not cancelled.
- The gauge package direction is documented.
- No Go code, tests, runtime behaviour, config schema, renderer code, or assets are changed.

## v3.2.1 gauge package loader checkpoints

- Load a gauge package from `assets/gauges/**/gauge.yaml`.
- Treat the dashboard `gauge` value as an asset-root relative directory path.
- Resolve `gauge.yaml` from that directory.
- Resolve layer image paths relative to the `gauge.yaml` directory.
- Allow relative reuse such as `../images/bezel.png` when the final path remains inside the asset root.
- Reject missing `gauge.yaml` files with clear errors.
- Reject invalid path traversal outside the asset root.
- Support `type: radial` in parsed gauge packages.

## v3.2.2 gauge widget support checkpoints

- Add `type: gauge` to dashboard widget config.
- Add `gauge` path and `scale` fields needed for gauge placement.
- Preserve all existing widget types and validation behaviour.
- For v3.2, gauge widgets do not define or override sensors.
- The sensor binding comes from `gauge.yaml`.

## v3.2.3 radial gauge scene model checkpoints

- Convert a loaded radial gauge package and current sensor state into scene parts.
- Preserve non-`ok` dashboard semantics: do not render live values for missing, unsupported, timeout, parse_error, error, stale, or unknown states.
- Include enough scene data for the display adapter to draw static layers and rotate/place the needle.
- Keep Fyne-specific rendering out of the dashboard scene model.

## v3.2.4 Fyne radial rendering checkpoints

- Render layer order: background, face, ticks, needle, overlay.
- Rotate the needle around `needle` pivot.
- Place the rotated needle pivot at the gauge `face` pivot.
- Use normalised pivot coordinates.
- Keep existing display adapter behaviour for existing widget types.

## v3.2.5 example gauge package checkpoints

- Add one small example gauge package under `assets/gauges/`.
- Use arbitrary directory names; only `gauge.yaml` is required.
- Demonstrate shared image references with relative paths if practical.
- Keep the example small enough for harness/manual verification.

## v3.2.6 harness verification checkpoints

- Exercise the gauge widget path through the existing v3 dashboard harness.
- Use fake sensor events through the real dashboard runtime path.
- Confirm fixed and sweep values behave as expected.
- Confirm non-`ok` states do not render fake live gauge values.

## v3.2.7 checkpoint checks

- Decide whether v3.1.7 dashboard event efficiency should resume as v3.1.7 or continue as v3.2.x.
- Decide whether v3.1.8 retirement readiness should resume as v3.1.8 or continue as v3.2.x.
- Decide whether clusters are still worth adding later.
- Decide whether sensor overrides or inheritance remain unnecessary.
