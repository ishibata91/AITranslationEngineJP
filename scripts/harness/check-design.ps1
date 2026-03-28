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
    @{
        File = (Join-Path $RepoRoot ".codex\README.md")
        Patterns = @("heavy", "light", "architect-direction", "light-direction")
    },
    @{
        File = (Join-Path $RepoRoot ".codex\agents\architect.toml")
        Patterns = @("architect", "heavy", "light", "Coder")
    },
    @{
        File = (Join-Path $RepoRoot ".codex\agents\research.toml")
        Patterns = @("Research agent", "Read-only", "facts", "inferences")
    },
    @{
        File = (Join-Path $RepoRoot ".codex\agents\coder.toml")
        Patterns = @("Coder agent", "implementation", "Do not spawn sub-agents")
    },
    @{
        File = (Join-Path $RepoRoot ".codex\skills\architect-direction\SKILL.md")
        Patterns = @("heavy", "light", "Research", "Coder", "Heavy plan")
    },
    @{
        File = (Join-Path $RepoRoot ".codex\skills\light-direction\SKILL.md")
        Patterns = @("Short Plan", "light-work", "light-review", "plan")
    },
    @{
        File = (Join-Path $RepoRoot "docs\spec.md")
        Patterns = @("LMStudio", "Gemini", "xAI", "BatchAPI", "xTranslator")
    },
    @{
        File = (Join-Path $RepoRoot "docs\architecture.md")
        Patterns = @("Dependency Inversion Principle", "UI Port / UseCase Input", "DTO", "SQLite", "Rust")
    },
    @{
        File = (Join-Path $RepoRoot "docs\tech-selection.md")
        Patterns = @("Tauri 2", "Rust", "Svelte 5", "SQLite", "sqlx")
    },
    @{
        File = (Join-Path $RepoRoot "docs\core-beliefs.md")
        Patterns = @("agent-first", "AGENTS.md", "exec-plans", ".codex/")
    },
    @{
        File = (Join-Path $RepoRoot "docs\index.md")
        Patterns = @("AGENTS.md", ".codex", "exec-plans", "quality-score.md")
    },
    @{
        File = (Join-Path $RepoRoot "docs\executable-specs.md")
        Patterns = @("test", "acceptance checks", "fixture", "validation")
    },
    @{
        File = (Join-Path $RepoRoot "4humans\quality-score.md")
        Patterns = @("Codex workflow source of truth", "Agent role contracts", "Executable specs and constraints")
    }
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
