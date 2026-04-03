from __future__ import annotations

import argparse
import json
import subprocess
import sys
from pathlib import Path

DOCKER_MCP_PROFILE = "codexmcps"


def normalize_repo_path(path: str | None, repo_root: Path) -> str | None:
    if path is None:
        return None
    raw = path.replace("\\", "/").strip()
    if not raw:
        return None

    candidate = Path(raw)
    if candidate.is_absolute():
        try:
            return candidate.resolve().relative_to(repo_root).as_posix()
        except ValueError:
            return candidate.resolve().as_posix()

    normalized = candidate.as_posix()
    while normalized.startswith("./"):
        normalized = normalized[2:]
    return normalized or None


def get_issue_repo_path(issue: dict, component_path_by_key: dict[str, str], repo_root: Path) -> str | None:
    component = issue.get("component")
    if isinstance(component, str) and ":" in component:
        return normalize_repo_path(component.split(":", 1)[1], repo_root)
    if isinstance(component, str) and component in component_path_by_key:
        return normalize_repo_path(component_path_by_key[component], repo_root)
    return None


def matches_owned_scope(repo_path: str | None, owned_paths: list[str]) -> bool:
    if not owned_paths:
        return True
    if repo_path is None:
        return False
    return any(repo_path == owned_path or repo_path.startswith(f"{owned_path}/") for owned_path in owned_paths)


def extract_json_payload(stdout: str) -> dict:
    payload_start = stdout.find("{")
    if payload_start == -1:
        raise ValueError("Docker MCP output does not contain a JSON object.")
    try:
        payload = json.loads(stdout[payload_start:])
    except json.JSONDecodeError as exc:
        raise ValueError(f"Failed to decode Docker MCP JSON payload: {exc}") from exc
    if not isinstance(payload, dict):
        raise ValueError("Docker MCP payload must be a JSON object.")
    return payload


def load_project_issues(project: str) -> dict:
    payload = json.dumps({"projects": [project], "page": 1, "pageSize": 500})
    completed = subprocess.run(
        [
            "docker",
            "mcp",
            "tools",
            "call",
            "search_sonar_issues_in_projects",
            "--gateway-arg=--profile",
            f"--gateway-arg={DOCKER_MCP_PROFILE}",
            payload,
        ],
        check=False,
        capture_output=True,
        text=True,
    )
    if completed.returncode != 0:
        stderr = completed.stderr.strip() or "docker mcp tools call failed without stderr output."
        raise RuntimeError(f"Docker MCP Sonar issue query failed for profile '{DOCKER_MCP_PROFILE}': {stderr}")

    payload_obj = extract_json_payload(completed.stdout)
    issues = payload_obj.get("issues")
    if not isinstance(issues, list):
        raise ValueError("Docker MCP Sonar payload must contain an 'issues' list.")
    return payload_obj


def main() -> int:
    parser = argparse.ArgumentParser(description="List open Sonar issues for the given project.")
    parser.add_argument("--project", required=True)
    parser.add_argument("--owned-paths", nargs="*", default=[])
    args = parser.parse_args()
    repo_root = Path.cwd().resolve()

    try:
        payload = load_project_issues(args.project)
    except (RuntimeError, ValueError) as exc:
        sys.stderr.write(f"{exc}\n")
        return 1

    component_path_by_key: dict[str, str] = {}
    normalized_owned_paths = [
        normalized
        for normalized in (normalize_repo_path(path, repo_root) for path in args.owned_paths)
        if normalized is not None
    ]

    open_issues: list[dict] = []
    for issue in payload.get("issues", []):
        if (
            not isinstance(issue, dict)
            or issue.get("project") != args.project
            or issue.get("status") != "OPEN"
        ):
            continue

        repo_path = get_issue_repo_path(issue, component_path_by_key, repo_root)
        if not matches_owned_scope(repo_path, normalized_owned_paths):
            continue

        open_issues.append(
            {
                "key": issue.get("key"),
                "rule": issue.get("rule"),
                "severity": issue.get("severity"),
                "status": issue.get("status"),
                "issueStatus": issue.get("issueStatus"),
                "message": issue.get("message"),
                "component": issue.get("component"),
                "repoPath": repo_path,
                "textRange": issue.get("textRange"),
            }
        )

    result = {
        "project": args.project,
        "totalIssues": len(payload.get("issues", [])),
        "openIssueCount": len(open_issues),
        "openIssues": open_issues,
    }
    json.dump(result, sys.stdout, ensure_ascii=False, indent=2)
    sys.stdout.write("\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
