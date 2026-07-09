# `backlash`

Applies to: odometer.

Status: **not implemented**.

## What it would do

When direction changes, the wheel would show a tiny amount of slack before movement takes up in the new direction.

## What it simulates in real life

Backlash is mechanical play in gears, shafts, or drive couplings. When a mechanism reverses direction, the driving part may move slightly before the driven part starts moving because the clearance has to be taken up first.

For an odometer, this would simulate worn or imperfect gearing where direction-change slack is visible before the number drum follows.

## Good result

The change of direction feels slightly mechanical without making the reading confusing.

## Bad result

The wheel loses numeric correctness, delays too long, or appears disconnected from the value change.

## Implementation note

Earlier v3.5 docs/checklists claimed `backlash` was implemented, but current `main` does not contain a `realism.backlash` config key or odometer runtime behaviour. Treat this as future work until deliberately implemented.
