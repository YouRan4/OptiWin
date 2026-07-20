//go:build windows

package utils

import (
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
	k, _, _ := registry.CreateKey(key, path, registry.SET_VALUE)
	defer k.Close()
	return k.SetStringValue(name, value) == nil
}

func RegDeleteValue(key registry.Key, path, name string) {
	k, err := registry.OpenKey(key, path, registry.SET_VALUE)
	if err != nil {
		return
	}
	defer k.Close()
	k.DeleteValue(name)
}

func RegSetDWord(key registry.Key, path, name string, value uint32) {
	k, err := registry.OpenKey(key, path, registry.SET_VALUE)
	if err != nil {
		k, _, _ = registry.CreateKey(key, path, registry.SET_VALUE)
	}
	defer k.Close()
	k.SetDWordValue(name, value)
}

func RegSetDWordBool(key registry.Key, path, name string, value uint32) bool {
	RegSetDWord(key, path, name, value)
	return true
}
