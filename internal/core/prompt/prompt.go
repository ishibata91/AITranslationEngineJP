// Package prompt は翻訳プロンプトの組み立て（base 指示・口調指示・原文の合成）を行う純粋ルール。
// 翻訳 Run と結果取得の実プロンプト再構成が同じ関数を使い、同じ入力から同じプロンプトを得る。
package prompt

import (
	"strings"

	"aitranslationenginejp/internal/core/runtimetag"
	"aitranslationenginejp/internal/provider"
)

// FillVariables は directive の指示文中の変数トークン（{name}）を vars の値へ差し込む純粋関数。
// MECE モデルの「directive ＝ 指示文 ＋ 実行時に埋める変数」の変数差し込みを 1 関数に閉じる。
// vars のキーはトークン全体（"{traits}" など）で渡す。未宣言のトークンはそのまま残す。
// 変数なしの directive（固有名・定型句・文体）は vars が空で、instruction をそのまま返す経路になる。副作用を持たない。
func FillVariables(instruction string, vars map[string]string) string {
	if len(vars) == 0 {
		return instruction
	}
	filled := instruction
	for token, value := range vars {
		filled = strings.ReplaceAll(filled, token, value)
	}
	return filled
}

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
