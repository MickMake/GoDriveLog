# GoDriveLog Implementation Records

This is the index for documentation explaining how designs became code.

Implementation records describe packages touched, implementation approach, important trade-offs, deviations from design, limitations and relevant tests.

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

| Area | Purpose | Location |
|---|---|---|
| Configuration | Configuration format, loading and validation | [`Designs/Configuration/`](Designs/Configuration/) |
| Dashboard | Dashboard architecture, rendering and composition | [`Designs/Dashboard/`](Designs/Dashboard/) |
| Logging | Logging, JSONL replay and event models | [`Designs/Logging/`](Designs/Logging/) |
| Realism Behaviour | Canonical gauge realism behaviour definitions | [`Designs/RealismBehaviour/realism-behaviour-guide.md`](Designs/RealismBehaviour/realism-behaviour-guide.md) |
| Runtime | Runtime architecture, kiosk mode and MQTT | [`Designs/Runtime/`](Designs/Runtime/) |

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

