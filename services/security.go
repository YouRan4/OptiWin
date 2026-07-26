//go:build windows

package services

import (
	"OptiWin/utils"
	"syscall"

	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func GetSecurityHealthServiceStatus() bool {
	mgr, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_CONNECT)
	if err != nil {
		return false
	}
	defer windows.CloseServiceHandle(mgr)

	svcName, _ := syscall.UTF16PtrFromString("SecurityHealthService")
	svc, err := windows.OpenService(mgr, svcName, windows.SERVICE_QUERY_STATUS)
	if err != nil {
		return false
	}
	defer windows.CloseServiceHandle(svc)

	var status windows.SERVICE_STATUS
	err = windows.QueryServiceStatus(svc, &status)
	if err != nil {
		return false
	}
	return status.CurrentState == windows.SERVICE_RUNNING
}

func RestoreDefender() bool {
	return utils.SuperExecute(utils.RestoreDefenderScript)
}

func DisableAllServices() bool {
	return utils.SuperExecute(utils.DisableDefenderScript)
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
	if exeName == "" {
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
	Doh       string `json:"doh"`
}

var dnsMap = map[string]DnsConfig{
	"isp":        {Primary: "", Secondary: "", Doh: ""},
	"ali":        {Primary: "223.5.5.5", Secondary: "223.6.6.6", Doh: "https://dns.alidns.com/dns-query"},
	"tencent":    {Primary: "119.29.29.29", Secondary: "119.28.28.28", Doh: "https://doh.pub/dns-query"},
	"baidu":      {Primary: "180.76.76.76", Secondary: "180.76.76.76", Doh: ""},
	"google":     {Primary: "8.8.8.8", Secondary: "8.8.4.4", Doh: "https://dns.google/dns-query"},
	"cloudflare": {Primary: "1.1.1.1", Secondary: "1.0.0.1", Doh: "https://cloudflare-dns.com/dns-query"},
	"quad9":      {Primary: "9.9.9.9", Secondary: "149.112.112.112", Doh: "https://dns1.quad9.net/dns-query"},
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
Remove-DnsClientDohServerAddress -ServerAddress * -ErrorAction SilentlyContinue
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
Set-DnsClientServerAddress -InterfaceIndex $activeAdapter.InterfaceIndex -ServerAddresses ('%s', '%s')`, config.Primary, config.Secondary)

	if config.Doh != "" {
		psScript += fmt.Sprintf(`
$existing = Get-DnsClientDohServerAddress -ServerAddress '%s' -ErrorAction SilentlyContinue
if (-not $existing) {
    Add-DnsClientDohServerAddress -ServerAddress '%s' -DohTemplate '%s' -AllowFallbackToUdp $true -AutoUpgrade $true -ErrorAction SilentlyContinue
} else {
    Set-DnsClientDohServerAddress -ServerAddress '%s' -DohTemplate '%s' -AllowFallbackToUdp $true -AutoUpgrade $true -ErrorAction SilentlyContinue
}`, config.Primary, config.Primary, config.Doh, config.Primary, config.Doh)

		if config.Secondary != "" {
			psScript += fmt.Sprintf(`
$existingSec = Get-DnsClientDohServerAddress -ServerAddress '%s' -ErrorAction SilentlyContinue
if (-not $existingSec) {
    Add-DnsClientDohServerAddress -ServerAddress '%s' -DohTemplate '%s' -AllowFallbackToUdp $true -AutoUpgrade $true -ErrorAction SilentlyContinue
} else {
    Set-DnsClientDohServerAddress -ServerAddress '%s' -DohTemplate '%s' -AllowFallbackToUdp $true -AutoUpgrade $true -ErrorAction SilentlyContinue
}`, config.Secondary, config.Secondary, config.Doh, config.Secondary, config.Doh)
		}
	}

	psScript += `
Clear-DnsClientCache
Write-Output "OK"`
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
$dohTemplate = ""
$dohStatus = "未启用"

if ($activeAdapter) {
    $dnsInfo = $activeAdapter | Get-DnsClientServerAddress -AddressFamily IPv4
    $servers = $dnsInfo.ServerAddresses
    
    if ($servers.Count -gt 0) { $primaryDns = $servers[0] }
    if ($servers.Count -gt 1) { $secondaryDns = $servers[1] }

    foreach ($serverIP in $servers) {
        $dohConfig = Get-DnsClientDohServerAddress -ServerAddress $serverIP -ErrorAction SilentlyContinue
        if ($dohConfig -and $dohConfig.DohTemplate) {
            $dohTemplate = $dohConfig.DohTemplate
            if ($dohConfig.AutoUpgrade -eq $true) {
                $dohStatus = "已启用"
                break
            }
        }
    }
}

[PSCustomObject]@{
    adapter      = if ($activeAdapter) { $activeAdapter.InterfaceDescription } else { "" }
    primaryDns   = $primaryDns
    secondaryDns = $secondaryDns
    doh          = $dohTemplate
    dohStatus    = $dohStatus
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
