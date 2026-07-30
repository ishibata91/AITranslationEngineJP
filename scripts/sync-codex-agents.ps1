[CmdletBinding()]
param(
    [switch]$Check
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$sourceRoot = Join-Path $repoRoot ".claude\agents"
$destinationRoot = Join-Path $repoRoot ".codex\agents"

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
    throw "Claude agents directory was not found: $sourceRoot"
}

if (-not (Test-Path -LiteralPath $destinationRoot -PathType Container)) {
    if ($Check) {
        throw "Codex agents directory was not found: $destinationRoot"
    }
    New-Item -ItemType Directory -Path $destinationRoot -Force | Out-Null
}

$expectedFiles = @{}
$differences = [System.Collections.Generic.List[string]]::new()

foreach ($sourceFile in Get-ChildItem -LiteralPath $sourceRoot -Filter "*.md" -File) {
    $content = Get-Content -LiteralPath $sourceFile.FullName -Raw
    if ($content -notmatch '(?s)\A---\r?\n(?<frontmatter>.*?)\r?\n---\r?\n(?<body>.*)\z') {
        throw "Invalid Claude agent frontmatter: $($sourceFile.FullName)"
    }

    $frontmatter = $Matches.frontmatter
    $body = ConvertTo-CodexText $Matches.body
    $name = ([regex]::Match($frontmatter, '(?m)^name:\s*(.+)$')).Groups[1].Value.Trim()
    $description = ConvertTo-CodexText (
        ([regex]::Match($frontmatter, '(?m)^description:\s*(.+)$')).Groups[1].Value.Trim().Trim('"')
    )

    if (-not $name -or -not $description) {
        throw "Claude agent must define name and description: $($sourceFile.FullName)"
    }

    $sandboxMode = if ($name -eq "conflict_resolver") { "workspace-write" } else { "read-only" }
    $generated = @(
        "name = $($name | ConvertTo-Json -Compress)"
        "description = $($description | ConvertTo-Json -Compress)"
        "sandbox_mode = $($sandboxMode | ConvertTo-Json -Compress)"
        "developer_instructions = $($body.TrimEnd() | ConvertTo-Json -Compress)"
        ""
    ) -join "`n"

    $destinationFile = Join-Path $destinationRoot "$name.toml"
    $expectedFiles[$destinationFile] = $true
    $matches = (Test-Path -LiteralPath $destinationFile -PathType Leaf) -and
        ((Get-Content -LiteralPath $destinationFile -Raw) -eq $generated)

    if ($matches) {
        continue
    }

    $differences.Add("$name.toml")
    if (-not $Check) {
        Set-Content -LiteralPath $destinationFile -Value $generated -Encoding utf8NoBOM -NoNewline
    }
}

foreach ($staleFile in Get-ChildItem -LiteralPath $destinationRoot -Filter "*.toml" -File) {
    if ($expectedFiles.ContainsKey($staleFile.FullName)) {
        continue
    }

    $differences.Add($staleFile.Name)
    if (-not $Check) {
        Remove-Item -LiteralPath $staleFile.FullName -Force
    }
}

if ($Check -and $differences.Count -gt 0) {
    Write-Error "Codex agents are out of sync:`n$($differences -join "`n")"
    exit 1
}

Write-Output "Codex agents are synchronized with .claude/agents."
