package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"OptiWin/services"
	"OptiWin/utils"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetCurrentVersion() string { return CurrentVersion }
func (a *App) GetSystemInfo() string     { return services.GetSystemInfo() }
func (a *App) GetProxyInfo() string      { return utils.GetProxyInfo() }

const StatusUpToDate = "same"

func (a *App) CheckUpdate() string {
	client := utils.GetHTTPClient()
	req, err := http.NewRequest("GET", "https://api.github.com/repos/YouRan4/OptiWin/releases/latest", nil)
	if err != nil {
		return "err:创建请求失败"
	}
	req.Header.Set("User-Agent", "OptiWin/1.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "err:无法连接 GitHub"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "err:服务器返回 " + fmt.Sprint(resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "err:读取响应失败"
	}

	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.Unmarshal(b, &release); err != nil {
		return "err:解析版本信息失败"
	}

	if release.TagName == "" {
		return "err:未获取到版本号"
	}

	if compareVersion(CurrentVersion, release.TagName) >= 0 {
		return StatusUpToDate
	}

	result, err := json.Marshal(map[string]string{
		"version": release.TagName,
		"body":    release.Body,
		"url":     release.HTMLURL,
	})
	if err != nil {
		return "err:序列化失败"
	}
	return string(result)
}

// --- 更新通道 ---
func (a *App) GetUpdateChannel() string             { return services.GetUpdateChannel() }
func (a *App) SetUpdateChannel(channel string) bool { return services.SetUpdateChannel(channel) }

// --- 暂停更新 ---
func (a *App) GetPauseUpdatesStatus() bool { return services.GetPauseUpdatesStatus() }
func (a *App) EnablePauseUpdates() bool    { return services.EnablePauseUpdates() }
func (a *App) DisablePauseUpdates() bool   { return services.DisablePauseUpdates() }

// --- 更新页面可见性 ---
func (a *App) GetVisibilityStatus() bool { return services.GetVisibilityStatus() }
func (a *App) EnableVisibility() bool    { return services.EnableVisibility() }
func (a *App) DisableVisibility() bool   { return services.DisableVisibility() }

// --- 驱动程序更新 ---
func (a *App) GetDriverUpdatesStatus() bool { return services.GetDriverUpdatesStatus() }
func (a *App) EnableDriverUpdates() bool    { return services.EnableDriverUpdates() }
func (a *App) DisableDriverUpdates() bool   { return services.DisableDriverUpdates() }

// --- 证书 & KGL ---
func (a *App) UpdateCertificates() string { return services.UpdateCertificates() }
func (a *App) UpdateKGL() string          { return services.UpdateKGL() }

// --- 电源计划 ---
func (a *App) GetUltimatePerformanceStatus() bool { return services.GetUltimatePerformanceStatus() }
func (a *App) EnableUltimatePerformance() bool    { return services.EnableUltimatePerformance() }
func (a *App) DisableUltimatePerformance() bool   { return services.DisableUltimatePerformance() }

// --- C-State ---
func (a *App) GetCStateStatus() bool { return services.GetCStateStatus() }
func (a *App) EnableCState() bool    { return services.EnableCState() }
func (a *App) DisableCState() bool   { return services.DisableCState() }

// --- Superfetch ---
func (a *App) GetSuperfetchStatus() bool { return services.GetSuperfetchStatus() }
func (a *App) EnableSuperfetch() bool    { return services.EnableSuperfetch() }
func (a *App) DisableSuperfetch() bool   { return services.DisableSuperfetch() }

// --- 全屏优化 ---
func (a *App) GetFullscreenOptimizationStatus() bool {
	return services.GetFullscreenOptimizationStatus()
}
func (a *App) EnableFullscreenOptimization() bool  { return services.EnableFullscreenOptimization() }
func (a *App) DisableFullscreenOptimization() bool { return services.DisableFullscreenOptimization() }

// --- 窗口优化 ---
func (a *App) GetWindowedOptimizationStatus() bool { return services.GetWindowedOptimizationStatus() }
func (a *App) EnableWindowedOptimization() bool    { return services.EnableWindowedOptimization() }
func (a *App) DisableWindowedOptimization() bool   { return services.DisableWindowedOptimization() }

// --- MPO ---
func (a *App) GetMpoStatus() bool { return services.GetMpoStatus() }
func (a *App) EnableMpo() bool    { return services.EnableMpo() }
func (a *App) DisableMpo() bool   { return services.DisableMpo() }

// --- 内存压缩 ---
func (a *App) GetMemoryCompressionStatus() bool { return services.GetMemoryCompressionStatus() }
func (a *App) EnableMemoryCompression() bool    { return services.EnableMemoryCompression() }
func (a *App) DisableMemoryCompression() bool   { return services.DisableMemoryCompression() }

// --- 着色器缓存 ---
func (a *App) ClearShaderCache() string { return services.ClearShaderCache() }

// --- Game Bar ---
func (a *App) GetGameBarStatus() bool { return services.GetGameBarStatus() }
func (a *App) EnableGameBar() bool    { return services.EnableGameBar() }
func (a *App) DisableGameBar() bool   { return services.DisableGameBar() }

// --- Defender ---
func (a *App) GetSecurityHealthServiceStatus() bool { return services.GetSecurityHealthServiceStatus() }
func (a *App) RestoreDefender() bool                { return services.RestoreDefender() }
func (a *App) DisableAllServices() bool             { return services.DisableAllServices() }

// --- UAC ---
func (a *App) GetUacStatus() bool { return services.GetUacStatus() }
func (a *App) EnableUac() bool    { return services.EnableUac() }
func (a *App) DisableUac() bool   { return services.DisableUac() }

// --- VBS ---
func (a *App) GetVbsStatus() bool { return services.GetVbsStatus() }
func (a *App) EnableVbs() bool    { return services.EnableVbs() }
func (a *App) DisableVbs() bool   { return services.DisableVbs() }

// --- MemoryIntegrity ---
func (a *App) GetMemoryIntegrityStatus() bool { return services.GetMemoryIntegrityStatus() }
func (a *App) EnableMemoryIntegrity() bool    { return services.EnableMemoryIntegrity() }
func (a *App) DisableMemoryIntegrity() bool   { return services.DisableMemoryIntegrity() }

// --- 通知 ---
func (a *App) GetNotificationStatus() string        { return services.GetNotificationStatus() }
func (a *App) SetNotificationMode(mode string) bool { return services.SetNotificationMode(mode) }

// --- 旧版气球通知 ---
func (a *App) GetLegacyBalloonStatus() bool  { return services.GetLegacyBalloonStatus() }
func (a *App) SetLegacyBalloon(on bool) bool { return services.SetLegacyBalloon(on) }

// --- 边缘滑动 ---
func (a *App) GetEdgeSwipeStatus() bool  { return services.GetEdgeSwipeStatus() }
func (a *App) SetEdgeSwipe(on bool) bool { return services.SetEdgeSwipe(on) }

// --- 新版上下文菜单 ---
func (a *App) GetNewContextMenuStatus() bool  { return services.GetNewContextMenuStatus() }
func (a *App) SetNewContextMenu(on bool) bool { return services.SetNewContextMenu(on) }

// --- Explorer 主页/图库 ---
func (a *App) GetExplorerHomeStatus() bool     { return services.GetExplorerHomeStatus() }
func (a *App) SetExplorerHome(on bool) bool    { return services.SetExplorerHome(on) }
func (a *App) GetExplorerGalleryStatus() bool  { return services.GetExplorerGalleryStatus() }
func (a *App) SetExplorerGallery(on bool) bool { return services.SetExplorerGallery(on) }

// --- 快捷方式 ---
func (a *App) GetRemoveShortcutArrowStatus() bool  { return services.GetRemoveShortcutArrowStatus() }
func (a *App) SetRemoveShortcutArrow(on bool) bool { return services.SetRemoveShortcutArrow(on) }
func (a *App) GetRemoveShortcutTextStatus() bool   { return services.GetRemoveShortcutTextStatus() }
func (a *App) SetRemoveShortcutText(on bool) bool  { return services.SetRemoveShortcutText(on) }
func (a *App) GetRemoveShieldStatus() bool         { return services.GetRemoveShieldStatus() }
func (a *App) SetRemoveShield(on bool) bool        { return services.SetRemoveShield(on) }

// --- 任务管理器 ---
func (a *App) GetOldTaskManagerStatus() bool { return services.GetOldTaskManagerStatus() }
func (a *App) SetOldTaskManager(enable bool) bool { return services.SetOldTaskManager(enable) }

// --- 休眠 ---
func (a *App) GetHibernateStatus() bool { return services.GetHibernateStatus() }
func (a *App) EnableHibernate() bool    { return services.EnableHibernate() }
func (a *App) DisableHibernate() bool   { return services.DisableHibernate() }

// --- 快速启动 ---
func (a *App) GetFastStartupStatus() bool { return services.GetFastStartupStatus() }
func (a *App) EnableFastStartup() bool    { return services.EnableFastStartup() }
func (a *App) DisableFastStartup() bool   { return services.DisableFastStartup() }

// --- 照片查看器 ---
func (a *App) GetPhotoViewerStatus() bool { return services.GetPhotoViewerStatus() }
func (a *App) EnablePhotoViewer() bool    { return services.EnablePhotoViewer() }
func (a *App) DisablePhotoViewer() bool   { return services.DisablePhotoViewer() }

// --- Edge / WebView2 ---
func (a *App) UninstallEdge() string      { return services.UninstallEdge() }
func (a *App) GetWebView2Version() string { return services.GetWebView2Version() }
func (a *App) InstallWebView2() string    { return services.InstallWebView2() }

// --- 安全模式 & 重启 ---
func (a *App) SetSafeBoot(mode string) bool { return services.SetSafeBoot(mode) }
func (a *App) RebootSystem() bool           { return services.RebootSystem() }
func (a *App) RebootToBios() bool           { return services.RebootToBios() }

// --- 系统主题色 ---
func (a *App) GetSystemAccentColor() string {
	return utils.GetSystemAccentColor()
}

// --- 重启资源管理器 ---
func (a *App) RestartExplorer() string {
	cmd := exec.Command("taskkill", "/f", "/im", "explorer.exe")
	utils.HideWindow(cmd)
	cmd.Run()
	go func() {
		cmd2 := exec.Command("explorer.exe")
		utils.HideWindow(cmd2)
		cmd2.Start()
	}()
	return "资源管理器已重启"
}
