//go:build windows
package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"
	"golang.org/x/sys/windows/registry"
	"OptiWin/utils"
)

func GetUltimatePerformanceStatus() bool {
	p := getActiveSchemeGuid()
	if p == nil {
		return false
	}
	defer freeGuid(p)
	return p.Data1 == ultimateGuid.Data1 &&
		p.Data2 == ultimateGuid.Data2 &&
		p.Data3 == ultimateGuid.Data3 &&
		p.Data4 == ultimateGuid.Data4
}

func EnableUltimatePerformance() bool {
	procPowerDuplicatePowerScheme.Call(0, uintptr(unsafe.Pointer(&ultimateGuid)), 0)
	procPowerSetActiveScheme.Call(0, uintptr(unsafe.Pointer(&ultimateGuid)))
	return true
}

func DisableUltimatePerformance() bool {
	procPowerSetActiveScheme.Call(0, uintptr(unsafe.Pointer(&balancedGuid)))
	return true
}

func GetCStateStatus() bool {
	p := getActiveSchemeGuid()
	if p == nil {
		return true
	}
	defer freeGuid(p)
	v, ok := cstateRead(p)
	if !ok {
		return true
	}
	return v == 1
}

func EnableCState() bool {
	p := getActiveSchemeGuid()
	if p == nil {
		return false
	}
	defer freeGuid(p)
	cstateWrite(p, 1)
	procPowerSetActiveScheme.Call(0, uintptr(unsafe.Pointer(p)))
	return true
}

func DisableCState() bool {
	p := getActiveSchemeGuid()
	if p == nil {
		return false
	}
	defer freeGuid(p)
	cstateWrite(p, 0)
	procPowerSetActiveScheme.Call(0, uintptr(unsafe.Pointer(p)))
	return true
}

func GetSuperfetchStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management\PrefetchParameters`,
		"EnableSuperfetch")
	return err == nil && v == 3
}

func EnableSuperfetch() bool {
	return utils.RegSetDWordBool(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management\PrefetchParameters`,
		"EnableSuperfetch", 3)
}

func DisableSuperfetch() bool {
	return utils.RegSetDWordBool(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management\PrefetchParameters`,
		"EnableSuperfetch", 0)
}

func GetFullscreenOptimizationStatus() bool {
	v, err := utils.RegReadDWord(registry.CURRENT_USER,
		`System\GameConfigStore`, "GameDVR_FSEBehaviorMode")
	return err != nil || v == 0
}

func fullscreenEnableForUser(hive registry.Key, path string) {
	utils.RegSetDWord(hive, path, "GameDVR_FSEBehaviorMode", 0)
	utils.RegDeleteValue(hive, path, "GameDVR_FSEBehavior")
	utils.RegDeleteValue(hive, path, "GameDVR_HonorUserFSEBehaviorMode")
	utils.RegDeleteValue(hive, path, "GameDVR_DXGIHonorFSEWindowsCompatible")
	utils.RegDeleteValue(hive, path, "GameDVR_EFSEFeatureFlags")
}

func fullscreenDisableForUser(hive registry.Key, path string) {
	utils.RegSetDWord(hive, path, "GameDVR_FSEBehaviorMode", 2)
	utils.RegSetDWord(hive, path, "GameDVR_HonorUserFSEBehaviorMode", 1)
	utils.RegSetDWord(hive, path, "GameDVR_DXGIHonorFSEWindowsCompatible", 1)
	utils.RegSetDWord(hive, path, "GameDVR_EFSEFeatureFlags", 0)
	utils.RegSetDWord(hive, path, "GameDVR_FSEBehavior", 2)
}

func EnableFullscreenOptimization() bool {
	fullscreenEnableForUser(registry.CURRENT_USER, `System\GameConfigStore`)
	fullscreenEnableForUser(registry.USERS, `.DEFAULT\System\GameConfigStore`)
	return true
}

func DisableFullscreenOptimization() bool {
	fullscreenDisableForUser(registry.CURRENT_USER, `System\GameConfigStore`)
	fullscreenDisableForUser(registry.USERS, `.DEFAULT\System\GameConfigStore`)
	return true
}

func GetWindowedOptimizationStatus() bool {
	val, _, err := utils.RegReadString(registry.CURRENT_USER,
		`Software\Microsoft\DirectX\UserGpuPreferences`,
		"DirectXUserGlobalSettings")
	if err != nil {
		return true
	}
	return !strings.Contains(val, "SwapEffectUpgradeEnable=0")
}

func windowedSetOpt(hive registry.Key, path, setting string) {
	val, _, err := utils.RegReadString(hive, path, "DirectXUserGlobalSettings")
	var newVal string
	if err != nil || val == "" {
		newVal = setting + ";"
	} else {
		cleaned := val
		cleaned = strings.ReplaceAll(cleaned, "SwapEffectUpgradeEnable=1;", "")
		cleaned = strings.ReplaceAll(cleaned, "SwapEffectUpgradeEnable=1", "")
		cleaned = strings.ReplaceAll(cleaned, "SwapEffectUpgradeEnable=0;", "")
		cleaned = strings.ReplaceAll(cleaned, "SwapEffectUpgradeEnable=0", "")
		newVal = setting + ";" + cleaned
	}
	utils.RegWriteString(hive, path, "DirectXUserGlobalSettings", newVal)
}

func EnableWindowedOptimization() bool {
	windowedSetOpt(registry.CURRENT_USER,
		`Software\Microsoft\DirectX\UserGpuPreferences`,
		"SwapEffectUpgradeEnable=1")
	return true
}

func DisableWindowedOptimization() bool {
	windowedSetOpt(registry.CURRENT_USER,
		`Software\Microsoft\DirectX\UserGpuPreferences`,
		"SwapEffectUpgradeEnable=0")
	return true
}

func GetMpoStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\Dwm`, "MPOEnabled")
	if err == nil && v == 0 {
		return false
	}
	v, err = utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode")
	if err == nil && v == 5 {
		return false
	}
	return true
}

func EnableMpo() bool {
	utils.RegDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "MPOEnabled")
	utils.RegDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode")
	return true
}

func DisableMpo() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "MPOEnabled", 0)
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode", 5)
	return true
}

func GetMemoryCompressionStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management`,
		"FeatureSettings")
	if err != nil {
		return true
	}
	return (v & 1) == 0
}

func EnableMemoryCompression() bool {
	cmd := exec.Command("powershell", "-Command", "Enable-MMAgent -mc")
	utils.HideWindow(cmd)
	return cmd.Run() == nil
}

func DisableMemoryCompression() bool {
	cmd := exec.Command("powershell", "-Command", "Disable-MMAgent -mc")
	utils.HideWindow(cmd)
	return cmd.Run() == nil
}

func ClearShaderCache() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return "无法获取 LOCALAPPDATA 路径"
	}

	dirs := []string{
		filepath.Join(localAppData, "NVIDIA", "DXCache"),
		filepath.Join(localAppData, "NVIDIA", "GLCache"),
		filepath.Join(localAppData, "AMD", "DxCache"),
	}

	total := 0
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			p := filepath.Join(dir, e.Name())
			os.RemoveAll(p)
			total++
		}
	}
	return fmt.Sprintf("已清理 %d 个着色器缓存文件", total)
}
