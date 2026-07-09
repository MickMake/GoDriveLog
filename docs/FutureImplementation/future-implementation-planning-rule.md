# Planning Rule

Origin: `docs/v3.7/PlannedFeatures.md`

Do not implement anything from this directory as part of an unrelated slice.

## Promotion rule

Promotion should be explicit:

1. choose one small candidate;
2. read or create the canonical behaviour definition under [`../RealismBehaviourGuide/`](../RealismBehaviourGuide/);
3. define its user-facing config;
4. define which gauge families support it;
5. add or update the FutureImplementation ticket for build planning;
6. add docs and prompt slice(s);
7. then implement it in a dedicated branch/PR.

## Boundary rule

```text
RealismBehaviourGuide = definition / behaviour / real-world simulation
FutureImplementation = implementation ticket / backlog / build plan
```

FutureImplementation files may link to behaviour guide pages, but should not become a second source of truth for realism behaviour definitions.
