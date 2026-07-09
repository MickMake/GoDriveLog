# `ghosting`

Applies to: numeric, segmented.

Status: **candidate / needs design**.

## What it would do

A previous displayed character may remain faintly visible for a short, bounded time after the display changes.

## What it simulates in real life

Ghosting appears in many real display technologies:

- LCD segments can retain a faint previous state;
- multiplexed LED displays can leave perceived persistence;
- VFD or gas-discharge elements may fade rather than cut instantly;
- aged displays may have slow decay or residual glow.

This option simulates short-lived previous-character persistence.

## Candidate visual model

When a digit slot changes from one character to another:

```text
previous glyph fades out
new glyph appears normally or fades in
```

The old glyph must disappear completely after the bounded ghosting duration.

## Good result

The display feels like a real physical/electronic display with mild persistence, while the current reading remains clear.

## Bad result

Ghosting makes the value unreadable, stacks many previous values, never fully clears, or turns into random flicker.

## Design notes

Ghosting should operate at digit-slot or character-slot level. Avoid requiring the renderer to infer segments from arbitrary glyph images unless a future display-mask abstraction exists.
