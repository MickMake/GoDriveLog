# GoDriveLog v3 directory structure

This is the intended v3 repo layout. Some runtime files may still differ while the codebase catches up to the documentation.

```text
GoDriveLog/
  cmd/
    GoDriveLog/
      main.go

  internal/
    config/
      config.go          # Root Config: vehicles, sensors, assets, logs, dashboards
      vehicle.go         # VehicleConfig and OBD endpoint config
      sensors.go         # SensorConfig and polling validation
      assets.go          # Asset family config structs
      logs.go            # Log subscriber config
      dashboard.go       # Dashboard and widget config structs
      load.go            # YAML loading
      validate.go        # Config validation
      resolve.go         # Runtime reference resolution
      errors.go

    vehicle/
      endpoint.go        # OBD-like endpoint abstraction
      elm327.go          # Serial/TCP ELM327-style endpoint support

    sensors/
      runtime.go         # Sensor polling runtime
      reader.go          # Sensor reader abstraction
      event.go           # SensorEvent and SensorStatus definitions
      state.go           # Latest known sensor state
      registry.go        # Sensor registration/resolution

    logger/
      jsonl.go           # JSONL event subscriber

    dashboard/
      assets/
        registry.go      # Asset registry and validation
        image.go         # Image asset loading
        digit_set.go     # Digit/character image asset sets
        bar_set.go       # Repeated cell bar assets
        frame_set.go     # Frame sequence assets
        indicator_set.go # On/off/unknown indicator assets

      widgets/
        image.go         # Static image widget
        digit_display.go # Formatted character display widget
        bar_display.go   # Cell bar widget
        frame_gauge.go   # Frame sequence gauge widget
        indicator.go     # Boolean/status indicator widget

      renderer/
        fyne/
          renderer.go
          image.go
          group.go

    ui/
      app.go
      display.go

  assets/
    dashboard/
      bttf/
        panel/
          background.png
        amber_digits/
          digit_back.png
          amber0.png
          amber1.png
          ...
          amber_minus.png
          amber_dp.png
          digit_glass.png
        green_digits/
          digit_back.png
          green0.png
          green1.png
          ...
          green_minus.png
          green_dp.png
          digit_glass.png
        throttle/
          throttle_back.png
          frame_000.png
          frame_001.png
          ...
          throttle_glass.png

      z31/
        panel/
          main_panel_backplate.png
        speed_bar/
          bar_back.png
          cell_off.png
          cell_green.png
          cell_yellow.png
          cell_red.png
          bar_glass.png
        rpm_boost/
          rpm_boost_back.png
          frame_000.png
          frame_001.png
          ...
          rpm_boost_glass.png

      s2000/
        panel/
          main_panel_backplate.png
        speed_digits/
          digit_back.png
          speed0.png
          speed1.png
          ...
          speed_minus.png
          speed_dp.png
          digit_glass.png
        rpm/
          rpm_back.png
          rpm_000.png
          rpm_001.png
          ...
          rpm_glass.png
        side_bar/
          bar_back.png
          cell_off.png
          cell_on.png
          cell_warning.png
          bar_glass.png

      warnings/
        engine_back.png
        engine_off.png
        engine_on.png
        engine_unknown.png
        engine_glass.png
        highbeam_back.png
        highbeam_off.png
        highbeam_on.png
        highbeam_unknown.png
        highbeam_glass.png

  docs/
    v3/
      README.md
      config.example.yaml
      config.full.yaml
      GoStructsConfig.md
      DirectoryStructure.md
      ImplementationGuardrails.md
      MigrationGuardrails.md
      PerformanceGuardrails.md
      examples/
        README.md
        nissan_300zx_z31_inspired.yaml
        honda_s2000_inspired.yaml

    archive/
      config.md
      dashboard/
        v2/
          README.md
          overview.md
          reference.md
          current-status.md
          decisions.md
          human-process.md
          prompts.md
          repo-structure-guardrails.md

  testdata/
    config/
      valid/
        minimal.yaml
        full.yaml
        single-vehicle.yaml
        multi-vehicle-explicit.yaml
      invalid/
        missing-vehicles.yaml
        multiple-vehicles-without-selection.yaml
        log-sensor-not-found.yaml
        dashboard-sensor-not-found.yaml
        dashboard-asset-not-found.yaml
        bad-obd-address.yaml
        poll-zero.yaml
        unknown-field.yaml

    dashboard/
      assets/
        tiny.png
        digits/
        bars/
        frames/
        indicators/
      configs/
        valid-minimal.yaml
        invalid-missing-asset.yaml
        invalid-missing-character.yaml
        invalid-indicator-state.yaml

  go.mod
  go.sum
  README.md
```

## Notes

- Archive docs are allowed to describe old behaviour.
- Active v3 docs should use the simplified top-level shape: `vehicles`, `sensors`, `assets`, `logs`, `dashboards`.
- `ImplementationGuardrails.md` is the implementation checklist for writing v3 code against the target model.
- `MigrationGuardrails.md` explains how to move from current code to the v3 target without leaking current assumptions into the v3 core.
- `PerformanceGuardrails.md` gives the current display speed work a safe lane without warping the v3 schema.
- Treat the documented v3 root sections as an allow-list; unknown root fields should fail validation during v3 implementation.
- Use an OBD-like endpoint address for both real hardware and bench/simulator endpoints.
- Sensor timing is `poll`; logs and dashboards subscribe to events.
- Dashboard config is widget-based, not decoder/block/layer/condition based.
