//go:build windows

package services

import (
	"fmt"
	"os/exec"
	"syscall"

	"OptiWin/utils"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func GetDefenderStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows Defender`, "DisableAntiSpyware")
	if err == nil && v == 1 {
		return false
	}

	mgr, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_CONNECT)
	if err != nil {
		return false
	}
	defer windows.CloseServiceHandle(mgr)

	svcName, _ := syscall.UTF16PtrFromString("WinDefend")
	svc, err := windows.OpenService(mgr, svcName, windows.SERVICE_QUERY_STATUS)
	if err != nil {
		return false
	}
	defer windows.CloseServiceHandle(svc)

	var status windows.SERVICE_STATUS
	err = windows.QueryServiceStatus(svc, &status)
	if err != nil {
		return true
	}
	if status.CurrentState != windows.SERVICE_RUNNING {
		return false
	}

	v, err = utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows Defender`, "DisableAntiSpyware")
	if err == nil && v == 1 {
		return false
	}
	return true
}

func disableService(name string) {
	keyPath := `HKLM\SYSTEM\CurrentControlSet\Services\` + name
	cmd := exec.Command(utils.GetPowerRunPath(), "reg.exe", "add", keyPath,
		"/v", "Start", "/t", "REG_DWORD", "/d", "4", "/f")
	utils.HideWindow(cmd)
	cmd.Run()
}

func enableService(name string) {
	keyPath := `HKLM\SYSTEM\CurrentControlSet\Services\` + name
	cmd := exec.Command(utils.GetPowerRunPath(), "reg.exe", "add", keyPath,
		"/v", "Start", "/t", "REG_DWORD", "/d", "2", "/f")
	utils.HideWindow(cmd)
	cmd.Run()
}

func regAddDWord(path, name string, value uint32) {
	utils.RunHide("reg", "add", `HKLM\`+path,
		"/v", name, "/t", "REG_DWORD", "/d", fmt.Sprintf("%d", value), "/f")
}

func regDeleteValue(path, name string) {
	utils.RunHide("reg", "delete", `HKLM\`+path, "/v", name, "/f")
}

func RestoreDefender() bool {
	defPath := `SOFTWARE\Policies\Microsoft\Windows Defender`
	rtPath := `SOFTWARE\Policies\Microsoft\Windows Defender\Real-Time Protection`

	regDeleteValue(defPath, "DisableAntiSpyware")
	regDeleteValue(defPath, "ServiceKeepAlive")
	regDeleteValue(defPath, "AllowFastServiceStartup")

	regDeleteValue(rtPath, "DisableRealtimeMonitoring")
	regDeleteValue(rtPath, "DisableBehaviorMonitoring")
	regDeleteValue(rtPath, "DisableOnAccessProtection")
	regDeleteValue(rtPath, "DisableIOAVProtection")
	regDeleteValue(rtPath, "DisableIntrusionPreventionSystem")
	regDeleteValue(rtPath, "LocalSettingOverrideDisableRealtimeMonitoring")
	regDeleteValue(rtPath, "LocalSettingOverrideDisableBehaviorMonitoring")
	regDeleteValue(rtPath, "LocalSettingOverrideDisableOnAccessProtection")
	regDeleteValue(rtPath, "LocalSettingOverrideDisableIOAVProtection")
	regDeleteValue(rtPath, "LocalSettingOverrideDisableIntrusionPreventionSystem")

	enableService("wscsvc")
	enableService("SecurityHealthService")
	enableService("WinDefend")
	enableService("WdBoot")
	enableService("WdFilter")
	enableService("WdNisDrv")
	enableService("WdNisSvc")
	enableService("Sense")
	enableService("MDCoreSvc")
	enableService("SgrmAgent")
	enableService("SgrmBroker")

	utils.CleanupPowerRun()
	return true
}

func DisableDefenderEngine() bool {
	defPath := `SOFTWARE\Policies\Microsoft\Windows Defender`
	rtPath := `SOFTWARE\Policies\Microsoft\Windows Defender\Real-Time Protection`

	regAddDWord(defPath, "DisableAntiSpyware", 1)
	regAddDWord(defPath, "ServiceKeepAlive", 0)
	regAddDWord(defPath, "AllowFastServiceStartup", 0)

	regAddDWord(rtPath, "DisableRealtimeMonitoring", 1)
	regAddDWord(rtPath, "DisableBehaviorMonitoring", 1)
	regAddDWord(rtPath, "DisableOnAccessProtection", 1)
	regAddDWord(rtPath, "DisableIOAVProtection", 1)
	regAddDWord(rtPath, "DisableIntrusionPreventionSystem", 1)
	regAddDWord(rtPath, "LocalSettingOverrideDisableRealtimeMonitoring", 0)
	regAddDWord(rtPath, "LocalSettingOverrideDisableBehaviorMonitoring", 0)
	regAddDWord(rtPath, "LocalSettingOverrideDisableOnAccessProtection", 0)
	regAddDWord(rtPath, "LocalSettingOverrideDisableIOAVProtection", 0)
	regAddDWord(rtPath, "LocalSettingOverrideDisableIntrusionPreventionSystem", 0)

	return true
}

func DisableAllServices() bool {
	DisableDefenderEngine()

	disableService("WinDefend")
	disableService("WdBoot")
	disableService("WdFilter")
	disableService("WdNisDrv")
	disableService("WdNisSvc")
	disableService("Sense")
	disableService("MDCoreSvc")
	disableService("SgrmAgent")
	disableService("SgrmBroker")
	disableService("wscsvc")
	disableService("SecurityHealthService")

	utils.CleanupPowerRun()
	return true
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
