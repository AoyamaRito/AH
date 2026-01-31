#!/usr/bin/env python3
import os, json, sys

"""
Scan a directory (e.g., AH/img/Chara) for character image files and
generate AH/board-3d/data/chars.json with a `files` array.

Filenames can encode fields via `__key=value` segments, e.g.:
  Amelda__id=amelda__x=2__y=3__hp=10__atk=3.png

Usage:
  python AH/board-3d/tools/scan_chars.py AH/img/Chara [out_json]
  (default out_json: AH/board-3d/data/chars.json)
"""

IMG_EXTS = {'.png', '.jpg', '.jpeg', '.webp'}

def main():
    if len(sys.argv) < 2:
        print("Usage: scan_chars.py <image_dir> [out_json]", file=sys.stderr)
        sys.exit(1)
    img_dir = sys.argv[1]
    out_path = sys.argv[2] if len(sys.argv) > 2 else os.path.join(os.path.dirname(__file__), '..', 'data', 'chars.json')

    if not os.path.isdir(img_dir):
        print(f"Not a directory: {img_dir}", file=sys.stderr)
        sys.exit(2)

    files = []
    for root, _, fnames in os.walk(img_dir):
        for fn in fnames:
            ext = os.path.splitext(fn)[1].lower()
            if ext in IMG_EXTS:
                rel = os.path.relpath(os.path.join(root, fn), os.path.join(os.path.dirname(__file__), '..'))
                # normalize to forward slashes for browser
                rel = rel.replace('\\', '/')
                # prefix with .. (from AH/board-3d to AH/img/...)
                if not rel.startswith('..'):
                    rel = '../' + rel
                files.append(rel)

    data = { 'files': sorted(files) }
    os.makedirs(os.path.dirname(out_path), exist_ok=True)
    with open(out_path, 'w', encoding='utf-8') as wf:
        json.dump(data, wf, ensure_ascii=False, indent=2)
    print(f"Wrote {len(files)} files to {out_path}")

if __name__ == '__main__':
    main()

