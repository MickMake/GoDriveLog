# GoDriveLog v3 frozen schema target

Status: v3.0.1 implementation slice output  
Target version: `v3.0.1`  
Branch: `v3.0.1-freeze-v3-docs-schema`

## 1. Purpose

This document freezes the active v3 docs as the implementation target for the next slice.

The next implementation slice, `v3.0.2`, should be able to build strict config loading and validation from these docs without inventing missing schema rules or preserving current config shapes.

## 2. Source documents reviewed

The frozen target is defined by these active v3 documents:

- `docs/v3/README.md`
- `docs/v3/GoStructsConfig.md`
- `docs/v3/ImplementationGuardrails.md`
- `docs/v3/MigrationGuardrails.md`
- `docs/v3/PerformanceGuardrails.md`
- `docs/v3/config.example.yaml`
- `docs/v3/config.full.yaml`
- `docs/v3/examples/README.md`
- `docs/v3/examples/simple_speed_warning.yaml`
- `docs/v3/examples/nissan_300zx_z31_inspired.yaml`
- `docs/v3/examples/honda_s2000_inspired.yaml`
- `docs/v3/WorkingCodeInventory.md`

Archive docs under `docs/archive/` are not active v3 schema sources.

## 3. Frozen root schema

The documented v3 root schema is exactly:

```yaml
vehicles:
sensors:
assets:
logs:
dashboards:
```

Rules:

- Treat these five root sections as an allow-list.
- Strict v3 loading must reject unknown root fields.
- Strict v3 loading must reject unknown nested fields at every documented level.
- Do not add compatibility aliases for old/current config fields.
- Do not auto-convert current config files into v3 config.

The old/current root fields are not v3 schema:

- `obd`
- `log`
- `vehicle`
- singular `dashboard`

## 4. Ownership model

The frozen ownership model is:

```text
selected vehicle
-> OBD endpoint
-> sensor polling runtime
-> sensor events
-> selected logs and dashboards as subscribers
```

Rules:

- Vehicles are runtime profiles.
- Vehicles choose one OBD-like endpoint.
- Vehicles select global log definitions by ID.
- Vehicles select global dashboard definitions by ID.
- Sensors are global catalogue entries.
- Assets are global catalogue entries.
- Logs reference global sensor IDs.
- Dashboard widgets reference global sensor IDs and global asset IDs.
- Vehicles do not directly list sensors.
- Vehicles do not directly list assets.
- Logs do not poll.
- Dashboards do not poll.
- Only sensors own polling cadence via `sensors.<id>.poll`.

## 5. Runtime vehicle selection

The config schema does not contain an active/current vehicle field.

Runtime vehicle selection is outside the config document:

- If exactly one vehicle exists, runtime may select it automatically.
- If multiple vehicles exist, runtime must require explicit selection, for example `--vehicle <id>`.

This is intentional. Do not add a top-level `vehicle`, `selected_vehicle`, or `default_vehicle` field for v3.0.2.

## 6. Asset path stance

Active v3 configs and examples use repository-root-relative asset paths.

Example:

```text
assets/dashboard/simple/panel/background.png
```

Rules:

- Do not add `asset_root` to v3 config.
- Do not teach config-file-relative paths in active v3 examples.
- Do not allow remote asset paths in the first v3 loader.
- Asset IDs only need to be unique within their asset family.
- Widget `type` determines the asset family used to resolve `asset`.

## 7. Active examples stance

These are active full-schema example config files and should validate under the same v3 rules as `config.example.yaml` and `config.full.yaml`:

- `docs/v3/examples/simple_speed_warning.yaml`
- `docs/v3/examples/nissan_300zx_z31_inspired.yaml`
- `docs/v3/examples/honda_s2000_inspired.yaml`

`docs/v3/examples/README.md` is explanatory prose, not a config file.

All active examples use the same root section model:

```text
vehicles, sensors, assets, logs, dashboards
```

## 8. Widget and asset family mapping

The frozen widget-to-asset-family mapping is:

| Widget type | Required asset family |
|---|---|
| `image` | `assets.image_sets` |
| `digit_display` | `assets.digit_sets` |
| `bar_display` | `assets.bar_sets` |
| `frame_gauge` | `assets.frame_sets` |
| `indicator` | `assets.indicator_sets` |

Rules:

- `image` widgets may omit `sensor`.
- Non-image widgets must reference an existing global sensor.
- Widgets must reference an existing asset in the correct asset family for the widget type.
- Widget IDs must be unique within a dashboard.
- Widget IDs do not need to be globally unique.
- Widget positions use `position`, not `rect`.

## 9. Validation target for v3.0.2

The next slice should implement strict v3 config loading and validation against this target.

Minimum v3.0.2 validation should include:

- root allow-list enforcement
- nested unknown-field rejection
- ID pattern validation: `^[a-z][a-z0-9_]*$`
- at least one vehicle
- explicit runtime vehicle selection when multiple vehicles exist
- vehicle OBD address validation for `serial://` and `tcp://`
- vehicle OBD timeout validation
- vehicle log/dashboard reference validation
- selected dashboard display-collision validation
- sensor `poll > 0`
- sensor `min < max` when both are present
- log sensor reference validation
- dashboard size validation
- widget ID/type/asset/position validation
- non-image widget sensor reference validation
- widget asset-family reference validation
- widget `min < max` when both are present
- digit display character/decimal-point validation where practical
- indicator set `off`, `on`, and `unknown` validation
- bar set `off` validation
- bar widget `cells > 0`
- bar zone ordering and cell-reference validation
- frame range `first <= last`

## 10. Implementation blockers found

No schema blocker was found that requires adding new v3 schema features before v3.0.2.

Important implementation notes for v3.0.2:

- The loader should live beside current config code until the v3 path replaces the old runtime.
- Do not make the v3 loader accept current config shapes.
- Do not add compatibility aliases.
- Do not let existing runtime wiring decide v3 schema shape.
- Keep validation errors clear enough to fix config mistakes without deep debugging.

## 11. Explicit non-goals

This slice does not:

- implement Go config structs
- implement YAML loading
- implement runtime vehicle selection
- implement endpoint adapters
- implement sensor polling
- implement logging subscribers
- implement dashboard rendering
- add config aliases
- change the v3 schema for old-code convenience

## 12. Summary decision

The active v3 schema target is frozen enough for `v3.0.2` strict config loading and validation.

Future implementation may tighten validation details, but schema expansion should happen only when a real blocker is found and documented before code is changed.
