from __future__ import annotations

import argparse
import json
import sys
from collections import Counter, defaultdict
from pathlib import Path
from typing import Any

REQUIRED_FIELDS = {
    "event_type",
    "runtime",
    "elapsed_ms_from_run_start",
    "phase",
    "status",
    "mechanical_summary",
}


def load_events(path: Path) -> tuple[list[dict[str, Any]], list[str]]:
    if not path.exists():
        return [], [f"missing telemetry file: {path}"]

    events: list[dict[str, Any]] = []
    gaps: list[str] = []
    for line_number, line in enumerate(path.read_text(encoding="utf-8").splitlines(), start=1):
        if not line.strip():
            continue
        try:
            event = json.loads(line)
        except json.JSONDecodeError as error:
            gaps.append(f"invalid json at line {line_number}: {error.msg}")
            continue
        if not isinstance(event, dict):
            gaps.append(f"non-object event at line {line_number}")
            continue
        missing = sorted(REQUIRED_FIELDS - event.keys())
        if missing:
            gaps.append(f"line {line_number} missing fields: {', '.join(missing)}")
        events.append(event)
    return events, gaps


def summarize(events: list[dict[str, Any]], gaps: list[str]) -> dict[str, Any]:
    response_count_by_runtime: Counter[str] = Counter()
    phase_elapsed: dict[str, int] = defaultdict(int)
    runtime_elapsed: dict[str, int] = defaultdict(int)
    blocked_elapsed = 0
    validation_elapsed = 0
    reroute_count = 0
    elapsed_values: list[int] = []

    for event in events:
        runtime = str(event.get("runtime") or "unknown")
        phase = str(event.get("phase") or "unknown")
        elapsed = coerce_int(event.get("elapsed_ms_from_run_start"))
        duration = coerce_int(event.get("duration_ms"))

        if event.get("event_type") == "assistant_response":
            response_count_by_runtime[runtime] += 1
        if elapsed is not None:
            elapsed_values.append(elapsed)
        if duration is not None:
            phase_elapsed[phase] += duration
            runtime_elapsed[runtime] += duration
            if event.get("blocked_reason"):
                blocked_elapsed += duration
            if "validation" in phase:
                validation_elapsed += duration
        if str(event.get("status")) in {"rerouted", "blocked_after_narrowing"}:
            reroute_count += 1

    return {
        "response_count_by_runtime": dict(sorted(response_count_by_runtime.items())),
        "elapsed_ms_total": max(elapsed_values) if elapsed_values else None,
        "phase_elapsed_ms": dict(sorted(phase_elapsed.items())),
        "runtime_elapsed_ms": dict(sorted(runtime_elapsed.items())),
        "blocked_elapsed_ms": blocked_elapsed or None,
        "reroute_count": reroute_count,
        "validation_elapsed_ms": validation_elapsed or None,
        "telemetry_gap": gaps or ["なし"],
    }


def coerce_int(value: Any) -> int | None:
    if isinstance(value, bool):
        return None
    if isinstance(value, int):
        return value
    if isinstance(value, str) and value.isdigit():
        return int(value)
    return None


def format_value(value: Any) -> str:
    if value is None:
        return "不明"
    if isinstance(value, dict):
        if not value:
            return "不明"
        return ", ".join(f"{key}: {item}" for key, item in value.items())
    if isinstance(value, list):
        return " / ".join(str(item) for item in value)
    return str(value)


def render_markdown(summary: dict[str, Any], telemetry_file: Path) -> str:
    lines = [
        "## Benchmark",
        "",
        f"- `telemetry_file`: `{telemetry_file}`",
        f"- `response_count_by_runtime`: `{format_value(summary['response_count_by_runtime'])}`",
        f"- `elapsed_ms_total`: `{format_value(summary['elapsed_ms_total'])}`",
        f"- `phase_elapsed_ms`: `{format_value(summary['phase_elapsed_ms'])}`",
        f"- `runtime_elapsed_ms`: `{format_value(summary['runtime_elapsed_ms'])}`",
        f"- `blocked_elapsed_ms`: `{format_value(summary['blocked_elapsed_ms'])}`",
        f"- `reroute_count`: `{format_value(summary['reroute_count'])}`",
        f"- `validation_elapsed_ms`: `{format_value(summary['validation_elapsed_ms'])}`",
        f"- `telemetry_gap`: `{format_value(summary['telemetry_gap'])}`",
        "- `benchmark_use`: `次回改善用。初期 close 判定には使わない。`",
    ]
    return "\n".join(lines) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser(description="Aggregate run telemetry JSONL into a benchmark markdown block.")
    parser.add_argument("telemetry_file", type=Path)
    parser.add_argument("--json", action="store_true", help="Print summary as JSON instead of Markdown.")
    args = parser.parse_args()

    events, gaps = load_events(args.telemetry_file)
    summary = summarize(events, gaps)
    if args.json:
        print(json.dumps(summary, ensure_ascii=False, indent=2, sort_keys=True))
    else:
        print(render_markdown(summary, args.telemetry_file), end="")
    return 0


if __name__ == "__main__":
    sys.exit(main())
