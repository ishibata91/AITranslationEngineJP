param(
    [string]$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Invoke-Step {
    param(
        [string]$Command,
        [string[]]$Arguments,
        [string]$WorkingDirectory,
        [ref]$Failures
    )

    if (-not (Get-Command $Command -ErrorAction SilentlyContinue)) {
        Write-Host "FAIL missing command: $Command" -ForegroundColor Red
        $Failures.Value++
        return
    }

    Write-Host "RUN $Command $($Arguments -join ' ')" -ForegroundColor Cyan
    Push-Location $WorkingDirectory
    try {
        & $Command @Arguments
    } finally {
        Pop-Location
    }
    if ($LASTEXITCODE -ne 0) {
        Write-Host "FAIL $Command $($Arguments -join ' ')" -ForegroundColor Red
        $Failures.Value++
    } else {
        Write-Host "PASS $Command $($Arguments -join ' ')" -ForegroundColor Green
    }
}

function Resolve-PackageManager {
    param([string]$Directory)

    if (Test-Path (Join-Path $Directory "pnpm-lock.yaml")) { return "pnpm" }
    if (Test-Path (Join-Path $Directory "package-lock.json")) { return "npm" }
    if (Test-Path (Join-Path $Directory "yarn.lock")) { return "yarn" }
    return "npm"
}

$failures = 0
$ranAnything = $false

$cargoTomls = Get-ChildItem -Path $RepoRoot -Recurse -File -Filter Cargo.toml |
    Where-Object { $_.FullName -notmatch '[\\/](target|node_modules|dist|build|coverage|\.git)[\\/]' }

foreach ($cargoToml in $cargoTomls) {
    $ranAnything = $true
    $dir = Split-Path -Parent $cargoToml.FullName
    Invoke-Step -Command "cargo" -Arguments @("fmt", "--all", "--check") -WorkingDirectory $dir -Failures ([ref]$failures)
    Invoke-Step -Command "cargo" -Arguments @("clippy", "--all-targets", "--all-features", "--", "-D", "warnings") -WorkingDirectory $dir -Failures ([ref]$failures)
    Invoke-Step -Command "cargo" -Arguments @("test", "--all-features") -WorkingDirectory $dir -Failures ([ref]$failures)
}

$packageJsons = Get-ChildItem -Path $RepoRoot -Recurse -File -Filter package.json |
    Where-Object { $_.FullName -notmatch '[\\/](node_modules|dist|build|coverage|target|\.git)[\\/]' }

foreach ($packageJson in $packageJsons) {
    $ranAnything = $true
    $dir = Split-Path -Parent $packageJson.FullName
    $packageManager = Resolve-PackageManager -Directory $dir
    $package = Get-Content -LiteralPath $packageJson.FullName -Raw | ConvertFrom-Json
    $scripts = @()

    if ($null -ne $package.scripts.lint) { $scripts += "lint" }
    if ($null -ne $package.scripts.test) { $scripts += "test" }
    if ($null -ne $package.scripts.build) { $scripts += "build" }

    if ($scripts.Count -eq 0) {
        Write-Host "SKIP no lint/test/build scripts in $($packageJson.FullName)" -ForegroundColor Yellow
        continue
    }

    foreach ($script in $scripts) {
        Invoke-Step -Command $packageManager -Arguments @("run", $script) -WorkingDirectory $dir -Failures ([ref]$failures)
    }
}

if (-not $ranAnything) {
    Write-Host "SKIP no Cargo.toml or package.json targets found. Execution harness is installed but has nothing to run yet." -ForegroundColor Yellow
}

if ($failures -gt 0) {
    Write-Error "Execution harness failed with $failures issue(s)."
}

Write-Host ""
Write-Host "Execution harness passed." -ForegroundColor Green
