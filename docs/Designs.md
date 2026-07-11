# GoDriveLog Designs

This is the index for permanent design documentation.

Design documents describe intent, behaviour, constraints, interactions and rejected alternatives. They remain valuable whether implementation is complete, partial, superseded or never started.

## Conventions

- Store designs under `docs/Designs/`.
- Group documents by area.
- Use lowercase kebab-case filenames.
- Pair significant designs with an implementation record at the same relative path under `docs/Implementation/`.
- Record current progress only in [`Status.md`](Status.md).
- Do not put implementation status into design documents.

## Design and implementation pairing

```text
docs/Designs/<area>/<name>.md
docs/Implementation/<area>/<name>.md
```

## Design index

| Area | Design | Purpose | Implementation record |
|---|---|---|---|
| _To classify_ | Existing documentation has not yet been migrated. | Preserve the current repository state while the documentation is reviewed. | _Not yet created_ |

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

[`Designs/RealismBehaviour/realism-behaviour-guide.md`](Designs/RealismBehaviour/realism-behaviour-guide.md) is the canonical definition of gauge realism behaviour.

The gauge capability matrix belongs with design documentation because it describes conceptual applicability. Current implementation support belongs in [`Status.md`](Status.md).

