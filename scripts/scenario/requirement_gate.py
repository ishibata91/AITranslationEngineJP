from __future__ import annotations

import argparse
import json
import re
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any

VALID_STATUSES = {"explicit", "derived", "not_applicable", "deferred", "needs_human_decision"}

EXPECTED_CANDIDATE_GENERATORS = {
    "actor-goal",
    "lifecycle",
    "state-transition",
    "failure",
    "external-integration",
    "operation-audit",
}

VALID_CANDIDATE_DECISIONS = {"adopted", "merged", "rejected", "conflicted", "needs_human_decision"}

VALID_CONFLICT_STATUSES = {"resolved", "unresolved"}

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


def question_title(question: dict[str, Any]) -> str:
    for key in ("question_title", "title", "unresolved_decision"):
        value = question.get(key)
        if isinstance(value, str) and value.strip():
            return value.strip()
    return "未決判断"


def recommended_option_text(question: dict[str, Any]) -> str:
    value = question.get("recommended_option")
    if isinstance(value, int):
        return str(value)
    if isinstance(value, str) and value.strip():
        return value.strip()
    return str(question.get("recommended", "")).strip()


def recommendation_reason_text(question: dict[str, Any]) -> str:
    value = question.get("recommendation_reason")
    if isinstance(value, str) and value.strip():
        return value.strip()
    return str(question.get("recommended", "")).strip()


def sorted_questions(questions: list[dict[str, Any]]) -> list[dict[str, Any]]:
    return sorted(questions, key=lambda question: str(question.get("id") or question.get("question_id") or ""))


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
                "question_title": detail.get("question_title") or requirement.get("title", ""),
                "source_requirement": requirement.get("source_requirement", ""),
                "detail_requirement_type": detail_type,
                "unresolved_decision": detail.get("unresolved_decision", ""),
                "user_goal": detail.get("user_goal") or requirement.get("source_requirement", ""),
                "reason": detail.get("reason", ""),
                "options": detail.get("options", []),
                "recommended_option": detail.get("recommended_option", ""),
                "recommended": detail.get("recommended", ""),
                "recommendation_reason": detail.get("recommendation_reason", ""),
                "uncertainty": detail.get("uncertainty", ""),
                "after_answer_generates": detail.get("after_answer_generates", [detail_type]),
            }
            questions.append(question)

            if not non_empty(question["unresolved_decision"]):
                findings.append(Finding("error", req_id, detail_type, "question requires unresolved_decision"))
            if not non_empty(question["reason"]):
                findings.append(Finding("error", req_id, detail_type, "question requires reason"))
            if option_count(question["options"]) != 3:
                findings.append(Finding("error", req_id, detail_type, "question requires 3 options before その他"))
            if not non_empty(question["recommended"]):
                findings.append(Finding("error", req_id, detail_type, "question requires recommended"))
            if not non_empty(question["user_goal"]):
                findings.append(Finding("error", req_id, detail_type, "question requires user_goal"))
            if not non_empty(question["uncertainty"]):
                findings.append(Finding("error", req_id, detail_type, "question requires uncertainty"))
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


def normalize_generators(value: Any) -> dict[str, dict[str, Any]]:
    if isinstance(value, dict):
        normalized: dict[str, dict[str, Any]] = {}
        for name, item in value.items():
            if isinstance(name, str):
                normalized[name] = item if isinstance(item, dict) else {"status": item}
        return normalized

    if isinstance(value, list):
        normalized = {}
        for item in value:
            if not isinstance(item, dict):
                continue
            name = item.get("name") or item.get("generator")
            if isinstance(name, str) and name.strip():
                normalized[name.strip()] = item
        return normalized

    return {}


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


def candidate_artifact_exists(base_dir: Path, artifact_path: Any) -> bool:
    if not isinstance(artifact_path, str) or not artifact_path.strip():
        return False
    path = Path(artifact_path)
    if not path.is_absolute():
        path = base_dir / path
    return path.exists() and path.is_file() and non_empty(path.read_text(encoding="utf-8"))


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
        "user_goal": question.get("user_goal") or question.get("source_requirement") or source_id,
        "reason": question.get("reason", ""),
        "options": question.get("options", []),
        "recommended_option": question.get("recommended_option", ""),
        "recommended": question.get("recommended", ""),
        "recommendation_reason": question.get("recommendation_reason", ""),
        "uncertainty": question.get("uncertainty", ""),
        "after_answer_generates": question.get("after_answer_generates", [detail_type]),
    }


def validate_question_shape(findings: list[Finding], question: dict[str, Any], source_id: str, detail_type: str) -> None:
    if not non_empty(question.get("unresolved_decision")):
        findings.append(Finding("error", source_id, detail_type, "question requires unresolved_decision"))
    if not non_empty(question.get("reason")):
        findings.append(Finding("error", source_id, detail_type, "question requires reason"))
    if option_count(question.get("options")) != 3:
        findings.append(Finding("error", source_id, detail_type, "question requires 3 options before その他"))
    if not non_empty(question.get("recommended")):
        findings.append(Finding("error", source_id, detail_type, "question requires recommended"))
    if not non_empty(question.get("user_goal")):
        findings.append(Finding("error", source_id, detail_type, "question requires user_goal"))
    if not non_empty(question.get("uncertainty")):
        findings.append(Finding("error", source_id, detail_type, "question requires uncertainty"))
    if not non_empty(question.get("after_answer_generates")):
        findings.append(Finding("error", source_id, detail_type, "question requires after_answer_generates"))


def validate_candidate_coverage(data: dict[str, Any], base_dir: Path) -> tuple[list[Finding], list[dict[str, Any]]]:
    findings: list[Finding] = []
    questions: list[dict[str, Any]] = []

    generators = normalize_generators(data.get("generators"))
    if not generators:
        findings.append(Finding("error", "candidate-coverage", "generators", "generators must be a non-empty list or object"))

    for generator_name in sorted(EXPECTED_CANDIDATE_GENERATORS):
        generator = generators.get(generator_name)
        if generator is None:
            findings.append(Finding("error", "candidate-coverage", generator_name, "required generator is missing"))
            continue
        if generator.get("status") != "completed":
            findings.append(Finding("error", "candidate-coverage", generator_name, "generator status must be completed"))
        if not candidate_artifact_exists(base_dir, generator.get("artifact_path")):
            findings.append(Finding("error", "candidate-coverage", generator_name, "generator artifact is missing or empty"))

    candidates = data.get("candidates")
    if not isinstance(candidates, list) or not candidates:
        findings.append(Finding("error", "candidate-coverage", "candidates", "candidates must be a non-empty list"))
        candidates = []

    question_lookup = normalize_questions(data.get("unresolved_questions"))
    seen_questions: set[str] = set()

    for index, candidate in enumerate(candidates, start=1):
        if not isinstance(candidate, dict):
            findings.append(Finding("error", f"candidate[{index}]", "-", "candidate must be an object"))
            continue

        candidate_id = str(candidate.get("candidate_id") or candidate.get("id") or f"candidate[{index}]")
        generator_name = candidate.get("generator")
        source_requirement_id = candidate.get("source_requirement_id")
        decision = candidate.get("decision")

        if generator_name not in EXPECTED_CANDIDATE_GENERATORS:
            findings.append(Finding("error", candidate_id, "generator", f"unknown generator: {generator_name}"))
        if not non_empty(source_requirement_id):
            findings.append(Finding("error", candidate_id, "source_requirement_id", "source_requirement_id is required"))
        if decision not in VALID_CANDIDATE_DECISIONS:
            findings.append(Finding("error", candidate_id, "decision", f"invalid candidate decision: {decision}"))
            continue

        if decision in {"adopted", "merged"} and not non_empty(candidate.get("final_scenario_id")):
            findings.append(Finding("error", candidate_id, "final_scenario_id", f"{decision} requires final_scenario_id"))
        if decision == "rejected" and not non_empty(candidate.get("decision_rationale")):
            findings.append(Finding("error", candidate_id, "decision_rationale", "rejected requires decision_rationale"))
        if decision in {"conflicted", "needs_human_decision"}:
            question_id = candidate.get("question_id")
            if not isinstance(question_id, str) or not question_id.strip():
                findings.append(Finding("error", candidate_id, "question_id", f"{decision} requires question_id"))
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
                    validate_question_shape(findings, question, candidate_id, "scenario_candidate_conflict")
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
            findings.append(Finding("error", conflict_id, "status", f"invalid conflict status: {status}"))
            continue
        if status == "resolved" and not non_empty(conflict.get("resolution_rationale")):
            findings.append(Finding("error", conflict_id, "resolution_rationale", "resolved conflict requires resolution_rationale"))
        if status == "unresolved":
            question_id = conflict.get("question_id")
            if not isinstance(question_id, str) or not question_id.strip():
                findings.append(Finding("error", conflict_id, "question_id", "unresolved conflict requires question_id"))
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
                    validate_question_shape(findings, question, conflict_id, "scenario_candidate_conflict")
            findings.append(Finding("error", conflict_id, "status", "unresolved conflict requires human decision before scenario completion"))

    return findings, questions


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

    lines.extend(["", "## Questions", ""])
    if questions:
        for question in sorted_questions(questions):
            lines.append(f"- `{question.get('id', '')}` {question_title(question)}")
    else:
        lines.append("- none")

    return "\n".join(lines) + "\n"


def render_questionnaire(questions: list[dict[str, Any]]) -> str:
    lines = ["# Human Decision Questionnaire", ""]
    if not questions:
        lines.append("- none")
        return "\n".join(lines) + "\n"

    for question in sorted_questions(questions):
        question_id = question.get("id") or question.get("question_id")
        options = question.get("options", [])
        lines.extend(
            [
                f"## [{question_id}] {question_title(question)}",
                "",
                "質問:",
                str(question.get("unresolved_decision", "")),
                "",
                "やりたいこと:",
                str(question.get("user_goal") or question.get("source_requirement", "")),
                "",
                "背景:",
                str(question.get("reason", "")),
                "",
                "選択肢:",
            ]
        )
        if isinstance(options, list) and options:
            for index, option in enumerate(options, start=1):
                if isinstance(option, dict):
                    label = option.get("label", "")
                    lines.append(f"{index}. {label}")
        else:
            lines.append("1. TODO: option is missing")
        lines.append("4. その他")
        lines.extend(
            [
                "",
                "AI推奨:",
                recommended_option_text(question),
                "",
                "推奨理由:",
                recommendation_reason_text(question),
                "",
                "不確実性:",
                str(question.get("uncertainty", "")),
                "",
                "回答形式:",
                "選択肢番号を選んでください。",
                "4 の場合は、採用したい業務ルールを1〜3文で記入してください。",
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
                    "question_count": len(questions),
                    "findings": [finding.__dict__ for finding in findings],
                    "questions": sorted_questions(questions),
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
