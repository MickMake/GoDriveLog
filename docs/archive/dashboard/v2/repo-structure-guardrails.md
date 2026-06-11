# Dashboard v2 Repo Structure Guardrails

## Purpose

This file defines where GoDriveLog dashboard v2 code, config, docs, assets, and test data should live.

The goal is to stop future implementation work from scattering files across the repo like someone sneezed during a scaffolding delivery.

These rules apply to all dashboard v2 implementation stages unless explicitly overridden by Mick.

---

## Core rule

```text
No new top-level directories without explicit approval.
```

If a branch needs a new top-level directory, stop and ask first.

Approved top-level directories are:

```text
cmd/
internal/
assets/
configs/
docs/
testdata/
```

Existing Go module files stay at the repo root:

```text
go.mod
go.sum
README.md
config.example.yaml
```

---

## Expected final shape

Target structure after the dashboard v2 rewrite:

```text
GoDriveLog/
  cmd/
    GoDriveLog/
      main.go

  internal/
    config/
      config.go
      sensors.go
      dashboard.go
      validation.go

    logger/
      jsonl.go

    sensors/
      reader.go
      mock_reader.go
      elmobd_reader.go
      state.go
      state_store.go

    dashboard/
      assets/
        registry.go
        image.go
        frameset.go
        charset.go

      decoders/
        registry.go
        normalize.go
        threshold.go
        frame_index.go
        format_number.go
        digits.go
        boolean.go

      scene/
        scene.go
        element.go
        condition.go
        block.go
        layout.go

      renderer/
        fyne/
          renderer.go
          image.go
          sprite_frame.go
          sprite_text.go
          group.go

  assets/
    dashboard/
      bttf/
        background.png
        overlays/
          rpm_box.png
          rpm_box_glow.png
        digits/
          yellow/
            0.png
            1.png
            2.png
            3.png
            4.png
            5.png
            6.png
            7.png
            8.png
            9.png
          green/
            0.png
            1.png
            2.png
            3.png
            4.png
            5.png
            6.png
            7.png
            8.png
            9.png
        throttle/
          frame_000.png
          frame_001.png
          frame_002.png

  configs/
    examples/
      dashboard-minimal.yaml
      dashboard-bttf-sprite.yaml
      dashboard-decoder-demo.yaml

  docs/
    dashboard/
      v2/
        README.md
        overview.md
        prompts.md
        reference.md
        repo-structure-guardrails.md
        current-status.md
        decisions.md
        human-process.md
        CHANGES.md

  testdata/
    dashboard/
      assets/
        tiny.png
        digits/
        frames/
      configs/
        valid-minimal.yaml
        invalid-missing-asset.yaml
        invalid-decoder-ref.yaml

  go.mod
  go.sum
  README.md
  config.example.yaml
```

This is the fence. Do not let the software goats eat the neighbour's laundry.

---

## Directory responsibilities

## `cmd/`

Application entry points only.

Allowed:

```text
cmd/GoDriveLog/main.go
```

Rules:

- Keep command startup here.
- Do not put dashboard engine logic here.
- Do not put config validation here.
- Do not put asset loading here.
- Do not put decoder logic here.

`main.go` should wire pieces together. It should not become the junk drawer of civilisation.

---

## `internal/config/`

YAML config structures, loading, defaults, and validation.

Allowed examples:

```text
internal/config/config.go
internal/config/sensors.go
internal/config/dashboard.go
internal/config/validation.go
```

Responsibilities:

- top-level config struct
- sensor config structs
- dashboard config structs
- config loading
- config validation
- config defaults
- active sensor extraction

Not allowed:

- Fyne rendering
- image loading
- decoder execution
- OBD reading
- JSONL logging
- dashboard drawing

If it needs Fyne imports, it does **not** belong here.

---

## `internal/sensors/`

Runtime sensor reading and latest state.

Allowed examples:

```text
internal/sensors/reader.go
internal/sensors/mock_reader.go
internal/sensors/elmobd_reader.go
internal/sensors/state.go
internal/sensors/state_store.go
```

Responsibilities:

- reader interfaces
- mock reader
- ELM327/OBD reader
- sensor reading model
- latest sensor state
- stale/error state
- state store

Not allowed:

- dashboard rendering
- scene elements
- image loading
- visual decoder logic
- layout logic

Sensors know values. They do not know what those values look like.

---

## `internal/logger/`

Logging only.

Allowed examples:

```text
internal/logger/jsonl.go
```

Responsibilities:

- JSONL log writing
- log rotation
- log path management

Not allowed:

- dashboard state logic
- UI updates
- visual formatting
- config schema changes unrelated to logging

Do not touch logger behaviour during dashboard stages unless compilation requires a tiny adjustment.

---

## `internal/dashboard/`

Dashboard v2 engine.

This is the home for dashboard-specific runtime code.

Approved subdirectories:

```text
internal/dashboard/assets/
internal/dashboard/decoders/
internal/dashboard/scene/
internal/dashboard/renderer/
```

Do not create sibling packages such as:

```text
internal/newdashboard/
internal/display2/
internal/visuals/
internal/engine/
internal/ui2/
```

without explicit approval.

That way lies `final_final_really_final.go`.

---

## `internal/dashboard/assets/`

Asset registry and asset loading.

Allowed examples:

```text
internal/dashboard/assets/registry.go
internal/dashboard/assets/image.go
internal/dashboard/assets/frameset.go
internal/dashboard/assets/charset.go
```

Responsibilities:

- image asset definitions
- frame set asset definitions
- charset asset definitions
- path resolution
- asset validation
- asset caching

Not allowed:

- decoder logic
- scene layout logic
- Fyne scene orchestration
- sensor polling

Asset loading may use Fyne/image primitives if needed, but keep it limited and boring.

---

## `internal/dashboard/decoders/`

Value-to-output logic.

Allowed examples:

```text
internal/dashboard/decoders/registry.go
internal/dashboard/decoders/normalize.go
internal/dashboard/decoders/threshold.go
internal/dashboard/decoders/frame_index.go
internal/dashboard/decoders/format_number.go
internal/dashboard/decoders/digits.go
internal/dashboard/decoders/boolean.go
```

Responsibilities:

- normalize values
- threshold states
- frame index selection
- number formatting
- digit extraction
- boolean conditions
- decoder registry
- decoder execution

Not allowed:

- image loading
- Fyne rendering
- OBD polling
- file system asset walking except where explicitly needed for tests

Decoders produce values. They do not draw.

---

## `internal/dashboard/scene/`

Renderer-neutral dashboard scene model.

Allowed examples:

```text
internal/dashboard/scene/scene.go
internal/dashboard/scene/element.go
internal/dashboard/scene/condition.go
internal/dashboard/scene/block.go
internal/dashboard/scene/layout.go
```

Responsibilities:

- scene model
- element model
- layer/z-order model
- groups
- reusable blocks
- conditions
- bindings
- renderer-neutral scene evaluation

Not allowed:

- Fyne-specific drawing
- OBD reading
- JSONL logging
- raw image loading unless represented as asset references

Scene code describes what should exist. It does not know how Fyne draws it.

---

## `internal/dashboard/renderer/`

Renderer implementations.

Approved renderer:

```text
internal/dashboard/renderer/fyne/
```

Allowed examples:

```text
internal/dashboard/renderer/fyne/renderer.go
internal/dashboard/renderer/fyne/image.go
internal/dashboard/renderer/fyne/sprite_frame.go
internal/dashboard/renderer/fyne/sprite_text.go
internal/dashboard/renderer/fyne/group.go
```

Responsibilities:

- Fyne-specific drawing
- converting scene elements to Fyne canvas objects
- updating displayed images/text/visibility
- z-order presentation
- handling renderer lifecycle

Not allowed:

- OBD polling
- config loading
- decoder definitions
- asset path validation
- logger changes

The renderer draws resolved scene state. It should not become the project oracle.

---

## `assets/`

User-facing dashboard visual assets.

Approved structure:

```text
assets/dashboard/<dashboard-name>/
```

Examples:

```text
assets/dashboard/bttf/background.png
assets/dashboard/bttf/overlays/rpm_box.png
assets/dashboard/bttf/overlays/rpm_box_glow.png
assets/dashboard/bttf/digits/yellow/0.png
assets/dashboard/bttf/throttle/frame_000.png
```

Rules:

- Dashboard assets live under `assets/dashboard/`.
- Each dashboard/theme gets its own folder.
- Do not place dashboard assets under `internal/`.
- Do not place dashboard assets beside Go source.
- Do not create random folders like `images/`, `pngs/`, `ui-assets/`, `resources/`, or `stuff/`.

Preferred asset grouping:

```text
backgrounds/
overlays/
digits/
frames/
icons/
panels/
```

Use whatever names fit, but keep them under the dashboard folder.

---

## `configs/`

Example and reusable config files.

Approved structure:

```text
configs/examples/
```

Examples:

```text
configs/examples/dashboard-minimal.yaml
configs/examples/dashboard-bttf-sprite.yaml
configs/examples/dashboard-decoder-demo.yaml
```

Rules:

- Example dashboard configs go under `configs/examples/`.
- Do not put examples in the repo root.
- Keep `config.example.yaml` at root only if it remains the primary quick-start example.
- If the root `config.example.yaml` becomes v2, it should reference or mirror one clean example.

Not allowed:

```text
examples/
sample-configs/
yaml/
dashboard-configs/
```

without explicit approval.

---

## `docs/`

Project documentation.

Dashboard v2 docs live here:

```text
docs/dashboard/v2/
```

Required planning docs:

```text
docs/dashboard/v2/README.md
docs/dashboard/v2/overview.md
docs/dashboard/v2/prompts.md
docs/dashboard/v2/reference.md
docs/dashboard/v2/repo-structure-guardrails.md
docs/dashboard/v2/current-status.md
docs/dashboard/v2/decisions.md
docs/dashboard/v2/human-process.md
docs/dashboard/v2/CHANGES.md
```

Rules:

- Dashboard v2 docs stay under `docs/dashboard/v2/`.
- Do not scatter planning docs in root.
- Do not add random `NOTES.md`, `PLAN.md`, or `TODO.md` files at root.
- If a new planning doc is needed, put it under `docs/dashboard/v2/`.

The root README should summarise and link. It should not become the municipal archive.

---

## `testdata/`

Test fixtures only.

Approved structure:

```text
testdata/dashboard/
  assets/
  configs/
```

Examples:

```text
testdata/dashboard/assets/tiny.png
testdata/dashboard/assets/digits/0.png
testdata/dashboard/assets/frames/frame_000.png
testdata/dashboard/configs/valid-minimal.yaml
testdata/dashboard/configs/invalid-missing-asset.yaml
testdata/dashboard/configs/invalid-decoder-ref.yaml
```

Rules:

- Test-only assets go under `testdata/`.
- Do not use production dashboard assets for unit tests unless intentionally testing those assets.
- Do not put test fixtures in `assets/`.
- Do not put test fixtures beside Go files unless there is a strong Go-package reason.

---

## What not to create

Do not create these without explicit approval:

```text
dashboard/
dashboard2/
new-dashboard/
newui/
ui2/
display/
display2/
visuals/
engine/
engine2/
renderer/
resources/
images/
pngs/
examples/
scratch/
tmp/
experimental/
prototype/
final/
final-final/
final_final_really_final/
```

A future chat may think these are helpful. It is wrong. It is trying its best, but so is mildew.

---

## Package naming rules

Prefer boring package names:

```text
assets
decoders
scene
renderer
config
sensors
logger
```

Avoid vague package names:

```text
common
utils
helpers
misc
shared
stuff
engine
core
```

If a `utils` package appears, something has probably lost its name tag.

---

## File naming rules

Use short lowercase filenames:

```text
state_store.go
frame_index.go
sprite_text.go
validation.go
```

Avoid:

```text
DashboardSceneRendererFinal.go
new_stuff.go
helpers2.go
misc.go
everything.go
```

One file should have one obvious responsibility.

---

## Stage-specific placement rules

## v2.0.x - New config schema

Expected areas touched:

```text
internal/config/
config.example.yaml
configs/examples/
docs/dashboard/v2/
```

Avoid touching:

```text
internal/dashboard/
assets/
internal/sensors/elmobd_reader.go
internal/logger/
```

---

## v2.1.x - Config validation only

Expected areas touched:

```text
internal/config/
configs/examples/
testdata/dashboard/configs/
docs/dashboard/v2/
```

Avoid touching:

```text
internal/dashboard/renderer/
assets/
internal/sensors/
internal/logger/
```

---

## v2.2.x - Sensor state boundary

Expected areas touched:

```text
internal/sensors/state.go
internal/sensors/state_store.go
cmd/GoDriveLog/main.go
internal/config/
```

Avoid touching:

```text
internal/dashboard/assets/
internal/dashboard/renderer/
assets/
configs/examples/ unless needed
```

---

## v2.3.x - Decoder engine

Expected areas touched:

```text
internal/dashboard/decoders/
testdata/dashboard/configs/
docs/dashboard/v2/
```

Avoid touching:

```text
internal/dashboard/renderer/
assets/
internal/sensors/elmobd_reader.go
internal/logger/
```

---

## v2.4.x - Asset registry

Expected areas touched:

```text
internal/dashboard/assets/
testdata/dashboard/assets/
testdata/dashboard/configs/
docs/dashboard/v2/
```

Avoid touching:

```text
internal/dashboard/renderer/
internal/sensors/
internal/logger/
```

---

## v2.5.x - Scene primitives

Expected areas touched:

```text
internal/dashboard/scene/
internal/dashboard/decoders/ only if integration is required
internal/dashboard/assets/ only if integration is required
testdata/dashboard/
docs/dashboard/v2/
```

Avoid touching:

```text
internal/dashboard/renderer/fyne/
internal/sensors/elmobd_reader.go
internal/logger/
```

---

## v2.6.x - Fyne scene renderer

Expected areas touched:

```text
internal/dashboard/renderer/fyne/
internal/dashboard/scene/
cmd/GoDriveLog/main.go
```

Avoid touching:

```text
internal/sensors/elmobd_reader.go
internal/logger/
configs/examples/ unless required for manual run
```

---

## v2.7.x - First real dashboard

Expected areas touched:

```text
assets/dashboard/
configs/examples/
docs/dashboard/v2/
README.md
```

Possibly touched:

```text
cmd/GoDriveLog/main.go
internal/dashboard/
```

Avoid touching:

```text
internal/sensors/elmobd_reader.go
internal/logger/
```

---

## v2.8.x - Remove old widgets

Expected areas touched:

```text
widgets/
internal/ui/
internal/config/
cmd/GoDriveLog/main.go
README.md
configs/examples/
docs/dashboard/v2/
```

Rules:

- Remove old widget references only after new dashboard path works.
- Do not keep compatibility shims.
- Do not leave old widget examples in docs.

---

## v2.9.x - Reusable block library

Expected areas touched:

```text
internal/dashboard/scene/
internal/dashboard/decoders/
configs/examples/
docs/dashboard/v2/
testdata/dashboard/
```

Avoid touching:

```text
internal/sensors/elmobd_reader.go
internal/logger/
```

---

## When to stop and ask

Stop and ask before doing any of these:

- creating a new top-level directory
- moving existing top-level directories
- renaming major packages
- introducing a new renderer other than Fyne
- adding a plugin system
- adding a scripting/expression language
- adding remote dashboard sync
- changing logger behaviour
- changing real OBD reader behaviour
- preserving legacy `display.widget` compatibility
- deleting large existing directories before the planned v2.8.x stage

Stopping is not failure. It is putting the circular saw down before adjusting the fence.

---

## Branch prompt guardrail snippet

Add this to every implementation prompt:

```text
Repo structure guardrail:
Read and obey docs/dashboard/v2/repo-structure-guardrails.md.

Do not create new top-level directories.
Do not create new dashboard-related packages outside internal/dashboard/.
Do not place dashboard assets outside assets/dashboard/.
Do not place dashboard examples outside configs/examples/.
Do not place dashboard v2 docs outside docs/dashboard/v2/.
Do not place test fixtures outside testdata/dashboard/.
If the approved structure does not fit the task, stop and ask before changing it.
```

---

## Review checklist

Before accepting a branch, check:

```text
[ ] No unapproved top-level directories added.
[ ] Dashboard code is under internal/dashboard/.
[ ] Sensor state code is under internal/sensors/.
[ ] Config code is under internal/config/.
[ ] Assets are under assets/dashboard/.
[ ] Example configs are under configs/examples/.
[ ] Test fixtures are under testdata/dashboard/.
[ ] Docs are under docs/dashboard/v2/.
[ ] No random utils/helpers/misc package was added.
[ ] No old display.widget compatibility was introduced.
[ ] No unrelated OBD/logger changes were made.
```

---

## Final rule

```text
If a future chat cannot explain why a file belongs where it placed it, the file is probably in the wrong place.
```

A tidy repo is not bureaucracy. It is sweeping the workshop so you can find the chisel before stepping on it.
