# Gauge Realism Map

Origin: `docs/v3.7/PlannedFeatures.md`

This map is a planning aid only. Do not treat it as implementation truth without checking the current code and completed release docs.

| Realism option | Numeric | Radial | Odometer | Indicator | Bar | Segmented |
| --- | --- | --- | --- | --- | --- | --- |
| `movement` | 🟡 parse only | 🟡 legacy `movement_policy` | ✅ `odometer.movement` (`instant`, `linear`, `ease_out`, `bell`) | ❌ | 🟡 via finite movement/damping policy only | ❌ |
| `wraparound` | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| `drum_slop` | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| `carry_drag` | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| `snap_settle` | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ |
| `backlash` | ❌ | ❌ | ❌ not implemented; stale v3.5 docs previously claimed it | ❌ | ❌ | ❌ |
| `hysteresis` | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ |
| `stiction` | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ |
| `damping` | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ |
| `overshoot` | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ |
| `peg_bounce` | ❌ | ✅ | ❌ | ❌ | ✅ | ❌ |
| `thermal_fade` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ✅ | ❌ | 🍺 potential candidate / needs beer thought |
| `per_digit_response_lag` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `leading_zero_behaviour` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `decimal_point_behaviour` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `needle_shadow` | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| `calibration_offset` | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| `segment_bleed` / `digit_bleed` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `ghosting` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `uneven_brightness` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `load_sag` | 🍺 potential candidate / needs beer thought | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought |
| `stepped_fill` | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought | 🍺 potential candidate / needs beer thought |
| `quantized_fill` | ❌ | ❌ | ❌ | ❌ | 🍺 potential candidate / needs beer thought | 🍺 potential candidate / needs beer thought |
