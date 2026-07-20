package services

// 安全页面相关注册表操作（待实现）
func ToggleDefender(on bool) bool {
	return SetToggle("defender", on)
}
