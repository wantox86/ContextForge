# build-windows.ps1 — Build ContextForge for Windows (amd64)
# Run in PowerShell on Windows: .\scripts\build-windows.ps1
# Optional: -Nsis to also generate an NSIS installer
param(
    [switch]$Nsis,
    [switch]$Clean
)

$ErrorActionPreference = "Stop"

$ProjectDir = Split-Path -Parent $PSScriptRoot
Set-Location $ProjectDir

Write-Host "▶ Building ContextForge for Windows (amd64)..." -ForegroundColor Cyan

$args = @("-platform", "windows/amd64", "-o", "ContextForge")
if ($Clean) { $args += "-clean" }
if ($Nsis)  { $args += "-nsis" }

& wails build @args

Write-Host ""
Write-Host "✅ Windows build ready in: $ProjectDir\build\bin" -ForegroundColor Green
Get-ChildItem "$ProjectDir\build\bin"
