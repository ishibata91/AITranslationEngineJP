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

function Test-PackageScript {
    param(
        [pscustomobject]$Package,
        [string]$ScriptName
    )

    return $null -ne $Package.scripts.PSObject.Properties[$ScriptName]
}

$failures = 0
$ranAnything = $false
$sonarScanned = $false
$rootExecutionGateRan = $false
$rootPackage = $null
$sonarProjectProperties = Join-Path $RepoRoot "sonar-project.properties"
$rootPackageJsonPath = Join-Path $RepoRoot "package.json"
$rootCargoTomlPath = Join-Path $RepoRoot "src-tauri\Cargo.toml"

$packageJsons = Get-ChildItem -Path $RepoRoot -Recurse -File -Filter package.json |
    Where-Object { $_.FullName -notmatch '[\\/](node_modules|dist|build|coverage|target|\.git)[\\/]' }

$rootPackageJson = $packageJsons | Where-Object { $_.FullName -eq $rootPackageJsonPath } | Select-Object -First 1
if ($null -ne $rootPackageJson) {
    $rootPackage = Get-Content -LiteralPath $rootPackageJson.FullName -Raw | ConvertFrom-Json
    if (Test-PackageScript -Package $rootPackage -ScriptName "gate:execution") {
        $ranAnything = $true
        $rootExecutionGateRan = $true
        Invoke-Step -Command (Resolve-PackageManager -Directory $RepoRoot) -Arguments @("run", "gate:execution") -WorkingDirectory $RepoRoot -Failures ([ref]$failures)
        if (Test-PackageScript -Package $rootPackage -ScriptName "scan:sonar") {
            $sonarScanned = $true
        }
    }
}

$cargoTomls = Get-ChildItem -Path $RepoRoot -Recurse -File -Filter Cargo.toml |
    Where-Object {
        $_.FullName -notmatch '[\\/](target|node_modules|dist|build|coverage|\.git)[\\/]' -and
        (-not $rootExecutionGateRan -or $_.FullName -ne $rootCargoTomlPath)
    }

foreach ($cargoToml in $cargoTomls) {
    $ranAnything = $true
    $dir = Split-Path -Parent $cargoToml.FullName
    Invoke-Step -Command "cargo" -Arguments @("fmt", "--all", "--check") -WorkingDirectory $dir -Failures ([ref]$failures)
    Invoke-Step -Command "cargo" -Arguments @("clippy", "--all-targets", "--all-features", "--", "-D", "warnings") -WorkingDirectory $dir -Failures ([ref]$failures)
    Invoke-Step -Command "cargo" -Arguments @("test", "--all-features") -WorkingDirectory $dir -Failures ([ref]$failures)
}

foreach ($packageJson in $packageJsons) {
    if ($rootExecutionGateRan -and $packageJson.FullName -eq $rootPackageJsonPath) {
        continue
    }

    $ranAnything = $true
    $dir = Split-Path -Parent $packageJson.FullName
    $packageManager = Resolve-PackageManager -Directory $dir
    $package = Get-Content -LiteralPath $packageJson.FullName -Raw | ConvertFrom-Json
    $scripts = @()

    if (Test-PackageScript -Package $package -ScriptName "lint") { $scripts += "lint" }
    if (Test-PackageScript -Package $package -ScriptName "test") { $scripts += "test" }
    if (Test-PackageScript -Package $package -ScriptName "build") { $scripts += "build" }

    if ($scripts.Count -eq 0) {
        Write-Host "SKIP no lint/test/build scripts in $($packageJson.FullName)" -ForegroundColor Yellow
        continue
    }

    foreach ($script in $scripts) {
        Invoke-Step -Command $packageManager -Arguments @("run", $script) -WorkingDirectory $dir -Failures ([ref]$failures)

        if (
            -not $sonarScanned -and
            $script -eq "lint" -and
            (Test-Path $sonarProjectProperties) -and
            ((Resolve-Path $dir).Path -eq (Resolve-Path $RepoRoot).Path)
        ) {
            if (Test-PackageScript -Package $package -ScriptName "scan:sonar") {
                Invoke-Step -Command $packageManager -Arguments @("run", "scan:sonar") -WorkingDirectory $RepoRoot -Failures ([ref]$failures)
            } else {
                Invoke-Step -Command "sonar-scanner" -Arguments @() -WorkingDirectory $RepoRoot -Failures ([ref]$failures)
            }
            $sonarScanned = $true
        }
    }
}

$hasPackageLint = $packageJsons |
    Where-Object {
        $package = Get-Content -LiteralPath $_.FullName -Raw | ConvertFrom-Json
        Test-PackageScript -Package $package -ScriptName "lint"
    } |
    Select-Object -First 1

if ((Test-Path $sonarProjectProperties) -and -not $sonarScanned -and -not $hasPackageLint) {
    $ranAnything = $true
    if ($null -ne $rootPackageJson) {
        $rootPackage = if ($null -ne $rootPackage) { $rootPackage } else { Get-Content -LiteralPath $rootPackageJson.FullName -Raw | ConvertFrom-Json }
        if (Test-PackageScript -Package $rootPackage -ScriptName "scan:sonar") {
            Invoke-Step -Command (Resolve-PackageManager -Directory $RepoRoot) -Arguments @("run", "scan:sonar") -WorkingDirectory $RepoRoot -Failures ([ref]$failures)
            $sonarScanned = $true
        }
    }

    if (-not $sonarScanned) {
        Invoke-Step -Command "sonar-scanner" -Arguments @() -WorkingDirectory $RepoRoot -Failures ([ref]$failures)
        $sonarScanned = $true
    }
}

if (-not $ranAnything) {
    Write-Host "SKIP no Cargo.toml, package.json, or sonar-project.properties targets found. Execution harness is installed but has nothing to run yet." -ForegroundColor Yellow
}

if ($failures -gt 0) {
    Write-Error "Execution harness failed with $failures issue(s)."
}

Write-Host ""
Write-Host "Execution harness passed." -ForegroundColor Green
