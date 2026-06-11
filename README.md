# GoDriveLog

GoDriveLog is a deliberately small Go/Fyne in-vehicle telemetry dashboard for Raspberry Pi-style installs.

The project is being reshaped toward a cleaner v3 runtime:

```text
vehicle endpoint
-> sensor polling runtime
-> sensor events
-> logs and dashboards as subscribers
```

The goal is simple: read vehicle telemetry, keep the runtime boring, log useful data, and render a dashboard that can look convincingly like real retro hardware instead of a web page wearing driving gloves.

## Current status

GoDriveLog currently contains working dashboard/runtime pieces from earlier versions, including Fyne rendering, OBD reader plumbing, JSONL logging, and dashboard asset experiments.

The v3 direction is documented under `docs/v3/` and is the target for cleanup and future implementation work. Some older config/runtime documents may still describe legacy v2 concepts while the repo is being migrated.

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
- If a config item exists, it is active.
- No `default_vehicle`.
- No `active_displays`.
- No separate top-level `displays` section.
- No `mock` / `real` source switch in GoDriveLog config.
- Bench testing should use an OBD-like endpoint, for example `tcp://127.0.0.1:35000`.

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
docs/v3/
docs/v3/examples/
docs/archive/
```

The examples under `docs/v3/examples/` are design examples, not final implementation contracts. They are there to argue with productively before code calcifies around the wrong idea like a fossilised ferret.

## Build

From the repository root:

```bash
go mod tidy
go build ./cmd/GoDriveLog
```

The binary will be written to the current directory as `GoDriveLog` unless you pass `-o`.

## Raspberry Pi notes

Fyne uses Go modules. The usual setup is:

```bash
go get fyne.io/fyne/v2@latest
go install fyne.io/tools/cmd/fyne@latest
```

On Raspberry Pi OS you may also need desktop/OpenGL build dependencies such as gcc, pkg-config, GL/X11 headers, xcursor, xrandr, xinerama, xi, and xxf86vm development packages.

## OBD transport

The intended v3 model is that GoDriveLog connects to an OBD-like endpoint:

```yaml
vehicles:
  vw_caddy:
    name: "VW Caddy"
    obd:
      address: "serial:///dev/ttyUSB0"
      timeout: 1000
```

For bench/simulator work, use the same connection path with a TCP endpoint:

```yaml
vehicles:
  bench:
    name: "Bench Simulator"
    obd:
      address: "tcp://127.0.0.1:35000"
      timeout: 1000
```

GoDriveLog should not need to know whether the endpoint is real hardware or a simulator. That is the point. The less the runtime knows, the less it can invent.

## Logging

The intended logging model is subscriber-based:

```yaml
logs:
  jsonl:
    path: "logs/godrivelog.jsonl"
    sensors:
      - speed
      - rpm
      - coolant_temperature
```

Current agreed/simple behaviour:

```text
first reading logs
value changes log
status changes log
unchanged duplicate readings do not spam logs
```

No log refresh. No per-sensor log refresh. No override knobs unless reality turns up with a receipt.

## Non-goals

Avoid these until proven necessary:

- plugin systems
- source orchestration
- live config reload
- dashboard scripting
- config inheritance
- generic event buses
- enable/disable flags everywhere
- mock/real branching leaking through the core runtime

The preferred architecture is still:

```text
clean boundaries, boring implementation
```

Boring is not an insult. Boring code is code that lets you sleep.
