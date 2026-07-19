// Package runtimetag は本文中の実行時タグ（<Alias=...> 等、ゲーム実行時に差し込まれるタグ）を
// AI 翻訳の前後で守る純粋ルール。副作用を持たず、入出力で完結する。
// タグを翻訳前にプレースホルダへ退避（Mask）し、翻訳後に原形へ復元（Unmask）する。
// 復元時に、モデルが落とした・書き換えたタグ（プレースホルダが出力に残らなかったもの）の件数を返す。
package runtimetag

import (
	"regexp"
	"strconv"
	"strings"
)

// tagPattern は保護対象の実行時タグ。山かっこで囲み、内側に山かっこを含まない 1 文字以上の並び。
// <Alias=...>・<Global=...> 等を捉える。空 <> や入れ子は捉えない。
var tagPattern = regexp.MustCompile(`<[^<>]+>`)

// placeholderPattern は退避後のプレースホルダ（⟦連番⟧）。Placeholders が出力から拾うのに使う。
var placeholderPattern = regexp.MustCompile(`⟦\d+⟧`)

// placeholder は i 番目のタグへ割り当てる退避トークンを返す。
// 記号 ⟦⟧（U+27E6/U+27E7）は本文・ゲームテキストにまず現れず、翻訳モデルが素通ししやすい。
// 万一モデルが書き換えても Unmask が欠落として検出する。
func placeholder(i int) string {
	return "⟦" + strconv.Itoa(i) + "⟧"
}

// Mask は text 中の実行時タグを出現順に ⟦0⟧⟦1⟧… へ退避し、退避後テキストと退避したタグ列（出現順）を返す。
// タグが無ければ text をそのまま返し、tags は nil。退避により、AI 翻訳も辞書機械置換もタグ内部へ触れなくなる。
func Mask(text string) (masked string, tags []string) {
	masked = tagPattern.ReplaceAllStringFunc(text, func(tag string) string {
		ph := placeholder(len(tags))
		tags = append(tags, tag)
		return ph
	})
	return masked, tags
}

// Unmask は text 中のプレースホルダ ⟦i⟧ を tags[i]（退避した原タグ）へ戻し、復元後テキストと欠落数を返す。
// 欠落数は、出力にプレースホルダが残らなかったタグの件数（モデルが削除・改変した実行時タグの数）。
// プレースホルダの並び替えや重複が起きても、連番で照合するため位置に依らず正しく復元する。
func Unmask(text string, tags []string) (restored string, lost int) {
	restored = text
	for i, tag := range tags {
		ph := placeholder(i)
		if !strings.Contains(restored, ph) {
			lost++ // このタグのプレースホルダが出力に残っていない＝実行時タグが失われた。
			continue
		}
		restored = strings.ReplaceAll(restored, ph, tag)
	}
	return restored, lost
}

// Placeholders は text 中に現れる退避トークン（⟦連番⟧）を出現順に返す。
// 決定的テストの fake provider が「タグを保持するモデル」を模すために、受け取った退避トークンをそのまま返すのに使う。
func Placeholders(text string) []string {
	return placeholderPattern.FindAllString(text, -1)
}

// HasPlaceholder は text に退避トークン（⟦連番⟧）が 1 つでも含まれるかを返す。
// プロンプト合成が、退避したタグを持つ本文にだけタグ保護指示を付けるための判定に使う。
func HasPlaceholder(text string) bool {
	return placeholderPattern.MatchString(text)
}

// GuardInstruction は退避トークンを保持させるためのプロンプト指示文を返す。
// モデルへ「⟦連番⟧ のトークンを一字一句変えずに訳文へ残す」ことを求め、実行時タグの欠落自体を減らす。
func GuardInstruction() string {
	return "本文中の ⟦0⟧ ⟦1⟧ のような角かっこ付きトークンは、ゲーム実行時に置き換わる印である。訳文へそのまま一字一句変えずに残すこと。翻訳・削除・記号の変更をしてはいけない。"
}
