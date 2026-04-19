from __future__ import annotations

import json
import re
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from dataclasses import dataclass
from pathlib import Path

from harness_common import (
    build_parser,
    default_repo_root,
    finalize_failures,
    find_command,
    load_json,
    report_fail,
    report_pass,
    report_run,
    report_skip,
    run_command,
    run_command_capture,
)

MINIMUM_COVERAGE = 70.0
MAX_SECURITY_ISSUES = 0
MAX_RELIABILITY_ISSUES = 0
MAX_MAINTAINABILITY_HIGH_ISSUES = 0
FRONTEND_SUMMARY_PATH = Path("test-results/frontend-coverage/coverage-summary.json")
FRONTEND_LCOV_PATH = Path("test-results/frontend-coverage/lcov.info")
BACKEND_SUMMARY_PATH = Path("test-results/backend-coverage/coverage-summary.txt")
BACKEND_COVERAGE_PATH = Path("test-results/backend-coverage/coverage.out")
MANIFEST_PATH = Path("test-results/coverage-manifest.json")
SONAR_PROPERTIES_PATH = Path("sonar-project.properties")
SONAR_REPORT_TASK_PATH = Path(".scannerwork/report-task.txt")
BACKEND_TOTAL_PATTERN = re.compile(r"total:\s+\(statements\)\s+(?P<pct>[0-9]+(?:\.[0-9]+)?)%")
BACKEND_PROFILE_PATTERN = re.compile(
    r"^(?P<file>[^:]+):(?P<start_line>\d+)\.\d+,(?P<end_line>\d+)\.\d+\s+\d+\s+(?P<count>\d+)$"
)


@dataclass(frozen=True)
class FrontendCoverage:
    statements_pct: float
    lines_to_cover: int
    covered_lines: int
    branches_to_cover: int
    covered_branches: int

    @property
    def line_pct(self) -> float:
        if self.lines_to_cover == 0:
            return 100.0
        return (self.covered_lines / self.lines_to_cover) * 100.0


@dataclass(frozen=True)
class BackendCoverage:
    statements_pct: float
    lines_to_cover: int
    covered_lines: int

    @property
    def line_pct(self) -> float:
        if self.lines_to_cover == 0:
            return 100.0
        return (self.covered_lines / self.lines_to_cover) * 100.0


@dataclass(frozen=True)
class OverallCoverage:
    coverage_pct: float
    line_coverage_pct: float
    lines_to_cover: int
    covered_lines: int
    branches_to_cover: int
    covered_branches: int


@dataclass(frozen=True)
class SonarIssues:
    security: int
    reliability: int
    maintainability_high: int


@dataclass(frozen=True)
class SonarCoverage:
    coverage_pct: float
    line_coverage_pct: float
    branch_coverage_pct: float
    lines_to_cover: int
    uncovered_lines: int


def resolve_package_manager(directory: Path) -> str:
    if (directory / "pnpm-lock.yaml").exists():
        return "pnpm"
    if (directory / "package-lock.json").exists():
        return "npm"
    if (directory / "yarn.lock").exists():
        return "yarn"
    return "npm"


def has_script(package: dict, script_name: str) -> bool:
    scripts = package.get("scripts")
    return isinstance(scripts, dict) and script_name in scripts


def invoke_step(command: str, arguments: list[str], working_directory: Path) -> int:
    if find_command(command) is None:
        report_fail(f"FAIL missing command: {command}")
        return 1

    rendered_command = f"{command} {' '.join(arguments)}"
    report_run(f"RUN {rendered_command}")
    exit_code = run_command(command, arguments, working_directory)
    if exit_code != 0:
        report_fail(f"FAIL {rendered_command}")
        return 1

    report_pass(f"PASS {rendered_command}")
    return 0


def invoke_step_capture(command: str, arguments: list[str], working_directory: Path) -> tuple[int, str]:
    if find_command(command) is None:
        report_fail(f"FAIL missing command: {command}")
        return 1, ""

    rendered_command = f"{command} {' '.join(arguments)}"
    report_run(f"RUN {rendered_command}")
    completed = run_command_capture(command, arguments, working_directory)
    output = completed.stdout + completed.stderr
    if output:
        print(output, end="" if output.endswith("\n") else "\n")
    if completed.returncode != 0:
        report_fail(f"FAIL {rendered_command}")
        return 1, output

    report_pass(f"PASS {rendered_command}")
    return 0, output


def load_frontend_coverage(summary_path: Path) -> FrontendCoverage | None:
    if not summary_path.exists():
        return None

    summary = json.loads(summary_path.read_text(encoding="utf-8"))
    total = summary.get("total")
    if not isinstance(total, dict):
        return None

    statements = total.get("statements")
    lines = total.get("lines")
    branches = total.get("branches")
    if not isinstance(statements, dict) or not isinstance(lines, dict) or not isinstance(branches, dict):
        return None

    statements_pct = statements.get("pct")
    lines_total = lines.get("total")
    lines_covered = lines.get("covered")
    branches_total = branches.get("total")
    branches_covered = branches.get("covered")
    if not isinstance(statements_pct, int | float):
        return None
    if not isinstance(lines_total, int) or not isinstance(lines_covered, int):
        return None
    if not isinstance(branches_total, int) or not isinstance(branches_covered, int):
        return None

    return FrontendCoverage(
        statements_pct=float(statements_pct),
        lines_to_cover=lines_total,
        covered_lines=lines_covered,
        branches_to_cover=branches_total,
        covered_branches=branches_covered,
    )


def load_backend_statements(summary_path: Path) -> float | None:
    if not summary_path.exists():
        return None

    match = BACKEND_TOTAL_PATTERN.search(summary_path.read_text(encoding="utf-8"))
    if match is None:
        return None

    return float(match.group("pct"))


def load_backend_line_coverage(profile_path: Path, statements_pct: float | None) -> BackendCoverage | None:
    if statements_pct is None or not profile_path.exists():
        return None

    coverable_by_file: dict[str, set[int]] = {}
    covered_by_file: dict[str, set[int]] = {}

    for raw_line in profile_path.read_text(encoding="utf-8").splitlines()[1:]:
        match = BACKEND_PROFILE_PATTERN.match(raw_line)
        if match is None:
            continue

        file_path = match.group("file")
        start_line = int(match.group("start_line"))
        end_line = int(match.group("end_line"))
        execution_count = int(match.group("count"))

        coverable_lines = coverable_by_file.setdefault(file_path, set())
        covered_lines = covered_by_file.setdefault(file_path, set())
        for line_number in range(start_line, end_line + 1):
            coverable_lines.add(line_number)
            if execution_count > 0:
                covered_lines.add(line_number)

    if not coverable_by_file:
        return None

    lines_to_cover = sum(len(line_numbers) for line_numbers in coverable_by_file.values())
    covered_lines = sum(len(covered_by_file.get(file_path, set())) for file_path in coverable_by_file)
    return BackendCoverage(
        statements_pct=statements_pct,
        lines_to_cover=lines_to_cover,
        covered_lines=covered_lines,
    )


def combine_coverage(frontend: FrontendCoverage | None, backend: BackendCoverage | None) -> OverallCoverage | None:
    if frontend is None or backend is None:
        return None

    lines_to_cover = frontend.lines_to_cover + backend.lines_to_cover
    covered_lines = frontend.covered_lines + backend.covered_lines
    branches_to_cover = frontend.branches_to_cover
    covered_branches = frontend.covered_branches

    line_coverage_pct = 100.0 if lines_to_cover == 0 else (covered_lines / lines_to_cover) * 100.0
    denominator = lines_to_cover + branches_to_cover
    coverage_pct = 100.0 if denominator == 0 else ((covered_lines + covered_branches) / denominator) * 100.0

    return OverallCoverage(
        coverage_pct=coverage_pct,
        line_coverage_pct=line_coverage_pct,
        lines_to_cover=lines_to_cover,
        covered_lines=covered_lines,
        branches_to_cover=branches_to_cover,
        covered_branches=covered_branches,
    )


def parse_properties(path: Path) -> dict[str, str]:
    properties: dict[str, str] = {}
    if not path.exists():
        return properties

    for raw_line in path.read_text(encoding="utf-8").splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        properties[key.strip()] = value.strip()
    return properties


def request_json(url: str) -> dict | None:
    try:
        with urllib.request.urlopen(url, timeout=30) as response:
            return json.loads(response.read().decode("utf-8"))
    except urllib.error.HTTPError:
        return None


def load_sonar_coverage(report_task_path: Path, timeout_seconds: int = 60) -> SonarCoverage | None:
    report_task = parse_properties(report_task_path)
    server_url = report_task.get("serverUrl")
    project_key = report_task.get("projectKey")
    if not server_url or not project_key:
        return None

    query = urllib.parse.urlencode(
        {
            "component": project_key,
            "metricKeys": "coverage,line_coverage,branch_coverage,lines_to_cover,uncovered_lines",
        }
    )
    measures_url = f"{server_url}/api/measures/component?{query}"

    payload = None
    deadline = time.monotonic() + timeout_seconds
    while time.monotonic() < deadline:
        payload = request_json(measures_url)
        component = None if payload is None else payload.get("component")
        measures = None if not isinstance(component, dict) else component.get("measures")
        if isinstance(measures, list) and measures:
            break
        time.sleep(2)

    if payload is None:
        return None
    component = payload.get("component")
    if not isinstance(component, dict):
        return None

    raw_measures = component.get("measures")
    if not isinstance(raw_measures, list):
        return None

    measures: dict[str, str] = {}
    for entry in raw_measures:
        if not isinstance(entry, dict):
            continue
        metric = entry.get("metric")
        value = entry.get("value")
        if isinstance(metric, str) and isinstance(value, str):
            measures[metric] = value

    try:
        return SonarCoverage(
            coverage_pct=float(measures["coverage"]),
            line_coverage_pct=float(measures["line_coverage"]),
            branch_coverage_pct=float(measures.get("branch_coverage", "0.0")),
            lines_to_cover=int(float(measures["lines_to_cover"])),
            uncovered_lines=int(float(measures["uncovered_lines"])),
        )
    except KeyError:
        return None


def load_sonar_issues(report_task_path: Path) -> SonarIssues | None:
    report_task = parse_properties(report_task_path)
    server_url = report_task.get("serverUrl")
    project_key = report_task.get("projectKey")
    if not server_url or not project_key:
        return None

    def count_issues(qualities: list[str], severities: list[str] | None = None) -> int | None:
        params: dict[str, str] = {
            "componentKeys": project_key,
            "impactSoftwareQualities": ",".join(qualities),
            "resolved": "false",
            "ps": "1",
        }
        if severities:
            params["impactSeverities"] = ",".join(severities)
        query = urllib.parse.urlencode(params)
        url = f"{server_url}/api/issues/search?{query}"
        payload = request_json(url)
        if payload is None:
            return None
        total = payload.get("total")
        return int(total) if isinstance(total, int | float) else None

    security = count_issues(["SECURITY"])
    reliability = count_issues(["RELIABILITY"])
    maintainability_high = count_issues(["MAINTAINABILITY"], ["HIGH"])

    if security is None or reliability is None or maintainability_high is None:
        return None

    return SonarIssues(
        security=security,
        reliability=reliability,
        maintainability_high=maintainability_high,
    )


def run_sonar_scan(repo_root: Path, package_manager: str, root_package: dict) -> SonarCoverage | None:
    if not (repo_root / SONAR_PROPERTIES_PATH).exists():
        return None

    if has_script(root_package, "scan:sonar"):
        failures, _ = invoke_step_capture(package_manager, ["run", "scan:sonar"], repo_root)
    else:
        failures, _ = invoke_step_capture("sonar-scanner", [], repo_root)
    if failures != 0:
        return None

    sonar_coverage = load_sonar_coverage(repo_root / SONAR_REPORT_TASK_PATH)
    if sonar_coverage is None:
        report_fail("FAIL Sonar coverage measures could not be loaded")
        return None

    report_pass(
        "PASS Sonar coverage summary "
        f"coverage={sonar_coverage.coverage_pct:.1f}% "
        f"line={sonar_coverage.line_coverage_pct:.1f}% "
        f"branch={sonar_coverage.branch_coverage_pct:.1f}%"
    )
    return sonar_coverage


def check_threshold(
    overall: OverallCoverage | None,
    sonar: SonarCoverage | None,
    issues: SonarIssues | None,
) -> int:
    failures = 0

    if sonar is not None:
        if sonar.coverage_pct < MINIMUM_COVERAGE:
            report_fail(
                "FAIL Sonar coverage "
                f"{sonar.coverage_pct:.1f}% < {MINIMUM_COVERAGE:.1f}% "
                f"(line={sonar.line_coverage_pct:.1f}%, branch={sonar.branch_coverage_pct:.1f}%)"
            )
            failures += 1
        else:
            report_pass(
                "PASS Sonar coverage "
                f"{sonar.coverage_pct:.1f}% >= {MINIMUM_COVERAGE:.1f}% "
                f"(line={sonar.line_coverage_pct:.1f}%, branch={sonar.branch_coverage_pct:.1f}%)"
            )
    elif overall is None:
        report_fail("FAIL Sonar-compatible coverage summary could not be parsed")
        failures += 1
    elif overall.coverage_pct < MINIMUM_COVERAGE:
        report_fail(
            "FAIL Sonar-compatible coverage "
            f"{overall.coverage_pct:.1f}% < {MINIMUM_COVERAGE:.1f}% "
            f"(line={overall.line_coverage_pct:.1f}%, branches={overall.covered_branches}/{overall.branches_to_cover})"
        )
        failures += 1
    else:
        report_pass(
            "PASS Sonar-compatible coverage "
            f"{overall.coverage_pct:.1f}% >= {MINIMUM_COVERAGE:.1f}% "
            f"(line={overall.line_coverage_pct:.1f}%, branches={overall.covered_branches}/{overall.branches_to_cover})"
        )

    # issues gate (Sonar から取得できた場合のみ)
    if issues is not None:
        if issues.security > MAX_SECURITY_ISSUES:
            report_fail(f"FAIL Sonar security issues {issues.security} > {MAX_SECURITY_ISSUES}")
            failures += 1
        else:
            report_pass(f"PASS Sonar security issues {issues.security} <= {MAX_SECURITY_ISSUES}")

        if issues.reliability > MAX_RELIABILITY_ISSUES:
            report_fail(f"FAIL Sonar reliability issues {issues.reliability} > {MAX_RELIABILITY_ISSUES}")
            failures += 1
        else:
            report_pass(f"PASS Sonar reliability issues {issues.reliability} <= {MAX_RELIABILITY_ISSUES}")

        if issues.maintainability_high > MAX_MAINTAINABILITY_HIGH_ISSUES:
            report_fail(
                f"FAIL Sonar maintainability HIGH issues {issues.maintainability_high} > {MAX_MAINTAINABILITY_HIGH_ISSUES}"
            )
            failures += 1
        else:
            report_pass(
                f"PASS Sonar maintainability HIGH issues {issues.maintainability_high} <= {MAX_MAINTAINABILITY_HIGH_ISSUES}"
            )

    return failures


def write_manifest(
    repo_root: Path,
    frontend: FrontendCoverage | None,
    backend: BackendCoverage | None,
    overall: OverallCoverage | None,
    sonar: SonarCoverage | None,
    issues: SonarIssues | None,
) -> None:
    manifest = {
        "minimum_sonar_coverage": MINIMUM_COVERAGE,
        "overall": {
            "coverage_pct": None if overall is None else overall.coverage_pct,
            "line_coverage_pct": None if overall is None else overall.line_coverage_pct,
            "lines_to_cover": None if overall is None else overall.lines_to_cover,
            "covered_lines": None if overall is None else overall.covered_lines,
            "branches_to_cover": None if overall is None else overall.branches_to_cover,
            "covered_branches": None if overall is None else overall.covered_branches,
        },
        "sonar": {
            "coverage_pct": None if sonar is None else sonar.coverage_pct,
            "line_coverage_pct": None if sonar is None else sonar.line_coverage_pct,
            "branch_coverage_pct": None if sonar is None else sonar.branch_coverage_pct,
            "lines_to_cover": None if sonar is None else sonar.lines_to_cover,
            "uncovered_lines": None if sonar is None else sonar.uncovered_lines,
        },
        "sonar_issues": {
            "security": None if issues is None else issues.security,
            "reliability": None if issues is None else issues.reliability,
            "maintainability_high": None if issues is None else issues.maintainability_high,
        },
        "frontend": {
            "statements_pct": None if frontend is None else frontend.statements_pct,
            "line_coverage_pct": None if frontend is None else frontend.line_pct,
            "lines_to_cover": None if frontend is None else frontend.lines_to_cover,
            "covered_lines": None if frontend is None else frontend.covered_lines,
            "branches_to_cover": None if frontend is None else frontend.branches_to_cover,
            "covered_branches": None if frontend is None else frontend.covered_branches,
            "summary_path": str(FRONTEND_SUMMARY_PATH),
            "sonar_report_path": str(FRONTEND_LCOV_PATH),
            "sonar_property": "sonar.javascript.lcov.reportPaths",
            "report_exists": (repo_root / FRONTEND_LCOV_PATH).exists(),
        },
        "backend": {
            "statements_pct": None if backend is None else backend.statements_pct,
            "line_coverage_pct": None if backend is None else backend.line_pct,
            "lines_to_cover": None if backend is None else backend.lines_to_cover,
            "covered_lines": None if backend is None else backend.covered_lines,
            "summary_path": str(BACKEND_SUMMARY_PATH),
            "sonar_report_path": str(BACKEND_COVERAGE_PATH),
            "sonar_property": "sonar.go.coverage.reportPaths",
            "report_exists": (repo_root / BACKEND_COVERAGE_PATH).exists(),
        },
    }
    (repo_root / MANIFEST_PATH).write_text(
        json.dumps(manifest, ensure_ascii=True, indent=2) + "\n",
        encoding="utf-8",
    )


def main() -> int:
    parser = build_parser("Run the unit coverage harness.", default_repo_root(__file__))
    args = parser.parse_args()
    repo_root = Path(args.repo_root).resolve()

    package_json_path = repo_root / "package.json"
    if not package_json_path.exists():
        report_skip(f"SKIP no package.json found at {package_json_path}")
        return finalize_failures("Coverage harness", 0)

    package = load_json(package_json_path)
    package_manager = resolve_package_manager(repo_root)
    failures = 0
    ran_anything = False

    frontend: FrontendCoverage | None = None
    backend: BackendCoverage | None = None

    if has_script(package, "test:frontend:coverage"):
        ran_anything = True
        failures += invoke_step(package_manager, ["run", "test:frontend:coverage"], repo_root)
        frontend = load_frontend_coverage(repo_root / FRONTEND_SUMMARY_PATH)
        if frontend is None:
            report_fail("FAIL frontend coverage summary could not be parsed")
            failures += 1
        else:
            report_pass(
                "PASS frontend coverage summary "
                f"statements={frontend.statements_pct:.1f}% lines={frontend.line_pct:.1f}% "
                f"branches={frontend.covered_branches}/{frontend.branches_to_cover}"
            )
    else:
        report_skip(f"SKIP no test:frontend:coverage script in {package_json_path}")

    if has_script(package, "test:backend:coverage"):
        ran_anything = True
        failures += invoke_step(package_manager, ["run", "test:backend:coverage"], repo_root)
        backend_statements_pct = load_backend_statements(repo_root / BACKEND_SUMMARY_PATH)
        backend = load_backend_line_coverage(repo_root / BACKEND_COVERAGE_PATH, backend_statements_pct)
        if backend is None:
            report_fail("FAIL backend coverage summary could not be parsed")
            failures += 1
        else:
            report_pass(
                "PASS backend coverage summary "
                f"statements={backend.statements_pct:.1f}% lines={backend.line_pct:.1f}%"
            )
    else:
        report_skip(f"SKIP no test:backend:coverage script in {package_json_path}")

    overall = combine_coverage(frontend, backend)
    sonar = None
    if ran_anything:
        if (repo_root / SONAR_PROPERTIES_PATH).exists():
            sonar = run_sonar_scan(repo_root, package_manager, package)
            if sonar is None:
                failures += 1
        issues = load_sonar_issues(repo_root / SONAR_REPORT_TASK_PATH) if sonar is not None else None
        failures += check_threshold(overall, sonar, issues)
        write_manifest(repo_root, frontend, backend, overall, sonar, issues)
        report_pass(f"PASS wrote coverage manifest: {repo_root / MANIFEST_PATH}")
    else:
        report_skip("SKIP no coverage scripts found. Coverage harness is installed but has nothing to run yet.")

    return finalize_failures("Coverage harness", failures)


if __name__ == "__main__":
    sys.exit(main())
