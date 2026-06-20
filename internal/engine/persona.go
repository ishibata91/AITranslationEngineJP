package engine

import (
	"strings"

	"aitranslationenginejp/internal/model"
)

// traitsToken は口調指示テンプレートで話者の性質列を差し込む位置を表す差し込み口。
const traitsToken = "{traits}"

// personaTraits は話者の口調 traits を箇条書き行へ写す。声質・種族の気質・所属の気風の順に並べる。
// 既知属性が 1 つも無ければ空 slice を返す（呼び出し側は口調指示なしとして扱う）。
func personaTraits(p model.SpeakerPersona) []string {
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
	return traits
}

// buildPersonaDirective は口調指示テンプレートの {traits} へ話者の性質列を差し込み、口調指示文を組む。
// system_requirements.md §3 のとおり、テンプレートで機械的に組み、AI を使わない。
// テンプレートは prompt_template から読んだ編集可能な雛形で、性質列の中身（属性 → 性質文）はハードコードのまま使う。
// traits が 1 つも無ければ空文字を返し、ペルソナ指示を注入しない。
func buildPersonaDirective(personaTemplate string, p model.SpeakerPersona) string {
	traits := personaTraits(p)
	if len(traits) == 0 {
		return ""
	}
	return strings.ReplaceAll(personaTemplate, traitsToken, strings.Join(traits, "\n"))
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
