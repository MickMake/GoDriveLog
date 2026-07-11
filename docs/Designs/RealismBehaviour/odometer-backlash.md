# Candidate: Odometer Backlash

Origin: `docs/v3.7/PlannedFeatures.md`

Status: **not implemented on `main`**.

## Canonical behaviour definition

The behaviour definition lives in [`../RealismBehaviourGuide/backlash.md`](../RealismBehaviourGuide/backlash.md).

Do not redefine backlash behaviour here. Use this file only as backlog/planning context for a future odometer backlash implementation ticket.

## Implementation truth

Earlier v3.5 docs/checklists claimed `backlash` was implemented. A code audit found no `realism.backlash` config field, no allowed YAML key, and no odometer runtime behaviour for direction-change slack.

Treat those older claims as stale documentation, not implementation truth.

## Implementation planning notes

If implemented, `backlash` should be its own odometer-only feature.

Existing odometer realism can create general mechanical feel, but does not fully create direction-change backlash:

| Existing option | Why it is not backlash |
| --- | --- |
| `drum_slop` | Static wheel alignment imperfection; does not care about direction changes. |
| `carry_drag` | Rollover coupling between wheels; not reverse-direction slack. |
| `snap_settle` | Landing/settle effect; not slack when reversing. |
| `movement: linear`, `ease_out`, `bell` | Movement curves; not mechanical play. |
| `wraparound` | Route choice across digit boundaries; not slack. |

## Suggested future implementation ticket

- Add odometer-only `realism.backlash` config validation and runtime behaviour.
- Test that direction changes briefly take up slack without changing final numeric correctness.
- Keep `backlash` distinct from `drum_slop`, `carry_drag`, `snap_settle`, and movement curves.
