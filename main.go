package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "OptiWin",
		Width:  1280,
		Height: 800,
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
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
