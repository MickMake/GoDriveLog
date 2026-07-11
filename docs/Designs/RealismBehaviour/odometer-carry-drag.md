# `carry_drag`

Applies to: odometer.

Status: **implemented**.

## What it does

A higher digit begins to creep as the lower digit approaches rollover.

## What it simulates in real life

In a mechanical odometer, the lower drum does not always leave the next drum perfectly untouched until the exact rollover point. The carry mechanism can start to load or nudge the neighbouring drum before the final click.

This option simulates light rollover coupling between adjacent number drums.

## Good result

The next wheel looks lightly dragged by the rolling lower wheel, then lands in the correct final digit.

## Bad result

The higher digit moves too early, moves too far, or appears to change value before the rollover is visually justified.
