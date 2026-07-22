$defPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender"
$rtPath  = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender\Real-Time Protection"

@("DisableAntiSpyware","ServiceKeepAlive","AllowFastServiceStartup") | ForEach-Object {
    Remove-ItemProperty -Path $defPath -Name $_ -ErrorAction SilentlyContinue
}

@("DisableRealtimeMonitoring","DisableBehaviorMonitoring","DisableOnAccessProtection","DisableIOAVProtection","DisableIntrusionPreventionSystem","LocalSettingOverrideDisableRealtimeMonitoring","LocalSettingOverrideDisableBehaviorMonitoring","LocalSettingOverrideDisableOnAccessProtection","LocalSettingOverrideDisableIOAVProtection","LocalSettingOverrideDisableIntrusionPreventionSystem") | ForEach-Object {
    Remove-ItemProperty -Path $rtPath -Name $_ -ErrorAction SilentlyContinue
}

foreach ($s in @("wscsvc","SecurityHealthService","WinDefend","WdBoot","WdFilter","WdNisDrv","WdNisSvc","Sense","MDCoreSvc","SgrmAgent","SgrmBroker")) {
    Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\$s" -Name "Start" -Value 2 -Force
}
