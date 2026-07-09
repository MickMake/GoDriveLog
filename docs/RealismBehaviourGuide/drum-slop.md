# `drum_slop`

Applies to: odometer.

Status: **implemented**.

## What it does

Each odometer wheel may sit with a small fixed alignment imperfection, as if the drums are not perfectly centred in the window.

## What it simulates in real life

Real mechanical odometer drums are not always perfectly indexed. Wear, manufacturing tolerance, gear lash, and imperfect assembly can leave each wheel sitting a fraction high or low in the viewing window.

This option simulates static per-wheel alignment imperfection. It is not movement, bounce, drift, or direction-change slack.

## Good result

Digits look slightly mechanical and imperfect while still being readable.

## Bad result

Digits become hard to read, offsets change between runs, or wheels drift while idle.
