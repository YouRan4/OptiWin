//go:build windows

package utils

import (
	"OptiWin/tilauncher"
	"encoding/base64"
	"encoding/binary"
	"os/exec"
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

func Execute(data []byte) bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-EncodedCommand", psEncode(string(data)))
	HideWindow(cmd)
	return cmd.Run() == nil
}

func SuperExecute(data []byte) bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-EncodedCommand", psEncode(string(data)))
	return tilauncher.RunAsTrustedInstaller(cmd) == nil
}

func ExecuteFile(scriptPath string) bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	HideWindow(cmd)
	return cmd.Run() == nil
}

func SuperExecuteFile(scriptPath string) bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	return tilauncher.RunAsTrustedInstaller(cmd) == nil
}
