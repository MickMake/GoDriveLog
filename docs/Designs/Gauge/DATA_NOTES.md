# Data and interpretation notes

## Authority and scope

The authoritative source is [gauge_display_research_catalog_v0.2.json](<_data/gauge_display_research_catalog_v0.2.json>) together with the three v0.2 indexes. The Markdown is generated for navigation and should not silently diverge from those files.

Mechanism-first research catalogue; generic framebuffer/HDMI displays excluded.

## Counting rules

- **Gauge-group counts:** Primary gauge_group only; every gauge contributes to exactly one primary count. Alternate memberships are listed separately.
- **Global quirk counts:** Distinct gauges containing the canonical quirk; a gauge is counted once per quirk.
- **Group quirk counts:** Primary gauge_group only; within a group, each gauge is counted once per canonical quirk.
- **Compound statements:** A source phrase may support multiple canonical quirks when it contains multiple behaviours.

A single preserved source phrase may therefore appear on more than one canonical quirk page. That is deliberate: the normalisation separates behaviours without discarding the original wording.

## Citation granularity

Sources and image links in v0.2 are attached to each gauge entry, not to individual quirk phrases. Quirk pages list the gauge-level references associated with the evidence, but they must not be read as a machine-verifiable statement-to-source mapping.

## Date confidence

| Code | Meaning |
|---|---|
| `H` | Dated primary, official or museum support |
| `M` | Approximate adoption/common-use range from reliable overviews |
| `L` | Uncertain or strongly regional |

## Normalised quirk record

Each gauge quirk record may contain:

- `quirk`: canonical controlled-vocabulary label.
- `qualifiers`: optional concise detail that changes how the canonical behaviour applies.
- `source_phrases`: exact preserved wording from the source catalogue.

The canonical label is useful for code and indexing; the qualifier and source wording carry the fine detail. Throwing those away would make the taxonomy neat in the same way a wood chipper makes a bookshelf neat.
