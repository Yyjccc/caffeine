package main

import (
	"caffeine/client"
	"caffeine/core"
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	clientApp := client.NewClientApp()

	// Store ctx in a variable accessible to callbacks
	var appCtx context.Context

	// 创建应用程序菜单
	appMenu := menu.NewMenu()

	// 文件菜单
	fileMenu := appMenu.AddSubmenu("文件")
	fileMenu.AddText("退出", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(appCtx)
	})

	// 设置菜单
	appMenu.AddSubmenu("设置")
	// 在这里可以添加设置相关的子菜单项

	// 主题菜单
	themeMenu := appMenu.AddSubmenu("主题")
	themeMenu.AddText("浅色", nil, func(_ *menu.CallbackData) {
		// 处理浅色主题切换
	})
	themeMenu.AddText("深色", nil, func(_ *menu.CallbackData) {
		// 处理深色主题切换
	})

	// 关于菜单
	aboutMenu := appMenu.AddSubmenu("关于")
	aboutMenu.AddText("版本信息", nil, func(_ *menu.CallbackData) {
		// 显示版本信息
	})

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ui",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{239, 239, 239, 1},
		OnStartup: func(ctx context.Context) {
			appCtx = ctx // Store the context
			core.AddHook(NewWailsLogHook(ctx))
		},
		Bind: []interface{}{
			app,
			clientApp,
		},
		Menu: appMenu, // 添加菜单到应用程序选项
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
