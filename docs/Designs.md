# GoDriveLog Designs

**Pillar 1:** Design — what the system should do and why.

This document defines the structure and conventions for permanent design documentation.

Design documents describe intent, behaviour, constraints, interactions and rejected alternatives. They remain valuable whether implementation is complete, partial, superseded or never started.

## Conventions

- Store designs under `docs/Designs/`.
- Group documents by area.
- Use lowercase kebab-case filenames.
- After implementation, pair Implementation doc under `docs/Implementation/` at the same relative path to the Design doc.
- Record current progress only in [`Status.md`](Status.md).
- Do not put implementation status into design or implementation documents.

## Design and implementation pairing

```text
docs/Designs/<area>/<name>.md
docs/Implementation/<area>/<name>.md
```

## Design index

Refer to the index file [`Designs/Index.md`](Designs/Index.md)

## Design principles

Design documents describe **what** the system should do and **why**.

They should remain useful even if:

- implementation changes;
- implementation never occurs;
- implementation is replaced;
- implementation is removed.

Implementation details belong under `docs/Implementation/`.

Current implementation status belongs only in `docs/Status.md`.

## Canonical behaviour documentation

The gauge capability matrix belongs with design documentation because it describes conceptual applicability. Current implementation support belongs in [`Status.md`](Status.md).

## Current design documentation

Current design documentation is organised by system area under `docs/Designs/`.

Use the design index to locate the relevant current design documents:

- [`Designs/Index.md`](Designs/Index.md)

Gauge mechanism designs live under `docs/Designs/Gauge/`.
Dashboard-level design concerns live under `docs/Designs/Dashboard/`.

Current implementation state belongs in [`Status.md`](Status.md).
