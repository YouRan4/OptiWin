//go:build windows

package services

import (
	"OptiWin/utils"
	"errors"
	"os"
	"time"

	"golang.org/x/sys/windows/registry"
)

// 通知模式（GetNotificationStatus / SetNotificationMode 的三态取值）
const (
	NotificationOn       = "0" // 开启通知
	NotificationToastOff = "1" // 仅关闭通知（保留操作中心）
	NotificationOff      = "2" // 完全关闭（含操作中心）
)

// 注册表路径常量
const (
	// 通知：ToastEnabled 主设置所在键（HKCU/HKLM 双写，控制 toast 开关）
	regPathPushNotifications = `Software\Microsoft\Windows\CurrentVersion\PushNotifications`
	// 通知：NoToastApplicationNotification 策略键，通过策略强制禁用应用 toast
	regPathPushNotificationsPol = `Software\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`
	// 通知/气球：Explorer 策略键，控制操作中心（DisableNotificationCenter）与旧版气球通知
	regPathExplorerPol = `SOFTWARE\Policies\Microsoft\Windows\Explorer`
	// 边缘滑动：EdgeUI 策略键，控制触屏边缘滑动手势（AllowEdgeSwipe）
	regPathEdgeUIPol = `SOFTWARE\Policies\Microsoft\Windows\EdgeUI`
	// 快捷方式箭头/盾牌：Shell 图标覆盖槽位所在键（数字值名 = 图标槽位号）
	regPathShellIcons = `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`
	// 快捷方式文本：新建快捷方式的命名模板
	regPathNamingTemplate = `Software\Microsoft\Windows\CurrentVersion\Explorer\NamingTemplate`
	// 右键菜单：该 CLSID 的 InprocServer32 默认值存在即启用 Win10 旧式右键菜单
	regPathOldMenuCLSID = `Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`
	// 主页：Home 项是否固定到资源管理器导航栏
	regPathHomeCLSID = `Software\Classes\CLSID\{f874310e-b6b7-47dc-bc84-b9e6b38f5903}`
	// 图库：Gallery 项是否固定到资源管理器导航栏
	regPathGalleryCLSID = `Software\Classes\CLSID\{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}`
)

// 注册表值名常量
const (
	// 通知：0=关闭 toast，1=开启
	regValueToastEnabled = "ToastEnabled"
	// 通知：1=禁止应用发送 toast 通知（策略）
	regValueNoToastAppNotif = "NoToastApplicationNotification"
	// 通知：1=禁用操作中心（通知中心）
	regValueDisableNotifCenter = "DisableNotificationCenter"
	// 气球通知：1=启用旧版气泡通知
	regValueLegacyBalloon = "EnableLegacyBalloonNotifications"
	// 边缘滑动：0=禁止边缘滑动；删除该值 = 恢复默认（允许）
	regValueAllowEdgeSwipe = "AllowEdgeSwipe"
	// 主页/图库：1=固定到导航栏
	regValuePinnedToNS = "System.IsPinnedToNameSpaceTree"
	// 快捷方式文本："%s.lnk" = 去掉新建快捷方式名称中的" - 快捷方式"后缀
	regValueShortcutNameTpl = "ShortcutNameTemplate"
	// 快捷方式箭头：29 号图标槽位，指向透明图标时隐藏箭头
	regValueShellIconArrow = "29"
	// 盾牌图标：77 号图标槽位，替换可执行文件上的 UAC 盾牌角标
	regValueShellIconShield = "77"
)

// 注册表写入值常量（value data，写入到指定值名的具体内容）
const (
	// 快捷方式箭头：imageres.dll 的 197 号空白图标，用于隐藏箭头
	regDataArrowOverlay = `%SystemRoot%\System32\imageres.dll,197`
	// 快捷方式文本：去掉新建快捷方式名称中的" - 快捷方式"后缀
	regDataShortcutNameTemplate = "%s.lnk"
	// 盾牌图标：imageres.dll 的 197 号空白图标，用于隐藏 UAC 盾牌角标
	regDataShieldIcon = `%SystemRoot%\System32\imageres.dll,197`
)

func deleteIgnoringMissing(hive registry.Key, path, name string) error {
	err := utils.RegDeleteValueE(hive, path, name)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	}
	return err
}

func dwordOn(hive registry.Key, path, name string) bool {
	v, err := utils.RegReadDWord(hive, path, name)
	return err == nil && v == 1
}

// setDWordToggle 通用开关：目标状态与当前一致则跳过，写入成功后重启 Explorer
func setDWordToggle(hive registry.Key, path, name string, on bool) bool {
	if dwordOn(hive, path, name) == on {
		return true
	}
	var v uint32
	if on {
		v = 1
	}
	if err := utils.RegSetDWordE(hive, path, name, v); err != nil {
		return false
	}
	utils.RestartExplorer()
	return true
}

// GetNotificationStatus 返回当前通知模式
// 返回值: "0"=开启, "1"=仅关闭通知, "2"=完全关闭
func GetNotificationStatus() string {
	toast, toastErr := utils.RegReadDWord(registry.CURRENT_USER, regPathPushNotifications, regValueToastEnabled)
	center, centerErr := utils.RegReadDWord(registry.CURRENT_USER, regPathExplorerPol, regValueDisableNotifCenter)
	if centerErr == nil && center == 1 {
		return NotificationOff
	}
	if toastErr == nil && toast == 0 {
		return NotificationToastOff
	}
	return NotificationOn
}

// SetNotificationMode 设置通知模式
// mode: "0"=开启, "1"=仅关闭通知, "2"=完全关闭
func SetNotificationMode(mode string) bool {
	before := GetNotificationStatus()

	switch mode {
	case NotificationOn:
		if err := deleteIgnoringMissing(registry.CURRENT_USER, regPathPushNotificationsPol, regValueNoToastAppNotif); err != nil {
			return false
		}
		if err := deleteIgnoringMissing(registry.CURRENT_USER, regPathExplorerPol, regValueDisableNotifCenter); err != nil {
			return false
		}
		if err := utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotifications, regValueToastEnabled, 1); err != nil {
			return false
		}
		if err := utils.RegSetDWordE(registry.LOCAL_MACHINE, regPathPushNotifications, regValueToastEnabled, 1); err != nil {
			return false
		}
	case NotificationToastOff:
		// 从"完全关闭"降级时必须清除 DisableNotificationCenter
		if err := deleteIgnoringMissing(registry.CURRENT_USER, regPathExplorerPol, regValueDisableNotifCenter); err != nil {
			return false
		}
		if err := utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotificationsPol, regValueNoToastAppNotif, 1); err != nil {
			return false
		}
		if err := utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotifications, regValueToastEnabled, 0); err != nil {
			return false
		}
		if err := utils.RegSetDWordE(registry.LOCAL_MACHINE, regPathPushNotifications, regValueToastEnabled, 0); err != nil {
			return false
		}
	case NotificationOff:
		if err := utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotificationsPol, regValueNoToastAppNotif, 1); err != nil {
			return false
		}
		if err := utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotifications, regValueToastEnabled, 0); err != nil {
			return false
		}
		if err := utils.RegSetDWordE(registry.LOCAL_MACHINE, regPathPushNotifications, regValueToastEnabled, 0); err != nil {
			return false
		}
		if err := utils.RegSetDWordE(registry.CURRENT_USER, regPathExplorerPol, regValueDisableNotifCenter, 1); err != nil {
			return false
		}
	default:
		return false
	}

	// 状态实际变化才重启 Explorer，避免无谓闪屏
	if GetNotificationStatus() != before {
		utils.RestartExplorer()
	}
	return true
}

func GetLegacyBalloonStatus() bool {
	return dwordOn(registry.CURRENT_USER, regPathExplorerPol, regValueLegacyBalloon)
}

func SetLegacyBalloon(on bool) bool {
	return setDWordToggle(registry.CURRENT_USER, regPathExplorerPol, regValueLegacyBalloon, on)
}

func GetEdgeSwipeStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE, regPathEdgeUIPol, regValueAllowEdgeSwipe)
	// 键不存在 = 未配置策略 = 默认允许滑动
	return err != nil || v != 0
}

func SetEdgeSwipe(on bool) bool {
	if GetEdgeSwipeStatus() == on {
		return true
	}
	if on {
		if err := deleteIgnoringMissing(registry.LOCAL_MACHINE, regPathEdgeUIPol, regValueAllowEdgeSwipe); err != nil {
			return false
		}
	} else {
		if err := utils.RegSetDWordE(registry.LOCAL_MACHINE, regPathEdgeUIPol, regValueAllowEdgeSwipe, 0); err != nil {
			return false
		}
	}
	utils.RestartExplorer()
	return true
}

func GetNewContextMenuStatus() bool {
	v, _, err := utils.RegReadString(registry.CURRENT_USER, regPathOldMenuCLSID, "")
	// 键不存在或值非空 = 新式右键菜单
	return err != nil || v != ""
}

func SetNewContextMenu(on bool) bool {
	if GetNewContextMenuStatus() == on {
		return true
	}
	if on {
		if err := deleteIgnoringMissing(registry.CURRENT_USER, regPathOldMenuCLSID, ""); err != nil {
			return false
		}
	} else {
		if err := utils.RegWriteStringE(registry.CURRENT_USER, regPathOldMenuCLSID, "", ""); err != nil {
			return false
		}
	}
	utils.RestartExplorer()
	return true
}

func GetExplorerHomeStatus() bool {
	return dwordOn(registry.CURRENT_USER, regPathHomeCLSID, regValuePinnedToNS)
}

func SetExplorerHome(on bool) bool {
	return setDWordToggle(registry.CURRENT_USER, regPathHomeCLSID, regValuePinnedToNS, on)
}

func GetExplorerGalleryStatus() bool {
	return dwordOn(registry.CURRENT_USER, regPathGalleryCLSID, regValuePinnedToNS)
}

func SetExplorerGallery(on bool) bool {
	return setDWordToggle(registry.CURRENT_USER, regPathGalleryCLSID, regValuePinnedToNS, on)
}

func GetRemoveShortcutArrowStatus() bool {
	v, _, err := utils.RegReadString(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconArrow)
	// 空字符串是无效残留（会导致图标渲染异常），视为未正确移除
	return err == nil && v != ""
}

func SetRemoveShortcutArrow(on bool) bool {
	if on {
		if v, _, err := utils.RegReadString(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconArrow); err == nil && v == regDataArrowOverlay {
			return true
		}
		if err := utils.RegWriteStringE(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconArrow, regDataArrowOverlay); err != nil {
			return false
		}
	} else {
		if _, _, err := utils.RegReadString(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconArrow); err != nil {
			return true
		}
		if err := deleteIgnoringMissing(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconArrow); err != nil {
			return false
		}
	}
	// 清除图标缓存，避免残留损坏的缓存条目导致图标黑块或不刷新
	utils.ExecuteString(`Remove-Item "$env:LOCALAPPDATA\IconCache.db","$env:LOCALAPPDATA\Microsoft\Windows\Explorer\iconcache_*" -Force -ErrorAction SilentlyContinue`)
	utils.RestartExplorer()
	return true
}

func GetRemoveShortcutTextStatus() bool {
	v, _, err := utils.RegReadString(registry.CURRENT_USER, regPathNamingTemplate, regValueShortcutNameTpl)
	return err == nil && v != ""
}

func SetRemoveShortcutText(on bool) bool {
	if GetRemoveShortcutTextStatus() == on {
		return true
	}
	if on {
		if err := utils.RegWriteStringE(registry.CURRENT_USER, regPathNamingTemplate, regValueShortcutNameTpl, regDataShortcutNameTemplate); err != nil {
			return false
		}
	} else {
		if err := deleteIgnoringMissing(registry.CURRENT_USER, regPathNamingTemplate, regValueShortcutNameTpl); err != nil {
			return false
		}
	}
	utils.RestartExplorer()
	return true
}

func GetRemoveShieldStatus() bool {
	v, _, err := utils.RegReadString(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconShield)
	return err == nil && v != ""
}

func SetRemoveShield(on bool) bool {
	if on {
		if v, _, err := utils.RegReadString(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconShield); err == nil && v == regDataShieldIcon {
			return true
		}
		if err := utils.RegWriteStringE(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconShield, regDataShieldIcon); err != nil {
			return false
		}
	} else {
		if _, _, err := utils.RegReadString(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconShield); err != nil {
			return true
		}
		if err := deleteIgnoringMissing(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconShield); err != nil {
			return false
		}
	}
	// 清除图标缓存，避免残留损坏的缓存条目导致图标黑块或不刷新
	utils.ExecuteString(`Remove-Item "$env:LOCALAPPDATA\IconCache.db","$env:LOCALAPPDATA\Microsoft\Windows\Explorer\iconcache_*" -Force -ErrorAction SilentlyContinue`)
	utils.RestartExplorer()
	return true
}

const taskManagerPath = `C:\Windows\SystemResources\Windows.UI.TaskManager`

func GetOldTaskManagerStatus() bool {
	_, err := os.Stat(taskManagerPath)
	if err == nil {
		return true // 文件夹存在 = Win11
	}
	// 原路径找不到，检查 .bak 是否存在（可能权限错误）
	_, err = os.Stat(taskManagerPath + ".bak")
	if err == nil {
		return false // 已重命名 = Win10
	}
	return true // 两边都查不到，默认 Win11
}

func SetOldTaskManager(enable bool) bool {
	if GetOldTaskManagerStatus() == enable {
		return true
	}
	var cmd string
	if enable {
		cmd = `Rename-Item -Path "C:\Windows\SystemResources\Windows.UI.TaskManager.bak" -NewName "Windows.UI.TaskManager" -Force`
	} else {
		cmd = `Rename-Item -Path "C:\Windows\SystemResources\Windows.UI.TaskManager" -NewName "Windows.UI.TaskManager.bak" -Force`
	}

	if !utils.SuperExecuteString(cmd) {
		return false
	}

	// TrustedInstaller 提权可能较慢，轮询最多 6 秒确认重命名生效
	source := taskManagerPath
	if enable {
		source = taskManagerPath + ".bak"
	}
	success := false
	for i := 0; i < 20; i++ {
		time.Sleep(300 * time.Millisecond)
		if _, err := os.Stat(source); os.IsNotExist(err) {
			success = true
			break
		}
	}
	utils.RestartExplorer()
	return success
}
