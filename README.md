# AH Board 3D (Railway Deploy)

## Local Dev

- Run Go server:
  - `go run board-3d/server/main.go -addr :8000 -root .` (from repo root)
- Open: `http://localhost:8000/AH/board-3d/`

## Characters (filename-encoded)
- Place transparent PNGs in `img/Chara/` named like:
  - `Name__id=amelda__x=2__y=3__hp=10__atk=3.png`
- Server auto-generates `files` from this folder at `/AH/board-3d/api/chars`.

## Railway

Two options:

1) Docker (recommended)
- Repository contains `Dockerfile` at repo root.
- On Railway, create a service from GitHub repo `AoyamaRito/AH`.
- Railway builds Dockerfile and runs `/server -addr :$PORT -root /app`.
- Public URL will serve `AH/board-3d/` at `https://.../AH/board-3d/`.

2) Nixpacks (manual start)
- Set Start Command: `go run board-3d/server/main.go -addr :$PORT -root .`
- Ensure Project Root is this repository root.

## Notes
- Board image: `img/Board_AH2.jpg` is used by default.
- API JSON: `/AH/board-3d/data/chars.json` (static if present) or `/AH/board-3d/api/chars` (dynamic scan).
- Large design sources (.psd/.ai/.tiff) are ignored by `.gitignore`.

