# Bar Realism Scope

Origin: `docs/v3.7/PlannedFeatures.md`

Bar gauges are linear fill/reveal gauges. Runtime realism should focus on the displayed fill edge moving toward the target, not on repainting the gauge artwork.

Before planning any bar realism beyond pointer markers, audit the current code and the completed v3.5 docs. Do not contradict completed v3.5 state from backlog notes alone.

Possible future bar candidates:

- `stepped_fill` for block-style bars;
- `quantized_fill` where the bar only visibly changes after the value crosses a display-resolution step;
- focused audits/fixes for already-documented bar realism options if code support is missing.

Both `stepped_fill` and `quantized_fill` need a clear config model before promotion.
