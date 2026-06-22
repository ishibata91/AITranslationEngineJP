package engine

import (
	"strings"
	"testing"

	"aitranslationenginejp/internal/engine/tone"
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

// buildToneTraits は性質文に種族訛りを重ね、範囲外の段階では空を返す。
func TestBuildToneTraits(t *testing.T) {
	// 丁寧×中（物腰やわ）＋ Khajiit → 口調行 ＋ 種族訛り行。
	traits := buildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudePolite, EmotionBand: tone.EmotionMid, RaceEDID: "KhajiitRace",
	})
	if len(traits) != 2 {
		t.Fatalf("Khajiit の口調指示行 = %d, want 2", len(traits))
	}
	if !strings.HasPrefix(traits[0], "- 口調: ") || !strings.HasPrefix(traits[1], "- 種族訛り: ") {
		t.Errorf("行の組み立てが想定外: %v", traits)
	}

	// 種族訛りを持たない種族は口調行のみ。
	plain := buildToneTraits(model.LinePersonaInput{
		AttitudeBand: tone.AttitudeArrogant, EmotionBand: tone.EmotionSuppressed, RaceEDID: "NordRace",
	})
	if len(plain) != 1 {
		t.Fatalf("Nord の口調指示行 = %d, want 1", len(plain))
	}

	// 段階が範囲外なら口調指示なし（空）。
	if got := buildToneTraits(model.LinePersonaInput{AttitudeBand: 9, EmotionBand: 9}); got != nil {
		t.Errorf("範囲外段階で traits = %v, want nil", got)
	}
}

// buildToneDirective は traits があればテンプレートの {traits} を置換し、無ければ空を返す。
func TestBuildToneDirective(t *testing.T) {
	tmpl := "口調指示:\n{traits}"
	got := buildToneDirective(tmpl, []string{"- 口調: 柔らかく丁寧"})
	if !strings.Contains(got, "- 口調: 柔らかく丁寧") || strings.Contains(got, "{traits}") {
		t.Errorf("{traits} の置換が想定外: %q", got)
	}
	if buildToneDirective(tmpl, nil) != "" {
		t.Error("traits が空で空文字を返さない")
	}
}

// buildToneLabel は基底口調セル名のチップ要約を返し、範囲外は空。
func TestBuildToneLabel(t *testing.T) {
	got := buildToneLabel(model.LinePersonaInput{AttitudeBand: tone.AttitudePolite, EmotionBand: tone.EmotionMid})
	if got != "口調: 物腰やわ" {
		t.Errorf("buildToneLabel = %q, want %q", got, "口調: 物腰やわ")
	}
	if buildToneLabel(model.LinePersonaInput{AttitudeBand: 9, EmotionBand: 9}) != "" {
		t.Error("範囲外段階で空文字を返さない")
	}
}
