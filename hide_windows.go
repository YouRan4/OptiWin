//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

// 隐藏外部进程的控制台窗口，避免弹 cmd 黑框
func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
