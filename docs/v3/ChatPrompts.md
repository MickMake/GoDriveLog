# GoDriveLog v3 chat prompts

Status: workflow guidance  
Applies to: using ChatGPT chats to implement and verify the v3 migration  
References: `MigrationState.md`, `MigrationGuardrails.md`, `ImplementationGuardrails.md`, `DirectoryStructure.md`, `GoStructsConfig.md`, `PerformanceGuardrails.md`

## 1. Purpose

This file contains the reusable prompts for the GoDriveLog v3 migration workflow.

The goal is simple phone-friendly operation:

```text
Implementation chat creates one PR.
Verification chat checks that PR.
Mick merges only after PASS.
The next implementation chat uses the next prompt.
```

Do not make Mick reconstruct the process from memory.

## 2. How Mick should use this file

For an implementation chat, paste a short instruction like:

```text
Hey mate. You are the implementation chat for GoDriveLog v3. Use the v3.0.0 implementation prompt from docs/v3/ChatPrompts.md.
```

For a verification chat, paste:

```text
Hey mate. You are the verification chat for GoDriveLog v3. Use the verifier prompt from docs/v3/ChatPrompts.md for PR <number>.
```

That is enough. The chat must read the repo files listed in the prompt and work from repo state, not from chat memory.

## 3. Universal rules for all chats

All implementation and verification chats must:

1. Read `docs/v3/MigrationState.md` first.
2. Read `docs/v3/MigrationGuardrails.md`.
3. Treat repo docs as source of truth.
4. Use the target version and branch prefix from the relevant prompt and migration state.
5. Never rely on previous chat memory.
6. Never merge their own PR.
7. Never expand the v3 schema for current-code convenience.
8. Never add undocumented config aliases.
9. Keep changes coherent, reviewable, and scoped to one migration slice.
10. Say clearly when something is blocked or uncertain.

If `MigrationState.md` and this file disagree about the current position, stop and report the conflict. Do not guess.

## 4. Implementation prompt: v3.0.0 — working-code inventory and seam plan

Use this prompt for the first real migration implementation slice.

```text
You are the implementation chat for GoDriveLog v3.0.0.

Repository: MickMake/GoDriveLog

Target version: v3.0.0
Target stage: working-code inventory and seam plan
Expected branch prefix: v3.0.0
Suggested branch: v3.0.0-working-code-inventory

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/DirectoryStructure.md

Task:
Create a docs-only working-code inventory and seam plan.

Inspect the current repository and map existing code to v3 roles.

Cover at least:
- config loading and config structs
- runtime startup and command flow
- OBD, ELM327, vehicle, endpoint, or adapter code
- sensor polling, state, cache, and status logic
- logging and JSONL writer behaviour
- dashboard, renderer, Fyne, display, and widget code
- asset loading and image handling
- tests that prove behaviour worth preserving

For each area, decide one of:
- reuse
- refactor
- wrap
- replace
- archive

For each decision, explain:
- why
- what v3 role it maps to
- what seam or boundary should protect it
- what must not leak into v3
- what tests should preserve useful behaviour

Do not implement runtime code.
Do not change v3 config schema.
Do not start v3.0.1 work.
Do not delete current code.

Before editing, report:
1. files read
2. branch name
3. exact docs to create or update
4. scope
5. non-goals
6. docs-only justification

Then wait for Mick's approval before editing.

After approval:
1. create the branch from latest main
2. add/update the inventory/seam-plan docs
3. update docs/v3/MigrationState.md if the PR enters review or state changes
4. open a PR to main
5. do not merge
```

## 5. Implementation prompt: v3.0.1 — frozen v3 docs and schema target

```text
You are the implementation chat for GoDriveLog v3.0.1.

Repository: MickMake/GoDriveLog

Target version: v3.0.1
Target stage: frozen v3 docs and schema target
Expected branch prefix: v3.0.1
Suggested branch: v3.0.1-freeze-v3-docs-schema

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/GoStructsConfig.md
5. docs/v3/config.example.yaml
6. docs/v3/config.full.yaml
7. docs/v3/examples/

Task:
Review and tighten the v3 docs so they are a stable implementation target.

Confirm:
- documented root sections are the schema allow-list
- vehicles select logs and dashboards by ID
- sensors and assets are global catalogues
- active examples use the same schema rules
- asset paths are repository-root relative
- implementation blockers are documented before code work begins

Do not implement loader/runtime code.
Do not add new schema features unless a real blocker is found and explained.
Do not preserve old config shapes for convenience.

Before editing, report files read, branch name, scope, non-goals, and expected docs changes. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 6. Implementation prompt: v3.0.2 — strict v3 config load and validation

```text
You are the implementation chat for GoDriveLog v3.0.2.

Repository: MickMake/GoDriveLog

Target version: v3.0.2
Target stage: strict v3 config load and validation
Expected branch prefix: v3.0.2
Suggested branch: v3.0.2-config-loader-validation

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/GoStructsConfig.md
5. docs/v3/config.example.yaml
6. docs/v3/config.full.yaml
7. docs/v3/examples/

Task:
Implement strict loading and validation for documented v3 config files.

Implement only enough runtime-independent code to load and validate v3 config.

Required behaviour:
- documented root sections only
- nested unknown fields rejected
- docs/v3/config.example.yaml loads
- docs/v3/config.full.yaml loads
- docs/v3/examples/ load
- vehicle log/dashboard references validate
- sensors and assets remain global catalogues
- bad references fail clearly

Do not wire the full runtime.
Do not add compatibility aliases.
Do not auto-convert current config into v3 config.
Do not expand schema for old code paths.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 7. Implementation prompt: v3.0.3 — RuntimePlan resolution

```text
You are the implementation chat for GoDriveLog v3.0.3.

Repository: MickMake/GoDriveLog

Target version: v3.0.3
Target stage: RuntimePlan resolution
Expected branch prefix: v3.0.3
Suggested branch: v3.0.3-runtime-plan

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/DirectoryStructure.md
5. docs/v3/GoStructsConfig.md

Task:
Implement explicit resolution from loaded config plus selected vehicle ID to a RuntimePlan-style boundary.

Required behaviour:
- selected vehicle resolves explicitly
- endpoint config resolves
- selected log definitions resolve
- selected dashboard definitions resolve
- selected dashboard display collisions validate
- unselected dashboard definitions are inert for that runtime plan
- runtime packages do not need to walk raw config maps

Do not implement endpoint connectors.
Do not implement sensor polling.
Do not implement dashboard rendering.
Do not let logs or dashboards choose their own vehicle.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 8. Implementation prompt: v3.0.4 — endpoint abstraction with serial/TCP simulator support

```text
You are the implementation chat for GoDriveLog v3.0.4.

Repository: MickMake/GoDriveLog

Target version: v3.0.4
Target stage: endpoint abstraction with serial/TCP simulator support
Expected branch prefix: v3.0.4
Suggested branch: v3.0.4-endpoint-abstraction

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/DirectoryStructure.md
5. current OBD/vehicle/adapter code

Task:
Create the endpoint/reader boundary for real hardware and bench simulator work.

Required behaviour:
- selected vehicle resolves to one endpoint address
- endpoint address validates according to v3 rules
- connector returns a reader usable by the later sensor runtime
- serial endpoint shape is supported or explicitly tracked
- TCP endpoint shape for simulator/bench work is supported or explicitly tracked
- useful existing OBD adapter code is wrapped where practical

Do not leak simulator identity into sensors, logs, dashboards, or assets.
Do not add endpoint-provider branching throughout the core runtime.
Do not implement sensor polling yet.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 9. Implementation prompt: v3.0.5 — sensor event spine and latest-state store

```text
You are the implementation chat for GoDriveLog v3.0.5.

Repository: MickMake/GoDriveLog

Target version: v3.0.5
Target stage: sensor event spine and latest-state store
Expected branch prefix: v3.0.5
Suggested branch: v3.0.5-sensor-event-spine

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/PerformanceGuardrails.md
5. current sensor/state/cache code

Task:
Implement the central event path:
endpoint reader -> sensor polling runtime -> sensor events -> state store.

Required behaviour:
- sensor event type
- sensor status type
- latest-state store
- per-sensor polling based on sensors.<id>.poll
- initial stale rule: max(poll * 3, 1000ms)
- events for first read, value change, status change, stale/error/recovery transitions

Do not create one polling loop per consumer.
Do not let dashboards or loggers read directly from endpoint code.
Do not encode errors as numeric values.
Do not hide stale/error/recovery transitions behind coalescing.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 10. Implementation prompt: v3.0.6 — selected JSONL logging

```text
You are the implementation chat for GoDriveLog v3.0.6.

Repository: MickMake/GoDriveLog

Target version: v3.0.6
Target stage: selected JSONL logging
Expected branch prefix: v3.0.6
Suggested branch: v3.0.6-jsonl-subscriber

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. current logging/JSONL code

Task:
Attach selected JSONL logging to the sensor event spine.

Required behaviour:
- selected vehicle controls which logs run
- logs selected from vehicles.<id>.logs
- logged sensors selected from logs.<id>.sensors
- JSONL output is driven only by sensor events
- first readings, value changes, and status changes are written
- unchanged duplicate values do not spam logs
- read timestamp and status are included
- useful existing JSONL behaviour is reused where it fits the event-subscriber model

Do not let the logger poll sensors.
Do not let the logger own sensor cadence.
Do not let the logger cause extra endpoint reads.
Do not add vehicle-owned sensor lists.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 11. Implementation prompt: v3.0.7 — minimal asset registry

```text
You are the implementation chat for GoDriveLog v3.0.7.

Repository: MickMake/GoDriveLog

Target version: v3.0.7
Target stage: minimal asset registry
Expected branch prefix: v3.0.7
Suggested branch: v3.0.7-minimal-asset-registry

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/GoStructsConfig.md
5. current asset/image-loading code

Task:
Implement the asset families required by the first useful dashboard.

Start with:
- image
- digit / character image sets
- indicator

Required behaviour:
- asset paths resolve as repository-root relative
- image, digit, and indicator asset references resolve
- missing assets fail clearly
- required indicator states validate
- digit display characters validate where practical
- decoded assets can be reused by widgets/renderers
- useful current asset loading code is reused where it fits the registry model

Do not implement bar/frame assets yet.
Do not add asset config scripts/rules.
Do not decode assets per widget in the hot render path.
Do not hide missing assets with nil fallback behaviour.
Do not add vehicle-owned asset lists.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 12. Implementation prompt: v3.0.8 — smallest selected dashboard

```text
You are the implementation chat for GoDriveLog v3.0.8.

Repository: MickMake/GoDriveLog

Target version: v3.0.8
Target stage: smallest selected dashboard
Expected branch prefix: v3.0.8
Suggested branch: v3.0.8-smallest-dashboard

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/PerformanceGuardrails.md
5. current dashboard/renderer/Fyne code

Task:
Prove a selected dashboard can render from sensor state without polling sensors.

Start with widgets:
- image
- digit_display
- indicator

Required behaviour:
- selected vehicle controls which dashboards render
- selected dashboards come from vehicles.<id>.dashboards
- static panel image can render
- numeric values render from latest sensor state
- indicator states render using sensor value and status
- stale/error/missing states are visible
- dashboard receives events/state from runtime, not endpoint code
- digit display slot/decimal behaviour follows documented rules

Do not implement bar_display or frame_gauge yet.
Do not let widgets read directly from endpoint code.
Do not add dashboard-level polling cadence.
Do not add YAML scripts/formulas/conditions.
Do not rebuild the full scene for every sensor tick if avoidable.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 13. Implementation prompt: v3.0.9 — richer asset registry

```text
You are the implementation chat for GoDriveLog v3.0.9.

Repository: MickMake/GoDriveLog

Target version: v3.0.9
Target stage: richer asset registry
Expected branch prefix: v3.0.9
Suggested branch: v3.0.9-richer-asset-registry

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/GoStructsConfig.md
5. current asset/image-loading code

Task:
Add the richer asset families needed by later dashboard widgets.

Add:
- bar
- frame

Required behaviour:
- bar set off cells validate
- bar zone cells and ordering validate
- frame ranges validate
- related image dimensions validate where required
- decoded bar/frame assets can be cached and reused
- richer examples can resolve their referenced assets

Do not add YAML formulas.
Do not add widget scripts.
Do not add hidden geometry languages.
Do not pre-optimise visual tricks before the event/render path needs them.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 14. Implementation prompt: v3.0.10 — richer dashboard widgets

```text
You are the implementation chat for GoDriveLog v3.0.10.

Repository: MickMake/GoDriveLog

Target version: v3.0.10
Target stage: richer dashboard widgets
Expected branch prefix: v3.0.10
Suggested branch: v3.0.10-richer-dashboard-widgets

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. docs/v3/PerformanceGuardrails.md
5. docs/v3/GoStructsConfig.md

Task:
Add richer dashboard widgets after the event/state path works.

Order:
1. bar_display
2. frame_gauge
3. additional indicator/status variants already supported by the asset model

Required behaviour:
- bar_display maps one sensor value to cells
- frame_gauge maps one sensor value to a frame sequence
- fancy curves and sweeps use frame sets, not YAML geometry
- widget behaviour stays in code, not config rules
- renderer remains event/state-driven
- no hidden polling is introduced

Do not add YAML scripts/formulas.
Do not let widgets poll sensors or endpoints.
Do not change sensor cadence from dashboard config.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 15. Implementation prompt: v3.0.11 — retire or archive replaced current paths

```text
You are the implementation chat for GoDriveLog v3.0.11.

Repository: MickMake/GoDriveLog

Target version: v3.0.11
Target stage: retire or archive replaced current paths
Expected branch prefix: v3.0.11
Suggested branch: v3.0.11-retire-replaced-paths

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/MigrationGuardrails.md
3. docs/v3/ImplementationGuardrails.md
4. current code paths marked as replaced by earlier v3 work

Task:
Remove or archive current paths after v3 paths replace them and tests cover wanted behaviour.

Required behaviour:
- v3 config is the active config path
- selected vehicle owns endpoint/log/dashboard runtime profile
- sensor runtime owns polling
- logs and dashboards are subscribers
- old runtime pieces are removed, archived, or clearly isolated
- compatibility bridges are removed after consumers move

Do not leave two active ways to configure the same thing.
Do not keep compatibility paths with no owner or removal plan.
Do not let old/current tests preserve obsolete behaviour without review.

Before editing, report files read, branch name, scope, non-goals, expected tests. Then wait for Mick's approval.

After approval, open a PR and do not merge.
```

## 16. Generic verifier prompt — any stage

Use this prompt for every verification chat.

```text
You are the verification chat for GoDriveLog v3.

Repository: MickMake/GoDriveLog

PR to verify: <PR_NUMBER>

Read first from latest main:
1. docs/v3/MigrationState.md
2. docs/v3/ChatPrompts.md
3. docs/v3/MigrationGuardrails.md
4. docs/v3/ImplementationGuardrails.md

Then fetch PR <PR_NUMBER>, including metadata, changed files, and diff.

Determine the target version from:
1. PR branch name
2. MigrationState.md
3. the matching implementation prompt in ChatPrompts.md

If those disagree, report BLOCKED or FAIL. Do not guess.

Verify:
- PR branch starts with the target version
- PR scope matches the matching implementation prompt
- PR does not perform later-stage work without clear justification
- PR follows the seam-based migration rule
- PR preserves useful behaviour where appropriate without preserving old architecture by default
- PR does not expand the v3 schema for current-code convenience
- PR does not add undocumented config aliases
- PR updates MigrationState.md correctly if migration state changes
- tests are suitable, or docs-only justification is valid
- repo remains coherent and reviewable

Return one verdict:
PASS — ready to merge
PASS WITH NOTES — acceptable, but follow-up work should be tracked
FAIL — changes required before merge
BLOCKED — cannot verify due to missing information/tooling

Include:
1. verdict
2. target version
3. PR branch
4. files reviewed
5. checks passed
6. issues found
7. required fixes, if any
8. whether MigrationState.md should advance after merge
9. next target version/action if the PR passes

Do not merge the PR unless Mick explicitly asks.
```

## 17. Minimal phone commands

Implementation chat:

```text
Hey mate. You are the implementation chat for GoDriveLog v3. Use the v3.0.0 implementation prompt from docs/v3/ChatPrompts.md.
```

Verification chat:

```text
Hey mate. You are the verification chat for GoDriveLog v3. Use the verifier prompt from docs/v3/ChatPrompts.md for PR <number>.
```

Next implementation chat:

```text
Hey mate. You are the implementation chat for GoDriveLog v3. Use the next implementation prompt from docs/v3/ChatPrompts.md and docs/v3/MigrationState.md.
```

If a chat cannot determine which prompt to use, it must stop and ask for the target version. It must not freestyle.
