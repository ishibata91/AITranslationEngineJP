from __future__ import annotations

import argparse
import json
import re
import sys
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class Finding:
    severity: str
    section: str
    message: str


def status_of(text: str) -> str:
    match = re.search(r"- `status`:\s*([^\n]+)", text)
    if not match:
        return ""
    return match.group(1).strip().strip("`")


def split_requirement_sections(text: str) -> list[tuple[str, str]]:
    matches = list(re.finditer(r"^###\s+(.+)$", text, flags=re.MULTILINE))
    sections: list[tuple[str, str]] = []
    for index, match in enumerate(matches):
        start = match.end()
        end = matches[index + 1].start() if index + 1 < len(matches) else len(text)
        title = match.group(1).strip()
        body = text[start:end]
        sections.append((title, body))
    return sections


def block_after(label: str, body: str) -> str:
    pattern = rf"^{re.escape(label)}\s*\n(.*?)(?=^\S.*:\s*$|^##\s|^###\s|\Z)"
    match = re.search(pattern, body, flags=re.MULTILINE | re.DOTALL)
    if not match:
        return ""
    return match.group(1).strip()


def has_spec_item(spec_block: str) -> bool:
    return any(line.strip().startswith("- ") and line.strip() != "- " for line in spec_block.splitlines())


def unresolved_items(unresolved_block: str) -> list[str]:
    lines = [line.strip() for line in unresolved_block.splitlines() if line.strip().startswith("- ")]
    if not lines:
        return []
    if len(lines) == 1 and lines[0] in {"- なし", "- `なし`"}:
        return []
    return lines


def validate_section(title: str, body: str) -> list[Finding]:
    findings: list[Finding] = []
    parent = block_after("親要件:", body)
    spec = block_after("仕様:", body)
    unresolved = block_after("未決:", body)
    answer = block_after("回答:", body)

    if not parent:
        findings.append(Finding("error", title, "親要件 が不足している"))
    if not spec:
        findings.append(Finding("error", title, "仕様 が不足している"))
    elif not has_spec_item(spec):
        findings.append(Finding("error", title, "仕様 に bullet がない"))
    if not unresolved:
        findings.append(Finding("error", title, "未決 が不足している"))
    if not answer:
        findings.append(Finding("error", title, "回答 が不足している"))

    unresolved = unresolved_items(unresolved)
    if unresolved and "未回答" in answer:
        findings.append(Finding("error", title, "未決に未回答が残っている"))

    return findings


def validate_file(path: Path) -> list[Finding]:
    text = path.read_text(encoding="utf-8")
    sections = split_requirement_sections(text)
    findings: list[Finding] = []

    if not sections:
        findings.append(Finding("error", "-", "親要件 section がない"))
        return findings

    for title, body in sections:
        findings.extend(validate_section(title, body))

    return findings


def render_report(path: Path, findings: list[Finding]) -> str:
    lines = [
        "# 詳細仕様差分 Gate Report",
        "",
        f"- `source`: `{path.as_posix()}`",
        f"- `status`: `{'fail' if findings else 'pass'}`",
        f"- `finding_count`: `{len(findings)}`",
        "",
        "## Findings",
        "",
    ]
    if findings:
        for finding in findings:
            lines.append(f"- `{finding.severity}` `{finding.section}`: {finding.message}")
    else:
        lines.append("- none")
    return "\n".join(lines) + "\n"


def write_if_requested(path: str | None, content: str) -> None:
    if not path:
        return
    output_path = Path(path)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(content, encoding="utf-8")


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Validate detail-spec-diff structure.")
    parser.add_argument("input", help="Path to detail-spec-diff.md")
    parser.add_argument("--report-out", help="Write a markdown gate report to this path")
    parser.add_argument("--json", action="store_true", help="Print machine-readable JSON result")
    return parser


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()
    input_path = Path(args.input).resolve()
    findings = validate_file(input_path)
    report = render_report(input_path, findings)
    write_if_requested(args.report_out, report)

    if args.json:
        print(
            json.dumps(
                {
                    "status": "fail" if findings else "pass",
                    "finding_count": len(findings),
                    "findings": [finding.__dict__ for finding in findings],
                },
                ensure_ascii=False,
                indent=2,
            )
        )
    else:
        print(report)

    return 1 if findings else 0


if __name__ == "__main__":
    sys.exit(main())
