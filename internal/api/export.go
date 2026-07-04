package api

import (
	"fmt"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// ExportXTranslatorXml は出力先フォルダを選ばせ、訳出済みの翻訳結果を xTranslator 互換 XML として書き出す。
// plugin ごとに 1 ファイルを出力先フォルダ直下へ書く。フォルダ選択・結果通知はいずれもネイティブダイアログで行う。
// 成功時は出力先と件数を情報ダイアログで、失敗時は原因をエラーダイアログで利用者へ知らせる。
// 利用者がフォルダ選択を閉じた場合は何もしない。書き出しに失敗した場合はダイアログ通知に加え error も返す。
func (a *App) ExportXTranslatorXml() error {
	dir, err := wailsruntime.OpenDirectoryDialog(a.baseCtx(), wailsruntime.OpenDialogOptions{
		Title: "書き出し先フォルダを選択",
	})
	if err != nil {
		return fmt.Errorf("書き出し先フォルダの選択: %w", err)
	}
	if dir == "" {
		return nil // 利用者がフォルダ選択を閉じた（キャンセル）。
	}

	paths, err := a.engine.ExportXTranslator(a.baseCtx(), dir)
	if err != nil {
		a.notify(wailsruntime.ErrorDialog, "書き出しに失敗しました", "書き出しに失敗しました。\n"+err.Error())
		return fmt.Errorf("xTranslator への書き出し: %w", err)
	}

	a.notify(wailsruntime.InfoDialog, "書き出しが完了しました",
		fmt.Sprintf("%d 件のファイルを次のフォルダへ書き出しました。\n%s", len(paths), dir))
	return nil
}

// notify は結果通知のネイティブダイアログを出す。ダイアログ表示自体の失敗は副次的で致命的でないため握る
// （書き出しの成否は呼び出し元の戻り値で扱う）。
func (a *App) notify(dialogType wailsruntime.DialogType, title, message string) {
	_, _ = wailsruntime.MessageDialog(a.baseCtx(), wailsruntime.MessageDialogOptions{ //nolint:errcheck // ダイアログ表示の失敗は副次的で無視する
		Type:    dialogType,
		Title:   title,
		Message: message,
	})
}
