package services

// 性能页面相关注册表操作（待实现）
func TogglePowerPlan(on bool) bool  { return SetToggle("power_plan", on) }
func ToggleSuperfetch(on bool) bool { return SetToggle("superfetch", on) }
func ToggleMemCompress(on bool) bool  { return SetToggle("mem_compress", on) }
