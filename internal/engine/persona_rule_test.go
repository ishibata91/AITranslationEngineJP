package engine

import (
	"testing"

	"aitranslationenginejp/internal/model"
)

// 声型 EditorID から声質の気質を引くこと。完全一致表と、命名規約からの推定の両方を確かめる。
func TestVoiceNature(t *testing.T) {
	cases := []struct {
		edid string
		want string
	}{
		{"MaleChild", "幼い少年の声"},         // 推定（Child + Male）
		{"FemaleOldGrumpy", "気難しい老女の声"}, // 推定（Old + Grumpy + Female）
		{"FemaleChildEager", "幼い少女の声"},  // 推定（Child + Female）
		{"MaleOldGrumpyA", "気難しい老人の声"},  // 推定（Old + Grumpy + Male）
		{"FemaleEvenToned", "女性の声"},     // 推定（性別のみ）
		{"MaleEvenToned", "男性の声"},       // 推定（性別のみ）
		{"", ""},                        // 声型なし
	}
	for _, c := range cases {
		if got := voiceNature(c.edid); got != c.want {
			t.Errorf("voiceNature(%q) = %q, want %q", c.edid, got, c.want)
		}
	}
}

// 識別子から口調 traits を組み、未知の識別子は畳むこと。
func TestPersonaFromIdentity(t *testing.T) {
	got := personaFromIdentity(model.SpeakerIdentity{
		RaceEDID:     "NordRace",
		VoiceEDID:    "MaleChild",
		FactionEDIDs: []string{"DarkBrotherhood", "UnknownFaction"},
	})
	if got.VoiceNature != "幼い少年の声" {
		t.Errorf("VoiceNature = %q", got.VoiceNature)
	}
	if got.RaceNature == "" {
		t.Errorf("既知の種族 NordRace の気質が空")
	}
	if len(got.FactionNatures) != 1 {
		t.Errorf("既知の勢力 1 件だけ残るはず: %v", got.FactionNatures)
	}
}

// 未知の識別子だけなら口調 traits が空になり、ペルソナ指示が出ないこと。
func TestPersonaFromIdentityUnknown(t *testing.T) {
	got := personaFromIdentity(model.SpeakerIdentity{RaceEDID: "MadeUpRace", VoiceEDID: ""})
	if got.RaceNature != "" || got.VoiceNature != "" || len(got.FactionNatures) != 0 {
		t.Errorf("未知の識別子で traits が残った: %+v", got)
	}
	if buildPersonaDirective(testPersonaTemplate, got) != "" {
		t.Errorf("空 traits で directive が出た")
	}
}
