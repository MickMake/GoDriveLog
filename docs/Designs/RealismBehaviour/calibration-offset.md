# `calibration_offset`

Applies to: radial.

Status: **implemented**.

## What it does

A radial needle may display with a small fixed angular offset, as if the needle is not perfectly aligned.

## What it simulates in real life

Real analogue gauges are not always perfectly calibrated. The needle might be installed slightly off, the mechanism may have a small fixed bias, or the dial artwork and pointer may not line up perfectly.

This option simulates a fixed display-only calibration imperfection. It must not change the source value or configured value range.

## Good result

The gauge looks slightly imperfect while still clearly representing the configured value range.

## Bad result

The offset changes source values, pushes the needle outside sensible visual bounds, or makes the gauge look broken.
