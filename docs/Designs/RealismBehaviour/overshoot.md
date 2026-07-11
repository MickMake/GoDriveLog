# `overshoot`

Applies to: radial, bar.

Status: **implemented**.

## What it does

The displayed position may pass the target slightly, then return and settle.

For radial gauges, this is angular needle overshoot. For bar gauges, this is fill/reveal extent overshoot.

## What it simulates in real life

Overshoot simulates momentum in a moving indicator. A needle, linkage, or damped mechanism can carry a little past the final reading before returning to rest. In a bar display, it simulates the displayed fill edge carrying slightly past the target extent before settling.

This should look like a small mechanical or damped response, not a cartoon spring.

## Good result

The movement gives a small sense of momentum without stealing attention.

## Bad result

The display swings too far, oscillates repeatedly, overshoots during tiny changes where it looks silly, or renders outside sensible visual bounds.
