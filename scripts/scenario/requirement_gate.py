from __future__ import annotations

import argparse
import json
import re
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any

VALID_STATUSES = {"explicit", "derived", "not_applicable", "deferred", "needs_human_decision"}

DETAIL_REQUIREMENT_TYPES = {
    "success_requirement",
    "alternative_success_requirement",
    "failure_handling_requirement",
    "boundary_requirement",
    "state_requirement",
    "data_requirement",
    "consistency_requirement",
    "authorization_requirement",
    "security_requirement",
    "concurrency_requirement",
    "idempotency_requirement",
    "observability_requirement",
    "recovery_requirement",
    "performance_requirement",
    "compatibility_requirement",
    "testability_requirement",
}

KIND_REQUIRED_TYPES = {
    "operation": {
        "success_requirement",
        "failure_handling_requirement",
        "boundary_requirement",
        "state_requirement",
        "data_requirement",
        "consistency_requirement",
        "authorization_requirement",
        "testability_requirement",
    },
    "persistence": {
        "success_requirement",
        "failure_handling_requirement",
        "boundary_requirement",
        "state_requirement",
        "data_requirement",
        "consistency_requirement",
        "concurrency_requirement",
        "idempotency_requirement",
        "recovery_requirement",
        "testability_requirement",
    },
    "display": {
        "success_requirement",
        "alternative_success_requirement",
        "failure_handling_requirement",
        "boundary_requirement",
        "state_requirement",
        "authorization_requirement",
        "compatibility_requirement",
        "testability_requirement",
    },
    "external_integration": {
        "success_requirement",
        "alternative_success_requirement",
        "failure_handling_requirement",
        "boundary_requirement",
        "state_requirement",
        "security_requirement",
        "idempotency_requirement",
        "observability_requirement",
        "recovery_requirement",
        "performance_requirement",
        "testability_requirement",
    },
    "authorization": {
        "success_requirement",
        "failure_handling_requirement",
        "state_requirement",
        "authorization_requirement",
        "security_requirement",
        "observability_requirement",
        "compatibility_requirement",
        "testability_requirement",
    },
    "security": {
        "failure_handling_requirement",
        "boundary_requirement",
        "state_requirement",
        "authorization_requirement",
        "security_requirement",
        "observability_requirement",
        "recovery_requirement",
        "compatibility_requirement",
        "testability_requirement",
    },
    "workflow": {
        "success_requirement",
        "alternative_success_requirement",
        "failure_handling_requirement",
        "boundary_requirement",
        "state_requirement",
        "data_requirement",
        "consistency_requirement",
        "idempotency_requirement",
        "observability_requirement",
        "recovery_requirement",
        "testability_requirement",
    },
    "non_functional": {
        "boundary_requirement",
        "observability_requirement",
        "recovery_requirement",
        "performance_requirement",
        "compatibility_requirement",
        "testability_requirement",
    },
}

DEFAULT_KIND = "operation"


@dataclass(frozen=True)
class Finding:
    severity: str
    requirement_id: str
    detail_type: str
    message: str


def read_json_coverage(markdown_path: Path) -> dict[str, Any]:
    text = markdown_path.read_text(encoding="utf-8")
    match = re.search(r"```json\s+requirement-coverage\s*\n(.*?)\n```", text, flags=re.DOTALL)
    if not match:
        raise ValueError("missing fenced block: ```json requirement-coverage")

    try:
        parsed = json.loads(match.group(1))
    except json.JSONDecodeError as exc:
        raise ValueError(f"invalid requirement-coverage JSON: {exc}") from exc

    if not isinstance(parsed, dict):
        raise ValueError("requirement-coverage JSON must be an object")
    return parsed


def non_empty(value: Any) -> bool:
    if isinstance(value, str):
        return bool(value.strip())
    if isinstance(value, list):
        return any(non_empty(item) for item in value)
    return value is not None


def option_count(value: Any) -> int:
    if not isinstance(value, list):
        return 0
    return len([item for item in value if isinstance(item, dict) and non_empty(item.get("label"))])


def requirement_id(requirement: dict[str, Any], index: int) -> str:
    value = requirement.get("id")
    if isinstance(value, str) and value.strip():
        return value.strip()
    return f"requirement[{index}]"


def requirement_kind(requirement: dict[str, Any]) -> str:
    value = requirement.get("kind")
    if isinstance(value, str) and value.strip():
        return value.strip()
    return DEFAULT_KIND


def detail_map(requirement: dict[str, Any]) -> dict[str, dict[str, Any]]:
    details = requirement.get("detail_requirements")
    if not isinstance(details, list):
        return {}

    mapped: dict[str, dict[str, Any]] = {}
    for detail in details:
        if not isinstance(detail, dict):
            continue
        detail_type = detail.get("type")
        if isinstance(detail_type, str) and detail_type.strip():
            mapped[detail_type.strip()] = detail
    return mapped


def required_types_for(requirement: dict[str, Any]) -> set[str]:
    kind = requirement_kind(requirement)
    required = set(KIND_REQUIRED_TYPES.get(kind, KIND_REQUIRED_TYPES[DEFAULT_KIND]))
    extra = requirement.get("required_detail_types")
    if isinstance(extra, list):
        required.update(item for item in extra if isinstance(item, str) and item in DETAIL_REQUIREMENT_TYPES)
    return required


def validate_requirement(requirement: dict[str, Any], index: int) -> tuple[list[Finding], list[dict[str, Any]]]:
    findings: list[Finding] = []
    questions: list[dict[str, Any]] = []
    req_id = requirement_id(requirement, index)
    kind = requirement_kind(requirement)
    details = detail_map(requirement)

    if not non_empty(requirement.get("source_requirement")):
        findings.append(Finding("error", req_id, "-", "source_requirement is required"))

    if kind not in KIND_REQUIRED_TYPES:
        findings.append(Finding("error", req_id, "-", f"unknown requirement kind: {kind}"))

    for required_type in sorted(required_types_for(requirement)):
        if required_type not in details:
            findings.append(Finding("error", req_id, required_type, "required detail requirement is missing"))

    for detail_type, detail in sorted(details.items()):
        if detail_type not in DETAIL_REQUIREMENT_TYPES:
            findings.append(Finding("error", req_id, detail_type, "unknown detail requirement type"))
            continue

        status = detail.get("status")
        if status not in VALID_STATUSES:
            findings.append(Finding("error", req_id, detail_type, f"invalid status: {status}"))
            continue

        source_or_rationale = detail.get("source_or_rationale")
        if status in {"explicit", "derived", "not_applicable", "deferred"} and not non_empty(source_or_rationale):
            findings.append(Finding("error", req_id, detail_type, f"{status} requires source_or_rationale"))

        if status == "deferred":
            if not non_empty(detail.get("owner")):
                findings.append(Finding("error", req_id, detail_type, "deferred requires owner"))
            if not non_empty(detail.get("recheck_condition")):
                findings.append(Finding("error", req_id, detail_type, "deferred requires recheck_condition"))

        if status == "needs_human_decision":
            findings.append(Finding("error", req_id, detail_type, "human decision is required before scenario completion"))
            question = {
                "id": detail.get("question_id") or f"Q-{req_id}-{detail_type}",
                "source_requirement": requirement.get("source_requirement", ""),
                "detail_requirement_type": detail_type,
                "unresolved_decision": detail.get("unresolved_decision", ""),
                "reason": detail.get("reason", ""),
                "options": detail.get("options", []),
                "recommended": detail.get("recommended", ""),
                "after_answer_generates": detail.get("after_answer_generates", [detail_type]),
            }
            questions.append(question)

            if not non_empty(question["unresolved_decision"]):
                findings.append(Finding("error", req_id, detail_type, "question requires unresolved_decision"))
            if not non_empty(question["reason"]):
                findings.append(Finding("error", req_id, detail_type, "question requires reason"))
            if not 2 <= option_count(question["options"]) <= 4:
                findings.append(Finding("error", req_id, detail_type, "question requires 2 to 4 options"))
            if not non_empty(question["recommended"]):
                findings.append(Finding("error", req_id, detail_type, "question requires recommended"))
            if not non_empty(question["after_answer_generates"]):
                findings.append(Finding("error", req_id, detail_type, "question requires after_answer_generates"))

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


def render_report(path: Path, findings: list[Finding], questions: list[dict[str, Any]]) -> str:
    lines = [
        "# Requirement Gate Report",
        "",
        f"- `source`: `{path.as_posix()}`",
        f"- `status`: `{'fail' if findings else 'pass'}`",
        f"- `finding_count`: `{len(findings)}`",
        f"- `question_count`: `{len(questions)}`",
        "",
        "## Findings",
        "",
    ]

    if findings:
        for finding in findings:
            lines.append(f"- `{finding.severity}` `{finding.requirement_id}` `{finding.detail_type}`: {finding.message}")
    else:
        lines.append("- none")

    lines.extend(["", "## Questionnaire", ""])
    if questions:
        lines.append(render_questionnaire(questions).rstrip())
    else:
        lines.append("- none")

    return "\n".join(lines) + "\n"


def render_questionnaire(questions: list[dict[str, Any]]) -> str:
    lines = ["# Human Decision Questionnaire", ""]
    if not questions:
        lines.append("- none")
        return "\n".join(lines) + "\n"

    for question in questions:
        lines.extend(
            [
                f"## `{question['id']}`",
                "",
                f"- `source_requirement`: {question.get('source_requirement', '')}",
                f"- `detail_requirement_type`: `{question.get('detail_requirement_type', '')}`",
                f"- `unresolved_decision`: {question.get('unresolved_decision', '')}",
                f"- `reason`: {question.get('reason', '')}",
                "- `options`:",
            ]
        )
        options = question.get("options", [])
        if isinstance(options, list) and options:
            for index, option in enumerate(options, start=1):
                if isinstance(option, dict):
                    label = option.get("label", "")
                    impact = option.get("impact", "")
                    lines.append(f"  {index}. {label}: {impact}")
        else:
            lines.append("  1. TODO: option is missing")
        lines.extend(
            [
                f"- `recommended`: {question.get('recommended', '')}",
                f"- `after_answer_generates`: {', '.join(question.get('after_answer_generates', []))}",
                "",
            ]
        )

    return "\n".join(lines)


def write_if_requested(path: str | None, content: str) -> None:
    if not path:
        return
    output_path = Path(path)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(content, encoding="utf-8")


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Validate scenario-design detail requirement coverage.")
    parser.add_argument("input", help="Path to scenario-design.md")
    parser.add_argument("--report-out", help="Write a markdown gate report to this path")
    parser.add_argument("--questionnaire-out", help="Write a markdown questionnaire to this path")
    parser.add_argument("--json", action="store_true", help="Print machine-readable JSON result")
    return parser


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()
    input_path = Path(args.input).resolve()

    try:
        coverage = read_json_coverage(input_path)
        findings, questions = validate_coverage(coverage)
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
                    "question_count": len(questions),
                    "findings": [finding.__dict__ for finding in findings],
                    "questions": questions,
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
