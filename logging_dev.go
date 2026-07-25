//go:build !production

package main

import (
	"log/slog"
	"os"
)

// setupLogging は dev 起動の観測ログを標準エラーへ出す（observability-logging.md）。
// dev 起動は scripts/dev/run-wails.sh が標準出力と標準エラーを tmp/logs/wails-dev.log へ
// 落とすため、ここではファイルへ書き分けない。
func setupLogging() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
}
