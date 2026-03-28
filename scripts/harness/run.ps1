param(
    [ValidateSet("all", "structure", "design", "execution")]
    [string]$Suite = "all",
    [string]$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Invoke-Harness {
    param([string]$ScriptPath)

    Write-Host ""
    Write-Host "== Running $([System.IO.Path]::GetFileName($ScriptPath)) ==" -ForegroundColor Cyan
    & powershell -ExecutionPolicy Bypass -File $ScriptPath -RepoRoot $RepoRoot
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }
}

$suiteOrder = switch ($Suite) {
    "structure" { @("check-structure.ps1") }
    "design" { @("check-design.ps1") }
    "execution" { @("check-execution.ps1") }
    default { @("check-structure.ps1", "check-design.ps1", "check-execution.ps1") }
}

foreach ($scriptName in $suiteOrder) {
    Invoke-Harness -ScriptPath (Join-Path $PSScriptRoot $scriptName)
}

Write-Host ""
Write-Host "All requested harness suites passed." -ForegroundColor Green
