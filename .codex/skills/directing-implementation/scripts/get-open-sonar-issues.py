from __future__ import annotations

import argparse
import json
import subprocess
import sys


def normalize_repo_path(path: str | None) -> str | None:
    if path is None:
        return None
    normalized = path.replace("\\", "/").lstrip("./").strip()
    return normalized or None


def get_issue_repo_path(issue: dict, component_path_by_key: dict[str, str]) -> str | None:
    component = issue.get("component")
    if isinstance(component, str) and ":" in component:
        return normalize_repo_path(component.split(":", 1)[1])
    if isinstance(component, str) and component in component_path_by_key:
        return normalize_repo_path(component_path_by_key[component])
    return None


def matches_owned_scope(repo_path: str | None, owned_paths: list[str]) -> bool:
    if not owned_paths:
        return True
    if repo_path is None:
        return False
    return any(repo_path == owned_path or repo_path.startswith(f"{owned_path}/") for owned_path in owned_paths)


def main() -> int:
    parser = argparse.ArgumentParser(description="List open Sonar issues for the given project.")
    parser.add_argument("--project", required=True)
    parser.add_argument("--owned-paths", nargs="*", default=[])
    args = parser.parse_args()

    completed = subprocess.run(
        ["sonar", "list", "issues", "--project", args.project, "--format", "json"],
        check=False,
        capture_output=True,
        text=True,
    )
    if completed.returncode != 0:
        sys.stderr.write(completed.stderr)
        return completed.returncode

    payload = json.loads(completed.stdout)
    component_path_by_key = {
        component["key"]: component["path"]
        for component in payload.get("components", [])
        if isinstance(component, dict)
        and isinstance(component.get("key"), str)
        and isinstance(component.get("path"), str)
        and component["path"].strip()
    }
    normalized_owned_paths = [
        normalized for normalized in (normalize_repo_path(path) for path in args.owned_paths) if normalized is not None
    ]

    open_issues: list[dict] = []
    for issue in payload.get("issues", []):
        if not isinstance(issue, dict) or issue.get("status") != "OPEN":
            continue

        repo_path = get_issue_repo_path(issue, component_path_by_key)
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
