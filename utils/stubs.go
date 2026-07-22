//go:build !windows

package utils

import (
	"net/http"
	"os/exec"
)

func HideWindow(cmd *exec.Cmd) {}
func RunHide(name string, args ...string) {}
func RestartExplorer() {}
func GetHTTPClient() *http.Client { return &http.Client{} }
func GetProxyInfo() string        { return "未使用代理" }
func GetPowerRunPath() string     { return "" }
func CleanupPowerRun()            {}
func GetSystemAccentColor() string { return "" }
