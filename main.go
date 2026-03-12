//go:build darwin || windows || linux

package main

import (
	"Content-Alchemist/backend"
	"embed"
	"fmt"
	"os"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// version 会被 ldflags 注入
var version = "dev"

func main() {
	// 设置多线程模式
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 创建应用实例
	app := backend.NewApp()

	// 应用配置
	opts := &options.App{
		Title:     "Content Alchemist",
		Width:     1400,
		Height:    900,
		MinWidth:  800,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		// macOS 特定配置
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 true,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.DefaultAppearance,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   fmt.Sprintf("Content Alchemist %s", version),
				Message: "© 2024 Content Alchemist Team\n本地优先的技术写作编辑器",
				Icon:    nil,
			},
		},
		// Windows 特定配置
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:          windows.RGB(30, 30, 30),
				DarkModeTitleBarInactive:  windows.RGB(40, 40, 40),
				DarkModeBorder:            windows.RGB(50, 50, 50),
				DarkModeBorderInactive:    windows.RGB(60, 60, 60),
				LightModeTitleBar:         windows.RGB(245, 245, 245),
				LightModeTitleBarInactive: windows.RGB(230, 230, 230),
				LightModeBorder:           windows.RGB(220, 220, 220),
				LightModeBorderInactive:   windows.RGB(200, 200, 200),
			},
		},
		// Linux 特定配置
		Linux: &linux.Options{
			Icon:                nil,
			WindowIsTranslucent: false,
		},
	}

	// 运行应用
	if err := wails.Run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
