# GoDriveLog generated v3.4 example assets

These assets are the committed output of the deterministic v3.4 example generator.

Source of truth:

- generator entry point: `go run ./scripts/generate-example-assets`
- committed themes: `framework-smoke`, `ornate-timber`
- dashboard configs using this tree: `examples/dashboards/*.yaml`

Regenerate from the repository root:

```bash
go run ./scripts/generate-example-assets -theme framework-smoke
go run ./scripts/generate-example-assets -theme ornate-timber
go run ./scripts/generate-example-assets -theme all
```

Conventions:

- generated theme assets live under `examples/assets/v3.4/<theme>/`
- generated dashboard configs live under `examples/dashboards/`
- generated PNGs are committed for review, but the script is the source of truth
- do not hand-edit generated PNGs unless you also update the generator

This directory is the workshop output rack, not the place where the chisels live.
