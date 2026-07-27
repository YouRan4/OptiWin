# Changelog

---

## v1.5.2

### 修复
- 修复移除快捷方式小箭头开关状态判断错误的问题

---

## v1.5.1

### 改进
- **TrustedInstaller 执行引擎重写**：移除外部 PowerRun.exe 依赖，改用原生 Go 实现，通过 SYSTEM 模拟（`ImpersonateLoggedOnUser`）获取 TI 令牌，全部使用文档化 Win32 API
- **脚本执行架构重构**：`Execute` / `SuperExecute` 改为写入临时文件 + `-File` 模式执行（解决 `-EncodedCommand` 在 TI 进程下命令行截断问题）
- 新增 `ExecuteString` / `SuperExecuteString`：支持 `-EncodedCommand` 模式执行短脚本字符串
- 临时文件安全处理：写入后设置**隐藏 + 只读**属性，执行完毕延迟 1 秒自动清理
- 脚本资源迁移：`ps1/` 目录脚本内嵌至 `scriptmgr/` 模块，通过 `GetScriptBytes` 按需加载
- Game Bar / Task Manager 脚本内联化：短脚本直接写入 Go 代码，消除外部文件依赖
- IFEO 输入校验增强：前端与后端同步验证 exe 名称格式、系统进程黑名单、调试器路径合法性（支持环境变量路径 `%windir%`、禁止 UNC 路径、路径穿越）
- 进程名匹配改为路径验证：`findSystemProcess` + `QueryFullProcessImageName` 防止同名进程欺骗
- `buildCommandLine` 使用 `syscall.EscapeArg` 正确处理命令行引号
- `GetExitCodeProcess` 退出码检查，脚本执行失败时返回明确错误
- 删除 PowerRun.exe 嵌入依赖，减小打包体积
- Vite 构建配置优化：target=chrome120、chunk 分割、cssCodeSplit=false

### 移除
- `ExecuteFile` / `SuperExecuteFile` 方法（功能合并至 `Execute` / `SuperExecute`）
- `utils/scripts.go` 全局脚本变量（替换为 `scriptmgr` 模块）
- `ps1/` 目录嵌入脚本（迁移至 `scriptmgr/scripts/`）
- 启动时临时目录创建/清理逻辑（改为单文件按需创建）

---

## v1.5.0

### 新增
- IFEO 镜像劫持管理：查看、添加、删除 IFEO 条目，支持系统进程黑名单
- APPX 应用管理：查看已安装应用、支持卸载、搜索过滤、用户/系统应用分类
- DNS 切换：一键切换 DNS 服务器，支持 ISP 默认、阿里、腾讯、百度、Google、Cloudflare、Quad9
- 当前 DNS 信息显示：网卡名称、主 DNS、副 DNS
- 系统遥测与推广阻止：一键阻止 Windows 遥测进程

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
