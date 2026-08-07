//go:build windows

package utils

import (
	"OptiWin/tilauncher"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

// psEncode 将脚本编码为 PowerShell -EncodedCommand 所需的 Base64(UTF-16LE)
func psEncode(script string) string {
	encoded := utf16.Encode([]rune(script))
	buf := make([]byte, len(encoded)*2)
	for i, v := range encoded {
		binary.LittleEndian.PutUint16(buf[i*2:], v)
	}
	return base64.StdEncoding.EncodeToString(buf)
}

func HideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}

func RunHide(name string, args ...string) {
	cmd := exec.Command(name, args...)
	HideWindow(cmd)
	cmd.Run()
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

const (
	wmExitExplorer = 0x5B4 // WM_USER + 436：explorer 内部优雅退出消息
	synchronize    = 0x00100000
	explorerWait   = 3 * time.Second
)

func shellTrayHwnd() uintptr {
	className, _ := syscall.UTF16PtrFromString("Shell_TrayWnd")
	hwnd, _, _ := user32.NewProc("FindWindowW").Call(
		uintptr(unsafe.Pointer(className)), 0,
	)
	return hwnd
}

func waitProcessExit(pid uint32, timeout time.Duration) bool {
	hProcess, _, _ := kernel32.NewProc("OpenProcess").Call(synchronize, 0, uintptr(pid))
	if hProcess == 0 {
		return true
	}
	defer kernel32.NewProc("CloseHandle").Call(hProcess)
	wait, _, _ := kernel32.NewProc("WaitForSingleObject").Call(hProcess, uintptr(timeout.Milliseconds()))
	return wait == 0 // WAIT_OBJECT_0 = 进程已退出
}

func RestartExplorer() {
	hwnd := shellTrayHwnd()
	var pid uint32
	user32.NewProc("GetWindowThreadProcessId").Call(
		hwnd, uintptr(unsafe.Pointer(&pid)),
	)
	if pid == 0 {
		return
	}

	user32.NewProc("PostMessageW").Call(hwnd, wmExitExplorer, 0, 0)
	waitProcessExit(pid, explorerWait)
	startExplorer()
}

func startExplorer() {
	explorer, _ := syscall.UTF16PtrFromString("explorer.exe")
	open, _ := syscall.UTF16PtrFromString("open")
	shell32.NewProc("ShellExecuteW").Call(
		0, uintptr(unsafe.Pointer(open)),
		uintptr(unsafe.Pointer(explorer)), 0, 0, 5,
	)
}

func ExecuteString(s string) bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-EncodedCommand", psEncode(s))
	HideWindow(cmd)
	return cmd.Run() == nil
}

func SuperExecuteString(s string) bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-EncodedCommand", psEncode(s))
	return tilauncher.RunAsTrustedInstaller(cmd) == nil
}

func SuperExecute(data []byte) bool {
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("sys_%d.ps1", time.Now().UnixNano()))
	if err := os.WriteFile(tempFile, data, 0600); err != nil {
		return false
	}
	setFileHiddenAndReadOnly(tempFile)
	defer os.Remove(tempFile)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", tempFile)
	return tilauncher.RunAsTrustedInstaller(cmd) == nil
}

// SuperExec 以 TrustedInstaller 令牌执行任意命令
func SuperExec(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	return tilauncher.RunAsTrustedInstaller(cmd) == nil
}

// SuperReg 以 TrustedInstaller 令牌执行 reg.exe 命令
func SuperReg(args ...string) bool {
	return SuperExec("reg.exe", args...)
}

func Execute(data []byte) bool {
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("sys_%d.ps1", time.Now().UnixNano()))
	if err := os.WriteFile(tempFile, data, 0600); err != nil {
		return false
	}
	setFileHiddenAndReadOnly(tempFile)
	defer os.Remove(tempFile)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", tempFile)
	HideWindow(cmd)
	return cmd.Run() == nil
}

func setFileHiddenAndReadOnly(filePath string) {
	ptr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return
	}
	attributes := uint32(0x2 | 0x4 | 0x1)
	_ = syscall.SetFileAttributes(ptr, attributes)
}
