//go:build windows

package utils

import (
	"OptiWin/tilauncher"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func RunHide(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	HideWindow(cmd)
	cmd.WaitDelay = 5 * time.Second
	return runLogged(fmt.Sprintf("%s %s", name, strings.Join(args, " ")), cmd)
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

const (
	// 	杀进程后等旧 explorer 全部退出 / 系统自动拉起 shell / 手动启动后任务栏出现的上限
	explorerExitWait  = 3 * time.Second
	explorerAutoWait  = 3 * time.Second
	explorerShellWait = 10 * time.Second
)

func shellTrayHwnd() uintptr {
	className, _ := syscall.UTF16PtrFromString("Shell_TrayWnd")
	hwnd, _, _ := user32.NewProc("FindWindowW").Call(
		uintptr(unsafe.Pointer(className)), 0,
	)
	return hwnd
}

func trayPid() (uintptr, uint32) {
	hwnd := shellTrayHwnd()
	if hwnd == 0 {
		return 0, 0
	}
	var pid uint32
	user32.NewProc("GetWindowThreadProcessId").Call(
		hwnd, uintptr(unsafe.Pointer(&pid)),
	)
	return hwnd, pid
}

func waitShellReady(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		if hwnd, _ := trayPid(); hwnd != 0 {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func waitExplorerGone(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		if left, err := CountProcesses("explorer.exe"); err == nil && left == 0 {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// RestartExplorer 直接结束 explorer 进程让系统重新拉起 shell：
// 必须先确认旧进程全部退出再判断任务栏，否则 FindWindow 找到的是旧进程
// 残留的任务栏窗口，会误判成功；旧进程没退干净就启动新的，新进程只会把
// 请求转交给旧进程后自己退出，表现就是“点了没反应 / 任务栏起不来”
func RestartExplorer() error {
	// 只杀带任务栏的那个进程会漏掉单独开的文件资源管理器窗口
	n, err := TerminateProcesses("explorer.exe")
	if err != nil {
		Logf("[explorer] 强制结束 explorer.exe 失败（已结束 %d 个）: %v", n, err)
	}
	if n > 0 || shellTrayHwnd() != 0 {
		if !waitExplorerGone(explorerExitWait) {
			left, _ := CountProcesses("explorer.exe")
			return fmt.Errorf("资源管理器进程未能结束（剩余 %d 个）", left)
		}
		// 强杀会触发系统自动拉起 shell；短等一下，免得再手动起一个、
		// 多开一个文件资源管理器窗口
		if waitShellReady(explorerAutoWait) {
			return nil
		}
	}
	if err := startExplorer(); err != nil {
		Logf("[explorer] 第一次启动失败: %v", err)
	}
	if waitShellReady(explorerShellWait) {
		return nil
	}
	if err := startExplorer(); err != nil {
		Logf("[explorer] 重试启动失败: %v", err)
	}
	if waitShellReady(explorerShellWait) {
		return nil
	}
	return errors.New("资源管理器未能启动，请在任务管理器里手动运行 explorer.exe")
}

// terminateProcess 直接结束指定进程：不经过 taskkill.exe，
// 避免命令行特征被安全软件当作批量杀进程拦截
func terminateProcess(pid uint32) error {
	const processTerminate = 0x0001
	h, _, err := kernel32.NewProc("OpenProcess").Call(processTerminate, 0, uintptr(pid))
	if h == 0 {
		return fmt.Errorf("OpenProcess: %w", err)
	}
	defer kernel32.NewProc("CloseHandle").Call(h)
	if ret, _, terr := kernel32.NewProc("TerminateProcess").Call(h, 1); ret == 0 {
		return fmt.Errorf("TerminateProcess: %w", terr)
	}
	return nil
}

func OpenLogDir() bool {
	dir := LogDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		Logf("[utils] 创建日志目录失败: %v", err)
		SetLastError(fmt.Sprintf("创建日志目录失败: %v", err))
		return false
	}
	// 资源管理器用 ShellExecute 打开目录：子进程方式带隐藏窗口标志时，
	// 新窗口会被 STARTUPINFO 的 SW_HIDE 吞掉，表现就是“点了没反应”
	if err := shellOpen(dir); err != nil {
		Logf("[utils] ShellExecute 打开目录失败: %v", err)
		if start := exec.Command("explorer.exe", dir).Start(); start != nil {
			Logf("[utils] 回退启动 explorer.exe 也失败: %v", start)
			SetLastError(fmt.Sprintf("打开日志目录失败: %v", err))
			return false
		}
	}
	return true
}

func shellOpen(target string) error {
	verb, _ := syscall.UTF16PtrFromString("open")
	path, _ := syscall.UTF16PtrFromString(target)
	ret, _, err := shell32.NewProc("ShellExecuteW").Call(
		0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(path)), 0, 0, 1,
	)
	if ret <= 32 {
		return fmt.Errorf("ShellExecute(%s): %v", target, err)
	}
	return nil
}

func startExplorer() error {
	explorer, _ := syscall.UTF16PtrFromString("explorer.exe")
	open, _ := syscall.UTF16PtrFromString("open")
	ret, _, err := shell32.NewProc("ShellExecuteW").Call(
		0, uintptr(unsafe.Pointer(open)),
		uintptr(unsafe.Pointer(explorer)), 0, 0, 5,
	)
	if ret <= 32 {
		return fmt.Errorf("ShellExecute explorer.exe 失败: %v", err)
	}
	return nil
}

func runLogged(label string, cmd *exec.Cmd) bool {
	return runLoggedCtx(nil, label, cmd)
}

func runLoggedCtx(ctx context.Context, label string, cmd *exec.Cmd) bool {
	out, err := cmd.CombinedOutput()
	if ctx != nil && ctx.Err() == context.DeadlineExceeded {
		Logf("[exec] %s 超时已终止 | 输出: %s", label, strings.TrimSpace(string(out)))
		SetLastError(fmt.Sprintf("%s 超时", label))
		return false
	}
	if err != nil {
		Logf("[exec] %s 失败: %v | 输出: %s", label, err, strings.TrimSpace(string(out)))
		SetLastError(fmt.Sprintf("%s: %v", label, err))
		return false
	}
	return true
}

func ExecuteString(s string) bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-EncodedCommand", psEncode(s))
	HideWindow(cmd)
	cmd.WaitDelay = 5 * time.Second
	return runLogged("ExecuteString", cmd)
}

// ExecuteStringTimeout 同 ExecuteString，但到点强制终止进程；
// 用于可能被系统服务卡住、不能让界面无限等待的命令
func ExecuteStringTimeout(s string, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-EncodedCommand", psEncode(s))
	HideWindow(cmd)
	// 子进程（如 Appx 部署）可能继承管道句柄；没有 WaitDelay 时即使杀掉了 powershell，
	// CombinedOutput 也会一直等管道关闭，表现为“超时了还卡住”
	cmd.WaitDelay = 5 * time.Second
	return runLoggedCtx(ctx, fmt.Sprintf("ExecuteStringTimeout(%s)", timeout), cmd)
}

// SuperExecuteString 以 TrustedInstaller 身份执行 PowerShell 脚本。
// 约定：脚本只能由受控常量拼接；任何拼入脚本的动态值必须先用 ValidPSArg 校验
// （SetOldTaskManager 已如此实现），否则它就是以 TI 权限执行任意命令的入口。
func SuperExecuteString(s string) bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-EncodedCommand", psEncode(s))
	if err := tilauncher.RunAsTrustedInstaller(cmd); err != nil {
		Logf("[exec] SuperExecuteString 失败: %v", err)
		return false
	}
	return true
}

func SuperExecute(data []byte) bool {
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("sys_%d.ps1", time.Now().UnixNano()))
	if err := os.WriteFile(tempFile, data, 0600); err != nil {
		Logf("[exec] SuperExecute 写入临时脚本失败: %v", err)
		return false
	}
	setFileHiddenAndReadOnly(tempFile)
	defer os.Remove(tempFile)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", tempFile)
	if err := tilauncher.RunAsTrustedInstaller(cmd); err != nil {
		Logf("[exec] SuperExecute 失败: %v", err)
		return false
	}
	return true
}

// SuperExec 以 TrustedInstaller 令牌执行任意命令。
// 参数以 argv 数组交给 CreateProcess，不经过 shell，因此不需要 ValidPSArg；
// 需要拼进 PowerShell 脚本的值请改用 SuperExecuteString + ValidPSArg。
func SuperExec(name string, args ...string) bool {
	if err := tilauncher.RunAsTrustedInstaller(exec.Command(name, args...)); err != nil {
		Logf("[exec] SuperExec %s %v 失败: %v", name, args, err)
		return false
	}
	return true
}

// SuperReg 以 TrustedInstaller 令牌执行 reg.exe 命令
func SuperReg(args ...string) bool {
	return SuperExec("reg.exe", args...)
}

func Execute(data []byte) bool {
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("sys_%d.ps1", time.Now().UnixNano()))
	if err := os.WriteFile(tempFile, data, 0600); err != nil {
		Logf("[exec] Execute 写入临时脚本失败: %v", err)
		return false
	}
	setFileHiddenAndReadOnly(tempFile)
	defer os.Remove(tempFile)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", tempFile)
	HideWindow(cmd)
	return runLogged("Execute", cmd)
}

func setFileHiddenAndReadOnly(filePath string) {
	ptr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return
	}
	attributes := uint32(0x2 | 0x4 | 0x1)
	_ = syscall.SetFileAttributes(ptr, attributes)
}
