package engine

import (
	"strings"

	"aitranslationenginejp/internal/engine/tone"
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

// buildToneTraits は注入入力（生成済み基底口調＋種族）から口調指示の箇条書き行を組む。
// 性質文（基底口調）と種族訛り（語彙マーカー）の順に並べる。性質文が無ければ空（口調指示なし）。
func buildToneTraits(in model.LinePersonaInput) []string {
	trait := toneTraitOf(in.AttitudeBand, in.EmotionBand)
	if trait == "" {
		return nil
	}
	traits := []string{"- 口調: " + trait}
	if m := raceMarkerTrait(in.RaceEDID); m != "" {
		traits = append(traits, "- 種族訛り: "+m)
	}
	return traits
}

// buildToneDirective は口調指示テンプレートの {traits} へ口調指示の箇条書きを差し込む。
// traits が空なら空文字を返し、ペルソナ指示を注入しない（既存の差し込み口 traitsToken を再利用）。
func buildToneDirective(personaTemplate string, traits []string) string {
	if len(traits) == 0 {
		return ""
	}
	return strings.ReplaceAll(personaTemplate, traitsToken, strings.Join(traits, "\n"))
}

// buildToneLabel は結果一覧の口調チップ用の短い要約。基底口調セル名を返す。段階が範囲外なら空文字。
func buildToneLabel(in model.LinePersonaInput) string {
	cell := tone.CellName(in.AttitudeBand, in.EmotionBand)
	if cell == "" {
		return ""
	}
	return "口調: " + cell
}
