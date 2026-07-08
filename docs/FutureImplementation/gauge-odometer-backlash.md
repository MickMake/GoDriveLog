# Candidate: Odometer Backlash

Origin: `docs/v3.7/PlannedFeatures.md`

Status: **not implemented on `main`**.

Earlier v3.5 docs/checklists claimed `backlash` was implemented. A code audit found no `realism.backlash` config field, no allowed YAML key, and no odometer runtime behaviour for direction-change slack. Treat those older claims as stale documentation, not implementation truth.

Treat odometer `backlash` as:

```text
candidate requiring implementation before use
```

## Candidate behaviour

`backlash` would model direction-change slack for odometer wheels.

Existing odometer realism can create general mechanical feel, but does not fully create direction-change backlash:

| Existing option | Why it is not backlash |
| --- | --- |
| `drum_slop` | Static wheel alignment imperfection; does not care about direction changes. |
| `carry_drag` | Rollover coupling between wheels; not reverse-direction slack. |
| `snap_settle` | Landing/settle effect; not slack when reversing. |
| `movement: linear`, `ease_out`, `bell` | Movement curves; not mechanical play. |
| `wraparound` | Route choice across digit boundaries; not slack. |

If implemented, `backlash` should be its own odometer-only feature.
