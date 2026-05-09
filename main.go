// Package main wires the desktop application entrypoint.
package main

import (
	"embed"
	"log"
	"log/slog"
	"os"

	"aitranslationenginejp/internal/bootstrap"
	infraruntime "aitranslationenginejp/internal/infra/runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var frontendAssets embed.FS

func main() {
	infraruntime.InstallDiagnosticLogger(os.Stderr, "backend", slog.LevelInfo)
	appController := bootstrap.NewAppController()
	appLifecycle := bootstrap.NewAppLifecycle(appController)

	err := wails.Run(&options.App{
		Title:  "AITranslationEngineJp",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: frontendAssets,
		},
		OnStartup:  appLifecycle.OnStartup,
		OnShutdown: appLifecycle.OnShutdown,
		Bind: []interface{}{
			appController,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
