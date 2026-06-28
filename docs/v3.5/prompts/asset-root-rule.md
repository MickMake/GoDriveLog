# Prompt: gauge asset-root rule

This is a separate follow-up task. Do not mix it with odometer movement work.

## Goal

Relax the gauge package asset-root rule without changing gauge movement behaviour.

## Desired rule

- A gauge package must live somewhere under an `assets` directory.
- Referenced assets must resolve under the same `assets` directory.
- Gauge packages do not have to be specifically under `assets/gauges`.

## Implementation note

If the code still has the old `assets/gauges` hard restriction, update asset-root discovery so it returns the nearest ancestor named `assets` rather than requiring the package path relative to `assets` to start with `gauges/`.

## Do not

- Do not change odometer movement.
- Do not change `movement` parsing or movement behaviour.
- Do not implement v3.5.7 or later gauge realism slices.
- Do not add asset-generation work.

## Validation

Run the existing dashboard/gauge package loading tests and any preview loading checks that cover gauge package asset paths.
