GoDriveLog/
  cmd/
    GoDriveLog/
      main.go

  internal/
    config/
      config.go
      vehicle.go
      sensors.go
      log.go
      display.go
      dashboard.go
      load.go
      validate.go
      resolve.go
      errors.go

    logger/
      jsonl.go

    sensors/
      reader.go
      cache.go
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

    ui/
      app.go
      display.go

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
            ...
          green/
            0.png
            1.png
            ...
        throttle/
          frame_000.png
          frame_001.png
          ...

  configs/
    config.example.yaml
    config.full.yaml
    examples/
      one-display.yaml
      two-displays.yaml
      race-driver.yaml
      dashboard-minimal.yaml

  docs/
    v3/
      README.md
      config-schema.md
      config-go-structs.md
      repo-layout.md
      config-validation.md
      config-testing.md
      decisions.md
      CHANGES.md

    archive/
      dashboard-v2/
        README.md
        overview.md
        prompts.md
        reference.md
        current-status.md
        decisions.md
        human-process.md
        CHANGES.md

  testdata/
    config/
      valid/
        minimal.yaml
        full.yaml
        one-display.yaml
        two-displays.yaml
        race-driver.yaml

      invalid/
        missing-default-vehicle.yaml
        default-vehicle-not-found.yaml
        active-display-not-found.yaml
        display-dashboard-not-found.yaml
        log-sensor-not-found.yaml
        dashboard-sensor-not-found.yaml
        dashboard-asset-not-found.yaml
        dashboard-decoder-not-found.yaml
        unknown-field.yaml
        bad-obd-source.yaml
        cache-zero.yaml
        refresh-zero.yaml

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