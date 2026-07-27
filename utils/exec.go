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

func RestartExplorer() {
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
