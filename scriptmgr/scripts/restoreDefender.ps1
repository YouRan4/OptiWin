# ==========================================
# 1. 清理/恢复策略路径与实时防护注册表
# ==========================================
$defPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender"
$rtPath  = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender\Real-Time Protection"
$spPath  = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender\SpyNet"

# 删除手动添加的策略键值
Remove-ItemProperty -Path $defPath -Name "DisableAntiSpyware" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $defPath -Name "ServiceKeepAlive" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $defPath -Name "AllowFastServiceStartup" -ErrorAction SilentlyContinue

Remove-ItemProperty -Path $rtPath -Name "DisableRealtimeMonitoring" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $rtPath -Name "DisableBehaviorMonitoring" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $rtPath -Name "DisableOnAccessProtection" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $rtPath -Name "DisableIOAVProtection" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $rtPath -Name "DisableIntrusionPreventionSystem" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableRealtimeMonitoring" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableBehaviorMonitoring" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableOnAccessProtection" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableIOAVProtection" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $rtPath -Name "LocalSettingOverrideDisableIntrusionPreventionSystem" -ErrorAction SilentlyContinue

Remove-ItemProperty -Path $spPath -Name "DisableBlockAtFirstSeen" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $spPath -Name "SpynetReporting" -ErrorAction SilentlyContinue
Remove-ItemProperty -Path $spPath -Name "SubmitSamplesConsent" -ErrorAction SilentlyContinue


# ==========================================
# 2. 恢复驱动与服务项启动类型 (Start)
# ==========================================
# 说明：0=Boot引导驱动, 2=自动启动, 3=手动启动
$services = @(
    @{ Name = "WinDefend"; Start = 2 },
    @{ Name = "WdBoot"; Start = 0 },
    @{ Name = "WdFilter"; Start = 0 },
    @{ Name = "WdNisDrv"; Start = 3 },
    @{ Name = "WdNisSvc"; Start = 3 },
    @{ Name = "Sense"; Start = 3 },
    @{ Name = "MDCoreSvc"; Start = 3 },
    @{ Name = "SgrmAgent"; Start = 2 },
    @{ Name = "SgrmBroker"; Start = 2 },
    @{ Name = "wscsvc"; Start = 2 },
    @{ Name = "SecurityHealthService"; Start = 3 }
)

foreach ($s in $services) {
    $servicePath = "HKLM:\SYSTEM\CurrentControlSet\Services\$($s.Name)"
    if (Test-Path $servicePath) {
        Set-ItemProperty -Path $servicePath -Name "Start" -Value $s.Start -Force -ErrorAction SilentlyContinue
    }
}


# ==========================================
# 3. 恢复并启用后台计划任务
# ==========================================
$tasks = @(
    "Microsoft\Windows\Windows Defender\Windows Defender Cache Maintenance",
    "Microsoft\Windows\Windows Defender\Windows Defender Cleanup",
    "Microsoft\Windows\Windows Defender\Windows Defender Scheduled Scan",
    "Microsoft\Windows\Windows Defender\Windows Defender Verification"
)

foreach ($task in $tasks) {
    # 启用计划任务
    schtasks /Change /TN $task /Enable 2>$null
}