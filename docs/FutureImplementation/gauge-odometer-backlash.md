# Candidate: Odometer Backlash

Origin: `docs/v3.7/PlannedFeatures.md`

`backlash` appears in earlier odometer realism planning and may need a focused implementation/audit slice.

Before implementation, verify the current code state. Do not rely only on old checklists or prompt files.

Treat odometer `backlash` as:

```text
candidate requiring audit before implementation
```

## Candidate behaviour

`backlash` would model direction-change slack for odometer wheels.

Existing odometer realism can create general mechanical feel, but may not fully create direction-change backlash:

| Existing option | Why it is not backlash |
| --- | --- |
| `drum_slop` | Static wheel alignment imperfection; does not care about direction changes. |
| `carry_drag` | Rollover coupling between wheels; not reverse-direction slack. |
| `snap_settle` | Landing/settle effect; not slack when reversing. |
| `movement: linear`, `ease_out`, `bell` | Movement curves; not mechanical play. |
| `wraparound` | Route choice across digit boundaries; not slack. |

If implemented, `backlash` should be its own odometer-only feature.
