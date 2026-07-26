# Changelog

## v1.5.0

### 新增
- IFEO 镜像劫持管理：查看、添加、删除 IFEO 条目
- APPX 应用管理：查看已安装应用、支持卸载、搜索过滤、用户/系统应用分类
- DNS 切换：一键切换 DNS 服务器，支持阿里、腾讯、百度、Google、Cloudflare、Quad9
- DNS over HTTPS (DoH) 支持：自动配置 DoH 模板
- 系统遥测与推广阻止：一键阻止 Windows 系统遥测进程
- 当前 DNS 信息显示：网卡、主DNS、副DNS、DoH 地址与状态

### Improvements
- IFEO management: view, add, delete IFEO entries
- APPX app management: view installed apps, uninstall support, search filter, user/system app categorization
- DNS switching: one-click DNS server switch with Ali, Tencent, Baidu, Google, Cloudflare, Quad9 support
- DNS over HTTPS (DoH) support: auto-configure DoH templates
- System telemetry and promotion blocking: one-click block Windows telemetry processes
- Current DNS info display: adapter, primary DNS, secondary DNS, DoH address and status

## v1.4.1

### 改进
- 统一 PowerShell 脚本执行引擎，所有系统级操作通过 PS1 脚本 + PowerRun 提权执行
- Game Bar 从开关改为独立的卸载/安装按钮，执行时弹出 WaitModal
- Defender 禁用/恢复重构为 PS1 脚本，减少 Go 冗余代码
- PowerRun 嵌入逻辑内联到公共执行器，简化代码结构

---

### Improvements
- Unified PowerShell script execution engine; all system operations via PS1 + PowerRun
- Game Bar reworked from toggle to separate uninstall/install buttons with WaitModal
- Defender disable/restore refactored to PS1 scripts, reduced Go boilerplate
- PowerRun embedding inlined into shared executor, simplified code structure

## v1.4.0

### 新增
- Win11 任务管理器切换（一键恢复 Win10 风格）
- Game Bar 开关（独立于 Xbox）
- 全局图标支持（Lucide 图标库）

### 改进
- 优化了 UI 交互体验

---

### New Features
- Win11 Task Manager toggle (one-click restore to Win10 style)
- Game Bar toggle (independent of Xbox)
- Global icon support (Lucide icon library)

### Improvements
- Improved UI interaction experience
