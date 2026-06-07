#!/usr/bin/env python3
"""
extractData.v2.pas が emit した JSON の構造検証と、新旧 emit の差分検出。

対象 schema は docs/conceptual-model.md (v17) に準拠した v2 出力。
v1 schema（dialogue_groups / npcs キー）の出力もある程度検知して legacy フラグを立てる。

用途:
  - 単体検証: `python3 scripts/validate_extraction.py <file.json>`
  - coverage 表示: `--coverage`
  - regression 検知: `--baseline <other.json>` で id 集合と response 数を比較

xTranslator 書き戻し用 type field の convention:
  - 各 entry の "type" field は "<SIG> <PRIMARY_SUBRECORD>" 形式（例: "WEAP FULL", "BOOK FULL"）。
  - 同一 entry の additional 翻訳対象 field の subrecord 対応は下記 convention：
      books[*].body          → BOOK DESC
      messages[*].title      → MESG FULL
      persons[*].short_name  → NPC_ SHRT
      equipment[*].description / consumables[*].description → "<kind> DESC"
      magic[*].description / enchantments[*].description / shouts[*].description → "<SIG> DESC"
      magic_effects[*].description → MGEF DNAM
      shouts[*].words[*].dragon_text → WOOP TNAM
      perks[*].description → PERK DESC
  - 上記以外の field は xTranslator 書き戻し対象ではない（id、source、edge ref 等）。
"""
from __future__ import annotations

import argparse
import json
import sys
from collections import Counter
from pathlib import Path
from typing import Any

# v2 top-level shape
ARRAY_KEYS_V2 = (
    "dialogues",
    "quests",
    "races",
    "factions",
    "locations",
    "items",
    "equipment",
    "consumables",
    "magic",
    "enchantments",
    "magic_effects",
    "shouts",
    "books",
    "messages",
    "loading_screens",
    "perks",
)
DICT_KEYS_V2 = ("persons",)
SCALAR_KEYS_V2 = ("target_plugin",)

# 各 type 文字列ごとの必須 field（primary type のみ enforced。secondary translatable
# field は convention 列に記載のみ）。
REQUIRED_FIELDS: dict[str, set[str]] = {
    "DIAL FULL": {"id", "editor_id", "type", "source", "player_text",
                  "category", "subtype", "is_services_branch",
                  "services_type", "quest_id", "responses"},
    "INFO NAM1": {"id", "editor_id", "type", "source", "order", "text",
                  "prompt", "topic_text", "menu_display_text",
                  "speaker_kind", "speaker_id", "voice_types",
                  "previous_id", "next_ids", "conditions"},
    "QUST FULL": {"id", "editor_id", "type", "source", "name",
                  "quest_type", "stages", "objectives"},
    "QUST CNAM": {"log_index", "type", "text"},
    "QUST NNAM": {"objective_index", "type", "display_text"},
    "NPC_ FULL": {"id", "editor_id", "type", "source", "name", "short_name",
                  "sex", "voice_type_id", "class_id", "race_id", "faction_ids"},
    "RACE FULL": {"id", "editor_id", "type", "source", "name"},
    "FACT FULL": {"id", "editor_id", "type", "source", "name"},
    "MGEF FULL": {"id", "editor_id", "type", "source", "name", "description"},
    "SPEL FULL": {"id", "editor_id", "type", "source", "name", "description",
                  "magic_effect_ids"},
    "ENCH FULL": {"id", "editor_id", "type", "source", "name", "description",
                  "magic_effect_ids"},
    "SHOU FULL": {"id", "editor_id", "type", "source", "name", "description",
                  "words"},
    "WOOP FULL": {"id", "editor_id", "type", "source", "name", "dragon_text"},
    "BOOK FULL": {"id", "editor_id", "type", "source", "title", "body"},
    "MESG DESC": {"id", "editor_id", "type", "source", "title", "body", "quest_id"},
    "LSCR DESC": {"id", "editor_id", "type", "source", "body"},
    "PERK FULL": {"id", "editor_id", "type", "source", "name", "description"},
}

# 装備・消耗・所持品・場所は kind ごとに type 文字列が変わる
ITEM_TYPE_KINDS = {"KEYM", "MISC", "LIGH", "CONT", "SLGM", "DOOR", "FLOR", "FURN"}
EQUIPMENT_KINDS = {"WEAP", "ARMO", "AMMO"}
CONSUMABLE_KINDS = {"SCRL", "ALCH", "INGR"}
LOCATION_KINDS = {"CELL", "LCTN", "WRLD"}


def required_fields_for(entry: dict[str, Any]) -> set[str] | None:
    """type 文字列または kind field から必須 field 集合を返す。未知 type は None。"""
    t = entry.get("type", "")
    if t in REQUIRED_FIELDS:
        return REQUIRED_FIELDS[t]
    # kind-based: "<KIND> FULL"
    if t.endswith(" FULL"):
        sig = t.split()[0]
        base = {"id", "editor_id", "type", "source", "name"}
        if sig in EQUIPMENT_KINDS:
            return base | {"description", "kind", "enchantment_id"}
        if sig in CONSUMABLE_KINDS:
            return base | {"description", "kind", "magic_effect_ids"}
        if sig in ITEM_TYPE_KINDS:
            return base | {"kind"}
        if sig in LOCATION_KINDS:
            return base | {"kind", "parent_id"}
    return None


class Report:
    def __init__(self) -> None:
        self.errors: list[str] = []
        self.warnings: list[str] = []
        self.info: list[str] = []

    def err(self, msg: str) -> None: self.errors.append(msg)
    def warn(self, msg: str) -> None: self.warnings.append(msg)
    def note(self, msg: str) -> None: self.info.append(msg)

    def print(self) -> int:
        for line in self.info:
            print(f"[info]  {line}")
        for line in self.warnings:
            print(f"[warn]  {line}")
        for line in self.errors:
            print(f"[ERROR] {line}")
        print()
        print(f"summary: {len(self.errors)} error / "
              f"{len(self.warnings)} warning / {len(self.info)} info")
        return 1 if self.errors else 0


def detect_legacy(data: dict[str, Any]) -> bool:
    """v1 出力の判別: 'dialogue_groups' / 'npcs' キーがあれば legacy。"""
    return "dialogue_groups" in data or "npcs" in data


def validate_structure_v2(data: dict[str, Any], report: Report) -> None:
    for k in SCALAR_KEYS_V2:
        if k not in data:
            report.err(f"top-level key missing: {k}")
        elif not isinstance(data[k], str):
            report.err(f"top-level {k} must be string")
    for k in ARRAY_KEYS_V2:
        if k not in data:
            report.err(f"top-level key missing: {k}")
        elif not isinstance(data[k], list):
            report.err(f"top-level {k} must be array")
    for k in DICT_KEYS_V2:
        if k not in data:
            report.err(f"top-level key missing: {k}")
        elif not isinstance(data[k], dict):
            report.err(f"top-level {k} must be object")


def validate_entry(cat: str, idx: int, entry: dict[str, Any], report: Report) -> None:
    if "type" not in entry:
        report.err(f"{cat}[{idx}] (id={entry.get('id', '?')}) missing 'type' field — "
                   "xTranslator 書き戻しに必須")
        return
    req = required_fields_for(entry)
    if req is None:
        report.warn(f"{cat}[{idx}] (id={entry.get('id', '?')}) unknown type "
                    f"{entry['type']!r}")
        return
    missing = req - set(entry.keys())
    if missing:
        report.err(f"{cat}[{idx}] (id={entry.get('id', '?')}) "
                   f"type={entry['type']} missing fields: {sorted(missing)}")


def validate_entries_v2(data: dict[str, Any], report: Report) -> None:
    for cat in ARRAY_KEYS_V2:
        for i, e in enumerate(data.get(cat) or []):
            if not isinstance(e, dict):
                report.err(f"{cat}[{i}] is not object")
                continue
            validate_entry(cat, i, e, report)

            # nested validations
            if cat == "dialogues":
                for j, r in enumerate(e.get("responses", []) or []):
                    if not isinstance(r, dict):
                        report.err(f"dialogues[{i}].responses[{j}] not object")
                        continue
                    validate_entry(f"dialogues[{i}].responses", j, r, report)
            elif cat == "quests":
                for j, s in enumerate(e.get("stages", []) or []):
                    if not isinstance(s, dict):
                        report.err(f"quests[{i}].stages[{j}] not object")
                        continue
                    for k, log in enumerate(s.get("log_entries", []) or []):
                        validate_entry(f"quests[{i}].stages[{j}].log_entries", k, log, report)
                for j, o in enumerate(e.get("objectives", []) or []):
                    if isinstance(o, dict):
                        validate_entry(f"quests[{i}].objectives", j, o, report)
            elif cat == "shouts":
                for j, w in enumerate(e.get("words", []) or []):
                    if isinstance(w, dict):
                        validate_entry(f"shouts[{i}].words", j, w, report)

    # persons dict
    for pid, p in (data.get("persons") or {}).items():
        if not isinstance(p, dict):
            report.err(f"persons[{pid}] not object")
            continue
        if p.get("id") != pid:
            report.warn(f"persons[{pid}].id ({p.get('id')!r}) differs from key")
        validate_entry(f"persons[{pid}]", 0, p, report)


def detect_suspicious(data: dict[str, Any], report: Report) -> None:
    dg = data.get("dialogues") or []
    empty = [d.get("id", "?") for d in dg if not (d.get("responses") or [])]
    if empty:
        report.note(f"dialogues with 0 responses: {len(empty)}/{len(dg)} "
                    f"(first 5: {empty[:5]})")

    persons = data.get("persons") or {}
    person_ids = set(persons.keys())
    dangling: list[tuple[str, str]] = []
    blank_resp = 0
    total_resp = 0
    for d in dg:
        for r in d.get("responses", []) or []:
            total_resp += 1
            sid = r.get("speaker_id", "")
            if not (r.get("text") or "") and not sid:
                blank_resp += 1
            if sid and sid not in person_ids:
                dangling.append((r.get("id", "?"), sid))
    report.note(f"responses total: {total_resp}")
    if blank_resp:
        report.warn(f"responses with empty text AND no speaker: {blank_resp}")
    if dangling:
        report.warn(f"responses referring to speaker_id not in persons: "
                    f"{len(dangling)} (master plugin 未取り込み候補、"
                    f"first 5: {dangling[:5]})")

    # MagicEffect 参照の dangling 検出
    mgef_ids = {e.get("id") for e in (data.get("magic_effects") or [])
                if isinstance(e, dict)}
    ench_ids = {e.get("id") for e in (data.get("enchantments") or [])
                if isinstance(e, dict)}
    race_ids = {e.get("id") for e in (data.get("races") or [])
                if isinstance(e, dict)}
    faction_ids = {e.get("id") for e in (data.get("factions") or [])
                   if isinstance(e, dict)}

    def count_dangling(items: list[dict], id_field: str, ref_set: set[str],
                       list_field: bool = False) -> int:
        n = 0
        for it in items:
            if not isinstance(it, dict):
                continue
            if list_field:
                for ref in it.get(id_field, []) or []:
                    if ref and ref not in ref_set:
                        n += 1
            else:
                v = it.get(id_field)
                if v and v not in ref_set:
                    n += 1
        return n

    dangling_mgef_from_spel = count_dangling(
        data.get("magic") or [], "magic_effect_ids", mgef_ids, list_field=True)
    if dangling_mgef_from_spel:
        report.warn(f"magic → magic_effect_ids dangling: {dangling_mgef_from_spel}")

    dangling_ench_from_equip = count_dangling(
        data.get("equipment") or [], "enchantment_id", ench_ids)
    if dangling_ench_from_equip:
        report.warn(f"equipment → enchantment_id dangling: {dangling_ench_from_equip}")

    dangling_race = 0
    dangling_faction = 0
    for p in persons.values():
        if isinstance(p, dict):
            rid = p.get("race_id")
            if rid and rid not in race_ids:
                dangling_race += 1
            for fid in p.get("faction_ids", []) or []:
                if fid and fid not in faction_ids:
                    dangling_faction += 1
    if dangling_race:
        report.warn(f"persons → race_id dangling: {dangling_race}")
    if dangling_faction:
        report.warn(f"persons → faction_id dangling: {dangling_faction}")


def show_coverage(data: dict[str, Any]) -> None:
    print("coverage:")
    print(f"  target_plugin: {data.get('target_plugin', '?')}")
    for cat in ARRAY_KEYS_V2:
        entries = data.get(cat) or []
        kinds = Counter()
        types = Counter()
        for e in entries:
            if isinstance(e, dict):
                if "kind" in e:
                    kinds[e["kind"]] += 1
                types[e.get("type", "?")] += 1
        extras = []
        if kinds:
            extras.append(f"kinds={dict(kinds)}")
        if types and len(types) > 1:
            extras.append(f"types={dict(types)}")
        line = f"  {cat}: total={len(entries)}"
        if extras:
            line += " " + " ".join(extras)
        print(line)
    dg = data.get("dialogues") or []
    with_resp = sum(1 for d in dg if d.get("responses"))
    resp_total = sum(len(d.get("responses", []) or []) for d in dg)
    print(f"  dialogues detail: with_responses={with_resp}, "
          f"total_responses={resp_total}")
    print(f"  persons: total={len(data.get('persons') or {})}")
    print()


def diff_against_baseline(new: dict[str, Any], baseline: dict[str, Any],
                          report: Report) -> None:
    def by_id(entries: list[dict]) -> dict[str, dict]:
        return {e.get("id"): e for e in entries
                if isinstance(e, dict) and e.get("id")}

    for cat in ARRAY_KEYS_V2:
        nv = by_id(new.get(cat) or [])
        bv = by_id(baseline.get(cat) or [])
        missing = set(bv.keys()) - set(nv.keys())
        if missing:
            report.err(f"{cat}: {len(missing)} ids dropped vs baseline "
                       f"(first 5: {list(missing)[:5]})")
        added = set(nv.keys()) - set(bv.keys())
        if added:
            report.note(f"{cat}: {len(added)} new ids vs baseline")

    new_dg = by_id(new.get("dialogues") or [])
    base_dg = by_id(baseline.get("dialogues") or [])
    regression: list[tuple[str, int, int]] = []
    for did, b in base_dg.items():
        n = new_dg.get(did)
        if n is None:
            continue
        bc = len(b.get("responses", []) or [])
        nc = len(n.get("responses", []) or [])
        if nc < bc:
            regression.append((did, bc, nc))
    if regression:
        report.err(f"dialogues: {len(regression)} response counts shrunk "
                   f"(first 5: {regression[:5]})")

    base_persons = set((baseline.get("persons") or {}).keys())
    new_persons = set((new.get("persons") or {}).keys())
    dropped = base_persons - new_persons
    if dropped:
        report.err(f"persons: {len(dropped)} ids dropped vs baseline "
                   f"(first 5: {list(dropped)[:5]})")


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("file", type=Path)
    p.add_argument("--baseline", type=Path, default=None)
    p.add_argument("--coverage", action="store_true")
    args = p.parse_args()

    try:
        data = json.loads(args.file.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as e:
        print(f"[ERROR] cannot load {args.file}: {e}", file=sys.stderr)
        return 2
    if not isinstance(data, dict):
        print(f"[ERROR] top-level must be object", file=sys.stderr)
        return 2

    report = Report()
    if detect_legacy(data):
        report.warn("v1 schema 検出（dialogue_groups / npcs キー）。本 script は v2 schema 用。"
                    "v1 を validate したい場合は git log で旧版 script を取得すること")
        return report.print()

    validate_structure_v2(data, report)
    validate_entries_v2(data, report)
    detect_suspicious(data, report)

    if args.coverage:
        show_coverage(data)

    if args.baseline is not None:
        try:
            base = json.loads(args.baseline.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError) as e:
            print(f"[ERROR] cannot load baseline {args.baseline}: {e}",
                  file=sys.stderr)
            return 2
        if detect_legacy(base):
            report.warn("baseline が v1 schema。schema 不一致のため id 集合比較のみ実施")
        diff_against_baseline(data, base, report)

    return report.print()


if __name__ == "__main__":
    sys.exit(main())
