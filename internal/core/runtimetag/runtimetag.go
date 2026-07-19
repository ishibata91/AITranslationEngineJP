// Package runtimetag は本文中の実行時タグ（<Alias=...> 等、ゲーム実行時に差し込まれるタグ）を
// 翻訳経路で守る純粋ルール。副作用を持たず、入出力で完結する。
//
// 守り方は 2 段。
//   - 辞書の機械置換から守る: タグを一時プレースホルダへ退避（Mask）して機械置換を通し、直後に原形へ戻す（Restore）。
//     置換がタグ内部の語（例 <Alias=Aventus> の Aventus）を書き換えるのを防ぐ。プレースホルダは AI へは送らない。
//   - AI からの改変・欠落を減らす・検出する: AI には意味の分かる生タグを見せ、保護指示を添える。
//     AI 出力に元のタグがそのまま残っているかを照合し、欠落数を返す（CountMissing）。
package runtimetag

import (
	"regexp"
	"strconv"
	"strings"
)

// tagPattern は保護対象の実行時タグ。山かっこで囲み、内側に山かっこを含まない 1 文字以上の並び。
// <Alias=...>・<Global=...>・<BribeCost> 等を捉える。空 <> や入れ子は捉えない。
var tagPattern = regexp.MustCompile(`<[^<>]+>`)

// placeholder は i 番目のタグへ割り当てる一時退避トークンを返す。
// 記号 ⟦⟧（U+27E6/U+27E7）は本文・辞書語にまず現れず、機械置換の照合に当たらない。AI へは送らず、Restore で必ず原形へ戻す。
func placeholder(i int) string {
	return "⟦" + strconv.Itoa(i) + "⟧"
}

// Mask は text 中の実行時タグを出現順に ⟦0⟧⟦1⟧… へ退避し、退避後テキストと退避したタグ列（出現順）を返す。
// 退避は機械置換（dictionary.Apply）からタグ内部を守るための一時措置で、AI へ送る前に Restore で戻す。
func Mask(text string) (masked string, tags []string) {
	masked = tagPattern.ReplaceAllStringFunc(text, func(tag string) string {
		ph := placeholder(len(tags))
		tags = append(tags, tag)
		return ph
	})
	return masked, tags
}

// Restore は text 中の退避トークン ⟦i⟧ を tags[i]（退避した原タグ）へ戻す。
// Mask→機械置換→Restore の順で使い、機械置換を通した本文へ生タグを復元してから AI へ送る。
func Restore(text string, tags []string) string {
	restored := text
	for i, tag := range tags {
		restored = strings.ReplaceAll(restored, placeholder(i), tag)
	}
	return restored
}

// Tags は text 中の実行時タグ（<...>）を出現順に返す。
// 決定的テストの fake provider が「タグを保持するモデル」を模すため、受け取った生タグをそのまま返すのに使う。
func Tags(text string) []string {
	return tagPattern.FindAllString(text, -1)
}

// HasTag は text に実行時タグ（<...>）が 1 つでも含まれるかを返す。
// プロンプト合成が、タグを持つ本文にだけタグ保護指示を付けるための判定に使う。
func HasTag(text string) bool {
	return tagPattern.MatchString(text)
}

// CountMissing は AI 出力 output に、元のタグ列 tags が原形で残っていない件数を返す。
// タグの出現回数で照合し、output 側の回数が足りない分を欠落として数える（モデルが削除・改変した実行時タグの数）。
func CountMissing(output string, tags []string) int {
	need := make(map[string]int, len(tags))
	for _, t := range tags {
		need[t]++
	}
	missing := 0
	for tag, n := range need {
		if got := strings.Count(output, tag); got < n {
			missing += n - got
		}
	}
	return missing
}

// GuardInstruction は生タグの保持を求めるプロンプト指示文を返す。
// モデルへ「<...> のタグを一字一句変えずに訳文へ残す」ことを求め、実行時タグの欠落・改変を減らす。
func GuardInstruction() string {
	return "本文中の <Alias=...> や <BribeCost> のような山かっこ < > で囲まれたタグは、ゲーム実行時に実際の値へ置き換わる印である。訳文へそのまま一字一句変えずに残すこと。タグの中身を翻訳したり、記号を変えたり、削除したりしてはいけない。"
}
