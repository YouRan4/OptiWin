package utils

import (
	"os"
	"path/filepath"
)

const appDirName = "OptiWin"

func AppDataDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, appDirName)
}

func WebViewDataDir() string { return filepath.Join(AppDataDir(), "WebView2") }

func LogPath() string { return filepath.Join(AppDataDir(), logFileName) }
