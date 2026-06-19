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
      vehicle.go         # VehicleConfig, OBD endpoint config, selected log/dashboard IDs
      sensors.go         # Global SensorConfig catalogue and polling validation
      assets.go          # Global asset family config structs
      logs.go            # Global log subscriber config
      dashboard.go       # Global dashboard and widget config structs
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
      simple/
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
        warnings/
          engine_back.png
          engine_off.png
          engine_on.png
          engine_unknown.png
          engine_glass.png

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
        amber_digits/
          digit_back.png
          amber0.png
          amber1.png
          ...
          amber_minus.png
          amberDP.png
          digit_glass.png
        green_digits/
          digit_back.png
          green0.png
          green1.png
          ...
          green_minus.png
          greenDP.png
          digit_glass.png
        speed_bar/
          bar_back.png
          cell_off.png
          cell_green.png
          cell_yellow.png
          cell_red.png
          bar_glass.png
        status_bar/
          bar_back.png
          cell_off.png
          cell_green.png
          bar_glass.png
        rpm_boost/
          rpm_boost_back.png
          frame_000.png
          frame_001.png
          ...
          rpm_boost_glass.png
        warnings/
          engine_back.png
          engine_off.png
          engine_on.png
          engine_unknown.png
          engine_glass.png

      s2000/
        panel/
          main_panel_backplate.png
        speed_digits/
          digit_back.png
          speed0.png
          speed1.png
          ...
          speed_minus.png
          speedDP.png
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

      warnings/
        engine_back.png
        engine_off.png
        engine_on.png
        engine_unknown.png
        engine_glass.png

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
        simple_speed_warning.yaml
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
          repo-structure-guardrails.md

  testdata/
    config/
      valid/
        minimal.yaml
        full.yaml
        single-vehicle.yaml
        multi-vehicle-with-runtime-selection.yaml
      invalid/
        missing-vehicles.yaml
        multiple-vehicles-without-selection.yaml
        vehicle-log-not-found.yaml
        vehicle-dashboard-not-found.yaml
        selected-dashboard-display-collision.yaml
        log-sensor-not-found.yaml
        dashboard-sensor-not-found.yaml
        dashboard-asset-not-found.yaml
        bad-obd-address.yaml
        timeout-zero.yaml
        poll-zero.yaml
        unknown-root-field.yaml
        unknown-nested-field.yaml
        sensor-min-greater-than-max.yaml
        duplicate-widget-id.yaml
        decimal-format-missing-decimal-point.yaml
        bar-set-missing-off.yaml
        bar-widget-missing-on.yaml
        unsorted-bar-zones.yaml

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
- Vehicles are runtime profiles: they select the OBD endpoint, log definitions, and dashboard definitions.
- Sensors and assets are global catalogues; vehicles do not directly list sensors or assets.
- Active v3 examples should validate against the same schema rules as `config.example.yaml` and `config.full.yaml`.
- Asset paths in active v3 configs are repository-root relative.
- `ImplementationGuardrails.md` is the implementation checklist for writing v3 code against the target model.
- `MigrationGuardrails.md` explains how to move from current code to the v3 target without leaking current assumptions into the v3 core.
- `PerformanceGuardrails.md` gives the current display speed work a safe lane without warping the v3 schema.
- Treat the documented v3 root sections as an allow-list; unknown fields should fail validation at every documented level during v3 implementation.
- Use an OBD-like endpoint address for both real hardware and bench/simulator endpoints.
- Sensor timing is `poll`; selected logs and dashboards subscribe to events.
- `multi-vehicle-with-runtime-selection.yaml` describes a valid multi-vehicle config only when the test/runtime supplies explicit vehicle selection; that selection is not encoded in config.
- Display collision validation applies to the dashboards selected by the selected vehicle, not every dashboard definition in the file.
- Dashboard config is widget-based, not decoder/block/layer/condition based.
