# GoDriveLog Dashboard Config Examples

These are brainstorming schemas, not implementation contracts. The goal is to test whether a dashboard can be built from a small engine plus large config and image assets.

Assumed renderer primitives:

- `image`: draw a static image.
- `sprite_text`: render formatted text using character sprites.
- `sprite_frame`: map a value/index to one image from a sequence.
- `threshold_image`: choose an image from named/ranged states.
- `group`: reusable block with named inputs.
- `visible_when`: conditionally draw an element.
- `opacity_when`: conditionally fade/dim an element.
- `z`: draw order.

Assumed data model:

- The display reads latest sensor state by sensor id, such as `rpm`, `speed`, `coolant_temp`, `battery_voltage`, `throttle_position`, and `fuel_level`.
- Display config can override visual min/max without changing canonical OBD sensor meaning.
- Missing/stale sensors are display states, not fatal errors. A typo should create a sad grey widget, not a dashboard-shaped crater.

---

## Example 1 — BTTF-style sprite dashboard

This example uses background art, seven-segment digit sprites, a squiggly throttle bar, and a redline glow overlay around RPM.

```yaml
schema_version: 1
kind: godrivelog_dashboard
id: bttf_sprite_main
name: BTTF Sprite Main

canvas:
  width: 800
  height: 480
  scale_mode: fit
  background_colour: "#000000"

assets:
  images:
    bg_main: dashboards/bttf/backgrounds/main_800x480.png
    glass_overlay: dashboards/bttf/overlays/glass_reflection.png
    scanlines: dashboards/bttf/overlays/scanlines_20pct.png
    rpm_box_normal: dashboards/bttf/panels/rpm_box_normal.png
    rpm_box_glow: dashboards/bttf/panels/rpm_box_red_glow.png
    speed_box_normal: dashboards/bttf/panels/speed_box_normal.png
    stale_mask: dashboards/common/overlays/stale_grey_mask.png
    mqtt_disconnected_banner: dashboards/common/banners/mqtt_disconnected.png

  charsets:
    yellow_7seg:
      base_path: dashboards/bttf/sprites/7seg/yellow
      width: 48
      height: 84
      spacing: 4
      map:
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
        " ": "blank.png"

    green_7seg:
      base_path: dashboards/bttf/sprites/7seg/green
      width: 48
      height: 84
      spacing: 4
      map:
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
        " ": "blank.png"

  frame_sets:
    throttle_squiggle:
      base_path: dashboards/bttf/sprites/throttle_squiggle
      pattern: frame_{index:03}.png
      frames: 101
      index_min: 0
      index_max: 100

    coolant_tube:
      base_path: dashboards/bttf/sprites/coolant_tube
      pattern: frame_{index:03}.png
      frames: 101
      index_min: 0
      index_max: 100

signals:
  rpm:
    sensor: rpm
    visual_min: 0
    visual_max: 6500
  speed:
    sensor: speed
    visual_min: 0
    visual_max: 180
  throttle:
    sensor: throttle_position
    visual_min: 0
    visual_max: 100
  coolant:
    sensor: coolant_temp
    visual_min: 40
    visual_max: 115

transforms:
  rpm_percent:
    type: normalise
    input: rpm.value
    min: rpm.visual_min
    max: rpm.visual_max
    clamp: true
    output_min: 0
    output_max: 100

  rpm_redline_state:
    type: threshold
    input: rpm.value
    states:
      normal: { lt: 5000 }
      warning: { gte: 5000, lt: 5700 }
      redline: { gte: 5700 }

  throttle_frame:
    type: frame_index
    input: throttle.value
    min: throttle.visual_min
    max: throttle.visual_max
    frames: 101
    clamp: true

  coolant_frame:
    type: frame_index
    input: coolant.value
    min: coolant.visual_min
    max: coolant.visual_max
    frames: 101
    clamp: true

blocks:
  sprite_number_panel:
    inputs:
      value: number
      format: string
      charset: charset
      panel_image: image
      stale_sensor: sensor_ref
    elements:
      - type: image
        id: panel
        image: panel_image
        x: 0
        y: 0
        z: 0
      - type: sprite_text
        id: digits
        value: value
        format: format
        charset: charset
        x: 24
        y: 18
        z: 10
      - type: image
        id: stale_dim
        image: stale_mask
        x: 0
        y: 0
        z: 90
        visible_when:
          sensor_state: stale_sensor
          in: [unseen, stale, unsupported, disconnected]

  rpm_redline_box:
    inputs:
      state: state
    elements:
      - type: image
        id: rpm_box_normal
        image: rpm_box_normal
        x: 0
        y: 0
        z: 0
      - type: image
        id: rpm_box_glow
        image: rpm_box_glow
        x: -8
        y: -8
        z: 1
        visible_when:
          state: state
          equals: redline

layers:
  - type: image
    id: background
    image: bg_main
    x: 0
    y: 0
    z: 0

  - type: group
    id: rpm_panel
    use: rpm_redline_box
    state: rpm_redline_state
    x: 46
    y: 58
    z: 10

  - type: group
    id: rpm_digits
    use: sprite_number_panel
    value: rpm.value
    format: "0000"
    charset: yellow_7seg
    panel_image: rpm_box_normal
    stale_sensor: rpm
    x: 58
    y: 70
    z: 20

  - type: group
    id: speed_digits
    use: sprite_number_panel
    value: speed.value
    format: "000"
    charset: green_7seg
    panel_image: speed_box_normal
    stale_sensor: speed
    x: 460
    y: 70
    z: 20

  - type: sprite_frame
    id: throttle_squiggle
    frame_set: throttle_squiggle
    index: throttle_frame
    x: 80
    y: 260
    z: 30

  - type: sprite_frame
    id: coolant_tube
    frame_set: coolant_tube
    index: coolant_frame
    x: 430
    y: 260
    z: 30

  - type: image
    id: glass
    image: glass_overlay
    x: 0
    y: 0
    z: 900

  - type: image
    id: scanlines
    image: scanlines
    x: 0
    y: 0
    z: 910

  - type: image
    id: mqtt_disconnected
    image: mqtt_disconnected_banner
    x: 180
    y: 210
    z: 1000
    visible_when:
      connection: mqtt
      equals: disconnected
```

Why this one matters: almost everything visual is art-driven. The renderer only knows how to draw images, render sprite text, choose frames, and apply conditions.

---

## Example 2 — Minimal reusable block library

This one is less flashy. It tests whether complex blocks can be reusable and sensor-agnostic.

```yaml
schema_version: 1
kind: godrivelog_dashboard
id: reusable_block_test
name: Reusable Block Test

canvas:
  width: 800
  height: 480
  scale_mode: fit
  background_colour: "#101010"

assets:
  images:
    bg_plain: dashboards/common/backgrounds/dark_grid.png
    rounded_box: dashboards/common/panels/rounded_box.png
    rounded_box_glow_red: dashboards/common/panels/rounded_box_glow_red.png
    rounded_box_glow_amber: dashboards/common/panels/rounded_box_glow_amber.png
    unavailable_overlay: dashboards/common/overlays/unavailable.png
    stale_overlay: dashboards/common/overlays/stale.png

  charsets:
    white_lcd:
      base_path: dashboards/common/sprites/lcd/white
      width: 36
      height: 60
      spacing: 2
      map:
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
        ".": "dot.png"
        "-": "dash.png"
        " ": "blank.png"
        "V": "V.png"
        "C": "C.png"

transforms:
  battery_state:
    type: threshold
    input: battery_voltage.value
    states:
      low: { lt: 12.0 }
      normal: { gte: 12.0, lt: 14.8 }
      high: { gte: 14.8 }

  coolant_state:
    type: threshold
    input: coolant_temp.value
    states:
      cold: { lt: 70 }
      normal: { gte: 70, lt: 100 }
      hot: { gte: 100 }

blocks:
  glowing_value_box:
    inputs:
      value: number
      label: string
      format: string
      charset: charset
      state: state
      warn_state: string
      danger_state: string
      sensor: sensor_ref
    elements:
      - type: image
        id: base_box
        image: rounded_box
        x: 0
        y: 0
        z: 0
      - type: image
        id: warn_glow
        image: rounded_box_glow_amber
        x: -6
        y: -6
        z: 1
        visible_when:
          state: state
          equals: warn_state
      - type: image
        id: danger_glow
        image: rounded_box_glow_red
        x: -8
        y: -8
        z: 2
        visible_when:
          state: state
          equals: danger_state
      - type: text
        id: label
        value: label
        x: 16
        y: 12
        font: default_bold
        size: 18
        colour: "#bbbbbb"
        z: 10
      - type: sprite_text
        id: value
        value: value
        format: format
        charset: charset
        x: 16
        y: 44
        z: 20
      - type: image
        id: stale
        image: stale_overlay
        x: 0
        y: 0
        z: 90
        visible_when:
          sensor_state: sensor
          in: [unseen, stale, disconnected]
      - type: image
        id: unavailable
        image: unavailable_overlay
        x: 0
        y: 0
        z: 91
        visible_when:
          sensor_state: sensor
          equals: unsupported

layers:
  - type: image
    id: background
    image: bg_plain
    x: 0
    y: 0
    z: 0

  - type: group
    id: coolant_box
    use: glowing_value_box
    value: coolant_temp.value
    label: COOLANT
    format: "000C"
    charset: white_lcd
    state: coolant_state
    warn_state: cold
    danger_state: hot
    sensor: coolant_temp
    x: 40
    y: 40
    z: 10

  - type: group
    id: battery_box
    use: glowing_value_box
    value: battery_voltage.value
    label: BATTERY
    format: "00.0V"
    charset: white_lcd
    state: battery_state
    warn_state: high
    danger_state: low
    sensor: battery_voltage
    x: 420
    y: 40
    z: 10

  - type: group
    id: rpm_box
    use: glowing_value_box
    value: rpm.value
    label: RPM
    format: "0000"
    charset: white_lcd
    state: rpm_state
    warn_state: high
    danger_state: redline
    sensor: rpm
    x: 40
    y: 220
    z: 10

  - type: group
    id: speed_box
    use: glowing_value_box
    value: speed.value
    label: SPEED
    format: "000"
    charset: white_lcd
    state: speed_state
    warn_state: fast
    danger_state: very_fast
    sensor: speed
    x: 420
    y: 220
    z: 10
```

Weak spot exposed by this example: `rpm_state` and `speed_state` are referenced but not defined. That could be allowed if the renderer treats missing transforms as `unknown`, but I would not. Better: config validation should fail clearly. Otherwise the dashboard becomes a haunted mirror with a YAML parser.

---

## Example 3 — Sloped bar and animated-ish gauge using frame sets

This tests the idea that many weird visual designs are just `value -> frame`.

```yaml
schema_version: 1
kind: godrivelog_dashboard
id: sloped_bar_gauge
name: Sloped Bar Gauge

canvas:
  width: 1024
  height: 600
  scale_mode: fit

assets:
  images:
    bg_carbon: dashboards/sloped/backgrounds/carbon_1024x600.png
    foreground_bezel: dashboards/sloped/overlays/bezel.png
    red_flash: dashboards/sloped/overlays/red_flash.png

  frame_sets:
    rpm_sloped_bar:
      base_path: dashboards/sloped/frames/rpm_bar
      pattern: rpm_{index:03}.png
      frames: 121
      index_min: 0
      index_max: 120

    speed_arc:
      base_path: dashboards/sloped/frames/speed_arc
      pattern: speed_{index:03}.png
      frames: 181
      index_min: 0
      index_max: 180

    fuel_squiggle:
      base_path: dashboards/sloped/frames/fuel_squiggle
      pattern: fuel_{index:03}.png
      frames: 101
      index_min: 0
      index_max: 100

transforms:
  rpm_bar_index:
    type: frame_index
    input: rpm.value
    min: 0
    max: 6000
    frames: 121
    clamp: true

  speed_arc_index:
    type: frame_index
    input: speed.value
    min: 0
    max: 180
    frames: 181
    clamp: true

  fuel_index:
    type: frame_index
    input: fuel_level.value
    min: 0
    max: 100
    frames: 101
    clamp: true

  aggressive_driving:
    type: boolean
    expression:
      any:
        - { input: rpm.value, gte: 5200 }
        - { input: throttle_position.value, gte: 90 }

layers:
  - type: image
    id: background
    image: bg_carbon
    x: 0
    y: 0
    z: 0

  - type: sprite_frame
    id: rpm_sloped_bar
    frame_set: rpm_sloped_bar
    index: rpm_bar_index
    x: 80
    y: 90
    z: 20

  - type: sprite_frame
    id: speed_arc
    frame_set: speed_arc
    index: speed_arc_index
    x: 360
    y: 70
    z: 20

  - type: sprite_frame
    id: fuel_squiggle
    frame_set: fuel_squiggle
    index: fuel_index
    x: 96
    y: 430
    z: 20

  - type: image
    id: red_flash
    image: red_flash
    x: 0
    y: 0
    z: 800
    opacity: 0.45
    visible_when:
      signal: aggressive_driving
      equals: true

  - type: image
    id: foreground_bezel
    image: foreground_bezel
    x: 0
    y: 0
    z: 900
```

This is the purest version of the idea. A complex sloped or squiggly graph is not a special widget. It is one frame set and one index transform.

---

## Example 4 — Full scene with background, reusable widgets, overlays, and foreground glass

This is closer to a complete dashboard scene.

```yaml
schema_version: 1
kind: godrivelog_dashboard
id: complete_scene_test
name: Complete Scene Test

canvas:
  width: 1280
  height: 720
  scale_mode: fit_crop
  safe_area:
    x: 20
    y: 20
    width: 1240
    height: 680

assets:
  images:
    bg_night: dashboards/complete/backgrounds/night_dash.png
    bg_day: dashboards/complete/backgrounds/day_dash.png
    vignette: dashboards/complete/overlays/vignette.png
    glass: dashboards/complete/overlays/glass.png
    dust: dashboards/complete/overlays/dust.png
    warning_bar: dashboards/complete/banners/warning_bar.png
    disconnected_bar: dashboards/complete/banners/disconnected_bar.png
    panel_large: dashboards/complete/panels/panel_large.png
    panel_small: dashboards/complete/panels/panel_small.png
    panel_large_glow: dashboards/complete/panels/panel_large_glow.png

  charsets:
    orange_digits:
      base_path: dashboards/complete/sprites/digits/orange
      width: 64
      height: 110
      spacing: 5
      map:
        "0": 0.png
        "1": 1.png
        "2": 2.png
        "3": 3.png
        "4": 4.png
        "5": 5.png
        "6": 6.png
        "7": 7.png
        "8": 8.png
        "9": 9.png
        ".": dot.png
        "-": dash.png
        " ": blank.png

    blue_digits:
      base_path: dashboards/complete/sprites/digits/blue
      width: 42
      height: 74
      spacing: 3
      map:
        "0": 0.png
        "1": 1.png
        "2": 2.png
        "3": 3.png
        "4": 4.png
        "5": 5.png
        "6": 6.png
        "7": 7.png
        "8": 8.png
        "9": 9.png
        ".": dot.png
        "-": dash.png
        " ": blank.png
        "V": V.png
        "C": C.png

  frame_sets:
    rpm_arc:
      base_path: dashboards/complete/frames/rpm_arc
      pattern: frame_{index:03}.png
      frames: 131
    temp_column:
      base_path: dashboards/complete/frames/temp_column
      pattern: frame_{index:03}.png
      frames: 101

transforms:
  rpm_state:
    type: threshold
    input: rpm.value
    states:
      normal: { lt: 4800 }
      high: { gte: 4800, lt: 5600 }
      redline: { gte: 5600 }

  rpm_arc_index:
    type: frame_index
    input: rpm.value
    min: 0
    max: 6500
    frames: 131
    clamp: true

  temp_index:
    type: frame_index
    input: coolant_temp.value
    min: 50
    max: 115
    frames: 101
    clamp: true

  show_global_warning:
    type: boolean
    expression:
      any:
        - { state: rpm_state, equals: redline }
        - { input: coolant_temp.value, gte: 105 }
        - { input: battery_voltage.value, lt: 11.8 }

  display_is_unhealthy:
    type: boolean
    expression:
      any:
        - { connection: mqtt, equals: disconnected }
        - { sensor_state: rpm, equals: stale }
        - { sensor_state: speed, equals: stale }

blocks:
  large_number_with_arc:
    inputs:
      value: number
      format: string
      charset: charset
      frame_set: frame_set
      frame_index: number
      state: state
      sensor: sensor_ref
    elements:
      - type: image
        image: panel_large
        x: 0
        y: 0
        z: 0
      - type: image
        image: panel_large_glow
        x: -14
        y: -14
        z: 1
        visible_when:
          state: state
          in: [high, redline]
      - type: sprite_frame
        frame_set: frame_set
        index: frame_index
        x: 24
        y: 20
        z: 10
      - type: sprite_text
        value: value
        format: format
        charset: charset
        x: 90
        y: 118
        z: 20
      - type: opacity
        target: self
        opacity: 0.35
        when:
          sensor_state: sensor
          in: [unseen, stale, unsupported, disconnected]

  small_metric:
    inputs:
      label: string
      value: number
      format: string
      charset: charset
      sensor: sensor_ref
    elements:
      - type: image
        image: panel_small
        x: 0
        y: 0
        z: 0
      - type: text
        value: label
        x: 14
        y: 10
        font: small_bold
        size: 18
        colour: "#888888"
        z: 10
      - type: sprite_text
        value: value
        format: format
        charset: charset
        x: 14
        y: 44
        z: 20
      - type: opacity
        target: self
        opacity: 0.35
        when:
          sensor_state: sensor
          in: [unseen, stale, unsupported, disconnected]

layers:
  - type: image
    id: bg_day
    image: bg_day
    x: 0
    y: 0
    z: 0
    visible_when:
      input: ambient_light.value
      gte: 40

  - type: image
    id: bg_night
    image: bg_night
    x: 0
    y: 0
    z: 0
    visible_when:
      input: ambient_light.value
      lt: 40

  - type: group
    id: rpm_main
    use: large_number_with_arc
    value: rpm.value
    format: "0000"
    charset: orange_digits
    frame_set: rpm_arc
    frame_index: rpm_arc_index
    state: rpm_state
    sensor: rpm
    x: 72
    y: 72
    z: 20

  - type: group
    id: speed_main
    use: small_metric
    label: SPEED
    value: speed.value
    format: "000"
    charset: blue_digits
    sensor: speed
    x: 850
    y: 80
    z: 20

  - type: group
    id: coolant_metric
    use: small_metric
    label: COOLANT
    value: coolant_temp.value
    format: "000C"
    charset: blue_digits
    sensor: coolant_temp
    x: 850
    y: 235
    z: 20

  - type: group
    id: battery_metric
    use: small_metric
    label: BATTERY
    value: battery_voltage.value
    format: "00.0V"
    charset: blue_digits
    sensor: battery_voltage
    x: 850
    y: 390
    z: 20

  - type: sprite_frame
    id: temp_column
    frame_set: temp_column
    index: temp_index
    x: 690
    y: 150
    z: 30

  - type: image
    id: warning_bar
    image: warning_bar
    x: 220
    y: 610
    z: 700
    visible_when:
      signal: show_global_warning
      equals: true

  - type: image
    id: disconnected_bar
    image: disconnected_bar
    x: 220
    y: 610
    z: 710
    visible_when:
      signal: display_is_unhealthy
      equals: true

  - type: image
    id: vignette
    image: vignette
    x: 0
    y: 0
    z: 880

  - type: image
    id: dust
    image: dust
    x: 0
    y: 0
    z: 890

  - type: image
    id: glass
    image: glass
    x: 0
    y: 0
    z: 900
```

This example shows the general rule: a complete dashboard is just a layered scene. Backgrounds, metrics, warnings, reflections, stale overlays, and foreground scratches are all the same kind of thing: elements with draw order and optional conditions.

---

## Example 5 — Asset-heavy, code-light dashboard with named decoder blocks

This tests your idea that coded logic can exist, but be reused as named building blocks.

```yaml
schema_version: 1
kind: godrivelog_dashboard
id: decoder_block_library_demo
name: Decoder Block Library Demo

canvas:
  width: 800
  height: 480

assets:
  images:
    bg: dashboards/library_demo/background.png
    rpm_normal: dashboards/library_demo/rpm/rpm_normal.png
    rpm_high: dashboards/library_demo/rpm/rpm_high.png
    rpm_redline: dashboards/library_demo/rpm/rpm_redline.png
    box_glow: dashboards/library_demo/overlays/box_glow.png

  frame_sets:
    throttle_weird_bar:
      base_path: dashboards/library_demo/throttle_weird_bar
      pattern: state_{index:03}.png
      frames: 101

  charsets:
    main_digits:
      base_path: dashboards/library_demo/digits
      width: 50
      height: 90
      spacing: 4
      map:
        "0": 0.png
        "1": 1.png
        "2": 2.png
        "3": 3.png
        "4": 4.png
        "5": 5.png
        "6": 6.png
        "7": 7.png
        "8": 8.png
        "9": 9.png
        " ": blank.png

# These are reusable decoder building blocks.
# Some could be built-ins. Some could be custom registered Go decoders later.
decoder_blocks:
  rpm_drive_mood:
    type: threshold_state
    inputs:
      value: rpm.value
    states:
      idle: { lt: 1000 }
      cruise: { gte: 1000, lt: 3200 }
      push: { gte: 3200, lt: 5200 }
      redline: { gte: 5200 }

  value_to_percent:
    type: normalise
    inputs:
      value: required
      min: required
      max: required
    clamp: true
    output_min: 0
    output_max: 100

  percent_to_101_frame:
    type: frame_index
    inputs:
      percent: required
    min: 0
    max: 100
    frames: 101
    clamp: true

# These are actual transform instances made from decoder blocks.
transforms:
  rpm_mood:
    use: rpm_drive_mood
    value: rpm.value

  throttle_percent:
    use: value_to_percent
    value: throttle_position.value
    min: 0
    max: 100

  throttle_frame:
    use: percent_to_101_frame
    percent: throttle_percent

blocks:
  mood_image:
    inputs:
      mood: state
    elements:
      - type: threshold_image
        state: mood
        states:
          idle: rpm_normal
          cruise: rpm_normal
          push: rpm_high
          redline: rpm_redline
        x: 0
        y: 0
        z: 0

  number_on_glow:
    inputs:
      value: number
      format: string
      charset: charset
      glow_when: state
    elements:
      - type: image
        image: box_glow
        x: -10
        y: -10
        z: 0
        visible_when:
          state: glow_when
          equals: redline
      - type: sprite_text
        value: value
        format: format
        charset: charset
        x: 0
        y: 0
        z: 10

layers:
  - type: image
    id: background
    image: bg
    x: 0
    y: 0
    z: 0

  - type: group
    id: rpm_mood_art
    use: mood_image
    mood: rpm_mood
    x: 50
    y: 50
    z: 10

  - type: group
    id: rpm_value
    use: number_on_glow
    value: rpm.value
    format: "0000"
    charset: main_digits
    glow_when: rpm_mood
    x: 90
    y: 88
    z: 20

  - type: sprite_frame
    id: throttle_bar
    frame_set: throttle_weird_bar
    index: throttle_frame
    x: 80
    y: 320
    z: 20
```

This is probably the most important pattern if you want “some coding, but reusable”. The renderer can ship with decoder blocks, and later you can add custom registered decoders without changing dashboard configs that use them.

---

## Useful schema pressure-test notes

1. **Groups/blocks need named inputs.** Otherwise reuse dies quickly.
2. **Transforms should be reusable and composable.** `rpm -> percent -> frame` is cleaner than every sprite frame element repeating min/max logic.
3. **Overlays should not be special.** They are just later layers with conditions.
4. **Backgrounds are just early layers.** Day/night backgrounds can be conditional too.
5. **Validation matters.** Missing transform names, missing asset ids, or bad frame counts should fail before the dashboard starts.
6. **Keep OBD correctness outside this.** This schema should not calculate real sensor truth. It should only visualise latest state.
7. **The config can be huge.** That is fine. The engine should stay boring.

The guiding sentence:

> GoDriveLog dashboards are layered visual scenes driven by sensor state, reusable decoders, and reusable visual blocks.

That is flexible without becoming a full scripting language wearing a fake moustache.
