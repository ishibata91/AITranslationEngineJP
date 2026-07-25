package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// exeDir は実行中の exe が置かれているフォルダを返す。
// 作業ディレクトリではなく exe の位置を基準にするのは、起動元（Mod Organizer など）が
// 作業ディレクトリを変える場合でも、実行ログの置き場所を exe の隣に定めるためである。
func exeDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("実行ファイルの場所を取得できない: %w", err)
	}
	return filepath.Dir(exe), nil
}

// appStateDir は利用者ごとの保存先（Windows は %LOCALAPPDATA%）配下の、本 app 専用フォルダを返す。
// WebView2 の作業データの置き場所に使う。exe の隣に置くと、書き込みできない場所へ配布された場合や
// 起動元がファイル操作を横取りする場合に起動できなくなるため、利用者ごとの保存先へ寄せる。
// 無ければ作る。
func appStateDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("利用者ごとの保存先を取得できない: %w", err)
	}
	dir := filepath.Join(base, "AITranslationEngineJp")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("保存先フォルダを作れない (%s): %w", dir, err)
	}
	return dir, nil
}
