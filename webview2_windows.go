//go:build windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

func ensureWebView2() {
	if isWebView2Installed() {
		return
	}

	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("OptiWin 需要 WebView2 运行时才能运行。\n即将打开安装程序，请完成安装后重新打开 OptiWin。"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("缺少组件"))),
		0x00000040, // MB_ICONINFORMATION
	)

	tmp := filepath.Join(os.TempDir(), "OptiWin_WebView2Setup.exe")
	os.WriteFile(tmp, webview2Installer, 0755)

	exec.Command(tmp).Start()
	os.Exit(0)
}

func isWebView2Installed() bool {
	guids := []string{
		`SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`,
		`Software\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`,
	}
	for _, hive := range []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER} {
		for _, g := range guids {
			k, err := registry.OpenKey(hive, g, registry.QUERY_VALUE)
			if err != nil {
				continue
			}
			defer k.Close()
			_, _, err = k.GetStringValue("pv")
			if err == nil {
				return true
			}
		}
	}
	return false
}
