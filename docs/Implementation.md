# GoDriveLog Implementation Records

**Pillar 3:** Implementation — how designs became code.

This document defines the structure and conventions for implementation records.

Implementation records describe packages touched, implementation approach, important trade-offs, deviations from design, limitations and relevant tests.

Implementation documents describe only behaviour verified in the current codebase. They must never document intended behaviour, proposed designs or assumed implementations. Where behaviour cannot be verified from code, the document must explicitly state that it could not be verified.

## Conventions

- Store records under `docs/Implementation/`.
- Match the relative path of the corresponding design document.
- Use lowercase kebab-case filenames.
- Link back to the paired design.
- Record current progress only in [`Status.md`](Status.md).

## Design and implementation pairing

```text
docs/Designs/<area>/<name>.md
docs/Implementation/<area>/<name>.md
```

## Implementation index

Refer to the index file [`Implementation/Index.md`](Implementation/Index.md)

## Implementation principles

Implementation records describe **how** a design became code.

They should record:

- packages touched;
- architectural decisions;
- implementation approach;
- trade-offs;
- deviations from the original design;
- limitations;
- important tests;
- follow-up work.

Every significant implementation should have a matching design document using the same relative path.

Example:

```text
docs/Designs/Logging/dashboard-jsonl-replay.md
docs/Implementation/Logging/dashboard-jsonl-replay.md
```

Implementation records should not redefine design intent.

Current implementation state belongs only in `docs/Status.md`.
