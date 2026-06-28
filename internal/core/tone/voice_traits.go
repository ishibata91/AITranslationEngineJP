package tone

import "strings"

// voiceTrait は汎用 voice_type の気質辞書 1 件。token は voice EditorID の気質部分（Male/Female を除く）、
// prior は対人軸の事前値（負＝尊大、正＝丁寧）、label は表示用の短い気質名。
// 出典は voice_type の EditorID 命名（制作側が付けた気質ラベル）で、本文と独立した源。
type voiceTrait struct {
	token string
	prior float64
	label string
}

// defaultVoiceTraits は PoC で較正した気質辞書（poc-tone-report.md「メタと本文の融合」）。
// voicePrior・voiceLabel は先頭から Contains で照合するため、並びは特定一致が先に来る順を保つ。
var defaultVoiceTraits = []voiceTrait{
	{"Condescending", -0.6, "横柄"},
	{"ElfHaughty", -0.5, "高慢"},
	{"Brute", -0.5, "粗暴"},
	{"Bandit", -0.5, "粗暴"},
	{"SlyCynical", -0.4, "皮肉"},
	{"OldGrumpy", -0.3, "気難老"},
	{"Shrill", -0.3, "刺々"},
	{"Commander", -0.1, "指揮"},
	{"Guard", -0.1, "衛兵"},
	{"Coward", 0.4, "臆病"},
	{"OldKindly", 0.4, "温厚老"},
	{"EvenToned", 0.0, "平静"},
	{"Sultry", 0.0, "色気"},
	{"YoungEager", 0.0, "前のめり"},
	{"Drunk", 0.0, "酔"},
	{"Commoner", 0.0, "庶民"},
	{"Soldier", 0.0, "兵"},
	{"Warlock", 0.0, "術士"},
}

// isUniqueVoice は固有・特殊 voice（キャラ専用・creature・特殊）かを判定する。気質 prior を持たない。
// Cr 接頭は creature（馬・犬等）、SPECIAL は特殊枠、空は voice 無しを表す。
func isUniqueVoice(voice string) bool {
	return voice == "" || strings.Contains(voice, "Unique") ||
		strings.HasPrefix(voice, "Cr") || strings.HasPrefix(voice, "SPECIAL")
}

// voicePrior は voice 気質から対人軸の事前値を返す。固有・特殊 voice は prior 無し（ok=false）。
// 汎用だが気質が中立の voice（Nord 等の種族系）は中立 0 を返す（ok=true）。
func (c *Classifier) voicePrior(voice string) (prior float64, ok bool) {
	if isUniqueVoice(voice) {
		return 0, false
	}
	for _, t := range c.voiceTraits {
		if strings.Contains(voice, t.token) {
			return t.prior, true
		}
	}
	return 0, true
}

// voiceLabel は voice 気質の表示名を返す。固有・特殊 voice は「固有」、汎用で気質中立は「—」。
func (c *Classifier) voiceLabel(voice string) string {
	if isUniqueVoice(voice) {
		return "固有"
	}
	for _, t := range c.voiceTraits {
		if strings.Contains(voice, t.token) {
			return t.label
		}
	}
	return "—"
}
