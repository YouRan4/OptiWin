//go:build windows

package services

import (
	"OptiWin/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows/registry"
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
	return reportErr("EnableSuperfetch", utils.RegSetDWordE(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management\PrefetchParameters`,
		"EnableSuperfetch", 3))
}

func DisableSuperfetch() bool {
	return reportErr("DisableSuperfetch", utils.RegSetDWordE(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management\PrefetchParameters`,
		"EnableSuperfetch", 0))
}

func GetFullscreenOptimizationStatus() bool {
	v, err := utils.RegReadDWord(registry.CURRENT_USER,
		`System\GameConfigStore`, "GameDVR_FSEBehaviorMode")
	return err != nil || v == 0
}

func fullscreenEnableForUser(hive registry.Key, path string) error {
	errs := []error{
		utils.RegSetDWordE(hive, path, "GameDVR_FSEBehaviorMode", 0),
		deleteIgnoringMissing(hive, path, "GameDVR_FSEBehavior"),
		deleteIgnoringMissing(hive, path, "GameDVR_HonorUserFSEBehaviorMode"),
		deleteIgnoringMissing(hive, path, "GameDVR_DXGIHonorFSEWindowsCompatible"),
		deleteIgnoringMissing(hive, path, "GameDVR_EFSEFeatureFlags"),
	}
	return errors.Join(errs...)
}

func fullscreenDisableForUser(hive registry.Key, path string) error {
	errs := []error{
		utils.RegSetDWordE(hive, path, "GameDVR_FSEBehaviorMode", 2),
		utils.RegSetDWordE(hive, path, "GameDVR_HonorUserFSEBehaviorMode", 1),
		utils.RegSetDWordE(hive, path, "GameDVR_DXGIHonorFSEWindowsCompatible", 1),
		utils.RegSetDWordE(hive, path, "GameDVR_EFSEFeatureFlags", 0),
		utils.RegSetDWordE(hive, path, "GameDVR_FSEBehavior", 2),
	}
	return errors.Join(errs...)
}

func EnableFullscreenOptimization() bool {
	err := fullscreenEnableForUser(registry.CURRENT_USER, `System\GameConfigStore`)
	if err != nil {
		return reportErr("EnableFullscreenOptimization", err)
	}
	// HKU\.DEFAULT 默认只允许 SYSTEM 写入，普通管理员写失败属于正常情况，不影响当前用户设置
	if err := fullscreenEnableForUser(registry.USERS, `.DEFAULT\System\GameConfigStore`); err != nil {
		utils.Logf("[services] EnableFullscreenOptimization 写入 .DEFAULT 失败（可忽略）: %v", err)
	}
	return true
}

func DisableFullscreenOptimization() bool {
	err := fullscreenDisableForUser(registry.CURRENT_USER, `System\GameConfigStore`)
	if err != nil {
		return reportErr("DisableFullscreenOptimization", err)
	}
	if err := fullscreenDisableForUser(registry.USERS, `.DEFAULT\System\GameConfigStore`); err != nil {
		utils.Logf("[services] DisableFullscreenOptimization 写入 .DEFAULT 失败（可忽略）: %v", err)
	}
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

func windowedSetOpt(hive registry.Key, path, setting string) error {
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
	return utils.RegWriteStringE(hive, path, "DirectXUserGlobalSettings", newVal)
}

func EnableWindowedOptimization() bool {
	err := windowedSetOpt(registry.CURRENT_USER,
		`Software\Microsoft\DirectX\UserGpuPreferences`,
		"SwapEffectUpgradeEnable=1")
	return reportErr("EnableWindowedOptimization", err)
}

func DisableWindowedOptimization() bool {
	err := windowedSetOpt(registry.CURRENT_USER,
		`Software\Microsoft\DirectX\UserGpuPreferences`,
		"SwapEffectUpgradeEnable=0")
	return reportErr("DisableWindowedOptimization", err)
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
	err := errors.Join(
		deleteIgnoringMissing(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "MPOEnabled"),
		deleteIgnoringMissing(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode"),
	)
	return reportErr("EnableMpo", err)
}

func DisableMpo() bool {
	err := errors.Join(
		utils.RegSetDWordE(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "MPOEnabled", 0),
		utils.RegSetDWordE(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode", 5),
	)
	return reportErr("DisableMpo", err)
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
	return utils.ExecuteString("Enable-MMAgent -mc")
}

func DisableMemoryCompression() bool {
	return utils.ExecuteString("Disable-MMAgent -mc")
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

// --- Game Bar ---

// waitFor 轮询 cond 直到为真或超时
func waitFor(timeout time.Duration, cond func() bool) bool {
	deadline := time.Now().Add(timeout)
	for {
		if cond() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func RemoveGameBar() bool {
	// 只针对当前用户卸载（与恢复路径一致，且不受其他用户配置影响）；
	// 原脚本不检查结果，这里补上状态确认和可读错误输出
	r := utils.ExecuteStringTimeout(`
		$ErrorActionPreference = 'Continue'
		try {
			Get-AppxPackage -Name Microsoft.XboxGamingOverlay | Remove-AppxPackage -Confirm:$false -ErrorAction Stop
		} catch {
			Write-Output ("Remove-AppxPackage 失败: " + $_.Exception.Message)
		}
		Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\GameDVR" -Name "AppCaptureEnabled" -Value 0 -Force
		Set-ItemProperty -Path "HKCU:\System\GameConfigStore" -Name "GameDVR_Enabled" -Value 0 -Force
		foreach ($protocol in @("ms-gamebar", "ms-gamebarservices", "ms-gamingoverlay")) {
			$path = "HKCU:\Software\Classes\$protocol"
			if (-not (Test-Path $path)) { New-Item -Path $path -Force | Out-Null }
			New-ItemProperty -Path $path -Name "NoOpenWith" -Value "" -PropertyType String -Force | Out-Null
		}
		$deadline = (Get-Date).AddSeconds(8)
		do {
			if (-not (Get-AppxPackage -Name Microsoft.XboxGamingOverlay)) { exit 0 }
			Start-Sleep -Milliseconds 250
		} while ((Get-Date) -lt $deadline)
		Write-Output "Game Bar 包仍在："
		Get-AppxPackage -Name Microsoft.XboxGamingOverlay | Select-Object -ExpandProperty PackageFullName
		exit 1
	`, 35*time.Second)
	if !r {
		utils.Logf("[services] RemoveGameBar 未成功：Game Bar 包仍在或命令超时")
	}
	return r
}

func RestoreGameBar() bool {
	// 重新注册包并确认状态；原脚本恒退 0，不确认就无法判断成败
	r := utils.ExecuteStringTimeout(`
		$manifest = (Get-AppxPackage -AllUsers Microsoft.XboxGamingOverlay).InstallLocation + "\AppXManifest.xml"
		if (Test-Path $manifest) { Add-AppxPackage -DisableDevelopmentMode -Register "$manifest" -ErrorAction SilentlyContinue }
		foreach ($protocol in @("ms-gamebar", "ms-gamebarservices", "ms-gamingoverlay")) {
			$path = "HKCU:\Software\Classes\$protocol"
			if (Test-Path $path) { Remove-Item -Path $path -Recurse -Force -ErrorAction SilentlyContinue }
		}
		Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\GameDVR" -Name "AppCaptureEnabled" -Value 1 -Force
		Set-ItemProperty -Path "HKCU:\System\GameConfigStore" -Name "GameDVR_Enabled" -Value 1 -Force
		$deadline = (Get-Date).AddSeconds(8)
		do {
			if (Get-AppxPackage Microsoft.XboxGamingOverlay) { exit 0 }
			Start-Sleep -Milliseconds 250
		} while ((Get-Date) -lt $deadline)
		exit 1
	`, 35*time.Second)
	if !r {
		utils.Logf("[services] RestoreGameBar 未成功：Game Bar 包仍未注册或命令超时")
	}
	return r
}
