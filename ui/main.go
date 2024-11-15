package main

import (
	"caffeine/client"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	clientApp := client.NewClientApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ui",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{239, 239, 239, 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			clientApp,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
