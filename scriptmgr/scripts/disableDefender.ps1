# ==========================================
# 1. 策略路径与实时防护注册表封锁
# ==========================================
$defPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender"
$rtPath  = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender\Real-Time Protection"
$spPath  = "HKLM:\SOFTWARE\Policies\Microsoft\Windows Defender\SpyNet"

# 确保核心策略目录存在
New-Item -Path $defPath -Force | Out-Null
New-Item -Path $rtPath -Force | Out-Null
New-Item -Path $spPath -Force | Out-Null

# 基础策略配置
Set-ItemProperty -Path $defPath -Name "DisableAntiSpyware" -Value 1 -Force -ErrorAction SilentlyContinue
Set-ItemProperty -Path $defPath -Name "ServiceKeepAlive" -Value 0 -Force -ErrorAction SilentlyContinue
Set-ItemProperty -Path $defPath -Name "AllowFastServiceStartup" -Value 0 -Force -ErrorAction SilentlyContinue

# 实时监控与行为感知禁用
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

# 云端上报与样本提交关闭（防云端联动唤醒）
Set-ItemProperty -Path $spPath -Name "DisableBlockAtFirstSeen" -Value 1 -Force
Set-ItemProperty -Path $spPath -Name "SpynetReporting" -Value 0 -Force
Set-ItemProperty -Path $spPath -Name "SubmitSamplesConsent" -Value 2 -Force


# ==========================================
# 2. 核心驱动与服务项强行置死 (Start = 4)
# ==========================================
$services = @(
    "WinDefend", 
    "WdBoot", 
    "WdFilter", 
    "WdNisDrv", 
    "WdNisSvc", 
    "Sense", 
    "MDCoreSvc", 
    "SgrmAgent", 
    "SgrmBroker", 
    "wscsvc", 
    "SecurityHealthService"
)

foreach ($s in $services) {
    $servicePath = "HKLM:\SYSTEM\CurrentControlSet\Services\$s"
    if (Test-Path $servicePath) {
        Set-ItemProperty -Path $servicePath -Name "Start" -Value 4 -Force -ErrorAction SilentlyContinue
    }
}


# ==========================================
# 3. 禁用相关的后台计划任务（防止残留唤醒）
# ==========================================
$tasks = @(
    "Microsoft\Windows\Windows Defender\Windows Defender Cache Maintenance",
    "Microsoft\Windows\Windows Defender\Windows Defender Cleanup",
    "Microsoft\Windows\Windows Defender\Windows Defender Scheduled Scan",
    "Microsoft\Windows\Windows Defender\Windows Defender Verification"
)

foreach ($task in $tasks) {
    # 尝试禁用任务计划，忽略可能不存在的任务报错
    schtasks /Change /TN $task /Disable 2>$null
}