package main

import (
	"caffeine/client"
	"caffeine/core"
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

//func main() {
//	// Create an instance of the app structure
//	app := application.New()
//	clientApp := client.NewClientApp()
//
//	// Store ctx in a variable accessible to callbacks
//	var appCtx context.Context
//
//	// 创建应用程序菜单
//	appMenu := menu.NewMenu()
//
//	// 文件菜单
//	fileMenu := appMenu.AddSubmenu("文件")
//	fileMenu.AddText("退出", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
//		runtime.Quit(appCtx)
//	})
//
//	// 设置菜单
//	appMenu.AddSubmenu("设置")
//	// 在这里可以添加设置相关的子菜单项
//
//	// 主题菜单
//	themeMenu := appMenu.AddSubmenu("主题")
//	themeMenu.AddText("浅色", nil, func(_ *menu.CallbackData) {
//		// 处理浅色主题切换
//	})
//	themeMenu.AddText("深色", nil, func(_ *menu.CallbackData) {
//		// 处理深色主题切换
//	})
//
//	// 关于菜单
//	aboutMenu := appMenu.AddSubmenu("关于")
//	aboutMenu.AddText("版本信息", nil, func(_ *menu.CallbackData) {
//		// 显示版本信息
//	})
//
//	// Create application with options
//	err := wails.Run(&options.App{
//		Title:  "ui",
//		Width:  1024,
//		Height: 768,
//		AssetServer: &assetserver.Options{
//			Assets: assets,
//		},
//		BackgroundColour: &options.RGBA{239, 239, 239, 1},
//		OnStartup: func(ctx context.Context) {
//			appCtx = ctx // Store the context
//			core.AddHook(NewWailsLogHook(ctx))
//		},
//		Bind: []interface{}{
//			app,
//			clientApp,
//		},
//		Menu: appMenu, // 添加菜单到应用程序选项
//	})
//
//	if err != nil {
//		println("Error:", err.Error())
//	}
//}

func main() {

	clientApp := client.GetClientApp()
	UiApp := NewApp(nil)
	app := application.New(application.Options{
		Name:        "caffeine",
		Description: "A tool of webshell management",
		Services: []application.Service{
			application.NewService(clientApp),
			application.NewService(UiApp),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	UiApp.ctx = app
	//添加日志回调
	core.AddHook(NewWailsLogHook(app))

	//设置菜单
	menu := app.NewMenu()
	//主菜单

	//appMenu := application.NewAppMenu()

	menu.AddRole(application.AppMenu)
	menu.AddRole(application.FileMenu)
	menu.AddRole(application.ViewMenu)
	menu.AddRole(application.EditMenu)

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "caffeine",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGBA(239, 239, 239, 1),
		Windows: application.WindowsWindow{
			Menu: menu,
		},
		URL: "/",
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
