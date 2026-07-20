package services

import "runtime"

// SetToggle 通用开关写入（待接入真实注册表）
func SetToggle(key string, on bool) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	return true
}

// GetToggle 通用开关读取（待接入真实注册表）
func GetToggle(key string) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	return false
}

// SetNotification 通知模式设置（待接入真实注册表）
func SetNotification(mode string) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	return true
}
