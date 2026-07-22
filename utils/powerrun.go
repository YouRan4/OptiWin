//go:build windows

package utils

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed PowerRun.exe
var powerRunBin []byte

func GetPowerRunPath() string {
	path := filepath.Join(os.TempDir(), "prun.exe")
	os.WriteFile(path, powerRunBin, 0755)
	return path
}

func CleanupPowerRun() {
	os.Remove(filepath.Join(os.TempDir(), "prun.exe"))
}
