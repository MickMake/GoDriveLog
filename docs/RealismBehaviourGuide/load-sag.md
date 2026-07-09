# `load_sag`

Applies to: numeric, segmented.

Status: **candidate / not implemented**.

## What it would do

The display would dim slightly when the currently shown value requires more lit segments or higher display load.

## What it simulates in real life

Some displays draw more current when more segments are lit. If the driver, supply, wiring, or display technology is marginal, high-load values may appear slightly dimmer.

For a seven-segment display, `888` can draw much more segment current than `111`, so the whole display or individual slots may sag in brightness.

## Candidate visual model

| Displayed value | Approximate lit segment load | Expected brightness |
|---|---|---|
| `111` | low | brighter |
| `777` | medium | normal-ish |
| `888` | high | dimmer |
| `888.8` | very high | dimmest |

Possible modes:

| Mode | Meaning |
|---|---|
| `display` | Whole display dims based on total lit-character load. |
| `slot` | Each digit slot dims based on that character's own load. |
| `hybrid` | Whole display sag plus slight per-slot sag. |

## Good result

Values with heavier segment load look slightly more strained, while the display remains readable and deterministic.

## Bad result

The display pumps brightness dramatically, changes numeric meaning, requires image-internal segment detection, or flickers randomly.

## Design notes

Prefer starting with a display-level model because it best matches a shared supply or driver sag. Character load can be configured or inferred from a known glyph/segment table rather than image internals.
