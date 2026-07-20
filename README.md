# OptiWin

<p align="center">
  <img src="frontend/src/assets/logo.png" width="128" height="128" alt="OptiWin Logo">
</p>

<p align="center">
  适用于任何 Windows 系统的个性化调整工具箱
</p>

## 功能

- **安全** — Windows Defender、UAC、VBS、内存完整性设置
- **性能** — 电源计划、Superfetch、内存压缩、C-State、全屏/窗口优化、MPO
- **个性化** — 通知、上下文菜单、Explorer 设置、快捷方式外观
- **实用工具** — 休眠、快速启动、Windows 照片查看器、Edge 卸载、WebView2
- **更新** — 证书更新、KGL 更新、Windows 更新策略

## 基于 [meetrevision/revision-tool](https://github.com/meetrevision/revision-tool) 二次开发

## 许可证

本项目基于 **GNU General Public License v3.0** 开源。

## 构建

```bash
# 开发
cd frontend && npm run dev

# 编译 Linux
wails build -tags webkit2_41

# 交叉编译 Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ wails build -tags webkit2_41
```

## 技术栈

| 层 | 技术 |
|----|------|
| 后端 | Go + Wails |
| 前端 | Vue 3 + Naive UI |
| 注册表 | golang.org/x/sys/windows/registry |
