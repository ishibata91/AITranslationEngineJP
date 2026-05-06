from __future__ import annotations

import argparse
import json
import re
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any

VALID_CONFLICT_STATUSES = {"resolved", "unresolved"}

@dataclass(frozen=True)
class Finding:
    severity: str
    requirement_id: str
    detail_type: str
    message: str


def default_coverage_path(markdown_path: Path) -> Path:
    return markdown_path.with_suffix(".requirement-coverage.json")


def default_candidate_coverage_path(markdown_path: Path) -> Path:
    return markdown_path.with_suffix(".candidate-coverage.json")


def read_json_coverage_file(coverage_path: Path) -> dict[str, Any]:
    try:
        parsed = json.loads(coverage_path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        raise ValueError(f"invalid requirement coverage JSON: {exc}") from exc

    if not isinstance(parsed, dict):
        raise ValueError("requirement coverage JSON must be an object")
    return parsed


def read_json_candidate_coverage(markdown_path: Path, coverage_path: Path | None = None) -> dict[str, Any]:
    sidecar_path = coverage_path or default_candidate_coverage_path(markdown_path)
    if not sidecar_path.exists():
        raise ValueError(f"missing candidate coverage JSON: {sidecar_path.as_posix()}")
    return read_json_coverage_file(sidecar_path)


def read_json_coverage(markdown_path: Path, coverage_path: Path | None = None) -> dict[str, Any]:
    sidecar_path = coverage_path or default_coverage_path(markdown_path)
    if sidecar_path.exists():
        return read_json_coverage_file(sidecar_path)

    text = markdown_path.read_text(encoding="utf-8")
    match = re.search(r"```json\s+requirement-coverage\s*\n(.*?)\n```", text, flags=re.DOTALL)
    if not match:
        raise ValueError(f"missing requirement coverage JSON: {sidecar_path.as_posix()}")

    try:
        parsed = json.loads(match.group(1))
    except json.JSONDecodeError as exc:
        raise ValueError(f"invalid requirement-coverage JSON: {exc}") from exc

    if not isinstance(parsed, dict):
        raise ValueError("requirement-coverage JSON must be an object")
    return parsed


def question_title(question: dict[str, Any]) -> str:
    for key in ("question_title", "title", "unresolved_decision"):
        value = question.get(key)
        if isinstance(value, str) and value.strip():
            return value.strip()
    return "未決判断"


def sorted_questions(questions: list[dict[str, Any]]) -> list[dict[str, Any]]:
    return sorted(questions, key=lambda question: str(question.get("id") or question.get("question_id") or ""))


def unique_questions(questions: list[dict[str, Any]]) -> list[dict[str, Any]]:
    result: list[dict[str, Any]] = []
    seen: set[str] = set()
    for question in sorted_questions(questions):
        question_id = str(question.get("id") or question.get("question_id") or "")
        if question_id in seen:
            continue
        seen.add(question_id)
        result.append(question)
    return result


def requirement_id(requirement: dict[str, Any], index: int) -> str:
    value = requirement.get("id")
    if isinstance(value, str) and value.strip():
        return value.strip()
    return f"requirement[{index}]"


def validate_requirement(requirement: dict[str, Any], index: int) -> tuple[list[Finding], list[dict[str, Any]]]:
    findings: list[Finding] = []
    questions: list[dict[str, Any]] = []
    req_id = requirement_id(requirement, index)

    details = requirement.get("detail_requirements")
    if not isinstance(details, list):
        findings.append(Finding("error", req_id, "detail_requirements", "detail_requirements must be a list"))
        return findings, questions

    for index, detail in enumerate(details, start=1):
        if not isinstance(detail, dict):
            findings.append(Finding("error", req_id, f"detail[{index}]", "detail requirement must be an object"))
            continue
        detail_type = str(detail.get("type") or f"detail[{index}]")
        status = detail.get("status")
        if status == "needs_human_decision":
            findings.append(Finding("error", req_id, detail_type, "human decision is required before scenario completion"))
            question = {
                "id": detail.get("question_id") or f"Q-{req_id}-{detail_type}",
                "question_title": detail.get("question_title") or requirement.get("title", ""),
                "source_requirement": requirement.get("source_requirement", ""),
                "detail_requirement_type": detail_type,
                "unresolved_decision": detail.get("unresolved_decision", ""),
                "premise": detail.get("premise") or requirement.get("source_requirement", ""),
                "undecided_reason": detail.get("undecided_reason") or detail.get("reason", ""),
                "options": detail.get("options", []),
                "recommended_option": detail.get("recommended_option", ""),
                "recommended": detail.get("recommended", ""),
                "recommendation_reason": detail.get("recommendation_reason", ""),
                "after_answer_generates": detail.get("after_answer_generates", [detail_type]),
            }
            questions.append(question)

    return findings, questions


def validate_coverage(data: dict[str, Any]) -> tuple[list[Finding], list[dict[str, Any]]]:
    requirements = data.get("requirements")
    if not isinstance(requirements, list) or not requirements:
        return [Finding("error", "-", "-", "requirements must be a non-empty list")], []

    all_findings: list[Finding] = []
    all_questions: list[dict[str, Any]] = []
    for index, requirement in enumerate(requirements, start=1):
        if not isinstance(requirement, dict):
            all_findings.append(Finding("error", f"requirement[{index}]", "-", "requirement must be an object"))
            continue
        findings, questions = validate_requirement(requirement, index)
        all_findings.extend(findings)
        all_questions.extend(questions)

    return all_findings, all_questions


def normalize_questions(value: Any) -> dict[str, dict[str, Any]]:
    if not isinstance(value, list):
        return {}

    normalized: dict[str, dict[str, Any]] = {}
    for item in value:
        if not isinstance(item, dict):
            continue
        question_id = item.get("question_id") or item.get("id")
        if isinstance(question_id, str) and question_id.strip():
            normalized[question_id.strip()] = {**item, "id": question_id.strip()}
    return normalized


def question_from_candidate(
    question_id: str,
    source_id: str,
    detail_type: str,
    question_lookup: dict[str, dict[str, Any]],
    fallback_title: str,
) -> dict[str, Any]:
    question = question_lookup.get(question_id, {})
    return {
        "id": question_id,
        "question_title": question.get("question_title") or question.get("title") or fallback_title,
        "source_requirement": question.get("source_requirement", source_id),
        "detail_requirement_type": detail_type,
        "unresolved_decision": question.get("unresolved_decision", ""),
        "premise": question.get("premise") or question.get("source_requirement") or source_id,
        "undecided_reason": question.get("undecided_reason") or question.get("reason", ""),
        "options": question.get("options", []),
        "recommended_option": question.get("recommended_option", ""),
        "recommended": question.get("recommended", ""),
        "recommendation_reason": question.get("recommendation_reason", ""),
        "after_answer_generates": question.get("after_answer_generates", [detail_type]),
    }


def validate_candidate_coverage(data: dict[str, Any], base_dir: Path) -> tuple[list[Finding], list[dict[str, Any]]]:
    findings: list[Finding] = []
    questions: list[dict[str, Any]] = []

    candidates = data.get("candidates")
    if not isinstance(candidates, list):
        findings.append(Finding("error", "candidate-coverage", "candidates", "candidates must be a list"))
        candidates = []

    question_lookup = normalize_questions(data.get("unresolved_questions"))
    seen_questions: set[str] = set()

    for index, candidate in enumerate(candidates, start=1):
        if not isinstance(candidate, dict):
            findings.append(Finding("error", f"candidate[{index}]", "-", "candidate must be an object"))
            continue

        candidate_id = str(candidate.get("candidate_id") or candidate.get("id") or f"candidate[{index}]")
        source_requirement_id = candidate.get("source_requirement_id")
        decision = candidate.get("decision")

        if decision in {"conflicted", "needs_human_decision"}:
            question_id = candidate.get("question_id")
            if not isinstance(question_id, str) or not question_id.strip():
                question_id = f"Q-{candidate_id}"
            else:
                question_id = question_id.strip()
            if question_id not in seen_questions:
                question = question_from_candidate(
                    question_id,
                    str(source_requirement_id or candidate_id),
                    "scenario_candidate_conflict",
                    question_lookup,
                    "scenario candidate conflict",
                )
                questions.append(question)
                seen_questions.add(question_id)
            findings.append(Finding("error", candidate_id, "decision", f"{decision} requires human decision before scenario completion"))

    conflicts = data.get("conflicts", [])
    if not isinstance(conflicts, list):
        findings.append(Finding("error", "candidate-coverage", "conflicts", "conflicts must be a list"))
        conflicts = []

    for index, conflict in enumerate(conflicts, start=1):
        if not isinstance(conflict, dict):
            findings.append(Finding("error", f"conflict[{index}]", "-", "conflict must be an object"))
            continue
        conflict_id = str(conflict.get("conflict_id") or conflict.get("id") or f"conflict[{index}]")
        status = conflict.get("status")
        if status not in VALID_CONFLICT_STATUSES:
            continue
        if status == "unresolved":
            question_id = conflict.get("question_id")
            if not isinstance(question_id, str) or not question_id.strip():
                question_id = f"Q-{conflict_id}"
            else:
                question_id = question_id.strip()
            if question_id not in seen_questions:
                question = question_from_candidate(
                    question_id,
                    conflict_id,
                    "scenario_candidate_conflict",
                    question_lookup,
                    "scenario candidate conflict",
                )
                questions.append(question)
                seen_questions.add(question_id)
            findings.append(Finding("error", conflict_id, "status", "unresolved conflict requires human decision before scenario completion"))

    return findings, questions


def render_report(path: Path, findings: list[Finding], questions: list[dict[str, Any]]) -> str:
    human_questions = unique_questions(questions)
    lines = [
        "# Requirement Gate Report",
        "",
        f"- `source`: `{path.as_posix()}`",
        f"- `status`: `{'fail' if findings else 'pass'}`",
        f"- `finding_count`: `{len(findings)}`",
        f"- `question_count`: `{len(human_questions)}`",
        "",
        "## Findings",
        "",
    ]

    if findings:
        for finding in findings:
            lines.append(f"- `{finding.severity}` `{finding.requirement_id}` `{finding.detail_type}`: {finding.message}")
    else:
        lines.append("- none")

    lines.extend(["", "## Questions", ""])
    if human_questions:
        for question in human_questions:
            lines.append(f"- `{question.get('id', '')}` {question_title(question)}")
    else:
        lines.append("- none")

    return "\n".join(lines) + "\n"


def render_questionnaire(questions: list[dict[str, Any]]) -> str:
    lines = ["# Human Decision Questionnaire", ""]
    human_questions = unique_questions(questions)
    if not human_questions:
        lines.append("- none")
        return "\n".join(lines) + "\n"

    lines.extend(
        [
            "gate は未回答 ID だけを列挙する。",
            "人間向けの質問本文は designer が仕様境界として再編集する。",
            "",
        ]
    )
    for question in human_questions:
        question_id = question.get("id") or question.get("question_id")
        lines.append(f"- `{question_id}` {question_title(question)}")

    return "\n".join(lines) + "\n"


def write_if_requested(path: str | None, content: str) -> None:
    if not path:
        return
    output_path = Path(path)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(content, encoding="utf-8")


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Validate scenario-design detail requirement coverage.")
    parser.add_argument("input", help="Path to scenario-design.md")
    parser.add_argument("--coverage", help="Path to requirement coverage JSON. Defaults to scenario-design.requirement-coverage.json")
    parser.add_argument(
        "--candidate-coverage",
        help="Path to scenario candidate coverage JSON. Defaults to scenario-design.candidate-coverage.json",
    )
    parser.add_argument("--report-out", help="Write a markdown gate report to this path")
    parser.add_argument("--questionnaire-out", help="Write a markdown questionnaire to this path")
    parser.add_argument("--json", action="store_true", help="Print machine-readable JSON result")
    return parser


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()
    input_path = Path(args.input).resolve()
    coverage_path = Path(args.coverage).resolve() if args.coverage else None
    candidate_coverage_path = Path(args.candidate_coverage).resolve() if args.candidate_coverage else None

    try:
        coverage = read_json_coverage(input_path, coverage_path)
        findings, questions = validate_coverage(coverage)
        candidate_coverage = read_json_candidate_coverage(input_path, candidate_coverage_path)
        candidate_findings, candidate_questions = validate_candidate_coverage(candidate_coverage, input_path.parent)
        findings.extend(candidate_findings)
        questions.extend(candidate_questions)
    except ValueError as exc:
        findings = [Finding("error", "-", "-", str(exc))]
        questions = []

    report = render_report(input_path, findings, questions)
    questionnaire = render_questionnaire(questions)
    write_if_requested(args.report_out, report)
    write_if_requested(args.questionnaire_out, questionnaire)

    if args.json:
        print(
            json.dumps(
                {
                    "status": "fail" if findings else "pass",
                    "finding_count": len(findings),
                    "question_count": len(unique_questions(questions)),
                    "findings": [finding.__dict__ for finding in findings],
                    "questions": unique_questions(questions),
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
