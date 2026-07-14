# v3.7 Release Plan

Status: planned

## Theme

Odometer and image-based numeric display realism.

All behaviour, configuration, constraints and non-goals are defined only in `docs/Designs/`.

## Release scope

| Slice | Name | Design |
|---|---|---|
| v3.7.0 | Release planning docs | Documentation-only |
| v3.7.1 | Odometer backlash | `../Designs/RealismBehaviour/odometer-backlash.md` |
| v3.7.2 | Per-digit response lag | `../Designs/RealismBehaviour/per-digit-response-lag.md` |
| v3.7.3 | Leading-zero behaviour | `../Designs/RealismBehaviour/numeric-leading-zero-behaviour.md` |
| v3.7.4 | Segment and digit bleed | `../Designs/RealismBehaviour/segment-bleed-digit-bleed.md` |
| v3.7.5 | Ghosting | `../Designs/RealismBehaviour/numeric-ghosting.md` |
| v3.7.6 | Uneven brightness | `../Designs/RealismBehaviour/uneven-brightness.md` |
| v3.7.7 | Load sag | `../Designs/RealismBehaviour/numeric-load-sag.md` |
| v3.7.8 | Tests, previews and docs checkpoint | Relevant v3.7 Design documents |

## Release boundaries

- One slice per branch and pull request.
- Do not implement later slices early.
- Do not add unrelated realism features.
- Do not redesign existing renderers.
- Any design change must be made in `docs/Designs/` before implementation.

## Completion

v3.7 is complete when every slice is implemented, tested, documented and marked complete in `ImplementationState.md`.
