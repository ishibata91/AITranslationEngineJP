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
// implementation-scope.md のとおり v1 はコード定数で持つ。一人称・語尾を確定する few-shot 例文は後続の精緻化で足す。
var toneTraits = [3][3]string{
	// 抑制（感情を抑える）
	{
		"相手を見下す冷たい口調。感情を抑え、丁寧さを欠いた突き放した言い方にする。", // 尊大
		"感情を交えない事務的な口調。敬語にも砕けにも寄らず淡々と述べる。",      // 中立
		"礼儀正しく端正な口調。敬語を保ち、感情を抑えて落ち着いて述べる。",      // 丁寧
	},
	// 中
	{
		"ぶっきらぼうで乱暴な口調。命令的に言い、相手を立てない。", // 尊大
		"飾らない率直な口調。過度な敬語も乱暴さもなく話す。",    // 中立
		"柔らかく丁寧な口調。相手を立てて穏やかに述べる。",     // 丁寧
	},
	// 激情
	{
		"激しく威圧する口調。感情を露わにし、相手を罵り見下す。", // 尊大
		"感情を高ぶらせた率直な口調。勢いよく言い切る。",     // 中立
		"感情を込めた丁寧な口調。熱意を持って訴える。",      // 丁寧
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
// 性質文（基底口調）→ 役割語（一人称・言い回し）→ 種族訛り の順に並べる。性質文が無ければ空（口調指示なし）。
// 役割語は roles テンプレートを 性別×年齢（race 由来）×セル で引く。roles が nil なら役割語は付けない。
func BuildToneTraits(in model.LinePersonaInput, roles *rolespeech.Table) []string {
	trait := toneTraitOf(in.AttitudeBand, in.EmotionBand)
	if trait == "" {
		return nil
	}
	traits := []string{"- 口調: " + trait}
	if line := roleSpeechLine(in, roles); line != "" {
		traits = append(traits, line)
	}
	if m := raceMarkerTrait(in.RaceEDID); m != "" {
		traits = append(traits, "- 種族訛り: "+m)
	}
	return traits
}

// roleSpeechLine は役割語テンプレートを引いて口調指示の 1 行へ整える。
// 照合は決定経路に依らず、性別・race 由来の区分・基底口調セルだけで決める（フォールバック・保留でも付く）。
// 一致が無い、または一人称も言い回しも空なら空文字を返す（役割語マーカーなし）。
func roleSpeechLine(in model.LinePersonaInput, roles *rolespeech.Table) string {
	cell := tone.CellName(in.AttitudeBand, in.EmotionBand)
	tmpl, ok := roles.Lookup(rolespeech.RoleClassOfRace(in.RaceEDID), strings.ToLower(in.Sex), cell)
	if !ok {
		return ""
	}
	switch {
	case tmpl.FirstPerson != "" && tmpl.Register != "":
		return "- 人称と言い回し: 一人称は「" + tmpl.FirstPerson + "」。" + tmpl.Register
	case tmpl.FirstPerson != "":
		return "- 人称と言い回し: 一人称は「" + tmpl.FirstPerson + "」。"
	case tmpl.Register != "":
		return "- 言い回し: " + tmpl.Register
	default:
		return ""
	}
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
