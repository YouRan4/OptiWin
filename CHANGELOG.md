# Changelog

## v1.5.1

### 新增
- IFEO 镜像劫持管理：查看、添加、删除 IFEO 条目，支持系统进程黑名单
- APPX 应用管理：查看已安装应用、支持卸载、搜索过滤、用户/系统应用分类
- DNS 切换：一键切换 DNS 服务器，支持 ISP 默认、阿里、腾讯、百度、Google、Cloudflare、Quad9
- 当前 DNS 信息显示：网卡名称、主 DNS、副 DNS
- 系统遥测与推广阻止：一键阻止 Windows 遥测进程

### 改进
- **TrustedInstaller 执行引擎重写**：移除外部 PowerRun.exe 依赖，改用原生 Go 实现，通过 SYSTEM 模拟（`ImpersonateLoggedOnUser`）获取 TI 令牌，全部使用文档化 Win32 API
- 脚本执行改为 `-EncodedCommand` 方式，不再写入本地临时文件，消除临时文件泄漏风险
- `ExecuteFile` / `SuperExecuteFile` 方法：支持直接执行指定路径的 PS1 脚本（适用于超大脚本）
- IFEO 输入校验：前端与后端同步验证 exe 名称格式、系统进程黑名单、调试器路径合法性（禁止 UNC 路径、路径穿越）
- 进程名匹配改为路径验证：`findSystemProcess` + `QueryFullProcessImageName` 防止同名进程欺骗
- `buildCommandLine` 使用 `syscall.EscapeArg` 正确处理命令行引号
- `GetExitCodeProcess` 退出码检查，脚本执行失败时返回明确错误
- 删除 PowerRun.exe 嵌入依赖，减小打包体积
- Vite 构建配置优化：target=chrome120、chunk 分割、cssCodeSplit=false

### 移除
- PowerRun.exe 外部依赖（已完全替换为原生 Go 实现）
- DNS over HTTPS (DoH) 功能（已简化为基础 DNS 切换）
- 临时文件写入/删除逻辑（改为 `-EncodedCommand` 不落盘执行）
- 未使用的 i18n 键和前端静态资源

---

## v1.4.1

### 改进
- Game Bar 从开关改为独立的卸载/安装按钮，执行时弹出 WaitModal
- Defender 禁用/恢复重构为 PS1 脚本，减少 Go 冗余代码

---

## v1.4.0

### 新增
- Win11 任务管理器切换（一键恢复 Win10 风格）
- Game Bar 开关（独立于 Xbox）
- 全局图标支持（Lucide 图标库）

### 改进
- 优化了 UI 交互体验
