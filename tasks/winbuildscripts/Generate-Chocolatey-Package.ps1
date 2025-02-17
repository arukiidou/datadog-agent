Param(
    [Parameter(Mandatory=$true,Position=0)]
    [ValidateSet("offline", "online")]
    [String]
    $installMethod
)

$ErrorActionPreference = 'Stop';
Set-Location c:\mnt

# Install chocolatey binary
$env:chocolateyUseWindowsCompression = 'true'; Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

# Install dev tools, including invoke
pip3 install -r requirements.txt

$outputDirectory = "c:\mnt\build-out"
$rawAgentVersion = (inv agent.version)
$copyright = "Datadog {0}" -f (Get-Date).Year

$releasePattern = "(\d+\.\d+\.\d+)"
$releaseCandidatePattern = "(\d+\.\d+\.\d+)-rc\.(\d+)"
$develPattern = "(\d+\.\d+\.\d+)-devel\+git\.\d+\.(.+)"

$nuspecFile = "c:\mnt\chocolatey\datadog-agent-online.nuspec"
$licensePath = "c:\mnt\chocolatey\tools-online\LICENSE.txt"
$installScript = "c:\mnt\chocolatey\tools-online\chocolateyinstall.ps1"

if ($installMethod -eq "offline") {
    $nuspecFile = "c:\mnt\chocolatey\datadog-agent-offline.nuspec"
    $licensePath = "c:\mnt\chocolatey\tools-offline\LICENSE.txt"
}

if ($rawAgentVersion -match $releaseCandidatePattern) {
    $agentVersionMatches = $rawAgentVersion | Select-String -Pattern $releaseCandidatePattern
    $agentVersion = "{0}-rc-{1}" -f $agentVersionMatches.Matches.Groups[1], $agentVersionMatches.Matches.Groups[2].Value
    # We don't have release notes for RCs but this way the user can always see what commits are included in this RC
    $releaseNotes = "https://github.com/DataDog/datadog-agent/releases/tag/{0}-rc.{1}" -f $agentVersionMatches.Matches.Groups[1], $agentVersionMatches.Matches.Groups[2]
    $url = "https://s3.amazonaws.com/dd-agent-mstesting/builds/beta/ddagent-cli-$($agentVersionMatches.Matches.Groups[1])-rc.$($agentVersionMatches.Matches.Groups[2]).msi"
} elseif ($rawAgentVersion -match $develPattern) {
    if ($installMethod -eq "online") {
        # We don't publish online chocolatey packages for dev branches, error out
        Write-Host "Chocolatey packages are not built for dev branches aborting"
        exit 2
    }
    $agentVersionMatches = $rawAgentVersion | Select-String -Pattern $develPattern
    $agentVersion = "{0}-devel-{1}" -f $agentVersionMatches.Matches.Groups[1], $agentVersionMatches.Matches.Groups[2].Value
    # We don't have release notes for devel, so point it to the generic url
    $releaseNotes = "https://github.com/DataDog/datadog-agent/releases"
} elseif ($rawAgentVersion -match $releasePattern) {
    $agentVersionMatches = $rawAgentVersion | Select-String -Pattern $releasePattern
    $agentVersion = $agentVersionMatches.Matches.Groups[1].Value
    $releaseNotes = "https://github.com/DataDog/datadog-agent/releases/tag/$agentVersion"
    $url = "https://s3.amazonaws.com/ddagent-windows-stable/ddagent-cli-$($agentVersion).msi"
} else {
    Write-Host "Unknown agent version '$rawAgentVersion', aborting"
    exit 3
}

Invoke-WebRequest -Uri "https://raw.githubusercontent.com/DataDog/datadog-agent/main/LICENSE" -OutFile $licensePath

Write-Host "Generating Chocolatey $installMethod package version $agentVersion in $outputDirectory"

if (!(Test-Path $outputDirectory)) {
    New-Item -ItemType Directory -Path $outputDirectory
}

if ($installMethod -eq "online") {
    # Set the $url in the install script
    (Get-Content $installScript).replace('$__url_from_ci__', '"' +  $url  + '"') | Set-Content $installScript
}

choco pack --out=$outputDirectory $nuspecFile package_version=$agentVersion release_notes=$releaseNotes copyright=$copyright
