# `wraparound`

Applies to: odometer.

Status: **implemented**.

## What it does

Odometer wheels roll cleanly through digit-strip boundaries, especially `9 -> 0` and `0 -> 9`.

## What it simulates in real life

A mechanical odometer drum is a continuous wheel, not ten disconnected images. When it passes from `9` to `0`, the strip continues around the drum rather than jumping across a flat image list.

This option simulates that continuous cylindrical path.

## Good result

A rollover looks like one continuous drum motion through the nearest boundary.

## Bad result

The wheel jumps, reverses unexpectedly, rolls the long way around, or briefly shows an impossible digit position.
