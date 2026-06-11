#!/usr/bin/env python3
"""抽出 JSON と xTranslator XML を REC:FIELD 単位の件数で比較する。

目的:
  v18 抽出 JSON（extractData.v2.pas 出力）を xTranslator が扱う翻訳文字列単位
  （REC:FIELD、例 INFO:NAM1 / BOOK:DESC）にフラット化し、対応する xTranslator
  XML の String 行件数と突き合わせる。翻訳対象でない record type（非抽出対象）は
  XML 側から除外したうえで比較する。

突き合わせの単位:
  - JSON 側: 各カテゴリの非空 翻訳 field を 1 文字列として数える。kind 付き
    カテゴリ（items / activators / equipment / consumables / locations）は
    kind を REC に使う。
  - XML 側: String 行を REC:FIELD で数える。

master 混入の扱い:
  Ensure* は master record（Skyrim.esm 等）も emit する。xTranslator XML は
  plugin 単位なので、JSON は source == 対象 plugin の record だけを数える
  （--source で上書き可、--no-source-filter で無効化）。

非抽出対象（翻訳不要・除外）:
  NON_TARGET_REC_BASES と NON_TARGET_REC_FIELDS に列挙した REC を XML 側で除外する。
"""

from __future__ import annotations

import argparse
import collections
import json
import sys
import xml.etree.ElementTree as ET
from pathlib import Path
from typing import Any

# ===== 非抽出対象（翻訳不要・XML から除外）=====
# record type 自体を翻訳対象にしない（内部名・非表示・ゲーム設定など）。
NON_TARGET_REC_BASES = {
    "AVIF",   # actor value（内部）
    "BPTD",   # body part data（内部）
    "CLAS",   # class（話者推定の手掛かりに使うが翻訳しない）
    "CLFM",   # color form（内部）
    "EXPL",   # explosion（非表示）
    "GMST",   # game setting（内部）
    "HDPT",   # head part（内部名）
    "PROJ",   # projectile（非表示）
    "REFR",   # 配置参照の name override（保留）
    "****",   # script fragment 等、REC を持たない行
}
# REC は対象だが特定 FIELD は非対象（PERK は FULL/DESC のみ対象）。
NON_TARGET_REC_FIELDS = {
    "PERK:EPF2",  # perk entry point button label（特殊）
    "PERK:EPFD",  # perk entry point description fragment（特殊）
    "WOOP:FULL",  # 力の言葉の龍語綴り（Source=Dest 固定、翻訳禁止）
}


def is_non_target(rec_field: str) -> bool:
    base = rec_field.split(":", 1)[0]
    return base in NON_TARGET_REC_BASES or rec_field in NON_TARGET_REC_FIELDS


# ===== JSON フラット化 =====

def _nonempty(v: Any) -> bool:
    return isinstance(v, str) and v.strip() != ""


def _rec_base(rec: dict[str, Any]) -> str:
    """kind があれば kind、なければ type の先頭 token を REC base とする。"""
    kind = rec.get("kind")
    if isinstance(kind, str) and kind:
        return kind.upper()
    t = (rec.get("type") or "")
    return t.split()[0].upper() if t else ""


def flatten_json(payload: dict[str, Any], source: str | None) -> collections.Counter:
    """JSON を Counter[REC:FIELD] にする。source 指定時はその plugin の record だけ数える。"""
    c: collections.Counter = collections.Counter()

    def want(rec: dict[str, Any]) -> bool:
        if source is None:
            return True
        return (rec.get("source") or "") == source

    def add(rec_field: str, text: Any) -> None:
        if _nonempty(text):
            c[rec_field] += 1

    # dialogues: player_text → DIAL:FULL、responses[].text → INFO:NAM1、prompt → INFO:RNAM
    # prompt（RNAM）は INFO 単位で 1 つ。同一 INFO id の複数 response に複製されるため、
    # id 単位で重複排除する（xTranslator も INFO:RNAM は record 1 行）。
    seen_rnam: set = set()
    for dlg in payload.get("dialogues", []):
        if want(dlg):
            add("DIAL:FULL", dlg.get("player_text"))
        for resp in dlg.get("responses", []):
            if not want(resp):
                continue
            add("INFO:NAM1", resp.get("text"))
            rid = resp.get("id")
            if rid not in seen_rnam:
                seen_rnam.add(rid)
                add("INFO:RNAM", resp.get("prompt"))

    # quests: name → QUST:FULL、log → QUST:CNAM、objective → QUST:NNAM
    for q in payload.get("quests", []):
        if not want(q):
            continue
        add("QUST:FULL", q.get("name"))
        for st in q.get("stages", []):
            for le in st.get("log_entries", []):
                add("QUST:CNAM", le.get("text"))
        for ob in q.get("objectives", []):
            add("QUST:NNAM", ob.get("display_text"))

    # races: name → RACE:FULL、description → RACE:DESC
    for r in payload.get("races", []):
        if not want(r):
            continue
        add("RACE:FULL", r.get("name"))
        add("RACE:DESC", r.get("description"))

    # factions: name → FACT:FULL
    for r in payload.get("factions", []):
        if want(r):
            add("FACT:FULL", r.get("name"))

    # items: name → <kind>:FULL
    for r in payload.get("items", []):
        if want(r):
            add(f"{_rec_base(r)}:FULL", r.get("name"))

    # activators: name → <kind>:FULL、activate_action → <kind>:RNAM
    for r in payload.get("activators", []):
        if not want(r):
            continue
        base = _rec_base(r)
        add(f"{base}:FULL", r.get("name"))
        add(f"{base}:RNAM", r.get("activate_action"))

    # equipment / consumables: name → <kind>:FULL、description → <kind>:DESC
    for cat in ("equipment", "consumables"):
        for r in payload.get(cat, []):
            if not want(r):
                continue
            base = _rec_base(r)
            add(f"{base}:FULL", r.get("name"))
            add(f"{base}:DESC", r.get("description"))

    # magic → SPEL、enchantments → ENCH（name/description）
    for r in payload.get("magic", []):
        if want(r):
            add("SPEL:FULL", r.get("name"))
            add("SPEL:DESC", r.get("description"))
    for r in payload.get("enchantments", []):
        if want(r):
            add("ENCH:FULL", r.get("name"))
            add("ENCH:DESC", r.get("description"))

    # magic_effects: name → MGEF:FULL、description → MGEF:DNAM
    for r in payload.get("magic_effects", []):
        if want(r):
            add("MGEF:FULL", r.get("name"))
            add("MGEF:DNAM", r.get("description"))

    # shouts: name → SHOU:FULL、description → SHOU:DESC、words → WOOP:TNAM
    # WOOP は複数 shout が同じ record を embed 参照するため、id 単位で重複排除する。
    # WOOP の source は word 自身が持つ（shout が Skyrim 定義 WOOP を参照しうる）。
    seen_woop: set = set()
    for r in payload.get("shouts", []):
        if not want(r):
            continue
        add("SHOU:FULL", r.get("name"))
        add("SHOU:DESC", r.get("description"))
        for w in r.get("words", []):
            wid = w.get("id")
            if source is not None and (w.get("source") or "") != source:
                continue
            if wid in seen_woop:
                continue
            seen_woop.add(wid)
            # WOOP:FULL（龍語綴り）は翻訳禁止のため数えない。WOOP:TNAM（意味の訳）のみ対象。
            add("WOOP:TNAM", w.get("translation"))

    # books: title → BOOK:FULL、body → BOOK:DESC、author → BOOK:CNAM
    for r in payload.get("books", []):
        if not want(r):
            continue
        add("BOOK:FULL", r.get("title"))
        add("BOOK:DESC", r.get("body"))
        add("BOOK:CNAM", r.get("author"))

    # locations: name → <kind>:FULL（LCTN/WRLD/CELL）
    for r in payload.get("locations", []):
        if want(r):
            add(f"{_rec_base(r)}:FULL", r.get("name"))

    # messages: title → MESG:FULL、body → MESG:DESC、buttons → MESG:ITXT
    for r in payload.get("messages", []):
        if not want(r):
            continue
        add("MESG:FULL", r.get("title"))
        add("MESG:DESC", r.get("body"))
        for b in r.get("buttons", []):
            add("MESG:ITXT", b.get("text"))

    # loading_screens: body → LSCR:DESC
    for r in payload.get("loading_screens", []):
        if want(r):
            add("LSCR:DESC", r.get("body"))

    # perks: name → PERK:FULL、description → PERK:DESC
    for r in payload.get("perks", []):
        if want(r):
            add("PERK:FULL", r.get("name"))
            add("PERK:DESC", r.get("description"))

    # speakers: name → <NPC_|TACT>:FULL、short_name → NPC_:SHRT
    speakers = payload.get("speakers", {})
    for r in (speakers.values() if isinstance(speakers, dict) else speakers):
        if not want(r):
            continue
        base = (r.get("type") or "NPC_ FULL").split()[0].upper()
        add(f"{base}:FULL", r.get("name"))
        if base == "NPC_":
            add("NPC_:SHRT", r.get("short_name"))

    # voice_types は識別子のみで翻訳文字列を持たないため数えない。
    return c


# ===== XML 件数 =====

def count_xml(path: Path) -> collections.Counter:
    c: collections.Counter = collections.Counter()
    for _, el in ET.iterparse(str(path), events=("end",)):
        if el.tag != "String":
            continue
        rec = (el.findtext("REC") or "").strip().upper()
        fld = (el.findtext("FIELD") or "").strip().upper()
        # Source（翻訳元）が空の行は翻訳対象テキストが無いので数えない。
        # xTranslator は CNAM 等が空の record も空行として登録するため、
        # これを数えると JSON 側の「空は出さない」方針と食い違い見かけの不足になる。
        src = (el.findtext("Source") or "").strip()
        if not src:
            el.clear()
            continue
        if ":" in rec:
            rf = rec
        elif rec and fld:
            rf = f"{rec}:{fld}"
        else:
            rf = rec or fld
        if rf:
            c[rf] += 1
        el.clear()
    return c


# ===== 比較・出力 =====

def build_report(json_path: Path, xml_path: Path, source: str | None) -> dict[str, Any]:
    payload = json.loads(json_path.read_text(encoding="utf-8"))
    if source == "":
        source = payload.get("target_plugin")
    json_counts = flatten_json(payload, source)
    xml_counts = count_xml(xml_path)

    xml_in = collections.Counter({k: v for k, v in xml_counts.items() if not is_non_target(k)})
    xml_excluded = collections.Counter({k: v for k, v in xml_counts.items() if is_non_target(k)})

    all_fields = sorted(set(json_counts) | set(xml_in))
    rows = []
    for rf in all_fields:
        j = json_counts.get(rf, 0)
        x = xml_in.get(rf, 0)
        rows.append({"rec_field": rf, "json": j, "xml": x, "delta": j - x})

    return {
        "source_filter": source,
        "totals": {
            "xml_rows_total": sum(xml_counts.values()),
            "xml_rows_excluded": sum(xml_excluded.values()),
            "xml_rows_in_scope": sum(xml_in.values()),
            "json_strings": sum(json_counts.values()),
            "delta_total": sum(json_counts.values()) - sum(xml_in.values()),
        },
        "rows": rows,
        "excluded": sorted(
            ({"rec_field": k, "xml": v} for k, v in xml_excluded.items()),
            key=lambda r: (-r["xml"], r["rec_field"]),
        ),
    }


def print_markdown(rep: dict[str, Any]) -> None:
    t = rep["totals"]
    print(f"# JSON vs XML 件数比較（source filter = {rep['source_filter']!r}）\n")
    print("## 合計")
    print(f"- XML String 行 合計: {t['xml_rows_total']}")
    print(f"- XML 非抽出対象（除外）: {t['xml_rows_excluded']}")
    print(f"- XML 抽出対象（比較対象）: {t['xml_rows_in_scope']}")
    print(f"- JSON 翻訳文字列: {t['json_strings']}")
    print(f"- 差分（JSON - XML 抽出対象）: {t['delta_total']:+d}")

    match = [r for r in rep["rows"] if r["delta"] == 0 and r["xml"] > 0]
    short = [r for r in rep["rows"] if r["delta"] < 0]   # JSON 不足（取りこぼし候補）
    over = [r for r in rep["rows"] if r["delta"] > 0]    # JSON 過剰（XML に無い / 多い）

    print("\n## JSON 不足（delta<0、取りこぼし候補）")
    if short:
        print("| REC:FIELD | json | xml | delta |")
        print("| --- | ---: | ---: | ---: |")
        for r in sorted(short, key=lambda r: r["delta"]):
            print(f"| {r['rec_field']} | {r['json']} | {r['xml']} | {r['delta']:+d} |")
    else:
        print("(なし)")

    print("\n## JSON 過剰（delta>0）")
    if over:
        print("| REC:FIELD | json | xml | delta |")
        print("| --- | ---: | ---: | ---: |")
        for r in sorted(over, key=lambda r: -r["delta"]):
            print(f"| {r['rec_field']} | {r['json']} | {r['xml']} | {r['delta']:+d} |")
    else:
        print("(なし)")

    print(f"\n## 一致（delta=0、{len(match)} 種）")
    print("| REC:FIELD | json=xml |")
    print("| --- | ---: |")
    for r in sorted(match, key=lambda r: r["rec_field"]):
        print(f"| {r['rec_field']} | {r['xml']} |")

    print("\n## 非抽出対象（XML から除外した行）")
    print("| REC:FIELD | xml |")
    print("| --- | ---: |")
    for r in rep["excluded"]:
        print(f"| {r['rec_field']} | {r['xml']} |")


def parse_args(argv: list[str]) -> argparse.Namespace:
    p = argparse.ArgumentParser(description=__doc__,
                                formatter_class=argparse.RawDescriptionHelpFormatter)
    p.add_argument("--json", default="dictionaries/Dawnguard.esm_Export.json",
                   help="抽出 JSON の path")
    p.add_argument("--xml", default="dictionaries/Dawnguard_english_japanese.xml",
                   help="xTranslator XML の path")
    p.add_argument("--source", default="",
                   help="JSON 側で数える source plugin。既定は JSON の target_plugin")
    p.add_argument("--no-source-filter", action="store_true",
                   help="source フィルタを無効化し全 record を数える")
    p.add_argument("--format", choices=("markdown", "json"), default="markdown")
    return p.parse_args(argv)


def main(argv: list[str]) -> int:
    # Windows 既定の cp932 stdout で日本語が化けるため UTF-8 に統一する。
    try:
        sys.stdout.reconfigure(encoding="utf-8")
    except (AttributeError, ValueError):
        pass
    a = parse_args(argv)
    source = None if a.no_source_filter else a.source
    rep = build_report(Path(a.json), Path(a.xml), source)
    if a.format == "json":
        print(json.dumps(rep, ensure_ascii=False, indent=2))
    else:
        print_markdown(rep)
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
