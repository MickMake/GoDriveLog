# v3.5 Implementation State

Status: v3.5.11 radial peg bounce implemented; v3.5.17+ radial visual/stat-marker tail prepared; deferred v3.5 realism slices remain pending

Current target: v3.5.17 radial needle visual polish

Current branch: `docs-v3.5-stat-markers-prompts`

## Scope

v3.5 is the gauge realism pass.

It adds believable gauge behaviour without changing the v3.4 gauge family model. The intent is to make gauges look more like real mechanisms when values change, while avoiding perpetual ambient effects.

The final v3.5 tail now includes small radial-only display refinements that need renderer/state support: optional needle drop shadow, optional display-only calibration offset, and rolling-window radial stat markers.

## Current decisions

- v3.5.12, v3.5.13, and v3.5.16 are temporarily deferred so the radial visual/stat-marker tail can be completed while the implementation context is fresh.
- The next-slice workflow should follow `Current target` when it is explicitly set, even if earlier unchecked slices remain.
- Most v3.5 realism options live under the `realism` key.
- `movement` is the exception: it is the single movement knob and should be accepted by any gauge type for now.
- Keep movement config collapsed as a scalar, not a nested object.
- Gauge types that do not yet have concrete movement behaviour may parse `movement` and use their current immediate behaviour until their movement slice defines more.
- Odometers use `odometer.movement` as the single source of truth for odometer wheel movement.
- Odometer `movement` supports `instant`, `linear`, `ease_out`, `bell`, `smooth`, and `click`.
- Odometer `instant` means digit display jumps immediately to the target value with no animation.
- Odometer `linear` means the wheel rolls from old digit position to target digit position at constant speed.
- Odometer `ease_out` means the wheel starts fast, then slows into the target.
- Odometer `bell` means the wheel starts slow, speeds up through the middle, then slows into the target.
- Odometer `smooth` is recognised only, reserved for future enhancement, and should warn then fall back to `instant`.
- Odometer `click` is recognised only, reserved for future stepped-click enhancement, and should warn then fall back to `instant`.
- `realism.movement_policy` is obsolete for odometer movement and must not be used or recommended for odometers.
- Existing top-level `movement` may remain supported for backwards compatibility where already present.
- Unknown movement values must fail configuration loading clearly unless a gauge type explicitly documents a recognised fallback.
- Unknown realism options must fail config loading.
- Known realism options used on unsupported gauge types must fail config loading.
- `realism.order` may optionally control the order of enabled realism behaviours.
- Do not rely on YAML key order to control behaviour order.
- Odometer movement should compose internally as `route -> lead_in -> travel -> settle -> rest`.
- The odometer phase model is internal implementation structure, not the public YAML shape.
- Do not expose `movement.pre`, `movement.primary`, or `movement.post` unless a later docs slice explicitly changes the public config model.
- `docs/v3.5/RealismBehaviourGuide.md` defines the intended visual feel of each realism option.
- Gauge Preview Mode is the simple visual viewer for one gauge at a time.
- Gauge Preview Mode CLI is `godrivelog dashboard preview <file>`.
- `<file>` is mandatory and positional.
- Preview files are normal YAML configs, not a special metadata system.
- Each single-feature preview file should enable one realism feature only.
- Each gauge type may also have one deliberate `99-all-options` preview file.
- Radial needle shadow is a static renderer feature, not dynamic parallax or lighting.
- Radial calibration offset is display-only and must not change input values.
- Radial `stat_markers` are display-only rolling-window markers and must not change input values, logs, exports, configured ranges, or source data.
- Radial `stat_markers` use the `realism.stat_markers` config key.
- `stat_markers.window` defines the trailing time range used to calculate enabled markers.
- `stat_markers.window: 0` means keep all stable displayed samples since runtime start.
- v3.5.18 implements radial `stat_markers.min` and `stat_markers.max` only.
- v3.5.19 implements radial `stat_markers.average` only.
- v3.5 stat marker assets are `needle_min.png`, `needle_max.png`, and `needle_average.png`.
- v3.5 stat markers are radial-only; bar gauge stat markers are a later feature unless a later slice explicitly promotes them.
- Hysteresis applies only to radial and bar gauges in v3.5.
- Indicator gauges support `thermal_fade` in v3.5.

## Approved v3.5 realism options

| Option | Applies to |
|---|---|
| `movement` | all gauge types for parsing; concrete behaviour defined per gauge type |
| `wraparound` | odometer |
| `drum_slop` | odometer |
| `carry_drag` | odometer |
| `snap_settle` | odometer |
| `backlash` | odometer |
| `hysteresis` | radial, bar |
| `stiction` | radial |
| `damping` | radial, bar |
| `overshoot` | radial, bar |
| `peg_bounce` | radial |
| `thermal_fade` | indicator |
| `needle_shadow` | radial |
| `calibration_offset` | radial |
| `stat_markers` | radial in v3.5; bar later |

## Scope boundaries

Allowed in v3.5:

- static imperfection;
- finite value-change movement;
- Gauge Preview Mode;
- deterministic, bounded behaviour;
- display-only realism options;
- small radial-only display refinements that need renderer support;
- rolling-window radial stat markers.

Not allowed in v3.5:

- idle needle vibration;
- random flicker or shimmer;
- gas-discharge jitter;
- LED multiplex flicker;
- power-on sweep;
- brownout dip;
- lazy power-off;
- dynamic parallax or gyro/light-driven shadow movement;
- general physics engine;
- generated artwork, videos, screenshot reports, or visual diff machinery.

## Slice checklist

- [x] v3.5.0 movement realism docs
- [x] v3.5.1 Gauge Preview Mode
- [x] v3.5.2 odometer wraparound
- [x] v3.5.3 odometer drum slop
- [x] v3.5.4 finite movement lifecycle
- [x] v3.5.5 shared movement policy groundwork
- [x] v3.5.6a document odometer movement goal / alignment
- [x] v3.5.6b implement odometer movement model
- [x] v3.5.7 odometer carry-drag / 9-drag
- [x] v3.5.8 radial damping
- [x] v3.5.9 radial stiction
- [x] v3.5.10 radial/bar overshoot
- [x] v3.5.11 radial peg bounce
- [ ] v3.5.12 indicator thermal fade
- [ ] v3.5.13 bar smoothing
- [x] v3.5.14 odometer snap / settle
- [x] v3.5.15 odometer backlash
- [ ] v3.5.16 display-only hysteresis
- [ ] v3.5.17 radial needle visual polish
- [ ] v3.5.18 radial stat markers min/max
- [ ] v3.5.19 radial stat marker average

## Next-slice workflow

When asked to do the next slice:

1. Read this file.
2. Use `Current target` if it is explicitly set; otherwise find the first unchecked allowed slice.
3. Read docs/v3.5/ReleasePlan.md.
4. Read docs/v3.5/RealismBehaviourGuide.md.
5. Read the matching prompt in docs/v3.5/prompts/.
6. Make only that slice’s changes.
7. Update this checklist and any relevant docs.
8. Do not implement later slices early.
9. Run the relevant local tests/checks.
10. Commit the completed slice with a clear message.
11. Push the branch to GitHub.
12. Raise a pull request against main.

Then enter the review-fix loop:

13. Wait for codex GitHub review feedback and CI results.
14. If CI fails, review requests changes, or unresolved review comments require code changes:
    * inspect the feedback;
    * make the smallest safe fixes only;
    * do not refactor unrelated code;
    * rerun relevant tests/checks;
    * commit and push the fixes.
15. Repeat the review-fix loop at most 3 times.

Stop when either:

* the PR exists;
* CI/checks are passing;
* there are no requested changes;
* there are no unresolved review comments requiring code changes;
* the PR is green and ready for human merge.

If the review-fix loop reaches 3 attempts, stop and leave a PR comment summarising:

* what was fixed;
* what remains unresolved;
* why it could not be safely completed automatically.

Do not merge the PR.
