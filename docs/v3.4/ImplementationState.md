# GoDriveLog v3.4 implementation state

Status: v3.4.0 planning docs branch in progress
Current target: v3.4 gauge/display package cleanup and expansion
Current branch: v3.4.0-gauge-type-docs

## Purpose

This file records the current implementation state for v3.4. Update it in every v3.4 slice PR.

## Gauge type decision

The gauge type direction is:

```text
numeric    = formatted value rendered through image character slots
radial     = value-to-angle needle/arc gauge
odometer   = rolling wheel gauge
indicator  = off/on state gauge
bar        = continuous fill/reveal/movement gauge
segmented  = stepped percent-threshold image gauge
```

Visual identity belongs to assets. Code should model behaviour only.

## Non-goals

- No `style` field.
- No `seven_segment` compatibility alias.
- No dot-matrix font/text renderer in this line.
- No merged `bar`/`segmented` supertype.
- No eager loading of all `segmented` percent images.

## Numeric rename

`seven_segment` is planned to become `numeric`.

The rename is intentionally a hard rename. This project does not need a compatibility layer for old local gauge YAML. If something breaks, it is cheaper to fix the package than to keep a small museum of aliases.

## Segmented percent model

`segmented` value layers use `{percent}`:

```yaml
layers:
  segments: levels/rpm_{percent:03}.png
```

The renderer discovers files such as:

```text
rpm_000.png
rpm_010.png
rpm_030.png
```

Those files are valid sparse percent thresholds. The renderer selects the highest discovered percent less than or equal to the current normalized value.

Discovery counts filenames only. Image decoding must stay lazy.

## Baseline dashboard

The v3.4 baseline remains conceptually based on the reusable baseline config:

```text
examples/baseline-dashboard.yaml
```

The current baseline workload remains useful because it exercises numeric displays and radial RPM through the active Ebiten path.

## Completed slices

| Version | Status | Notes |
|---|---|---|
| v3.4.0 | in progress | Planning docs and prompt set for gauge type cleanup and expansion. |

## Pending slices

| Version | Status | Next action |
|---|---|---|
| v3.4.1 | not started | Rename `seven_segment` to `numeric` in code and examples. |
| v3.4.2 | not started | Add odometer config/scene model. |
| v3.4.3 | not started | Add indicator gauge behaviour. |
| v3.4.4 | not started | Add first bar gauge behaviour. |
| v3.4.5 | not started | Add segmented percent-threshold image discovery and rendering. |

## Update rule

Every v3.4 implementation PR must update this file with:

- completed version;
- current branch;
- next target;
- any changed decisions or deferrals.
