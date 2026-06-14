package engine

import (
	"strings"

	"aitranslationenginejp/internal/model"
)

// personaFromIdentity は話者の事実上の識別子（EditorID）を口調 traits へ写す最小ルール。
// T2 のルール 1 系統で、声型・種族・所属勢力の EditorID から機械的に気質を引く。AI を使わない。
// ルールの永続化と編集 UI は後続 task（T4）で扱う。未知の識別子は空にして畳む。
func personaFromIdentity(id model.SpeakerIdentity) model.SpeakerPersona {
	persona := model.SpeakerPersona{
		VoiceNature: voiceNature(id.VoiceEDID),
		RaceNature:  raceNature(id.RaceEDID),
	}
	for _, edid := range id.FactionEDIDs {
		if n := factionNatureByEDID[edid]; n != "" {
			persona.FactionNatures = append(persona.FactionNatures, n)
		}
	}
	return persona
}

// voiceNature は声型 EditorID（VTYP）から声質の気質を引く。完全一致の表に無くても、
// Skyrim の VTYP 命名（性別＋年齢＋気性）から機械的に推定する。空 EditorID は空を返す。
func voiceNature(edid string) string {
	if edid == "" {
		return ""
	}
	if n, ok := voiceNatureByEDID[edid]; ok {
		return n
	}
	female := strings.HasPrefix(edid, "Female")
	gender := "男性"
	if female {
		gender = "女性"
	}
	switch {
	case strings.Contains(edid, "Child"):
		if female {
			return "幼い少女の声"
		}
		return "幼い少年の声"
	case strings.Contains(edid, "Old") && strings.Contains(edid, "Grumpy"):
		if female {
			return "気難しい老女の声"
		}
		return "気難しい老人の声"
	case strings.Contains(edid, "Old") && strings.Contains(edid, "Kindly"):
		if female {
			return "穏やかな老女の声"
		}
		return "好々爺の声"
	case strings.Contains(edid, "Old"):
		return "年老いた" + gender + "の声"
	case strings.Contains(edid, "Young") || strings.Contains(edid, "Eager"):
		return "若々しい" + gender + "の声"
	case strings.Contains(edid, "Drunk"):
		return "酔った" + gender + "の声"
	case strings.Contains(edid, "Coward"):
		return "気弱な" + gender + "の声"
	case strings.Contains(edid, "Haughty") || strings.Contains(edid, "Condescending"):
		return "高慢な" + gender + "の声"
	case strings.Contains(edid, "Sultry"):
		return "艶のある女性の声"
	default:
		return gender + "の声"
	}
}

// voiceNatureByEDID は特に効かせたい代表 VTYP の完全一致表。命名推定より優先する。
var voiceNatureByEDID = map[string]string{
	"MaleChild":       "幼い少年の声",
	"FemaleChild":     "幼い少女の声",
	"FemaleOldGrumpy": "気難しい老女の声",
	"FemaleOldKindly": "穏やかな老女の声",
	"MaleOldGrumpy":   "気難しい老人の声",
	"MaleOldKindly":   "好々爺の声",
}

// raceNature は種族 EditorID（RACE）から種族の気質を引く。完全一致の表に無くても、
// 子供種族（"…RaceChild"）は基底種族の気質に子供を加えて推定する。未知は空。
func raceNature(edid string) string {
	if edid == "" {
		return ""
	}
	if n, ok := raceNatureByEDID[edid]; ok {
		return n
	}
	if base, ok := strings.CutSuffix(edid, "Child"); ok {
		if n, ok := raceNatureByEDID[base]; ok {
			if i := strings.Index(n, "（"); i >= 0 {
				return n[:i] + "の子供" + n[i:]
			}
			return n + "の子供"
		}
		return "子供"
	}
	return ""
}

// raceNatureByEDID は種族 EditorID（RACE）から種族の気質を引く最小表。未登録は空。
var raceNatureByEDID = map[string]string{
	"NordRace":     "ノルド（実直で粘り強い北方の気質）",
	"BretonRace":   "ブレトン（如才ない交渉上手の気質）",
	"ImperialRace": "インペリアル（規律を重んじる帝国の気質）",
	"RedguardRace": "レッドガード（誇り高い砂漠の戦士の気質）",
	"DarkElfRace":  "ダンマー（用心深い灰の民の気質）",
	"HighElfRace":  "アルトマー（高慢な魔法の民の気質）",
	"WoodElfRace":  "ボズマー（森に生きる狩人の気質）",
	"OrcRace":      "オルシマー（武骨な氏族の気質）",
	"ArgonianRace": "アルゴニアン（沼地の異邦人の気質）",
	"KhajiitRace":  "カジート（抜け目ない隊商の気質）",
}

// factionNatureByEDID は所属勢力 EditorID（FACT）から気風を引く最小表。未登録は畳む。
var factionNatureByEDID = map[string]string{
	"CompanionsFaction":          "戦士団の気風（武勇と仲間意識）",
	"CollegeofWinterholdFaction": "魔法学院の気風（探究と理知）",
	"ThievesGuildFaction":        "盗賊ギルドの気風（抜け目なさ）",
	"DarkBrotherhood":            "闇の一党の気風（冷徹さ）",
}
