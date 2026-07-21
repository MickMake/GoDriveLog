# Overview

## Purpose

GoDriveLog is heading toward a **gauge behaviour simulation platform** rather than a growing collection of one-off gauge implementations.

The intent is still to simulate the **physical behaviour** of gauges. What is changing is the practical method used to get there.

Instead of treating each real-world gauge as a separate implementation target, the system should be built around:

- a small set of **gauge types**
- configurable **assets** that define the visible look
- configurable **quirks** that define the behaviour or feel

This keeps the architecture grounded in the existing vocabulary and avoids turning the project into a large, fragile collection of special cases.

## Core idea

A final gauge is formed from four parts:

1. **Gauge type**  
   The physical mechanism or topology being simulated.

2. **Assets**  
   The visible presentation supplied to that gauge type.

3. **Quirks**  
   The behaviour parameters attached to that gauge type.

4. **Runtime code**  
   The logic that applies the quirks over time and produces the final displayed behaviour.

In short:

```text
Gauge Type + Assets + Quirks + Runtime Logic = Final Gauge
```

## Mantra

The project should follow a simple rule:

> **Assets give the gauge its visible look. Code gives it its behaviour or feel.**

That rule is important because it prevents assets from becoming a hidden behaviour system.

- Assets should describe visible things such as faces, bezels, overlays, masks, drum art, segment art, glow art, glass effects, and similar presentation layers.
- Quirks should describe behaviour such as lag, slop, carry drag, backlash, ghosting, thermal fade, load sag, and related display behaviour.

## What this is not

This should **not** be described as a brand new architecture direction.

It is better understood as **honing in on the intended architecture**:

- keeping the current vocabulary
- separating research from runtime design
- collapsing duplicated behaviour ideas into common runtime mechanisms
- making the platform configurable enough to simulate many real gauges without implementing each one separately

## Why this matters

The project has research covering many real gauges and physical mechanisms, but the real-world space is vastly larger than the current sample.

That means the code cannot reasonably be organised as:

- one gauge researched
- one implementation written
- repeat forever

That approach would take too long and produce brittle code.

The better approach is to treat the research as a way to discover:

- the major **gauge types**
- the major **quirk vocabularies**
- the major **behaviour patterns**
- sensible parameter ranges and presets

The runtime platform then uses those discoveries to simulate a much wider behaviour space.

## Separation of concerns

There are now two distinct concerns that should stay separate.

### 1. Gauge research

This is the catalogue of real physical gauge mechanisms, observed quirks, examples, and reference material.

This should live separately, for example under a directory such as:

```text
GaugeResearch/
```

That material is reference data. It should not be treated as the runtime architecture.

The research Markdown files are expected to be generated from JSON later. That generation process is outside the scope of these docs.

### 2. Runtime design

This is the code-facing design for how GoDriveLog simulates gauges.

This includes:

- what a gauge type means in the runtime
- what assets are expected
- how quirks are configured
- how multiple quirks work together
- how behaviour is applied consistently across supported gauge types

## Short example

```text
Gauge type: rolling_drum_counter
Assets: cream drum digits, worn bezel, dusty glass
Quirks: lag, drum_slop, carry_drag, wraparound
```

This means:

- the **gauge type** defines the physical mechanism
- the **assets** define what it looks like
- the **quirks** define how it feels when driven

The same runtime platform can then be used to make that gauge feel tight, worn, lazy, twitchy, or sloppy simply by adjusting quirk values.

## Direction from here

The next design steps should stay simple:

- keep gauge types small and stable
- keep assets presentation-only
- keep quirks as the behaviour vocabulary
- define how quirks combine safely
- document a minimal plan for evolving the runtime without a giant rewrite
