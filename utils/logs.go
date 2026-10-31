package utils

import (
	"os"
	"path/filepath"
)

const (
	logFileName   = "app.log"
	logMaxBytes   = 1 << 20
	logKeepRolled = 5
)

// LogDir 返回 %LOCALAPPDATA%\OptiWin\logs；取不到缓存目录时退回临时目录
func LogDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "OptiWin", "logs")
}
