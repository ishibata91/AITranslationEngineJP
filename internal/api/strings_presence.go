package api

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// StringsPresenceView は対象 plugin が置かれた Data フォルダにある Strings ファイルの言語別有無。
// 翻訳実行画面の警告（片側 Strings 欠け＝既存訳の対を作れず全文 AI 翻訳になる）の判定材料。
// 既訳の供給は Data フォルダ全 plugin の走査で立つため、警告が指す状態は
// 「Data フォルダに片方の言語の Strings が 1 本も無い」であり、選んだ plugin 自身の Strings の有無ではない。
type StringsPresenceView struct {
	English  bool `json:"english"`
	Japanese bool `json:"japanese"`
}

// stringsExtensions は Strings ファイルの拡張子 3 種（通常・DL・IL）。
var stringsExtensions = map[string]bool{".strings": true, ".dlstrings": true, ".ilstrings": true}

// GetStringsPresence は plugin と同じ Data フォルダ配下の strings/ を調べ、english / japanese の
// Strings ファイルの有無を返す。抽出前の判定のため plugin の中身は読まない（ファイル名だけで判定する）。
// strings フォルダ自体が無い場合は両方 false を返す（エラーにしない。非 localized 構成もあり得るため）。
func (a *App) GetStringsPresence(pluginPath string) (StringsPresenceView, error) {
	dir := stringsDir(filepath.Dir(pluginPath))
	view := StringsPresenceView{}
	if dir == "" {
		return view, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return view, nil //nolint:nilerr // 読めないフォルダは「Strings 無し」と同じ扱いにする（警告表示の判定材料のため）。
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	view = classifyStringsPresence(names)
	// 実行時にしか確定しない判定結果を残す（画面警告が出た・出ない理由の切り分け用）。
	slog.InfoContext(a.baseCtx(), "Strings の言語別有無を判定した",
		slog.String("event", "strings_presence"),
		slog.String("where", "api.GetStringsPresence"),
		slog.Bool("english", view.English),
		slog.Bool("japanese", view.Japanese),
	)
	return view, nil
}

// stringsDir は Data フォルダ配下の Strings フォルダのパスを返す。大文字小文字の揺れ
// （strings / Strings）を許し、どちらも無ければ空を返す。
func stringsDir(dataFolder string) string {
	for _, name := range []string{"strings", "Strings"} {
		p := filepath.Join(dataFolder, name)
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p
		}
	}
	return ""
}

// classifyStringsPresence はファイル名群から english / japanese の Strings の有無を判定する純粋関数。
// 対象は `*_english.<ext>` / `*_japanese.<ext>`（ext は .strings / .dlstrings / .ilstrings、大文字小文字は区別しない）。
func classifyStringsPresence(names []string) StringsPresenceView {
	view := StringsPresenceView{}
	for _, name := range names {
		lower := strings.ToLower(name)
		ext := filepath.Ext(lower)
		if !stringsExtensions[ext] {
			continue
		}
		base := strings.TrimSuffix(lower, ext)
		switch {
		case strings.HasSuffix(base, "_english"):
			view.English = true
		case strings.HasSuffix(base, "_japanese"):
			view.Japanese = true
		}
	}
	return view
}
