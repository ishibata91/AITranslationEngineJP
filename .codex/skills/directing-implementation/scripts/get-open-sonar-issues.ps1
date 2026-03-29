param(
    [Parameter(Mandatory = $true)]
    [string]$Project,

    [string[]]$OwnedPaths = @()
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Normalize-RepoPath {
    param(
        [AllowNull()]
        [string]$Path
    )

    if ([string]::IsNullOrWhiteSpace($Path)) {
        return $null
    }

    return ($Path -replace "\\", "/").TrimStart(".").TrimStart("/")
}

function Get-IssueRepoPath {
    param(
        $Issue,
        [hashtable]$ComponentPathByKey
    )

    if ($Issue.PSObject.Properties.Name -contains "component") {
        $component = [string]$Issue.component
        if ($component.Contains(":")) {
            return Normalize-RepoPath -Path ($component.Split(":", 2)[1])
        }

        if ($ComponentPathByKey.ContainsKey($component)) {
            return Normalize-RepoPath -Path $ComponentPathByKey[$component]
        }
    }

    return $null
}

$rawJson = & sonar list issues --project $Project --format json
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}

$payload = $rawJson | ConvertFrom-Json
$componentPathByKey = @{}

if ($payload.PSObject.Properties.Name -contains "components") {
    foreach ($component in @($payload.components)) {
        if (
            ($component.PSObject.Properties.Name -contains "key") -and
            ($component.PSObject.Properties.Name -contains "path") -and
            -not [string]::IsNullOrWhiteSpace([string]$component.path)
        ) {
            $componentPathByKey[[string]$component.key] = [string]$component.path
        }
    }
}

$normalizedOwnedPaths = @($OwnedPaths | ForEach-Object { Normalize-RepoPath -Path $_ } | Where-Object { $_ })
$openIssues = foreach ($issue in @($payload.issues)) {
    if ([string]$issue.status -ne "OPEN") {
        continue
    }

    $repoPath = Get-IssueRepoPath -Issue $issue -ComponentPathByKey $componentPathByKey
    $matchesOwnedScope = $true

    if ($normalizedOwnedPaths.Count -gt 0) {
        $matchesOwnedScope = $false
        foreach ($ownedPath in $normalizedOwnedPaths) {
            if ($repoPath -eq $ownedPath -or ($repoPath -and $repoPath.StartsWith("$ownedPath/"))) {
                $matchesOwnedScope = $true
                break
            }
        }
    }

    if (-not $matchesOwnedScope) {
        continue
    }

    $item = [ordered]@{
        key = $issue.key
        rule = $issue.rule
        severity = $issue.severity
        status = $issue.status
        issueStatus = $issue.issueStatus
        message = $issue.message
        component = $issue.component
        repoPath = $repoPath
        textRange = $issue.textRange
    }

    [pscustomobject]$item
}

$result = [ordered]@{
    project = $Project
    totalIssues = @($payload.issues).Count
    openIssueCount = @($openIssues).Count
    openIssues = @($openIssues)
}

$result | ConvertTo-Json -Depth 100
