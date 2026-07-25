//go:build production

package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// setupLogging は配布ビルドの実行ログをファイルへ残す。
//
// 配布ビルドは console を持たないため、標準出力と標準エラーの書き込み先が無い。画面が出る前に
// 落ちる事象（WebView2 の起動失敗など）は、ファイルに残さないと後から原因を追えない。
//
// slog だけでなく標準出力と標準エラーも同じファイルへ向ける。Wails と WebView2 の失敗の一部は
// slog を通さず fmt で直接書くため（go-webview2 の `[WebView2 Error]` 行など）、両方を拾う。
//
// 置き場所は exe の隣にする。配布フォルダを開けばログが並んでいる形にして、利用者が渡しやすくする。
// 場所を取得できない場合や開けない場合は標準エラーへ出す従来どおりの動きに戻し、起動そのものは止めない。
func setupLogging() {
	dir, err := exeDir()
	if err != nil {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
		slog.Warn("実行ログのファイル出力を用意できない", "where", "main.setupLogging", "reason", err.Error())
		return
	}

	path := filepath.Join(dir, "app.log")
	// 追記で開く。起動ごとに区切り行を書き、どの起動の記録かを読み分けられるようにする。
	// プロセスの生存期間ずっと使うため閉じない。
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
		slog.Warn("実行ログのファイルを開けない", "where", "main.setupLogging", "reason", err.Error())
		return
	}

	os.Stdout = f
	os.Stderr = f
	slog.SetDefault(slog.New(slog.NewJSONHandler(f, nil)))
	fmt.Fprintf(f, "=== start %s ===\n", time.Now().Format(time.RFC3339))
	// 画面が出ずに終わる代表的な原因を、ログを開いた人がその場で解けるよう見出しへ添える。
	// Mod Organizer 経由の起動では、仮想ファイルシステムが WebView2 の子プロセスへ hook を掛け、
	// Chromium 側の hook と衝突してコントローラ生成が失敗する（docs/build-windows.md）。
	fmt.Fprintln(f, "（画面が出ずに終わる場合: Mod Organizer 経由なら msedgewebview2.exe を"+
		" MO2 の実行ファイル ブラックリストへ追加する。詳細は docs/build-windows.md）")
}
