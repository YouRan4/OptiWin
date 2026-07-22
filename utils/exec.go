//go:build windows

package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

func HideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

func RunHide(name string, args ...string) {
	cmd := exec.Command(name, args...)
	HideWindow(cmd)
	cmd.Run()
}

func RestartExplorer() {
	cacheDir := os.Getenv("LOCALAPPDATA") + `\Microsoft\Windows\Explorer`
	matches, _ := filepath.Glob(cacheDir + `\iconcache*`)
	for _, m := range matches {
		os.Remove(m)
	}

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
