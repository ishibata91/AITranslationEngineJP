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

function Assert-PathState {
    param(
        [string]$Path,
        [bool]$ShouldExist,
        [string]$Label,
        [ref]$Failures
    )

    $exists = Test-Path -LiteralPath $Path
    if ($exists -ne $ShouldExist) {
        $expected = if ($ShouldExist) { "exist" } else { "be absent" }
        Write-Host "FAIL semantic check '$Label': expected '$Path' to $expected" -ForegroundColor Red
        $Failures.Value++
    } else {
        $state = if ($ShouldExist) { "exists" } else { "is absent" }
        Write-Host "PASS semantic check '$Label': '$Path' $state" -ForegroundColor Green
    }
}

function Assert-CanonicalPhrase {
    param(
        [string]$FilePath,
        [string]$ExpectedPhrase,
        [string[]]$ForbiddenPhrases,
        [ref]$Failures
    )

    $content = Get-Content -LiteralPath $FilePath -Raw

    if ($content -notmatch [regex]::Escape($ExpectedPhrase)) {
        Write-Host "FAIL semantic phrase '$ExpectedPhrase' missing in $FilePath" -ForegroundColor Red
        $Failures.Value++
    } else {
        Write-Host "PASS semantic phrase '$ExpectedPhrase' found in $FilePath" -ForegroundColor Green
    }

    foreach ($forbiddenPhrase in $ForbiddenPhrases) {
        if ($content -match [regex]::Escape($forbiddenPhrase)) {
            Write-Host "FAIL semantic phrase '$forbiddenPhrase' should not appear in $FilePath" -ForegroundColor Red
            $Failures.Value++
        } else {
            Write-Host "PASS semantic phrase '$forbiddenPhrase' absent from $FilePath" -ForegroundColor Green
        }
    }
}

$checks = @(
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\README.md"); Patterns = @("directing-implementation", "designing-implementation", "directing-fixes", "architecting-tests", "UI", "Scenario", "Logic", "single-pass", "changes/", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\task_designer.toml"); Patterns = @("UI", "Scenario", "Logic", "task-local design", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\ctx_loader.toml"); Patterns = @("facts", "constraints", "gaps", "exec-plan") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\workplan_builder.toml"); Patterns = @("UI", "Scenario", "Logic", "Implementation Plan", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\test_architect.toml"); Patterns = @("test_architect", "spec-aligned tests", "fixtures", "validation commands") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\implementer.toml"); Patterns = @("allowed scope", "validation", "touched files") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\fault_tracer.toml"); Patterns = @("root-cause hypotheses", "observation plan", "temporary logging") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\log_instrumenter.toml"); Patterns = @("temporary log statements", "\[tracing-fixes\]", "Remove any temporary instrumentation") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\agents\review_cycler.toml"); Patterns = @("single-pass", "spec deviation", "exception handling", "resource cleanup", "missing tests", "pass", "reroute") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\directing-implementation\SKILL.md"); Patterns = @("docs/exec-plans/templates/impl-plan.md", "designing-implementation", "UI", "Scenario", "Logic", "architecting-tests", "reviewing-implementation", "reroute", "changes/", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\designing-implementation\SKILL.md"); Patterns = @("UI", "Scenario", "Logic", "task-local design", "changes/", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\directing-fixes\SKILL.md"); Patterns = @("docs/exec-plans/templates/fix-plan.md", "distilling-fixes", "tracing-fixes", "logging-fixes", "architecting-tests", "reviewing-fixes", "tasks.md") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\architecting-tests\SKILL.md"); Patterns = @("Acceptance Checks", "failing tests", "fixtures", "validation commands", "closeout notes", "updating-docs") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\reviewing-implementation\SKILL.md"); Patterns = @("Review Scope", "pass", "reroute", "score") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot ".codex\skills\reviewing-fixes\SKILL.md"); Patterns = @("Review Scope", "pass", "reroute") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\exec-plans\templates\impl-plan.md"); Patterns = @("Decision Basis", "UI", "Scenario", "Logic", "Implementation Plan", "Required Evidence", "4humans Sync", "Outcome") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\exec-plans\templates\fix-plan.md"); Patterns = @("Decision Basis", "Known Facts", "Trace Plan", "Fix Plan", "Required Evidence", "4humans Sync", "Outcome") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\spec.md"); Patterns = @("LMStudio", "Gemini", "xAI", "BatchAPI", "xTranslator") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\architecture.md"); Patterns = @("Dependency Inversion Principle", "UI Port / UseCase Input", "DTO", "SQLite", "Rust") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\tech-selection.md"); Patterns = @("Tauri 2", "Rust", "Svelte 5", "SQLite", "sqlx") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\core-beliefs.md"); Patterns = @("agent-first", "directing-implementation", "directing-fixes", "single-pass", "evidence") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "docs\index.md"); Patterns = @("AGENTS.md", ".codex", "directing-implementation", "directing-fixes", "designing-implementation", "quality-score.md", "tests / acceptance checks / validation commands") },
    [pscustomobject]@{ File = (Join-Path $RepoRoot "4humans\quality-score.md"); Patterns = @("Codex workflow source of truth", "Design harness", "directing-implementation", "directing-fixes") }
)

$failures = 0

foreach ($check in $checks) {
    Assert-Patterns -FilePath $check.File -Patterns $check.Patterns -Failures ([ref]$failures)
}

$semanticPathChecks = @(
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src\application"); ShouldExist = $true; Label = "frontend application layer root" },
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src\gateway"); ShouldExist = $true; Label = "frontend gateway root" },
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src\shared"); ShouldExist = $true; Label = "frontend shared DTO root" },
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src\ui"); ShouldExist = $true; Label = "frontend ui root" },
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src\domain"); ShouldExist = $false; Label = "frontend domain forbidden during bootstrap" },
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src\infra"); ShouldExist = $false; Label = "frontend infra forbidden during bootstrap" },
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src-tauri\src\application"); ShouldExist = $true; Label = "backend application layer root" },
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src-tauri\src\domain"); ShouldExist = $true; Label = "backend domain layer root" },
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src-tauri\src\infra"); ShouldExist = $true; Label = "backend infra layer root" },
    [pscustomobject]@{ Path = (Join-Path $RepoRoot "src-tauri\src\gateway"); ShouldExist = $true; Label = "backend gateway root" }
)

foreach ($check in $semanticPathChecks) {
    Assert-PathState -Path $check.Path -ShouldExist $check.ShouldExist -Label $check.Label -Failures ([ref]$failures)
}

$semanticPhraseChecks = @(
    [pscustomobject]@{
        File = (Join-Path $RepoRoot "docs\core-beliefs.md")
        ExpectedPhrase = "tests / acceptance checks / validation commands"
        ForbiddenPhrases = @("test / acceptance checks / validation commands")
    },
    [pscustomobject]@{
        File = (Join-Path $RepoRoot "docs\architecture.md")
        ExpectedPhrase = "tests / acceptance checks / validation commands"
        ForbiddenPhrases = @()
    },
    [pscustomobject]@{
        File = (Join-Path $RepoRoot "docs\index.md")
        ExpectedPhrase = "tests / acceptance checks / validation commands"
        ForbiddenPhrases = @()
    }
)

foreach ($check in $semanticPhraseChecks) {
    Assert-CanonicalPhrase -FilePath $check.File -ExpectedPhrase $check.ExpectedPhrase -ForbiddenPhrases $check.ForbiddenPhrases -Failures ([ref]$failures)
}

if ($failures -gt 0) {
    Write-Error "Design harness failed with $failures issue(s)."
}

Write-Host ""
Write-Host "Design harness passed." -ForegroundColor Green

