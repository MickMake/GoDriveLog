# GoDriveLog v3.1 implementation prompts

Use one prompt file per implementation slice.

Planning work is not a numbered v3.1 slice. Implementation starts at `v3.1.0`.

## Prompt index

- `v3.1.0-runnable-command-path.md`
- `v3.1.1-display-adapter.md`
- `v3.1.2-dashboard-gauge-test-harness.md`
- `v3.1.3-dashboard-update-performance.md`
- `v3.1.4-jsonl-rotation-decision.md`
- `v3.1.5-typed-sensor-values.md`
- `v3.1.6-unsupported-missing-semantics.md`
- `v3.1.7-dashboard-event-efficiency.md`
- `v3.1.8-retirement-readiness.md`

## v3.1 implementation focus

The main v3.1 goal is to wire the existing v3 foundation into a practical app path.

The preferred pipeline remains: selected vehicle -> endpoint -> sensor polling runtime -> events -> selected logs and dashboards.
