# `thermal_fade`

Applies to: indicator.

Status: **implemented**.

## What it does

An indicator turns on and off with a soft incandescent-style response.

## What it simulates in real life

Incandescent bulbs do not appear and disappear instantly. The filament warms up, glows, then cools down after power is removed. The off transition can feel different from the on transition.

This option simulates that warm-up and cool-down response. It does not simulate random flicker, bloom, lens dirt, weak bulb tint, or ageing. Those are artwork or future display-layer concerns.

## Good result

On-state appears to warm in rather than appearing instantly. Off-state fades away softly and then settles fully off.

## Bad result

The indicator flickers, pulses, randomly changes brightness, or remains partly on after it should be off.
