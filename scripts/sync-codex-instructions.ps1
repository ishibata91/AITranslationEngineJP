[CmdletBinding()]
param(
    [switch]$Check
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$sourceFile = Join-Path $repoRoot "CLAUDE.md"
$destinationFile = Join-Path $repoRoot "AGENTS.md"

if (-not (Test-Path -LiteralPath $sourceFile -PathType Leaf)) {
    throw "Claude instruction file was not found: $sourceFile"
}

$generated = (Get-Content -LiteralPath $sourceFile -Raw).
    Replace(".claude/skills/", ".agents/skills/").
    Replace(".claude/agents/", ".codex/agents/").
    Replace(".claude/settings.json", ".codex/config.toml").
    Replace(".claude", ".codex").
    Replace("CLAUDE.md", "AGENTS.md").
    Replace("Claude", "Codex").
    Replace("claude", "codex")

$matches = (Test-Path -LiteralPath $destinationFile -PathType Leaf) -and
    ((Get-Content -LiteralPath $destinationFile -Raw) -ceq $generated)

if ($Check -and -not $matches) {
    Write-Error "AGENTS.md is out of sync with CLAUDE.md."
    exit 1
}

if (-not $Check -and -not $matches) {
    Set-Content -LiteralPath $destinationFile -Value $generated -Encoding utf8NoBOM -NoNewline
}

Write-Output "AGENTS.md is synchronized with CLAUDE.md."
