//go:build !windows

package main

// 非 Windows 平台空实现

func (a *App) GetPauseUpdatesStatus() bool          { return false }
func (a *App) EnablePauseUpdates() bool             { return false }
func (a *App) DisablePauseUpdates() bool            { return false }
func (a *App) GetVisibilityStatus() bool            { return false }
func (a *App) EnableVisibility() bool               { return false }
func (a *App) DisableVisibility() bool              { return false }
func (a *App) GetDriverUpdatesStatus() bool         { return true }
func (a *App) EnableDriverUpdates() bool            { return false }
func (a *App) DisableDriverUpdates() bool           { return false }
func (a *App) GetCStateStatus() bool                { return true }
func (a *App) EnableCState() bool                   { return false }
func (a *App) DisableCState() bool                  { return false }
func (a *App) GetUltimatePerformanceStatus() bool   { return false }
func (a *App) EnableUltimatePerformance() bool      { return false }
func (a *App) DisableUltimatePerformance() bool     { return false }
func (a *App) GetSuperfetchStatus() bool            { return true }
func (a *App) EnableSuperfetch() bool               { return false }
func (a *App) DisableSuperfetch() bool              { return false }
func (a *App) GetFullscreenOptimizationStatus() bool { return true }
func (a *App) EnableFullscreenOptimization() bool   { return false }
func (a *App) DisableFullscreenOptimization() bool  { return false }
func (a *App) GetWindowedOptimizationStatus() bool  { return true }
func (a *App) EnableWindowedOptimization() bool     { return false }
func (a *App) DisableWindowedOptimization() bool    { return false }
func (a *App) GetMpoStatus() bool                   { return true }
func (a *App) EnableMpo() bool                      { return false }
func (a *App) DisableMpo() bool                     { return false }
func (a *App) GetMemoryCompressionStatus() bool     { return true }
func (a *App) EnableMemoryCompression() bool        { return false }
func (a *App) DisableMemoryCompression() bool       { return false }
func (a *App) GetDefenderStatus() bool             { return true }
func (a *App) EnableDefender() bool                { return false }
func (a *App) DisableDefender() bool               { return false }
func (a *App) GetUacStatus() bool                  { return true }
func (a *App) EnableUac() bool                     { return false }
func (a *App) DisableUac() bool                    { return false }
func (a *App) GetVbsStatus() bool                  { return true }
func (a *App) EnableVbs() bool                     { return false }
func (a *App) DisableVbs() bool                    { return false }
func (a *App) GetMemoryIntegrityStatus() bool      { return true }
func (a *App) EnableMemoryIntegrity() bool         { return false }
func (a *App) DisableMemoryIntegrity() bool        { return false }
func (a *App) GetHibernateStatus() bool            { return false }
func (a *App) EnableHibernate() bool               { return false }
func (a *App) DisableHibernate() bool              { return false }
func (a *App) GetFastStartupStatus() bool           { return true }
func (a *App) EnableFastStartup() bool              { return false }
func (a *App) DisableFastStartup() bool             { return false }
func (a *App) GetNotificationStatus() string        { return "开启" }
func (a *App) SetNotificationMode(mode string) bool { return false }
func (a *App) GetLegacyBalloonStatus() bool         { return false }
func (a *App) SetLegacyBalloon(on bool) bool        { return false }
func (a *App) GetEdgeSwipeStatus() bool             { return true }
func (a *App) SetEdgeSwipe(on bool) bool            { return false }
func (a *App) GetNewContextMenuStatus() bool        { return true }
func (a *App) SetNewContextMenu(on bool) bool       { return false }
func (a *App) GetExplorerHomeStatus() bool          { return true }
func (a *App) SetExplorerHome(on bool) bool         { return false }
func (a *App) GetExplorerGalleryStatus() bool       { return true }
func (a *App) SetExplorerGallery(on bool) bool      { return false }
func (a *App) GetRemoveShortcutArrowStatus() bool  { return false }
func (a *App) SetRemoveShortcutArrow(on bool) bool { return false }
func (a *App) GetRemoveShortcutTextStatus() bool   { return false }
func (a *App) SetRemoveShortcutText(on bool) bool  { return false }
func (a *App) GetRemoveShieldStatus() bool         { return false }
func (a *App) SetRemoveShield(on bool) bool        { return false }
func (a *App) GetPhotoViewerStatus() bool          { return false }
func (a *App) EnablePhotoViewer() bool             { return false }
func (a *App) DisablePhotoViewer() bool            { return false }
func (a *App) UninstallEdge() string                 { return "仅支持 Windows" }
func (a *App) GetWebView2Version() string             { return "未知" }
func (a *App) InstallWebView2() string               { return "仅支持 Windows" }
func (a *App) UpdateKGL() string                     { return "仅支持 Windows" }
