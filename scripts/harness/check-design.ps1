param(
    [string]$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Assert-Patterns {
    param(
        [string]$FilePath,
        [string[]]$Patterns,
        [ref]$Failures
    )

    $content = Get-Content -LiteralPath $FilePath -Raw
    foreach ($pattern in $Patterns) {
        if ($content -notmatch $pattern) {
            Write-Host "FAIL missing pattern '$pattern' in $FilePath" -ForegroundColor Red
            $Failures.Value++
        } else {
            Write-Host "PASS pattern '$pattern' found in $FilePath" -ForegroundColor Green
        }
    }
}

$checks = @(
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\README.md"); Patterns = @("directing-implementation", "directing-fixes", "architecting-tests", "UI", "Scenario", "Logic", "single-pass", "changes/", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\ctx_loader.toml"); Patterns = @("facts", "constraints", "gaps", "exec-plan") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\workplan_builder.toml"); Patterns = @("UI", "Scenario", "Logic", "Implementation Plan", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\test_architect.toml"); Patterns = @("test_architect", "spec-aligned tests", "fixtures", "validation commands") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\implementer.toml"); Patterns = @("allowed scope", "validation", "touched files") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\fault_tracer.toml"); Patterns = @("root-cause hypotheses", "observation plan", "temporary logging") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\log_instrumenter.toml"); Patterns = @("temporary log statements", "\[tracing-fixes\]", "Remove any temporary instrumentation") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\review_cycler.toml"); Patterns = @("single-pass", "spec deviation", "exception handling", "resource cleanup", "missing tests", "pass", "reroute") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\directing-implementation\SKILL.md"); Patterns = @("docs/exec-plans/templates/impl-plan.md", "UI", "Scenario", "Logic", "architecting-tests", "reviewing-implementation", "reroute", "changes/", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\directing-fixes\SKILL.md"); Patterns = @("docs/exec-plans/templates/fix-plan.md", "distilling-fixes", "tracing-fixes", "logging-fixes", "architecting-tests", "reviewing-fixes", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\architecting-tests\SKILL.md"); Patterns = @("Acceptance Checks", "failing tests", "fixtures", "validation commands", "docs sync") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\reviewing-implementation\SKILL.md"); Patterns = @("Review Scope", "pass", "reroute", "score") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\reviewing-fixes\SKILL.md"); Patterns = @("Review Scope", "pass", "reroute") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\exec-plans\templates\impl-plan.md"); Patterns = @("Decision Basis", "UI", "Scenario", "Logic", "Implementation Plan", "Required Evidence", "Docs Sync", "Outcome") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\exec-plans\templates\fix-plan.md"); Patterns = @("Decision Basis", "Known Facts", "Trace Plan", "Fix Plan", "Required Evidence", "Docs Sync", "Outcome") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\spec.md"); Patterns = @("LMStudio", "Gemini", "xAI", "BatchAPI", "xTranslator") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\architecture.md"); Patterns = @("Dependency Inversion Principle", "UI Port / UseCase Input", "DTO", "SQLite", "Rust") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\tech-selection.md"); Patterns = @("Tauri 2", "Rust", "Svelte 5", "SQLite", "sqlx") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\core-beliefs.md"); Patterns = @("agent-first", "directing-implementation", "directing-fixes", "single-pass", "evidence") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\index.md"); Patterns = @("AGENTS.md", ".codex", "directing-implementation", "directing-fixes", "quality-score.md", "tests / acceptance checks / validation commands") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "4humans\quality-score.md"); Patterns = @("Codex workflow source of truth", "Design harness", "directing-implementation", "directing-fixes") }
)

$failures = 0

foreach ($check in $checks) {
    Assert-Patterns -FilePath $check.File -Patterns $check.Patterns -Failures ([ref]$failures)
}

if ($failures -gt 0) {
    Write-Error "Design harness failed with $failures issue(s)."
}

Write-Host ""
Write-Host "Design harness passed." -ForegroundColor Green

