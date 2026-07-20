//go:build !windows

package utils

import "os/exec"

func HideWindow(cmd *exec.Cmd) {}
func ScHide(args ...string)    {}
func RunHide(name string, args ...string) {}
func RestartExplorer() {}
func LogError(msg string) {}
