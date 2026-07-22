//go:build !windows

package services

func UpdateCertificates() string { return "仅支持 Windows" }

func GetUpdateChannel() string               { return "retail" }
func SetUpdateChannel(channel string) bool    { return false }
func GetPauseUpdatesStatus() bool             { return false }
func EnablePauseUpdates() bool                { return false }
func DisablePauseUpdates() bool               { return false }
func GetVisibilityStatus() bool               { return false }
func EnableVisibility() bool                  { return false }
func DisableVisibility() bool                 { return false }
func GetDriverUpdatesStatus() bool            { return true }
func EnableDriverUpdates() bool               { return false }
func DisableDriverUpdates() bool              { return false }
func UpdateKGL() string                       { return "仅支持 Windows" }

func GetUltimatePerformanceStatus() bool      { return false }
func EnableUltimatePerformance() bool         { return false }
func DisableUltimatePerformance() bool        { return false }
func GetCStateStatus() bool                   { return true }
func EnableCState() bool                      { return false }
func DisableCState() bool                     { return false }
func GetSuperfetchStatus() bool               { return true }
func EnableSuperfetch() bool                  { return false }
func DisableSuperfetch() bool                 { return false }
func GetFullscreenOptimizationStatus() bool   { return true }
func EnableFullscreenOptimization() bool      { return false }
func DisableFullscreenOptimization() bool     { return false }
func GetWindowedOptimizationStatus() bool     { return true }
func EnableWindowedOptimization() bool        { return false }
func DisableWindowedOptimization() bool       { return false }
func GetMpoStatus() bool                      { return true }
func EnableMpo() bool                         { return false }
func DisableMpo() bool                        { return false }
func GetMemoryCompressionStatus() bool        { return true }
func EnableMemoryCompression() bool           { return false }
func DisableMemoryCompression() bool          { return false }
func ClearShaderCache() string                { return "仅支持 Windows" }
func GetGameBarStatus() bool                  { return false }
func EnableGameBar() bool                     { return false }
func DisableGameBar() bool                    { return false }

func GetSecurityHealthServiceStatus() bool    { return true }
func RestoreDefender() bool                   { return false }
func DisableAllServices() bool                { return false }
func GetUacStatus() bool                      { return true }
func EnableUac() bool                         { return false }
func DisableUac() bool                        { return false }
func GetVbsStatus() bool                      { return true }
func EnableVbs() bool                         { return false }
func DisableVbs() bool                        { return false }
func GetMemoryIntegrityStatus() bool          { return true }
func EnableMemoryIntegrity() bool             { return false }
func DisableMemoryIntegrity() bool            { return false }

func GetNotificationStatus() string           { return "开启" }
func SetNotificationMode(mode string) bool    { return false }
func GetLegacyBalloonStatus() bool            { return false }
func SetLegacyBalloon(on bool) bool           { return false }
func GetEdgeSwipeStatus() bool                { return true }
func SetEdgeSwipe(on bool) bool               { return false }
func GetNewContextMenuStatus() bool           { return true }
func SetNewContextMenu(on bool) bool          { return false }
func GetExplorerHomeStatus() bool             { return true }
func SetExplorerHome(on bool) bool            { return false }
func GetExplorerGalleryStatus() bool          { return true }
func SetExplorerGallery(on bool) bool         { return false }
func GetRemoveShortcutArrowStatus() bool      { return false }
func SetRemoveShortcutArrow(on bool) bool     { return false }
func GetRemoveShortcutTextStatus() bool       { return false }
func SetRemoveShortcutText(on bool) bool      { return false }
func GetRemoveShieldStatus() bool             { return false }
func SetRemoveShield(on bool) bool            { return false }
func GetOldTaskManagerStatus() bool           { return false }
func SetOldTaskManager(enable bool) bool      { return false }

func GetHibernateStatus() bool                { return false }
func EnableHibernate() bool                   { return false }
func DisableHibernate() bool                  { return false }
func GetFastStartupStatus() bool              { return true }
func EnableFastStartup() bool                 { return false }
func DisableFastStartup() bool                { return false }
func GetPhotoViewerStatus() bool              { return false }
func EnablePhotoViewer() bool                 { return false }
func DisablePhotoViewer() bool                { return false }
func UninstallEdge() string                   { return "仅支持 Windows" }
func GetWebView2Version() string              { return "未知" }
func InstallWebView2() string                 { return "仅支持 Windows" }
func SetSafeBoot(mode string) bool            { return false }
func RebootSystem() bool                      { return false }
func RebootToBios() bool                      { return false }
func GetSystemInfo() string                   { return `{"os":"","build":"","cpu":"","ram":"","ipv4":"","ipv6":""}` }

