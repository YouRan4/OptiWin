//go:build windows
package services

import (
	"golang.org/x/sys/windows/registry"
	"OptiWin/utils"
)

func GetNotificationStatus() string {
	toast, _ := utils.RegReadDWord(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled")
	center, _ := utils.RegReadDWord(registry.CURRENT_USER,
		`SOFTWARE\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter")
	if center == 1 {
		return "完全关闭"
	}
	if toast == 0 {
		return "仅关闭通知"
	}
	return "开启"
}

func SetNotificationMode(mode string) bool {
	switch mode {
	case "开启":
		utils.RegDeleteValue(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`, "NoToastApplicationNotification")
		utils.RegDeleteValue(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter")
		utils.RegDeleteValue(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled")
		utils.RegDeleteValue(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled")
	case "仅关闭通知":
		utils.RegSetDWord(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\CurrentVersion\PushNotifications`, "NoToastApplicationNotification", 1)
		utils.RegSetDWord(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
		utils.RegSetDWord(registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", 0)
	case "完全关闭":
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
