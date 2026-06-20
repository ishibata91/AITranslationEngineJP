package engine

import (
	"strings"

	"aitranslationenginejp/internal/provider"
)

// ComposePrompt は base 指示・口調指示・機械置換済み原文から完成プロンプトを組む純粋関数。
// 翻訳実行（Run）と結果取得（api の実プロンプト再構成）の両方がこの 1 関数を使い、同じ入力から同じプロンプトを得る。
// personaDirective が空なら system は base 指示だけにする。source はそのまま user メッセージへ入れる。
func ComposePrompt(baseDirective, personaDirective, source string) provider.Prompt {
	system := baseDirective
	if strings.TrimSpace(personaDirective) != "" {
		system = baseDirective + "\n\n" + personaDirective
	}
	return provider.Prompt{System: system, User: source}
}

// RenderPrompt は完成プロンプトを結果行の実プロンプト表示用の 1 文字列へ描く。
// 実際に送る system メッセージと user メッセージを役割見出し付きで連結し、口調指示の合成を目視で確かめられるようにする。
func RenderPrompt(p provider.Prompt) string {
	return "system:\n" + p.System + "\n\nuser:\n" + p.User
}
