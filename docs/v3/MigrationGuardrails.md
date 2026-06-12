# GoDriveLog v3 migration guardrails

Status: transition guidance  
Applies to: moving the current codebase toward the v3 docs  
References: `README.md`, `config.full.yaml`, `GoStructsConfig.md`, `ImplementationGuardrails.md`, `PerformanceGuardrails.md`

## 1. Purpose

This document explains how to move from the code that exists now to the v3 target without turning the repo into a compatibility swamp or throwing away working knowledge for the sake of a clean-looking rewrite.

It is not the final v3 implementation design. It is the bridge between current working code and the intended v3 shape.

The current code may contain useful runtime pieces, renderer experiments, config loaders, asset handling, logging behaviour, and OBD plumbing. Those pieces are not automatically wrong. They are also not automatically v3 just because they already exist.

The migration posture is:

```text
Refactor through seams.
Replace only after proving the replacement.
```

Do not preserve old architecture by default. Do preserve working behaviour and proven code where they fit behind v3 seams.

## 2. Core migration rule

The target is:

```text
selected vehicle
-> OBD endpoint
-> sensor polling runtime
-> sensor events
-> selected logs and dashboards as subscribers
```

The selected vehicle is the runtime profile. It chooses the OBD endpoint, log definitions, and dashboard definitions to run.

Sensors and assets remain global catalogues. Vehicles do not directly list sensors or assets.

Migration work should move code toward that model in small, versioned slices.

Do not bend the v3 docs around old code unless the old code reveals a real requirement. Convenience is not a requirement. A passing old test is not a requirement. A renderer goblin with a tiny clipboard is definitely not a requirement.

## 3. Migration release line

The v3.0.x line is the foundation and migration line.

Each v3.0.x version should leave the repository in a coherent, buildable, reviewable state. No v3.0.x slice should require the whole grand plan to be complete before it has value.

```text
v3.0.0  working-code inventory and seam plan
v3.0.1  frozen v3 docs and schema target
v3.0.2  strict v3 config load/validation
v3.0.3  RuntimePlan resolution
v3.0.4  endpoint abstraction with serial/TCP simulator support
v3.0.5  sensor event spine and latest-state store
v3.0.6  selected JSONL logging
v3.0.7  minimal asset registry: image, digit, indicator
v3.0.8  smallest selected dashboard: image + digit_display + indicator
v3.0.9  richer asset registry: bar and frame assets
v3.0.10 richer dashboard widgets: bar_display and frame_gauge
v3.0.11 retire or archive replaced current paths
```

Version rule:

```text
v3.0.x = building the v3 foundation
v3.1.0 = first useful v3 runtime
```

The expected v3.1.0 threshold is the first useful v3 runtime:

```text
selected vehicle
-> endpoint
-> sensors
-> JSONL
-> one dashboard with image + digit_display + indicator
```

## 4. Branch naming rule

Branches for v3 migration work should start with the target version number.

Examples:

```text
v3.0.0-docs-migration-seams
v3.0.2-config-loader-validation
v3.0.3-runtime-plan
v3.0.4-endpoint-abstraction
v3.0.8-smallest-dashboard
```

If a branch prepares or clarifies the versioned migration process itself, treat it as `v3.0.0` work.

If the target version is unclear, decide the target version before creating the branch. Versionless branches invite goblins with clipboards.

## 5. Current state versus target state

| Area | Current code may contain | v3 target |
|---|---|---|
| Config | earlier config structs/loaders | strict v3 root schema: `vehicles`, `sensors`, `assets`, `logs`, `dashboards` |
| Vehicle/OBD | existing reader/adapter plumbing | selected vehicle connects to an OBD-like endpoint and selects logs/dashboards |
| Sensors | reader/state/cache-style concepts | global sensor catalogue plus polling runtime emits sensor events |
| Logging | current JSONL writer behaviour | global log definitions selected by vehicle; log subscribers receive selected sensor events |
| Dashboard | current Fyne/dashboard renderer pieces | global dashboard definitions selected by vehicle; widget-driven dashboard subscribers |
| Assets | current asset experiments | global repo-root-relative asset catalogue with digit, bar, frame, indicator, and image families |
| Performance | current display path may be slow | optimise locally without changing the v3 schema |

This table is not a complaint list. It is a migration map.

## 6. Working-code inventory and seams

Before replacing current pieces, identify what exists and what should happen to it.

Every current subsystem should receive one decision:

| Existing thing | Default migration decision |
|---|---|
| Current config shape | replace or archive if it conflicts with v3 |
| OBD/ELM327 adapter code | reuse behind the new endpoint/reader seam where practical |
| JSONL writer behaviour | refactor or reuse as an event subscriber |
| Sensor polling/cache logic | refactor or replace depending on coupling |
| Dashboard renderer experiments | reuse ideas/code behind the new widget/renderer seam where practical |
| Asset loading code | refactor or reuse if it fits the repo-root asset registry |
| Old dashboard config model | archive or replace |
| Tests | keep if they prove behaviour still wanted; rewrite if they protect obsolete shape |

Rules:

- Do not start with deletion.
- Do not start with a complete rewrite.
- Map current code to v3 roles first.
- Keep working behaviour only when it still belongs in v3.
- Wrap useful code at boundaries before moving it into the v3 path.
- Do not promote old architecture just because old code already exists.

Useful summary:

```text
Salvage behaviour and proven implementation knowledge.
Do not preserve the old machine.
```

## 7. Migration adapters

Migration adapters are allowed at boundaries.

Migration behaviour must not leak into the v3 core model.

Allowed boundary adapters:

- a wrapper that exposes an existing OBD reader through the v3 endpoint/reader interface
- a temporary adapter that feeds existing sensor values into the new sensor event store
- a dashboard compatibility layer used only while replacing current renderer pieces
- a small bridge from existing asset loading into the v3 asset registry
- a wrapper around current JSONL behaviour if it can consume v3 sensor events cleanly

Not allowed:

- making the v3 config loader accept undocumented shapes
- spreading compatibility branches across the runtime
- letting old renderer concepts define the v3 dashboard model
- letting the logger become the hidden scheduler
- making dashboards poll sensors because the current renderer wants values directly
- making vehicles directly own sensor or asset definitions
- keeping compatibility paths with no owner or removal condition

## 8. v3.0.0 — working-code inventory and seam plan

Goal: identify which current pieces should be reused, refactored, wrapped, replaced, or archived before v3 implementation work begins.

Allowed:

- inspect current config, runtime, OBD, logging, dashboard, renderer, and asset code
- map existing code to v3 roles
- identify behaviour that should survive into v3
- identify code that can be wrapped behind a v3 seam
- identify old shapes that should be archived rather than carried forward
- choose the first implementation seam

Not allowed:

- promoting old config shape into v3 by accident
- rewriting major runtime pieces before seams are clear
- deleting working code before its useful behaviour has been reviewed
- treating current runtime architecture as the v3 target

Exit criteria:

- current config, runtime, OBD, logging, dashboard, renderer, and asset code are mapped to v3 roles
- each subsystem has a reuse/refactor/replace/archive decision
- useful working behaviour is identified
- old architecture is not automatically promoted into v3
- first implementation seam is chosen

## 9. v3.0.1 — frozen v3 docs and schema target

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
- Vehicles select logs and dashboards by ID.
- Sensors and assets are global catalogues.
- All active v3 examples should validate against the same schema rules.
- Asset paths are repository-root relative.
- Avoid schema churn once implementation starts.
- If implementation finds a real blocker, update docs first, then code.

Exit criteria:

- v3 schema shape is clear
- vehicle runtime-profile ownership is clear
- migration release line is documented
- implementation order is documented
- migration guardrails are accepted
- performance constraints are acknowledged without warping the schema
- examples are schema-compliant, not just illustrative confetti

## 10. v3.0.2 — strict v3 config load and validation

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
- vehicles resolve their selected log IDs
- vehicles resolve their selected dashboard IDs
- sensors and assets remain global catalogues
- unknown fields fail at root and nested levels
- references validate
- repository still builds while v3 config work is staged

## 11. v3.0.3 — RuntimePlan resolution

Goal: turn a loaded config plus selected vehicle ID into an explicit runtime plan.

Suggested boundary:

```text
Resolve(config, selectedVehicleID) -> RuntimePlan
```

The exact Go shape may change during implementation. The important point is that selected-vehicle resolution becomes explicit before runtime wiring spreads through the code.

Allowed:

- resolve the selected vehicle
- resolve the endpoint config
- resolve selected log definitions
- resolve selected dashboard definitions
- validate selected dashboard display collisions
- resolve referenced sensors and assets where practical

Not allowed:

- making runtime packages repeatedly walk raw config maps
- letting dashboards or logs choose their own vehicle
- letting unselected dashboard definitions behave as active runtime participants

Exit criteria:

- multiple vehicles require explicit runtime vehicle selection
- selected vehicle resolves to exactly one runtime profile
- selected logs and dashboards are explicit
- display collision validation applies only to selected dashboards
- runtime packages can receive a resolved plan instead of raw config sprawl

## 12. v3.0.4 — endpoint abstraction with serial/TCP simulator support

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

Exit criteria:

- selected vehicle resolves to one endpoint address
- endpoint address validates according to the v3 rules
- endpoint connector returns a reader usable by the sensor runtime
- serial and TCP endpoint shapes are both supported or explicitly tracked as staged work

## 13. v3.0.5 — sensor event spine and latest-state store

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

## 14. v3.0.6 — selected JSONL logging

Goal: attach selected logging to the event spine.

Allowed:

- implement JSONL log subscriber
- select log definitions from `vehicles.<id>.logs`
- select logged sensors from `logs.<id>.sensors`
- write first readings, value changes, and status changes
- include read timestamp and status
- reuse current JSONL behaviour where it fits the event-subscriber model

Not allowed:

- logger polling sensors
- logger owning sensor cadence
- logger causing extra endpoint reads
- vehicle-owned sensor lists

Exit criteria:

- selected vehicle controls which logs run
- JSONL output is driven only by sensor events
- unchanged duplicate values do not spam logs
- status changes are logged clearly

## 15. v3.0.7 — minimal asset registry

Goal: validate and load the asset families required by the first useful dashboard.

Start with:

- `image`
- `digit` / character image sets
- `indicator`

Allowed:

- implement minimal asset family structs and registries
- resolve asset paths as repository-root relative
- validate asset references used by the smallest dashboard
- validate required indicator states
- validate digit display characters where practical
- cache loaded/decoded assets
- reuse current asset loading code where it fits the registry model

Not allowed:

- asset config rules/scripts
- per-widget asset decoding in the hot render path
- nil fallback behaviour that hides missing assets
- multiple active path dialects
- vehicle-owned asset lists
- implementing bar/frame asset support before the smallest dashboard needs it

Exit criteria:

- image, digit, and indicator asset references resolve
- missing assets fail clearly
- decoded assets can be reused by widgets/renderers
- active minimal examples use the same path convention

## 16. v3.0.8 — smallest selected dashboard

Goal: prove a selected dashboard can render from sensor state without polling sensors.

Start with:

- `image`
- `digit_display`
- `indicator`

Allowed:

- select dashboard definitions from `vehicles.<id>.dashboards`
- render a static panel image
- render formatted numeric values from latest sensor state
- render indicator states using sensor value and status
- update only changed widgets where practical
- reuse current renderer experiments where they fit the widget/renderer seam

Not allowed:

- full scene rebuild for every sensor tick if avoidable
- direct endpoint reads from widgets
- conditions/scripts/formulas in dashboard YAML
- dashboard-level polling cadence
- vehicle-owned asset lists

Exit criteria:

- selected vehicle controls which dashboards render
- one dashboard displays image + digit_display + indicator
- stale/error/missing states are visible
- dashboard receives events/state from runtime, not endpoint code
- digit display slot/decimal behaviour follows the documented rules

## 17. v3.0.9 — richer asset registry

Goal: add asset families needed by the richer dashboard widgets after the smallest dashboard path works.

Add:

- `bar`
- `frame`

Allowed:

- validate bar set `off` cells
- validate bar zone cells and ordering
- validate frame ranges
- validate related image dimensions where required
- cache decoded bar/frame assets

Not allowed:

- YAML formulas
- widget scripts
- hidden geometry languages
- pre-optimising every visual trick before the basic event path works

Exit criteria:

- bar and frame assets validate
- missing/invalid bar and frame assets fail clearly
- richer examples can resolve their referenced assets

## 18. v3.0.10 — richer dashboard widgets

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

## 19. v3.0.11 — retire or archive replaced current paths

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
- selected vehicle owns endpoint/log/dashboard runtime profile
- sensor runtime owns polling
- logs and dashboards are subscribers
- old runtime pieces are either removed, archived, or clearly isolated

## 20. Migration decision rules

When touching existing code, ask:

1. Which v3.0.x target version does this belong to?
2. Does the branch name start with that target version?
3. Is this working-code inventory, current-runtime stabilisation, or v3 migration?
4. Does this move code toward the documented v3 shape?
5. Is this a boundary adapter or a leak into the core model?
6. Does this preserve dashboard/log subscriber boundaries?
7. Does this add schema surface area?
8. Is the performance fix local, measurable, and removable?
9. Would this change make future implementation simpler or just make today easier?

If the answer is fuzzy, write a note or open an issue before changing the shape.

## 21. Compatibility policy

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

## 22. Testing during migration

Every migration slice should add or preserve tests for the boundary it changes.

Useful test groups:

- v3 config loading and validation
- nested unknown field rejection
- ID naming and duplicate widget ID validation
- endpoint selection and reader creation
- vehicle log/dashboard reference validation
- selected vehicle dashboard display collision validation
- RuntimePlan resolution
- sensor event emission
- stale and recovery transition emission
- log subscriber behaviour
- asset registry validation
- widget rendering state transitions
- current display performance baseline, where practical

Avoid giant end-to-end-only testing. Small boundary tests catch goblins before they learn teamwork.

## 23. Performance during migration

Performance matters, especially for the display path.

The performance rule is:

```text
Optimise the current display path where needed, but do not let those fixes define the v3 architecture.
```

If a performance lesson applies to v3, document it in `PerformanceGuardrails.md`.

If a performance fix is temporary, keep it local and make the removal condition obvious.

## 24. Definition of done for migration docs

These migration guardrails are doing their job when contributors can answer:

- what the target shape is
- what versioned migration slice they are working on
- what branch naming convention applies
- what the current code may still contain
- what can be wrapped temporarily
- what must not leak into v3 core
- how performance work fits without warping the schema
- what order to implement changes in
- how staged work is tracked
- why vehicles select logs and dashboards
- why sensors and assets remain global catalogues
- why all active examples validate against one schema

If the docs create confusion, fix the docs before fixing the wrong code.
