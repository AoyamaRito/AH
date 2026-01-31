# Characters Data Format (AH Board 3D)

- Location: `AH/board-3d/data/chars.json`
- Primary schema:

```
{
  "board": { "cols": 10, "rows": 10, "image": "../img/Board_AH2.jpg" },
  "characters": [
    {
      "id": "amelda",
      "name": "Amelda",
      "imgUrl": "../img/Chara/Amelda.png",
      "pos": { "x": 1, "y": 1 },
      "faction": "Foresta",
      "hp": 10,
      "atk": 3,
      "range": 1
    }
  ]
}
```

- Alternative: filename-encoded fields via `files` array. Each entry is a relative image path and can encode fields as `__key=value`.

```
{
  "files": [
    "../img/Chara/amelda__id=amelda__x=2__y=3__hp=10__atk=3.png"
  ]
}
```

- Filename encoding
  - Separator: double underscore `__`
  - Key=Value pairs: `id`, `name`, `x` (col), `y` (row), `faction`, `hp`, `atk`, `range`
  - Example: `Amelda__id=amelda__x=2__y=3__faction=Foresta__hp=10__atk=3__range=1.png`
  - Coordinates `x` and `y` are required to place the sprite. Entries without both are ignored.

- Auto-generate manifest from filenames
  - Script: `AH/board-3d/tools/scan_chars.py`
  - Usage: `python AH/board-3d/tools/scan_chars.py AH/img/Chara`
  - Output: updates `AH/board-3d/data/chars.json` with a `files` array.

- Tile coordinates: `{x, y}` map to world as center-of-tile. `x` in `[0..cols-1]`, `y` in `[0..rows-1]`.
- Image recommendations: PNG with transparent background. Place under `AH/img/Chara/`.
- If an image fails to load, a placeholder circle with the initial letter is rendered.

Board image
- By default, `../img/Board_AH2.jpg` is used. Override with `board.image` (relative path from `AH/board-3d/`).
- The plane aspect is auto-adjusted from the image width/height while preserving `board.cols`.
