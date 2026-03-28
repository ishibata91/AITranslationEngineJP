param(
    [string]$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Assert-PathExists {
    param(
        [string]$Path,
        [ref]$Failures
    )

    if (-not (Test-Path -LiteralPath $Path)) {
        Write-Host "FAIL missing path: $Path" -ForegroundColor Red
        $Failures.Value++
    } else {
        Write-Host "PASS path exists: $Path" -ForegroundColor Green
    }
}

function Remove-CodeBlocks {
    param([string]$Text)
    return [regex]::Replace($Text, '(?ms)```.*?```', '')
}

function Resolve-LinkTarget {
    param(
        [string]$SourceFile,
        [string]$Target,
        [string]$RepoRoot
    )

    $cleanTarget = $Target.Split('#')[0]
    if ([string]::IsNullOrWhiteSpace($cleanTarget)) {
        return $null
    }

    if ($cleanTarget -match '^[a-zA-Z][a-zA-Z0-9+.-]*:' -and $cleanTarget -notmatch '^[a-zA-Z]:[\\/]') {
        return $null
    }

    if ($cleanTarget.StartsWith('/')) {
        return $null
    }

    if ($cleanTarget -match '^[a-zA-Z]:[\\/]') {
        return $cleanTarget
    }

    $baseDir = Split-Path -Parent $SourceFile
    return [System.IO.Path]::GetFullPath((Join-Path $baseDir $cleanTarget))
}

$requiredPaths = @(
    (Join-Path $RepoRoot "AGENTS.md"),
    (Join-Path $RepoRoot ".codex\README.md"),
    (Join-Path $RepoRoot ".codex\agents\architect.toml"),
    (Join-Path $RepoRoot ".codex\agents\research.toml"),
    (Join-Path $RepoRoot ".codex\agents\coder.toml"),
    (Join-Path $RepoRoot ".codex\skills\architect-direction\SKILL.md"),
    (Join-Path $RepoRoot ".codex\skills\architect-direction\agents\openai.yaml"),
    (Join-Path $RepoRoot ".codex\skills\light-direction\SKILL.md"),
    (Join-Path $RepoRoot ".codex\skills\light-direction\agents\openai.yaml"),
    (Join-Path $RepoRoot ".codex\skills\light-work\SKILL.md"),
    (Join-Path $RepoRoot ".codex\skills\light-work\agents\openai.yaml"),
    (Join-Path $RepoRoot ".codex\skills\light-review\SKILL.md"),
    (Join-Path $RepoRoot ".codex\skills\light-review\agents\openai.yaml"),
    (Join-Path $RepoRoot "docs\index.md"),
    (Join-Path $RepoRoot "docs\core-beliefs.md"),
    (Join-Path $RepoRoot "docs\spec.md"),
    (Join-Path $RepoRoot "docs\architecture.md"),
    (Join-Path $RepoRoot "docs\tech-selection.md"),
    (Join-Path $RepoRoot "docs\er-draft.md"),
    (Join-Path $RepoRoot "docs\tech-debt-tracker.md"),
    (Join-Path $RepoRoot "docs\quality-score.md"),
    (Join-Path $RepoRoot "docs\references\index.md"),
    (Join-Path $RepoRoot "docs\executable-specs.md"),
    (Join-Path $RepoRoot "docs\exec-plans\active\README.md"),
    (Join-Path $RepoRoot "docs\exec-plans\completed\README.md"),
    (Join-Path $RepoRoot "docs\exec-plans\templates\heavy-plan.md"),
    (Join-Path $RepoRoot "docs\exec-plans\templates\light-plan.md"),
    (Join-Path $RepoRoot "scripts\harness\run.ps1"),
    (Join-Path $RepoRoot "scripts\harness\check-structure.ps1"),
    (Join-Path $RepoRoot "scripts\harness\check-design.ps1"),
    (Join-Path $RepoRoot "scripts\harness\check-execution.ps1")
)

$failures = 0

Write-Host "== Required paths ==" -ForegroundColor Cyan
foreach ($path in $requiredPaths) {
    Assert-PathExists -Path $path -Failures ([ref]$failures)
}

Write-Host ""
Write-Host "== Markdown links ==" -ForegroundColor Cyan

$markdownFiles = Get-ChildItem -Path $RepoRoot -Recurse -File -Filter *.md |
    Where-Object { $_.FullName -notmatch '[\\/]\.git[\\/]' }

$linkPattern = '!?'
$linkPattern += '\[[^\]]*\]\((?<target>[^)]+)\)'

foreach ($file in $markdownFiles) {
    $content = Remove-CodeBlocks -Text (Get-Content -LiteralPath $file.FullName -Raw)
    $matches = [regex]::Matches($content, $linkPattern)

    foreach ($match in $matches) {
        $target = $match.Groups["target"].Value.Trim()
        if ($target.StartsWith("#")) {
            continue
        }

        $resolved = Resolve-LinkTarget -SourceFile $file.FullName -Target $target -RepoRoot $RepoRoot
        if ($null -eq $resolved) {
            continue
        }

        if (-not (Test-Path -LiteralPath $resolved)) {
            Write-Host "FAIL broken link: $($file.FullName) -> $target" -ForegroundColor Red
            $failures++
        }
    }
}

if ($failures -gt 0) {
    Write-Error "Structure harness failed with $failures issue(s)."
}

Write-Host ""
Write-Host "Structure harness passed." -ForegroundColor Green
