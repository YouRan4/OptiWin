# Changelog

---

## v1.7.1

### 修复
- **UninstallEdge 不再误杀工具箱自身 UI**：移除 `taskkill msedgewebview2.exe`——它是 WebView2 Runtime 进程，工具箱界面同样基于 WebView2，原命令会连同工具箱一起终止；Edge 浏览器本体为 `msedge.exe`，卸载仅需结束该进程

---

## v1.7.0

### 新增
- **Windows Defender 实时保护管理**：禁用后停止实时扫描释放资源（仅改组策略键，不动服务，可完全恢复）；操作前检测篡改防护与实时保护状态，未关闭时弹窗引导前往安全中心
- **Windows 小组件（Widgets）开关**：通过 `Dsh` 组策略键禁用/恢复任务栏小组件按钮与面板
- **Microsoft Edge 安装/恢复**：卸载 Edge 后可通过官方安装器一键重装，同时清除阻止复活策略
- **更新检查多源回退**：官方 GitHub API 失败时自动尝试 ghproxy 系国内镜像，解决中国大陆用户无法访问 GitHub 导致获取不到更新的问题

### 变更
- **许可证从 GPL v3 变更为 MIT**
- **README 新增"设计理念"章节**：确立"安全、可逆、克制"原则——所有修改可恢复、不做破坏性操作、不提供 Defender 完全禁用（避免破坏 WdBoot 导致生物识别失效）
- **HomePage 重构**：移除"基于 revision-tool 二次开发"表述，改为鸣谢列表（灵感来源），许可证展示更新为 MIT
- i18n 清理：移除未使用键，压缩冗余文案

### 修复
- 小组件开关状态检测语义反转（"打开不了也恢复不了"）
- security.go 方法嵌套与冗余注释压缩

---

## v1.6.0

### 移除
- **Windows 安全中心(Defender)禁用功能整体下线**：移除安全模式全自动禁用/恢复方案（`PrepareDefenderDisable` / `PrepareDefenderRestore` / `SafeModeApply`），该功能依赖文件备份且可靠性不足
- 删除 scriptmgr 内嵌脚本包及 disableDefender/restoreDefender PowerShell 脚本
- 清理前端安全中心开关卡片、wailsjs 绑定及对应多语言文案

### 修复
- **RestartExplorer 重写**：采用 Dism++ 同款消息驱动方案（`PostMessage WM_EXITEXPLORER` 优雅退出 → 等待进程结束 → ShellExecute 重启），彻底移除 `taskkill` 调用，避免被安全软件行为拦截
- **UninstallAppx 增强**：增加 `Remove-AppxPackage -AllUsers` 与 `Remove-AppxProvisionedPackage -Online`，系统预装应用可彻底移除且不再随更新恢复
- 个性化设置模块多项修复：Edge 滑动手势状态判断、通知模式降级、快捷方式图标/盾牌常量修正、TaskManager 轮询超时调整
- 注册表 API 增加 `*E` 变体（返回 error），供需要错误传播的场景使用

### 其他
- `build.sh` 改为本地构建脚本，不再纳入版本库管理

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
