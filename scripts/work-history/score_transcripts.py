from __future__ import annotations

import argparse
import json
import re
from collections import Counter
from dataclasses import dataclass
from datetime import UTC, datetime
from pathlib import Path
from typing import Any

DEFAULT_OUTPUT_ROOT = Path("work_history/runs")
ANALYSIS_DIR_NAME = "analysis"
SCORE_JSON_NAME = "benchmark-score.json"
TRANSCRIPT_REFS_NAME = "transcript_refs.json"
RUN_TITLE_NAME = "run-title.txt"
MAX_FOLDER_NAME_CHARS = 90
MAX_TEXT_CHARS = 220
DEFAULT_LONG_IDLE_MS = 10 * 60 * 1000

USER_CORRECTION_KEYWORDS = (
    "違う",
    "ちがう",
    "だめ",
    "ダメ",
    "取れてなさそう",
    "なくない",
    "停止",
    "failed",
    "エラー",
    "やっぱ",
    "いやー",
)
SUBAGENT_KEYWORDS = ("subagent", "runSubagent", "spawn_agent", "SubagentStop", "サブエージェント")


@dataclass(frozen=True)
class TranscriptRef:
    runtime: str
    transcript_path: Path


def parse_timestamp(value: Any) -> datetime | None:
    if isinstance(value, int | float):
        seconds = float(value)
        if seconds > 10_000_000_000:
            seconds = seconds / 1000
        return datetime.fromtimestamp(seconds, UTC)
    if not isinstance(value, str) or not value.strip():
        return None
    try:
        parsed = datetime.fromisoformat(value.replace("Z", "+00:00"))
    except ValueError:
        return None
    if parsed.tzinfo is None:
        return parsed.replace(tzinfo=UTC)
    return parsed.astimezone(UTC)


def format_timestamp(value: datetime | None) -> str | None:
    if value is None:
        return None
    return value.astimezone(UTC).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def normalize_text(value: str | None) -> str:
    if not value:
        return ""
    return " ".join(value.split())


def truncate_text(value: str, limit: int = MAX_TEXT_CHARS) -> str:
    text = normalize_text(value)
    if len(text) <= limit:
        return text
    return text[: limit - 1].rstrip() + "…"


def coerce_text(value: Any) -> str:
    if isinstance(value, str):
        return value
    if isinstance(value, list):
        return "\n".join(part for part in (coerce_text(item) for item in value) if part)
    if isinstance(value, dict):
        for key in ("text", "content", "message", "output"):
            if key in value:
                text = coerce_text(value[key])
                if text:
                    return text
    return ""


def is_human_prompt(text: str) -> bool:
    normalized = normalize_text(text)
    if not normalized:
        return False
    blocked_prefixes = (
        "# AGENTS.md instructions",
        "The following is the Codex agent history",
        "<skill>",
        "<subagent_notification>",
        "<turn_aborted>",
        "<permissions instructions>",
        "## Memory",
        "<collaboration_mode>",
        "<apps_instructions>",
        "<skills_instructions>",
        "<plugins_instructions>",
    )
    return not any(normalized.startswith(prefix) for prefix in blocked_prefixes)


def parse_jsonl(path: Path) -> tuple[list[tuple[int, dict[str, Any]]], list[str]]:
    events: list[tuple[int, dict[str, Any]]] = []
    gaps: list[str] = []
    try:
        lines = path.read_text(encoding="utf-8").splitlines()
    except OSError as error:
        return [], [f"{path}: cannot read transcript: {error}"]
    for line_number, line in enumerate(lines, start=1):
        if not line.strip():
            continue
        try:
            parsed = json.loads(line)
        except json.JSONDecodeError as error:
            gaps.append(f"{path}:{line_number}: invalid json: {error.msg}")
            continue
        if isinstance(parsed, dict):
            events.append((line_number, parsed))
        else:
            gaps.append(f"{path}:{line_number}: non-object jsonl event")
    return events, gaps


def source_ref(path: Path, line_number: int) -> str:
    return f"{path}:{line_number}"


def evidence(path: Path, line_number: int, timestamp: datetime | None, text: str) -> dict[str, Any]:
    return {
        "timestamp": format_timestamp(timestamp),
        "source_ref": source_ref(path, line_number),
        "excerpt": truncate_text(text),
    }


def bounded_score(value: float) -> int:
    return max(0, min(100, round(value)))


def command_key(text: str) -> str:
    normalized = normalize_text(text)
    normalized = re.sub(r"\s+", " ", normalized)
    return normalized[:300]


def is_nonzero_tool_result(text: str) -> bool:
    match = re.search(r"Process exited with code\s+(-?\d+)", text)
    if match:
        return int(match.group(1)) != 0
    return False


def message_content(payload: dict[str, Any]) -> str:
    content = payload.get("content")
    if isinstance(content, list):
        return normalize_text("\n".join(coerce_text(item) for item in content if coerce_text(item)))
    return normalize_text(coerce_text(content or payload.get("message") or payload.get("text")))


def new_metrics() -> dict[str, int]:
    return {
        "duration_ms_total": 0,
        "active_duration_ms_total": 0,
        "user_turns": 0,
        "assistant_turns": 0,
        "tool_calls": 0,
        "subagent_calls": 0,
        "nonzero_tool_results": 0,
        "user_corrections": 0,
        "long_idle_gaps": 0,
        "repeated_tool_commands": 0,
    }


def merge_metrics(target: dict[str, int], source: dict[str, int]) -> None:
    for key, value in source.items():
        target[key] = target.get(key, 0) + value


def empty_evidence_refs() -> dict[str, list[dict[str, Any]]]:
    return {
        "long_idle_gaps": [],
        "nonzero_tool_results": [],
        "user_corrections": [],
        "repeated_tool_commands": [],
    }


def merge_evidence_refs(target: dict[str, list[dict[str, Any]]], source: dict[str, list[dict[str, Any]]]) -> None:
    for key, values in source.items():
        target.setdefault(key, []).extend(values)


def add_timestamp(
    timestamps: list[tuple[datetime, str]],
    timestamp: datetime | None,
    source: str,
    started_at: datetime | None,
    ended_at: datetime | None,
) -> tuple[datetime | None, datetime | None]:
    if timestamp is None:
        return started_at, ended_at
    timestamps.append((timestamp, source))
    started_at = min(started_at, timestamp) if started_at else timestamp
    ended_at = max(ended_at, timestamp) if ended_at else timestamp
    return started_at, ended_at


def collect_idle_gaps(
    timestamps: list[tuple[datetime, str]],
    threshold_ms: int,
) -> tuple[int, list[dict[str, Any]]]:
    refs: list[dict[str, Any]] = []
    previous: tuple[datetime, str] | None = None
    for current in sorted(timestamps, key=lambda item: item[0]):
        if previous is not None:
            gap_ms = round((current[0] - previous[0]).total_seconds() * 1000)
            if gap_ms >= threshold_ms:
                refs.append(
                    {
                        "from_timestamp": format_timestamp(previous[0]),
                        "to_timestamp": format_timestamp(current[0]),
                        "gap_ms": gap_ms,
                        "from_source_ref": previous[1],
                        "to_source_ref": current[1],
                    }
                )
        previous = current
    return len(refs), refs


def active_duration_ms(duration_ms: int, idle_gap_refs: list[dict[str, Any]]) -> int:
    idle_ms = sum(int(item.get("gap_ms") or 0) for item in idle_gap_refs)
    return max(0, duration_ms - idle_ms)


def build_session(runtime: str, session_id: str | None, path: Path, started_at: datetime | None, ended_at: datetime | None) -> dict[str, Any]:
    duration_ms = 0
    if started_at is not None and ended_at is not None:
        duration_ms = max(0, round((ended_at - started_at).total_seconds() * 1000))
    return {
        "runtime": runtime,
        "session_id": session_id,
        "transcript_path": str(path),
        "started_at": format_timestamp(started_at),
        "ended_at": format_timestamp(ended_at),
        "duration_ms": duration_ms,
    }


def extract_codex(ref: TranscriptRef, long_idle_ms: int) -> tuple[dict[str, Any], dict[str, int], dict[str, list[dict[str, Any]]], list[str], str]:
    path = ref.transcript_path
    events, gaps = parse_jsonl(path)
    has_response_items = any(event.get("type") == "response_item" for _, event in events)
    session_id: str | None = None
    first_user_prompt = ""
    started_at: datetime | None = None
    ended_at: datetime | None = None
    timestamps: list[tuple[datetime, str]] = []
    metrics = new_metrics()
    evidence_refs = empty_evidence_refs()
    commands: Counter[str] = Counter()
    seen_messages: set[tuple[str, str | None, str]] = set()

    for line_number, event in events:
        timestamp = parse_timestamp(event.get("timestamp"))
        started_at, ended_at = add_timestamp(timestamps, timestamp, source_ref(path, line_number), started_at, ended_at)
        event_type = event.get("type")
        payload = event.get("payload") if isinstance(event.get("payload"), dict) else {}

        if event_type == "session_meta":
            if isinstance(payload.get("id"), str):
                session_id = payload["id"]
            meta_started = parse_timestamp(payload.get("timestamp"))
            started_at, ended_at = add_timestamp(timestamps, meta_started, source_ref(path, line_number), started_at, ended_at)
            continue

        if event_type == "event_msg":
            if has_response_items:
                continue
            payload_type = str(payload.get("type") or "")
            if payload_type == "user_message":
                text = normalize_text(coerce_text(payload.get("message")))
                key = ("user", format_timestamp(timestamp), text)
                if text and is_human_prompt(text) and key not in seen_messages:
                    seen_messages.add(key)
                    metrics["user_turns"] += 1
                    if not first_user_prompt:
                        first_user_prompt = text
                    if any(keyword in text for keyword in USER_CORRECTION_KEYWORDS):
                        metrics["user_corrections"] += 1
                        evidence_refs["user_corrections"].append(evidence(path, line_number, timestamp, text))
            elif payload_type == "agent_message":
                text = normalize_text(coerce_text(payload.get("message")))
                key = ("assistant", format_timestamp(timestamp), text)
                if text and key not in seen_messages:
                    seen_messages.add(key)
                    metrics["assistant_turns"] += 1
            elif "exec_command" in payload_type:
                metrics["tool_calls"] += 1
                command = normalize_text(coerce_text(payload.get("command") or payload.get("message") or payload.get("aggregated_output")))
                if command:
                    commands[command_key(command)] += 1
                if any(keyword.lower() in command.lower() for keyword in SUBAGENT_KEYWORDS):
                    metrics["subagent_calls"] += 1
                exit_code = payload.get("exit_code")
                text = normalize_text(coerce_text(payload.get("aggregated_output") or payload.get("message") or payload.get("command")))
                if exit_code not in {None, 0, "0"} or is_nonzero_tool_result(text):
                    metrics["nonzero_tool_results"] += 1
                    evidence_refs["nonzero_tool_results"].append(evidence(path, line_number, timestamp, text or command))

        if event_type == "response_item":
            payload_type = payload.get("type")
            role = payload.get("role")
            if payload_type == "message" and role == "user":
                text = message_content(payload)
                key = ("user", format_timestamp(timestamp), text)
                if text and is_human_prompt(text) and key not in seen_messages:
                    seen_messages.add(key)
                    metrics["user_turns"] += 1
                    if not first_user_prompt:
                        first_user_prompt = text
                    if any(keyword in text for keyword in USER_CORRECTION_KEYWORDS):
                        metrics["user_corrections"] += 1
                        evidence_refs["user_corrections"].append(evidence(path, line_number, timestamp, text))
            elif payload_type == "message" and role == "assistant":
                text = message_content(payload)
                key = ("assistant", format_timestamp(timestamp), text)
                if text and key not in seen_messages:
                    seen_messages.add(key)
                    metrics["assistant_turns"] += 1
            elif payload_type == "function_call":
                metrics["tool_calls"] += 1
                name = str(payload.get("name") or "tool")
                arguments = coerce_text(payload.get("arguments"))
                command = f"{name} {arguments}".strip()
                commands[command_key(command)] += 1
                if any(keyword.lower() in command.lower() for keyword in SUBAGENT_KEYWORDS):
                    metrics["subagent_calls"] += 1
            elif payload_type == "function_call_output":
                text = coerce_text(payload.get("output"))
                if is_nonzero_tool_result(text):
                    metrics["nonzero_tool_results"] += 1
                    evidence_refs["nonzero_tool_results"].append(evidence(path, line_number, timestamp, text))

    metrics["long_idle_gaps"], evidence_refs["long_idle_gaps"] = collect_idle_gaps(timestamps, long_idle_ms)
    repeated = [(command, count) for command, count in commands.items() if count > 1]
    metrics["repeated_tool_commands"] = sum(count - 1 for _, count in repeated)
    for command, count in sorted(repeated, key=lambda item: item[1], reverse=True)[:20]:
        evidence_refs["repeated_tool_commands"].append({"count": count, "command_excerpt": truncate_text(command)})

    session = build_session(ref.runtime, session_id, path, started_at, ended_at)
    metrics["duration_ms_total"] = session["duration_ms"]
    metrics["active_duration_ms_total"] = active_duration_ms(session["duration_ms"], evidence_refs["long_idle_gaps"])
    return session, metrics, evidence_refs, gaps, first_user_prompt


def extract_copilot(ref: TranscriptRef, long_idle_ms: int) -> tuple[dict[str, Any], dict[str, int], dict[str, list[dict[str, Any]]], list[str], str]:
    path = ref.transcript_path
    events, gaps = parse_jsonl(path)
    session_id: str | None = None
    first_user_prompt = ""
    started_at: datetime | None = None
    ended_at: datetime | None = None
    timestamps: list[tuple[datetime, str]] = []
    metrics = new_metrics()
    evidence_refs = empty_evidence_refs()
    commands: Counter[str] = Counter()

    for line_number, event in events:
        request_events = event.get("v") if event.get("k") == ["requests"] and isinstance(event.get("v"), list) else None
        if request_events is not None:
            for request in request_events:
                if not isinstance(request, dict):
                    continue
                timestamp = parse_timestamp(request.get("timestamp"))
                started_at, ended_at = add_timestamp(timestamps, timestamp, source_ref(path, line_number), started_at, ended_at)
                result = request.get("result") if isinstance(request.get("result"), dict) else {}
                metadata = result.get("metadata") if isinstance(result.get("metadata"), dict) else {}
                if isinstance(metadata.get("sessionId"), str):
                    session_id = metadata["sessionId"]
                message = request.get("message") if isinstance(request.get("message"), dict) else {}
                text = normalize_text(coerce_text(message.get("text") or message.get("content")))
                if text and is_human_prompt(text):
                    metrics["user_turns"] += 1
                    if not first_user_prompt:
                        first_user_prompt = text
                    if any(keyword in text for keyword in USER_CORRECTION_KEYWORDS):
                        metrics["user_corrections"] += 1
                        evidence_refs["user_corrections"].append(evidence(path, line_number, timestamp, text))

                assistant_text_parts: list[str] = []
                for item in request.get("response") or []:
                    if isinstance(item, dict):
                        assistant_text_parts.append(coerce_text(item.get("value") or item.get("content") or item.get("text")))
                assistant_text = normalize_text("\n".join(part for part in assistant_text_parts if part))
                if assistant_text:
                    metrics["assistant_turns"] += 1

                for round_data in request.get("toolCallRounds") or []:
                    if not isinstance(round_data, dict):
                        continue
                    for tool_call in round_data.get("toolCalls") or []:
                        if not isinstance(tool_call, dict):
                            continue
                        metrics["tool_calls"] += 1
                        command = normalize_text(coerce_text(tool_call.get("name") or tool_call.get("arguments")))
                        if command:
                            commands[command_key(command)] += 1
                        if any(keyword.lower() in command.lower() for keyword in SUBAGENT_KEYWORDS):
                            metrics["subagent_calls"] += 1
                    for tool_result in (round_data.get("toolCallResults") or {}).values():
                        result_text = normalize_text(coerce_text(tool_result))
                        if is_nonzero_tool_result(result_text):
                            metrics["nonzero_tool_results"] += 1
                            evidence_refs["nonzero_tool_results"].append(evidence(path, line_number, timestamp, result_text))

                if isinstance(result.get("errorDetails"), dict):
                    error_text = normalize_text(coerce_text(result["errorDetails"]))
                    if error_text:
                        metrics["nonzero_tool_results"] += 1
                        evidence_refs["nonzero_tool_results"].append(evidence(path, line_number, timestamp, error_text))
            continue

        timestamp = parse_timestamp(event.get("timestamp") or event.get("createdAt") or event.get("time"))
        started_at, ended_at = add_timestamp(timestamps, timestamp, source_ref(path, line_number), started_at, ended_at)
        event_type = str(event.get("type") or event.get("event") or "")

        if event_type == "session.start":
            for key in ("sessionId", "session_id", "id"):
                if isinstance(event.get(key), str):
                    session_id = event[key]
                    break
            continue

        if event_type == "user.message":
            text = normalize_text(coerce_text(event.get("message") or event.get("text") or event.get("content")))
            if text and is_human_prompt(text):
                metrics["user_turns"] += 1
                if not first_user_prompt:
                    first_user_prompt = text
                if any(keyword in text for keyword in USER_CORRECTION_KEYWORDS):
                    metrics["user_corrections"] += 1
                    evidence_refs["user_corrections"].append(evidence(path, line_number, timestamp, text))
        elif event_type == "assistant.message":
            text = normalize_text(coerce_text(event.get("message") or event.get("text") or event.get("content")))
            if text:
                metrics["assistant_turns"] += 1
        elif event_type.startswith("tool.execution"):
            metrics["tool_calls"] += 1
            text = normalize_text(coerce_text(event.get("command") or event.get("name") or event.get("message") or event.get("output")))
            if text:
                commands[command_key(text)] += 1
            if any(keyword.lower() in text.lower() for keyword in SUBAGENT_KEYWORDS):
                metrics["subagent_calls"] += 1
            status = str(event.get("status") or "").lower()
            exit_code = event.get("exitCode") if "exitCode" in event else event.get("exit_code")
            if status in {"failed", "error", "timeout"} or exit_code not in {None, 0, "0"} or is_nonzero_tool_result(text):
                metrics["nonzero_tool_results"] += 1
                evidence_refs["nonzero_tool_results"].append(evidence(path, line_number, timestamp, text or event_type))

    metrics["long_idle_gaps"], evidence_refs["long_idle_gaps"] = collect_idle_gaps(timestamps, long_idle_ms)
    repeated = [(command, count) for command, count in commands.items() if count > 1]
    metrics["repeated_tool_commands"] = sum(count - 1 for _, count in repeated)
    for command, count in sorted(repeated, key=lambda item: item[1], reverse=True)[:20]:
        evidence_refs["repeated_tool_commands"].append({"count": count, "command_excerpt": truncate_text(command)})

    session = build_session(ref.runtime, session_id, path, started_at, ended_at)
    metrics["duration_ms_total"] = session["duration_ms"]
    metrics["active_duration_ms_total"] = active_duration_ms(session["duration_ms"], evidence_refs["long_idle_gaps"])
    return session, metrics, evidence_refs, gaps, first_user_prompt


def sanitize_folder_part(value: str) -> str:
    normalized = normalize_text(value)
    allowed: list[str] = []
    previous_dash = False
    for char in normalized:
        if char.isalnum() or char in {"-", "_", "."}:
            allowed.append(char)
            previous_dash = False
        else:
            if not previous_dash:
                allowed.append("-")
                previous_dash = True
    sanitized = "".join(allowed).strip("-._")
    return sanitized[:MAX_FOLDER_NAME_CHARS].strip("-._") or "untitled"


def select_run_title(results: list[tuple[dict[str, Any], dict[str, int], dict[str, list[dict[str, Any]]], list[str], str]]) -> str:
    candidates: list[tuple[datetime, str]] = []
    for session, _, _, _, prompt in results:
        if not prompt:
            continue
        started_at = parse_timestamp(session.get("started_at")) or datetime.max.replace(tzinfo=UTC)
        candidates.append((started_at, prompt))
    if not candidates:
        return "untitled run"
    return sorted(candidates, key=lambda item: item[0])[0][1]


def score(metrics: dict[str, int]) -> dict[str, int]:
    assistant_turns = max(metrics["assistant_turns"], 1)
    hours = metrics["active_duration_ms_total"] / 3_600_000
    total_turns = metrics["user_turns"] + metrics["assistant_turns"]
    tool_calls_per_assistant = metrics["tool_calls"] / assistant_turns
    return {
        "time_cost": bounded_score(hours / 6 * 100),
        "interaction_cost": bounded_score(total_turns / 120 * 100),
        "tool_churn": bounded_score(tool_calls_per_assistant / 8 * 100),
        "rework_cost": bounded_score(
            metrics["user_corrections"] * 8
            + metrics["nonzero_tool_results"] * 4
            + metrics["repeated_tool_commands"] * 1.5
        ),
    }


def load_existing_refs(path: Path) -> list[TranscriptRef]:
    if not path.exists():
        return []
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError):
        return []
    refs: list[TranscriptRef] = []
    for item in data.get("transcripts", []):
        if isinstance(item, dict) and item.get("runtime") in {"codex", "copilot"} and item.get("transcript_path"):
            refs.append(TranscriptRef(str(item["runtime"]), Path(str(item["transcript_path"]))))
    return refs


def write_json(path: Path, data: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, ensure_ascii=False, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def build_score(refs: list[TranscriptRef], long_idle_ms: int) -> tuple[dict[str, Any], str]:
    results: list[tuple[dict[str, Any], dict[str, int], dict[str, list[dict[str, Any]]], list[str], str]] = []
    for ref in refs:
        if ref.runtime == "codex":
            results.append(extract_codex(ref, long_idle_ms))
        else:
            results.append(extract_copilot(ref, long_idle_ms))

    run_title = select_run_title(results)
    sessions: list[dict[str, Any]] = []
    metrics = new_metrics()
    metrics["duration_ms_total"] = 0
    evidence_refs = empty_evidence_refs()
    transcript_gaps: list[str] = []
    for session, session_metrics, session_evidence, gaps, _ in results:
        sessions.append(session)
        merge_metrics(metrics, session_metrics)
        merge_evidence_refs(evidence_refs, session_evidence)
        transcript_gaps.extend(gaps)

    score_data = {
        "run_title": run_title,
        "sessions": sorted(sessions, key=lambda item: item.get("started_at") or ""),
        "metrics": metrics,
        "scores": score(metrics),
        "evidence_refs": evidence_refs,
        "transcript_gaps": transcript_gaps,
        "scoring_notes": [
            "スコアは close gate ではなく次回改善用の機械指標である。",
            "script は原因推定、責務判断、改善案を出さない。",
            "必要な時は evidence_refs の source_ref から transcript 原文へ戻る。",
        ],
    }
    return score_data, run_title


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Score Codex / Copilot transcripts for work_history benchmark input.")
    parser.add_argument("--codex-transcript", action="append", default=[], help="Codex transcript JSONL path. Can be repeated.")
    parser.add_argument("--copilot-transcript", action="append", default=[], help="Copilot transcript JSONL path. Can be repeated.")
    parser.add_argument("--output-root", type=Path, default=DEFAULT_OUTPUT_ROOT, help="Output root. Default: work_history/runs")
    parser.add_argument("--long-idle-ms", type=int, default=DEFAULT_LONG_IDLE_MS, help="Long idle gap threshold in milliseconds.")
    parser.add_argument("--print-run-folder", action="store_true", help="Print generated run folder path.")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    input_refs = [TranscriptRef("codex", Path(path).expanduser()) for path in args.codex_transcript]
    input_refs += [TranscriptRef("copilot", Path(path).expanduser()) for path in args.copilot_transcript]
    if not input_refs:
        raise SystemExit("at least one --codex-transcript or --copilot-transcript is required")

    preliminary_score, run_title = build_score(input_refs, args.long_idle_ms)
    first_started = preliminary_score["sessions"][0].get("started_at") if preliminary_score["sessions"] else None
    date_part = (parse_timestamp(first_started) or datetime.now(UTC)).date().isoformat()
    run_folder = args.output_root / f"{date_part}-{sanitize_folder_part(run_title)}-run"
    refs_path = run_folder / TRANSCRIPT_REFS_NAME

    merged_refs = load_existing_refs(refs_path) + input_refs
    deduped: dict[tuple[str, str], TranscriptRef] = {}
    for ref in merged_refs:
        deduped[(ref.runtime, str(ref.transcript_path))] = ref
    refs = list(deduped.values())
    score_data, run_title = build_score(refs, args.long_idle_ms)

    run_folder.mkdir(parents=True, exist_ok=True)
    (run_folder / RUN_TITLE_NAME).write_text(run_title + "\n", encoding="utf-8")
    write_json(refs_path, {"transcripts": [{"runtime": ref.runtime, "transcript_path": str(ref.transcript_path)} for ref in refs]})
    write_json(run_folder / ANALYSIS_DIR_NAME / SCORE_JSON_NAME, score_data)

    if args.print_run_folder:
        print(run_folder.resolve())
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
