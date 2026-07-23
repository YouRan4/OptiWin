//go:build windows

package utils

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"github.com/google/uuid"
)

//go:embed PowerRun.exe
var powerRunBin []byte

func HideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
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
	id := uuid.New()
	scriptPath := filepath.Join(os.TempDir(), "optiwin", id.String()+".ps1")
	os.WriteFile(scriptPath, data, 0644)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	HideWindow(cmd)
	return cmd.Run() == nil
}

func SuperExecute(data []byte) bool {
	id := uuid.New()
	scriptPath := filepath.Join(os.TempDir(), "optiwin", id.String()+".ps1")
	os.WriteFile(scriptPath, data, 0644)
	prPath := filepath.Join(os.TempDir(), "optiwin", id.String()+".exe")
	os.WriteFile(prPath, powerRunBin, 0755)
	cmd := exec.Command(prPath, "/SW:0", "powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	r := cmd.Run() == nil
	time.Sleep(500 * time.Millisecond)
	return r
}
