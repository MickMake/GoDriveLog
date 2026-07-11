# `peg_bounce`

Applies to: radial, bar.

Status: **implemented**.

## What it does

When the displayed position reaches the configured minimum or maximum stop, it should appear to tap the stop, rebound slightly, then settle quickly.

For radial gauges, this appears as a needle tapping a min/max physical peg. For bar gauges, this appears as an end-stop bounce on the displayed fill/reveal extent. The config key remains `realism.peg_bounce` even though bars do not have literal pegs.

## What it simulates in real life

Many analogue gauges have physical stop pegs or hard limits. When a needle reaches the stop with momentum, it may make a tiny rebound before settling. A bar gauge does not literally hit a peg, but the same idea can apply visually to the fill edge or reveal extent reaching its display limit.

## Good result

The bounce is small, quick, deterministic, and only visible at the configured visual stop.

## Bad result

The display bounces during ordinary in-range movement, passes through the stop, keeps bouncing, or changes source values.
