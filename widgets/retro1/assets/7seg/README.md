# retro1 7-seg segment assets

Suggested destination inside the repository:

```text
widgets/retro1/assets/7seg/
```

Each color folder contains:

- `*_sheet.png` - original full source sheet
- `off_row.png` - top/off row
- `on_row.png` - bottom/on row
- `a_on.png` / `a_off.png`
- `b_on.png` / `b_off.png`
- `c_on.png` / `c_off.png`
- `d_on.png` / `d_off.png`
- `e_on.png` / `e_off.png`
- `f_on.png` / `f_off.png`
- `g_on.png` / `g_off.png`

These are extracted from the supplied 8888 source sheets and are intended for
use as composited per-segment sprites in `retro1_7seg3`.

Example embed directive:

```go
package retro1

import "embed"

//go:embed assets/7seg/*/*.png
var sevenSegAssets embed.FS
```
