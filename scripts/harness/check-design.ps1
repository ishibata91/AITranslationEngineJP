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
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\README.md"); Patterns = @("flow light, gate heavy", "Plan Stabilization Loop", "blocking unknown", "workflow-gate") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\architect.toml"); Patterns = @("architect", "Plan Stabilization Loop", "Workflow Gate", "blocking") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\research.toml"); Patterns = @("Research agent", "Read-only", "facts", "blocking") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\coder.toml"); Patterns = @("Coder agent", "implementation", "workflow gate", "Do not spawn sub-agents") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\architect-direction\SKILL.md"); Patterns = @("Plan Stabilization Loop", "blocking", "workflow-gate", "Heavy plan") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\light-direction\SKILL.md"); Patterns = @("Short Plan", "workflow-gate", "blocking unknown", "reroute") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\workflow-gate\SKILL.md"); Patterns = @("decision", "missing_evidence", "contract_breaks", "docs_sync", "recheck") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\light-review\SKILL.md"); Patterns = @("supplemental", "workflow-gate", "standard gate") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\exec-plans\templates\heavy-plan.md"); Patterns = @("Decision Basis", "Unknown Classification", "Plan Ready Criteria", "Required Evidence", "Reroute Trigger", "Docs Sync") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\exec-plans\templates\light-plan.md"); Patterns = @("Decision Basis", "Required Evidence", "Reroute Trigger", "Docs Sync") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\spec.md"); Patterns = @("LMStudio", "Gemini", "xAI", "BatchAPI", "xTranslator") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\architecture.md"); Patterns = @("Dependency Inversion Principle", "UI Port / UseCase Input", "DTO", "SQLite", "Rust") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\tech-selection.md"); Patterns = @("Tauri 2", "Rust", "Svelte 5", "SQLite", "sqlx") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\core-beliefs.md"); Patterns = @("agent-first", "Workflow Gate", "Plan Stabilization Loop", "evidence") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\index.md"); Patterns = @("AGENTS.md", ".codex", "workflow-gate", "quality-score.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\executable-specs.md"); Patterns = @("test", "acceptance checks", "Required Evidence", "workflow gate") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "4humans\quality-score.md"); Patterns = @("Codex workflow source of truth", "Design harness", "workflow gate") }
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
