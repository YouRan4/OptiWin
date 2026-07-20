//go:build windows
package services

import (
	"syscall"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"OptiWin/utils"
)

func GetDefenderStatus() bool {
	mgr, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_CONNECT)
	if err != nil {
		return true
	}
	defer windows.CloseServiceHandle(mgr)

	svcName, _ := syscall.UTF16PtrFromString("WinDefend")
	svc, err := windows.OpenService(mgr, svcName, windows.SERVICE_QUERY_STATUS)
	if err != nil {
		return true
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

	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows Defender`, "DisableAntiSpyware")
	if err == nil && v == 1 {
		return false
	}
	return true
}

func EnableDefender() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`, "DisableAntiSpyware", 0)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`, "DisableAntiVirus", 0)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows Defender`, "DisableAntiSpyware", 0)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows Defender`, "DisableAntiVirus", 0)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender\Features`, "TamperProtection", 0)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender\Features`, "TamperProtectionSource", 0)

	utils.RunHide("gpupdate", "/target:computer", "/force")

	utils.ScHide("config", "WinDefend", "start=", "auto")
	utils.ScHide("start", "WinDefend")
	utils.ScHide("config", "SecurityHealthService", "start=", "auto")
	utils.ScHide("start", "SecurityHealthService")
	utils.ScHide("config", "MDCoreSvc", "start=", "auto")
	utils.ScHide("config", "wscsvc", "start=", "auto")
	utils.ScHide("start", "wscsvc")
	return true
}

func DisableDefender() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender\Features`, "TamperProtection", 0)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender\Features`, "TamperProtectionSource", 0)

	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`, "DisableAntiSpyware", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`, "DisableAntiVirus", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows Defender`, "DisableAntiSpyware", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows Defender`, "DisableAntiVirus", 1)

	utils.RunHide("gpupdate", "/target:computer", "/force")

	utils.ScHide("stop", "WinDefend")
	utils.ScHide("config", "WinDefend", "start=", "disabled")
	utils.ScHide("stop", "SecurityHealthService")
	utils.ScHide("config", "SecurityHealthService", "start=", "disabled")
	utils.ScHide("stop", "MDCoreSvc")
	utils.ScHide("config", "MDCoreSvc", "start=", "disabled")
	utils.ScHide("stop", "wscsvc")
	utils.ScHide("config", "wscsvc", "start=", "disabled")
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
