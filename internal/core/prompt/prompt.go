// Package prompt は翻訳プロンプトの組み立て（base 指示・口調指示・原文の合成）を行う純粋ルール。
// 翻訳 Run と結果取得の実プロンプト再構成が同じ関数を使い、同じ入力から同じプロンプトを得る。
package prompt

import (
	"strings"

	"aitranslationenginejp/internal/core/runtimetag"
	"aitranslationenginejp/internal/provider"
)

// ComposePrompt は base 指示・口調指示・機械置換済み原文から完成プロンプトを組む純粋関数。
// 翻訳実行（Run）と結果取得（api の実プロンプト再構成）の両方がこの 1 関数を使い、同じ入力から同じプロンプトを得る。
// personaDirective が空なら system は base 指示だけにする。source はそのまま user メッセージへ入れる。
// source に実行時タグ（<...>）があれば、system へタグ保護指示を足してモデルにタグの原形保持を求める。
func ComposePrompt(baseDirective, personaDirective, source string) provider.Prompt {
	system := baseDirective
	if strings.TrimSpace(personaDirective) != "" {
		system = baseDirective + "\n\n" + personaDirective
	}
	if runtimetag.HasTag(source) {
		system += "\n\n" + runtimetag.GuardInstruction()
	}
	return provider.Prompt{System: system, User: source}
}

// RenderPrompt は完成プロンプトを結果行の実プロンプト表示用の 1 文字列へ描く。
// 実際に送る system メッセージと user メッセージを役割見出し付きで連結し、口調指示の合成を目視で確かめられるようにする。
func RenderPrompt(p provider.Prompt) string {
	return "system:\n" + p.System + "\n\nuser:\n" + p.User
}
