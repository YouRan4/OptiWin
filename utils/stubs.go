//go:build !windows

package utils

import (
	"net/http"
	"os/exec"
	"time"
)

func HideWindow(cmd *exec.Cmd)                                  {}
func RunHide(name string, args ...string) bool                  { return false }
func ExecuteStringTimeout(s string, timeout time.Duration) bool { return false }
func Logf(format string, args ...any)                           {}
func RestartExplorer() error                                    { return nil }
func OpenLogDir() bool                                          { return false }
func TerminateProcesses(imageName string) (int, error)          { return 0, nil }
func GetHTTPClient() *http.Client                               { return &http.Client{} }
func GetProxyInfo() string                                      { return "未使用代理" }
func GetSystemAccentColor() string                              { return "" }
