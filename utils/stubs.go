//go:build !windows

package utils

import (
	"net/http"
	"os/exec"
)

var (
	RemoveGameBarScript      []byte
	RestoreGameBarScript     []byte
	DisableDefenderScript    []byte
	RestoreDefenderScript    []byte
	EnableTaskManagerScript  []byte
	DisableTaskManagerScript []byte
)

func HideWindow(cmd *exec.Cmd)                   {}
func RunHide(name string, args ...string)        {}
func RestartExplorer()                           {}
func GetHTTPClient() *http.Client                { return &http.Client{} }
func GetProxyInfo() string                       { return "未使用代理" }
func GetSystemAccentColor() string               { return "" }
func Execute(data []byte) bool           { return false }
func SuperExecute(data []byte) bool      { return false }
func ExecuteFile(path string) bool       { return false }
func SuperExecuteFile(path string) bool  { return false }
