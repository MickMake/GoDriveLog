# GoDriveLog v3 docs

This directory describes the intended v3 direction for GoDriveLog.

The codebase may still contain current-runtime or earlier-dashboard pieces while migration is underway. These docs define where the project is heading and how to get there without dragging old assumptions into the new shape.

## Read these first

1. `config.example.yaml`  
   Small starter config using the v3 shape.

2. `config.full.yaml`  
   Heavily-commented full config reference with simple/BTTF, Nissan 300ZX Z31-inspired, and Honda S2000-inspired dashboard examples.

3. `GoStructsConfig.md`  
   Intended Go struct shape for the v3 config model.

4. `ImplementationGuardrails.md`  
   Rules for writing v3 code against the documented model.

5. `MigrationGuardrails.md`  
   How to move from current code to the v3 target safely.

6. `PerformanceGuardrails.md`  
   How to fix current display speed and design v3 renderer performance without warping the schema.

7. `DirectoryStructure.md`  
   Intended repo/package layout for v3.

8. `examples/`  
   Standalone dashboard examples that should validate against the same v3 schema rules as `config.example.yaml` and `config.full.yaml`.

## Target runtime model

```text
vehicle endpoint
-> sensor polling runtime
-> sensor events
-> logs and dashboards as subscribers
```

## Target config shape

```yaml
vehicles:
sensors:
assets:
logs:
dashboards:
```

Treat those documented root sections as an allow-list.

Strict v3 config loading should reject unknown fields at every documented level, not only at the root.

## Example stance

Every file under `docs/v3/examples/` should validate against the same v3 schema rules as `config.example.yaml` and `config.full.yaml`.

Active v3 examples use repository-root-relative asset paths, for example:

```text
assets/dashboard/simple/panel/background.png
```

Do not teach multiple active asset path dialects.

## Migration stance

Current code may be useful. Current code is not automatically the target.

Use boundary adapters while migrating, but do not let migration behaviour leak into the v3 core model.

Useful rule:

```text
Adapters at boundaries are allowed.
Old assumptions inside the v3 core are not.
```

## Performance stance

The current display path may need tactical speed work before full v3 migration is complete.

That is allowed, but it should remain local to the current renderer/runtime path unless the lesson cleanly applies to the v3 widget renderer.

Useful rule:

```text
Optimise rendering locally.
Keep the v3 model clean.
```

Performance work must not hide stale, error, or recovery transitions.

## Archive docs

Older/current-dashboard documents belong under `docs/archive/`.

Archive docs may describe old behaviour. Active v3 docs should describe the target shape or the migration path toward it.
