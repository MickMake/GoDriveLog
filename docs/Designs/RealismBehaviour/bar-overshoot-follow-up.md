# Bar Gauge Overshoot Follow-Up

Index: 1

Area: `gauge/bar`, `realism.overshoot`, animation

Effort: 2-4 Codex hours

`v3.5.10` is currently being treated as the radial overshoot slice because the active prompt/spec only defines radial overshoot behaviour. Bar gauge overshoot remains approved as a follow-up idea, but it must not be pulled into the radial overshoot implementation by inference from older `radial/bar overshoot` wording.

Bar gauges should eventually support `realism.overshoot`, but this was intentionally left out of the radial overshoot slice to avoid ambiguous behaviour and accidental scope creep.

## Notes

- Display-only.
- Bounded pass-and-settle movement.
- Should compose cleanly with bar damping/smoothing.
- Do not copy radial behaviour blindly; bar movement has its own visual semantics.
- A bar overshoot should affect the displayed fill/level extent, not mutate source sensor values.
- Clamp final settled display to the real target/range after the overshoot tail completes.
- Consider vertical and horizontal bars, plus different origins, when defining the later prompt.
- Keep radial overshoot behaviour unchanged when this is implemented.
