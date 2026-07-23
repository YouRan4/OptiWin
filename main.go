package main

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	os.MkdirAll(filepath.Join(os.TempDir(), "optiwin"), 0755)
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "OptiWin",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless: true, // 无边框窗口
		Windows: &windows.Options{
			WebviewIsTransparent: true,         // WebView2 透明背景
			WindowIsTranslucent:  true,         // 窗口半透明
			BackdropType:         windows.Mica, // 亚克力毛玻璃效果
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("错误:", err.Error())
	}
	os.RemoveAll(filepath.Join(os.TempDir(), "optiwin"))
}
