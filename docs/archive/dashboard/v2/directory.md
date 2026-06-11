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

