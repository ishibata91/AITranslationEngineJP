#!/usr/bin/env python3
"""Compare xEdit JSON translation fields with xTranslator XML string rows.

The script compares records at three levels:
- strict key: EDID/FormID + REC:FIELD
- strict text: REC:FIELD + XML Source/Dest text against JSON field text
- base key: EDID/FormID + REC

Dawnguard_english_japanese.xml contains English Source and Japanese Dest.
Dawnguard.esm_Export.json currently contains Japanese strings, so text matching uses
both XML Source and XML Dest.
"""

from __future__ import annotations

import argparse
import collections
import json
import re
import sys
import xml.etree.ElementTree as ET
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable

FORM_ID_RE = re.compile(r"\[....:([0-9A-Fa-f]{8})\]")
HEX_8_RE = re.compile(r"([0-9A-Fa-f]{8})")
WHITESPACE_RE = re.compile(r"\s+")


@dataclass(frozen=True)
class TranslationCandidate:
    rec_field: str
    rec_base: str
    edid: str
    formid: str
    json_field: str
    json_path: str
    text: str


@dataclass(frozen=True)
class XMLStringRow:
    rec_field: str
    rec_base: str
    edid: str
    formid: str
    source: str
    dest: str


def normalize_text(value: object) -> str:
    if value is None:
        return ""
    return str(value).strip()


def normalize_match_text(value: object) -> str:
    text = normalize_text(value)
    if not text:
        return ""
    return WHITESPACE_RE.sub(" ", text).casefold()


def normalize_edid(value: object) -> str:
    return normalize_text(value).casefold()


def normalize_formid(value: object) -> str:
    text = normalize_text(value)
    if not text:
        return ""
    match = FORM_ID_RE.search(text)
    if match:
        return match.group(1).upper()
    if text.lower().startswith("0x"):
        text = text[2:]
    matches = HEX_8_RE.findall(text)
    if matches:
        return matches[-1].upper()
    return ""


def type_to_rec_field(type_value: object) -> str:
    text = normalize_text(type_value).upper()
    if not text:
        return ""
    if ":" in text:
        return text
    parts = text.split()
    if len(parts) >= 2:
        return f"{parts[0]}:{parts[1]}"
    return text


def rec_base(rec_field: str) -> str:
    return rec_field.split(":", 1)[0]


def description_rec_field(type_value: object) -> str:
    base = rec_base(type_to_rec_field(type_value))
    if not base:
        return ""
    return f"{base}:DESC"


def message_title_rec_field(type_value: object) -> str:
    base = rec_base(type_to_rec_field(type_value))
    if not base:
        return ""
    return f"{base}:FULL"


def add_candidate(
    candidates: list[TranslationCandidate],
    *,
    record: dict[str, object],
    rec_field: str,
    json_field: str,
    source_text: object,
    json_path: str,
    fallback_edid: object = "",
    fallback_formid: object = "",
) -> None:
    text = normalize_text(source_text)
    if not rec_field or not text:
        return
    edid = normalize_edid(record.get("editor_id") or fallback_edid)
    formid = normalize_formid(record.get("id") or fallback_formid)
    candidates.append(
        TranslationCandidate(
            rec_field=rec_field,
            rec_base=rec_base(rec_field),
            edid=edid,
            formid=formid,
            json_field=json_field,
            json_path=json_path,
            text=text,
        )
    )


def iter_json_candidates(payload: dict[str, object]) -> Iterable[TranslationCandidate]:
    candidates: list[TranslationCandidate] = []

    for index, record in enumerate(payload.get("dialogue_groups", [])):
        if not isinstance(record, dict):
            continue
        add_candidate(
            candidates,
            record=record,
            rec_field=type_to_rec_field(record.get("type")),
            json_field="player_text",
            source_text=record.get("player_text"),
            json_path=f"dialogue_groups[{index}].player_text",
        )
        for response_index, response in enumerate(record.get("responses", [])):
            if not isinstance(response, dict):
                continue
            path_prefix = f"dialogue_groups[{index}].responses[{response_index}]"
            add_candidate(
                candidates,
                record=response,
                rec_field=type_to_rec_field(response.get("type")),
                json_field="text",
                source_text=response.get("text"),
                json_path=f"{path_prefix}.text",
            )
            for field_name in ("prompt", "topic_text", "menu_display_text"):
                add_candidate(
                    candidates,
                    record=response,
                    rec_field=f"INFO:{field_name.upper()}",
                    json_field=field_name,
                    source_text=response.get(field_name),
                    json_path=f"{path_prefix}.{field_name}",
                )

    for index, record in enumerate(payload.get("quests", [])):
        if not isinstance(record, dict):
            continue
        add_candidate(
            candidates,
            record=record,
            rec_field=type_to_rec_field(record.get("type")),
            json_field="name",
            source_text=record.get("name"),
            json_path=f"quests[{index}].name",
        )
        for stage_index, stage in enumerate(record.get("stages", [])):
            if not isinstance(stage, dict):
                continue
            add_candidate(
                candidates,
                record=stage,
                rec_field=type_to_rec_field(stage.get("type")),
                json_field="text",
                source_text=stage.get("text"),
                json_path=f"quests[{index}].stages[{stage_index}].text",
                fallback_edid=stage.get("parent_editor_id") or record.get("editor_id"),
                fallback_formid=stage.get("parent_id") or record.get("id"),
            )
        for objective_index, objective in enumerate(record.get("objectives", [])):
            if not isinstance(objective, dict):
                continue
            add_candidate(
                candidates,
                record=objective,
                rec_field=type_to_rec_field(objective.get("type")),
                json_field="text",
                source_text=objective.get("text"),
                json_path=f"quests[{index}].objectives[{objective_index}].text",
                fallback_edid=objective.get("parent_editor_id") or record.get("editor_id"),
                fallback_formid=objective.get("parent_id") or record.get("id"),
            )

    for category in ("items", "magic", "locations", "system", "messages", "load_screens"):
        records = payload.get(category, [])
        if not isinstance(records, list):
            continue
        for index, record in enumerate(records):
            if not isinstance(record, dict):
                continue
            path_prefix = f"{category}[{index}]"
            primary_rec_field = type_to_rec_field(record.get("type"))
            add_candidate(
                candidates,
                record=record,
                rec_field=primary_rec_field,
                json_field="name",
                source_text=record.get("name"),
                json_path=f"{path_prefix}.name",
            )
            if rec_base(primary_rec_field) == "BOOK":
                add_candidate(
                    candidates,
                    record=record,
                    rec_field="BOOK:DESC",
                    json_field="text",
                    source_text=record.get("text"),
                    json_path=f"{path_prefix}.text",
                )
            else:
                add_candidate(
                    candidates,
                    record=record,
                    rec_field=primary_rec_field,
                    json_field="text",
                    source_text=record.get("text"),
                    json_path=f"{path_prefix}.text",
                )
            if normalize_text(record.get("description")) != normalize_text(record.get("text")):
                add_candidate(
                    candidates,
                    record=record,
                    rec_field=description_rec_field(record.get("type")),
                    json_field="description",
                    source_text=record.get("description"),
                    json_path=f"{path_prefix}.description",
                )
            add_candidate(
                candidates,
                record=record,
                rec_field=message_title_rec_field(record.get("type")),
                json_field="title",
                source_text=record.get("title"),
                json_path=f"{path_prefix}.title",
            )

    npcs = payload.get("npcs", {})
    npc_records = npcs.values() if isinstance(npcs, dict) else npcs
    for index, record in enumerate(npc_records):
        if not isinstance(record, dict):
            continue
        add_candidate(
            candidates,
            record=record,
            rec_field=type_to_rec_field(record.get("type")),
            json_field="name",
            source_text=record.get("name"),
            json_path=f"npcs[{index}].name",
        )

    return candidates


def child_text(element: ET.Element, name: str) -> str:
    child = element.find(name)
    if child is None or child.text is None:
        return ""
    return child.text.strip()


def iter_xml_rows(path: Path) -> Iterable[XMLStringRow]:
    for _, element in ET.iterparse(path, events=("end",)):
        if element.tag != "String":
            continue
        rec = child_text(element, "REC").upper()
        field = child_text(element, "FIELD").upper()
        if ":" in rec:
            rec_field = rec
        elif rec and field:
            rec_field = f"{rec}:{field}"
        else:
            rec_field = rec or field
        source = child_text(element, "Source")
        dest = child_text(element, "Dest")
        if rec_field and source:
            yield XMLStringRow(
                rec_field=rec_field,
                rec_base=rec_base(rec_field),
                edid=normalize_edid(child_text(element, "EDID")),
                formid=normalize_formid(child_text(element, "FORMID")),
                source=source,
                dest=dest,
            )
        element.clear()


def strict_keys_for(candidate: TranslationCandidate) -> set[tuple[str, str, str]]:
    keys: set[tuple[str, str, str]] = set()
    if candidate.edid:
        keys.add(("edid", candidate.edid, candidate.rec_field))
    if candidate.formid:
        keys.add(("formid", candidate.formid, candidate.rec_field))
    return keys


def base_keys_for(candidate: TranslationCandidate) -> set[tuple[str, str, str]]:
    keys: set[tuple[str, str, str]] = set()
    if candidate.edid:
        keys.add(("edid", candidate.edid, candidate.rec_base))
    if candidate.formid:
        keys.add(("formid", candidate.formid, candidate.rec_base))
    return keys


def strict_keys_for_xml(row: XMLStringRow) -> set[tuple[str, str, str]]:
    keys: set[tuple[str, str, str]] = set()
    if row.edid:
        keys.add(("edid", row.edid, row.rec_field))
    if row.formid:
        keys.add(("formid", row.formid, row.rec_field))
    return keys


def base_keys_for_xml(row: XMLStringRow) -> set[tuple[str, str, str]]:
    keys: set[tuple[str, str, str]] = set()
    if row.edid:
        keys.add(("edid", row.edid, row.rec_base))
    if row.formid:
        keys.add(("formid", row.formid, row.rec_base))
    return keys


def text_keys_for(candidate: TranslationCandidate) -> set[tuple[str, str]]:
    text = normalize_match_text(candidate.text)
    if not text:
        return set()
    return {(candidate.rec_field, text)}


def base_text_keys_for(candidate: TranslationCandidate) -> set[tuple[str, str]]:
    text = normalize_match_text(candidate.text)
    if not text:
        return set()
    return {(candidate.rec_base, text)}


def text_keys_for_xml(row: XMLStringRow) -> set[tuple[str, str]]:
    keys: set[tuple[str, str]] = set()
    for text in (row.source, row.dest):
        normalized = normalize_match_text(text)
        if normalized:
            keys.add((row.rec_field, normalized))
    return keys


def base_text_keys_for_xml(row: XMLStringRow) -> set[tuple[str, str]]:
    keys: set[tuple[str, str]] = set()
    for text in (row.source, row.dest):
        normalized = normalize_match_text(text)
        if normalized:
            keys.add((row.rec_base, normalized))
    return keys


def category_for_rec_base(base: str) -> str:
    item = {"WEAP", "ARMO", "AMMO", "ALCH", "INGR", "KEYM", "MISC", "LIGH", "CONT", "SLGM", "BOOK", "FURN", "DOOR", "FLOR"}
    magic = {"SPEL", "MGEF", "ENCH", "SCRL", "SHOU"}
    location = {"LCTN", "WRLD", "CELL"}
    if base in item:
        return "item"
    if base in magic:
        return "magic"
    if base in location:
        return "location"
    if base == "PERK":
        return "system"
    if base == "MESG":
        return "message"
    if base == "LSCR":
        return "load_screen"
    if base == "QUST":
        return "quest"
    if base in {"DIAL", "INFO"}:
        return "dialogue"
    if base == "NPC_":
        return "npc"
    return "other"


def top(counter: collections.Counter[str], limit: int) -> list[dict[str, object]]:
    return [{"key": key, "count": count} for key, count in counter.most_common(limit)]


def build_report(json_path: Path, xml_path: Path, top_limit: int) -> dict[str, object]:
    payload = json.loads(json_path.read_text(encoding="utf-8"))
    json_candidates = list(iter_json_candidates(payload))
    xml_rows = list(iter_xml_rows(xml_path))

    json_strict_keys = set().union(*(strict_keys_for(candidate) for candidate in json_candidates)) if json_candidates else set()
    json_base_keys = set().union(*(base_keys_for(candidate) for candidate in json_candidates)) if json_candidates else set()
    json_text_keys = set().union(*(text_keys_for(candidate) for candidate in json_candidates)) if json_candidates else set()
    json_base_text_keys = set().union(*(base_text_keys_for(candidate) for candidate in json_candidates)) if json_candidates else set()

    strict_key_matched_rows = []
    strict_text_matched_rows = []
    strict_any_matched_rows = []
    base_key_matched_rows = []
    base_any_matched_rows = []
    base_only_by_key_rows = []
    unmatched_by_base_any_rows = []

    for row in xml_rows:
        strict_key_match = bool(strict_keys_for_xml(row) & json_strict_keys)
        strict_text_match = bool(text_keys_for_xml(row) & json_text_keys)
        strict_any_match = strict_key_match or strict_text_match
        base_key_match = bool(base_keys_for_xml(row) & json_base_keys)
        base_text_match = bool(base_text_keys_for_xml(row) & json_base_text_keys)
        base_any_match = base_key_match or base_text_match

        if strict_key_match:
            strict_key_matched_rows.append(row)
        if strict_text_match:
            strict_text_matched_rows.append(row)
        if strict_any_match:
            strict_any_matched_rows.append(row)
        if base_key_match:
            base_key_matched_rows.append(row)
        if base_any_match:
            base_any_matched_rows.append(row)
        if base_key_match and not strict_key_match:
            base_only_by_key_rows.append(row)
        if not base_any_match:
            unmatched_by_base_any_rows.append(row)

    xml_by_rec_field = collections.Counter(row.rec_field for row in xml_rows)
    json_by_rec_field = collections.Counter(candidate.rec_field for candidate in json_candidates)
    base_only_by_key_by_rec_field = collections.Counter(row.rec_field for row in base_only_by_key_rows)
    unmatched_by_base_any_by_rec_field = collections.Counter(row.rec_field for row in unmatched_by_base_any_rows)

    category_summary: dict[str, dict[str, int]] = {}
    for row in xml_rows:
        category = category_for_rec_base(row.rec_base)
        category_summary.setdefault(
            category,
            {
                "xml_rows": 0,
                "json_fields": 0,
                "strict_key_matched_xml_rows": 0,
                "strict_text_matched_xml_rows": 0,
                "strict_any_matched_xml_rows": 0,
                "base_key_only_xml_rows": 0,
                "base_any_unmatched_xml_rows": 0,
            },
        )
        category_summary[category]["xml_rows"] += 1
    for candidate in json_candidates:
        category = category_for_rec_base(candidate.rec_base)
        category_summary.setdefault(
            category,
            {
                "xml_rows": 0,
                "json_fields": 0,
                "strict_key_matched_xml_rows": 0,
                "strict_text_matched_xml_rows": 0,
                "strict_any_matched_xml_rows": 0,
                "base_key_only_xml_rows": 0,
                "base_any_unmatched_xml_rows": 0,
            },
        )
        category_summary[category]["json_fields"] += 1
    for row in strict_key_matched_rows:
        category_summary[category_for_rec_base(row.rec_base)]["strict_key_matched_xml_rows"] += 1
    for row in strict_text_matched_rows:
        category_summary[category_for_rec_base(row.rec_base)]["strict_text_matched_xml_rows"] += 1
    for row in strict_any_matched_rows:
        category_summary[category_for_rec_base(row.rec_base)]["strict_any_matched_xml_rows"] += 1
    for row in base_only_by_key_rows:
        category_summary[category_for_rec_base(row.rec_base)]["base_key_only_xml_rows"] += 1
    for row in unmatched_by_base_any_rows:
        category_summary[category_for_rec_base(row.rec_base)]["base_any_unmatched_xml_rows"] += 1

    xml_rec_fields = set(xml_by_rec_field)
    json_rec_fields = set(json_by_rec_field)

    return {
        "paths": {"json": str(json_path), "xml": str(xml_path)},
        "totals": {
            "json_translation_fields": len(json_candidates),
            "xml_string_rows": len(xml_rows),
            "xml_rows_strict_key_matched": len(strict_key_matched_rows),
            "xml_rows_strict_text_matched": len(strict_text_matched_rows),
            "xml_rows_strict_key_or_text_matched": len(strict_any_matched_rows),
            "xml_rows_base_key_matched": len(base_key_matched_rows),
            "xml_rows_base_key_or_text_matched": len(base_any_matched_rows),
            "xml_rows_base_key_only": len(base_only_by_key_rows),
            "xml_rows_unmatched_by_base_key_or_text": len(unmatched_by_base_any_rows),
        },
        "rec_field_sets": {
            "xml_only": sorted(xml_rec_fields - json_rec_fields),
            "json_only": sorted(json_rec_fields - xml_rec_fields),
            "common": sorted(xml_rec_fields & json_rec_fields),
        },
        "category_summary": dict(sorted(category_summary.items())),
        "top_xml_by_rec_field": top(xml_by_rec_field, top_limit),
        "top_json_by_rec_field": top(json_by_rec_field, top_limit),
        "top_base_key_only_xml_by_rec_field": top(base_only_by_key_by_rec_field, top_limit),
        "top_unmatched_by_base_key_or_text_xml_by_rec_field": top(unmatched_by_base_any_by_rec_field, top_limit),
    }


def print_markdown(report: dict[str, object]) -> None:
    totals = report["totals"]
    print("## Totals")
    for key, value in totals.items():
        print(f"- {key}: {value}")

    print("\n## REC Field Sets")
    sets = report["rec_field_sets"]
    for key in ("common", "xml_only", "json_only"):
        values = sets[key]
        joined = ", ".join(values) if values else "(none)"
        print(f"- {key}: {joined}")

    print("\n## Category Summary")
    print(
        "| category | xml_rows | json_fields | strict_key | strict_text | strict_any | base_key_only | base_any_unmatched |"
    )
    print("| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |")
    for category, row in report["category_summary"].items():
        print(
            f"| {category} | {row['xml_rows']} | {row['json_fields']} | "
            f"{row['strict_key_matched_xml_rows']} | {row['strict_text_matched_xml_rows']} | "
            f"{row['strict_any_matched_xml_rows']} | {row['base_key_only_xml_rows']} | "
            f"{row['base_any_unmatched_xml_rows']} |"
        )

    for section in (
        "top_xml_by_rec_field",
        "top_json_by_rec_field",
        "top_base_key_only_xml_by_rec_field",
        "top_unmatched_by_base_key_or_text_xml_by_rec_field",
    ):
        print(f"\n## {section}")
        for item in report[section]:
            print(f"- {item['key']}: {item['count']}")


def parse_args(argv: list[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--json", default="dictionaries/Dawnguard.esm_Export.json", help="Path to xEdit export JSON")
    parser.add_argument("--xml", default="dictionaries/Dawnguard_english_japanese.xml", help="Path to xTranslator XML")
    parser.add_argument("--top", type=int, default=30, help="Top N REC fields to print")
    parser.add_argument("--format", choices=("markdown", "json"), default="markdown", help="Output format")
    return parser.parse_args(argv)


def main(argv: list[str]) -> int:
    args = parse_args(argv)
    report = build_report(Path(args.json), Path(args.xml), args.top)
    if args.format == "json":
        print(json.dumps(report, ensure_ascii=False, indent=2))
    else:
        print_markdown(report)
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
