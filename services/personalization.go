//go:build windows

package services

import (
	"OptiWin/utils"
	"os"
	"time"

	"golang.org/x/sys/windows/registry"
)

// GetNotificationStatus 返回当前通知模式
// 返回值: "0"=开启, "1"=仅关闭通知, "2"=完全关闭
func GetNotificationStatus() string {
	toast, toastErr := utils.RegReadDWord(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled")
	center, centerErr := utils.RegReadDWord(registry.CURRENT_USER,
		`SOFTWARE\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter")
	if centerErr == nil && center == 1 {
		return "2"
	}
	if toastErr == nil && toast == 0 {
		return "1"
	}
	return "0"
}

// SetNotificationMode 设置通知模式
// mode: "0"=开启, "1"=仅关闭通知, "2"=完全关闭
func SetNotificationMode(mode string) bool {
	switch mode {
	case "0":
		utils.RegDeleteValue(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`, "NoToastApplicationNotification")
		utils.RegDeleteValue(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter")
		utils.RegSetDWord(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 1)
		utils.RegSetDWord(registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 1)
	case "1":
		utils.RegSetDWord(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`, "NoToastApplicationNotification", 1)
		utils.RegSetDWord(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
		utils.RegSetDWord(registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
	case "2":
		utils.RegSetDWord(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`, "NoToastApplicationNotification", 1)
		utils.RegSetDWord(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
		utils.RegSetDWord(registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
		utils.RegSetDWord(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter", 1)
	}
	utils.RestartExplorer()
	return true
}

func GetLegacyBalloonStatus() bool {
	v, _ := utils.RegReadDWord(registry.CURRENT_USER,
		`Software\Policies\Microsoft\Windows\Explorer`, "EnableLegacyBalloonNotifications")
	return v == 1
}

func SetLegacyBalloon(on bool) bool {
	if on {
		utils.RegSetDWord(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\Explorer`, "EnableLegacyBalloonNotifications", 1)
	} else {
		utils.RegSetDWord(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\Explorer`, "EnableLegacyBalloonNotifications", 0)
	}
	utils.RestartExplorer()
	return true
}

func GetEdgeSwipeStatus() bool {
	v, _ := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Policies\Microsoft\Windows\EdgeUI`, "AllowEdgeSwipe")
	return v != 0
}

func SetEdgeSwipe(on bool) bool {
	if on {
		utils.RegDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\EdgeUI`, "AllowEdgeSwipe")
	} else {
		utils.RegSetDWord(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Windows\EdgeUI`, "AllowEdgeSwipe", 0)
	}
	utils.RestartExplorer()
	return true
}

func GetNewContextMenuStatus() bool {
	v, _, err := utils.RegReadString(registry.CURRENT_USER,
		`Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`, "")
	return err != nil || v != ""
}

func SetNewContextMenu(on bool) bool {
	if on {
		utils.RegDeleteValue(registry.CURRENT_USER,
			`Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`, "")
	} else {
		utils.RegWriteString(registry.CURRENT_USER,
			`Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`, "", "")
	}
	utils.RestartExplorer()
	return true
}

func GetExplorerHomeStatus() bool {
	v, _ := utils.RegReadDWord(registry.CURRENT_USER,
		`Software\Classes\CLSID\{f874310e-b6b7-47dc-bc84-b9e6b38f5903}`, "System.IsPinnedToNameSpaceTree")
	return v == 1
}

func SetExplorerHome(on bool) bool {
	if on {
		utils.RegSetDWord(registry.CURRENT_USER,
			`Software\Classes\CLSID\{f874310e-b6b7-47dc-bc84-b9e6b38f5903}`, "System.IsPinnedToNameSpaceTree", 1)
	} else {
		utils.RegSetDWord(registry.CURRENT_USER,
			`Software\Classes\CLSID\{f874310e-b6b7-47dc-bc84-b9e6b38f5903}`, "System.IsPinnedToNameSpaceTree", 0)
	}
	utils.RestartExplorer()
	return true
}

func GetExplorerGalleryStatus() bool {
	v, _ := utils.RegReadDWord(registry.CURRENT_USER,
		`Software\Classes\CLSID\{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}`, "System.IsPinnedToNameSpaceTree")
	return v == 1
}

func SetExplorerGallery(on bool) bool {
	if on {
		utils.RegSetDWord(registry.CURRENT_USER,
			`Software\Classes\CLSID\{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}`, "System.IsPinnedToNameSpaceTree", 1)
	} else {
		utils.RegSetDWord(registry.CURRENT_USER,
			`Software\Classes\CLSID\{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}`, "System.IsPinnedToNameSpaceTree", 0)
	}
	utils.RestartExplorer()
	return true
}

func GetRemoveShortcutArrowStatus() bool {
	v, _, err := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "29")
	return err == nil && v != ""
}

func SetRemoveShortcutArrow(on bool) bool {
	if on {
		utils.RegWriteString(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "29", "")
	} else {
		utils.RegDeleteValue(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "29")
	}
	utils.RestartExplorer()
	return true
}

func GetRemoveShortcutTextStatus() bool {
	v, _, err := utils.RegReadString(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Explorer\NamingTemplate`, "ShortcutNameTemplate")
	return err == nil && v != ""
}

func SetRemoveShortcutText(on bool) bool {
	if on {
		utils.RegWriteString(registry.CURRENT_USER,
			`Software\Microsoft\Windows\CurrentVersion\Explorer\NamingTemplate`, "ShortcutNameTemplate", "%s.lnk")
	} else {
		utils.RegDeleteValue(registry.CURRENT_USER,
			`Software\Microsoft\Windows\CurrentVersion\Explorer\NamingTemplate`, "ShortcutNameTemplate")
	}
	utils.RestartExplorer()
	return true
}

func GetRemoveShieldStatus() bool {
	v, _, err := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "77")
	return err == nil && v != ""
}

func SetRemoveShield(on bool) bool {
	if on {
		utils.RegWriteString(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "77",
			`%SystemRoot%\System32\shell32.dll,-278`)
	} else {
		utils.RegDeleteValue(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Shell Icons`, "77")
	}
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
	var script []byte
	var name string
	if enable {
		script = utils.EnableTaskManagerScript
		name = "enableTaskManager.ps1"
	} else {
		script = utils.DisableTaskManagerScript
		name = "disableTaskManager.ps1"
	}

	if !utils.SuperExecute(script, name) {
		return false
	}

	source := taskManagerPath
	if enable {
		source = taskManagerPath + ".bak"
	}
	for i := 0; i < 10; i++ {
		time.Sleep(300 * time.Millisecond)
		_, err := os.Stat(source)
		if os.IsNotExist(err) {
			return true
		}
	}
	return false
}
