$manifest = (Get-AppxPackage -AllUsers Microsoft.XboxGamingOverlay).InstallLocation + "\AppXManifest.xml"
if (Test-Path $manifest) { Add-AppxPackage -DisableDevelopmentMode -Register "$manifest" -ErrorAction SilentlyContinue }
foreach ($protocol in @("ms-gamebar", "ms-gamebarservices", "ms-gamingoverlay")) {
    $path = "HKCU:\Software\Classes\$protocol"
    if (Test-Path $path) { Remove-Item -Path $path -Recurse -Force -ErrorAction SilentlyContinue }
}
Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\GameDVR" -Name "AppCaptureEnabled" -Value 1 -Force
Set-ItemProperty -Path "HKCU:\System\GameConfigStore" -Name "GameDVR_Enabled" -Value 1 -Force