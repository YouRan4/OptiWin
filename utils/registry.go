//go:build windows

package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

// Logf 追加一条日志到 %LOCALAPPDATA%\OptiWin\logs\app.log；日志写入失败时静默跳过
func Logf(format string, args ...any) {
	logDir := LogDir()
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return
	}
	logPath := filepath.Join(logDir, logFileName)
	rotateLog(logDir, logPath)
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
}

// rotateLog 在 app.log 超过 1MB 时重命名为带时间戳的 app-<20060102-150405>.log，
// 并只保留最近 5 个滚动文件；重命名失败（如文件被占用）时继续追写原文件
func rotateLog(logDir, logPath string) {
	info, err := os.Stat(logPath)
	if err != nil || info.Size() < logMaxBytes {
		return
	}
	rotated := "app-" + time.Now().Format("20060102-150405") + ".log"
	if err := os.Rename(logPath, filepath.Join(logDir, rotated)); err != nil {
		return
	}
	pruneRolledLogs(logDir)
}

func pruneRolledLogs(logDir string) {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return
	}
	var rolled []string
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() && strings.HasPrefix(name, "app-") && strings.HasSuffix(name, ".log") {
			rolled = append(rolled, name)
		}
	}
	// 文件名里的时间戳可直接按字典序比较，等价于按时间排序
	sort.Strings(rolled)
	for _, name := range rolled[:max(0, len(rolled)-logKeepRolled)] {
		_ = os.Remove(filepath.Join(logDir, name))
	}
}

func RegReadDWord(key registry.Key, path, name string) (uint64, error) {
	k, err := registry.OpenKey(key, path, registry.QUERY_VALUE)
	if err != nil {
		return 0, err
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue(name)
	return val, err
}

func RegReadString(key registry.Key, path, name string) (string, uint32, error) {
	k, err := registry.OpenKey(key, path, registry.QUERY_VALUE)
	if err != nil {
		return "", 0, err
	}
	defer k.Close()
	return k.GetStringValue(name)
}

func RegWriteStringE(key registry.Key, path, name, value string) error {
	k, _, err := registry.CreateKey(key, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue(name, value)
}

func RegDeleteValueE(key registry.Key, path, name string) error {
	k, err := registry.OpenKey(key, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.DeleteValue(name)
}

// RegDeleteKeyE 递归删除注册表键（含所有子键），键不存在视为成功
func RegDeleteKeyE(hive registry.Key, path string) error {
	k, err := registry.OpenKey(hive, path, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return nil
		}
		return err
	}
	subKeys, err := k.ReadSubKeyNames(-1)
	k.Close()
	if err != nil {
		return err
	}
	for _, sub := range subKeys {
		if err := RegDeleteKeyE(hive, path+`\`+sub); err != nil {
			return err
		}
	}
	return registry.DeleteKey(hive, path)
}

func RegSetDWordE(key registry.Key, path, name string, value uint32) error {
	k, err := registry.OpenKey(key, path, registry.SET_VALUE)
	if err != nil {
		k, _, err = registry.CreateKey(key, path, registry.SET_VALUE)
		if err != nil {
			return err
		}
	}
	defer k.Close()
	return k.SetDWordValue(name, value)
}

func GetSystemAccentColor() string {
	val, err := RegReadDWord(registry.CURRENT_USER, `Software\Microsoft\Windows\DWM`, "AccentColor")
	if err != nil {
		return ""
	}
	r := byte(val)
	g := byte(val >> 8)
	b := byte(val >> 16)
	// Windows 对极暗或极亮的颜色会自动调整，注册表值不准确
	lum := float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114
	if lum < 40 || lum > 220 {
		return ""
	}
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}
