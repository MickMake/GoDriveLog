# GoDriveLog Renderer v2 Improvement Plan

Version: v2.0.x planning pack  
Project: GoDriveLog  
Purpose: Replace the old config-driven scene renderer as the main runtime display path with a faster fixed instrument dashboard renderer.

## 1. Direction

The new renderer work should be treated as the primary display direction for GoDriveLog.

The existing config-driven dashboard renderer is not to be preserved as a normal runtime preference. If a fallback is needed, use Git history and a tag, not long-term production code.

Recommended safety tag before starting:

```bash
git checkout main
git pull
git tag legacy-config-scene-dashboard
git push origin legacy-config-scene-dashboard
```

Optional version tag:

```bash
git tag v1.x-legacy-config-scene-dashboard
git push origin v1.x-legacy-config-scene-dashboard
```

Once tagged, move forward with the fast renderer decisively.

## 2. Why change renderer architecture?

The current dashboard renderer is flexible but too generic for a live vehicle dashboard. It performs this style of pipeline:

```text
sensor state -> decoders -> scene evaluation -> renderer diff/cache -> Fyne canvas object updates
```

That is useful for layout experiments, but it is overbuilt for a production-style instrument display that needs frequent, predictable updates.

The new dashboard should use:

```text
sensor state -> instrument model -> direct updates to fixed Fyne canvas objects
```

The key difference is that the new renderer creates the object graph once and then mutates existing `canvas.Text`, `canvas.Rectangle`, and `canvas.Image` objects directly.

## 3. Target release sequence

### v2.0.0 — Fast renderer skeleton

Purpose: establish the new fast instrument display path.

Deliver:
- Tag old renderer baseline before work starts.
- Add new 1920x480 fast instrument dashboard.
- Use fixed pre-created Fyne objects.
- Display basic race-demo fields.
- Make the new fast renderer the primary app display path.
- Do not preserve the old renderer as a user-facing preference.

### v2.0.1 — Visual dashboard polish

Purpose: make the display visually useful and readable.

Deliver:
- Proper 1920x480 dashboard layout.
- Large RPM and speed displays.
- Throttle, oil temp, oil pressure, coolant, gear, warning, failure, reset states.
- Obvious warning/critical/failure states.
- Yellow/green/red dashboard styling where practical.
- 10 Hz or better perceived update behaviour.

### v2.0.2 — Cleanup and old renderer isolation/removal

Purpose: prevent codebase pollution.

Deliver:
- Remove old config-scene renderer from normal runtime.
- Remove or isolate old dashboard docs/config references.
- Keep fallback via Git tag, not production code.
- Remove misleading config/preferences.
- Keep only useful shared primitives.

### v2.0.3 — Optional asset/performance refinement

Purpose: tune on target hardware after the main renderer works.

Deliver:
- Profile on target hardware.
- Optimize SVG/image loading if needed.
- Introduce sprite digit assets or raster drawing if beneficial.
- Tune refresh cadence.
- Add instrumentation/debug overlay if useful.

## 4. Global guardrails

These apply to every pass.

Do not:
- Change live OBD sensor behaviour unless the pass explicitly requires it.
- Modify `RaceDemoScenario` unless a real missing output blocks the display.
- Duplicate `RaceDemoScenario`.
- Rebuild the whole app architecture.
- Preserve old renderer as an equal runtime preference.
- Add config complexity just to keep the old way alive.
- Touch unrelated logging, OBD transport, asset generation, invoice/business files, or unrelated projects.
- Rename broad packages unnecessarily.
- Refactor unrelated code “while here”.
- Add new dependencies unless clearly justified.

Do:
- Keep each pass focused.
- Prefer additive changes until the new renderer is proven.
- Keep the fast display path simple and explicit.
- Update only the files needed for that pass.
- Run `go test ./...`.
- State clearly in PR body what was and was not changed.
- Treat live OBD behaviour as protected.

## 5. Primary design

The fast renderer should look like this conceptually:

```go
type InstrumentDashboard struct {
    root fyne.CanvasObject

    rpmText         *canvas.Text
    speedText       *canvas.Text
    throttleBar     *canvas.Rectangle
    oilTempText     *canvas.Text
    oilPressureText *canvas.Text
    gearText        *canvas.Text
    warningText     *canvas.Text
    failureOverlay  *canvas.Rectangle

    store *sensors.StateStore
}
```

It should:
- Construct all canvas objects once.
- On each update tick, read current state snapshot.
- Update only changed object fields.
- Call `Refresh()` only on changed objects or the smallest necessary container.
- Avoid scene graph rebuilding.
- Avoid generic block traversal.
- Avoid signature comparison.
- Avoid render throttling for the instrument dashboard.

## 6. Display target

Canvas:

```yaml
width: 1920
height: 480
```

Suggested layout:
- Left: RPM and throttle.
- Centre: speed and gear.
- Right: oil temp, oil pressure, coolant, warning/failure/reset.
- Bottom strip: status line / phase / alert message if available.
- Warning and failure states must be visually obvious.

## 7. Acceptance summary

The work is successful when:
- Race demo launches into a 1920x480 fast dashboard.
- The display visibly changes at useful cadence.
- RPM, speed, throttle, oil temp, oil pressure, gear, warning, failure, and reset state are readable.
- The thrown rod event is obvious.
- Live OBD behaviour remains protected.
- The old renderer is no longer the preferred/main runtime path.
