package main

import (
	"embed"
	"log"

	controllerwails "aitranslationenginejp/internal/controller/wails"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var frontendAssets embed.FS

func main() {
	appController := controllerwails.NewAppController()

	err := wails.Run(&options.App{
		Title:  "AITranslationEngineJp",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: frontendAssets,
		},
		OnStartup: appController.OnStartup,
		OnShutdown: appController.OnShutdown,
		Bind: []interface{}{
			appController,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
