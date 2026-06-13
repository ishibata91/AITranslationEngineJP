// AITranslationEngineJp の Wails app entry。composition root で配線した api.App を Bind する。
package main

import (
	"embed"
	"log"

	"aitranslationenginejp/internal/bootstrap"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app, st, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("bootstrap 失敗: %v", err)
	}
	defer func() { _ = st.Close() }()

	if err := wails.Run(&options.App{
		Title:  "AITranslationEngineJp",
		Width:  1200,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		Bind:      []any{app},
	}); err != nil {
		log.Fatalf("wails 実行失敗: %v", err)
	}
}
