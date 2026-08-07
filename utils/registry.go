//go:build windows

package utils

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows/registry"
)

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

func RegWriteString(key registry.Key, path, name, value string) bool {
	return RegWriteStringE(key, path, name, value) == nil
}

func RegWriteStringE(key registry.Key, path, name, value string) error {
	k, _, err := registry.CreateKey(key, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue(name, value)
}

func RegDeleteValue(key registry.Key, path, name string) {
	_ = RegDeleteValueE(key, path, name)
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

func RegSetDWord(key registry.Key, path, name string, value uint32) {
	_ = RegSetDWordE(key, path, name, value)
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

func RegSetDWordBool(key registry.Key, path, name string, value uint32) bool {
	RegSetDWord(key, path, name, value)
	return true
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
