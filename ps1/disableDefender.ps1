$defPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender"
$rtPath  = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender\Real-Time Protection"

New-Item -Path $defPath -Force | Out-Null
New-Item -Path $rtPath -Force | Out-Null

Set-ItemProperty -Path $defPath -Name "DisableAntiSpyware" -Value 1 -Force
Set-ItemProperty -Path $defPath -Name "ServiceKeepAlive" -Value 0 -Force
Set-ItemProperty -Path $defPath -Name "AllowFastServiceStartup" -Value 0 -Force

Set-ItemProperty -Path $rtPath -Name "DisableRealtimeMonitoring" -Value 1 -Force
Set-ItemProperty -Path $rtPath -Name "DisableBehaviorMonitoring" -Value 1 -Force
Set-ItemProperty -Path $rtPath -Name "DisableOnAccessProtection" -Value 1 -Force
Set-ItemProperty -Path $rtPath -Name "DisableIOAVProtection" -Value 1 -Force
Set-ItemProperty -Path $rtPath -Name "DisableIntrusionPreventionSystem" -Value 1 -Force
Set-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableRealtimeMonitoring" -Value 0 -Force
Set-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableBehaviorMonitoring" -Value 0 -Force
Set-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableOnAccessProtection" -Value 0 -Force
Set-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableIOAVProtection" -Value 0 -Force
Set-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableIntrusionPreventionSystem" -Value 0 -Force

foreach ($s in @("WinDefend","WdBoot","WdFilter","WdNisDrv","WdNisSvc","Sense","MDCoreSvc","SgrmAgent","SgrmBroker","wscsvc","SecurityHealthService")) {
    Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\$s" -Name "Start" -Value 4 -Force
}
