# PNG seven-segment digit config example

This example shows the intended config shape for realistic PNG seven-segment digits.

The PNG digit renderer uses full-cell digit images instead of procedural drawing or SVG segment manipulation.

## Asset layout

```text
assets/dashboard/bttf/digits/normal/
  0.png
  1.png
  2.png
  3.png
  4.png
  5.png
  6.png
  7.png
  8.png
  9.png
  dash.png
  blank.png
  dp.png
```

Rules:

- `blank.png` is the shared unlit/background cell.
- `0.png` to `9.png` and `dash.png` are complete cells using the same background as `blank.png`.
- `dp.png` is a transparent full-cell overlay with only the decimal point visible.
- All PNGs in the set must use the same canvas size.
- Do not make `dp.png` a tiny cropped dot unless the renderer is changed to handle manual alignment.

## Minimal config snippet

```yaml
dashboard:
  asset_root: assets/dashboard/bttf
  assets:
    - id: normal_png_digits
      type: png_digit_set
      glyphs:
        "0": digits/normal/0.png
        "1": digits/normal/1.png
        "2": digits/normal/2.png
        "3": digits/normal/3.png
        "4": digits/normal/4.png
        "5": digits/normal/5.png
        "6": digits/normal/6.png
        "7": digits/normal/7.png
        "8": digits/normal/8.png
        "9": digits/normal/9.png
        dash: digits/normal/dash.png
        blank: digits/normal/blank.png
        dp: digits/normal/dp.png

  decoders:
    - id: oil_temp_text
      type: format_number
      sensor: oil_temperature
      format: "000.0"

    - id: oil_temp_digits
      type: digits
      input: oil_temp_text
      asset: normal_png_digits

  blocks:
    - id: oil_temp_display
      type: seven_segment_number
      asset: normal_png_digits
      decoder: oil_temp_digits
      geometry:
        x: 1240
        y: 72
        width: 280
        height: 92
```

## Decimal point behaviour

The decimal point belongs to the digit before it.

```text
82.5  => 8.png | 2.png + dp.png overlay | 5.png
14.7  => 1.png | 4.png + dp.png overlay | 7.png
0.9   => 0.png + dp.png overlay | 9.png
-5.0  => dash.png | 5.png + dp.png overlay | 0.png
```

No extra decimal-point cell is inserted. The next digit is not shifted. This avoids asset duplication and the small but energetic goblin called manual alignment.

## Full-cell overlay requirement

`dp.png` should be the same size as every digit image. It should be transparent except for the decimal point, already positioned in the lower-right of the cell.

That keeps scaling predictable and avoids edge bleeding or per-digit offset fiddling.
