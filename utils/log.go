//go:build windows

package utils

import (
	"os"
	"path/filepath"
)

func LogError(msg string) {
	logPath := filepath.Join(os.TempDir(), "revitool.log")
	f, _ := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(msg + "\n")
}
