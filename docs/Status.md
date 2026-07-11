# GoDriveLog Documentation Status

Pillar 2 - The "current code truth".

This file is the single source of truth for the current state of documented features.

Design documents define intent. Implementation records explain how designs became code. This register records the current state.

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
| Documentation | Documentation structure | [`Designs.md`](Designs.md) | [`Implementation.md`](Implementation.md) | In Progress | Core structure exists. Remaining work includes classification, status population, legacy-document retirement, naming cleanup and link validation. |
| Configuration | Documentation classification | [`Designs/Configuration/`](Designs/Configuration/) | [`Implementation/Configuration/`](Implementation/Configuration/) | In Progress | Area directories exist; documents and implementation records still require audit and classification. |
| Dashboard | Documentation classification | [`Designs/Dashboard/`](Designs/Dashboard/) | [`Implementation/Dashboard/`](Implementation/Dashboard/) | In Progress | Area directories exist; documents and implementation records still require audit and classification. |
| Logging | Documentation classification | [`Designs/Logging/`](Designs/Logging/) | [`Implementation/Logging/`](Implementation/Logging/) | Partial | Several design documents and one implementation boundary record have been classified. Pairing and feature status still require audit. |
| Realism Behaviour | Documentation classification | [`Designs/RealismBehaviour/realism-behaviour-guide.md`](Designs/RealismBehaviour/realism-behaviour-guide.md) | [`Implementation/RealismBehaviour/`](Implementation/RealismBehaviour/) | Partial | Canonical design guide and behaviour documents are classified. Implementation records and verified feature states still require audit. |
| Runtime | Documentation classification | [`Designs/Runtime/`](Designs/Runtime/) | [`Implementation/Runtime/`](Implementation/Runtime/) | Partial | Runtime design documents have been classified. Matching implementation records and feature states still require audit. |

## Rules

- Update this register when feature state changes.
- Do not duplicate implementation status in design documents.
- Link both design and implementation records where they exist.
- Use one row per significant documented feature once its state has been verified.
- Preserve superseded and rejected designs; record why their status changed.
- Do not mark a feature Implemented solely because code or a design document exists; verify the implementation against the design.
