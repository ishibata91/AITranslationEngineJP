package engine

import (
	"strings"

	"aitranslationenginejp/internal/model"
)

// buildPersonaDirective は話者の口調 traits を口調指示文（翻訳ディレクティブ）へ変換する。
// system_requirements.md §3 のとおり、テンプレートで機械的に組み、AI を使わない。
// traits が 1 つも無ければ空文字を返し、ペルソナ指示を注入しない。
func buildPersonaDirective(p model.SpeakerPersona) string {
	var traits []string
	if p.VoiceNature != "" {
		traits = append(traits, "- 声質: "+p.VoiceNature)
	}
	if p.RaceNature != "" {
		traits = append(traits, "- 種族の気質: "+p.RaceNature)
	}
	if len(p.FactionNatures) > 0 {
		traits = append(traits, "- 所属の気風: "+strings.Join(p.FactionNatures, "、"))
	}
	if len(traits) == 0 {
		return ""
	}
	return "この台詞の話者の人物像:\n" +
		strings.Join(traits, "\n") +
		"\nこの人物像に合う口調と人称で訳すこと。"
}

// buildPersonaLabel は結果一覧の口調チップ用の短い要約を返す。最も口調に効く声質を優先するが、
// 声質が性別だけの一般的な声（固有声などで個性が出ない場合）なら、より具体的な種族の気質を採る。
// それも無ければ一般的な声質、所属の気風の順で採る。traits が無ければ空文字。
func buildPersonaLabel(p model.SpeakerPersona) string {
	switch {
	case p.VoiceNature != "" && !isGenericVoice(p.VoiceNature):
		return "声質: " + p.VoiceNature
	case p.RaceNature != "":
		return "種族: " + p.RaceNature
	case p.VoiceNature != "":
		return "声質: " + p.VoiceNature
	case len(p.FactionNatures) > 0:
		return "所属: " + p.FactionNatures[0]
	default:
		return ""
	}
}

// isGenericVoice は声質が性別だけの一般的な声（個性が出ない）かを判定する。
func isGenericVoice(v string) bool {
	return v == "男性の声" || v == "女性の声"
}
