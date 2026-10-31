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

func LogDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, appDirName)
}
