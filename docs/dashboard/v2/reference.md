# Dashboard Implementation Reference - v2.x.x series

## Core design rule

```text
Old: PID owns display widget
New: Dashboard scene binds to sensor state
```

This is the main cut. Protect it.

If a change makes a sensor know how it should look, it is probably dragging the project backwards by the ankle.

---

## Proposed top-level config shape

```yaml
mock_mode: true
obd_address: serial:///dev/ttyUSB0
obd_debug: false

log:
  rotate: daily
  directory: ./log

vehicle:
  name: "VW Caddy"

sensors:
  rpm:
    type: obd
    pid: "010C"
    unit: rpm
    refresh: 250
    min: 0
    max: 7000
    log: true

  throttle_position:
    type: obd
    pid: "0111"
    unit: "%"
    refresh: 500
    min: 0
    max: 100
    log: true

dashboard:
  canvas:
    width: 800
    height: 480

assets: {}

decoders: {}

blocks: {}

layers: []
```

---

## Vocabulary

| Term | Meaning |
|---|---|
| Sensor | A configured data source such as RPM, speed, throttle, voltage |
| State | Latest runtime value/status for each sensor |
| Decoder | Reusable logic block converting state into percent, frame index, text, or state |
| Asset | Image, frame set, charset, or later font/nine-slice |
| Scene element | A drawable configured item |
| Layer | Scene element with z-order |
| Block | Reusable visual composition with named inputs |
| Condition | Rule that controls visibility or state |
| Renderer | Fyne-specific drawing layer |

---

## Minimal decoder set

### normalize

```yaml
rpm_percent:
  type: normalize
  input: sensor.rpm.value
  min: 0
  max: 6000
  clamp: true
```

Output:

```text
0.0 to 1.0
```

### threshold

```yaml
rpm_redline_state:
  type: threshold
  input: sensor.rpm.value
  states:
    normal: { lt: 5200 }
    redline: { gte: 5200 }
```

Output:

```text
normal | redline
```

### frame_index

```yaml
throttle_frame:
  type: frame_index
  input: sensor.throttle_position.value
  min: 0
  max: 100
  frames: 101
```

Output:

```text
0..100
```

### format_number

```yaml
rpm_text:
  type: format_number
  input: sensor.rpm.value
  format: "0000"
```

Output:

```text
"0000".."7000"
```

### digits

```yaml
rpm_digits:
  type: digits
  input: decoder.rpm_text
```

Output:

```text
["0", "3", "5", "0"]
```

### boolean

```yaml
is_moving:
  type: boolean
  input: sensor.speed.value
  gte: 1
```

Output:

```text
true | false
```

---

## Minimal asset set

### image

```yaml
assets:
  bttf_background:
    type: image
    path: assets/dashboards/bttf/background.png
```

### frame_set

```yaml
assets:
  throttle_squiggle:
    type: frame_set
    pattern: assets/throttle/frame_{index:03}.png
    frames: 101
```

### charset

```yaml
assets:
  yellow_7seg:
    type: charset
    path: assets/digits/yellow
    chars:
      "0": "0.png"
      "1": "1.png"
      "2": "2.png"
      "3": "3.png"
      "4": "4.png"
      "5": "5.png"
      "6": "6.png"
      "7": "7.png"
      "8": "8.png"
      "9": "9.png"
      "-": "dash.png"
      ".": "dot.png"
```

---

## Minimal scene primitives

### image

```yaml
- id: background
  type: image
  asset: bttf_background
  x: 0
  y: 0
  width: 800
  height: 480
  z: 0
```

### sprite_frame

```yaml
- id: throttle_bar
  type: sprite_frame
  asset: throttle_squiggle
  frame: decoder.throttle_frame
  x: 120
  y: 350
  z: 20
```

### sprite_text

```yaml
- id: rpm_digits
  type: sprite_text
  asset: yellow_7seg
  value: sensor.rpm.value
  format: "0000"
  x: 80
  y: 110
  z: 30
```

### conditional overlay

```yaml
- id: rpm_redline_glow
  type: image
  asset: rpm_glow
  x: 40
  y: 80
  z: 25
  visible_when:
    decoder: rpm_redline_state
    equals: redline
```

---

## First vertical slice config sketch

```yaml
assets:
  background:
    type: image
    path: assets/dashboard/background.png

  rpm_box:
    type: image
    path: assets/dashboard/rpm_box.png

  rpm_box_glow:
    type: image
    path: assets/dashboard/rpm_box_glow.png

  yellow_7seg:
    type: charset
    path: assets/digits/yellow
    chars:
      "0": "0.png"
      "1": "1.png"
      "2": "2.png"
      "3": "3.png"
      "4": "4.png"
      "5": "5.png"
      "6": "6.png"
      "7": "7.png"
      "8": "8.png"
      "9": "9.png"

  throttle_frames:
    type: frame_set
    pattern: assets/throttle/frame_{index:03}.png
    frames: 101

decoders:
  rpm_redline_state:
    type: threshold
    input: sensor.rpm.value
    states:
      normal: { lt: 5200 }
      redline: { gte: 5200 }

  throttle_frame:
    type: frame_index
    input: sensor.throttle_position.value
    min: 0
    max: 100
    frames: 101

dashboard:
  canvas:
    width: 800
    height: 480

  layers:
    - id: bg
      type: image
      asset: background
      x: 0
      y: 0
      width: 800
      height: 480
      z: 0

    - id: rpm_box
      type: image
      asset: rpm_box
      x: 40
      y: 80
      z: 10

    - id: rpm_glow
      type: image
      asset: rpm_box_glow
      x: 40
      y: 80
      z: 11
      visible_when:
        decoder: rpm_redline_state
        equals: redline

    - id: rpm_digits
      type: sprite_text
      asset: yellow_7seg
      value: sensor.rpm.value
      format: "0000"
      x: 80
      y: 115
      z: 20

    - id: throttle_bar
      type: sprite_frame
      asset: throttle_frames
      frame: decoder.throttle_frame
      x: 80
      y: 350
      z: 20
```

---

## Validation checklist

### Config

- Top-level IDs are unique.
- Sensor IDs are valid.
- Dashboard canvas width/height are positive.
- Asset references exist.
- Decoder references exist.
- Scene element IDs are unique within their scope.
- Layer geometry is valid.
- z-order is deterministic.
- Conditions reference real values.

### Sensors/state

- Sensor value updates latest state.
- Sensor error updates state.
- Stale state can be represented.
- Dashboard can read state without PID config.
- Logging remains independent.

### Decoders

- Unknown input fails.
- Invalid ranges fail.
- Frame index clamps correctly.
- Threshold order is deterministic.
- Number formatting is predictable.
- Decoder output type is known.

### Assets

- Missing files fail.
- Frame sets verify count.
- Charsets verify required chars.
- Asset paths resolve relative to config.
- Images are cached.

### Rendering

- Background draws first.
- z-order is respected.
- Hidden elements do not draw.
- sprite_text updates with state.
- sprite_frame updates with decoder output.
- Error/stale display is visible.

---

## Things not to do yet

```text
legacy support
remote dashboard sync
visual editor
SVG import
full expression language
plugin system
advanced animation system
theme marketplace
any sentence containing "while we are here" unless it is followed by "we stop"
```
