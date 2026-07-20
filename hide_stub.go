//go:build !windows

package main

import "os/exec"

// 非 Windows 平台空实现
func hideWindow(cmd *exec.Cmd) {}
