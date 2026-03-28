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
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\README.md"); Patterns = @("impl-direction", "fix-direction", "test-architect", "UI", "Scenario", "Logic", "single-pass", "changes/", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\ctx_loader.toml"); Patterns = @("facts", "constraints", "gaps", "exec-plan") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\workplan_builder.toml"); Patterns = @("UI", "Scenario", "Logic", "Implementation Plan", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\test_architect.toml"); Patterns = @("test_architect", "spec-aligned tests", "fixtures", "validation commands") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\implementer.toml"); Patterns = @("allowed scope", "validation", "touched files") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\fault_tracer.toml"); Patterns = @("root-cause hypotheses", "observation plan", "temporary logging") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\log_instrumenter.toml"); Patterns = @("temporary log statements", "\[fix-trace\]", "Remove any temporary instrumentation") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\review_cycler.toml"); Patterns = @("single-pass", "spec deviation", "exception handling", "resource cleanup", "missing tests", "pass", "reroute") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\impl-direction\SKILL.md"); Patterns = @("docs/exec-plans/templates/impl-plan.md", "UI", "Scenario", "Logic", "test-architect", "impl-review", "reroute", "changes/", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\fix-direction\SKILL.md"); Patterns = @("docs/exec-plans/templates/fix-plan.md", "fix-distill", "fix-trace", "fix-logging", "test-architect", "fix-review", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\test-architect\SKILL.md"); Patterns = @("Acceptance Checks", "failing tests", "fixtures", "validation commands", "docs sync") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\impl-review\SKILL.md"); Patterns = @("Review Scope", "pass", "reroute", "score") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\fix-review\SKILL.md"); Patterns = @("Review Scope", "pass", "reroute") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\exec-plans\templates\impl-plan.md"); Patterns = @("Decision Basis", "UI", "Scenario", "Logic", "Implementation Plan", "Required Evidence", "Docs Sync", "Outcome") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\exec-plans\templates\fix-plan.md"); Patterns = @("Decision Basis", "Known Facts", "Trace Plan", "Fix Plan", "Required Evidence", "Docs Sync", "Outcome") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\spec.md"); Patterns = @("LMStudio", "Gemini", "xAI", "BatchAPI", "xTranslator") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\architecture.md"); Patterns = @("Dependency Inversion Principle", "UI Port / UseCase Input", "DTO", "SQLite", "Rust") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\tech-selection.md"); Patterns = @("Tauri 2", "Rust", "Svelte 5", "SQLite", "sqlx") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\core-beliefs.md"); Patterns = @("agent-first", "impl-direction", "fix-direction", "single-pass", "evidence") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\index.md"); Patterns = @("AGENTS.md", ".codex", "impl-direction", "fix-direction", "quality-score.md", "tests / acceptance checks / validation commands") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "4humans\quality-score.md"); Patterns = @("Codex workflow source of truth", "Design harness", "impl-direction", "fix-direction") }
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
