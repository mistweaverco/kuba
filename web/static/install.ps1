#Requires -Version 5.1

<#
.SYNOPSIS
    Kuba Installation Script for Windows

.DESCRIPTION
    Automatically downloads and installs the latest version of Kuba from GitHub releases.
    Detects the user's Windows architecture and installs to the appropriate location.

.PARAMETER SystemWide
    Install to system-wide location (requires Administrator privileges)

.EXAMPLE
    .\install.ps1

.EXAMPLE
    .\install.ps1 -SystemWide

.NOTES
    Requires PowerShell 5.1 or later
    Requires Internet connection to download from GitHub
#>

param(
    [switch]$SystemWide
)

# Set error action preference
$ErrorActionPreference = "Stop"

# GitHub repository details
$Repo = "mistweaverco/kuba"
$BinaryName = "kuba.exe"

# Colors for output
$Colors = @{
    Red    = "Red"
    Green  = "Green"
    Yellow = "Yellow"
    Blue   = "Blue"
    White  = "White"
}

# Function to print colored output
function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Colors.Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Colors.Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Colors.Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Colors.Red
}

# Function to detect Windows architecture
function Get-WindowsArchitecture {
    if ([Environment]::Is64BitOperatingSystem) {
        return "amd64"
    } else {
        return "386"
    }
}

# Function to get latest release version
function Get-LatestVersion {
    try {
        # GitHub automatically redirects /latest to the latest release
        # We can extract the version from the redirect URL
        $response = Invoke-WebRequest -Uri "https://github.com/$Repo/releases/latest" -MaximumRedirection 0 -ErrorAction SilentlyContinue

        if ($response.StatusCode -eq 302 -or $response.StatusCode -eq 301) {
            $redirectUrl = $response.Headers.Location
            $version = $redirectUrl -replace ".*/releases/tag/", ""
            return $version
        } else {
            throw "Unexpected response from GitHub"
        }
    }
    catch {
        Write-Error "Failed to get latest version from GitHub: $($_.Exception.Message)"
        exit 1
    }
}

# Function to download binary
function Download-Binary {
    param(
        [string]$Version,
        [string]$Architecture
    )

    $downloadUrl = "https://github.com/$Repo/releases/download/$Version/kuba-windows-$Architecture.exe"

    Write-Status "Downloading $BinaryName $Version for windows-$Architecture..."

    # Create temporary directory
    $tempDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }
    $binaryPath = Join-Path $tempDir.FullName $BinaryName

    try {
        # Download binary
        Invoke-WebRequest -Uri $downloadUrl -OutFile $binaryPath
        Write-Success "Download completed successfully"
        return $binaryPath
    }
    catch {
        Write-Error "Failed to download binary from $downloadUrl"
        Remove-Item $tempDir.FullName -Recurse -Force -ErrorAction SilentlyContinue
        exit 1
    }
}

# Function to determine install location
function Get-InstallLocation {
    if ($SystemWide) {
        # Check if running as Administrator
        if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
            Write-Error "System-wide installation requires Administrator privileges. Please run PowerShell as Administrator or remove the -SystemWide parameter."
            exit 1
        }
        return "C:\Program Files\kuba\$BinaryName"
    } else {
        # Install to user's local directory
        $userBin = Join-Path $env:USERPROFILE "AppData\Local\Microsoft\WinGet\Packages"
        $kubaDir = Join-Path $userBin "kuba"
        if (-not (Test-Path $kubaDir)) {
            New-Item -ItemType Directory -Path $kubaDir -Force | Out-Null
        }
        return Join-Path $kubaDir $BinaryName
    }
}

# Function to add to PATH
function Add-ToPath {
    param([string]$InstallPath)

    if ($SystemWide) {
        $installDir = Split-Path $InstallPath -Parent
        $currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")

        if ($currentPath -notlike "*$installDir*") {
            $newPath = "$currentPath;$installDir"
            [Environment]::SetEnvironmentVariable("PATH", $newPath, "Machine")
            Write-Status "Added $installDir to system PATH"
            Write-Warning "You may need to restart your terminal or computer for PATH changes to take effect"
        }
    } else {
        $installDir = Split-Path $InstallPath -Parent
        $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")

        if ($currentPath -notlike "*$installDir*") {
            $newPath = "$currentPath;$installDir"
            [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
            Write-Status "Added $installDir to user PATH"
            Write-Warning "You may need to restart your terminal for PATH changes to take effect"
        }
    }
}

# Function to install binary
function Install-Binary {
    param(
        [string]$SourcePath,
        [string]$InstallPath
    )

    Write-Status "Installing $BinaryName to $InstallPath..."

    # Create backup if binary already exists
    if (Test-Path $InstallPath) {
        $backupPath = "$InstallPath.backup.$(Get-Date -Format 'yyyyMMdd_HHmmss')"
        Write-Warning "Backing up existing binary to $backupPath"
        Copy-Item $InstallPath $backupPath
    }

    # Create directory if it doesn't exist
    $installDir = Split-Path $InstallPath -Parent
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }

    # Install binary
    try {
        Copy-Item $SourcePath $InstallPath
        Write-Success "$BinaryName installed successfully to $InstallPath"

        # Add to PATH
        Add-ToPath -InstallPath $InstallPath

    } catch {
        Write-Error "Failed to install $BinaryName: $($_.Exception.Message)"
        exit 1
    }
}

# Function to verify installation
function Test-Installation {
    param([string]$InstallPath)

    if (Test-Path $InstallPath) {
        Write-Success "Installation verified successfully!"
        Write-Status "You can now run: $BinaryName --version"
    } else {
        Write-Error "Installation verification failed"
        exit 1
    }
}

# Function to check PowerShell version
function Test-PowerShellVersion {
    $psVersion = $PSVersionTable.PSVersion
    if ($psVersion.Major -lt 5 -or ($psVersion.Major -eq 5 -and $psVersion.Minor -lt 1)) {
        Write-Error "PowerShell 5.1 or later is required. Current version: $psVersion"
        Write-Error "Please update PowerShell or use the manual installation method."
        exit 1
    }
}

# Function to check internet connectivity
function Test-InternetConnection {
    try {
        $response = Invoke-WebRequest -Uri "https://www.google.com" -TimeoutSec 10 -ErrorAction Stop
        return $true
    } catch {
        Write-Error "No internet connection detected. Please check your network connection."
        exit 1
    }
}

# Main installation function
function Main {
    Write-Status "Installing $BinaryName..."

    # Check PowerShell version
    Test-PowerShellVersion

    # Check internet connection
    Test-InternetConnection

    # Detect architecture
    $architecture = Get-WindowsArchitecture
    Write-Status "Detected architecture: $architecture"

    # Get latest version
    $version = Get-LatestVersion
    Write-Status "Latest version: $version"

    # Download binary
    $tempBinary = Download-Binary -Version $version -Architecture $architecture

    # Determine install location
    $installPath = Get-InstallLocation

    # Install binary
    Install-Binary -SourcePath $tempBinary -InstallPath $installPath

    # Clean up temporary files
    Remove-Item $tempBinary -Force -ErrorAction SilentlyContinue
    Remove-Item (Split-Path $tempBinary -Parent) -Recurse -Force -ErrorAction SilentlyContinue

    # Verify installation
    Test-Installation -InstallPath $installPath

    Write-Success "$BinaryName installation completed successfully!"

    if (-not $SystemWide) {
        Write-Status "To use $BinaryName from any location, restart your terminal or run:"
        Write-Host "refreshenv" -ForegroundColor $Colors.Yellow
    }
}

# Check if running on Windows
if ($env:OS -ne "Windows_NT") {
    Write-Error "This script is designed for Windows systems only."
    Write-Error "For Unix-like systems (Linux/macOS), please use the bash installation script:"
    Write-Error "curl -sSL https://kuba.mwco.app/install.sh | sh"
    exit 1
}

# Run main function
try {
    Main
} catch {
    Write-Error "Installation failed: $($_.Exception.Message)"
    exit 1
}
