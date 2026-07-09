# `stiction`

Applies to: radial, bar.

Status: **implemented**.

## What it does

The displayed position resists very small changes, then releases once the change is large enough to notice.

For radial gauges, the needle may hold briefly before making a catch-up movement. For bar gauges, the displayed fill/reveal extent may hold briefly before making a catch-up movement.

## What it simulates in real life

Stiction is static friction: the extra resistance that must be overcome before a resting mechanism starts moving. A sticky needle, linkage, bearing, or sliding display can ignore tiny input changes until enough force builds up to break it free.

This option simulates that thresholded release.

## Good result

Tiny value changes may not move the display immediately. When it does move, it makes a small catch-up movement and settles.

## Bad result

The display sticks during large changes, jumps violently, behaves unpredictably, or changes source values.
