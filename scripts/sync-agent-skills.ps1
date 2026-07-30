[CmdletBinding()]
param(
    [switch]$Check
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$sourceRoot = Join-Path $repoRoot ".claude\skills"
$destinationRoot = Join-Path $repoRoot ".agents\skills"

function ConvertTo-CodexText {
    param([Parameter(Mandatory)][string]$Content)

    return $Content.
        Replace(".claude/skills/", ".agents/skills/").
        Replace(".claude/agents/", ".codex/agents/").
        Replace(".claude/settings.json", ".codex/config.toml").
        Replace(".claude", ".codex").
        Replace("CLAUDE.md", "AGENTS.md").
        Replace("Claude", "Codex").
        Replace("claude", "codex")
}

if (-not (Test-Path -LiteralPath $sourceRoot -PathType Container)) {
    throw "Claude skills directory was not found: $sourceRoot"
}

if (-not (Test-Path -LiteralPath $destinationRoot -PathType Container)) {
    if ($Check) {
        throw "Codex skills directory was not found: $destinationRoot"
    }
    New-Item -ItemType Directory -Path $destinationRoot | Out-Null
}

$sourceFiles = Get-ChildItem -LiteralPath $sourceRoot -Recurse -File
$sourceRelativePaths = @{}
$differences = [System.Collections.Generic.List[string]]::new()

foreach ($sourceFile in $sourceFiles) {
    $relativePath = [System.IO.Path]::GetRelativePath($sourceRoot, $sourceFile.FullName)
    $sourceRelativePaths[$relativePath] = $true
    $destinationFile = Join-Path $destinationRoot $relativePath
    $generated = ConvertTo-CodexText (Get-Content -LiteralPath $sourceFile.FullName -Raw)

    $matches = (Test-Path -LiteralPath $destinationFile -PathType Leaf) -and
        ((Get-Content -LiteralPath $destinationFile -Raw) -ceq $generated)

    if ($matches) {
        continue
    }

    $differences.Add($relativePath)
    if (-not $Check) {
        $destinationDirectory = Split-Path -Parent $destinationFile
        New-Item -ItemType Directory -Path $destinationDirectory -Force | Out-Null
        Set-Content -LiteralPath $destinationFile -Value $generated -Encoding utf8NoBOM -NoNewline
    }
}

$staleFiles = Get-ChildItem -LiteralPath $destinationRoot -Recurse -File |
    Where-Object {
        $relativePath = [System.IO.Path]::GetRelativePath($destinationRoot, $_.FullName)
        -not $sourceRelativePaths.ContainsKey($relativePath)
    }

foreach ($staleFile in $staleFiles) {
    $relativePath = [System.IO.Path]::GetRelativePath($destinationRoot, $staleFile.FullName)
    $differences.Add($relativePath)
    if (-not $Check) {
        Remove-Item -LiteralPath $staleFile.FullName -Force
    }
}

if ($Check -and $differences.Count -gt 0) {
    Write-Error "Codex skills are out of sync:`n$($differences -join "`n")"
    exit 1
}

if (-not $Check) {
    Get-ChildItem -LiteralPath $destinationRoot -Recurse -Directory |
        Sort-Object FullName -Descending |
        Where-Object { -not (Get-ChildItem -LiteralPath $_.FullName -Force) } |
        Remove-Item -Force
}

Write-Output "Codex skills are synchronized with .claude/skills."
