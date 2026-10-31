//go:build windows

package services

import (
	"OptiWin/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

// reportErr 记录操作失败并返回是否成功，供不带前缀的开关方法汇总写入结果
func reportErr(op string, err error) bool {
	if err != nil {
		utils.Logf("[services] %s 失败: %v", op, err)
		utils.SetLastError(fmt.Sprintf("%s: %v", op, err))
		return false
	}
	return true
}

func dwordOn(hive registry.Key, path, name string) bool {
	v, err := utils.RegReadDWord(hive, path, name)
	return err == nil && v == 1
}

// setDWordToggle 通用开关：目标状态与当前一致则跳过，写入成功后重启 Explorer；op 用于失败日志
func setDWordToggle(op string, hive registry.Key, path, name string, on bool) bool {
	if dwordOn(hive, path, name) == on {
		return true
	}
	var v uint32
	if on {
		v = 1
	}
	if err := utils.RegSetDWordE(hive, path, name, v); err != nil {
		return reportErr(op, err)
	}
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
	const op = "SetNotificationMode"
	var err error
	switch mode {
	case NotificationOn:
		err = errors.Join(
			deleteIgnoringMissing(registry.CURRENT_USER, regPathPushNotificationsPol, regValueNoToastAppNotif),
			deleteIgnoringMissing(registry.CURRENT_USER, regPathExplorerPol, regValueDisableNotifCenter),
			utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotifications, regValueToastEnabled, 1),
			utils.RegSetDWordE(registry.LOCAL_MACHINE, regPathPushNotifications, regValueToastEnabled, 1),
		)
	case NotificationToastOff:
		// 从"完全关闭"降级时必须清除 DisableNotificationCenter
		err = errors.Join(
			deleteIgnoringMissing(registry.CURRENT_USER, regPathExplorerPol, regValueDisableNotifCenter),
			utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotificationsPol, regValueNoToastAppNotif, 1),
			utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotifications, regValueToastEnabled, 0),
			utils.RegSetDWordE(registry.LOCAL_MACHINE, regPathPushNotifications, regValueToastEnabled, 0),
		)
	case NotificationOff:
		err = errors.Join(
			utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotificationsPol, regValueNoToastAppNotif, 1),
			utils.RegSetDWordE(registry.CURRENT_USER, regPathPushNotifications, regValueToastEnabled, 0),
			utils.RegSetDWordE(registry.LOCAL_MACHINE, regPathPushNotifications, regValueToastEnabled, 0),
			utils.RegSetDWordE(registry.CURRENT_USER, regPathExplorerPol, regValueDisableNotifCenter, 1),
		)
	default:
		return reportErr(op, fmt.Errorf("未知通知模式 %q", mode))
	}
	if !reportErr(op, err) {
		return false
	}
	return true
}

func GetLegacyBalloonStatus() bool {
	return dwordOn(registry.CURRENT_USER, regPathExplorerPol, regValueLegacyBalloon)
}

func SetLegacyBalloon(on bool) bool {
	return setDWordToggle("SetLegacyBalloon", registry.CURRENT_USER, regPathExplorerPol, regValueLegacyBalloon, on)
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
			return reportErr("SetEdgeSwipe", err)
		}
	} else {
		if err := utils.RegSetDWordE(registry.LOCAL_MACHINE, regPathEdgeUIPol, regValueAllowEdgeSwipe, 0); err != nil {
			return reportErr("SetEdgeSwipe", err)
		}
	}
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
			return reportErr("SetNewContextMenu", err)
		}
	} else {
		if err := utils.RegWriteStringE(registry.CURRENT_USER, regPathOldMenuCLSID, "", ""); err != nil {
			return reportErr("SetNewContextMenu", err)
		}
	}
	return true
}

func GetExplorerHomeStatus() bool {
	return dwordOn(registry.CURRENT_USER, regPathHomeCLSID, regValuePinnedToNS)
}

func SetExplorerHome(on bool) bool {
	return setDWordToggle("SetExplorerHome", registry.CURRENT_USER, regPathHomeCLSID, regValuePinnedToNS, on)
}

func GetExplorerGalleryStatus() bool {
	return dwordOn(registry.CURRENT_USER, regPathGalleryCLSID, regValuePinnedToNS)
}

func SetExplorerGallery(on bool) bool {
	return setDWordToggle("SetExplorerGallery", registry.CURRENT_USER, regPathGalleryCLSID, regValuePinnedToNS, on)
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
			return reportErr("SetRemoveShortcutArrow", err)
		}
	} else {
		if _, _, err := utils.RegReadString(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconArrow); err != nil {
			return true
		}
		if err := deleteIgnoringMissing(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconArrow); err != nil {
			return reportErr("SetRemoveShortcutArrow", err)
		}
	}
	// 清除图标缓存，避免残留损坏的缓存条目导致图标黑块或不刷新
	utils.ExecuteString(`Remove-Item "$env:LOCALAPPDATA\IconCache.db","$env:LOCALAPPDATA\Microsoft\Windows\Explorer\iconcache_*" -Force -ErrorAction SilentlyContinue`)
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
			return reportErr("SetRemoveShortcutText", err)
		}
	} else {
		if err := deleteIgnoringMissing(registry.CURRENT_USER, regPathNamingTemplate, regValueShortcutNameTpl); err != nil {
			return reportErr("SetRemoveShortcutText", err)
		}
	}
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
			return reportErr("SetRemoveShield", err)
		}
	} else {
		if _, _, err := utils.RegReadString(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconShield); err != nil {
			return true
		}
		if err := deleteIgnoringMissing(registry.LOCAL_MACHINE, regPathShellIcons, regValueShellIconShield); err != nil {
			return reportErr("SetRemoveShield", err)
		}
	}
	// 清除图标缓存，避免残留损坏的缓存条目导致图标黑块或不刷新
	utils.ExecuteString(`Remove-Item "$env:LOCALAPPDATA\IconCache.db","$env:LOCALAPPDATA\Microsoft\Windows\Explorer\iconcache_*" -Force -ErrorAction SilentlyContinue`)
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
	source, target := taskManagerPath, taskManagerPath+".bak"
	rename := func(from, to string) bool {
		cmd := fmt.Sprintf(`Rename-Item -LiteralPath '%s' -NewName '%s' -Force`, from, filepath.Base(to))
		return utils.SuperExecuteString(cmd)
	}
	if enable {
		source, target = target, source
	}

	if !utils.ValidPSArg(source) || !utils.ValidPSArg(target) {
		utils.Logf("[services] SetOldTaskManager 拒绝执行：路径含非法字符")
		return false
	}

	if !rename(source, target) {
		utils.Logf("[services] SetOldTaskManager 重命名失败: %s -> %s", source, target)
		return false
	}

	// TrustedInstaller 提权可能较慢，轮询最多 6 秒确认重命名生效
	if waitFor(6*time.Second, func() bool {
		_, err := os.Stat(source)
		return os.IsNotExist(err)
	}) {
		return true
	}

	// 超时后重新读取实际状态，已生效就不报失败（否则用户再点一次会把结果改回去）
	if GetOldTaskManagerStatus() == enable {
		return true
	}

	utils.Logf("[services] SetOldTaskManager 超时且未生效，尝试回滚")
	rename(target, source) // 失败或被打断时反向重命名，避免留下半成品
	return false
}

// widgetsPolicyPath Windows 小组件组策略键（AllowNewsAndInterests=0 禁用）
const widgetsPolicyPath = `SOFTWARE\Policies\Microsoft\Dsh`

// GetWidgetsStatus 检测小组件是否启用：AllowNewsAndInterests=0(已禁用) 返回 false
func GetWidgetsStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE, widgetsPolicyPath, "AllowNewsAndInterests")
	return err != nil || v != 0
}

// DisableWidgets 禁用 Windows 小组件：写组策略键
func DisableWidgets() bool {
	return reportErr("DisableWidgets", utils.RegSetDWordE(registry.LOCAL_MACHINE, widgetsPolicyPath, "AllowNewsAndInterests", 0))
}

// EnableWidgets 恢复 Windows 小组件：删除组策略键（键本来不存在也算成功）
func EnableWidgets() bool {
	return reportErr("EnableWidgets", deleteIgnoringMissing(registry.LOCAL_MACHINE, widgetsPolicyPath, "AllowNewsAndInterests"))
}
