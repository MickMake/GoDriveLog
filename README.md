# GoDriveLog

GoDriveLog is a deliberately small Go/Ebiten in-vehicle telemetry dashboard for Raspberry Pi-style installs.

The project is being reshaped around the active v3 runtime:

```text
vehicle endpoint
-> sensor polling runtime
-> sensor events
-> logs and dashboards as subscribers
-> Ebiten dashboard renderer
```

The goal is simple: read vehicle telemetry, keep the runtime boring, log useful data, and render a dashboard that can look convincingly like real retro hardware instead of a web page wearing driving gloves.

## Current status

GoDriveLog now uses Ebiten for the active v3 dashboard command path. Earlier Fyne dashboard code lives in the v3.2.x line only; v3.3.x and later are Ebiten-first.

The current v3.3 implementation state is documented under `docs/v3.3/`. Some older config/runtime documents may still describe legacy concepts while the repo is being migrated.

## v3 direction

The intended v3 config shape is:

```yaml
vehicles: {}
sensors: {}
assets: {}
logs: {}
dashboards: {}
```

The important design rules are:

- Sensors own polling cadence using `poll`.
- Logs and dashboards subscribe to sensor events.
- Logs do not poll sensors independently.
- Dashboards do not fetch OBD values directly.
- If a documented config item exists, it is active.
- Each dashboard owns its physical/logical display target.
- GoDriveLog connects to an OBD-like endpoint address.
- Bench testing should use an OBD-like endpoint, for example `tcp://127.0.0.1:35000`.
- Unknown config fields should fail validation during v3 implementation.

Boring boundaries are intentional. Cleverness is allowed only when it pays rent and does not bring a YAML demon as a lodger.

## Runtime model

The intended v3 runtime is:

```text
load config
resolve vehicle
connect to the vehicle OBD-like endpoint
start sensor runtime
poll sensors according to sensors.<id>.poll
emit sensor reading/status events
logs receive selected sensor events
dashboards receive sensor events implied by widgets
render dashboard updates from event state
```

Sensor readings should carry their original read timestamp. Log writers may add their own write timestamp, but the sensor timestamp is the source of truth.

Sensor status should distinguish real values from trouble states such as:

```text
ok
stale
error
missing/unsupported
```

Do not use `0` as an error value. Zero is a perfectly respectable number and should not be framed for crimes committed by the transport layer.

## Dashboard asset direction

The v3 dashboard direction is asset-driven and photoreal-friendly.

Common render pattern:

```text
asset background
+ value/state-driven dynamic layer
+ optional foreground/glass/bezel overlay
= rendered widget
```

The intended asset families are:

```yaml
assets:
  digit_sets: {}
  bar_sets: {}
  frame_sets: {}
  indicator_sets: {}
  image_sets: {}
```

Typical widget types are expected to be:

```yaml
- type: digit_display
- type: bar_display
- type: frame_gauge
- type: indicator
- type: image
```

For PNG digit displays, formatted output should resolve characters rather than only numeric digits. A digit asset set should be able to provide characters such as:

```text
0 1 2 3 4 5 6 7 8 9 -
```

A blank slot should normally mean: draw the digit background only. Decimal points are overlays.

## Documentation

Useful docs live under:

```text
docs/v3.3/
docs/v3.2/
docs/archive/
```

The v3.2 docs describe the final supported Fyne dashboard line. The active v3.3 docs describe the Ebiten migration and current renderer path.

## Build

From the repository root:

```bash
go mod tidy
go build ./cmd/GoDriveLog
```

The binary will be written to the current directory as `GoDriveLog` unless you pass `-o`.

## Raspberry Pi notes

The active v3.3 dashboard renderer is Ebiten. Raspberry Pi builds should focus on Go, graphics/display dependencies needed by Ebiten, and the selected kiosk/display setup.

## OBD transport

The intended v3 model is that GoDriveLog connects to an OBD-like endpoint:

```yaml
vehicles:
  vw_caddy:
    name: "VW Caddy"
```
