package services

import "strings"

// ifeoBlockedExecutables Windows 核心进程，禁止通过 IFEO 劫持（会导致系统无法启动）
// 键必须统一小写：查询方用 strings.ToLower(name) 查表
var ifeoBlockedExecutables = map[string]bool{
	"winlogon.exe":      true,
	"csrss.exe":         true,
	"lsass.exe":         true,
	"smss.exe":          true,
	"services.exe":      true,
	"svchost.exe":       true,
	"explorer.exe":      true,
	"dwm.exe":           true,
	"wininit.exe":       true,
	"fontdrvhost.exe":   true,
	"sihost.exe":        true,
	"runtimebroker.exe": true,
}

// isValidIfeoExeName 校验 IFEO 条目的可执行文件名：必须以 .exe 结尾、不含路径字符、不在核心进程黑名单内
func isValidIfeoExeName(name string) bool {
	if name == "" || !strings.HasSuffix(strings.ToLower(name), ".exe") {
		return false
	}
	if strings.ContainsAny(name, `\/:;*?<>|`) {
		return false
	}
	if ifeoBlockedExecutables[strings.ToLower(name)] {
		return false
	}
	return true
}

// isValidDebuggerPath 校验 IFEO 调试器路径：空值合法，否则须为绝对路径或环境变量路径，且不含 UNC 与路径穿越
func isValidDebuggerPath(path string) bool {
	if path == "" {
		return true
	}
	lower := strings.ToLower(path)
	// 允许环境变量路径（%windir%、%systemroot% 等）或盘符绝对路径
	if !strings.Contains(lower, `:\`) && !strings.Contains(lower, `%`) {
		return false
	}
	if strings.HasPrefix(lower, `\\`) {
		return false
	}
	if strings.Contains(lower, "..") {
		return false
	}
	return true
}
