# Gauge Realism Map

Origin: `docs/v3.7/PlannedFeatures.md`

Status: historical planning note / superseded by canonical behaviour guide

## Current source of truth

Gauge realism behaviour definitions and current option status live in [`../RealismBehaviourGuide/`](../RealismBehaviourGuide/).

Do not use this historical map as implementation truth. It is retained only to preserve the v3.7 planning extraction trail.

## Implementation planning rule

When implementing a realism option:

1. Read the matching behaviour definition in [`../RealismBehaviourGuide/`](../RealismBehaviourGuide/).
2. Check current code and completed release docs.
3. Create or use a specific FutureImplementation ticket for the implementation work.
4. Do not copy behaviour definitions back into this file.

## Historical note

The original map mixed implemented options, not-implemented options, and possible candidates in one table. That made it too easy for FutureImplementation to become a second source of truth. The canonical map is now the Realism Behaviour Guide index.
