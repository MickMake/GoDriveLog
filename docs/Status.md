# GoDriveLog Documentation Status

This file is the single source of truth for the current state of documented features.

## Allowed states

- **Designed** — design intent is documented; implementation has not started.
- **In Progress** — implementation work is actively underway.
- **Partial** — some of the design is implemented, with known remaining work.
- **Implemented** — the documented design is implemented.
- **Superseded** — replaced by a newer design or implementation.
- **Rejected** — deliberately not proceeding.

## Status register

| Area | Feature | Design | Implementation | State | Notes |
|---|---|---|---|---|---|
| Documentation | Documentation structure | [`Designs.md`](Designs.md) | [`Implementation.md`](Implementation.md) | In Progress | Structure created; existing documents still require classification and migration. |

## Rules

- Update this register when feature state changes.
- Do not duplicate implementation status in design documents.
- Link both design and implementation records where they exist.
- Preserve superseded and rejected designs; record why their status changed.
