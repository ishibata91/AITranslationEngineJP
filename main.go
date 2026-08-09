// AITranslationEngineJp の Wails app entry。composition root で配線した api.App を Bind する。
package main

import (
	"embed"
	"log"
	"path/filepath"

	"aitranslationenginejp/internal/bootstrap"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 観測ログの出力先を決める。dev 起動は標準エラー、配布ビルドは exe の隣のファイルへ残す
	//（observability-logging.md、logging_dev.go / logging_production.go）。
	setupLogging()

	// WebView2 の作業データの置き場所を、利用者ごとの保存先へ固定する。既定では exe の隣に作られるが、
	// 書き込みできない場所へ配布された場合や、起動元がファイル操作を横取りする場合に作れず、
	// 画面が出ないまま終了する。取得できない場合は空のままにして Wails の既定に任せる。
	webviewDataPath := ""
	if dir, err := appStateDir(); err == nil {
		webviewDataPath = filepath.Join(dir, "webview2")
	}

	app, closer, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("bootstrap 失敗: %v", err)
	}
	defer func() { _ = closer.Close() }()

	if err := wails.Run(&options.App{
		Title:  "AITranslationEngineJp",
		Width:  1200,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		Bind:      []any{app},
		Windows: &windows.Options{
			WebviewUserDataPath: webviewDataPath,
		},
	}); err != nil {
		log.Fatalf("wails 実行失敗: %v", err)
	}
}
