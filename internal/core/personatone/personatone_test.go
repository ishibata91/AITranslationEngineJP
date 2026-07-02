package personatone

import (
	"strings"
	"testing"

	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/core/tone"
	"aitranslationenginejp/internal/model"
)

// toneTraitOf は 9 セル全てで性質文を持ち、範囲外は空を返す。
func TestToneTraitOf(t *testing.T) {
	for emo := range 3 {
		for att := range 3 {
			if toneTraitOf(att, emo) == "" {
				t.Errorf("基底口調セル (対人%d, 感情%d) の性質文が空", att, emo)
			}
		}
	}
	if toneTraitOf(tone.AttitudeNeutral, 3) != "" {
		t.Error("感情段階が範囲外で性質文が空でない")
	}
	if toneTraitOf(3, tone.EmotionMid) != "" {
		t.Error("対人段階が範囲外で性質文が空でない")
	}
}

// raceMarkerTrait は種族訛りを持つ種族にだけ注記を返す。
func TestRaceMarkerTrait(t *testing.T) {
	if m := raceMarkerTrait("KhajiitRace"); !strings.Contains(m, "三人称") {
		t.Errorf("Khajiit の種族訛り = %q, 三人称の注記を期待", m)
	}
	if raceMarkerTrait("ArgonianRace") == "" {
		t.Error("Argonian の種族訛りが空")
	}
	if raceMarkerTrait("NordRace") != "" {
		t.Error("種族訛りを持たない Nord で注記が空でない")
	}
}

// buildToneTraits は性質文に種族訛りを重ね、範囲外の段階では空を返す（役割語テンプレート無し）。
func TestBuildToneTraits(t *testing.T) {
	// 丁寧×中（物腰やわ）＋ Khajiit → 口調行 ＋ 種族訛り行。役割語テンプレート nil。
	traits := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudePolite, EmotionBand: tone.EmotionMid, RaceEDID: "KhajiitRace",
	}, nil)
	if len(traits) != 2 {
		t.Fatalf("Khajiit の口調指示行 = %d, want 2", len(traits))
	}
	if !strings.HasPrefix(traits[0], "- 口調: ") || !strings.HasPrefix(traits[1], "- 種族訛り: ") {
		t.Errorf("行の組み立てが想定外: %v", traits)
	}

	// 種族訛りを持たない種族は口調行のみ。
	plain := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudeArrogant, EmotionBand: tone.EmotionSuppressed, RaceEDID: "NordRace",
	}, nil)
	if len(plain) != 1 {
		t.Fatalf("Nord の口調指示行 = %d, want 1", len(plain))
	}

	// 段階が範囲外なら口調指示なし（空）。
	if got := BuildToneTraits(model.LinePersonaInput{AttitudeBand: 9, EmotionBand: 9}, nil); got != nil {
		t.Errorf("範囲外段階で traits = %v, want nil", got)
	}
}

// buildToneTraits は役割語テンプレートがあれば、性質文と種族訛りの間へ一人称・言い回しの行を挟む。
func TestBuildToneTraitsWithRoleSpeech(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader(
		"elder\tfemale\t*\tわたし\t年配の女性らしく落ち着いて。\n"))
	if err != nil {
		t.Fatalf("ParseRoleSpeech: %v", err)
	}
	// 老女（ElderRace + Female）・物腰やわ（丁寧×中）。口調 ＋ 人称と言い回し の 2 行（種族訛り無し）。
	traits := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudePolite, EmotionBand: tone.EmotionMid,
		Sex: "Female", RaceEDID: "ElderRace",
	}, roles)
	if len(traits) != 2 {
		t.Fatalf("老女の口調指示行 = %d, want 2: %v", len(traits), traits)
	}
	if !strings.HasPrefix(traits[1], "- 人称と言い回し: ") || !strings.Contains(traits[1], "わたし") {
		t.Errorf("役割語の行が想定外: %q", traits[1])
	}

	// テンプレートに当たらない成人男（adult/male、行なし）は役割語を付けず口調行のみ。
	none := BuildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudeNeutral, EmotionBand: tone.EmotionMid,
		Sex: "Male", RaceEDID: "NordRace",
	}, roles)
	if len(none) != 1 {
		t.Fatalf("成人男の口調指示行 = %d, want 1: %v", len(none), none)
	}
}

// buildFreeToneTraits は自由記述の口調へ、感情段階の助言と性別の一人称・語尾を重ねる（汎用・PC 用）。
func TestBuildFreeToneTraits(t *testing.T) {
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader(
		"adult\tfemale\t*\tわたし\t女性らしいやわらかな言い回し。\n"))
	if err != nil {
		t.Fatalf("ParseRoleSpeech: %v", err)
	}

	// 自由記述が空なら口調指示なし（空）。
	if got := BuildFreeToneTraits("", tone.EmotionMid, "Female", roles); got != nil {
		t.Errorf("自由記述が空で traits = %v, want nil", got)
	}

	// 感情が中（1）は助言を出さない。女性は一人称・語尾が付く。口調 ＋ 人称 の 2 行。
	mid := BuildFreeToneTraits("衛兵の汎用台詞。", tone.EmotionMid, "Female", roles)
	if len(mid) != 2 {
		t.Fatalf("中・女性の口調指示行 = %d, want 2: %v", len(mid), mid)
	}
	if !strings.HasPrefix(mid[0], "- 口調: ") || !strings.Contains(mid[1], "わたし") {
		t.Errorf("行の組み立てが想定外: %v", mid)
	}

	// 感情が激情（2）は助言を出す。性別なしは一人称・語尾を付けない。口調 ＋ 感情 の 2 行。
	intense := BuildFreeToneTraits("汎用台詞。", tone.EmotionIntense, "", roles)
	if len(intense) != 2 {
		t.Fatalf("激情・性別なしの口調指示行 = %d, want 2: %v", len(intense), intense)
	}
	if !strings.HasPrefix(intense[1], "- 感情: ") {
		t.Errorf("感情助言の行が想定外: %v", intense)
	}

	// 抑制（0）＋ 性別なし ＋ roles nil は口調 ＋ 感情 の 2 行（一人称・語尾なし）。
	suppressed := BuildFreeToneTraits("PC の選択肢。", tone.EmotionSuppressed, "Male", nil)
	if len(suppressed) != 2 {
		t.Fatalf("抑制の口調指示行 = %d, want 2: %v", len(suppressed), suppressed)
	}

	// テンプレートに当たらない成人男（adult/male、行なし）は一人称・語尾を付けず、中なら口調行のみ。
	maleMid := BuildFreeToneTraits("汎用台詞。", tone.EmotionMid, "Male", roles)
	if len(maleMid) != 1 {
		t.Fatalf("中・成人男の口調指示行 = %d, want 1: %v", len(maleMid), maleMid)
	}
}

// freeRoleSpeechLine（formatRoleSpeech）は一人称のみ・言い回しのみ・両方空を出し分ける。
func TestFreeRoleSpeechFormatting(t *testing.T) {
	firstOnly, _ := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tmale\t*\t俺\t\n"))
	if got := freeRoleSpeechLine("Male", firstOnly); got != "- 人称と言い回し: 一人称は「俺」。" {
		t.Errorf("一人称のみ = %q", got)
	}
	registerOnly, _ := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tmale\t*\t\t男性らしい言い回し。\n"))
	if got := freeRoleSpeechLine("Male", registerOnly); got != "- 言い回し: 男性らしい言い回し。" {
		t.Errorf("言い回しのみ = %q", got)
	}
	bothEmpty, _ := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tmale\t*\t\t\n"))
	if got := freeRoleSpeechLine("Male", bothEmpty); got != "" {
		t.Errorf("一人称も言い回しも空で = %q, want 空", got)
	}
}

// buildToneDirective は traits があればテンプレートの {traits} を置換し、無ければ空を返す。
func TestBuildToneDirective(t *testing.T) {
	tmpl := "口調指示:\n{traits}"
	got := BuildToneDirective(tmpl, []string{"- 口調: 柔らかく丁寧"})
	if !strings.Contains(got, "- 口調: 柔らかく丁寧") || strings.Contains(got, "{traits}") {
		t.Errorf("{traits} の置換が想定外: %q", got)
	}
	if BuildToneDirective(tmpl, nil) != "" {
		t.Error("traits が空で空文字を返さない")
	}
}

// buildToneLabel は基底口調セル名のチップ要約を返し、範囲外は空。
func TestBuildToneLabel(t *testing.T) {
	got := BuildToneLabel(model.LinePersonaInput{AttitudeBand: tone.AttitudePolite, EmotionBand: tone.EmotionMid})
	if got != "口調: 物腰やわ" {
		t.Errorf("buildToneLabel = %q, want %q", got, "口調: 物腰やわ")
	}
	if BuildToneLabel(model.LinePersonaInput{AttitudeBand: 9, EmotionBand: 9}) != "" {
		t.Error("範囲外段階で空文字を返さない")
	}
}
