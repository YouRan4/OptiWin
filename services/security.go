//go:build windows

package services

import (
	"OptiWin/utils"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// realtimeScanPolicyPath Defender 实时保护组策略键路径
const realtimeScanPolicyPath = `SOFTWARE\Policies\Microsoft\Windows Defender\Real-Time Protection`

// defenderPolicyPath Defender 主策略键（DisableAntiSpyware 触发"由管理员管理"横幅）
const defenderPolicyPath = `SOFTWARE\Policies\Microsoft\Windows Defender`

// realtimeProtectionValues 实时保护全部组策略开关（1=禁用）
var realtimeProtectionValues = []string{
	"DisableRealtimeMonitoring",     // 实时监控
	"DisableBehaviorMonitoring",     // 行为监控
	"DisableOnAccessProtection",     // 访问时保护
	"DisableIOAVProtection",         // 下载/附件扫描
	"DisableScanOnRealtimeEnable",   // 实时保护启用时扫描
	"DisableIntrusionPreventionSystem", // 入侵防护系统
}

// setRealtimeProtection 禁用/恢复实时保护（写/删全部组策略键 + DisableAntiSpyware）
// DisableAntiSpyware 触发"由管理员管理"横幅；必须在 TP 关闭后执行，TP 开启时引擎会忽略策略
func setRealtimeProtection(disable bool) bool {
	path := defenderPolicyPath
	vals := realtimeProtectionValues
	if disable {
		if utils.RegSetDWordE(registry.LOCAL_MACHINE, path, "DisableAntiSpyware", 1) != nil {
			return false
		}
		for _, n := range vals {
			if utils.RegSetDWordE(registry.LOCAL_MACHINE, realtimeScanPolicyPath, n, 1) != nil {
				return false
			}
		}
		return true
	}
	if utils.RegDeleteValueE(registry.LOCAL_MACHINE, path, "DisableAntiSpyware") != nil {
		return false
	}
	for _, n := range vals {
		if utils.RegDeleteValueE(registry.LOCAL_MACHINE, realtimeScanPolicyPath, n) != nil {
			return false
		}
	}
	return true
}

// getMpStatus 查询 Get-MpComputerStatus 的指定布尔属性
func getMpStatus(prop string) bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive",
		"-Command", "(Get-MpComputerStatus)."+prop)
	utils.HideWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "True"
}

// GetTamperProtectionStatus 篡改防护是否开启（开启时需先关闭才能禁实时保护）
func GetTamperProtectionStatus() bool {
	return getMpStatus("IsTamperProtected")
}

// GetRealtimeProtectionEnabled 实时保护是否运行中
func GetRealtimeProtectionEnabled() bool {
	return getMpStatus("RealTimeProtectionEnabled")
}

// OpenWindowsSecurity 打开安全中心篡改防护设置页
func OpenWindowsSecurity() bool {
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler",
		"windowsdefender://threatsettings")
	utils.HideWindow(cmd)
	return cmd.Start() == nil
}

// GetCoreServicesDisabled 实时保护是否已禁用（DisableAntiSpyware 存在且=1）
func GetCoreServicesDisabled() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE, defenderPolicyPath, "DisableAntiSpyware")
	return err == nil && v == 1
}

// DisableCoreServices 禁用实时保护释放资源（不动服务，MsMpEng.exe 保留但停止扫描）
func DisableCoreServices() bool {
	return setRealtimeProtection(true)
}

// EnableCoreServices 恢复实时保护
func EnableCoreServices() bool {
	return setRealtimeProtection(false)
}

func GetUacStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableLUA")
	return err == nil && v == 1
}

func EnableUac() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableLUA", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableVirtualization", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableInstallerDetection", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "PromptOnSecureDesktop", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "ConsentPromptBehaviorAdmin", 2)
	return true
}

func DisableUac() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableLUA", 0)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "ConsentPromptBehaviorAdmin", 0)
	return true
}

func GetVbsStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Policies\Microsoft\Windows\DeviceGuard`, "EnableVirtualizationBasedSecurity")
	if err == nil && v == 1 {
		return true
	}
	v, err = utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\ControlSet001\Control\DeviceGuard`, "EnableVirtualizationBasedSecurity")
	return err == nil && v == 1
}

func EnableVbs() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\DeviceGuard`, "EnableVirtualizationBasedSecurity", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Control\DeviceGuard`, "EnableVirtualizationBasedSecurity", 1)
	return true
}

func DisableVbs() bool {
	utils.RegDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\DeviceGuard`, "EnableVirtualizationBasedSecurity")
	utils.RegDeleteValue(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Control\DeviceGuard`, "EnableVirtualizationBasedSecurity")
	return true
}

func GetMemoryIntegrityStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\ControlSet001\Control\DeviceGuard\Scenarios\HypervisorEnforcedCodeIntegrity`, "Enabled")
	return err == nil && v == 1
}

func EnableMemoryIntegrity() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE,
		`SYSTEM\ControlSet001\Control\DeviceGuard\Scenarios\HypervisorEnforcedCodeIntegrity`, "Enabled", 1)
	return true
}

func DisableMemoryIntegrity() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE,
		`SYSTEM\ControlSet001\Control\DeviceGuard\Scenarios\HypervisorEnforcedCodeIntegrity`, "Enabled", 0)
	return true
}

// -----IEFO
const ifeoPath = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Image File Execution Options`

var ifeoBlockedExecutables = map[string]bool{
	"winlogon.exe":      true,
	"csrss.exe":         true,
	"lsass.exe":         true,
	"smss.exe":          true,
	"services.exe":      true,
	"svchost.exe":       true,
	"explorer.exe":      true,
	"dwm.exe":           true,
	"wininit.exe":       true,
	"fontdrvhost.exe":   true,
	"sihost.exe":        true,
	"RuntimeBroker.exe": true,
}

func isValidIfeoExeName(name string) bool {
	if name == "" || !strings.HasSuffix(strings.ToLower(name), ".exe") {
		return false
	}
	if strings.ContainsAny(name, `\/:;*?<>|`) {
		return false
	}
	if ifeoBlockedExecutables[strings.ToLower(name)] {
		return false
	}
	return true
}

func isValidDebuggerPath(path string) bool {
	if path == "" {
		return true
	}
	lower := strings.ToLower(path)
	// 允许环境变量路径（%windir%、%systemroot% 等）或盘符绝对路径
	if !strings.Contains(lower, `:\`) && !strings.Contains(lower, `%`) {
		return false
	}
	if strings.HasPrefix(lower, `\\`) {
		return false
	}
	if strings.Contains(lower, "..") {
		return false
	}
	return true
}

type ifeoEntry struct {
	Name       string `json:"name"`
	Debugger   string `json:"debugger"`
	GlobalFlag uint64 `json:"globalFlag"`
}

func ListIfeoEntries() string {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, ifeoPath, registry.READ)
	if err != nil {
		return "[]"
	}
	defer key.Close()

	subKeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return "[]"
	}

	var entries []ifeoEntry
	for _, name := range subKeys {
		subKey, err := registry.OpenKey(key, name, registry.READ)
		if err != nil {
			continue
		}

		debuger, _, _ := subKey.GetStringValue("Debugger")
		globalFlag, _, _ := subKey.GetIntegerValue("GlobalFlag")
		subKey.Close()

		if debuger == "" && globalFlag == 0 {
			continue
		}

		entries = append(entries, ifeoEntry{
			Name:       name,
			Debugger:   debuger,
			GlobalFlag: globalFlag,
		})
	}

	b, _ := json.Marshal(entries)
	return string(b)
}

func AddIfeoEntry(exeName, debugger string) bool {
	if !isValidIfeoExeName(exeName) {
		return false
	}
	if !isValidDebuggerPath(debugger) {
		return false
	}

	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, ifeoPath+`\`+exeName, registry.SET_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	if debugger != "" {
		err = key.SetStringValue("Debugger", debugger)
		return err == nil
	}
	return true
}

func RemoveIfeoEntry(exeName string) bool {
	if exeName == "" {
		return false
	}

	parentKey, err := registry.OpenKey(registry.LOCAL_MACHINE, ifeoPath, registry.ALL_ACCESS)
	if err != nil {
		return false
	}
	defer parentKey.Close()

	err = registry.DeleteKey(parentKey, exeName)
	return err == nil
}

func GetRunningProcesses() string {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command",
		"Get-Process | Where-Object { $_.ProcessName -ne 'Idle' } | Select-Object -ExpandProperty ProcessName | Sort-Object -Unique | ConvertTo-Json -Compress")
	utils.HideWindow(cmd)
	output, err := cmd.Output()
	if err != nil {
		return "[]"
	}
	result := strings.TrimSpace(string(output))
	if result == "" || result == "null" {
		return "[]"
	}
	return result
}

type DnsConfig struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}

var dnsMap = map[string]DnsConfig{
	"isp":        {Primary: "", Secondary: ""},
	"ali":        {Primary: "223.5.5.5", Secondary: "223.6.6.6"},
	"tencent":    {Primary: "119.29.29.29", Secondary: "119.28.28.28"},
	"baidu":      {Primary: "180.76.76.76", Secondary: "180.76.76.76"},
	"google":     {Primary: "8.8.8.8", Secondary: "8.8.4.4"},
	"cloudflare": {Primary: "1.1.1.1", Secondary: "1.0.0.1"},
	"quad9":      {Primary: "9.9.9.9", Secondary: "149.112.112.112"},
}

func SetDns(code string) bool {
	config, ok := dnsMap[code]
	if !ok {
		return false
	}

	getActiveAdapter := `
$activeAdapter = Get-NetAdapter | Where-Object { 
    $_.Status -eq 'Up' -and 
    $_.InterfaceDescription -notmatch 'Virtual|VPN|Loopback|Bluetooth|Hyper-V|TAP|TUN' 
} | Sort-Object LinkSpeed -Descending | Select-Object -First 1
`

	if code == "isp" {
		psScript := getActiveAdapter + `
Set-DnsClientServerAddress -InterfaceIndex $activeAdapter.InterfaceIndex -ResetServerAddresses
Write-Output "OK"`
		cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
		utils.HideWindow(cmd)
		out, err := cmd.Output()
		if err != nil {
			return false
		}
		return strings.Contains(string(out), "OK")
	}

	psScript := getActiveAdapter + fmt.Sprintf(`
Set-DnsClientServerAddress -InterfaceIndex $activeAdapter.InterfaceIndex -ServerAddresses ('%s', '%s')
Clear-DnsClientCache
Write-Output "OK"`, config.Primary, config.Secondary)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
	utils.HideWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "OK")
}

func GetCurrentDns() string {
	psScript := `
chcp 65001 >$null
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

$activeAdapter = Get-NetAdapter | Where-Object { 
    $_.Status -eq 'Up' -and 
    $_.InterfaceDescription -notmatch 'Virtual|VPN|Loopback|Bluetooth|Hyper-V|TAP|TUN' 
} | Sort-Object LinkSpeed -Descending | Select-Object -First 1

$primaryDns = ""
$secondaryDns = ""

if ($activeAdapter) {
    $dnsInfo = $activeAdapter | Get-DnsClientServerAddress -AddressFamily IPv4
    $servers = $dnsInfo.ServerAddresses
    
    if ($servers.Count -gt 0) { $primaryDns = $servers[0] }
    if ($servers.Count -gt 1) { $secondaryDns = $servers[1] }
}

[PSCustomObject]@{
    adapter    = if ($activeAdapter) { $activeAdapter.InterfaceDescription } else { "" }
    primaryDns = $primaryDns
    secondaryDns = $secondaryDns
} | ConvertTo-Json -Compress
`
	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
	utils.HideWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return "{}"
	}
	result := strings.TrimSpace(string(out))
	if result == "" || result == "null" {
		return "{}"
	}
	return result
}
