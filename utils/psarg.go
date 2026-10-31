package utils

// ValidPSArg 校验要拼进 PowerShell 脚本字符串的参数值（如路径、文件名）是否安全。
// 只允许字母、数字和空格、点、冒号、反斜杠、% 与环境变量占位符、下划线、连字符；
// 任何需要拼进脚本的用户输入都必须先过它（或用 psEncode 的参数化写法），否则就是命令注入。
// 注意：参数以 argv 数组传递的调用（SuperExec/SuperReg）不经过 shell，不需要它。
func ValidPSArg(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' ||
			r == ' ' || r == '.' || r == ':' || r == '\\' || r == '%' || r == '_' || r == '-') {
			return false
		}
	}
	return true
}
