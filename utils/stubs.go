//go:build !windows

package utils

import (
	"net/http"
	"os/exec"
)

func HideWindow(cmd *exec.Cmd) {}
func ScHide(args ...string)    {}
func RunHide(name string, args ...string) {}
func RestartExplorer() {}
func LogError(msg string) {}
func GetHttpClient() *http.Client { return &http.Client{} }
func GetProxyInfo() string        { return "未使用代理" }
