//go:build windows

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
	"OptiWin/services"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func logError(msg string) {
	logPath := filepath.Join(os.TempDir(), "revitool.log")
	f, _ := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(msg + "\n")
}

// GetPauseUpdatesStatus 读取暂停更新状态
func (a *App) GetPauseUpdatesStatus() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	_, _, err = k.GetStringValue("PauseUpdatesExpiryTime")
	return err == nil
}

// EnablePauseUpdates 启用暂停更新（暂停到 2126 年）
func (a *App) EnablePauseUpdates() bool {
	k, _, _ := registry.CreateKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`, registry.SET_VALUE)
	defer k.Close()

	k.SetDWordValue("FlightSettingsMaxPauseDays", 37200)
	k.SetStringValue("PauseFeatureUpdatesStartTime", "2023-08-17T12:47:51Z")
	k.SetStringValue("PauseFeatureUpdatesEndTime", "2126-01-01T00:00:00Z")
	k.SetStringValue("PauseQualityUpdatesStartTime", "2023-08-17T12:47:51Z")
	k.SetStringValue("PauseQualityUpdatesEndTime", "2126-01-01T00:00:00Z")
	k.SetStringValue("PauseUpdatesStartTime", "2023-08-17T12:47:51Z")
	k.SetStringValue("PauseUpdatesExpiryTime", "2126-01-01T00:00:00Z")
	return true
}

// DisablePauseUpdates 禁用暂停更新
func (a *App) DisablePauseUpdates() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`, registry.SET_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	for _, v := range []string{
		"FlightSettingsMaxPauseDays",
		"PauseFeatureUpdatesStartTime",
		"PauseFeatureUpdatesEndTime",
		"PauseQualityUpdatesStartTime",
		"PauseQualityUpdatesEndTime",
		"PauseUpdatesStartTime",
		"PauseUpdatesExpiryTime",
	} {
		k.DeleteValue(v)
	}
	return true
}

// GetVisibilityStatus 读取 Windows 更新页面可见性状态
// 返回 true 表示隐藏（不可见），false 表示可见
func (a *App) GetVisibilityStatus() bool {
	val, _, err := regReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"SettingsPageVisibility")
	if err != nil || val == "" {
		return false
	}
	for i := 0; i+18 <= len(val); i++ {
		if val[i:i+18] == "hide:windowsupdate" {
			return true
		}
	}
	return false
}

// EnableVisibility 隐藏 Windows 更新页面
func (a *App) EnableVisibility() bool {
	existing, _, err := regReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"SettingsPageVisibility")
	if err != nil {
		existing = ""
	}
	if existing != "" {
		existing += ";hide:windowsupdate"
	} else {
		existing = "hide:windowsupdate"
	}
	cmd := exec.Command("reg", "add", `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"/v", "SettingsPageVisibility", "/t", "REG_SZ", "/d", existing, "/f")
	hideWindow(cmd)
	if out, err := cmd.CombinedOutput(); err != nil {
		logError(fmt.Sprintf("EnableVisibility 失败: %v\n输出: %s", err, string(out)))
		return false
	}
	return true
}

// DisableVisibility 恢复 Windows 更新页面可见
func (a *App) DisableVisibility() bool {
	existing, _, err := regReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"SettingsPageVisibility")
	if err != nil {
		return true
	}

	var parts []string
	start := 0
	for start < len(existing) {
		end := start
		for end < len(existing) && existing[end] != ';' {
			end++
		}
		part := existing[start:end]
		if part != "hide:windowsupdate" && part != "" {
			parts = append(parts, part)
		}
		start = end + 1
	}

	if len(parts) == 0 {
		cmd := exec.Command("reg", "delete", `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
			"/v", "SettingsPageVisibility", "/f")
		hideWindow(cmd)
		if out, err := cmd.CombinedOutput(); err != nil {
			logError(fmt.Sprintf("DisableVisibility 删除失败: %v\n输出: %s", err, string(out)))
			return false
		}
		return true
	}

	newVal := parts[0]
	for i := 1; i < len(parts); i++ {
		newVal += ";" + parts[i]
	}
	cmd := exec.Command("reg", "add", `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"/v", "SettingsPageVisibility", "/t", "REG_SZ", "/d", newVal, "/f")
	hideWindow(cmd)
	if out, err := cmd.CombinedOutput(); err != nil {
		logError(fmt.Sprintf("DisableVisibility 更新失败: %v\n输出: %s", err, string(out)))
		return false
	}
	return true
}

// GetDriverUpdatesStatus 读取驱动程序更新状态
// 返回 true 表示允许通过 Windows Update 安装驱动
func (a *App) GetDriverUpdatesStatus() bool {
	val, err := regReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Device Metadata`,
		"PreventDeviceMetadataFromNetwork")
	if err != nil || val == 0 {
		return true
	}
	return false
}

// EnableDriverUpdates 允许通过 Windows Update 安装驱动程序
func (a *App) EnableDriverUpdates() bool {
	regSetDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Device Metadata`,
		"PreventDeviceMetadataFromNetwork", 0)
	regDeleteValue(registry.LOCAL_MACHINE,
		`SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate`,
		"ExcludeWUDriversInQualityUpdate")
	regDeleteValue(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`,
		"ExcludeWUDriversInQualityUpdate")
	return true
}

// DisableDriverUpdates 禁止通过 Windows Update 安装驱动程序
func (a *App) DisableDriverUpdates() bool {
	regSetDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Device Metadata`,
		"PreventDeviceMetadataFromNetwork", 1)
	regSetDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate`,
		"ExcludeWUDriversInQualityUpdate", 1)
	regSetDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`,
		"ExcludeWUDriversInQualityUpdate", 1)
	return true
}

// 注册表辅助函数
func regReadDWord(key registry.Key, path, name string) (uint64, error) {
	k, err := registry.OpenKey(key, path, registry.QUERY_VALUE)
	if err != nil {
		return 0, err
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue(name)
	return val, err
}

func regReadString(key registry.Key, path, name string) (string, uint32, error) {
	k, err := registry.OpenKey(key, path, registry.QUERY_VALUE)
	if err != nil {
		return "", 0, err
	}
	defer k.Close()
	return k.GetStringValue(name)
}

func regWriteString(key registry.Key, path, name, value string) bool {
	k, _, _ := registry.CreateKey(key, path, registry.SET_VALUE)
	defer k.Close()
	return k.SetStringValue(name, value) == nil
}

func regDeleteValue(key registry.Key, path, name string) {
	k, err := registry.OpenKey(key, path, registry.SET_VALUE)
	if err != nil {
		return
	}
	defer k.Close()
	k.DeleteValue(name)
}

func regSetDWord(key registry.Key, path, name string, value uint32) {
	k, err := registry.OpenKey(key, path, registry.SET_VALUE)
	if err != nil {
		k, _, _ = registry.CreateKey(key, path, registry.SET_VALUE)
	}
	defer k.Close()
	k.SetDWordValue(name, value)
}

// UpdateKGL 获取最新 KGL 数据并写入注册表
func (a *App) UpdateKGL() string {
	data, errMsg := services.FetchKGL()
	if data == nil {
		return errMsg
	}

	// 写入 HKLM\SOFTWARE\Microsoft\KGL\OneSettings
	kglKey, _, _ := registry.CreateKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\KGL\OneSettings`, registry.SET_VALUE)
	defer kglKey.Close()

	kglKey.SetStringValue("ActivateOnUpdate", fmt.Sprintf("%v", data.ActivateOnUpdate))
	kglKey.SetStringValue("Hash", data.Hash)
	kglKey.SetStringValue("URI", data.URI)
	kglKey.SetStringValue("Version", data.Version)
	kglKey.SetDWordValue("VersionCheckTimeout", uint32(data.VersionCheckTimeout))

	// 写入 HKCU\Software\Microsoft\Windows\CurrentVersion\GameDVR
	userKey, _, _ := registry.CreateKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\GameDVR`, registry.SET_VALUE)
	defer userKey.Close()

	userKey.SetStringValue("KGLRevision", data.Version)
	userKey.SetStringValue("KGLToGCSUpdatedRevision", data.Version)
	return "KGL 更新成功"
}

// --- Power API via powrprof.dll ---

var (
	procPowerGetActiveScheme       = powrProf.NewProc("PowerGetActiveScheme")
	procPowerSetActiveScheme       = powrProf.NewProc("PowerSetActiveScheme")
	procPowerDuplicatePowerScheme  = powrProf.NewProc("PowerDuplicatePowerScheme")
	procPowerReadACValueIndex      = powrProf.NewProc("PowerReadACValueIndex")
	procPowerWriteACValueIndex     = powrProf.NewProc("PowerWriteACValueIndex")
)

var (
	ultimateGuid = windows.GUID{0xe9a42b02, 0xd5df, 0x448d, [8]byte{0xaa, 0x00, 0x03, 0xf1, 0x47, 0x49, 0xeb, 0x61}}
	balancedGuid = windows.GUID{0x381b4222, 0xf694, 0x41f0, [8]byte{0x96, 0x85, 0xff, 0x5b, 0xb2, 0x60, 0xdf, 0x2e}}
	subProcGuid  = windows.GUID{0x54533251, 0x82be, 0x4824, [8]byte{0x96, 0xc1, 0x47, 0xb6, 0x0b, 0x74, 0x0d, 0x00}}
	cStateGuid   = windows.GUID{0x0cc5b647, 0xc1df, 0x4637, [8]byte{0x89, 0x1a, 0xde, 0xc3, 0x5c, 0x31, 0x85, 0x83}}
)

func getActiveSchemeGuid() *windows.GUID {
	var p *windows.GUID
	ret, _, _ := procPowerGetActiveScheme.Call(0, uintptr(unsafe.Pointer(&p)))
	if ret != 0 {
		return nil
	}
	return p
}

func freeGuid(p *windows.GUID) {
	windows.LocalFree(windows.Handle(unsafe.Pointer(p)))
}

// GetUltimatePerformanceStatus 返回当前是否为卓越性能电源计划
func (a *App) GetUltimatePerformanceStatus() bool {
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

// EnableUltimatePerformance 启用卓越性能电源计划
func (a *App) EnableUltimatePerformance() bool {
	procPowerDuplicatePowerScheme.Call(0, uintptr(unsafe.Pointer(&ultimateGuid)), 0)
	procPowerSetActiveScheme.Call(0, uintptr(unsafe.Pointer(&ultimateGuid)))
	return true
}

// DisableUltimatePerformance 恢复平衡电源计划
func (a *App) DisableUltimatePerformance() bool {
	procPowerSetActiveScheme.Call(0, uintptr(unsafe.Pointer(&balancedGuid)))
	return true
}

func cstateRead(guid *windows.GUID) (uint32, bool) {
	var val uint32
	ret, _, _ := procPowerReadACValueIndex.Call(0, uintptr(unsafe.Pointer(guid)),
		uintptr(unsafe.Pointer(&subProcGuid)), uintptr(unsafe.Pointer(&cStateGuid)), uintptr(unsafe.Pointer(&val)))
	return val, ret == 0
}

func cstateWrite(guid *windows.GUID, val uint32) {
	procPowerWriteACValueIndex.Call(0, uintptr(unsafe.Pointer(guid)),
		uintptr(unsafe.Pointer(&subProcGuid)), uintptr(unsafe.Pointer(&cStateGuid)), uintptr(val))
}

// GetCStateStatus 返回 C-State 当前状态，true=允许（默认），false=禁用
func (a *App) GetCStateStatus() bool {
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

// EnableCState 启用 C-State
func (a *App) EnableCState() bool {
	p := getActiveSchemeGuid()
	if p == nil {
		return false
	}
	defer freeGuid(p)
	cstateWrite(p, 1)
	procPowerSetActiveScheme.Call(0, uintptr(unsafe.Pointer(p)))
	return true
}

// DisableCState 禁用 C-State
func (a *App) DisableCState() bool {
	p := getActiveSchemeGuid()
	if p == nil {
		return false
	}
	defer freeGuid(p)
	cstateWrite(p, 0)
	procPowerSetActiveScheme.Call(0, uintptr(unsafe.Pointer(p)))
	return true
}

// Superfetch - 预加载
func (a *App) GetSuperfetchStatus() bool {
	v, err := regReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management\PrefetchParameters`,
		"EnableSuperfetch")
	return err == nil && v == 3
}
func (a *App) EnableSuperfetch() bool { return regSetDWordBool(registry.LOCAL_MACHINE,
	`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management\PrefetchParameters`,
	"EnableSuperfetch", 3) }
func (a *App) DisableSuperfetch() bool { return regSetDWordBool(registry.LOCAL_MACHINE,
	`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management\PrefetchParameters`,
	"EnableSuperfetch", 0) }

// FullscreenOptimization - 全屏优化
// GameDVR_FSEBehaviorMode: 0=开, 2=关
func (a *App) GetFullscreenOptimizationStatus() bool {
	v, err := regReadDWord(registry.CURRENT_USER,
		`System\GameConfigStore`, "GameDVR_FSEBehaviorMode")
	return err != nil || v == 0
}

func fullscreenEnableForUser(hive registry.Key, path string) {
	regSetDWord(hive, path, "GameDVR_FSEBehaviorMode", 0)
	regDeleteValue(hive, path, "GameDVR_FSEBehavior")
	regDeleteValue(hive, path, "GameDVR_HonorUserFSEBehaviorMode")
	regDeleteValue(hive, path, "GameDVR_DXGIHonorFSEWindowsCompatible")
	regDeleteValue(hive, path, "GameDVR_EFSEFeatureFlags")
}

func fullscreenDisableForUser(hive registry.Key, path string) {
	regSetDWord(hive, path, "GameDVR_FSEBehaviorMode", 2)
	regSetDWord(hive, path, "GameDVR_HonorUserFSEBehaviorMode", 1)
	regSetDWord(hive, path, "GameDVR_DXGIHonorFSEWindowsCompatible", 1)
	regSetDWord(hive, path, "GameDVR_EFSEFeatureFlags", 0)
	regSetDWord(hive, path, "GameDVR_FSEBehavior", 2)
}

func (a *App) EnableFullscreenOptimization() bool {
	fullscreenEnableForUser(registry.CURRENT_USER, `System\GameConfigStore`)
	fullscreenEnableForUser(registry.USERS, `.DEFAULT\System\GameConfigStore`)
	return true
}
func (a *App) DisableFullscreenOptimization() bool {
	fullscreenDisableForUser(registry.CURRENT_USER, `System\GameConfigStore`)
	fullscreenDisableForUser(registry.USERS, `.DEFAULT\System\GameConfigStore`)
	return true
}

// WindowedOptimization - 窗口优化
// DirectXUserGlobalSettings 是一个 REG_SZ 字符串，包含 SwapEffectUpgradeEnable=0/1
func (a *App) GetWindowedOptimizationStatus() bool {
	val, _, err := regReadString(registry.CURRENT_USER,
		`Software\Microsoft\DirectX\UserGpuPreferences`,
		"DirectXUserGlobalSettings")
	if err != nil {
		return true
	}
	return !strings.Contains(val, "SwapEffectUpgradeEnable=0")
}

func windowedSetOpt(hive registry.Key, path, setting string) {
	val, _, err := regReadString(hive, path, "DirectXUserGlobalSettings")
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
	regWriteString(hive, path, "DirectXUserGlobalSettings", newVal)
}

func (a *App) EnableWindowedOptimization() bool {
	windowedSetOpt(registry.CURRENT_USER,
		`Software\Microsoft\DirectX\UserGpuPreferences`,
		"SwapEffectUpgradeEnable=1")
	return true
}
func (a *App) DisableWindowedOptimization() bool {
	windowedSetOpt(registry.CURRENT_USER,
		`Software\Microsoft\DirectX\UserGpuPreferences`,
		"SwapEffectUpgradeEnable=0")
	return true
}

// MPO - 多平面覆盖
func (a *App) GetMpoStatus() bool {
	v, err := regReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\Dwm`, "MPOEnabled")
	if err == nil && v == 0 {
		return false
	}
	v, err = regReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode")
	if err == nil && v == 5 {
		return false
	}
	return true
}
func (a *App) EnableMpo() bool {
	regDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "MPOEnabled")
	regDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode")
	return true
}
func (a *App) DisableMpo() bool {
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "MPOEnabled", 0)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode", 5)
	return true
}

// MemoryCompression - 内存压缩
func (a *App) GetMemoryCompressionStatus() bool {
	v, err := regReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management`,
		"FeatureSettings")
	if err != nil {
		return true
	}
	return (v & 1) == 0 // 第 0 位=0 表示启用压缩
}
func (a *App) EnableMemoryCompression() bool {
	cmd := exec.Command("powershell", "-Command", "Enable-MMAgent -mc")
	hideWindow(cmd)
	return cmd.Run() == nil
}
func (a *App) DisableMemoryCompression() bool {
	cmd := exec.Command("powershell", "-Command", "Disable-MMAgent -mc")
	hideWindow(cmd)
	return cmd.Run() == nil
}

// -------- 个性化 --------

// 通知模式: "开启" "仅关闭通知" "完全关闭"
func (a *App) GetNotificationStatus() string {
	toast, _ := regReadDWord(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled")
	center, _ := regReadDWord(registry.CURRENT_USER,
		`SOFTWARE\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter")
	if center == 1 {
		return "完全关闭"
	}
	if toast == 0 {
		return "仅关闭通知"
	}
	return "开启"
}

func (a *App) SetNotificationMode(mode string) bool {
	switch mode {
	case "开启":
		regDeleteValue(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`, "NoToastApplicationNotification")
		regDeleteValue(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter")
		regDeleteValue(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled")
		regDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled")
	case "仅关闭通知":
		regSetDWord(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`, "NoToastApplicationNotification", 1)
		regSetDWord(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
		regSetDWord(registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
	case "完全关闭":
		regSetDWord(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`, "NoToastApplicationNotification", 1)
		regSetDWord(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
		regSetDWord(registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
		regSetDWord(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter", 1)
	}
	restartExplorer()
	return true
}

// 旧版气球通知
func (a *App) GetLegacyBalloonStatus() bool {
	v, _ := regReadDWord(registry.CURRENT_USER,
		`Software\Policies\Microsoft\Windows\Explorer`, "EnableLegacyBalloonNotifications")
	return v == 1
}
func (a *App) SetLegacyBalloon(on bool) bool {
	if on {
		regSetDWord(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\Explorer`, "EnableLegacyBalloonNotifications", 1)
	} else {
		regSetDWord(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\Explorer`, "EnableLegacyBalloonNotifications", 0)
	}
	restartExplorer()
	return true
}

// 屏幕边缘滑动
func (a *App) GetEdgeSwipeStatus() bool {
	v, _ := regReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Policies\Microsoft\Windows\EdgeUI`, "AllowEdgeSwipe")
	return v != 0
}
func (a *App) SetEdgeSwipe(on bool) bool {
	if on {
		regDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\EdgeUI`, "AllowEdgeSwipe")
	} else {
		regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\EdgeUI`, "AllowEdgeSwipe", 0)
	}
	restartExplorer()
	return true
}

// 新版上下文菜单
func (a *App) GetNewContextMenuStatus() bool {
	v, _, err := regReadString(registry.CURRENT_USER,
		`Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`, "")
	return err != nil || v != ""
}
func (a *App) SetNewContextMenu(on bool) bool {
	if on {
		regDeleteValue(registry.CURRENT_USER,
			`Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`, "")
	} else {
		regWriteString(registry.CURRENT_USER,
			`Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`, "", "")
	}
	restartExplorer()
	return true
}

// Explorer 主页
func (a *App) GetExplorerHomeStatus() bool {
	v, _ := regReadDWord(registry.CURRENT_USER,
		`Software\Classes\CLSID\{f874310e-b6b7-47dc-bc84-b9e6b38f5903}`, "System.IsPinnedToNameSpaceTree")
	return v == 1
}
func (a *App) SetExplorerHome(on bool) bool {
	if on {
		regSetDWord(registry.CURRENT_USER,
			`Software\Classes\CLSID\{f874310e-b6b7-47dc-bc84-b9e6b38f5903}`, "System.IsPinnedToNameSpaceTree", 1)
	} else {
		regSetDWord(registry.CURRENT_USER,
			`Software\Classes\CLSID\{f874310e-b6b7-47dc-bc84-b9e6b38f5903}`, "System.IsPinnedToNameSpaceTree", 0)
	}
	restartExplorer()
	return true
}

// Explorer 图库
func (a *App) GetExplorerGalleryStatus() bool {
	v, _ := regReadDWord(registry.CURRENT_USER,
		`Software\Classes\CLSID\{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}`, "System.IsPinnedToNameSpaceTree")
	return v == 1
}
func (a *App) SetExplorerGallery(on bool) bool {
	if on {
		regSetDWord(registry.CURRENT_USER,
			`Software\Classes\CLSID\{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}`, "System.IsPinnedToNameSpaceTree", 1)
	} else {
		regSetDWord(registry.CURRENT_USER,
			`Software\Classes\CLSID\{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}`, "System.IsPinnedToNameSpaceTree", 0)
	}
	restartExplorer()
	return true
}

// 移除快捷方式小箭头
func (a *App) GetRemoveShortcutArrowStatus() bool {
	v, _, err := regReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "29")
	return err == nil && v != ""
}
func (a *App) SetRemoveShortcutArrow(on bool) bool {
	if on {
		regWriteString(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "29", "")
	} else {
		regDeleteValue(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "29")
	}
	restartExplorer()
	return true
}

// 移除"快捷方式"文字
func (a *App) GetRemoveShortcutTextStatus() bool {
	v, _, err := regReadString(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Explorer\NamingTemplate`, "ShortcutNameTemplate")
	return err == nil && v != ""
}
func (a *App) SetRemoveShortcutText(on bool) bool {
	if on {
		regWriteString(registry.CURRENT_USER,
			`Software\Microsoft\Windows\CurrentVersion\Explorer\NamingTemplate`, "ShortcutNameTemplate", "%s.lnk")
	} else {
		regDeleteValue(registry.CURRENT_USER,
			`Software\Microsoft\Windows\CurrentVersion\Explorer\NamingTemplate`, "ShortcutNameTemplate")
	}
	restartExplorer()
	return true
}

// 移除可执行文件小盾牌
func (a *App) GetRemoveShieldStatus() bool {
	v, _, err := regReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "77")
	return err == nil && v != ""
}
func (a *App) SetRemoveShield(on bool) bool {
	if on {
		regWriteString(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "77",
			`%SystemRoot%\System32\shell32.dll,-278`)
	} else {
		regDeleteValue(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "77")
	}
	restartExplorer()
	return true
}

func restartExplorer() {
	exec.Command("taskkill", "/f", "/im", "explorer.exe").Run()
	time.Sleep(1500 * time.Millisecond)

	dll := syscall.MustLoadDLL("shell32.dll")
	proc := dll.MustFindProc("ShellExecuteW")
	proc.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("open"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("explorer.exe"))),
		0, 0,
		5,
	)
}

// Windows 照片查看器
var imgExts = []string{".jpg", ".jpeg", ".png", ".bmp", ".gif", ".tiff", ".tif", ".ico", ".webp"}

func (a *App) GetPhotoViewerStatus() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\.jpg\OpenWithProgids`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	vals, err := k.ReadValueNames(0)
	if err != nil {
		return false
	}
	for _, v := range vals {
		if strings.Contains(v, "PhotoViewer") {
			return true
		}
	}
	return false
}

func (a *App) EnablePhotoViewer() bool {
	for _, ext := range imgExts {
		k, _, _ := registry.CreateKey(registry.CURRENT_USER,
			`Software\Classes\`+ext+`\OpenWithProgids`, registry.SET_VALUE)
		k.SetStringValue("PhotoViewer.FileAssoc.Tiff", "")
		k.Close()
	}
	return true
}

func (a *App) DisablePhotoViewer() bool {
	entries := []string{
		".jpg", ".jpeg", ".png", ".bmp", ".gif",
		".tiff", ".tif", ".ico", ".webp",
	}
	for _, ext := range entries {
		k, err := registry.OpenKey(registry.CURRENT_USER,
			`Software\Classes\`+ext+`\OpenWithProgids`, registry.SET_VALUE)
		if err != nil {
			continue
		}
		k.DeleteValue("PhotoViewer.FileAssoc.Tiff")
		k.Close()
	}
	return true
}

// 增强版卸载 Edge（不删 EdgeUpdate / EdgeCore，不影响 WebView2）
func (a *App) UninstallEdge() string {
	// 杀 Edge 进程
	exec.Command("taskkill", "/f", "/im", "msedge.exe").Run()

	// 找 setup.exe
	matches, err := filepath.Glob(`C:\Program Files (x86)\Microsoft\Edge\Application\*\Installer\setup.exe`)
	if err != nil || len(matches) == 0 {
		return "未找到 Edge 安装路径"
	}

	// 卸载
	cmd := exec.Command(matches[0], "--uninstall", "--force-uninstall", "--system-level")
	hideWindow(cmd)
	cmd.Run()

	// 删 Edge 应用目录（不删 EdgeCore）
	edgeAppDir := filepath.Dir(filepath.Dir(matches[0])) // Application\版本号\Installer\setup.exe → Application\版本号
	os.RemoveAll(edgeAppDir)

	// 删快捷方式
	shortcuts := []string{
		os.Getenv("PUBLIC") + `\Desktop\Microsoft Edge.lnk`,
		os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs\Microsoft Edge.lnk`,
	}
	for _, s := range shortcuts {
		os.Remove(s)
	}

	return "Edge 已卸载（保留 EdgeUpdate/EdgeCore，WebView2 不受影响）"
}

// WebView2 版本
func (a *App) GetWebView2Version() string {
	guids := []string{
		`Software\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`,
		`SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`,
	}
	for _, g := range guids {
		v, _, err := regReadString(registry.LOCAL_MACHINE, g, "pv")
		if err == nil && v != "" {
			return v
		}
	}
	return "未安装"
}

// 安装/升级 WebView2（Evergreen Bootstrapper）
func (a *App) InstallWebView2() string {
	tmp := os.TempDir() + `\MicrosoftEdgeWebView2Setup.exe`
	url := "https://go.microsoft.com/fwlink/p/?LinkId=2124703"

	out, err := http.Get(url)
	if err != nil {
		return "下载失败: " + err.Error()
	}
	defer out.Body.Close()

	f, err := os.Create(tmp)
	if err != nil {
		return "创建文件失败: " + err.Error()
	}
	io.Copy(f, out.Body)
	f.Close()

	cmd := exec.Command(tmp, "/silent", "/install")
	hideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return "安装失败: " + err.Error()
	}
	return "WebView2 安装/升级完成"
}

// -------- 电源管理 --------

var powrProf = syscall.NewLazyDLL("powrprof.dll")
var procCallNtPower = powrProf.NewProc("CallNtPowerInformation")

const systemReserveHiberFile = 10

func (a *App) GetHibernateStatus() bool {
	v, err := regReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Power`, "HibernateEnabled")
	return err == nil && v == 1
}

func (a *App) EnableHibernate() bool {
	var v byte = 1
	ret, _, _ := procCallNtPower.Call(systemReserveHiberFile,
		uintptr(unsafe.Pointer(&v)), 1, 0, 0)
	if ret != 0 {
		return false
	}
	regSetDWord(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Power`, "HibernateEnabled", 1)
	return true
}

func (a *App) DisableHibernate() bool {
	var v byte = 0
	ret, _, _ := procCallNtPower.Call(systemReserveHiberFile,
		uintptr(unsafe.Pointer(&v)), 1, 0, 0)
	if ret != 0 {
		return false
	}
	regSetDWord(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Power`, "HibernateEnabled", 0)
	return true
}

// FastStartup
func (a *App) GetFastStartupStatus() bool {
	v, err := regReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Power`, "HiberbootEnabled")
	return err == nil && v == 1
}

func (a *App) EnableFastStartup() bool {
	regSetDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Power`, "HiberbootEnabled", 1)
	return true
}

func (a *App) DisableFastStartup() bool {
	regSetDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Power`, "HiberbootEnabled", 0)
	return true
}

// -------- 安全模块 --------

func scHide(args ...string) {
	cmd := exec.Command("sc", args...)
	hideWindow(cmd)
	cmd.Run()
}

func runHide(name string, args ...string) {
	cmd := exec.Command(name, args...)
	hideWindow(cmd)
	cmd.Run()
}

func (a *App) GetDefenderStatus() bool {
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

	v, err := regReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows Defender`, "DisableAntiSpyware")
	if err == nil && v == 1 {
		return false
	}
	return true
}

func (a *App) EnableDefender() bool {
	// 显式设为 0（启用），比删除更可靠
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`, "DisableAntiSpyware", 0)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`, "DisableAntiVirus", 0)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows Defender`, "DisableAntiSpyware", 0)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows Defender`, "DisableAntiVirus", 0)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender\Features`, "TamperProtection", 0)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender\Features`, "TamperProtectionSource", 0)

	// 刷新策略
	runHide("gpupdate", "/target:computer", "/force")

	// 启用服务
	scHide("config", "WinDefend", "start=", "auto")
	scHide("start", "WinDefend")
	scHide("config", "SecurityHealthService", "start=", "auto")
	scHide("start", "SecurityHealthService")
	scHide("config", "MDCoreSvc", "start=", "auto")
	scHide("config", "wscsvc", "start=", "auto")
	scHide("start", "wscsvc")
	return true
}

func (a *App) DisableDefender() bool {
	// 关篡改保护
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender\Features`, "TamperProtection", 0)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender\Features`, "TamperProtectionSource", 0)

	// 写禁用标记
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`, "DisableAntiSpyware", 1)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows Defender`, "DisableAntiVirus", 1)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows Defender`, "DisableAntiSpyware", 1)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows Defender`, "DisableAntiVirus", 1)

	// 刷新策略
	runHide("gpupdate", "/target:computer", "/force")

	// 停服务
	scHide("stop", "WinDefend")
	scHide("config", "WinDefend", "start=", "disabled")
	scHide("stop", "SecurityHealthService")
	scHide("config", "SecurityHealthService", "start=", "disabled")
	scHide("stop", "MDCoreSvc")
	scHide("config", "MDCoreSvc", "start=", "disabled")
	scHide("stop", "wscsvc")
	scHide("config", "wscsvc", "start=", "disabled")
	return true
}

// UAC

func (a *App) GetUacStatus() bool {
	v, err := regReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableLUA")
	return err == nil && v == 1
}

func (a *App) EnableUac() bool {
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableLUA", 1)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableVirtualization", 1)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableInstallerDetection", 1)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "PromptOnSecureDesktop", 1)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "ConsentPromptBehaviorAdmin", 2)
	return true
}

func (a *App) DisableUac() bool {
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "EnableLUA", 0)
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "ConsentPromptBehaviorAdmin", 0)
	return true
}

// VBS

func (a *App) GetVbsStatus() bool {
	v, err := regReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Policies\Microsoft\Windows\DeviceGuard`, "EnableVirtualizationBasedSecurity")
	if err == nil && v == 1 {
		return true
	}
	v, err = regReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\ControlSet001\Control\DeviceGuard`, "EnableVirtualizationBasedSecurity")
	return err == nil && v == 1
}

func (a *App) EnableVbs() bool {
	regSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\DeviceGuard`, "EnableVirtualizationBasedSecurity", 1)
	regSetDWord(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Control\DeviceGuard`, "EnableVirtualizationBasedSecurity", 1)
	return true
}

func (a *App) DisableVbs() bool {
	regDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\DeviceGuard`, "EnableVirtualizationBasedSecurity")
	regDeleteValue(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Control\DeviceGuard`, "EnableVirtualizationBasedSecurity")
	return true
}

// MemoryIntegrity

func (a *App) GetMemoryIntegrityStatus() bool {
	v, err := regReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\ControlSet001\Control\DeviceGuard\Scenarios\HypervisorEnforcedCodeIntegrity`, "Enabled")
	return err == nil && v == 1
}

func (a *App) EnableMemoryIntegrity() bool {
	regSetDWord(registry.LOCAL_MACHINE,
		`SYSTEM\ControlSet001\Control\DeviceGuard\Scenarios\HypervisorEnforcedCodeIntegrity`, "Enabled", 1)
	return true
}

func (a *App) DisableMemoryIntegrity() bool {
	regSetDWord(registry.LOCAL_MACHINE,
		`SYSTEM\ControlSet001\Control\DeviceGuard\Scenarios\HypervisorEnforcedCodeIntegrity`, "Enabled", 0)
	return true
}

// 注册表辅助（返回 bool）
func regSetDWordBool(key registry.Key, path, name string, value uint32) bool {
	regSetDWord(key, path, name, value)
	return true
}
func regDeleteIfExists(key registry.Key, path, name string) bool {
	regDeleteValue(key, path, name)
	return true
}


