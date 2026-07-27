// Package personatone は注入入力（基底口調・性別・種族）から翻訳プロンプトの口調指示を組む純粋ルール。
// 基底口調の性質文・役割語（rolespeech）・種族訛りを織り合わせ、副作用を持たない。
package personatone

import (
	"strings"

	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/core/tone"
	"aitranslationenginejp/internal/model"
)

// traitsToken は口調指示テンプレートで話者の性質列を差し込む位置を表す差し込み口。
const traitsToken = "{traits}"

// toneTraits は基底口調セル（[感情段階][対人段階]）の性質文。翻訳プロンプトの口調指示へ入れる。
// implementation-scope.md のとおり v1 はコード定数で持つ。
// 一人称・語尾を確定する例文は assets/role-speech-examples.tsv が持ち、役割語と同じキーで引いて口調指示へ足す。
//
// 対人段階が丁寧の 3 セルから、敬語を求める語を外してある（2026-07-28）。
// 変える前は「敬語を保ち」「柔らかく丁寧な口調」「感情を込めた丁寧な口調」と書いていた。
// 実測で、柔和・気遣いと慇懃・端正の 2 セルだけが訳文の敬体 87.4% と 81.7% を作っており、
// 同じ行の公式日本語訳は 5.7% と 2.8% しか敬体を使っていなかった。日本語では丁寧な人物でも
// 常体で話せるので、丁寧さを言い表す語から敬体を求める部分だけを外す。
// 測定の記録は docs/exec-plans/active/dialogue-tone-naturalness/loop-log.md の回 3 が持つ。
var toneTraits = [3][3]string{
	// 抑制（感情を抑える）
	{
		"相手を見下す冷たい口調。感情を抑え、丁寧さを欠いた突き放した言い方にする。", // 尊大
		"感情を交えない事務的な口調。敬語にも砕けにも寄らず淡々と述べる。",      // 中立
		"礼儀正しく端正な口調。言葉を選び、感情を抑えて落ち着いて述べる。",      // 丁寧
	},
	// 中
	{
		"ぶっきらぼうで乱暴な口調。命令的に言い、相手を立てない。", // 尊大
		"飾らない率直な口調。過度な敬語も乱暴さもなく話す。",    // 中立
		"柔らかく気遣う口調。相手を立てて穏やかに述べる。",     // 丁寧
	},
	// 激情
	{
		"激しく威圧する口調。感情を露わにし、相手を罵り見下す。", // 尊大
		"感情を高ぶらせた率直な口調。勢いよく言い切る。",     // 中立
		"感情を込めて訴える口調。相手を立てつつ熱意を持って言う。",      // 丁寧
	},
}

// toneTraitOf は生成済みの基底口調（段階）から性質文を引く。段階が範囲外なら空（口調指示なし）。
func toneTraitOf(attitudeBand, emotionBand int) string {
	if emotionBand < 0 || emotionBand >= len(toneTraits) {
		return ""
	}
	row := toneTraits[emotionBand]
	if attitudeBand < 0 || attitudeBand >= len(row) {
		return ""
	}
	return row[attitudeBand]
}

// sexTrait は話者の性別から口調指示の 1 行を返す。役割語（一人称・言い回し）とは別の行にして、
// 性別そのものを翻訳モデルへ伝える。役割語テンプレートは成人男性と性別不明が同じ出力になるため、
// 引いた結果だけでは男女の区別がプロンプトに現れない。
// 性別を取れない話者（空・未知の値）は行を返さない。取れない話者へ既定の性別を当てない。
// 名指し話者と汎用台詞・PC 発話で同じ形の行を出す（呼ぶ側が違っても文面を変えない）。
func sexTrait(sex string) string {
	switch strings.ToLower(strings.TrimSpace(sex)) {
	case "male":
		return "- 性別: 男性。男性の話者として訳す。"
	case "female":
		return "- 性別: 女性。女性の話者として訳す。"
	default:
		return ""
	}
}

// raceMarkerTrait は種族 EditorID から種族訛りの注記を返す。基底口調へ重ねる語彙マーカー（R7、軸と直交）。
// 古風さは skyrim に一貫使用のキャラがおらず雑音のため入れない（persona-design.md）。未対応の種族は空。
func raceMarkerTrait(raceEDID string) string {
	switch {
	case strings.Contains(raceEDID, "Khajiit"):
		return "カジートの訛り。三人称で自称する（「この者は」など）。"
	case strings.Contains(raceEDID, "Argonian"):
		return "アルゴニアンの訛り。沼地の異邦人らしい硬さを残す。"
	default:
		return ""
	}
}

// BuildToneTraits は注入入力（生成済み基底口調＋性別＋種族）から口調指示の箇条書き行を組む。
// 性質文（基底口調）→ 性別 → 役割語（一人称・言い回し）→ 種族訛り の順に並べる。性質文が無ければ空（口調指示なし）。
// 役割語は roles テンプレートを 性別×年齢（race 由来）×セル で引く。roles が nil なら役割語は付けない。
func BuildToneTraits(in model.LinePersonaInput, roles *rolespeech.Table) []string {
	trait := toneTraitOf(in.AttitudeBand, in.EmotionBand)
	if trait == "" {
		return nil
	}
	traits := []string{"- 口調: " + trait}
	if line := sexTrait(in.Sex); line != "" {
		traits = append(traits, line)
	}
	traits = append(traits, roleSpeechLines(in, roles)...)
	if m := raceMarkerTrait(in.RaceEDID); m != "" {
		traits = append(traits, "- 種族訛り: "+m)
	}
	if line := lineEmotionLine(in.EmotionType); line != "" {
		traits = append(traits, line)
	}
	return traits
}

// roleSpeechLines は役割語テンプレートを引いて口調指示の行へ整える。
// 照合は決定経路に依らず、性別・race 由来の区分・基底口調セルだけで決める（フォールバック・保留でも付く）。
// 一致が無い、または一人称・言い回し・例文がすべて空なら行を返さない（役割語マーカーなし）。
func roleSpeechLines(in model.LinePersonaInput, roles *rolespeech.Table) []string {
	cell := tone.CellName(in.AttitudeBand, in.EmotionBand)
	tmpl, ok := roles.Lookup(rolespeech.RoleClassOfRace(in.RaceEDID), strings.ToLower(in.Sex), cell)
	if !ok {
		return nil
	}
	return formatRoleSpeech(tmpl)
}

// formatRoleSpeech は役割語テンプレート（一人称・言い回し・例文）を口調指示の行へ整える。
// 例文は説明文だけでは揺れる一人称・語尾を訳例で固定するために置く。
// 一人称も言い回しも例文も空なら行を返さない。名指し話者と汎用・PC で共通の整形。
func formatRoleSpeech(tmpl rolespeech.Template) []string {
	var lines []string
	switch {
	case tmpl.FirstPerson != "" && tmpl.Register != "":
		lines = append(lines, "- 人称と言い回し: 一人称は「"+tmpl.FirstPerson+"」。"+tmpl.Register)
	case tmpl.FirstPerson != "":
		lines = append(lines, "- 人称と言い回し: 一人称は「"+tmpl.FirstPerson+"」。")
	case tmpl.Register != "":
		lines = append(lines, "- 言い回し: "+tmpl.Register)
	}
	if !tmpl.Example.IsZero() {
		lines = append(lines, "- 例: "+tmpl.Example.Source+" → "+tmpl.Example.Dest)
	}
	return lines
}

// emotionAdvice は感情段階を口調指示へ重ねる助言文を返す。上書き命令でなく助言にとどめ、
// 翻訳モデルが本文から読み取る余地を残す。中（1）は雑音を避けるため助言を出さない。
func emotionAdvice(emotionBand int) string {
	switch emotionBand {
	case tone.EmotionSuppressed:
		return "落ち着いた、感情を抑えた調子にする。"
	case tone.EmotionIntense:
		return "感情が高ぶった、勢いのある調子にする。"
	default:
		return ""
	}
}

// lineEmotionWords は台詞感情型（TRDT）→ 日本語感情語。Neutral・型なしは持たず、加算しない。
// 対人段階×感情段階の基底口調（9 セル）とは別軸で、台詞ごとの情動をプロンプトへ重ねる。
var lineEmotionWords = map[string]string{
	tone.LineEmotionAnger:    "怒り",
	tone.LineEmotionDisgust:  "嫌悪",
	tone.LineEmotionFear:     "恐れ",
	tone.LineEmotionSad:      "悲しみ",
	tone.LineEmotionHappy:    "喜び",
	tone.LineEmotionSurprise: "驚き",
	tone.LineEmotionPuzzled:  "戸惑い",
}

// lineEmotionLine は台詞感情型から口調指示の 1 行を組む。Neutral・型なし・未知型は空（加算なし）。
// トーンは助言（上書きでなく助言）に揃え、翻訳モデルが本文からも情動を読む余地を残す。
func lineEmotionLine(emotionType string) string {
	word, ok := lineEmotionWords[emotionType]
	if !ok {
		return ""
	}
	return "- 感情: この台詞は" + word + "を込めた口調で話す。"
}

// freeRoleSpeechLines は性別だけから一人称・語尾・例文の行を引く（汎用・PC 用）。
// 汎用・PC は対人段階・セルを持たないため、年齢区分は成人、セルはワイルドカードで照合する。
// 性別が空でも引く。空文字は性別列のワイルドカード行に一致し、性別を取れない話者の既定行へ落ちる。
// 一致が無いなら行を返さない（一人称・語尾の指定なし）。
func freeRoleSpeechLines(sex string, roles *rolespeech.Table) []string {
	tmpl, ok := roles.Lookup(rolespeech.RoleClassOfRace(""), strings.ToLower(strings.TrimSpace(sex)), "")
	if !ok {
		return nil
	}
	return formatRoleSpeech(tmpl)
}

// BuildFreeToneTraits は汎用台詞・PC 発話の口調指示の箇条書きを組む。
// 自由記述の口調（baseText）→ 性別 → 台詞の感情（TRDT 種別、無ければ本文 1 行の感情段階の助言）
// → 性別の一人称・語尾・例文 の順に並べる。
// 性別の行は名指し話者（BuildToneTraits）と同じ形にする。性別を取れない話者と PC 性別の未設定は行を出さない。
// 性別が空でも役割語は引く（性別列のワイルドカード行へ落ちる）。
// baseText が空なら空（口調指示なし）。対人段階・セルは持たず、台詞の感情と性別だけを重ねる。
func BuildFreeToneTraits(baseText string, emotionBand int, sex string, roles *rolespeech.Table, emotionType string) []string {
	base := strings.TrimSpace(baseText)
	if base == "" {
		return nil
	}
	traits := []string{"- 口調: " + base}
	if line := sexTrait(sex); line != "" {
		traits = append(traits, line)
	}
	// TRDT 種別があれば台詞の感情行を優先し、本文推定の強度助言は出さない（二重回避）。無ければ従来の強度助言を出す。
	if line := lineEmotionLine(emotionType); line != "" {
		traits = append(traits, line)
	} else if adv := emotionAdvice(emotionBand); adv != "" {
		traits = append(traits, "- 感情: "+adv)
	}
	traits = append(traits, freeRoleSpeechLines(sex, roles)...)
	return traits
}

// BuildToneDirective は口調指示テンプレートの {traits} へ口調指示の箇条書きを差し込む。
// traits が空なら空文字を返し、ペルソナ指示を注入しない（既存の差し込み口 traitsToken を再利用）。
func BuildToneDirective(personaTemplate string, traits []string) string {
	if len(traits) == 0 {
		return ""
	}
	return strings.ReplaceAll(personaTemplate, traitsToken, strings.Join(traits, "\n"))
}

// BuildToneLabel は結果一覧の口調チップ用の短い要約。基底口調セル名を返す。段階が範囲外なら空文字。
func BuildToneLabel(in model.LinePersonaInput) string {
	cell := tone.CellName(in.AttitudeBand, in.EmotionBand)
	if cell == "" {
		return ""
	}
	return "口調: " + cell
}

// PersonaMetaOf は注入入力から、UI 表示用の基底口調セル名と性質文を引く（判定結果の表示）。
func PersonaMetaOf(in model.LinePersonaInput) (cell, trait string) {
	return tone.CellName(in.AttitudeBand, in.EmotionBand), toneTraitOf(in.AttitudeBand, in.EmotionBand)
}
