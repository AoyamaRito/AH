#!/usr/bin/env python3
import csv, json, sys, os

"""
Convert a characters CSV into chars.json for AH Board 3D

Usage:
  python AH/board-3d/tools/csv2json.py <csv_path> [out_json]

CSV headers (flexible):
  id, name, img | filename, x | col, y | row, faction, hp, atk, range

Output JSON schema written to AH/board-3d/data/chars.json by default.
"""

def to_int(v, default=None):
    try:
        if v is None or v == "":
            return default
        return int(float(v))
    except Exception:
        return default

def row_to_char(row):
    img = row.get('img') or row.get('filename') or row.get('file')
    if not img:
        return None
    x = to_int(row.get('x') or row.get('col'))
    y = to_int(row.get('y') or row.get('row'))
    if x is None or y is None:
        # position required for placement
        return None
    ch = {
        'id': (row.get('id') or os.path.splitext(os.path.basename(img))[0]).lower(),
        'name': row.get('name') or row.get('Name') or row.get('id') or os.path.splitext(os.path.basename(img))[0],
        'imgUrl': img,
        'pos': { 'x': x, 'y': y }
    }
    # optional stats
    for k in ['faction','hp','atk','range']:
        v = row.get(k) or row.get(k.capitalize())
        if v is not None and v != '':
            if k in ['hp','atk','range']:
                iv = to_int(v)
                if iv is not None:
                    ch[k] = iv
            else:
                ch[k] = v
    return ch

def main():
    if len(sys.argv) < 2:
        print("Usage: csv2json.py <csv_path> [out_json]", file=sys.stderr)
        sys.exit(1)
    csv_path = sys.argv[1]
    out_path = sys.argv[2] if len(sys.argv) > 2 else os.path.join(os.path.dirname(__file__), '..', 'data', 'chars.json')
    chars = []
    with open(csv_path, newline='', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            ch = row_to_char(row)
            if ch:
                chars.append(ch)
    data = { 'board': { 'cols': 10, 'rows': 10 }, 'characters': chars }
    os.makedirs(os.path.dirname(out_path), exist_ok=True)
    with open(out_path, 'w', encoding='utf-8') as wf:
        json.dump(data, wf, ensure_ascii=False, indent=2)
    print(f"Wrote {len(chars)} characters to {out_path}")

if __name__ == '__main__':
    main()

