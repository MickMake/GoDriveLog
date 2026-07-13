# v3.7 Release Plan

Status: planned slice list complete

v3.7 follows the completed v3.6 pointer marker release.

v3.7 is the odometer and image-based numeric display realism release. It adds the remaining planned odometer `backlash` behaviour and a focused set of realistic effects for the existing image-based numeric/seven-segment renderer.

The planned v3.7 slice list is the release contract. Keep each slice small, extend the existing renderer directly, and move unrelated realism ideas to v3.8+ or a separate follow-up issue/PR.

## Theme

Complete the selected realism work using the KISS principle.

v3.7 extends existing renderers. It does not replace or redesign them.

Every realism option should own its implementation unless sharing a small helper is obviously simpler than keeping the code local. Do not create a generic realism engine, a new display abstraction, or a shared runtime framework merely because multiple options happen to affect digits.

All v3.7 realism remains:

- optional;
- display-only;
- deterministic;
- bounded and subtle;
- disabled unless configured;
- non-mutating for source values, logs, exports, configured ranges, or input data.

## Existing renderer ownership

The existing numeric renderer remains responsible for:

- formatting the source value;
- splitting formatted text into digit slots;
- assigning decimal points;
- resolving character assets;
- composing background, character, decimal-point, foreground, and overlay layers.

v3.7 adds realism effects to that existing flow. It must not introduce a replacement numeric renderer.

The existing whole-image `segmented` gauge renderer remains separate. Its existing hysteresis behaviour and threshold-image selection are not to be moved into a new shared numeric-display system.

## Scope

### Odometer

- `backlash`

### Image-based numeric / seven-segment display

- `per_digit_response_lag`
- `leading_zero_behaviour`
- `segment_bleed`
- `digit_bleed`
- `ghosting`
- `uneven_brightness`
- `load_sag`

## Odometer `backlash`

`backlash` applies only to odometer gauges.

When an odometer value reverses direction, wheel movement may show a small bounded amount of mechanical slack before following the new direction and settling exactly on the correct rendered target.

Rules:

- disabled unless configured;
- applies only when direction changes;
- remains bounded and subtle;
- must settle exactly on the correct target;
- must not alter source values or exported values;
- must not interfere with forward-only movement;
- must continue to work with existing odometer movement options such as wraparound, carry drag, snap settle, and drum slop.

## Numeric realism options

### `per_digit_response_lag`

Digit slots may update with small controlled timing differences instead of changing at the exact same instant.

Rules:

- slot order and delay must be deterministic;
- delays must remain short enough that the display stays readable;
- all slots must settle on the correct current value;
- do not create random update order;
- do not alter formatting or source values;
- keep any required state local to this feature unless another completed feature clearly benefits from a tiny shared helper.

### `leading_zero_behaviour`

Leading zero slots may be presented deliberately rather than only as a formatting side effect.

Supported presentation modes should follow the realism guide and may include:

- `show`;
- `blank`;
- `dim`;
- `placeholder`, where the final config design defines a clear asset or display rule.

Rules:

- presentation only;
- preserve significant zeroes;
- do not change the numeric meaning;
- do not change exported or formatted source values;
- keep slot alignment stable.

### `segment_bleed`

Inactive seven-segment shapes may remain faintly visible beneath the active character image, simulating inactive segment visibility through a lens, mask, or display material.

Rules:

- use the supplied image-based display assets;
- do not inspect arbitrary pixel content to infer segment geometry;
- keep bleed subtle;
- render deterministically;
- preserve the active character as the dominant readable image.

### `digit_bleed`

An inactive digit-slot image or mask may remain faintly visible beneath the active digit.

Rules:

- keep the effect slot-local;
- use explicit package assets or an existing documented layer where possible;
- do not synthesize a new procedural digit system;
- keep the active digit readable;
- do not affect blank slots unless configured.

### `ghosting`

A previous character image may fade out briefly after a digit changes instead of disappearing instantly.

Rules:

- previous-character state belongs to the ghosting feature;
- fade timing must be finite and deterministic;
- ghosting must fully settle;
- do not retain unbounded history;
- do not make long-lived mixed values unreadable;
- do not require unrelated numeric gauges to maintain runtime history.

### `uneven_brightness`

Digit slots may render with small stable brightness differences.

Rules:

- start with per-slot brightness, not per-segment analysis;
- variation must be deterministic and stable for a package/slot;
- do not vary randomly each frame;
- keep the range subtle and readable;
- active character, decimal-point, bleed, and ghost layers should only share brightness behaviour where the feature design explicitly requires it.

### `load_sag`

The display may dim slightly when the currently displayed value represents a heavier segment load.

Rules:

- start with a display-level model unless the slice proves a slot-level model is simpler;
- use an explicit known character-to-segment-load table or configured load values;
- do not inspect image pixels to count lit segments;
- keep brightness changes subtle;
- avoid visible pumping or flicker;
- remain deterministic;
- do not change the displayed characters or numeric meaning.

## Rendering rules

- Extend the existing numeric scene generation and rendering path.
- Do not replace the numeric renderer.
- Do not route ordinary numeric gauges through a new generic realism engine.
- Preserve the existing layer order unless a feature explicitly needs an additional realism layer.
- Keep decimal-point handling explicit and separate from the character image.
- Use existing `ScenePart` slot, character, position, alpha, and asset fields where practical.
- Add new scene data only when a feature cannot be represented clearly with the existing model.
- When a realism option is absent or disabled, the gauge should continue through the existing behaviour with no unrelated visual changes.

## Slice plan

| Slice | Name | Intent |
| --- | --- | --- |
| v3.7.0 | Release planning docs | Activate v3.7, lock the KISS scope, define the slice list, and prepare prompt/state documents. |
| v3.7.1 | Odometer backlash | Add bounded direction-change slack to odometer movement without altering source values or existing forward movement. |
| v3.7.2 | Per-digit response lag | Add small deterministic per-slot update delays directly to the existing numeric display path. |
| v3.7.3 | Leading-zero behaviour | Add explicit show, blank, dim, and any approved placeholder presentation to leading digit slots. |
| v3.7.4 | Segment and digit bleed | Add faint inactive segment/digit imagery using explicit image assets and existing layer composition. |
| v3.7.5 | Ghosting | Add finite previous-character fade behaviour for changed digit slots. |
| v3.7.6 | Uneven brightness | Add stable deterministic per-slot brightness variation. |
| v3.7.7 | Load sag | Add subtle brightness reduction based on known displayed segment load. |
| v3.7.8 | Tests, previews, docs checkpoint | Verify config, rendering, state settling, regression safety, preview packages, and realism-guide status. |

## Success criteria

v3.7 is complete when:

- every planned config key is parsed and validated;
- unsupported shapes and unknown fields fail clearly;
- every option produces visible but subtle behaviour in a dedicated preview;
- every stateful effect settles exactly and does not retain unbounded history;
- existing odometer and numeric gauges still render correctly when new options are absent;
- existing whole-image segmented gauge behaviour remains unchanged;
- tests cover config parsing, runtime behaviour, scene output, and disabled-option behaviour;
- the Realism Behaviour Guide accurately marks the implemented state;
- the implementation checklist is complete.

## Non-goals

- Do not redesign the numeric renderer.
- Do not redesign scene generation.
- Do not introduce a generic realism engine.
- Do not introduce a shared numeric-display runtime framework.
- Do not merge the numeric and whole-image segmented renderers.
- Do not infer segment geometry by analysing image pixels.
- Do not mutate source values, logs, exports, configured ranges, or input data.
- Do not add non-deterministic frame-to-frame variation.
- Do not implement later slices while completing an earlier slice.
- Do not let v3.7 become a realism wishlist drawer. Drawers are where simple features go to breed hinges.
