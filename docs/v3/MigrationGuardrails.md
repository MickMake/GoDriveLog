# GoDriveLog v3 migration guardrails

Status: transition guidance  
Applies to: moving the current codebase toward the v3 docs  
References: `README.md`, `config.full.yaml`, `GoStructsConfig.md`, `ImplementationGuardrails.md`, `PerformanceGuardrails.md`

## 1. Purpose

This document explains how to move from the code that exists now to the v3 target without turning the repo into a compatibility swamp.

It is not the final v3 implementation design. It is the bridge between current working code and the intended v3 shape.

The current code may contain useful runtime pieces, renderer experiments, config loaders, asset handling, and OBD plumbing. Those pieces are not automatically wrong. They are also not automatically v3 just because they already exist.

## 2. Core migration rule

The target remains:

```text
vehicle endpoint
-> sensor polling runtime
-> sensor events
-> logs and dashboards as subscribers
```

Migration work should move code toward that model in small slices.

Do not bend the v3 docs around old code unless the old code reveals a real requirement. Convenience is not a requirement. A passing old test is not a requirement. A renderer goblin with a tiny clipboard is definitely not a requirement.

## 3. Current state versus target state

| Area | Current code may contain | v3 target |
|---|---|---|
| Config | earlier config structs/loaders | strict v3 root schema: `vehicles`, `sensors`, `assets`, `logs`, `dashboards` |
| Vehicle/OBD | existing reader/adapter plumbing | selected vehicle connects to an OBD-like endpoint address |
| Sensors | reader/state/cache-style concepts | sensor polling runtime emits sensor events |
| Logging | current JSONL writer behaviour | JSONL subscriber receives selected sensor events |
| Dashboard | current Fyne/dashboard renderer pieces | widget-driven dashboard subscriber |
| Assets | current asset experiments | repo-root-relative asset paths and asset families: digit, bar, frame, indicator, image |
| Performance | current display path may be slow | optimise locally without changing the v3 schema |

This table is not a complaint list. It is a migration map.

## 4. Migration adapters

Migration adapters are allowed at boundaries.

Migration behaviour must not leak into the v3 core model.

Allowed boundary adapters:

- a wrapper that exposes an existing OBD reader through the v3 endpoint/reader interface
- a temporary adapter that feeds existing sensor values into the new sensor event store
- a dashboard compatibility layer used only while replacing current renderer pieces
- a small bridge from existing asset loading into the v3 asset registry

Not allowed:

- making the v3 config loader accept undocumented shapes
- spreading compatibility branches across the runtime
- letting old renderer concepts define the v3 dashboard model
- letting the logger become the hidden scheduler
- making dashboards poll sensors because the current renderer wants values directly

## 5. Phase 0 — current display performance stabilisation

The current display path is known to need speed work before or during migration.

This is allowed, but it is tactical current-runtime work unless it directly supports the v3 target.

Allowed:

- profile current rendering hot spots
- cache decoded images/assets
- avoid re-decoding assets during render
- avoid rebuilding a full scene when one sensor value changes
- reduce unnecessary UI invalidation/redraw churn
- batch UI updates where the toolkit benefits from it
- keep temporary performance fixes local to the current renderer/runtime boundary
- capture useful performance lessons for the v3 renderer

Not allowed:

- changing the v3 config shape to match current renderer limitations
- turning current renderer concepts into v3 schema concepts by accident
- adding dashboard polling to solve renderer slowness
- adding global timing knobs outside the sensor event model without design review
- hiding stale/error/recovery transitions to improve frame times
- spreading temporary optimisation branches through the future v3 core

Exit criteria:

- current display is usable enough for dashboard iteration
- speed fixes are isolated
- no new schema concepts were added solely for current-renderer performance
- stale/error/recovery display behaviour remains visible
- lessons that apply to v3 are captured in `PerformanceGuardrails.md`

## 6. Phase 1 — freeze the v3 docs as the target

Goal: make docs the stable target before major code movement.

Inputs:

- `config.full.yaml`
- `config.example.yaml`
- all files under `docs/v3/examples/`
- `GoStructsConfig.md`
- `ImplementationGuardrails.md`
- `PerformanceGuardrails.md`

Rules:

- Treat documented v3 root sections as an allow-list.
- Unknown fields should fail at every documented level during v3 implementation.
- All active v3 examples should validate against the same schema rules.
- Asset paths are repository-root relative.
- Avoid schema churn once implementation starts.
- If implementation finds a real blocker, update docs first, then code.

Exit criteria:

- v3 schema shape is clear
- implementation order is documented
- migration guardrails are accepted
- performance constraints are acknowledged without warping the schema
- examples are schema-compliant, not just illustrative confetti

## 7. Phase 2 — v3 config structs and strict loading

Goal: load v3 config files strictly without needing the rest of the v3 runtime.

Allowed:

- add new v3 config structs beside current structs
- add strict YAML loading for v3 docs/examples
- reject unknown fields at all documented levels
- add validation for documented fields
- add tests for valid minimal, full, and standalone example configs
- add tests for unknown fields and bad references

Not allowed:

- silently accepting undocumented root fields
- silently accepting undocumented nested fields
- auto-converting current config files into v3 config behind the user's back
- expanding the schema just to satisfy old code paths

Exit criteria:

- `docs/v3/config.example.yaml` loads
- `docs/v3/config.full.yaml` loads
- all files under `docs/v3/examples/` load
- unknown fields fail at root and nested levels
- references validate
- current runtime can still build while v3 config work is staged

## 8. Phase 3 — vehicle endpoint abstraction

Goal: connect to real hardware or a bench simulator through endpoint addresses.

Allowed:

- wrap existing OBD adapter code behind a new endpoint/reader boundary
- add serial endpoint support
- add TCP endpoint support for bench/simulator work
- keep old reader code behind a compatibility boundary while replacing it

Not allowed:

- endpoint-provider branching throughout the core runtime
- leaking simulator identity into sensors, logs, dashboards, or asset config
- making config consumers care whether data came from real hardware or a simulator

If endpoint support is staged, the missing endpoint type must be tracked as a visible issue or migration task, not left as implied future work.

Exit criteria:

- selected vehicle resolves to one endpoint address
- endpoint address validates according to the v3 rules
- endpoint connector returns a reader usable by the sensor runtime
- serial and TCP endpoint shapes are both supported or explicitly tracked as staged work

## 9. Phase 4 — sensor event spine

Goal: build the central runtime path.

```text
endpoint reader -> sensor polling runtime -> sensor events -> state store
```

Allowed:

- implement sensor event type
- implement sensor status type
- implement latest-state store
- implement per-sensor polling based on `sensors.<id>.poll`
- implement initial stale rule: `max(poll * 3, 1000ms)`
- emit events for first read, value change, status change, stale/error/recovery transitions

Not allowed:

- one polling loop per consumer
- dashboards reading directly from endpoint/reader code
- logger reading directly from endpoint/reader code
- encoding errors as numeric values
- hiding stale/error/recovery transitions behind coalescing

Exit criteria:

- sensor polling works without dashboard or logger attached
- state store updates from events
- timestamps preserve original read time
- tests prove first/value/status/stale/recovery events

## 10. Phase 5 — JSONL subscriber

Goal: attach logging to the event spine.

Allowed:

- implement JSONL log subscriber
- select logged sensors from `logs.<id>.sensors`
- write first readings, value changes, and status changes
- include read timestamp and status

Not allowed:

- logger polling sensors
- logger owning sensor cadence
- logger causing extra endpoint reads

Exit criteria:

- JSONL output is driven only by sensor events
- unchanged duplicate values do not spam logs
- status changes are logged clearly

## 11. Phase 6 — asset registry

Goal: validate and load assets before widget rendering.

Allowed:

- implement asset family structs and registries
- resolve asset paths as repository-root relative
- validate asset references
- validate required indicator states
- validate bar set `off` cells
- validate bar zone cells and ordering
- validate frame ranges
- validate related image dimensions where required
- cache loaded/decoded assets

Not allowed:

- asset config rules/scripts
- per-widget asset decoding in the hot render path
- nil fallback behaviour that hides missing assets
- multiple active path dialects

Exit criteria:

- config references resolve to asset definitions
- missing assets fail clearly
- decoded assets can be reused by widgets/renderers
- active examples use the same path convention

## 12. Phase 7 — smallest dashboard subscriber

Goal: prove a dashboard can render from sensor state without polling sensors.

Start with:

- `image`
- `digit_display`
- `indicator`

Allowed:

- render a static panel image
- render formatted numeric values from latest sensor state
- render indicator states using sensor value and status
- update only changed widgets where practical

Not allowed:

- full scene rebuild for every sensor tick if avoidable
- direct endpoint reads from widgets
- conditions/scripts/formulas in dashboard YAML
- dashboard-level polling cadence

Exit criteria:

- one dashboard displays image + digit_display + indicator
- stale/error/missing states are visible
- dashboard receives events/state from runtime, not endpoint code
- digit display slot/decimal behaviour follows the documented rules

## 13. Phase 8 — richer dashboard widgets

Goal: add the remaining v3 widget types after the event path works.

Order:

1. `bar_display`
2. `frame_gauge`
3. any additional indicator/status variants that are already supported by the asset model

Rules:

- `bar_display` maps one sensor value to cells.
- `frame_gauge` maps one sensor value to a frame sequence.
- Fancy curves and sweeps use frame sets, not YAML geometry.
- Keep widget behaviour in code, not config rules.
- Bar zone behaviour follows `GoStructsConfig.md`.

Exit criteria:

- full examples can be represented by supported widget types
- renderer remains event/state-driven
- no hidden polling is introduced

## 14. Phase 9 — retire replaced current paths

Goal: remove or archive current paths after v3 paths replace them.

Allowed:

- delete old code once v3 tests cover the behaviour
- move old docs to archive
- remove compatibility bridges after consumers move
- keep small stable adapters only where they are genuine boundaries

Not allowed:

- leaving two active ways to configure the same thing
- keeping compatibility paths with no owner or removal plan
- letting old/current tests keep obsolete behaviour alive without review

Exit criteria:

- v3 config is the active config path
- sensor runtime owns polling
- logs and dashboards are subscribers
- old runtime pieces are either removed, archived, or clearly isolated

## 15. Migration decision rules

When touching existing code, ask:

1. Is this current-runtime stabilisation or v3 migration?
2. Does this move code toward the documented v3 shape?
3. Is this a boundary adapter or a leak into the core model?
4. Does this preserve dashboard/log subscriber boundaries?
5. Does this add schema surface area?
6. Is the performance fix local, measurable, and removable?
7. Would this change make future implementation simpler or just make today easier?

If the answer is fuzzy, write a note or open an issue before changing the shape.

## 16. Compatibility policy

Compatibility is temporary unless explicitly promoted.

Allowed compatibility:

- local wrappers around current code
- bridge packages at package boundaries
- tests that prove current behaviour while replacement is built
- clear TODOs tied to migration phases

Disallowed compatibility:

- undocumented config aliases
- global flags that choose old/new internals everywhere
- new schema fields whose only purpose is current-code convenience
- compatibility paths with no removal condition

## 17. Testing during migration

Every migration slice should add or preserve tests for the boundary it changes.

Useful test groups:

- v3 config loading and validation
- nested unknown field rejection
- ID naming and duplicate widget ID validation
- endpoint selection and reader creation
- sensor event emission
- stale and recovery transition emission
- log subscriber behaviour
- asset registry validation
- widget rendering state transitions
- current display performance baseline, where practical

Avoid giant end-to-end-only testing. Small boundary tests catch goblins before they learn teamwork.

## 18. Performance during migration

Performance matters, especially for the current display path.

The performance rule is:

```text
Optimise the current display path where needed, but do not let those fixes define the v3 architecture.
```

If a performance lesson applies to v3, document it in `PerformanceGuardrails.md`.

If a performance fix is temporary, keep it local and make the removal condition obvious.

## 19. Definition of done for migration docs

These migration guardrails are doing their job when contributors can answer:

- what the target shape is
- what the current code may still contain
- what can be wrapped temporarily
- what must not leak into v3 core
- how performance work fits without warping the schema
- what order to implement changes in
- how staged work is tracked
- why all active examples validate against one schema

If the docs create confusion, fix the docs before fixing the wrong code.
